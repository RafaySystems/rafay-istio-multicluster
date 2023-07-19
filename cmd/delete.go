package cmd

import (
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/commands"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/log"
	"github.com/spf13/cobra"
)

func newDeleteCmd(o commands.CmdOptions, logger log.Logger) *cobra.Command {
	// cmd represents the delete command
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete various resources in Console",
		Long:  `Delete clusters, namespaces, workloads, and other resources in your current project`,
		Args:  o.Validate,
		RunE:  o.Run,
	}

	// add subcommands here
	cmd.AddCommand()

	o.AddFlags(cmd)
	return cmd
}
