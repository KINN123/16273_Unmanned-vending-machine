package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"code":    -1,
		"message": err.Error(),
	})
}

func SuccessResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "success",
		"data":    data,
	})
}
