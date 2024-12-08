package controllers

import (
	"car_system/vehicle_service/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

// Reserve Vehicle
func CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservation models.Reservation

	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	// Validate user_id
	if reservation.UserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User ID is missing or invalid",
		})
		return
	}

	log.Printf("Reservation attempt by User ID: %d", reservation.UserID)

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

	// Save reservation
	if err := models.CreateReservation(&reservation); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Failed to create reservation",
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Reservation created successfully",
		"data":    reservation,
	})
}

// GetLatestReservation fetches the latest reservation for a user by their ID
func GetLatestReservation(w http.ResponseWriter, r *http.Request) {
	userIDParam := r.URL.Query().Get("user_id")
	if userIDParam == "" {
		http.Error(w, `{"message":"User ID is required"}`, http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		http.Error(w, `{"message":"Invalid User ID"}`, http.StatusBadRequest)
		return
	}

	reservation, err := models.GetLatestReservationByUserID(userID)
	if err != nil {
		log.Printf("Error retrieving latest reservation for user ID %d: %v", userID, err)
		http.Error(w, `{"message":"Failed to retrieve reservation"}`, http.StatusInternalServerError)
		return
	}

	if reservation == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "No reservations found for the user",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Latest reservation fetched successfully",
		"data":    reservation,
	})
}
