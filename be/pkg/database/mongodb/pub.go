package mongodb

import (
	"app/be/internal/models"
	"context"

	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gopkg.in/yaml.v3"
)

var (
	currDb *mongo.Database
)

func Init(ctx context.Context) (err error) {
	data, err := os.ReadFile("../cfg.yml")
	if err != nil {
		log.Print("[Mongo] Read cfg.yml failed: ", err)
	}
	var cfg models.MongoCfg
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		log.Print("[Mongo] Unmarshal cfg.yml failed: ", err)
		return
	}
	opts := options.Client().ApplyURI(cfg.Mongo.URI)
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Print("[Mongo] Connect to MongoDB failed: ", err)
		return
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Print("[Mongo] Ping MongoDB failed: ", err)
		return
	}
	log.Print("[Mongo] Ping MongoDB success")
	currDb = client.Database(cfg.Mongo.DB)
	log.Print("[Mongo] Connected to MongoDB: ", cfg.Mongo.URI)
	return
}

func GetDB() *mongo.Database {
	return currDb
}
