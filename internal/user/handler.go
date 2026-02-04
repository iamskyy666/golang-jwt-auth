package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Handlers
func (h *Handler) RegisterUser(ctx *gin.Context) {
	var input RegisterInput
	if err:=ctx.ShouldBindJSON(&input); err!=nil{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":"⚠️ Invalid JSON body!",
			"status_code":http.StatusBadRequest,
		})
	}

	output,err:= h.svc.Register(ctx.Request.Context(),input)
	if err !=nil{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":err.Error(),
			"status_code":http.StatusBadRequest,
		})
		return
	}

	ctx.JSON(http.StatusCreated, output)
}


// 07: