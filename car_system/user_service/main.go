package main

import (
	"car_system/user_service/config"
	"car_system/user_service/controllers"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
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

	// Routes
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/rental-records", controllers.DisplayRentalRecords).Methods("GET")
	router.HandleFunc("/membership-details", controllers.DisplayUserMembership).Methods("GET")
	router.HandleFunc("/view-details", controllers.DisplayUserDetails).Methods("GET")
	router.HandleFunc("/update-details", controllers.UpdateUserDetails).Methods("PUT")
	router.HandleFunc("/get-session-user-id", controllers.GetSessionUserID).Methods("GET")

	// Middleware
	router.Use(loggingMiddleware)
	router.Use(recoverMiddleware)

	// Serve static files
	staticDir := "./static/" // Directory where your static files are located
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	// Enable CORS for cross-origin requests
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8081"}), // Adjust for your setup
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	// Start the server
	log.Println("User-service running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", cors(router)))
}

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Middleware for recovering from panics
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
