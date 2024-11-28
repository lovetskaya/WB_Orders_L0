package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
	"wb_service_order/service/order"
	"wb_service_order/web/handlers"
)

func TestGetItem(t *testing.T) {
	// Инициализируем логгер
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("Ошибка при инициализации логгера: %s", err)
	}
	defer logger.Sync() // Отложить синхронизацию

	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock базы данных: %s", err)
	}
	defer db.Close()

	// Создаем тестовый заказ
	cell := order.Order{
		OrderUID:          "b563feb7b2b84b6test",
		TrackNumber:       "track123",
		Entry:             "entry",
		Delivery:          order.Delivery{Name: "John Doe", Phone: "1234567890", Zip: "12345", City: "City", Address: "123 Main St", Region: "Region", Email: "john@example.com"},
		Payment:           order.Payment{Transaction: "trans123", RequestID: "req123", Currency: "USD", Provider: "Visa", Amount: 100.0, PaymentDt: time.Now(), Bank: "Bank", DeliveryCost: 10.0, GoodsTotal: 90.0, CustomFee: 0.0},
		Items:             []order.Item{{ChrtID: 1, Name: "Item 1", Price: 100.0, TotalPrice: 100.0}},
		Locale:            "en",
		InternalSignature: "signature",
		CustomerID:       "cust123",
		DeliveryService:   "service",
		ShardKey:         "shardkey",
		SMID:             1,
		DateCreated:      time.Now(),
		OofShard:         "oof_shard",
	}

	// Ожидаем, что будет выполнен SQL-запрос на получение заказа
	mock.ExpectQuery(`SELECT order_uid, track_number, entry, delivery_name`).
		WithArgs(cell.OrderUID).
		WillReturnRows(sqlmock.NewRows([]string{
			"order_uid", "track_number", "entry", "delivery_name", "delivery_phone",
			"delivery_zip", "delivery_city", "delivery_address", "delivery_region",
			"delivery_email", "payment_transaction", "payment_request_id",
			"payment_currency", "payment_provider", "payment_amount",
			"payment_payment_dt", "payment_bank", "payment_delivery_cost",
			"payment_goods_total", "payment_custom_fee", "items", "locale",
			"internal_signature", "customer_id", "delivery_service", "shardkey",
			"sm_id", "date_created", "oof_shard",
		}).AddRow(
		cell.OrderUID, cell.TrackNumber, cell.Entry, cell.Delivery.Name, cell.Delivery.Phone,
		cell.Delivery.Zip, cell.Delivery.City, cell.Delivery.Address, cell.Delivery.Region,
		cell.Delivery.Email, cell.Payment.Transaction, cell.Payment.RequestID,
		cell.Payment.Currency, cell.Payment.Provider, cell.Payment.Amount,
		cell.Payment.PaymentDt, cell.Payment.Bank, cell.Payment.DeliveryCost,
		cell.Payment.GoodsTotal, cell.Payment.CustomFee, "[]", cell.Locale,
		cell.InternalSignature, cell.CustomerID, cell.DeliveryService, cell.ShardKey,
		cell.SMID, cell.DateCreated, cell.OofShard,
	))


	// Создаем тестовый HTTP-запрос
	req, err := http.NewRequest("GET", "/order?id="+cell.OrderUID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetOrder(logger, db)) // ваш обработчик

	handler.ServeHTTP(rr, req)

	// Проверяем статус ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем содержимое ответа
	// Здесь мы ожидаем, что HTML-ответ будет содержать определенные значения.
	// Для этого мы можем использовать rr.Body.String() для извлечения HTML-контента.
	responseBody := rr.Body.String()

	// Проверяем, что ответ содержит ожидаемые значения
	if !contains(responseBody, cell.OrderUID) {
		t.Errorf("handler returned unexpected OrderUID: got %v want %v", responseBody, cell.OrderUID)
	}
	if !contains(responseBody, cell.TrackNumber) {
		t.Errorf("handler returned unexpected TrackNumber: got %v want %v", responseBody, cell.TrackNumber)
	}
	if !contains(responseBody, cell.Delivery.Name) {
		t.Errorf("handler returned unexpected Delivery Name: got %v want %v", responseBody, cell.Delivery.Name)
	}
	if !contains(responseBody, cell.Payment.Transaction) {
		t.Errorf("handler returned unexpected Payment Transaction: got %v want %v", responseBody, cell.Payment.Transaction)
	}
}


// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}