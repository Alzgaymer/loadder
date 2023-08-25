package routes

import (
	"fmt"
	"loadder/internal/config"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func Routes(c *config.Config) ([]*httputil.ReverseProxy, error) {
	for _, service := range c.Services {
		ports, err := parsePorts(service.Ports, service.Exclude)
		if err != nil {
			return nil, err
		}

		proxies := make([]*httputil.ReverseProxy, len(ports))
		for i, port := range ports {
			u, err := url.Parse(port)
			if err != nil {
				return nil, err
			}
			proxy := httputil.NewSingleHostReverseProxy(u)

			proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {

			}

			proxies[i] = proxy

		}
	}

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
