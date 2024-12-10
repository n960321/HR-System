package cmd

import (
	"HRSystem/internal/config"
	"HRSystem/pkg/logger"
	"HRSystem/pkg/server"
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "HRSystem",
	Short: "Run a HR system for managing employee attendance records",
	Long: `Run a simple HR backend system with the following features:
		1. User login
		2. Admin creates employee accounts 
		3. User changes password
		4. User clock in/out
		5. User retrieves attendance records
	`,
	Run: RunServer,
}

var (
	configFile string
	local      bool
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "configs/config.yaml", "The config file.")
	rootCmd.PersistentFlags().BoolVarP(&local, "local", "l", false, "Run on local (true or false)")
}

func RunServer(cmd *cobra.Command, args []string) {
	logger.SetLogger(local)
	config := config.GetConfig(configFile)

	svr := server.NewServer(config.Http)

	svr.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	svr.Shutdown(ctx)

	log.Info().Msg("shutting down")
	os.Exit(0)

}
