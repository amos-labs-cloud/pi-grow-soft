package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stianeikeland/go-rpio/v4"
	"os"
	osSignal "os/signal"
	"syscall"
	"time"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := make(chan os.Signal)
		osSignal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			os.Exit(1)
		}()

		err := rpio.Open()
		if err != nil {
			cmd.PrintErrf("unable to open gpio: %s", err)
		}

		for {
			//client.controllerService.FanControl()
			client.controllerService.WaterControl()
			client.controllerService.LightsControl()
			time.Sleep(time.Second * 20)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
