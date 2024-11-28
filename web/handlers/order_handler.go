package handlers

import (
	"database/sql"
	"encoding/json"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"wb_service_order/service/cache"
	"wb_service_order/service/order"
)
var logger *zap.Logger // Объявляем логгер
func InitLogger(l *zap.Logger) {
	logger = l
}

func GetOrder(logger *zap.Logger, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUID := r.URL.Query().Get("id")
		logger.Info("Получен запрос на заказ", zap.String("orderUID", orderUID))

		// Проверка наличия заказа в кэше
		if res, exists := cache.OrderCache.Load(orderUID); exists {
			logger.Info("Заказ найден в кэше", zap.String("orderUID", orderUID))
			order, ok := res.(order.Order)
			if !ok {
				logger.Warn("Неверный тип данных в кэше для заказа", zap.String("orderUID", orderUID))
				http.Error(w, "Invalid order data", http.StatusInternalServerError)
				return
			}
			renderOrder(logger, w, order)
			return
		}

		// Если заказа нет в кэше, пытаемся получить его из базы данных
		order := getOrderFromDB(orderUID, db)
		if order.OrderUID == "" {
			http.Error(w, "Order not found", http.StatusNotFound)
			logger.Warn("Заказ не найден в базе данных", zap.String("orderUID", orderUID))
			return
		}

		// Сохраняем заказ в кэше
		cache.OrderCache.Store(orderUID, order)
		renderOrder(logger, w, order)
	}
}

func renderOrder(logger *zap.Logger, w http.ResponseWriter, order order.Order) {
	tmpl, err := template.ParseFiles("/Users/polinaloveckaya/GolandProjects/WB_Orders/wb_service_order/web/templates/index.html")
	if err != nil {
		logger.Error("Failed to load template", zap.Error(err))
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, order); err != nil {
		logger.Error("Failed to execute template", zap.Error(err))
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}

func getOrderFromDB(orderUID string, db *sql.DB) order.Order {
	var order order.Order
	var itemsJSON string

	row := db.QueryRow(`SELECT 
		order_uid, track_number, entry, delivery_name, delivery_phone, 
		delivery_zip, delivery_city, delivery_address, delivery_region, 
		delivery_email, payment_transaction, payment_request_id, 
		payment_currency, payment_provider, payment_amount, 
		payment_payment_dt, payment_bank, payment_delivery_cost, 
		payment_goods_total, payment_custom_fee, items, locale, 
		internal_signature, customer_id, delivery_service, shardkey, 
		sm_id, date_created, oof_shard 
	FROM orders WHERE order_uid = $1`, orderUID)

	err := row.Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry,
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestID,
		&order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount,
		&order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal, &order.Payment.CustomFee, &itemsJSON,
		&order.Locale, &order.InternalSignature, &order.CustomerID,
		&order.DeliveryService, &order.ShardKey, &order.SMID,
		&order.DateCreated, &order.OofShard,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("Заказ не найден для", zap.String("orderUID", orderUID))
			return order
		}
		logger.Error("Ошибка при получении заказа из БД", zap.Error(err))
		return order
	}
	if err := json.Unmarshal([]byte(itemsJSON), &order.Items); err != nil {
		logger.Error("Ошибка при разборе JSON элементов", zap.Error(err))
		return order
	}

	return order
}