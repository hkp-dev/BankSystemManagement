package database

import (
	"app/be/internal/models"
	"app/be/pkg/database/mongodb"
	"context"
)

func NewDB(ctx context.Context) (hub *models.DbHub, err error) {
	if err := mongodb.Init(ctx); err != nil {
		panic(err)
	}
	return &models.DbHub{
		DbcUser:  mongodb.NewUserRepo(),
		DbcTrans: mongodb.NewTransRepo(),
	}, nil
}
