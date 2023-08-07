package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			if viper.GetBool("devices.fans.enabled") && viper.GetBool("thSensor.enabled") {
				client.controllerService.FanControl()
			} else {
				log.Info().Msg("fans or temp humidity sensor are not enabled, skipping fan control")
			}
			if viper.GetBool("devices.lights.enabled") {
				client.controllerService.LightsControl()
			} else {
				log.Info().Msg("lights are not enabled, skipping light control")
			}
			ms := client.controllerService.Sensors.MoistureSensors()
			if viper.GetBool("moistureSensors.enabled") && len(ms) >= 1 {
				client.controllerService.WaterControl()
			} else {
				log.Info().Msg("moisture sensors are not enabled, or none are loaded, skipping water control")
			}
			time.Sleep(viper.GetDuration("controller.checkPeriod"))
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
