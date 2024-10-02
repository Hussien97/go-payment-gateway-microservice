package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"payment-gateway/db"
	"payment-gateway/internal/models"
	"payment-gateway/internal/redis"
	"testing"
)

func TestMain(m *testing.M) {
	redis.InitRedis()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbURL := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	db.InitializeDB(dbURL)

	code := m.Run()
	os.Exit(code)
}

// Test DepositHandler
func TestValidDepositJSON(t *testing.T) {
	handler := http.HandlerFunc(DepositHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := models.TransactionRequest{
		Amount: 100.0,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	res, err := http.Post(server.URL+"/deposit", "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %v", res.Status)
	}
}

func TestValidDepositSOAP(t *testing.T) {
	handler := http.HandlerFunc(DepositHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := `<?xml version="1.0" encoding="UTF-8"?>
<transaction>
    <amount>100.0</amount>
</transaction>`
	res, err := http.Post(server.URL+"/deposit", "text/xml", bytes.NewBuffer([]byte(reqBody)))
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %v", res.Status)
	}
}

func TestInvalidDepositMissingAmount(t *testing.T) {
	handler := http.HandlerFunc(DepositHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := models.TransactionRequest{}
	reqBodyBytes, _ := json.Marshal(reqBody)

	res, err := http.Post(server.URL+"/deposit", "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request, got %v", res.Status)
	}
}

func TestInvalidDepositMissingAmountSOAP(t *testing.T) {
	handler := http.HandlerFunc(DepositHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := `<?xml version="1.0" encoding="UTF-8"?>
<transaction>
</transaction>`
	res, err := http.Post(server.URL+"/deposit", "text/xml", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request, got %v", res.Status)
	}
}

// Test WithdrawalHandler
func TestValidWithdrawalJSON(t *testing.T) {
	handler := http.HandlerFunc(WithdrawalHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := models.TransactionRequest{
		Amount: 50.0,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	res, err := http.Post(server.URL+"/withdrawal", "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %v", res.Status)
	}
}

func TestValidWithdrawalSOAP(t *testing.T) {
	handler := http.HandlerFunc(WithdrawalHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := `<?xml version="1.0" encoding="UTF-8"?>
<transaction>
    <amount>50.0</amount>
</transaction>`
	res, err := http.Post(server.URL+"/withdrawal", "text/xml", bytes.NewBuffer([]byte(reqBody)))
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %v", res.Status)
	}
}

func TestInvalidWithdrawalMissingAmount(t *testing.T) {
	handler := http.HandlerFunc(WithdrawalHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := models.TransactionRequest{}
	reqBodyBytes, _ := json.Marshal(reqBody)

	res, err := http.Post(server.URL+"/withdrawal", "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request, got %v", res.Status)
	}
}

func TestInvalidWithdrawalMissingAmountSOAP(t *testing.T) {
	handler := http.HandlerFunc(WithdrawalHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := `<?xml version="1.0" encoding="UTF-8"?>
<transaction>
</transaction>`
	res, err := http.Post(server.URL+"/withdrawal", "text/xml", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request, got %v", res.Status)
	}
}

// Test CallbackHandler
func TestValidCallbackJSON(t *testing.T) {
	handler := http.HandlerFunc(CallbackHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	redis.SetTransactionStatus("trans1", "pending")

	reqBody := models.TransactionRequest{
		TransactionID: "trans1",
		Status:        "completed",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	res, err := http.Post(server.URL+"/callback", "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %v", res.Status)
	}
}

func TestValidCallbackSOAP(t *testing.T) {
	handler := http.HandlerFunc(CallbackHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	redis.SetTransactionStatus("trans1", "pending")

	reqBody := `<?xml version="1.0" encoding="UTF-8"?>
<transaction>
    <transaction_id>trans1</transaction_id>
    <status>completed</status>
</transaction>`
	res, err := http.Post(server.URL+"/callback", "text/xml", bytes.NewBuffer([]byte(reqBody)))
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %v", res.Status)
	}
}

func TestInvalidCallbackMissingTransactionID(t *testing.T) {
	handler := http.HandlerFunc(CallbackHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := models.TransactionRequest{
		Status: "completed",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	res, err := http.Post(server.URL+"/callback", "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request, got %v", res.Status)
	}
}

func TestInvalidCallbackMissingTransactionIDSOAP(t *testing.T) {
	handler := http.HandlerFunc(CallbackHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := `<?xml version="1.0" encoding="UTF-8"?>
<transaction>
    <status>completed</status>
</transaction>`
	res, err := http.Post(server.URL+"/callback", "text/xml", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request, got %v", res.Status)
	}
}

func TestInvalidCallbackTransactionNotFound(t *testing.T) {
	handler := http.HandlerFunc(CallbackHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := models.TransactionRequest{
		TransactionID: "nonexistent-transactions",
		Status:        "completed",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	res, err := http.Post(server.URL+"/callback", "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request for nonexistent transaction, got %v", res.Status)
	}
}

func TestInvalidCallbackTransactionNotFoundSOAP(t *testing.T) {
	handler := http.HandlerFunc(CallbackHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := `<?xml version="1.0" encoding="UTF-8"?>
<transaction>
    <transaction_id>nonexistent-transactionss</transaction_id>
    <status>completed</status>
</transaction>`
	res, err := http.Post(server.URL+"/callback", "text/xml", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request for nonexistent transaction, got %v", res.Status)
	}
}

func TestInvalidCallbackMissingStatus(t *testing.T) {
	handler := http.HandlerFunc(CallbackHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	redis.SetTransactionStatus("trans1", "pending")

	reqBody := models.TransactionRequest{
		TransactionID: "trans1",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	res, err := http.Post(server.URL+"/callback", "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request for missing status, got %v", res.Status)
	}
}

func TestInvalidCallbackMissingStatusSOAP(t *testing.T) {
	handler := http.HandlerFunc(CallbackHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	redis.SetTransactionStatus("trans1", "pending")

	reqBody := `<?xml version="1.0" encoding="UTF-8"?>
<transaction>
    <transaction_id>trans1</transaction_id>
</transaction>`
	res, err := http.Post(server.URL+"/callback", "text/xml", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request for missing status, got %v", res.Status)
	}
}

// Test Invalid Content Type for All Handlers
func TestInvalidContentTypeDeposit(t *testing.T) {
	handler := http.HandlerFunc(DepositHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := models.TransactionRequest{
		Amount: 100.0,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	res, err := http.Post(server.URL+"/deposit", "application/xml", bytes.NewBuffer(reqBodyBytes))
	if err != nil || res.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("Expected status 415 Unsupported Media Type, got %v", res.Status)
	}
}

func TestInvalidContentTypeWithdrawal(t *testing.T) {
	handler := http.HandlerFunc(WithdrawalHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := models.TransactionRequest{
		Amount: 50.0,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	res, err := http.Post(server.URL+"/withdrawal", "application/xml", bytes.NewBuffer(reqBodyBytes))
	if err != nil || res.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("Expected status 415 Unsupported Media Type, got %v", res.Status)
	}
}

func TestInvalidContentTypeCallback(t *testing.T) {
	handler := http.HandlerFunc(CallbackHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	reqBody := models.TransactionRequest{
		TransactionID: "trans1",
		Status:        "completed",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	res, err := http.Post(server.URL+"/callback", "application/xml", bytes.NewBuffer(reqBodyBytes))
	if err != nil || res.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("Expected status 415 Unsupported Media Type, got %v", res.Status)
	}
}
