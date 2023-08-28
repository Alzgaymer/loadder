package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	"io"
	"loadder/internal/config"
	lb "loadder/internal/domain/load_balancer"
	"loadder/platform/service"
	"os"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts loadder.",
	Long:  `Starts loadder with specified port.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// parse config
		file, err := os.Open(configPath)
		if err != nil {
			return err
		}
		defer file.Close()

		cfg, err := ParseConfig(file)
		if err != nil {
			return err
		}

		services, err := service.Parse(cfg)
		if err != nil {
			return err
		}

		// configure load balancer
		loadBalancer := lb.NewLoadBalancer()

		for _, backendService := range services {
			loadBalancer.Add(backendService)
		}

		// start load balancer
		if err = loadBalancer.Run(cmd.Context()); err != nil {
			return err
		}

		return nil
	},
}

var (
	configPath string
)

func StartCommand(flags *pflag.FlagSet) {
	flags.StringVarP(&configPath, "config", "c", ".loadder.yml", "Defines config file")
}

// TODO move Parse to config package and refactor to provide yaml, json unmarshal
func ParseConfig(r io.Reader) (*config.Config, error) {
	cfg := &config.Config{}

	text, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(text, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
