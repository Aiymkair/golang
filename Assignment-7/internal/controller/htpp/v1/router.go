package v1

import (
	"Assignment-7/internal/usecase"
	"Assignment-7/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, userUsecase usecase.UserInterface, log logger.Logger) {
	v1 := router.Group("/v1")
	{
		NewUserRoutes(v1, userUsecase, log)
	}
}
