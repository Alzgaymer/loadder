package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops executing",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stop called")
	},
}
