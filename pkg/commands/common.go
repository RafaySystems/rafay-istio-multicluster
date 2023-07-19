package commands

import (
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/constants"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/context"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	YamlConfigFlag           = "config-file"
	YamlConfigShorthandFlag  = "f"
	VersionNameFlag          = "version-name"
	VersionNameShorthandFlag = "s"
)

func AddYamlConfigFlag(cmd *cobra.Command, flagHelp string) {
	cmd.PersistentFlags().StringP(YamlConfigFlag, YamlConfigShorthandFlag, "", flagHelp)
	cmd.MarkPersistentFlagFilename(YamlConfigFlag, "yml", "yaml")
}

func AddYamlConfigFlagVar(ptr *string, cmd *cobra.Command, flagHelp string) {
	cmd.PersistentFlags().StringVarP(ptr, YamlConfigFlag, YamlConfigShorthandFlag, "", flagHelp)
	cmd.MarkPersistentFlagFilename(YamlConfigFlag, "yml", "yaml")
}

type CmdOptions interface {
	// Validate is used to validate arguments and flags.
	// The function will validate args without calling REST APIs.
	// This means validating if files exist, if there are duplicate arguments provided,
	// if the proper flags are provided, etc.
	// It is also where config files are parsed
	Validate(cmd *cobra.Command, args []string) error
	// Run runs the command action
	Run(cmd *cobra.Command, args []string) error
	// AddFlags adds flags to the supplied cobra command
	AddFlags(cmd *cobra.Command)
}

// GlobalOptions is a struct to hold the values for global options
type GlobalOptions struct {
	Verbose,
	Debug bool
	Output,
	wait bool
	logger log.Logger
}

func NewGlobalOptions(log log.Logger) *GlobalOptions {
	globalOptions := new(GlobalOptions)
	globalOptions.logger = log
	return globalOptions
}

func (g *GlobalOptions) Validate(_ *cobra.Command, _ []string) error {
	return nil
}

func (g *GlobalOptions) Run(cmd *cobra.Command, _ []string) error {
	if cmd.Name() == "init" {
		return nil
	}
	isVerbose, err := cmd.Flags().GetBool(constants.VERBOSE_FLAG_NAME)
	if err != nil {
		isVerbose = false
	}

	isDebug, err := cmd.Flags().GetBool(constants.DEBUG_FLAG_NAME)
	if err != nil {
		isDebug = false
	}

	// set the desired logging level
	// by default, the log level is set to error
	if isVerbose {
		log.SetLevel(zap.InfoLevel)
	}
	if isDebug {
		log.SetLevel(zap.DebugLevel)
	}

	isStructuredOutput, err := cmd.Flags().GetBool(constants.STRUCTURED_OUTPUT_FLAG_NAME)
	if err != nil {
		isStructuredOutput = false
	}

	cliCtx := context.GetContext()

	cliCtx.Verbose = isVerbose
	cliCtx.Debug = isDebug
	cliCtx.StructuredOutput = isStructuredOutput

	log.GetLogger().Debugf("Prerun")

	// check if the wait flag is provided
	if cmd.Flags().Changed("wait") {
		g.wait = true
	}
	return nil
}

func (g *GlobalOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&g.Verbose, "verbose", "v", false, "Verbose mode. A lot more information output.")
	cmd.PersistentFlags().BoolVarP(&g.Debug, "debug", "d", false, "Enable debug logs")
	cmd.PersistentFlags().MarkHidden("v3")
	cmd.PersistentFlags().BoolVarP(&g.wait, "wait", "", false, "Wait for the operation to complete")
}
