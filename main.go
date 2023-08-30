package main

import (
	"flag"
	"go.uber.org/zap"
	"loadder/config"
	logging "loadder/logger"
	"log"
	"os"
)

var configPath = flag.String("config", ".loadder.yml", "Specifies services for load balancer")

func main() {
	flag.Parse()

	logger, err := logging.New("")
	if err != nil {
		log.Fatal(err)
	}
	// open config file
	file, err := os.Open(*configPath)
	if err != nil {
		logger.Fatal("failed to open config file", zap.Error(err), zap.String("config", *configPath))
	}
	defer file.Close()
	// parse config file
	cfg, err := config.Parse(file)
	if err != nil {
		logger.Panic("invalid config file structure", zap.Error(err), zap.String("config", *configPath))
	}

}
