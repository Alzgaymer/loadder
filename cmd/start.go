package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"loadder/internal/config"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts loadder.",
	Long:  `Starts loadder with specified port.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// parse config
		_, err := ParseConfig(cmd)
		if err != nil {
			return err
		}
		// set config
		// configure load balancer
		// start load balancer
	},
}

var (
	port       = new(uint16)
	configFile string
)

func StartCommand(flags *pflag.FlagSet) {
	flags.Uint16VarP(port, "port", "p", 8080, "Specifies application running port")
	configFile = *flags.String("config", ".loadder.yml", "Defines config file")
}

func InvalidValue(key, v, t string) error {
	return fmt.Errorf("invalid parametr `%s` with value:%s, type: %s", key, v, t)
}

func ParseConfig(cmd *cobra.Command) (*config.Config, error) {
	cfg := &config.Config{}

	v := cmd.Flag("port").Value
	if v.Type() != "uint16" {
		return nil, InvalidValue("port", v.String(), v.Type())
	}

	cfg.Port = v.String()
}
