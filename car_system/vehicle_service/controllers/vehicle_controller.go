package controllers

import (
	"car_system/vehicle_service/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// GetAvailableVehicles retrieves all vehicles from the database
func GetAvailableVehicles(w http.ResponseWriter, r *http.Request) {
	// Fetch available vehicles from the database
	vehicles, err := models.GetAvailableVehicles()
	if err != nil {
		http.Error(w, "Failed to fetch available vehicles", http.StatusInternalServerError)
		return
	}

	// Respond with the list of available vehicles
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Available vehicles fetched successfully",
		"vehicles": vehicles,
	})
}

// Retrieve user ID
func FetchUserIDFromSession() (int, error) {
	url := "http://localhost:8080/get-session-user-id"

	// Set up the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %v", err)
	}

	// Ensure the cookie is included
	req.Header.Set("Cookie", "user-session=<your-session-value>") // Replace with the actual session value

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error calling user_service: %v", err)
	}
	defer resp.Body.Close()

	// Log the response status and any errors
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("user_service returned status: %v", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %v", err)
	}

	userID, ok := result["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user_id format in response")
	}

	return int(userID), nil
}

// Reserve Vehicle
func CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservation models.Reservation

	// Decode the reservation details from the request body
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	// Validate the start_time and end_time formats
	if reservation.StartTime.IsZero() || reservation.EndTime.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Start time or end time is missing or invalid",
		})
		return
	}

	// Ensure duration is between 1 hour and 3 days
	duration := reservation.EndTime.Sub(reservation.StartTime)
	if duration < time.Hour || duration > 72*time.Hour {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Reservation duration must be between 1 hour and 3 days.",
		})
		return
	}

	// Fetch user_id from session
	userID, err := FetchUserIDFromSession()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Unauthorized user",
			"error":   err.Error(),
		})
		return
	}
	reservation.UserID = userID

	// Log the user ID in the terminal
	log.Printf("Reservation attempt by User ID: %d", userID)

	// Check vehicle availability
	available, err := models.IsVehicleAvailable(reservation.VehicleID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Error checking vehicle availability",
			"error":   err.Error(),
		})
		return
	}
	if !available {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Vehicle not available for the selected time range",
		})
		return
	}

	// Save the reservation in the database
	if err := models.CreateReservation(&reservation); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Failed to create reservation",
			"error":   err.Error(),
		})
		return
	}

	// Respond with success, including user_id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Reservation created successfully",
		"data":    reservation,
		"user_id": userID, // Include user_id in response
	})
}
