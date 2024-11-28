package cache

import (
	"database/sql"
	"encoding/json"
	"go.uber.org/zap"
	"sync"
	"wb_service_order/service/order"
)

var (
	OrderCache sync.Map
	logger     *zap.Logger
)

// Инициализация логгера
func InitLogger(l *zap.Logger) {
	logger = l
}

func RestoreCache(db *sql.DB) {
	rows, err := db.Query(`SELECT order_uid, track_number, entry, delivery_name, delivery_phone, delivery_zip, delivery_city, delivery_address, delivery_region, delivery_email, payment_transaction, payment_request_id, payment_currency, payment_provider, payment_amount, payment_payment_dt, payment_bank, payment_delivery_cost, payment_goods_total, payment_custom_fee, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders`)
	if err != nil {
		logger.Fatal("Ошибка при выполнении запроса к базе данных", zap.Error(err))
	}
	defer rows.Close()

	for rows.Next() {
		var order order.Order
		var itemsJSON string

		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
			&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
			&order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestID,
			&order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount,
			&order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee, &itemsJSON,
			&order.Locale, &order.InternalSignature, &order.CustomerID,
			&order.DeliveryService, &order.ShardKey, &order.SMID, &order.DateCreated, &order.OofShard)

		if err != nil {
			logger.Warn("Ошибка при сканировании строки", zap.Error(err))
			continue
		}

		if err := json.Unmarshal([]byte(itemsJSON), &order.Items); err != nil {
			logger.Warn("Ошибка при десериализации JSON элементов", zap.Error(err))
			continue
		}

		OrderCache.Store(order.OrderUID, order)
		logger.Info("Заказ успешно восстановлен из базы данных", zap.String("orderUID", order.OrderUID))
	}
}