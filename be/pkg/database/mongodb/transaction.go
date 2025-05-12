package mongodb

import (
	"app/be/internal/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Transaction struct {
	*Base
}

func NewTransRepo() *Transaction {
	return &Transaction{
		Base: &Base{
			collection: GetDB().Collection("transaction"),
		},
	}
}
func (ins *Transaction) AddTransaction(ctx context.Context, object interface{}) (
	string, error) {
	transactionId, err := ins.Base.Add(ctx, object)
	if err != nil {
		return "", err
	}
	return transactionId, nil
}

func (ins *Transaction) GetTransactionByAccountNumber(ctx context.Context, AccNumber string) (
	[]*models.Transaction, error) {
	var trans []*models.Transaction
	filter := bson.M{"$or": bson.A{
		bson.M{"senderAccountNumber": AccNumber},
		bson.M{"recipientAccountNumber": AccNumber},
	}}
	err := ins.Base.GetAll(ctx, filter, &trans)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}
	return trans, nil
}
