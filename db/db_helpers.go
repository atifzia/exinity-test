package db

import (
	"database/sql"
	"fmt"
	"log"
	"payment-gateway/internal/models/common"
	"payment-gateway/internal/models/postgres"
	"payment-gateway/internal/util"
	"time"

	_ "github.com/lib/pq"
)

type (
	DB struct {
		db *sql.DB
	}

	Idb interface {
		GetSupportedGatewaysByCountry(countryID int) ([]*common.Gateway, error)
		CreateTransaction(tx *postgres.Transaction) error
		UpdateTxStatus(txID int64, status string) error
	}
)

func New(dsn string) (Idb, error) {
	var d DB
	err := util.RetryOperation(func() error {
		dbInst, err := sql.Open("postgres", dsn)
		if err != nil {
			return err
		}

		if err = dbInst.Ping(); err != nil {
			return err
		}
		d.db = dbInst

		return nil
	}, 5)
	if err != nil {
		log.Printf("Could not connect to the database: %v\n", err)

		return nil, err
	}

	return &d, nil
}

func (d *DB) CreateTransaction(transaction *postgres.Transaction) error {
	query := `INSERT INTO transactions (amount, type, status, gateway_id, country_id, user_id, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := d.db.QueryRow(query, transaction.Amount, transaction.Type, transaction.Status, transaction.GatewayID, transaction.CountryID, transaction.UserID, time.Now()).Scan(&transaction.ID)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}
	return nil
}

func (d *DB) UpdateTxStatus(txId int64, status string) error {
	query := `UPDATE transactions SET status = $1 WHERE id = $2`
	_, err := d.db.Exec(query, status, txId)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %v", err)
	}

	return nil
}

func (d *DB) GetSupportedGatewaysByCountry(countryId int) ([]*common.Gateway, error) {
	query := `
		SELECT g.id, g.name, g.data_format_supported, g.priority
		FROM gateways g
		INNER JOIN gateway_countries gc ON g.id = gc.gateway_id
		WHERE gc.country_id = $1
		ORDER BY g.name
	`

	rows, err := d.db.Query(query, countryId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gateways for country %d: %v", countryId, err)
	}
	defer rows.Close()

	var gateways []*common.Gateway
	for rows.Next() {
		var gt common.Gateway
		if err = rows.Scan(&gt.ID, &gt.Name, &gt.DataFormatSupported, &gt.Priority); err != nil {
			return nil, fmt.Errorf("failed to scan gateway: %v", err)
		}
		gateways = append(gateways, &gt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	return gateways, nil
}
