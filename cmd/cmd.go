package cmd

import (
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	ex "github.com/dsrvlabs/vatz/manager/executor"
	health "github.com/dsrvlabs/vatz/manager/healthcheck"
)

const (
	defaultFlagConfig = "default.yaml"
	defaultFlagLog    = ""
	defaultRPC        = "http://localhost:19091"
	defaultPromPort   = "18080"
	defaultHomePath   = "~/.vatz"
)

var (
	healthChecker          = health.GetHealthChecker()
	executor               = ex.NewExecutor()
	dispatchers            []dp.Dispatcher
	defaultVerifyInterval  = 15
	defaultExecuteInterval = 30

	configFile string
	logfile    string
	vatzRPC    string
	promPort   string
)

// CreateRootCommand creates root command of Cobra.
func CreateRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{}
	cmd := createInitCommand(tp.LIVE)
	cmd.PersistentFlags().StringP("all", "a", "", "Set default config.yaml include default vatz-plugin-sysutil")
	viper.BindPFlag("initialize", installCommand.PersistentFlags().Lookup("all"))

	rootCmd.AddCommand(cmd)
	rootCmd.AddCommand(createStartCommand())
	rootCmd.AddCommand(createPluginCommand())
	rootCmd.AddCommand(createVersionCommand())

	return rootCmd
}
