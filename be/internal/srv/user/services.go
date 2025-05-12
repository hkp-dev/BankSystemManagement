package user

import (
	"app/be/internal/hub"
	"app/be/internal/models"
	"context"

	// "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	appHub *hub.App
}

func NewAuthService(appHub *hub.App) *AuthService {
	return &AuthService{appHub: appHub}
}
func (ins *AuthService) Register(ctx context.Context,
	identifier, email,
	password, firstName, lastName, phoneNumber,
	address, city, zipCode, idCardFront, idCardBack string) (
	id string, err error) {
	// Validate required fields
	if identifier == "" || email == "" || password == "" {
		log.Printf("Validation failed: identifier, email, or password missing")
		return "", fmt.Errorf("username, email, and password are required")
	}

	// Check existing username
	existingUser, err := ins.appHub.DbcUser.GetUserByIdentifier(ctx, identifier)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("No existing user found for username: %s", identifier)
		} else {
			log.Printf("Error checking existing username: %v", err)
			return "", fmt.Errorf("error checking existing username: %v", err)
		}
	} else if existingUser != nil {
		log.Printf("Username already exists: %s", identifier)
		return "", fmt.Errorf("username already exists")
	}

	// Check existing email
	existingEmailUser, err := ins.appHub.DbcUser.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("No existing user found for email: %s", email)
		} else {
			log.Printf("Error checking existing email: %v", err)
			return "", fmt.Errorf("error checking existing email: %v", err)
		}
	} else if existingEmailUser != nil {
		log.Printf("Email already exists: %s", email)
		return "", fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", fmt.Errorf("error hashing password: %v", err)
	}

	// Generate a random 10-digit account number
	rand.Seed(time.Now().UnixNano())
	accountNumber := fmt.Sprintf("%010d", rand.Intn(1e10))

	// Create user model
	user := &models.User{
		Identifier:    identifier,
		Email:         email,
		PasswordHash:  string(hashedPassword),
		FirstName:     firstName,
		LastName:      lastName,
		PhoneNumber:   phoneNumber,
		Address:       address,
		City:          city,
		ZipCode:       zipCode,
		IsVerified:    false,
		Status:        models.StatusActive,
		Role:          models.UserRole,
		AccountNumber: accountNumber,
		Balance:       0,
		IDCardFront:   idCardFront,
		IDCardBack:    idCardBack,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		TOTPSecret:    "",
	}

	id, err = ins.appHub.DbcUser.AddUser(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return "", fmt.Errorf("error creating user: %v", err)
	}

	log.Printf("Successfully created user with ID: %s", id)
	return id, nil
}

func (ins *AuthService) LoginStep1(ctx context.Context, identifier, password string) (id string, is2FAEnabled bool, err error) {
	user, err := ins.appHub.DbcUser.GetUserByIdentifier(ctx, identifier)
	if err != nil {
		return "", false, fmt.Errorf("user not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", false, errors.New("invalid password")
	}

	// if !user.IsVerified || user.Status != models.StatusActive {
	// 	return "", false, errors.New("account is not active or pending admin approval")
	// }

	return user.ID.Hex(), user.TOTPSecret != "", nil
}

func (ins *AuthService) GetTOTPSetup(ctx context.Context, userID string) (qrCodeURL, secret string, err error) {
	user, err := ins.appHub.DbcUser.GetUserById(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("user not found: %w", err)
	}

	if user.TOTPSecret != "" {
		return "", "", fmt.Errorf("2FA already enabled")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "BankingSystem",
		AccountName: user.FirstName + user.LastName,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	qrBytes, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate QR code: %w", err)
	}
	qrCodeBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrBytes)

	return qrCodeBase64, key.Secret(), nil
}

func (ins *AuthService) VerifyFirstTimeOTP(ctx context.Context, userID, otp, secret string) error {
	if secret == "" {
		return fmt.Errorf("missing secret")
	}

	if !totp.Validate(otp, secret) {
		return fmt.Errorf("invalid OTP code")
	}

	err := ins.appHub.DbcUser.UpdateTOTPSecret(ctx, userID, secret)
	if err != nil {
		return fmt.Errorf("failed to save 2FA secret: %w", err)
	}

	return nil
}

func (ins *AuthService) VerifyLoginOTP(ctx context.Context, userID, otp string) (
	userId string, err error) {
	user, err := ins.appHub.DbcUser.GetUserById(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	if user.TOTPSecret == "" {
		return "", fmt.Errorf("2FA not setup for user")
	}

	if !totp.Validate(otp, user.TOTPSecret) {
		return "", fmt.Errorf("invalid OTP code")
	}
	userId = user.ID.Hex()
	return
}

func (ins *AuthService) GetUsrInfo(ctx context.Context, userId string) (fullName, email, accountNumber string, balance float64, err error) {
	user, err := ins.appHub.DbcUser.GetUserById(ctx, userId)
	if err != nil {
		log.Print("userId %w not found: %w", userId, err)
		return "", "", "", 0, fmt.Errorf("user not found: %w", err)
	}
	fullName = user.FirstName + " " + user.LastName
	email = user.Email
	accountNumber = user.AccountNumber
	balance = user.Balance
	return
}
