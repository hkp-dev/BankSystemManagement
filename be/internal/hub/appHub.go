package hub

import (
	"app/be/internal/models"
	"app/be/pkg/database"
	"context"
)

type App struct {
	*models.DbHub
}

var (
	currIns *App
)

func Current() *App {
	return currIns
}
func New(ctx context.Context) (hub *App, err error) {
	hub = &App{}
	hub.DbHub, err = database.NewDB(ctx)
	if err != nil {
		return nil, err
	}
	return hub, nil
}
