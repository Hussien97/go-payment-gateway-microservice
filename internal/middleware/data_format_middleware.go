package middleware

import (
	"net/http"
	"payment-gateway/internal/models"
	"payment-gateway/internal/services"
)

// checks if the request content type is supported and can easily have more data types supported by simply adding them in here
func DataFormatMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if !services.IsSupportedContentType(contentType) {
			services.RespondWithTransaction(w, models.APIResponse{
				StatusCode: http.StatusUnsupportedMediaType,
				Message:    "Unsupported content type",
			}, contentType)
			return
		}
		next.ServeHTTP(w, r)
	})
}
