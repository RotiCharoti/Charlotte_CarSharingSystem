package models

import (
	"car_system/billing_service/config"
	"time"
)

// Billing represents a billing record in the database
type Billing struct {
	BillID        int       `json:"bill_id"`
	UserID        int       `json:"user_id"`
	ReservationID int       `json:"reservation_id"`
	PromoID       *int      `json:"promo_id"` // Pointer to allow NULL values
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// InsertBilling inserts a new billing record into the database
func InsertBilling(billing *Billing) error {
	query := `
		INSERT INTO Billing (user_id, reservation_id, promo_id, amount, status)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := config.DB.Exec(query, billing.UserID, billing.ReservationID, billing.PromoID, billing.Amount, billing.Status)
	return err
}
