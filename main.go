package main

import (
	"fmt"
	"github.com/reinaldocomputer/BTC_Billionaire/cmd/api/routes"
	"github.com/reinaldocomputer/BTC_Billionaire/internal/platform/mongodb"
	"os"
)

func init() {
	err := mongodb.Connect()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	routes.API()
}
