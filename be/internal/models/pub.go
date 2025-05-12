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
	Balance       float64       `bson:"balance" json:"balance"`
	AccountNumber string        `bson:"accountNumber" json:"accountNumber"`
	CreatedAt     time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time     `bson:"updatedAt" json:"updatedAt"`
	TOTPSecret    string        `bson:"totpSecret,omitempty" json:"-"`
}

type Transaction struct {
	ID                     bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SenderAccountNumber    string        `bson:"senderAccountNumber" json:"senderAccountNumber"`
	RecipientAccountNumber string        `bson:"recipientAccountNumber" json:"recipientAccountNumber"`
	Date                   time.Time     `bson:"date" json:"date"`
	Amount                 float64       `bson:"amount" json:"amount"`
	Description            string        `bson:"description,omitempty" json:"description,omitempty"`
	CreateAt               time.Time     `bson:"create_at" json:"createAt"`
}
