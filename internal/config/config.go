package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI   string
	MongoDBName string
	JWTSecret  string
}

func Load() (Config, error) {
	_ = godotenv.Load()
	
	cfg:=Config{
		MongoURI: strings.TrimSpace(os.Getenv("MONGODB_URI")),
		MongoDBName: strings.TrimSpace(os.Getenv("MONGODB_NAME")),
		JWTSecret: strings.TrimSpace(os.Getenv("JWT_SECRET")),
	}

	// Checks
	if cfg.MongoURI == ""{
		return Config{},fmt.Errorf("⚠️ Missing MONGO-URI!")
	}

	if cfg.MongoDBName == ""{
		return Config{},fmt.Errorf("⚠️ Missing MONGODB-NAME!")
	}

	if cfg.JWTSecret == ""{
		return Config{},fmt.Errorf("⚠️ Missing JWT-SECRET!")
	}

	return cfg,nil

}