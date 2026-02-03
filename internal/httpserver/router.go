package httpserver

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine{

	// health-check 
	r:=gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/health",HealthCheck)

	return r

}