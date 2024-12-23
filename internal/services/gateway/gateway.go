package gateway

import (
	"errors"
	"fmt"
	"payment-gateway/db"
	"payment-gateway/internal/models/common"
	"payment-gateway/internal/models/postgres"
)

type (
	SvcGateway struct {
		db db.Idb
	}

	ISvcGateway interface {
		SelectGateway(countryID int) (*common.Gateway, error)
		SendTxToGateway(tx postgres.Transaction) (interface{}, error)
	}
)

func NewSvcGateway(db db.Idb) ISvcGateway {
	return &SvcGateway{db: db}
}

// SelectGateway chooses a payment gateway dynamically according to priority and country.
func (g SvcGateway) SelectGateway(countryID int) (*common.Gateway, error) {
	gateways, err := g.db.GetSupportedGatewaysByCountry(countryID)
	if err != nil {
		fmt.Printf("failed to query gateways: %v\n", err)

		return nil, err
	}

	if len(gateways) == 0 {
		return nil, errors.New("no gateways available for the specified country")
	}

	sortGatewaysASC(gateways)

	for _, gateway := range gateways {
		if isGatewayHealthy(gateway) {
			return gateway, nil
		}
	}

	return nil, errors.New("gateways are unhealthy/unavailable")
}

func (g SvcGateway) SendTxToGateway(tx postgres.Transaction) (interface{}, error) {
	// having gateway_id in tx struct, we will select the required gateway and send tx for processing.
	return nil, nil
}

// sortGatewaysASC sort gateways in ascending order by priority
func sortGatewaysASC(gateways []*common.Gateway) {
	for i := 0; i < len(gateways); i++ {
		for j := i + 1; j < len(gateways); j++ {
			if gateways[i].Priority > gateways[j].Priority {
				gateways[i], gateways[j] = gateways[j], gateways[i]
			}
		}
	}
}

func isGatewayHealthy(gateway *common.Gateway) bool {
	// check gateway health
	return true
}
