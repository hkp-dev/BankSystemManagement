package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role string
type Status bool
type TransactionType string
type TransactionStatus string

const (
	TransferSent     TransactionType = "TRANSFER_SENT"
	TransferReceived TransactionType = "TRANSFER_RECEIVED"
)

const (
	StatusSuccess TransactionStatus = "SUCCESS"
	StatusPending TransactionStatus = "PENDING"
	StatusFailed  TransactionStatus = "FAILED"
)

var (
	UserRole       Role = "user"
	AdminRole      Role = "admin"
	SuperAdminRole Role = "super-admin"
	StatusActive   Status
	StatusLocked   Status
)

type User struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Identifier    string        `bson:"identifier" json:"identifier"`
	Email         string        `bson:"email" json:"email"`
	PasswordHash  string        `bson:"passwordHash" json:"-"`
	FirstName     string        `bson:"firstName" json:"firstName"`
	LastName      string        `bson:"lastName" json:"lastName"`
	PhoneNumber   string        `bson:"phoneNumber" json:"phoneNumber"`
	Address       string        `bson:"address" json:"address"`
	City          string        `bson:"city" json:"city"`
	ZipCode       string        `bson:"zipCode" json:"zipCode"`
	IDCardFront   string        `bson:"idCardFront" json:"idCardFront"`
	IDCardBack    string        `bson:"idCardBack" json:"idCardBack"`
	IsVerified    bool          `bson:"verified" json:"verified"`
	Status        Status        `bson:"status" json:"status"`
	Role          Role          `bson:"role" json:"role"`
	TermsAccepted bool          `bson:"termsAccepted" json:"termsAccepted"`
	CreatedAt     time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time     `bson:"updatedAt" json:"updatedAt"`

	TOTPSecret string `bson:"totpSecret,omitempty" json:"-"`
}

type Transaction struct {
	ID               bson.ObjectID     `bson:"_id" json:"id"`
	UserID           bson.ObjectID     `bson:"userId" json:"userId"`
	Date             time.Time         `bson:"date" json:"date"`
	Type             TransactionType   `bson:"type" json:"type"`
	Amount           float64           `bson:"amount" json:"amount"`
	RecipientAccount *string           `bson:"recipientAccount" json:"recipientAccount"`
	Description      *string           `bson:"description" json:"description"`
	Status           TransactionStatus `bson:"status" json:"status"`
	CreatedAt        time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time         `bson:"updatedAt" json:"updatedAt"`
}
