package services

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"payment-gateway/internal/models"
)

// Supported content types can easily be extended by adding more types here
var supportedContentTypes = map[string]bool{
	"application/json": true,
	"text/xml":         true,
	"application/xml":  true,
}

// checks if the content type is supported
func IsSupportedContentType(contentType string) bool {
	return supportedContentTypes[contentType]
}

// decodes the incoming request based on content type
func DecodeRequest(r *http.Request, request *models.TransactionRequest) error {
	contentType := r.Header.Get("Content-Type")
	if !IsSupportedContentType(contentType) {
		return fmt.Errorf("unsupported content type")
	}

	switch contentType {
	case "application/json":
		return json.NewDecoder(r.Body).Decode(request)
	case "text/xml":
		return xml.NewDecoder(r.Body).Decode(request)
	case "application/xml":
		return xml.NewDecoder(r.Body).Decode(request)
	default:
		return fmt.Errorf("unsupported content type")
	}
}

// ends the transaction response in the appropriate format
func RespondWithTransaction(w http.ResponseWriter, response models.APIResponse, contentType string) {

	w.Header().Set("Content-Type", contentType)

	w.WriteHeader(response.StatusCode)

	// encode and generate the response based on the content type to insure the response is in the correct format for the payment gateways
	switch contentType {
	case "application/json":
		json.NewEncoder(w).Encode(response)
	case "text/xml":
		xml.NewEncoder(w).Encode(response)
	default:
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
	}
}

// creates a standardized error response for better errorhandling
func RespondWithError(w http.ResponseWriter, statusCode int, message string, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	response := models.APIResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
	}

	// encode and generate the response based on the content type to insure the response is in the correct format for the payment gateways
	switch contentType {
	case "application/json":
		json.NewEncoder(w).Encode(response)
	case "text/xml":
		xml.NewEncoder(w).Encode(response)
	}
}
