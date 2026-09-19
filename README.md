# Real Nest Auth 🔐

A secure, production-grade Full-Stack Authentication System built with **Go (Gin framework)**, **MongoDB**, and a modern **Glassmorphic Web UI**.

Featuring JSON Web Tokens (JWT), Bcrypt password hashing, IP-based brute-force protection, and an automated OTP password reset workflow integrated with the **Brevo (Sendinblue) Transactional Email API**.

---

## 📑 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Project Architecture](#-project-architecture)
- [Security Architecture](#-security-architecture)
- [API Documentation](#-api-documentation)
  - [1. User Registration](#1-user-registration)
  - [2. User Login](#2-user-login)
  - [3. Forgot Password (OTP Request)](#3-forgot-password-otp-request)
  - [4. Reset Password (OTP Verification)](#4-reset-password-otp-verification)
  - [5. User Logout](#5-user-logout)
- [Environment Variables](#-environment-variables)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation & Setup](#installation--setup)
  - [Running the Frontend](#running-the-frontend)
- [Frontend Overview](#-frontend-overview)
- [License](#-license)

---

## ✨ Features

- **Robust Authentication**:
  - User signup with regex email validation and duplicate checks.
  - User login issuing 24-hour signed **JWT tokens**.
  - Password encryption using **bcrypt** with a secure work factor of 14.
- **Brute-Force & Rate-Limiting Protection**:
  - Custom in-memory IP rate limiter on login (`LoginLimiter`): blocks IP addresses for 2 minutes after 5 consecutive attempts.
  - Password reset OTP throttling: 1-minute cooldown between requests and a hard limit of 3 requests per hour.
- **Self-Service Password Recovery (OTP-Based)**:
  - Cryptographically secure 6-digit OTP generation using Go's `crypto/rand`.
  - OTP hashed with bcrypt before saving to MongoDB (no plain OTPs stored in database).
  - 5-minute OTP expiration window with atomic one-time consumption (cleared immediately upon reset).
  - Transactional email dispatch powered by **Brevo REST API**.
- **Interactive Modern UI**:
  - Full-screen ambient video background with dark overlay.
  - Sleek glassmorphism UI card with smooth backdrop-filter blur effects.
  - Seamless animated tab switching between Login, Registration, and Password Reset forms.
- **Production Ready**:
  - Structured Go codebase with clear separation of concerns (`handlers`, `middleware`, `models`, `db`, `utils`).
  - CORS-enabled for cross-origin frontend communication.
  - Clean error handling with informative HTTP status codes.

---

## 🛠 Tech Stack

| Layer | Technologies Used |
|---|---|
| **Backend Framework** | [Go 1.25](https://golang.org/) & [Gin Web Framework](https://github.com/gin-gonic/gin) |
| **Database** | [MongoDB](https://www.mongodb.com/) (using official [mongo-go-driver](https://go.mongodb.org/mongo-driver)) |
| **Authentication & Tokens** | [Golang-JWT v5](https://github.com/golang-jwt/jwt) |
| **Cryptography** | [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) |
| **Email Service** | [Brevo REST API v3](https://developers.brevo.com/) (Transactional Email / SMTP) |
| **Configuration** | [godotenv](https://github.com/joho/godotenv) |
| **Frontend** | Vanilla HTML5, Modern CSS3 (Glassmorphism), JavaScript (Fetch API) |

---

## 📂 Project Architecture

```plaintext
Login_Page/
├── README.md                     # Project documentation
└── login-api/
    ├── main.go                   # Server initialization, middleware & routes
    ├── go.mod                    # Go module dependencies
    ├── go.sum                    # Checksums for Go dependencies
    ├── .env.example              # Template for environment variables
    ├── .gitignore                # Git ignore rules
    ├── auth-frontend/
    │   └── index.html            # Single-page glassmorphism frontend application
    ├── db/
    │   ├── mongo.go              # MongoDB client setup & connection pooling
    │   └── user-db.go            # User database operations (queries, OTP tracking)
    ├── handlers/
    │   ├── auth.go               # Login & Signup handlers
    │   ├── forgot_password.go    # OTP generation & Brevo email dispatch
    │   ├── reset_password.go     # OTP verification & password update
    │   └── logout.go             # Session termination handler
    ├── middleware/
    │   └── login_limit.go        # In-memory IP rate-limiting middleware
    ├── models/
    │   └── user.go               # MongoDB BSON & JSON User model definitions
    ├── routes/
    │   └── routes.go             # Route definition blueprints
    └── utils/
        ├── brevo_email.go        # Brevo transactional API client
        ├── hash.go               # Bcrypt password hashing & validation helpers
        ├── jwt.go                # JWT token creation & signing
        └── validation.go         # Regex-based email validation helper
```

---

## 🛡 Security Architecture

1. **Password Hashing**: Passwords are never stored in plaintext. They are hashed using `bcrypt` with cost `14`.
2. **Hashed OTP Storage**: When an OTP is issued, it is hashed via bcrypt before writing to MongoDB (`reset_otp_hash`). Even if the database is accessed, plaintext OTPs cannot be extracted.
3. **Short Expiration Windows**: OTP codes expire strictly after 5 minutes.
4. **Anti-Spam & Cooldown**: Users must wait at least 60 seconds between OTP requests, and are capped at 3 requests per hour.
5. **Single-Use OTP**: Successful password reset executes `db.ClearOTP()`, invalidating the token immediately.
6. **Login Throttling**: IP addresses exceeding 5 login attempts within a short period are blocked for 2 minutes with HTTP `429 Too Many Requests`.

---

## 📡 API Documentation

Base URL (Local): `http://localhost:8080`  
Base URL (Production Demo): `https://login-api-m00t.onrender.com`

### 1. User Registration

Creates a new user account with hashed credentials.

- **Endpoint**: `POST /signup`
- **Headers**: `Content-Type: application/json`
- **Request Body**:
  ```json
  {
    "name": "John",
    "surname": "Doe",
    "email": "john.doe@example.com",
    "password": "SecurePassword123!",
    "phone_number": "+1234567890"
  }
  ```
- **Responses**:
  - `201 Created`:
    ```json
    {
      "message": "Signup successful"
    }
    ```
  - `400 Bad Request`: `{"error": "Invalid email format"}` or `{"error": "Invalid request"}`
  - `409 Conflict`: `{"error": "User already exists"}`

---

### 2. User Login

Validates user credentials, applies IP rate-limiting, and returns a signed JWT.

- **Endpoint**: `POST /login`
- **Rate Limit**: Max 5 attempts per IP per 2 minutes
- **Headers**: `Content-Type: application/json`
- **Request Body**:
  ```json
  {
    "email": "john.doe@example.com",
    "password": "SecurePassword123!"
  }
  ```
- **Responses**:
  - `200 OK`:
    ```json
    {
      "message": "Login successful",
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```
  - `401 Unauthorized`: `{"error": "User not found"}` or `{"error": "Wrong password", "action": "forgot-password"}`
  - `429 Too Many Requests`: `{"error": "Too many login attempts. Try again later."}`

---

### 3. Forgot Password (OTP Request)

Generates a secure 6-digit OTP, stores its bcrypt hash in MongoDB, and dispatches it via email.

- **Endpoint**: `POST /forgot-password`
- **Rate Limit**: Minimum 1 minute cooldown, max 3 OTP requests / hour
- **Headers**: `Content-Type: application/json`
- **Request Body**:
  ```json
  {
    "email": "john.doe@example.com"
  }
  ```
- **Responses**:
  - `200 OK`:
    ```json
    {
      "message": "OTP sent to your email (valid for 5 minutes)"
    }
    ```
  - `404 Not Found`: `{"error": "User not found"}`
  - `429 Too Many Requests`: `{"error": "Wait 1 minute before requesting OTP again"}` or `{"error": "OTP limit reached. Try after 1 hour"}`

---

### 4. Reset Password (OTP Verification)

Validates the submitted OTP against the stored hash and updates the account password.

- **Endpoint**: `POST /reset-password`
- **Headers**: `Content-Type: application/json`
- **Request Body**:
  ```json
  {
    "email": "john.doe@example.com",
    "otp": "123456",
    "new_password": "NewSecurePassword456!"
  }
  ```
- **Requirements**:
  - `new_password` must be at least 8 characters long.
  - OTP must match and not be expired (within 5 minutes of request).
- **Responses**:
  - `200 OK`:
    ```json
    {
      "message": "Password reset successful"
    }
    ```
  - `400 Bad Request`: `{"error": "Invalid OTP"}` or `{"error": "OTP expired"}` or `{"error": "Password must be at least 8 characters long"}`
  - `404 Not Found`: `{"error": "User not found"}`

---

### 5. User Logout

Terminates the user's active session.

- **Endpoint**: `POST /logout`
- **Headers**: `Content-Type: application/json`
- **Responses**:
  - `200 OK`:
    ```json
    {
      "message": "Logout successful"
    }
    ```

---

## ⚙️ Environment Variables

Create a `.env` file in the `login-api/` directory by copying `.env.example`:

```bash
cp login-api/.env.example login-api/.env
```

| Variable | Description | Default / Example | Required |
|---|---|---|---|
| `PORT` | Port the Gin HTTP server listens on | `8080` | No |
| `MONGODB_URI` | MongoDB connection URI (local or MongoDB Atlas cluster) | `mongodb+srv://user:pass@cluster0.mongodb.net/` | **Yes** |
| `BREVO_API_KEY` | API Key from Brevo dashboard for sending transactional emails | `xkeysib-...` | **Yes** (for OTP) |
| `BREVO_FROM_EMAIL` | Verified sender email address configured in Brevo | `noreply@yourdomain.com` | **Yes** (for OTP) |
| `BREVO_FROM_NAME` | Sender display name shown in email client | `Real Nest Auth` | No |

---

## 🚀 Getting Started

### Prerequisites

Make sure you have the following installed on your machine:
- [Go](https://go.dev/dl/) (version 1.21 or higher, project uses Go 1.25)
- [MongoDB](https://www.mongodb.com/) (running locally or a free [MongoDB Atlas](https://www.mongodb.com/atlas) cloud cluster)
- A [Brevo](https://www.brevo.com/) account with an active API Key and verified sender email

### Installation & Setup

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/ujjwal0563/Login_Page.git
   cd Login_Page/login-api
   ```

2. **Configure Environment Variables**:
   ```bash
   cp .env.example .env
   ```
   Open `.env` and fill in your MongoDB connection string and Brevo API credentials.

3. **Install Go Dependencies**:
   ```bash
   go mod download
   ```

4. **Run the API Server**:
   ```bash
   go run main.go
   ```
   The server will start at `http://localhost:8080`. You should see:
   ```text
   MongoDB connected successfully
   [GIN-debug] Listening and serving HTTP on :8080
   ```

---

### Running the Frontend

The frontend is a self-contained single-page web client located in `login-api/auth-frontend/index.html`.

You can run it using any local static file server or simply open the file in your browser:

#### Option A: Using Python (Quickest)
```bash
cd login-api/auth-frontend
python3 -m http.server 3000
```
Then visit `http://localhost:3000` in your web browser.

#### Option B: Using VS Code / IDE Live Server
Right-click `login-api/auth-frontend/index.html` and select **"Open with Live Server"**.

> 💡 **Tip for Local API Testing**:
> By default, `index.html` connects to the deployed Render instance. To connect to your local backend, update line 243 of `login-api/auth-frontend/index.html`:
> ```javascript
> const API_BASE_URL = "http://localhost:8080";
> ```

---

## 🎨 Frontend Overview

The web client includes:
- **Responsive Glassmorphic Design**: Frosted glass container (`backdrop-filter: blur(15px)`), translucent gradients, and custom modern typography (`Poppins`).
- **Dynamic Video Background**: High-definition looping network data visualization background video.
- **Feedback Notifications**: Floating status/error toast banner for instant user feedback.
- **Client-Side View Switching**: Seamlessly toggle between:
  - `Login`: Simple email and password form with rate-limit handling.
  - `Register`: Name, username, phone, email, and password registration.
  - `Reset Password`: 2-step OTP verification and password update flow.

---

## 📄 License

This project is licensed under the MIT License — feel free to modify and use it for your personal or commercial applications.
