package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// masks the sensitive information before publishing to Kafka using 64encoding but can be improved by using stronger algorithms and match it with the secret key.
func MaskData(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// generates a digital signature for the given data using the secret key to be used for verification later
func CreateSignature(data string, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	signature := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}

// checks if the provided signature matches the generated signature.
func VerifySignature(data string, secretKey string, signature string) bool {
	expectedSignature := CreateSignature(data, secretKey)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
