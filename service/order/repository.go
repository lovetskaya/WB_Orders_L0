package order

import (
	"database/sql"
	"encoding/json"
	"go.uber.org/zap"
)

var logger *zap.Logger

// Инициализация логгера
func InitLoggerRep(l *zap.Logger) {
	logger = l
}

func (order *Order) SaveToDB(db *sql.DB) {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		logger.Warn("Ошибка при маршализации элементов заказа", zap.String("orderUID", order.OrderUID), zap.Error(err))
		return
	}

	_, err = db.Exec(`INSERT INTO orders (
		order_uid, track_number, entry, delivery_name, delivery_phone,
		delivery_zip, delivery_city, delivery_address, delivery_region,
		delivery_email, payment_transaction, payment_request_id,
		payment_currency, payment_provider, payment_amount,
		payment_payment_dt, payment_bank, payment_delivery_cost,
		payment_goods_total, payment_custom_fee, items, locale,
		internal_signature, customer_id, delivery_service, shardkey,
		sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5,
		$6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
		$16, $17, $18, $19, $20, $21, $22, $23, $24,
		$25, $26, $27, $28, $29)`,
		order.OrderUID, order.TrackNumber, order.Entry,
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region,
		order.Delivery.Email, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee, itemsJSON,
		order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.ShardKey, order.SMID,
		order.DateCreated, order.OofShard)

	if err != nil {
		logger.Error("Ошибка при вставке заказа в базу данных", zap.String("orderUID", order.OrderUID), zap.Error(err))
		return // или обработка ошибки
	}

	logger.Info("Заказ успешно сохранен в базу данных", zap.String("orderUID", order.OrderUID))
}