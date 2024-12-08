package controllers

import (
	"bytes"
	"car_system/user_service/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

// Declare session store
var store *sessions.CookieStore

// InitializeSessionStore initializes the session store with a secret key from .env
func InitializeSessionStore() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("SESSION_SECRET")
	if secretKey == "" {
		log.Fatal("SESSION_SECRET is not set in the environment")
	}

	store = sessions.NewCookieStore([]byte(secretKey))
	log.Println("Session store initialized successfully")
}

// Response structure for API responses
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// sendErrorResponse is a helper function to send error messages
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{Message: message})
}

// RegisterUser handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if user.Name == "" || user.Email == "" || user.PhoneNo == "" || user.Password == "" || user.DOB == "" {
		sendErrorResponse(w, "All fields (name, email, phone_no, password, dob) are required", http.StatusBadRequest)
		return
	}

	// Check for duplicate email or phone number
	exists, err := models.IsUserExists(user.Email, user.PhoneNo)
	if err != nil {
		log.Printf("Error checking for duplicate user: %v", err)
		sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		sendErrorResponse(w, "Email or phone number already exists", http.StatusConflict)
		return
	}

	// Register the user
	err = models.RegisterUser(&user)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		sendErrorResponse(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "User registered successfully"})
}

// User Login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Println("Error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	// Authenticate user
	user, err := models.LoginUser(credentials.Email, credentials.Password)
	if err != nil || user == nil {
		log.Printf("Invalid login attempt for email: %s\n", credentials.Email)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid email or password",
		})
		return
	}

	// Create a session
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println("Error creating session:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Could not create session",
			"error":   err.Error(),
		})
		return
	}

	// Set session values
	session.Values["user_id"] = user.UserID
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
	}

	// Save session
	if err := session.Save(r, w); err != nil {
		log.Println("Error saving session:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Could not save session",
			"error":   err.Error(),
		})
		return
	}

	// Retrieve the session cookie
	cookie, err := r.Cookie("user-session")
	if err != nil {
		log.Println("Error retrieving session cookie:", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Session cookie is missing. Please log in again.",
		})
		return
	}

	// Log session details in the terminal
	log.Printf("Login successful: User ID: %d, Session Cookie Value: %s\n", user.UserID, cookie.Value)

	// Log session details in the app console
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Login successful",
		"user_id":    user.UserID,
		"session_id": cookie.Value, // Pass the session cookie value to the frontend
	})
}

// Display Rental Records of User
func DisplayRentalRecords(w http.ResponseWriter, r *http.Request) {
	// Retrieve the session
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println("Error retrieving session:", err)
		http.Error(w, "Session error. Please log in again.", http.StatusUnauthorized)
		return
	}

	// Retrieve the user ID from the session
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("Invalid or missing user ID in session")
		http.Error(w, "Unauthorized access. Please log in again.", http.StatusUnauthorized)
		return
	}

	// Fetch rental records for the user
	rentals, err := models.GetRentalsByUserID(userID)
	if err != nil {
		log.Printf("Error fetching rental records for user_id %d: %v\n", userID, err)
		http.Error(w, "Failed to fetch rental records", http.StatusInternalServerError)
		return
	}

	// Respond with rental records
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Rental records fetched successfully",
		"data":    rentals,
	})
}

// Display Membership Details of User
func DisplayUserMembership(w http.ResponseWriter, r *http.Request) {
	// Retrieve session
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println("Error retrieving session:", err)
		http.Error(w, "Session error. Please log in again.", http.StatusUnauthorized)
		return
	}

	// Retrieve user ID from session
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("Invalid or missing user ID in session")
		http.Error(w, "Unauthorized access. Please log in again.", http.StatusUnauthorized)
		return
	}

	// Fetch membership tier details from the model
	membership, err := models.GetUserMembershipDetails(userID)
	if err != nil {
		log.Printf("Error fetching membership details for user_id %d: %v\n", userID, err)
		http.Error(w, "Failed to fetch membership details", http.StatusInternalServerError)
		return
	}

	// Respond with membership details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Membership details fetched successfully",
		"data":    membership,
	})
}

// DisplayUserDetails fetches and returns the logged-in user's details
func DisplayUserDetails(w http.ResponseWriter, r *http.Request) {
	// Retrieve session
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println("Error retrieving session:", err)
		http.Error(w, "Session error. Please log in again.", http.StatusUnauthorized)
		return
	}

	// Retrieve user ID from session
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("Invalid or missing user ID in session")
		http.Error(w, "Unauthorized access. Please log in again.", http.StatusUnauthorized)
		return
	}

	// Fetch user details from the model
	user, err := models.GetUserDetailsByID(userID)
	if err != nil {
		log.Printf("Error fetching user details for user_id %d: %v\n", userID, err)
		http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
		return
	}

	// Respond with user details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User details fetched successfully",
		"data":    user,
	})
}

// UpdateUserDetails updates the logged-in user's details
func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println("Error retrieving session:", err)
		http.Error(w, `{"message":"Session error. Please log in again."}`, http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("Invalid or missing user ID in session")
		http.Error(w, `{"message":"Unauthorized access. Please log in again."}`, http.StatusUnauthorized)
		return
	}

	var userDetails models.User
	if err := json.NewDecoder(r.Body).Decode(&userDetails); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, `{"message":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	err = models.UpdateUserDetails(userID, &userDetails)
	if err != nil {
		log.Printf("Error updating user details for user_id %d: %v\n", userID, err)
		http.Error(w, `{"message":"Failed to update user details"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User details updated successfully",
	})
}

// ProxyAvailableVehicles fetches available vehicles from the vehicle_service
func ProxyAvailableVehicles(w http.ResponseWriter, r *http.Request) {
	vehicleServiceURL := "http://localhost:8081/available-vehicles"

	// Forward the request to the vehicle_service
	resp, err := http.Get(vehicleServiceURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch available vehicles: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Forward the response body directly to the frontend
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Failed to forward response", http.StatusInternalServerError)
	}
}

// ProxyCreateReservation proxies reservation creation requests to vehicle_service
func ProxyCreateReservation(w http.ResponseWriter, r *http.Request) {
	vehicleServiceURL := "http://localhost:8081/create-reservation"

	// Retrieve session
	session, err := store.Get(r, "user-session")
	if err != nil || session.Values["user_id"] == nil {
		log.Printf("Error retrieving session: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Unauthorized user. Please log in.",
		})
		return
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Printf("Invalid session data: %v", session.Values)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Session data is corrupted",
		})
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Inject user_id into the request payload
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	payload["user_id"] = userID

	// Forward the request to vehicle_service
	proxyBody, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to marshal request payload", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", vehicleServiceURL, bytes.NewReader(proxyBody))
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header // Forward headers
	req.Header.Set("Content-Type", "application/json")

	// Forward session cookies
	if cookie, err := r.Cookie("user-session"); err == nil {
		req.AddCookie(cookie)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to communicate with vehicle_service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Log headers and response
	log.Printf("Forwarding headers: %+v", req.Header)
	log.Printf("Response from vehicle_service: %d", resp.StatusCode)

	// Forward the response from vehicle_service
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from vehicle_service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	if _, err := w.Write(respBody); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// ProxyGetLatestReservation proxies the request to fetch the latest reservation for the logged-in user
func ProxyGetLatestReservation(w http.ResponseWriter, r *http.Request) {
	vehicleServiceURL := "http://localhost:8081/latest-reservation"

	// Retrieve session
	session, err := store.Get(r, "user-session")
	if err != nil || session.Values["user_id"] == nil {
		log.Printf("Error retrieving session: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Unauthorized user. Please log in.",
		})
		return
	}

	// Get user_id from session
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Printf("Invalid session data: %v", session.Values)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Session data is corrupted",
		})
		return
	}

	// Log for debugging
	log.Printf("ProxyGetLatestReservation: Retrieved User ID: %d", userID)

	// Prepare the request to vehicle_service
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?user_id=%d", vehicleServiceURL, userID), nil)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header // Forward headers

	// Forward session cookies
	if cookie, err := r.Cookie("user-session"); err == nil {
		req.AddCookie(cookie)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to communicate with vehicle_service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Log headers and response status for debugging
	log.Printf("Forwarding headers: %+v", req.Header)
	log.Printf("Response from vehicle_service: %d", resp.StatusCode)

	// Forward the response from vehicle_service
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from vehicle_service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	if _, err := w.Write(respBody); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// ProxyCalculateRentalFee calculates the total fee based on vehicle rental rate and reservation duration
func ProxyCalculateRentalFee(w http.ResponseWriter, r *http.Request) {
	billingServiceURL := "http://localhost:8082/calculate-rental-fee"
	vehicleServiceURL := "http://localhost:8081/get-vehicle-details" // Assuming an endpoint to fetch vehicle details

	// Read and parse the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"message":"Failed to read request body"}`, http.StatusBadRequest)
		return
	}

	var payload struct {
		ReservationID int    `json:"reservation_id"`
		StartTime     string `json:"start_time"`
		EndTime       string `json:"end_time"`
		VehicleID     int    `json:"vehicle_id"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, `{"message":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Validate payload
	if payload.VehicleID == 0 || payload.StartTime == "" || payload.EndTime == "" {
		http.Error(w, `{"message":"Vehicle ID, start time, and end time are required"}`, http.StatusBadRequest)
		return
	}

	// Fetch vehicle details to get the rental rate
	vehicleDetails, err := fetchVehicleDetails(vehicleServiceURL, payload.VehicleID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"Failed to fetch vehicle details: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Parse the start and end times
	startTime, err := time.Parse(time.RFC3339, payload.StartTime)
	if err != nil {
		http.Error(w, `{"message":"Invalid start time format"}`, http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, payload.EndTime)
	if err != nil {
		http.Error(w, `{"message":"Invalid end time format"}`, http.StatusBadRequest)
		return
	}

	// Calculate duration in hours
	duration := endTime.Sub(startTime).Hours()
	if duration <= 0 {
		http.Error(w, `{"message":"Invalid time range"}`, http.StatusBadRequest)
		return
	}

	// Calculate total fee
	totalFee := duration * vehicleDetails.RentalRate

	// Prepare the final payload for billing service
	billingPayload := map[string]interface{}{
		"reservation_id": payload.ReservationID,
		"start_time":     payload.StartTime,
		"end_time":       payload.EndTime,
		"rental_rate":    vehicleDetails.RentalRate,
		"total_fee":      totalFee,
	}

	// Forward to billing service
	respBody, err := forwardToBillingService(billingServiceURL, billingPayload)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":"Failed to communicate with billing service: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Forward the response from billing service
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

// fetchVehicleDetails fetches the rental rate of the vehicle from vehicle_service
func fetchVehicleDetails(vehicleServiceURL string, vehicleID int) (*struct {
	RentalRate float64 `json:"rental_rate"`
}, error) {
	url := fmt.Sprintf("%s?vehicle_id=%d", vehicleServiceURL, vehicleID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch vehicle details, status: %d", resp.StatusCode)
	}

	var vehicleDetails struct {
		RentalRate float64 `json:"rental_rate"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&vehicleDetails); err != nil {
		return nil, err
	}

	return &vehicleDetails, nil
}

// forwardToBillingService forwards the calculated payload to billing_service
func forwardToBillingService(billingServiceURL string, payload map[string]interface{}) ([]byte, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal billing payload: %v", err)
	}

	req, err := http.NewRequest("POST", billingServiceURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request to billing service: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to communicate with billing service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("billing service error: %s", string(respBody))
	}

	return io.ReadAll(resp.Body)
}
