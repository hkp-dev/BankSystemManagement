package transaction

import (
	"app/be/internal/hub"
	"app/be/internal/models"
	"context"
)

type TransactionService struct {
	appHub *hub.App
}

func NewTransactionService(appHub *hub.App) *TransactionService {
	return &TransactionService{
		appHub: appHub,
	}
}

func (ins *TransactionService) GetAllTransaction(ctx context.Context, identifier string) (
	trans []*models.Transaction, err error) {
	trans, err = ins.appHub.DbcTrans.GetTransactionByUserId(ctx, identifier)
	if err != nil {
		return nil, err
	}
	return trans, nil
}

// func (ins *TransactionService) AddTransaction(ctx context.Context, )
