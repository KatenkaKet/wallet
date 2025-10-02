package handler

import (
	"github.com/KatenkaKet/wallet/pkg/service"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/KatenkaKet/wallet/docs"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	r := router.Group("/api/v1")
	{
		r.POST("/wallet", h.createWalletTransaction)
		r.GET("/wallets/:id", h.getWalletBalance)
	}

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
