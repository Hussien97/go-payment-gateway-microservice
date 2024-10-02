package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
)

type TransactionRequest struct {
	TransactionID string  `json:"transaction_id" xml:"transaction_id"`
	Amount        float64 `json:"amount" xml:"amount"`
	Status        string  `json:"status" xml:"status"`
}

type CallbackRequest struct {
	TransactionID string `xml:"transaction_id"`
	Status        string `xml:"status"`
}

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

	callbackData := CallbackRequest{
		TransactionID: transaction.TransactionID,
		Status:        "completed",
	}

	xmlData, err := XMLMarshaller(callbackData)
	if err != nil {
		log.Printf("Error marshalling callback data to XML: %v", err)
		return
	}

	callbackURL := "http://payment_gateway_app:8080/callback"
	resp, err := http.Post(callbackURL, "text/xml", bytes.NewBuffer(xmlData))
	if err != nil {
		log.Printf("Error calling callback endpoint: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Processed transaction: %s with status: %s", transaction.TransactionID, transaction.Status)
}

func UnmaskData(maskedData string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(maskedData)
}

func XMLMarshaller(req CallbackRequest) ([]byte, error) {
	return xml.Marshal(req)
}
