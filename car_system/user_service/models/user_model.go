package models

import (
	"car_system/user_service/config"
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// User structure
type User struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	PhoneNo  string `json:"phone_no"`
	Password string `json:"password,omitempty"`
	DOB      string `json:"dob"`
}

type Rental struct {
	HistoryID int     `json:"history_id"`
	VehicleID int     `json:"vehicle_id"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	Cost      float64 `json:"cost"`
	Status    string  `json:"status"`
}

type Membership struct {
	Tier               string  `json:"tier"`
	HourlyRateDiscount float64 `json:"hourly_rate_discount"`
	PriorityAccess     bool    `json:"priority_access"`
	BookingLimit       int     `json:"booking_limit"`
}

// RegisterUser inserts a new user into the database
func RegisterUser(user *User) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	user.Password = string(hashedPassword)

	// Insert the user into the database
	query := "INSERT INTO User (name, email, phone_no, password, dob) VALUES (?, ?, ?, ?, ?)"
	_, err = config.DB.Exec(query, user.Name, user.Email, user.PhoneNo, user.Password, user.DOB)
	if err != nil {
		return fmt.Errorf("failed to register user: %v", err)
	}
	return nil
}

// IsUserExists checks if a user with the given email or phone number already exists
func IsUserExists(email, phoneNo string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM User WHERE email = ? OR phone_no = ?"
	err := config.DB.QueryRow(query, email, phoneNo).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %v", err)
	}
	return count > 0, nil
}

// LoginUser authenticates a user by email and password
func LoginUser(email, password string) (*User, error) {
	var user User
	query := "SELECT user_id, name, email, phone_no, password, dob FROM User WHERE email = ?"
	err := config.DB.QueryRow(query, email).Scan(&user.UserID, &user.Name, &user.Email, &user.PhoneNo, &user.Password, &user.DOB)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil
	}

	// Omit the password from the response
	user.Password = ""
	return &user, nil
}

// GetRentalsByUserID fetches rental records for a specific user.
func GetRentalsByUserID(userID int) ([]Rental, error) {
	query := `
		SELECT history_id, vehicle_id, start_time, end_time, cost, status
		FROM Rental_History
		WHERE user_id = ?
	`
	rows, err := config.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error fetching rental records for user_id %d: %v\n", userID, err)
		return nil, err
	}
	defer rows.Close()

	var rentals []Rental
	for rows.Next() {
		var rental Rental
		if err := rows.Scan(&rental.HistoryID, &rental.VehicleID, &rental.StartTime, &rental.EndTime, &rental.Cost, &rental.Status); err != nil {
			log.Printf("Error scanning rental record: %v\n", err)
			return nil, err
		}
		rentals = append(rentals, rental)
	}

	return rentals, nil
}

// GetUserMembershipDetails fetches the membership tier and details for a specific user
func GetUserMembershipDetails(userID int) (*Membership, error) {
	var membership Membership

	query := `
        SELECT m.membership_tier, m.hourly_rate_discount, m.priority_access, m.booking_limit
        FROM Membership m
        INNER JOIN User u ON u.membership_tier = m.membership_tier
        WHERE u.user_id = ?
    `

	err := config.DB.QueryRow(query, userID).Scan(
		&membership.Tier,
		&membership.HourlyRateDiscount,
		&membership.PriorityAccess,
		&membership.BookingLimit,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching membership details: %v", err)
	}

	return &membership, nil
}

// GetUserDetailsByID fetches user details by their ID
func GetUserDetailsByID(userID int) (*User, error) {
	var user User
	query := "SELECT user_id, name, email, phone_no, dob FROM User WHERE user_id = ?"
	err := config.DB.QueryRow(query, userID).Scan(&user.UserID, &user.Name, &user.Email, &user.PhoneNo, &user.DOB)
	if err != nil {
		return nil, fmt.Errorf("error fetching user details: %v", err)
	}
	return &user, nil
}

// UpdateUserDetails updates the user details in the database
func UpdateUserDetails(userID int, user *User) error {
	var updateFields []interface{}
	query := "UPDATE User SET name = ?, email = ?, phone_no = ?, dob = ?"
	updateFields = append(updateFields, user.Name, user.Email, user.PhoneNo, user.DOB)

	// If password is provided, hash and include it
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}
		query += ", password = ?"
		updateFields = append(updateFields, string(hashedPassword))
	}

	query += " WHERE user_id = ?"
	updateFields = append(updateFields, userID)

	_, err := config.DB.Exec(query, updateFields...)
	if err != nil {
		return fmt.Errorf("failed to update user details: %v", err)
	}
	return nil
}
