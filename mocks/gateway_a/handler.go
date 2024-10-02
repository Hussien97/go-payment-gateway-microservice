package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

type TransactionRequest struct {
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
}

type CallbackRequest struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

// processes the received transaction from Kafka.
func ProcessTransaction(data []byte) {

	unmaskedData, err := UnmaskData(string(data))
	if err != nil {
		log.Printf("Error unmasking data: %v", err)
		return
	}

	var transaction TransactionRequest
	err = json.Unmarshal(unmaskedData, &transaction)
	if err != nil {
		log.Printf("Error unmarshalling transaction data: %v", err)
		return
	}

	callbackRequest := CallbackRequest{
		TransactionID: transaction.TransactionID,
		Status:        "completed",
	}

	callbackData, err := json.Marshal(callbackRequest)
	if err != nil {
		log.Printf("Error marshalling callback data: %v", err)
		return
	}

	callbackURL := "http://payment_gateway_app:8080/callback"
	resp, err := http.Post(callbackURL, "application/json", bytes.NewBuffer(callbackData))
	if err != nil {
		log.Printf("Error calling callback endpoint: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Processed transaction from Gateway A: %s with status: %s", transaction.TransactionID, transaction.Status)
}

func UnmaskData(maskedData string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(maskedData)
}
