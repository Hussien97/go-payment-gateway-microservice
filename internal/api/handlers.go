package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"payment-gateway/db"
	"payment-gateway/internal/models"
	"payment-gateway/internal/security"
	"payment-gateway/internal/services"

	"github.com/google/uuid"
)

// handles deposit requests
func DepositHandler(w http.ResponseWriter, r *http.Request) {
	handleTransaction(w, r, "deposit")
}

// handles withdrawal requests
func WithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	handleTransaction(w, r, "withdrawal")
}

// processes both deposit and withdrawal requests in same logic
func handleTransaction(w http.ResponseWriter, r *http.Request, transactionType string) {
	var request models.TransactionRequest

	// Decode the request based on content type
	if err := services.DecodeRequest(r, &request); err != nil {
		services.RespondWithTransaction(w, models.APIResponse{
			StatusCode: http.StatusUnsupportedMediaType,
			Message:    "Invalid request format",
		}, r.Header.Get("Content-Type"))
		return
	}

	// Validate the transaction request (data fields)
	if err := services.ValidateTransactionRequest(request); err != nil {
		services.RespondWithTransaction(w, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}, r.Header.Get("Content-Type"))
		return
	}

	// Generate a unique transaction ID and create a transaction object
	transactionID := uuid.New().String()

	transaction := models.Transaction{
		TransactionID: transactionID,
		Amount:        request.Amount,
		Type:          transactionType,
		Status:        "pending",
		DataFormat:    r.Header.Get("Content-Type"),
	}

	// Save transaction in the database and Redis concurrently
	if err := services.SaveTransaction(transaction); err != nil {
		services.RespondWithTransaction(w, models.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save transaction",
		}, r.Header.Get("Content-Type"))
		return
	}

	// Serialize transaction data for Kafka
	transactionData, err := json.Marshal(transaction)
	if err != nil {
		services.RespondWithTransaction(w, models.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to serialize transaction data",
		}, r.Header.Get("Content-Type"))
		return
	}

	// Mask the transaction data before sending it to Kafka for better security
	maskedData := security.MaskData(transactionData)

	// Publish transaction to Kafka
	ctx := context.Background()
	if err := services.PublishTransaction(ctx, transaction.TransactionID, []byte(maskedData), r.Header.Get("Content-Type")); err != nil {
		log.Printf("failed to publish transaction to Kafka: %v", err)

		// Mark the transaction as failed after the circuit breaker threshold to prevent any duplicate processing
		db.UpdateTransactionStatus(transaction.TransactionID, "failed")

		services.RespondWithTransaction(w, models.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to process request due to publishing error",
		}, r.Header.Get("Content-Type"))
		return
	}

	// Success Response
	services.RespondWithTransaction(w, models.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Transaction processed successfully",
		Data:       transaction,
	}, r.Header.Get("Content-Type"))
}

// handles callbacks from payment gateways
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	var transactionRequest models.TransactionRequest

	// Decode the callback request
	if err := services.DecodeRequest(r, &transactionRequest); err != nil {
		services.RespondWithTransaction(w, models.APIResponse{
			StatusCode: http.StatusUnsupportedMediaType,
			Message:    "Invalid request",
		}, r.Header.Get("Content-Type"))
		return
	}

	// Validate the callback request
	if err := services.ValidateCallbackRequest(transactionRequest); err != nil {
		services.RespondWithTransaction(w, models.APIResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}, r.Header.Get("Content-Type"))
		return
	}

	// Update the transaction status in the database and Redis concurrently
	if err := services.UpdateTransactionStatus(transactionRequest.TransactionID, transactionRequest.Status); err != nil {
		services.RespondWithTransaction(w, models.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to process request",
		}, r.Header.Get("Content-Type"))
		return
	}

	// Success Response
	services.RespondWithTransaction(w, models.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Callback processed successfully",
		Data:       transactionRequest,
	}, r.Header.Get("Content-Type"))
}
