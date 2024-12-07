package controllers

import (
	"car_system/user_service/models"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

// Declare session store
var store *sessions.CookieStore

// InitializeSessionStore initializes the session store with a secret key from .env
func InitializeSessionStore() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the secret key
	secretKey := os.Getenv("SESSION_SECRET")
	if secretKey == "" {
		log.Fatal("SESSION_SECRET is not set in the environment")
	}

	// Initialize the session store
	store = sessions.NewCookieStore([]byte(secretKey))
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

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := models.LoginUser(credentials.Email, credentials.Password)
	if err != nil || user == nil {
		log.Printf("Invalid login attempt for email: %s\n", credentials.Email)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println("Error creating session:", err)
		http.Error(w, "Could not create session", http.StatusInternalServerError)
		return
	}

	session.Values["user_id"] = user.UserID
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
	}

	if err := session.Save(r, w); err != nil {
		log.Println("Error saving session:", err)
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	log.Printf("Session after save: %v", session.Values)

	// Debugging cookies
	cookies := r.Cookies()
	for _, cookie := range cookies {
		log.Printf("Cookie: %s = %s", cookie.Name, cookie.Value)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"user_id": user.UserID,
		"session": session.Values, // Debugging session content
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

	// Decode the updated details from the request body
	var userDetails models.User
	if err := json.NewDecoder(r.Body).Decode(&userDetails); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update user details in the model
	err = models.UpdateUserDetails(userID, &userDetails)
	if err != nil {
		log.Printf("Error updating user details for user_id %d: %v\n", userID, err)
		http.Error(w, "Failed to update user details", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User details updated successfully",
	})
}

// For vehicle_service to retrieve user id
func GetSessionUserID(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user-session")
	if err != nil {
		log.Printf("Error retrieving cookie: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "No active session or user not logged in",
		})
		return
	}

	log.Printf("Raw cookie value: %s", cookie.Value)

	session, err := store.Get(r, "user-session")
	if err != nil || session.Values["user_id"] == nil {
		log.Printf("Error decoding session: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid or expired session",
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

	log.Printf("Successfully retrieved user ID: %d", userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
	})
}
