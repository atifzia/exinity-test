package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payment-gateway/db"
	"payment-gateway/internal/kafka"
	"payment-gateway/internal/models/common"
	"payment-gateway/internal/models/postgres"
	"payment-gateway/internal/models/request"
	"testing"
)

func TestWithdrawHandler(t *testing.T) {

	// Mock request payload
	withdrawRequest := request.Transaction{
		Amount:    50.0,
		UserID:    1,
		CountryID: 840,
		Currency:  "USD",
	}

	requestBody, err := json.Marshal(withdrawRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	a := API{db: &db.MockDB{
		GetSupportedGatewaysByCountryFunc: func(countryID int) ([]*common.Gateway, error) {
			return []*common.Gateway{
				{ID: 1, Name: "Mock Gateway", Priority: 1},
			}, nil
		},
		CreateTransactionFunc: func(tx *postgres.Transaction) error {
			tx.ID = 12345
			return nil
		},
		UpdateTxStatusFunc: func(txID int64, status string) error {
			return nil
		},
	}}
	a.SetupServices(&kafka.MockKafkaProducer{})
	handler := http.HandlerFunc(a.WithdrawalHandler)
	handler.ServeHTTP(rr, req)

	// Assert response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Assert response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if response["message"] != "tx under processing" {
		t.Errorf("Expected message 'Transaction processed successfully', got %s", response["message"])
	}
}
