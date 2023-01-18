package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/reinaldocomputer/BTC_Billionaire/internal/btc"
	"net/http"
)

func HandleSendBTC(c *gin.Context) {
	btcTransactions := []btc.SendBTCRequest{}
	if err := c.ShouldBindJSON(&btcTransactions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, element := range btcTransactions {
		transaction := btc.NewTransaction(element)
		// Valid method - checks data and convert to UTC Time Zone
		err := transaction.Valid()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = transaction.SendBTC()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "BTC Received",
	})
}

func HandleGetHistory(c *gin.Context) {
	req := btc.HistoryRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h := btc.NewHistory(req)
	if err := h.Valid(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transactions, err := h.GetHistory()
	fmt.Println("transactions: ", transactions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": transactions,
	})
}
