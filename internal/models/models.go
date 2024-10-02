package models

import "time"

// a transaction model
type Transaction struct {
	ID            int       `json:"id"`
	TransactionID string    `json:"transaction_id"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"` // deposit or withdrawal
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	DataFormat    string    `json:"data_format"`
}

// a standard request structure for the APIs
type TransactionRequest struct {
	Type          string  `json:"type" xml:"type"`
	Amount        float64 `json:"amount" xml:"amount"`
	TransactionID string  `json:"transaction_id,omitempty" xml:"transaction_id,omitempty"`
	Status        string  `json:"status,omitempty" xml:"status,omitempty"`
}

// a standard response structure for the APIs
type APIResponse struct {
	StatusCode int         `json:"status_code" xml:"status_code"`
	Message    string      `json:"message" xml:"message"`
	Data       interface{} `json:"data,omitempty" xml:"data,omitempty"`
}
