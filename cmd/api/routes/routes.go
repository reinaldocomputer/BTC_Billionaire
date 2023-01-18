package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/reinaldocomputer/BTC_Billionaire/cmd/api/handlers"
)

func API() {
	r := gin.Default()
	r.POST("/sendBTC", handlers.HandleSendBTC)
	r.POST("/getHistory", handlers.HandleGetHistory)
	r.Run()
}
