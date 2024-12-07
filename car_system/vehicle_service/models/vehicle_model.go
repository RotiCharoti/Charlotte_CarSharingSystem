package models

import (
	"car_system/vehicle_service/config"
	"log"
	"time"
)

// Reservation represents a vehicle reservation
type Reservation struct {
	ReservationID       int       `json:"reservation_id"`
	VehicleID           int       `json:"vehicle_id"`
	UserID              int       `json:"user_id"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	ExpectedChargeLevel float64   `json:"expected_charge_level"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
}

// Vehicle represents a vehicle in the database
type Vehicle struct {
	VehicleID          int     `json:"vehicle_id"`
	LicensePlate       string  `json:"license_plate"`
	Model              string  `json:"model"`
	ChargeLevel        float64 `json:"charge_level"`
	Location           string  `json:"location"`
	RentalRate         float64 `json:"rental_rate"`
	Mileage            int     `json:"mileage"`
	Status             string  `json:"status"`
	BatteryCapacityKWH float64 `json:"battery_capacity_kwh,omitempty"`
	ReservationStatus  string  `json:"reservation_status"`
	Cleanliness        string  `json:"cleanliness"`
}

// GetAllVehicles retrieves all vehicles from the database
func GetAvailableVehicles() ([]Vehicle, error) {
	query := `
		SELECT 
			vehicle_id, license_plate, model, charge_level, location, rental_rate, mileage, status, battery_capacity_kwh, reservation_status
		FROM Vehicle
		WHERE reservation_status = 'Available'
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		log.Printf("Error querying available vehicles: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var v Vehicle
		if err := rows.Scan(&v.VehicleID, &v.LicensePlate, &v.Model, &v.ChargeLevel, &v.Location, &v.RentalRate, &v.Mileage, &v.Status, &v.BatteryCapacityKWH, &v.ReservationStatus); err != nil {
			log.Printf("Error scanning vehicle: %v\n", err)
			return nil, err
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

// IsVehicleAvailable checks if a vehicle is available for a specific time range
func IsVehicleAvailable(vehicleID int, startTime, endTime time.Time) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM Reservation
		WHERE vehicle_id = ?
		  AND ((start_time < ? AND end_time > ?)
		    OR (start_time < ? AND end_time > ?))
	`
	var count int
	err := config.DB.QueryRow(query, vehicleID, endTime, startTime, endTime, startTime).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// CreateReservation inserts a new reservation into the database
func CreateReservation(reservation *Reservation) error {
	query := `
		INSERT INTO Reservation (vehicle_id, user_id, start_time, end_time, expected_charge_level, status)
		VALUES (?, ?, ?, ?, ?, 'Active')
	`
	_, err := config.DB.Exec(query, reservation.VehicleID, reservation.UserID, reservation.StartTime, reservation.EndTime, reservation.ExpectedChargeLevel)
	return err
}
