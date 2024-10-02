package db

import (
	"database/sql"
	"fmt"
	"log"
	"payment-gateway/internal/models"
	"payment-gateway/internal/resilience"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitializeDB initializes the database connection
func InitializeDB(dataSourceName string) {
	var err error

	// Retry connecting to the database in case of failure up to 5 times
	err = resilience.RetryOperation(func() error {
		db, err = sql.Open("postgres", dataSourceName)
		if err != nil {
			return err
		}

		return db.Ping()
	}, 5)

	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	log.Println("Successfully connected to the database.")
}

// GetTransactionByID retrieves a transaction by its ID
func GetTransactionByID(transactionID string) (models.Transaction, error) {
	var transaction models.Transaction

	query := `SELECT transaction_id, amount, type, status, data_format FROM transactions WHERE transaction_id = $1`
	err := resilience.RetryOperation(func() error {
		return db.QueryRow(query, transactionID).Scan(&transaction.TransactionID, &transaction.Amount, &transaction.Type, &transaction.Status, &transaction.DataFormat)
	}, 3)

	if err != nil {
		if err == sql.ErrNoRows {
			return transaction, fmt.Errorf("transaction not found")
		}
		return transaction, err
	}

	return transaction, nil
}
