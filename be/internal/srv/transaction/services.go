package transaction

import (
	"app/be/internal/hub"
	"app/be/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pquerna/otp/totp"
)

type TransactionService struct {
	appHub *hub.App
}

func NewTransactionService(appHub *hub.App) *TransactionService {
	return &TransactionService{
		appHub: appHub,
	}
}
func (ins *TransactionService) GetTransactions(ctx context.Context, userId string) (
	trans []*models.Transaction, total int, err error) {
	user, err := ins.appHub.DbcUser.GetUserById(ctx, userId)
	if err != nil {
		log.Printf("userId %s not found: %v", userId, err)
		return
	}
	trans, err = ins.appHub.DbcTrans.GetTransactionByAccountNumber(ctx, user.AccountNumber)
	if err != nil {
		log.Printf("get transaction of user %s failed: %v", user.ID, err)
		return
	}

	total = len(trans)
	return
}

func (ins *TransactionService) PostTransaction(ctx context.Context, userId, accountNumberSend,
	recipientAccountNumber string, amount float64, description, otp string) (transId string, createTime string, err error) {
	user, err := ins.appHub.DbcUser.GetUserByAccNumber(ctx, accountNumberSend)
	if err != nil {
		return "", "", fmt.Errorf("user not found: %w", err)
	}

	if user.TOTPSecret == "" {
		return "", "", fmt.Errorf("2FA not setup for user")
	}

	if !totp.Validate(otp, user.TOTPSecret) {
		return "", "", fmt.Errorf("invalid OTP code")
	}
	if user.Balance < amount {
		return "", "", fmt.Errorf("insufficient balance")
	}
	reciptUser, err := ins.appHub.DbcUser.GetUserByAccNumber(ctx, recipientAccountNumber)
	if err != nil {
		return "", "", fmt.Errorf("recipient user not found: %w", err)
	}
	balanceUser := user.Balance - amount
	balanceReciptUser := reciptUser.Balance + amount
	err = ins.appHub.DbcUser.UpdateBalance(ctx, userId, balanceUser)
	if err != nil {
		return "", "", fmt.Errorf("failed to update balance of user send: %w", err)
	}
	err = ins.appHub.DbcUser.UpdateBalance(ctx, reciptUser.ID.Hex(), balanceReciptUser)
	if err != nil {
		return "", "", fmt.Errorf("failed to update balance of user recipt: %w", err)
	}
	trans := &models.Transaction{
		SenderAccountNumber:    accountNumberSend,
		RecipientAccountNumber: recipientAccountNumber,
		Amount:                 amount,
		Description:            description,
		CreateAt:               time.Now(),
	}
	transId, err = ins.appHub.DbcTrans.AddTransaction(ctx, trans)
	if err != nil {
		return "", "", fmt.Errorf("failed to add transaction: %w", err)
	}
	createTime = trans.CreateAt.Format(time.RFC3339)
	return
}
