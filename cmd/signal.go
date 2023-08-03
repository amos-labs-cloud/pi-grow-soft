package cmd

import (
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/pin"
	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
	"os"

	"github.com/spf13/cobra"
)

// signalCommand represents the pin command
var signalCommand = &cobra.Command{
	Use:   "signal",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := rpio.Open()
		if err != nil {
			cmd.PrintErrf("unable to open gpio: %s", err)
		}
		defer func() {
			err := rpio.Close()
			if err != nil {
				cmd.PrintErrf("unable to close gpio: %s", err)
			}
		}()
		thePin := pin.NewPin(pinNumber)
		log.Debug().Msgf("signal is: %t", signal)
		if signal {
			thePin.High()
		} else {
			thePin.Low()
		}
	},
}

var (
	pinNumber int
	signal    bool
)

func init() {
	rootCmd.AddCommand(signalCommand)
	signalCommand.Flags().IntVarP(&pinNumber, "number", "n", 0, "gpio pin to execute on")
	signalCommand.Flags().BoolVarP(&signal, "signal", "s", false, "whether to send high")
	signalCommand.MarkFlagsRequiredTogether("number", "signal")

	err := signalCommand.MarkFlagRequired("number")
	if err != nil {
		signalCommand.Help()
		os.Exit(1)
	}
}
