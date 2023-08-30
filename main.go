package main

import (
	"context"
	"flag"
	"go.uber.org/zap"
	"loadder/config"
	"loadder/load_balancer"
	logging "loadder/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	services, err := lb.Parse(cfg)
	if err != nil {
		logger.Panic("invalid config file structure", zap.Error(err), zap.String("config", *configPath))
	}

	// configure load balancer
	loadBalancer, err := lb.NewLoadBalancer(&http.Server{}, logger, services, lb.WithAlgorithm(cfg.Algorithm), lb.WithPort(cfg.LoadBalancerAddress))
	if err != nil {
		logger.Panic("could not create load balancer", zap.Error(err), zap.String("config", *configPath))
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// start load balancer
	if err = loadBalancer.Run(ctx); err != nil {
		logger.Panic("failed to run load balancer", zap.Error(err), zap.String("port", cfg.LoadBalancerAddress))
	}
}
