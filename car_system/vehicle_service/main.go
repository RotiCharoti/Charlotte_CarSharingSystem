package main

import (
	"car_system/vehicle_service/config"
	"car_system/vehicle_service/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func fly() {
	// Connect to the database
	config.ConnectDB()
	defer config.DB.Close()

	// Set up router
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/available-vehicles", controllers.GetAvailableVehicles).Methods("GET")
	router.HandleFunc("/create-reservation", controllers.CreateReservation).Methods("POST")

	// Serve static files
	staticDir := "./static/" // Directory where your static files are located
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	// Start the server
	log.Println("Vehicle Service running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", router))
}
