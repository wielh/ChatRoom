package main

import (
	"common"
	"fmt"
	"serviceClient"

	"github.com/gin-gonic/gin"
)

func main() {
	common.ConfigInit()
	err := serviceClient.MicroServiceClientInit()
	if err != nil {
		common.ErrorLogger("gate", "main", "create client connection to each micro-service failed", err)
		return
	}
	router := gin.Default()
	setRouter(router)
	common.InfoLogger("gate", "main", fmt.Sprintf("running gateway on port %d", common.GatePort))
	router.Run(fmt.Sprintf(":%d", common.GatePort))
}
