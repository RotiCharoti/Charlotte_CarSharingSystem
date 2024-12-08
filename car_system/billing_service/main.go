package main

import (
	"car_system/billing_service/config"
	"car_system/billing_service/controllers"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Connect to the database
	config.ConnectDB()
	defer config.DB.Close()

	// Set up router
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/calculate-rental-fee", controllers.CalculateRentalFee).Methods("POST")
	router.HandleFunc("/billing", controllers.InsertBillingHandler).Methods("POST")

	// Serve static files if needed (adjust directory as per your frontend setup)
	staticDir := "./static/" // Directory where your static files are located
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	// Enable CORS for cross-origin requests
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8080"}), // Frontend origin
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	log.Println("Billing service running on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", cors(router)))
}
