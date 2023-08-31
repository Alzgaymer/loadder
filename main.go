package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
	"loadder/config"
	lb "loadder/load_balancer"
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

	// make services
	services, err := lb.Parse(cfg)
	if err != nil {
		logger.Panic("failed to parse config file", zap.Error(err), zap.String("config", *configPath))
	}

	algorithm, err := lb.DefineAlgorithm(cfg.Algorithm)
	if err != nil {
		logger.Panic("failed to define algorithm", zap.Error(err), zap.String("algorithm", cfg.Algorithm))
	}

	balancer, err := lb.NewBalancer(
		lb.WithAlgorithm(algorithm),
		lb.WithAddress(cfg.LoadBalancerAddress),
		lb.WithServices(services...),
	)
	if err != nil {
		logger.Panic("failed to create load balancer", zap.Error(err))
	}

	ctx := context.Background()
	if err = balancer.Run(ctx); err != nil {
		logger.Panic("failed to run load balancer", zap.Error(err))
	}

}
