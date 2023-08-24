package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	"io"
	"loadder/internal/config"
	lb "loadder/internal/domain/load_balancer"
	"loadder/internal/domain/routes"
	"net/http"
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

		cfg, err := ParseConfig(cmd, file)
		if err != nil {
			return err
		}

		router, err := routes.ExtractRoutes(cfg)
		if err != nil {
			return err
		}

		server := &http.Server{
			Addr:    ":" + cfg.LoadBalancerPort,
			Handler: router,
		}

		// configure load balancer
		loadBalancer := lb.NewLoadBalancer(server)

		// start load balancer
		if err = loadBalancer.Run(cmd.Context()); err != nil {
			return err
		}

		return nil
	},
}

var (
	port       = new(uint16)
	configPath string
)

func StartCommand(flags *pflag.FlagSet) {
	flags.Uint16VarP(port, "port", "p", 8080, "Specifies application running port")
	flags.StringVarP(&configPath, "config", "c", ".loadder.yml", "Defines config file")
}

func InvalidValue(key, v, t string) error {
	return fmt.Errorf("invalid parametr `%s` with value:%s, type: %s", key, v, t)
}

// TODO move Parse to config package and refactor to provide yaml, json unmarshal
func ParseConfig(cmd *cobra.Command, r io.Reader) (*config.Config, error) {
	cfg := &config.Config{}

	v := cmd.Flag("port").Value
	if v.Type() != "uint16" {
		return nil, InvalidValue("port", v.String(), v.Type())
	}

	cfg.LoadBalancerPort = v.String()

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
