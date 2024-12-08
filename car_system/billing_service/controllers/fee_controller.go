package controllers

import (
	"car_system/billing_service/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// CalculateRentalFee calculates the total rental fee based on reservation details
func CalculateRentalFee(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ReservationID int     `json:"reservation_id"`
		StartTime     string  `json:"start_time"`
		EndTime       string  `json:"end_time"`
		RentalRate    float64 `json:"rental_rate"`
	}

	// Decode the request payload
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"message":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Parse start and end times
	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		http.Error(w, `{"message":"Invalid start time format"}`, http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, request.EndTime)
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
	totalFee := duration * request.RentalRate
	log.Printf("CalculateRentalFee: Duration: %.2f hours, Rate: %.2f, Total Fee: %.2f", duration, request.RentalRate, totalFee)

	// Respond with the calculated fee
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Rental fee calculated successfully",
		"total_fee": totalFee,
	})

	log.Printf("Received payload for fee calculation: %+v", request)
}

// InsertBillingHandler handles inserting a new billing record
func InsertBillingHandler(w http.ResponseWriter, r *http.Request) {
	var billingRequest struct {
		UserID        int     `json:"user_id"`
		ReservationID int     `json:"reservation_id"`
		PromoID       *int    `json:"promo_id"` // Optional field for promotion ID
		Amount        float64 `json:"amount"`
		Status        string  `json:"status"` // 'Pending', 'Paid', 'Refunded'
	}

	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&billingRequest); err != nil {
		http.Error(w, `{"message":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if billingRequest.UserID == 0 || billingRequest.ReservationID == 0 || billingRequest.Amount <= 0 || billingRequest.Status == "" {
		http.Error(w, `{"message":"Missing or invalid fields in request"}`, http.StatusBadRequest)
		return
	}

	// Create Billing object
	billing := models.Billing{
		UserID:        billingRequest.UserID,
		ReservationID: billingRequest.ReservationID,
		PromoID:       billingRequest.PromoID,
		Amount:        billingRequest.Amount,
		Status:        billingRequest.Status,
	}

	// Insert into the database
	if err := models.InsertBilling(&billing); err != nil {
		http.Error(w, `{"message":"Failed to insert billing record"}`, http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Billing record inserted successfully",
	})
}
