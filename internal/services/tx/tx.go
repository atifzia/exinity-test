package tx

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"payment-gateway/db"
	"payment-gateway/internal/kafka"
	"payment-gateway/internal/models/postgres"
	"payment-gateway/internal/models/request"
	"payment-gateway/internal/models/response"
	svcGateway "payment-gateway/internal/services/gateway"
	"payment-gateway/internal/util"
	"time"
)

type (
	SvcTx struct {
		db            db.Idb
		kafkaProducer kafka.IProducer
	}

	ISvcTx interface {
		ProcessTransaction(req request.Transaction, iSvcGateway svcGateway.ISvcGateway, transactionType string) (response.APIResponse, error)
		ProcessCallBack(txId int64, status string) error
	}
)

func NewSvcTx(db db.Idb, kafkaProducer kafka.IProducer) ISvcTx {
	return &SvcTx{db: db, kafkaProducer: kafkaProducer}
}

// ProcessTransaction handles deposit or withdrawal transactions.
func (t SvcTx) ProcessTransaction(req request.Transaction, iSvcGateway svcGateway.ISvcGateway, transactionType string) (response.APIResponse, error) {
	if req.Amount <= 0 {
		return response.APIResponse{}, errors.New("invalid amount, must be greater than zero")
	}

	if req.UserID <= 0 {
		return response.APIResponse{}, errors.New("invalid user_id, must be a positive integer")
	}

	// Step 1: select gateway dynamically based on country_id
	gateway, err := iSvcGateway.SelectGateway(req.CountryID)
	if err != nil {
		return response.APIResponse{}, err
	}

	// Step 2: prepare tx data for processing
	tx := postgres.Transaction{
		UserID:    req.UserID,
		Amount:    req.Amount,
		GatewayID: gateway.ID,
		CountryID: req.CountryID,
		Status:    "pending",
		Type:      transactionType,
	}

	// Step 3: save tx to the database
	err = t.db.CreateTransaction(&tx)
	if err != nil {
		return response.APIResponse{}, errors.New("failed to save tx to database")
	}

	// Step 4: send tx to selected gateway using retry mechanism
	if err = util.RetryOperation(func() error {
		_, err = iSvcGateway.SendTxToGateway(tx)
		return err
	}, 5); err != nil {

		if err = t.db.UpdateTxStatus(tx.ID, "failed"); err != nil {
			return response.APIResponse{}, errors.New("failed to update tx status to db")
		}

		return response.APIResponse{}, errors.New("failed to send tx to gateway")
	}

	// Step 5: publish the tx to Kafka with Circuit Breaker
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	txMsg, err := json.Marshal(tx)
	if err != nil {
		log.Printf("failed to marshal tx for Kafka: %v", err)

		return response.APIResponse{}, errors.New("failed to marshal tx")
	}

	dataFormat := "application/json" // Defaulting to JSON format
	err = t.kafkaProducer.PubTx(ctx, tx.ID, txMsg, dataFormat)
	if err != nil {
		log.Printf("failed to publish tx to Kafka: %v", err)

		return response.APIResponse{}, errors.New("failed to publish tx to Kafka")
	}

	// Step 6: Prepare and return response
	return response.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "tx under processing",
		Data: map[string]interface{}{
			"transaction_id": tx.ID,
			"gateway_id":     tx.GatewayID,
			"status":         tx.Status,
		},
	}, nil
}

func (t SvcTx) ProcessCallBack(txId int64, status string) error {
	if err := t.db.UpdateTxStatus(txId, status); err != nil {
		return errors.New("failed to update transaction status in database")
	}

	return nil
}
