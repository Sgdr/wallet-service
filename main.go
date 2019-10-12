package main

import (
	"github.com/Sgdr/wallet-service/internal/logger"
	"github.com/go-kit/kit/log/level"
)

func main() {
	log := logger.Init()
	level.Info(log).Log("msg", "wallet's service is starting...")

}
