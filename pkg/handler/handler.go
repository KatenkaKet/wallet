package handler

import (
	"github.com/KatenkaKet/wallet/pkg/service"
	"github.com/gin-gonic/gin"
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

		//r.GET("/wallets/all", h.getWallets)
		//r.GET("/wallets/transactions", h.getWalletTransactions)
		//r.GET("/wallet/:id/transactions", h.getWalletTransactionsById)
	}

	return router
}
