# Charlotte_CarSharingSystem
This project is a fully functional electric car-sharing system, designed to be efficient, scalable, and user-friendly. It consists of three distinct microservices, each handling specific aspects of the system: User Service, Vehicle Service, and Billing Service. These services communicate seamlessly through RESTful APIs, enabling a cohesive experience for users while maintaining modularity for future scalability and maintenance.

# Microservices Overview
## 1. User Service
The User Service manages all user-related functionality, including user registration, login, account management, and membership upgrades. Below are the key features:

### User Registration:
- Handles new user accounts.
- Passwords are hashed using bcrypt for secure storage.

### User Login:
- User-provided passwords are verified by comparing them to the hashed version in the database using bcrypt.
- Ensures secure authentication through encryption.

### Membership Tiers:
- Accounts are initialized with the "Basic" membership tier.
- Tiers can be upgraded to "Premium" or "VIP" based on monthly rental spending.
- Benefits of higher tiers include:
- Reduced hourly rental rates.
- Priority access to vehicles.
- Increased booking limits.

### Dashboard:
- Users can view their membership status, rental history, and other key details.
- Accessible after login.

### Account Management:
- Users can update personal information on the Profile Management page.

<br>

## 2. Vehicle Service
The Vehicle Service manages vehicle-related functionality, including the display of available vehicles and the reservation process. It ensures efficient usage of the vehicle fleet while maintaining fairness for all users.  Below are the key features:

### Vehicle Availability:
Displays a list of available vehicles with relevant details such as:
- Model
- Charge level
- Location
- Rental rate
- Mileage

### Reservation System:
Users can select a vehicle and proceed to make a reservation. The reservation process includes:
- Selecting a start and end date.
- Ensuring the reservation duration is between 1 hour and 3 days to meet system requirements.
- A straightforward Reserve button initiates the reservation.

### Fair Access:
- Restrictions on reservation duration ensure that all users have a fair opportunity to access vehicles.

<br>

## 3. Billing Service
The Billing Service handles payment processing and rental fee calculations for user reservations. It ensures accurate pricing and supports secure payment workflows. Below are the key features:

### Rental Fee Calculation:
Computes fees based on:
- Vehicle rental rate.
- Membership tier benefits.
- Reservation duration.

### Payment Processing:

<br>

# Instructions On Running Car System Microservices
## 1. Open 3 different terminals
- Terminal Type can be Shell/Git Bash
  
## 2. In each terminal, change directory to each micro services:
   ### user_service:
   - cd car_system/user_service
   - Runs on port 8080

   ### vehicle_service:
   - cd car_system/vehicle_service
   - Runs on port 8081
  
   ### billing_service:
   - cd car_system/billing_service
   - Runs on port 8082

<br> 

# Architecture Diagram of Car Rental System
  ![image](https://github.com/user-attachments/assets/9b49281d-9ea3-443a-9671-b35238250a3a)


