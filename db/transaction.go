package db

import (
	"payment-gateway/internal/models"

	_ "github.com/lib/pq"
)

// Saves a transaction
func SaveTransaction(transaction models.Transaction) error {
	query := `
        INSERT INTO transactions (transaction_id, amount, type, status, created_at, data_format)
        VALUES ($1, $2, $3, $4, NOW(), $5)`

	_, err := db.Exec(query, transaction.TransactionID, transaction.Amount, transaction.Type, transaction.Status, transaction.DataFormat)
	return err
}

// updates the status of a transaction based on the transaction ID
func UpdateTransactionStatus(transactionID string, status string) error {
	query := `
        UPDATE transactions
        SET status = $1
        WHERE transaction_id = $2`

	_, err := db.Exec(query, status, transactionID)
	return err
}
