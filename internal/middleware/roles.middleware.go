package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Admin-level restrictions MW..
// Check if logged-in/authenticated user is ADMIN or not

func RequireAdmin()gin.HandlerFunc{
	return func(ctx *gin.Context) {
		role,ok:= GetUserRole(ctx)

		if !ok{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":"Unauthorized!",
			})
			return 
		}

		if !strings.EqualFold(role,"admin"){
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":"FORBIDDEN.. This route can only be accessed by an ADMIN!",
			})
			return 
		}
		ctx.Next()
	}
}