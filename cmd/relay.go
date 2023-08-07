package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// relayCmd represents the relay command
var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		device, err := client.controllerService.Relays.DeviceByRelay(relayNumber)
		if err != nil {
			cmd.PrintErr("unable to find relay: %w", err)
			os.Exit(1)
		}
		if on {
			device.On()
		} else {
			device.Off()
		}
	},
}

var (
	relayNumber int
	on          bool
	off         bool
)

func init() {
	rootCmd.AddCommand(relayCmd)
	relayCmd.Flags().IntVarP(&relayNumber, "number", "n", 0, "relay to execute on")
	relayCmd.Flags().BoolVarP(&on, "on", "o", false, "turn on relay")
	relayCmd.Flags().BoolVarP(&off, "off", "f", false, "turn off relay")
	relayCmd.MarkFlagsMutuallyExclusive("on", "off")
	err := relayCmd.MarkFlagRequired("number")
	if err != nil {
		os.Exit(1)
	}
}
