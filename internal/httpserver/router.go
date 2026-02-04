package httpserver

import (
	"github.com/callmeskyy111/golang-jwt-auth/internal/app"
	"github.com/callmeskyy111/golang-jwt-auth/internal/user"
	"github.com/gin-gonic/gin"
)

func NewRouter(a *app.App) *gin.Engine{

	// health-check 
	r:=gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/health",HealthCheck)

	userRepo:= user.NewRepo(a.DB)
	userSvc:=user.NewService(userRepo, a.Config.JWTSecret)
	userHandler:=user.NewHandler(userSvc)

	r.POST("/register",userHandler.RegisterUser)

	return r
}

