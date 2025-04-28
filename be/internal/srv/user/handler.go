package user

import (
	"app/be/internal/utility"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	utility.MapCor(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log payload size
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	log.Printf("Payload size: %d bytes", len(body))
	r.Body = io.NopCloser(bytes.NewReader(body))

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Create data/img directory if it doesn't exist
	imgDir := "data/img"
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		log.Printf("Error creating directory %s: %v", imgDir, err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Save ID card images
	frontPath := ""
	backPath := ""
	if req.IDCardFront != "" {
		frontPath, err = utility.SaveImage(req.IDCardFront, imgDir, req.Identifier+"_front")
		if err != nil {
			log.Printf("Error saving front ID image: %v", err)
			http.Error(w, "Error saving front ID image", http.StatusBadRequest)
			return
		}
	}
	if req.IDCardBack != "" {
		backPath, err = utility.SaveImage(req.IDCardBack, imgDir, req.Identifier+"_back")
		if err != nil {
			log.Printf("Error saving back ID image: %v", err)
			http.Error(w, "Error saving back ID image", http.StatusBadRequest)
			return
		}
	}

	log.Printf("Processing registration for username: %s, email: %s", req.Identifier, req.Email)

	// Add 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	id, err := h.authService.Register(ctx,
		req.Identifier, req.Email, req.Password,
		req.FirstName, req.LastName, req.PhoneNumber,
		req.Address, req.City, req.ZipCode,
		frontPath, backPath)

	if err != nil {
		log.Printf("Registration failed for username %s: %v", req.Identifier, err)
		// Delete images if registration fails
		if frontPath != "" {
			os.Remove(frontPath)
		}
		if backPath != "" {
			os.Remove(backPath)
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Successfully registered user %s with ID: %s", req.Identifier, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "User registered successfully, awaiting admin approval",
	})
}

func (ins *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	utility.MapCor(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userId, is2FAEnabled, err := ins.authService.LoginStep1(r.Context(), req.Identifier, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"userId":       userId,
		"is2FAEnabled": is2FAEnabled,
	})
}

func (h *AuthHandler) GetTOTPSetup(w http.ResponseWriter, r *http.Request) {
	utility.MapCor(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	qrCodeUrl, secret, err := h.authService.GetTOTPSetup(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"qrCodeUrl": qrCodeUrl,
		"secret":    secret,
	})
}

func (h *AuthHandler) VerifyFirstTimeOTP(w http.ResponseWriter, r *http.Request) {
	utility.MapCor(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID    string `json:"userId"`
		OTP       string `json:"otp"`
		SecretKey string `json:"secKey"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.authService.VerifyFirstTimeOTP(r.Context(), req.UserID, req.OTP, req.SecretKey); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "2FA setup success",
	})
}

func (h *AuthHandler) VerifyLoginOTP(w http.ResponseWriter, r *http.Request) {
	utility.MapCor(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string `json:"userId"`
		OTP    string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.authService.VerifyLoginOTP(r.Context(), req.UserID, req.OTP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set Access Token
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   15 * 60, // 15 minutes
		SameSite: http.SameSiteStrictMode,
	})

	// Set Refresh Token
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   7 * 24 * 60 * 60,
		SameSite: http.SameSiteStrictMode,
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login success",
	})
}

func (c *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
