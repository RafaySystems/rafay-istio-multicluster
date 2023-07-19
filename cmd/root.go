package cmd

import (
	"fmt"
	"os"

	"github.com/RafaySystems/rafay-istio-multicluster/pkg/commands"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/constants"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/log"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/output"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRootCmd() *cobra.Command {
	logger := log.GetLogger()
	o := commands.NewGlobalOptions(logger)
	// this cmd represents the base command when called without any subcommands
	cmd := &cobra.Command{
		Use:               constants.CLI_NAME,
		Short:             "A CLI tool to manage istio mesh across multiple Kubernetes clusters.",
		Long:              `A CLI tool to manage istio mesh across multiple Kubernetes clusters.`,
		TraverseChildren:  true,
		SilenceUsage:      true,
		PersistentPreRunE: o.Run,
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			output.Exit()
		},
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	o.AddFlags(cmd)

	// add subcommands here
	cmd.AddCommand(
		// new commands
		newDeleteCmd(commands.NewDeleteOnAllOptions(logger), logger),
		newApplyFileCmd(commands.NewApplyOnAllOptions(logger)),
		newDownloadCmd(logger),

		// version command
		newVersionCmd(),
	)

	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		if err := c.Help(); err != nil {
			return err
		}
		return err
	})

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// if any of the subcommands run into an error, it shows up here
func Execute() {
	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// run when each command's execute method is called
	// do command wide inits here
	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(home)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	log.GetLogger().Infof("Using config file:", viper.ConfigFileUsed())
	// }
}
