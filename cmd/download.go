package cmd

import (
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/commands"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/log"
	"github.com/spf13/cobra"
)

func newDownloadCmd(logger log.Logger) *cobra.Command {
	// createCmd represents the create command
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download various resources in Console",
		Long:  `Download various resources in Console`,
	}

	// add subcommands here
	cmd.AddCommand(
		newDownloadIstioCTLCmd(commands.NewDownloadIstioCTLOptions(logger)),
	)

	return cmd
}

func newDownloadIstioCTLCmd(o commands.CmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "istio",
		Short:   "Download istio",
		Long:    "Download istio",
		Aliases: []string{"i", "di"},
		Args:    o.Validate,
		RunE:    o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
