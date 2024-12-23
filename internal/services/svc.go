package services

import (
	"payment-gateway/internal/services/gateway"
	"payment-gateway/internal/services/tx"
)

type Service struct {
	gateway.ISvcGateway
	tx.ISvcTx
}
