package middleware

import (
	"net/http"
	"strings"

	"github.com/callmeskyy111/golang-jwt-auth/internal/auth"
	"github.com/gin-gonic/gin"
)

// store -> auth data-info -> gin ctxt.

const (
	ctxUserIDkey ="auth.userId"
	ctxRolekey = "auth.role"
)

func AuthRequired(jwtSecret string)gin.HandlerFunc{
	return func(ctx *gin.Context) {
		authHeader := strings.TrimSpace(ctx.GetHeader("Authorization"))
		if authHeader==""{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":"⚠️ Missing Authorization-Token!",
			})
		return 	
		}

		parts:=strings.SplitN(authHeader, " ",2)
		if len(parts) !=2{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":"⚠️ Invalid Authorization-format!",
			})
		return 	
		}

		scheme:=strings.TrimSpace(parts[0])
		tokenStr:=strings.TrimSpace(parts[1])

		if !strings.EqualFold(scheme, "Bearer"){
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":"⚠️ Authorization-scheme must be BEARER!",
			})
		return
		}

		if tokenStr==""{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":"⚠️ Missing token here!",
			})
		return
		}

		claims,err:=auth.ParseToken(jwtSecret, tokenStr) // jwt.go

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":"⚠️ Invalid/Expired token!",
			})
		return
		}
		ctx.Set(ctxUserIDkey, claims.Subject)
		ctx.Set(ctxRolekey, claims.Role)

		// Next MW..
		ctx.Next()
	}
}

func GetUserID(ctx *gin.Context)(string,bool){
	result,ok:=ctx.Get(ctxUserIDkey)
	if !ok{
		return "",false
	}

	userID, ok:= result.(string)

	return userID,ok
}

func GetUserRole(ctx *gin.Context)(string,bool){
	result,ok:=ctx.Get(ctxRolekey)
	if !ok{
		return "",false
	}

	role, ok:= result.(string)

	return role,ok
}
