package tx

import (
	"errors"
	"payment-gateway/db"
	"payment-gateway/internal/kafka"
	"payment-gateway/internal/models/common"
	"payment-gateway/internal/models/postgres"
	"payment-gateway/internal/models/request"
	"payment-gateway/internal/services/gateway"
	"testing"
)

func TestProcessTransaction_Success(t *testing.T) {
	// Mock database implementation
	mockDB := &db.MockDB{
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
	}

	// Test data
	requestPayload := request.Transaction{
		Amount:    100.0,
		UserID:    1,
		GatewayID: 101,
		CountryID: 840,
		Currency:  "USD",
	}

	// Call the function
	response, err := NewSvcTx(mockDB, &kafka.MockKafkaProducer{}).ProcessTransaction(requestPayload, gateway.NewSvcGateway(mockDB), "deposit")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	// Assertions
	if response.StatusCode != 200 {
		t.Errorf("expected status code 200, got %d", response.StatusCode)
	}
	if response.Message != "tx under processing" {
		t.Errorf("unexpected message: %s", response.Message)
	}
	if response.Data["transaction_id"] != int64(12345) {
		t.Errorf("unexpected transaction ID: %v", response.Data["transaction_id"])
	}
}

// TestProcessTransaction_Failure tests scenarios where ProcessTransaction fails
func TestProcessTransaction_Failure(t *testing.T) {
	// Case 1: No available gateways
	t.Run("NoAvailableGateways", func(t *testing.T) {
		// Mock database
		mockDB := &db.MockDB{
			GetSupportedGatewaysByCountryFunc: func(countryID int) ([]*common.Gateway, error) {
				return nil, errors.New("no gateways available")
			},
			CreateTransactionFunc: func(tx *postgres.Transaction) error { return nil },
			UpdateTxStatusFunc:    func(txID int64, status string) error { return nil },
		}

		requestPayload := request.Transaction{
			Amount:    100.0,
			UserID:    1,
			CountryID: 840,
			Currency:  "USD",
		}

		_, err := NewSvcTx(mockDB, &kafka.MockKafkaProducer{}).ProcessTransaction(requestPayload, gateway.NewSvcGateway(mockDB), "deposit")
		if err == nil || err.Error() != "no gateways available" {
			t.Errorf("expected error 'no gateways available', got %v", err)
		}
	})

	// Case 2: Database failure during transaction creation
	t.Run("CreateTransactionFailure", func(t *testing.T) {
		// Mock database
		mockDB := &db.MockDB{
			GetSupportedGatewaysByCountryFunc: func(countryID int) ([]*common.Gateway, error) {
				return []*common.Gateway{{ID: 1, Name: "Mock Gateway", Priority: 1}}, nil
			},
			CreateTransactionFunc: func(tx *postgres.Transaction) error {
				return errors.New("failed to save tx to database")
			},
			UpdateTxStatusFunc: func(txID int64, status string) error { return nil },
		}

		requestPayload := request.Transaction{
			Amount:    100.0,
			UserID:    1,
			CountryID: 840,
			Currency:  "USD",
		}

		_, err := NewSvcTx(mockDB, &kafka.MockKafkaProducer{}).ProcessTransaction(requestPayload, gateway.NewSvcGateway(mockDB), "deposit")
		if err == nil || err.Error() != "failed to save tx to database" {
			t.Errorf("expected error 'failed to save tx to database', got %v", err)
		}
	})

	// Case 3: Gateway failure during transaction processing
	t.Run("GatewayFailure", func(t *testing.T) {
		// Mock database
		mockDB := &db.MockDB{
			GetSupportedGatewaysByCountryFunc: func(countryID int) ([]*common.Gateway, error) {
				return []*common.Gateway{{ID: 1, Name: "Mock Gateway", Priority: 1}}, nil
			},
			CreateTransactionFunc: func(tx *postgres.Transaction) error {
				tx.ID = 12345
				return nil
			},
			UpdateTxStatusFunc: func(txID int64, status string) error {
				return errors.New("failed to update transaction status")
			},
		}

		// Mock gateway processing
		mockGatewayProcessor := &gateway.MockGatewayProcessor{
			SendTxToGatewayFunc: func(tx postgres.Transaction) (interface{}, error) {
				return nil, errors.New("gateway error")
			},
			SelectGatewayFunc: func(countryID int) (*common.Gateway, error) {
				return nil, errors.New("gateway error")
			},
		}

		requestPayload := request.Transaction{
			Amount:    100.0,
			UserID:    1,
			CountryID: 840,
			Currency:  "USD",
		}

		_, err := NewSvcTx(mockDB, &kafka.MockKafkaProducer{}).ProcessTransaction(requestPayload, mockGatewayProcessor, "deposit")
		if err == nil || err.Error() != "gateway error" {
			t.Errorf("expected error 'gateway error', got %v", err)
		}
	})
}
