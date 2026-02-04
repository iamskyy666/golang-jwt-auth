package db

import (
	"context"
	"fmt"
	"time"

	"github.com/callmeskyy111/golang-jwt-auth/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	DB *mongo.Database
}

func ConnectDB(ctx context.Context, cfg config.Config)(*Mongo, error){
	connectCtx, cancel:=context.WithTimeout(ctx, 8 * time.Second)
	defer cancel()

	// Create client-options
	clientOptns:=options.Client().ApplyURI(cfg.MongoURI)

	client,err:=mongo.Connect(connectCtx, clientOptns)

	if err!=nil{
		return nil,fmt.Errorf("⚠️ Mongo-connection failed: %w",err)
	}

	database:=client.Database(cfg.MongoDBName)

	// return Mongo struct instance and err.
	return &Mongo{
		Client: client,
		DB: database,
	}, nil
}