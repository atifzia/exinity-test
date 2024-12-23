package db

import (
	"payment-gateway/internal/models/common"
	"payment-gateway/internal/models/postgres"
)

// MockDB implements the DB interface for testing
type MockDB struct {
	GetSupportedGatewaysByCountryFunc func(countryID int) ([]*common.Gateway, error)
	CreateTransactionFunc             func(tx *postgres.Transaction) error
	UpdateTxStatusFunc                func(txID int64, status string) error
}

func (m *MockDB) GetSupportedGatewaysByCountry(countryID int) ([]*common.Gateway, error) {
	return m.GetSupportedGatewaysByCountryFunc(countryID)
}

func (m *MockDB) CreateTransaction(tx *postgres.Transaction) error {
	return m.CreateTransactionFunc(tx)
}

func (m *MockDB) UpdateTxStatus(txID int64, status string) error {
	return m.UpdateTxStatusFunc(txID, status)
}
