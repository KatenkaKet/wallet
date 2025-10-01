package handler

import (
	"net/http"

	"github.com/KatenkaKet/wallet"
	"github.com/gin-gonic/gin"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

func (h *Handler) createWalletTransaction(c *gin.Context) {
	var WT wallet.WalletTransactions

	if err := c.ShouldBindJSON(&WT); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if WT.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		return
	}

	if err := h.service.Wallet.UpdateBalance(c.Request.Context(), WT); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"walletId":      WT.ValletId.UUID.String(),
		"operationType": WT.OperationType,
		"amount":        WT.Amount,
		"status":        "success",
	})

}

func (h *Handler) getWalletBalance(c *gin.Context) {
	idParam := c.Param("id")

	var walletID uuid.UUID
	if err := walletID.Scan(idParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet id"})
		return
	}

	balance, err := h.service.Wallet.GetBalance(c.Request.Context(), walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"walletId": walletID.UUID.String(),
		"balance":  balance,
	})
}

//func (h *Handler) getWallets(context *gin.Context) {
//
//}
//
//func (h *Handler) getWalletTransactions(context *gin.Context) {
//
//}
//
//func (h *Handler) getWalletTransactionsById(context *gin.Context) {
//
//}
