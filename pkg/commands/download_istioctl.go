package commands

import (
	"fmt"
	"os"

	"github.com/RafaySystems/rafay-istio-multicluster/pkg/istioctl"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/log"
	"github.com/spf13/cobra"
)

const (
	DownloadIstioCTLVersion = "version"
)

type DownloadIstioCTLOptions struct {
	Version        string
	DownloadFolder string
	logger         log.Logger
}

func NewDownloadIstioCTLOptions(logger log.Logger) *DownloadIstioCTLOptions {
	o := new(DownloadIstioCTLOptions)
	o.logger = logger
	return o
}

func (c *DownloadIstioCTLOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.ExactArgs(0)(cmd, args)
}

func (c *DownloadIstioCTLOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Infof("Start [%s]", cmd.CommandPath())

	version, _ := cmd.Flags().GetString("version")
	folder, _ := cmd.Flags().GetString("folder")

	if version == "" {
		fmt.Println("version cannot be empty")
		os.Exit(1)
	}
	istioctlPath, err := istioctl.DownloadIsctioCtl(version, folder)
	if err != nil {
		log.GetLogger().Errorf("failed to get istioctl", err)
		os.Exit(1)
	}

	log.GetLogger().Debugf("istioctlPath %s", istioctlPath)
	log.GetLogger().Infof("End [%s]", cmd.CommandPath())
	return nil
}

func (c *DownloadIstioCTLOptions) AddFlags(cmd *cobra.Command) {
	// add flags
	flagSet := cmd.Flags()
	flagSet.StringVarP(&c.Version, "version", "", "", "Set the istio version to download")
	flagSet.StringVar(&c.DownloadFolder, "folder", "", "Set the folder to download istio")
}
