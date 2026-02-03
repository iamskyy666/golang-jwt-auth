package app

import (
	"context"
	"fmt"
	"time"

	"github.com/callmeskyy111/golang-jwt-auth/internal/config"
	"github.com/callmeskyy111/golang-jwt-auth/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Config config.Config
	MongoClient *mongo.Client
	DB *mongo.Database
}

// Now, join all the individual parts - .env-loading, connection, setup, server, etc.
func NewApp(ctx context.Context)(*App, error){

	// Load env
	cfg, err:=config.Load()

	if err != nil {
		return nil, err
	}

	// Connect to DB
	mongoConnection, err:=db.ConnectDB(ctx,cfg)

	if err != nil {
		return nil, err
	}

	// Finally..
	return &App{
		Config: cfg,
		MongoClient: mongoConnection.Client,
		DB: mongoConnection.DB,
	},nil
}

func (a *App) CloseMongo(ctx context.Context)error{
	if a.MongoClient==nil{
		return nil
	}

	closeCtx,cancel:=context.WithTimeout(ctx,time.Second * 5)
	defer cancel()

	if err:=a.MongoClient.Disconnect(closeCtx); err!=nil{
		return fmt.Errorf("⚠️ Failed to disconnect mongo: %w",err)
	}

	return nil
}