package main

import (
	"github.com/go-kit/kit/log/level"
	"github.com/sgdr/wallet-service/internal/logger"
)

func main() {
	log := logger.Init()
	level.Info(log).Log("msg", "wallet's service is starting...")

}
