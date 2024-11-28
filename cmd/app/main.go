package main

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"wb_service_order/service/cache"
	"wb_service_order/service/config"
	"wb_service_order/service/kafka"
	"wb_service_order/service/order"
	"wb_service_order/web/handlers"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	_ "go.uber.org/zap/zapcore"
)

var (
	db     *sql.DB
	logger *zap.Logger
)

func init() {
	var err error

	// Создаем папку для логов
	logDir := "logger"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания лог каталога: %v\n", err)
		os.Exit(1)
	}

	// Открываем файл для записи логов
	logFilePath := filepath.Join(logDir, "app.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка открытия лог файла: %v\n", err)
		os.Exit(1)
	}

	// Настройка кодера для записи в файл
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapcore.InfoLevel)

	logger = zap.New(core)
	defer func() {
		if err := logger.Sync(); err != nil {
			logFatal("Ошибка синхронизации логгера", err) // Обработка ошибки синхронизации логгера
		}
	}()

	cfg, err := config.LoadConfig()
	if err != nil {
		logFatal("Ошибка загрузки конфигурации", err)
	}

	db, err = sql.Open("postgres", cfg.DBConnectionString)
	if err != nil {
		logFatal("Ошибка подключения к базе данных", err)
	}

	// Инициализация логгера для пакетов
	kafka.InitLogger(logger)
	cache.InitLogger(logger)
	order.InitLoggerRep(logger)
	order.InitLoggerSer(logger)
	handlers.InitLogger(logger)
	cache.RestoreCache(db)
}

func logFatal(msg string, err error) {
	logger.Fatal(msg, zap.Error(err))
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logFatal("Ошибка загрузки конфигурации", err)
	}

	kafka.CreateTopic(cfg.KafkaBrokers)
	http.HandleFunc("/order", handlers.GetOrder(logger, db))

	srv := &http.Server{
		Addr: cfg.ServerAddress,
	}

	order.LoadOrdersFromJSON("model.json", db)

	go func() {
		logger.Info("Сервер работает", zap.String("address", cfg.ServerAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Ошибка прослушивания", zap.Error(err))
		}
	}()

	go kafka.ConsumeMessages(db, cfg.KafkaBrokers)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Выключение сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Сервер вынужден завершить работу", zap.Error(err))
	}
	logger.Info("Сервер закрывается")
}