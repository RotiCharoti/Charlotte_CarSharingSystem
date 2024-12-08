package main

import (
	"car_system/user_service/config"
	"car_system/user_service/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Connect to the database
	config.ConnectDB()
	defer config.DB.Close()

	// Initialize session store globally in controllers
	controllers.InitializeSessionStore()

	// Set up router
	router := mux.NewRouter()

	// API Routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	api.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	api.HandleFunc("/rental-records", controllers.DisplayRentalRecords).Methods("GET")
	api.HandleFunc("/membership-details", controllers.DisplayUserMembership).Methods("GET")
	api.HandleFunc("/view-details", controllers.DisplayUserDetails).Methods("GET")
	api.HandleFunc("/update-details", controllers.UpdateUserDetails).Methods("PUT")
	api.HandleFunc("/proxy-available-vehicles", controllers.ProxyAvailableVehicles).Methods("GET")
	api.HandleFunc("/proxy-create-reservation", controllers.ProxyCreateReservation).Methods("POST")
	api.HandleFunc("/proxy-get-latest-reservation", controllers.ProxyGetLatestReservation).Methods("GET")
	api.HandleFunc("/proxy-calculate-rental-fee", controllers.ProxyCalculateRentalFee).Methods("POST")

	// Serve static files
	staticDir := "./static/"
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	// Start the server
	log.Println("User-service running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
