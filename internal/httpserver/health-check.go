package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HealthCheck(ctx *gin.Context){
	ctx.JSON(http.StatusOK, gin.H{
		"ok":true,
		"service":"go-auth 🛡️",
		"time":time.Now().UTC(),
	})
}