package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequestResponse(ctx *gin.Context, payload interface{}) {
	WriteJsonResponse(ctx, http.StatusBadRequest, gin.H{
		"error":   true,
		"message": payload,
	})
}

func InternalServerJsonResponse(ctx *gin.Context, payload interface{}) {
	WriteJsonResponse(ctx, http.StatusInternalServerError, gin.H{
		"error":   true,
		"message": payload,
	})
}

func NotFoundResponse(ctx *gin.Context, payload interface{}) {
	WriteJsonResponse(ctx, http.StatusNotFound, gin.H{
		"error":   true,
		"message": payload,
	})
}

func WriteJsonResponse(ctx *gin.Context, status int, payload interface{}) {
	ctx.JSON(status, payload)
}

func UnauthorizedResponse(ctx *gin.Context, payload interface{}) {
	WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
		"error":   true,
		"message": payload,
	})
}
