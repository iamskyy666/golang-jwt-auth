package httpserver

import (
	"net/http"

	"github.com/callmeskyy111/golang-jwt-auth/internal/app"
	"github.com/callmeskyy111/golang-jwt-auth/internal/middleware"
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

	// unauth/public routes
	r.POST("/register",userHandler.RegisterUser)
	r.POST("/login",userHandler.LoginUser)

	// auth/protected routes - RBA
	api:=r.Group("/api")

	api.Use(middleware.AuthRequired(a.Config.JWTSecret)) // authenticated ✅

	api.GET("/files", func(ctx *gin.Context) {

		userID,_:=middleware.GetUserID(ctx) // if we want the userID
		ctx.JSON(http.StatusOK, gin.H{
			"ok":true,
			"userID":userID,
			"files":[]any{},
		})
	})

	 api.GET("/products", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"ok":true,
			"products":[]any{},
		})
	 })

	 //! RBA implementation
	admin:= api.Group("/admin")
	admin.Use(middleware.RequireAdmin())
	
	admin.GET("/protected",func(ctx *gin.Context) {
		role,_:=middleware.GetUserRole(ctx) // if we want the userID
		ctx.JSON(http.StatusOK, gin.H{
			"ok":true,
			"role":role,
		})
	})


	return r
}

