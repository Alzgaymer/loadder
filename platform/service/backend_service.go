package service

import (
	"context"
	"fmt"
	"loadder/internal/config"
	lb "loadder/internal/domain/load_balancer"
	"loadder/platform/backend"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type BackendService struct { //nolint:govet
	name string

	server *http.Server
}

func NewBackendService(name string, server *http.Server) *BackendService {
	return &BackendService{name: name, server: server}
}

func (b *BackendService) Name() string {
	return b.name
}

func (b *BackendService) Shutdown(ctx context.Context) error {
	return b.server.Shutdown(ctx)
}

func (b *BackendService) ListenAndServe() error {
	return b.server.ListenAndServe()
}

func Parse(c *config.Config) ([]lb.Service, error) {
	services := make([]lb.Service, 0, len(c.Services))

	for _, service := range c.Services {
		ports, err := parsePorts(service.Ports, service.Exclude)
		if err != nil {
			return nil, err
		}

		urls, err := parsURLs(parseHostPorts(service.Address, ports))
		if err != nil {
			return nil, err
		}

		backends := backend.CreateBackends(urls...)

		algorithm, err := DefineAlgorithm(service.Algorithm, backends...)
		if err != nil {
			return nil, err
		}

		server := &http.Server{
			Addr:    ":" + service.ProxyPort,
			Handler: algorithm,
		}

		services = append(services, NewBackendService(service.Name, server))
	}

	return services, nil
}

func parsURLs(addresses []string) ([]*url.URL, error) {
	var (
		n   = len(addresses)
		res = make([]*url.URL, n)
		err error
	)

	for i := 0; i < n; i++ {
		res[i], err = url.Parse(addresses[i])
		if err != nil {
			return nil, err
		}
	}

	return res, err
}

func parseHostPorts(host string, ports []string) []string {
	n := len(ports)
	res := make([]string, n)
	sb := &strings.Builder{}

	for i := 0; i < n; i++ {
		sb.WriteString(host)
		sb.WriteString(":")
		sb.WriteString(ports[i])

		res[i] = sb.String()

		sb.Reset()
	}

	return res
}

func parsePorts(portsStr string, toExclude []string) ([]string, error) {
	portRanges := strings.Split(portsStr, "...")
	if len(portRanges) != 2 {
		return nil, fmt.Errorf("invalid port format: %s", portsStr)
	}

	startPort, err := strconv.Atoi(portRanges[0])
	if err != nil {
		return nil, err
	}

	endPort, err := strconv.Atoi(portRanges[1])
	if err != nil {
		return nil, err
	}

	if (startPort == endPort) || (startPort <= 0 || endPort <= 0) || (startPort > 65535 || endPort > 65535) {
		return nil, fmt.Errorf("invalid port format: %s", portsStr)
	}

	ports := []int{}
	for p := startPort; p <= endPort; p++ {
		ports = append(ports, p)
	}

	return resolveExcludedPorts(ports, toExclude), nil
}

func resolveExcludedPorts(ports []int, excluded []string) []string {
	excludedPorts := map[int]bool{}

	for _, ex := range excluded {
		port, _ := strconv.Atoi(ex)
		excludedPorts[port] = true
	}

	var filteredPorts []string

	for _, p := range ports {
		if !excludedPorts[p] {
			filteredPorts = append(filteredPorts, strconv.Itoa(p))
		}
	}

	return filteredPorts
}
