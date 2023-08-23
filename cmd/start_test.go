package cmd_test

import (
	"loadder/cmd"
	"loadder/internal/config"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestParseConfig(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		input       string
		expectedCfg *config.Config
		wantedErr   bool
	}{
		{
			name: "ValidConfig",
			args: []string{"--port", "8080"},
			input: `
services:
  example-service:
    service-name: "Example"
    address: "http://example.com"
    ports:
      - "80"
    requests:
      http:
        GET:
          /hello-world : RR
`,
			expectedCfg: &config.Config{
				LoadBalancerPort: "8080",
				Services: map[string]*config.Service{
					"example-service": {
						Name:    "Example",
						Address: "http://example.com",
						Ports:   []string{"80"},
						Requests: &config.Requests{
							HTTP: config.HTTPRequest{
								GET: map[string]string{
									"/hello-world": "RR",
								},
							},
						},
					},
				},
			},
			wantedErr: false,
		},
		{
			name: "Add second method request",
			args: []string{"--port", "8081"},
			input: `
services:
  example-service:
    service-name: "Example"
    address: "http://example.com"
    ports:
      - "80"
    requests:
      http:
        GET:
          /hello-world : RR
          /stop-the-world/{s} : RR
        POST:
          /create-smth : RR
`,
			expectedCfg: &config.Config{
				LoadBalancerPort: "8081",
				Services: map[string]*config.Service{
					"example-service": {
						Name:    "Example",
						Address: "http://example.com",
						Ports:   []string{"80"},
						Requests: &config.Requests{
							HTTP: config.HTTPRequest{
								GET: map[string]string{
									"/hello-world":        "RR",
									"/stop-the-world/{s}": "RR",
								},
								POST: map[string]string{
									"/create-smth": "RR",
								},
							},
						},
					},
				},
			},
			wantedErr: false,
		},
		{
			name: "Invalid request method",
			args: []string{"--port", "8081"},
			input: `
services:
  example-service:
    service-name: "Example"
    address: "http://example.com"
    ports:
      - "80"
    requests:
      http:
        GET:
          - /hello-world : RR
          - /stop-the-world/{s} : RR
        POST:
          - /create-smth : RR
`,
			expectedCfg: nil,
			wantedErr:   true,
		},
	}
	c := &cobra.Command{}
	c.Flags().Uint16("port", 8080, "Port number")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c.SetArgs(tc.args)
			_ = c.Execute()

			reader := strings.NewReader(tc.input)

			cfg, err := cmd.ParseConfig(c, reader)

			if (tc.wantedErr == true && err != nil) && !configEqual(cfg, tc.expectedCfg) {
				t.Errorf("Expected config: %+v, but got: %+v", tc.expectedCfg, cfg)
			}
		})
	}
}

func configEqual(cfg, cfg2 *config.Config) bool {
	return reflect.DeepEqual(cfg, cfg2)
}
