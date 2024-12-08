package main

import (
	"car_system/vehicle_service/config"
	"car_system/vehicle_service/controllers"
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
	router.HandleFunc("/available-vehicles", controllers.GetAvailableVehicles).Methods("GET")
	router.HandleFunc("/create-reservation", controllers.CreateReservation).Methods("POST")
	router.HandleFunc("/latest-reservation", controllers.GetLatestReservation).Methods("GET")

	// Serve static files
	staticDir := "./static/" // Directory where your static files are located
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	// Enable CORS for cross-origin requests
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8080"}), // Frontend origin
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	log.Println("User-service running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", cors(router)))

}
