package gateway

import (
	"payment-gateway/internal/models/common"
	"payment-gateway/internal/models/postgres"
)

// MockGatewayProcessor implements the GatewayProcessor interface for testing.
type MockGatewayProcessor struct {
	SelectGatewayFunc   func(countryID int) (*common.Gateway, error)
	SendTxToGatewayFunc func(tx postgres.Transaction) (interface{}, error)
}

func (m *MockGatewayProcessor) SelectGateway(countryID int) (*common.Gateway, error) {
	return m.SelectGatewayFunc(countryID)
}

func (m *MockGatewayProcessor) SendTxToGateway(tx postgres.Transaction) (interface{}, error) {
	return m.SendTxToGatewayFunc(tx)
}
