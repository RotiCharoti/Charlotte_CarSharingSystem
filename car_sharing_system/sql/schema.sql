CREATE DATABASE car_system;

USE car_system;

--================================================================================================================
-- USER SERVICE -- 

-- User Table
CREATE TABLE User (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone_no VARCHAR(15) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    dob DATE NOT NULL,
    license_id INT UNIQUE,                         -- References DriverLicense table
    membership_tier VARCHAR(50) DEFAULT 'Basic',   -- References Membership table
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (license_id) REFERENCES DriverLicense(license_id),
    FOREIGN KEY (membership_tier) REFERENCES Membership(membership_tier)
);

-- Membership Table
CREATE TABLE Membership (
    membership_tier VARCHAR(50) PRIMARY KEY,     
    hourly_rate_discount DECIMAL(5, 2) NOT NULL, 
    priority_access BOOLEAN DEFAULT FALSE,       
    booking_limit INT NOT NULL                  
);

CREATE TABLE Rental_History (
    history_id SERIAL PRIMARY KEY,                   
    user_id INT NOT NULL,                            -- Reference to User table
    vehicle_id INT,                                  -- Reference to Vehicle table
    start_time TIMESTAMP NOT NULL,                 
    end_time TIMESTAMP NOT NULL,                     
    cost DECIMAL(10, 2) NOT NULL,                   
    status ENUM('Completed', 'Refunded', 'Cancelled') DEFAULT 'Completed', 
    FOREIGN KEY (user_id) REFERENCES User(user_id),  
    FOREIGN KEY (vehicle_id) REFERENCES Vehicle(vehicle_id)
);

-- Driver License
CREATE TABLE DriverLicense (
    license_id SERIAL PRIMARY KEY,               
    license_no VARCHAR(50) UNIQUE NOT NULL,       
    license_issue_date DATE NOT NULL,             
    license_expiry_date DATE NOT NULL,            
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

-- Admin

--================================================================================================================
-- VEHICLE SERVICE -- 

-- Add potential insurance attribute for renter and vehicle
-- Vehicle Table
CREATE TABLE Vehicle (
    vehicle_id SERIAL PRIMARY KEY,
    license_plate VARCHAR(20) UNIQUE NOT NULL,
    model VARCHAR(100) NOT NULL,
    charge_level DECIMAL(5, 2) NOT NULL,  
    location VARCHAR(255) NOT NULL, -- Location here refers to pick-up and return location
    rental_rate DECIMAL(10, 2) NOT NULL, 
    mileage INT NOT NULL,
    status ENUM('Operational', 'Decommissioned', 'Under Maintenance') DEFAULT 'Operational',
    battery_capacity_kwh DECIMAL(5, 2) DEFAULT NULL -- Optional battery capacity in kWh
);

-- Reservation Table
-- Validation 1: Rental of maximum 3 days, Min 1 hour (Trigger)
-- Validation 2: Let vehicle be rental free for at least 1 hour after reservation to recharge
-- Add logic in application layer/database triggers to prevent overlapping reservations
CREATE TABLE Reservation (
    reservation_id SERIAL PRIMARY KEY,          
    vehicle_id INT NOT NULL,                       -- References Vehicle table
    user_id INT NOT NULL,                          -- References User table
    start_time DATETIME NOT NULL,                 
    end_time DATETIME NOT NULL,                   
    expected_charge_level DECIMAL(5, 2) NOT NULL,  -- Predicted charge level after the trip
    status ENUM('Active', 'Completed', 'Cancelled') NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    FOREIGN KEY (vehicle_id) REFERENCES Vehicle(vehicle_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    CHECK (start_time < end_time)                  -- Ensure valid time range
)

-- Rental Table
CREATE TABLE Rental (
    rental_id SERIAL PRIMARY KEY,             
    reservation_id INT NOT NULL,                -- References the Reservation table
    start_date DATETIME NOT NULL,              
    end_date DATETIME NOT NULL,                 
    rental_fee DECIMAL(10, 2) NOT NULL,        
    payment_status ENUM('Pending', 'Paid', 'Refunded') NOT NULL, 
    payment_amount DECIMAL(10, 2) DEFAULT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    FOREIGN KEY (reservation_id) REFERENCES Reservation(reservation_id)
);

-- Return Location Table

-- Pick-up Location Table



--================================================================================================================
-- BILLING SERVICE -- 

-- Billing Table
CREATE TABLE Billing (
    bill_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,           -- Reference to User table
    reservation_id INT NOT NULL,    -- Reference to Reservation table
    promo_id INT DEFAULT NULL,      -- Reference to Promotion table
    amount DECIMAL(10, 2) NOT NULL,
    status ENUM('Pending', 'Paid', 'Refunded') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User(user_id),
    FOREIGN KEY (reservation_id) REFERENCES Reservation(reservation_id),
    FOREIGN KEY (promo_id) REFERENCES Promotion(promo_id) 
);

-- Promotion Table
CREATE TABLE Promotion (
    promo_id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(255) NOT NULL,
    discount_rate DECIMAL(5, 2) NOT NULL,
    valid_from TIMESTAMP NOT NULL,
    valid_to TIMESTAMP NOT NULL,
    usage_limit INT NOT NULL
);

-- PaymentMethod

-- GiftCard










