package handler

import (
	"net/http"
	"strings"

	"github.com/KatenkaKet/wallet"
	"github.com/gin-gonic/gin"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

// createWalletTransaction godoc
// @Summary Пополнить или снять деньги с кошелька, а также записать историю выполненных операций
// @Tags wallet
// @Accept json
// @Produce json
// @Param transaction body wallet.WalletTransactions true "Данные транзакции"
// @Success 200 {object} map[string]string "status: success"
// @Failure 400 {object} map[string]string "Ошибка валидации или неверные данные"
// @Failure 500 {object} map[string]string "Ошибка при обновлении баланса"
// @Router /wallet [post]
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
		"status": "success",
	})
}

// getWalletBalance godoc
// @Summary Получить баланс кошелька по ID
// @Tags wallet
// @Produce json
// @Param id path string true "ID кошелька"
// @Success 200 {object} map[string]float64 "Баланс кошелька"
// @Failure 400 {object} map[string]string "Неверный ID кошелька"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /wallets/{id} [get]
func (h *Handler) getWalletBalance(c *gin.Context) {
	idParam := c.Param("id")

	idParam = strings.TrimSpace(idParam)

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
		"balance": balance,
	})
}
