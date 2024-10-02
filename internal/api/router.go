package api

import (
	"net/http"
	"payment-gateway/internal/middleware"

	"github.com/gorilla/mux"
)

// initializes the router and routes
func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Apply the data format middleware to all routes and for future improvement we can apply verify_signature middleware on the callback route
	router.Handle("/deposit", middleware.DataFormatMiddleware(http.HandlerFunc(DepositHandler))).Methods("POST")
	router.Handle("/withdrawal", middleware.DataFormatMiddleware(http.HandlerFunc(WithdrawalHandler))).Methods("POST")
	router.Handle("/callback", middleware.DataFormatMiddleware(http.HandlerFunc(CallbackHandler))).Methods("POST")

	return router
}
