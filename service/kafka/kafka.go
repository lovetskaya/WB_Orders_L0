package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"wb_service_order/service/cache"
	"wb_service_order/service/order"
)

var logger *zap.Logger

// Инициализация логгера
func InitLogger(l *zap.Logger) {
	logger = l
}

func CreateTopic(kafkaBrokers []string) {
	conn, err := kafka.Dial("tcp", kafkaBrokers[0])
	if err != nil {
		logger.Fatal("Ошибка при подключении к Kafka", zap.Error(err))
	}
	defer conn.Close()

	topic := "orders"
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	controller, err := conn.Controller()
	if err != nil {
		logger.Fatal("Ошибка при получении контроллера", zap.Error(err))
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		logger.Fatal("Ошибка при подключении к контроллеру", zap.Error(err))
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(topicConfig)
	if err != nil {
		logger.Fatal("Ошибка при создании топика", zap.Error(err))
	}

	logger.Info("Топик успешно создан", zap.String("topic", topic))
}

func ConsumeMessages(db *sql.DB, kafkaBrokers []string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     kafkaBrokers,
		GroupID:     "order_service",
		Topic:       "orders",
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})

	defer r.Close()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			logger.Warn("Ошибка при чтении сообщения", zap.Error(err))
			continue
		}

		var order order.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			logger.Warn("Не удалось десериализовать сообщение", zap.Error(err))
			continue
		}

		order.SaveToDB(db)
		cache.OrderCache.Store(order.OrderUID, order)
		logger.Info("Сообщение успешно обработано", zap.String("orderUID", order.OrderUID))
	}
}