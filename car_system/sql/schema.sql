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

-- Sample Data
-- Membership Data
INSERT INTO Membership (membership_tier, hourly_rate_discount, priority_access, booking_limit)
VALUES
('Basic', 0.00, FALSE, 5),
('Premium', 10.00, TRUE, 10),
('VIP', 20.00, TRUE, 20);


-- User Data
INSERT INTO User (name, email, phone_no, password, dob, membership_tier)
VALUES
('John Doe', 'john.doe@example.com', '1234567890', 'hashed_password_1', '1990-01-15', 'Basic'),
('Jane Smith', 'jane.smith@example.com', '0987654321', 'hashed_password_2', '1985-06-25', 'Premium'),
('Robert Brown', 'robert.brown@example.com', '1122334455', 'hashed_password_3', '1995-11-10', 'VIP'),
('Emily Davis', 'emily.davis@example.com', '5566778899', 'hashed_password_4', '1992-03-05', 'Basic');

-- Rental History Data
INSERT INTO Rental_History (user_id, vehicle_id, start_time, end_time, cost, status)
VALUES
(1, 101, '2024-12-01 08:00:00', '2024-12-01 12:00:00', 50.00, 'Completed'),
(2, 102, '2024-12-02 10:00:00', '2024-12-02 14:00:00', 80.00, 'Completed'),
(3, 103, '2024-12-03 09:00:00', '2024-12-03 17:00:00', 120.00, 'Refunded'),
(1, 104, '2024-12-04 15:00:00', '2024-12-04 18:00:00', 75.00, 'Cancelled');

-- Driver License Data
INSERT INTO DriverLicense (license_no, license_issue_date, license_expiry_date)
VALUES
('DL12345678', '2020-01-15', '2030-01-15'),
('DL87654321', '2018-06-10', '2028-06-10'),
('DL45678901', '2019-03-20', '2029-03-20'),
('DL23456789', '2021-08-05', '2031-08-05');

-- Select Statements
SELECT * FROM DriverLicense;
SELECT * FROM Membership;
SELECT * FROM User;
SELECT * FROM Rental_History;

--================================================================================================================
-- VEHICLE SERVICE -- 

DROP DATABASE IF EXISTS vehicle_service;
CREATE DATABASE vehicle_service;
USE vehicle_service;
-- Vehicle Table

CREATE TABLE Vehicle (
    vehicle_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    license_plate VARCHAR(20) UNIQUE NOT NULL,
    model VARCHAR(100) NOT NULL,
    charge_level DECIMAL(5, 2) NOT NULL,
    location VARCHAR(255) NOT NULL,
    rental_rate DECIMAL(10, 2) NOT NULL,
    mileage INT NOT NULL,
    status ENUM('Operational', 'Decommissioned', 'Under Maintenance') DEFAULT 'Operational',
    battery_capacity_kwh DECIMAL(5, 2) DEFAULT NULL,
    reservation_status VARCHAR(50) DEFAULT 'Available'
);

-- Reservation Table
CREATE TABLE Reservation (
    reservation_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    vehicle_id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NOT NULL, -- Use an unsigned integer for consistency
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    expected_charge_level DECIMAL(5, 2) NOT NULL,
    status ENUM('Active', 'Completed', 'Cancelled') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vehicle_id) REFERENCES Vehicle(vehicle_id),
    CHECK (start_time < end_time)
);

-- Rental Table
CREATE TABLE Rental (
    rental_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    reservation_id INT UNSIGNED NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    rental_fee DECIMAL(10, 2) NOT NULL,
    payment_status ENUM('Pending', 'Paid', 'Refunded') NOT NULL,
    payment_amount DECIMAL(10, 2) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reservation_id) REFERENCES Reservation(reservation_id)
);

-- Sample Data
-- Vehicle Data
INSERT INTO Vehicle (license_plate, model, charge_level, location, rental_rate, mileage, status, battery_capacity_kwh)
VALUES
('ABC123', 'Tesla Model 3', 80.00, 'Downtown Station', 50.00, 12000, 'Operational', 75.00),
('XYZ789', 'Nissan Leaf', 90.00, 'Airport Terminal', 40.00, 15000, 'Operational', 62.00),
('JKL456', 'Chevrolet Bolt', 60.00, 'Suburban Hub', 45.00, 18000, 'Operational', 65.00),
('DEF321', 'Hyundai Kona Electric', 50.00, 'City Center', 55.00, 20000, 'Operational', 64.00),
('GHI654', 'BMW i3', 70.00, 'Train Station', 60.00, 22000, 'Operational', 42.00);

-- Reservation Data
INSERT INTO Reservation (vehicle_id, user_id, start_time, end_time, expected_charge_level, status)
VALUES
(1, 1, '2024-12-10 08:00:00', '2024-12-10 12:00:00', 80.00, 'Active'),
(2, 1, '2024-12-11 09:00:00', '2024-12-11 15:00:00', 90.00, 'Completed'),
(3, 2, '2024-12-12 10:00:00', '2024-12-12 14:00:00', 75.00, 'Active'),
(4, 3, '2024-12-13 11:00:00', '2024-12-13 16:00:00', 50.00, 'Cancelled'),
(5, 4, '2024-12-14 07:00:00', '2024-12-14 10:00:00', 85.00, 'Completed');

-- Rental Data
INSERT INTO Rental (reservation_id, start_date, end_date, rental_fee, payment_status, payment_amount)
VALUES
(1, '2024-12-10 08:00:00', '2024-12-10 12:00:00', 200.00, 'Paid', 200.00),
(1, '2024-12-11 09:00:00', '2024-12-11 15:00:00', 240.00, 'Paid', 240.00),
(3, '2024-12-12 10:00:00', '2024-12-12 14:00:00', 180.00, 'Pending', NULL),
(4, '2024-12-13 11:00:00', '2024-12-13 16:00:00', 275.00, 'Refunded', 275.00),
(5, '2024-12-14 07:00:00', '2024-12-14 10:00:00', 165.00, 'Paid', 165.00);

-- Select Statements
SELECT * FROM Vehicle;
SELECT * FROM Reservation;
SELECT * FROM Rental;

--================================================================================================================
-- BILLING SERVICE -- 

DROP DATABASE IF EXISTS billing_service;
CREATE DATABASE billing_service;
USE billing_service;

-- Billing Table
CREATE TABLE Billing (
    bill_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,           
    reservation_id INT NOT NULL,   
    promo_id INT DEFAULT NULL,     
    amount DECIMAL(10, 2) NOT NULL,
    status ENUM('Pending', 'Paid', 'Refunded') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Promotion Table
CREATE TABLE Promotion (
    promo_id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(255) NOT NULL,
    discount_rate DECIMAL(5, 2) NOT NULL,
    valid_from TIMESTAMP NOT NULL,
    valid_to TIMESTAMP NOT NULL
);

-- Sample Data
-- Promotion Data
INSERT INTO Promotion (code, description, discount_rate, valid_from, valid_to)
VALUES
('HOLIDAY10', '10% off during holiday season', 10.00, '2024-12-01 00:00:00', '2024-12-31 23:59:59'),
('NEWYEAR20', '20% off for New Year', 20.00, '2024-12-25 00:00:00', '2025-01-05 23:59:59'),
('WEEKEND5', '5% off on weekends', 5.00, '2024-01-01 00:00:00', '2024-12-31 23:59:59'),
('VIP25', '25% discount for VIP members', 25.00, '2024-01-01 00:00:00', '2024-12-31 23:59:59'),
('SUMMER15', '15% off during summer', 15.00, '2024-06-01 00:00:00', '2024-08-31 23:59:59');

-- Billing Data
INSERT INTO Billing (user_id, reservation_id, promo_id, amount, status)
VALUES
(101, 1, 1, 180.00, 'Paid'),
(102, 2, 2, 192.00, 'Paid'),
(103, 3, NULL, 200.00, 'Pending'),
(104, 4, 3, 261.25, 'Refunded'),
(105, 5, 4, 123.75, 'Paid');

-- Select Statements
SELECT * FROM Billing;
SELECT * FROM Promotion;









