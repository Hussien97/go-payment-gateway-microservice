package services

import (
	"context"
	"fmt"
	"payment-gateway/db"
	"payment-gateway/internal/kafka"
	"payment-gateway/internal/models"
	"payment-gateway/internal/redis"
	"payment-gateway/internal/resilience"

	"sync"
)

// saves the transaction in the database and Redis concurrently for better performance
func SaveTransaction(transaction models.Transaction) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	// Save transaction in the database
	go func() {
		defer wg.Done()
		if err := db.SaveTransaction(transaction); err != nil {
			errChan <- err
		}
	}()

	// Save the transaction status in Redis
	go func() {
		defer wg.Done()
		redis.SetTransactionStatus(transaction.TransactionID, transaction.Status) // Assuming Redis set always succeeds
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// retrieves a transaction by its ID from the database
func GetTransactionByID(transactionID string) (models.Transaction, error) {
	return db.GetTransactionByID(transactionID) // Call the function in db package
}

// validates the transaction request (data fields)
func ValidateTransactionRequest(request models.TransactionRequest) error {

	if request.Amount == 0 {
		return fmt.Errorf("amount is required")
	}

	if request.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	return nil
}

// validates the callback request (data fields)
func ValidateCallbackRequest(request models.TransactionRequest) error {

	if request.TransactionID == "" {
		return fmt.Errorf("transaction ID is required")
	}

	if request.Status == "" {
		return fmt.Errorf("status is required")
	}

	// try to validate the transaction existence first from redis then from the database
	status, err := GetTransactionStatus(request.TransactionID)
	if err != nil {
		return fmt.Errorf("error retrieving transaction status: %v", err)
	}

	if status == "completed" {
		return fmt.Errorf("transaction has already been completed")
	}

	return nil
}

// retrieves the transaction status from Redis first if not found then will get from the database
func GetTransactionStatus(transactionID string) (string, error) {

	status, err := redis.GetTransactionStatus(transactionID)
	if err == nil {
		return status, nil
	}

	transaction, dbErr := db.GetTransactionByID(transactionID)
	if dbErr != nil {
		return "", dbErr
	}

	return transaction.Status, nil
}

// updates the transaction status in the database and Redis concurrently
func UpdateTransactionStatus(transactionID, status string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := db.UpdateTransactionStatus(transactionID, status); err != nil {
			errChan <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		redis.SetTransactionStatus(transactionID, status) // Assuming Redis set always succeeds
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// publishes a transaction to Kafka with retries and circuit breaker
func PublishTransaction(ctx context.Context, transactionID string, transactionData []byte, dataFormat string) error {
	return resilience.PublishWithCircuitBreaker(func() error {
		return kafka.PublishTransaction(ctx, transactionID, transactionData, dataFormat)
	})
}
