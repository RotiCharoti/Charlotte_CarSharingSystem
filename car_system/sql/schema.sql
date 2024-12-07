--================================================================================================================
-- USER SERVICE -- 

DROP DATABASE IF EXISTS user_service;
CREATE DATABASE user_service;
USE user_service;

-- Driver License
CREATE TABLE DriverLicense (
    license_id INT AUTO_INCREMENT PRIMARY KEY,
    license_no VARCHAR(50) UNIQUE NOT NULL,
    license_issue_date DATE NOT NULL,
    license_expiry_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Membership Table
CREATE TABLE Membership (
    membership_tier VARCHAR(50) PRIMARY KEY,
    hourly_rate_discount DECIMAL(5, 2) NOT NULL,
    priority_access BOOLEAN DEFAULT FALSE,
    booking_limit INT NOT NULL
);

-- User Table
CREATE TABLE User (
    user_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone_no VARCHAR(15) UNIQUE NOT NULL,
    password VARCHAR(500) NOT NULL,
    dob DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    membership_tier VARCHAR(50) DEFAULT 'Basic',
    FOREIGN KEY (membership_tier) REFERENCES Membership(membership_tier)
);

-- Rental History Table
CREATE TABLE Rental_History (
    history_id SERIAL PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,  -- Matches User.user_id
    vehicle_id INT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    cost DECIMAL(10, 2) NOT NULL,
    status ENUM('Completed', 'Refunded', 'Cancelled') DEFAULT 'Completed',
    FOREIGN KEY (user_id) REFERENCES User(user_id)
);

INSERT INTO membership (membership_tier, hourly_rate_discount, priority_access, booking_limit)
VALUES
('Basic', 0.00, FALSE, 5),
('Premium', 10.00, TRUE, 10),
('VIP', 20.00, TRUE, 20);

--================================================================================================================
-- VEHICLE SERVICE -- 

DROP DATABASE IF EXISTS vehicle_service;
CREATE DATABASE vehicle_service;
USE vehicle_service;

-- Vehicle Table
CREATE TABLE Vehicle (
    vehicle_id SERIAL PRIMARY KEY,
    license_plate VARCHAR(20) UNIQUE NOT NULL,
    model VARCHAR(100) NOT NULL,
    charge_level DECIMAL(5, 2) NOT NULL,
    location VARCHAR(255) NOT NULL,
    rental_rate DECIMAL(10, 2) NOT NULL,
    mileage INT NOT NULL,
    status ENUM('Operational', 'Decommissioned', 'Under Maintenance') DEFAULT 'Operational',
    battery_capacity_kwh DECIMAL(5, 2) DEFAULT NULL
);

-- Reservation
CREATE TABLE Reservation (
    reservation_id SERIAL PRIMARY KEY,
    vehicle_id INT NOT NULL,
    user_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    expected_charge_level DECIMAL(5, 2) NOT NULL,
    status ENUM('Active', 'Completed', 'Cancelled') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vehicle_id) REFERENCES Vehicle(vehicle_id),
    FOREIGN KEY (user_id) REFERENCES User(user_id),
    CHECK (start_time < end_time)
);

-- Rental Table
CREATE TABLE Rental (
    rental_id SERIAL PRIMARY KEY,
    reservation_id INT NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    rental_fee DECIMAL(10, 2) NOT NULL,
    payment_status ENUM('Pending', 'Paid', 'Refunded') NOT NULL,
    payment_amount DECIMAL(10, 2) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reservation_id) REFERENCES Reservation(reservation_id)
);

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










