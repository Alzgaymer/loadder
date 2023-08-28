package routes

import (
	"fmt"
	"loadder/internal/config"
	lb "loadder/internal/domain/load_balancer"
	"net/url"
	"strconv"
	"strings"
)

func Parse(c *config.Config) ([]*lb.Service, error) {
	services := make([]*lb.Service, 0, len(c.Services))

	for _, service := range c.Services {
		ports, err := parsePorts(service.Ports, service.Exclude)
		if err != nil {
			return nil, err
		}

		urls, err := parsURLs(parseHostPorts(service.Address, ports))
		if err != nil {
			return nil, err
		}

		backends := lb.ParseBackends(urls...)

		services = append(services, lb.NewService(backends))
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
