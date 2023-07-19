package commands

import (
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/common"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/constants"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/istioctl"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/log"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// flagpole
type DeleteOnAllOptions struct {
	YamlConfigFilePath string
	Yes                bool
	logger             log.Logger
}

const (
	ConfirmFlag = "confirm-delete"
)

func (o *DeleteOnAllOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *DeleteOnAllOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Debugf("apply start %s", o.YamlConfigFilePath)
	// retrieve the flags
	flagSet := cmd.Flags()

	// Nothing to do...
	if !flagSet.Changed(YamlConfigFlag) {
		log.GetLogger().Errorf("failed to get config flag")
		return nil
	}

	// read the config file YAML
	fileBytes, err := utils.ReadYAMLFileContents(o.YamlConfigFilePath)
	if err != nil {
		log.GetLogger().Errorf("failed to read config flag %s error %s", o.YamlConfigFilePath, err)
		return err
	}

	return deleteFileBytes(cmd, o, fileBytes)
}

func (o *DeleteOnAllOptions) AddFlags(cmd *cobra.Command) {
	AddYamlConfigFlagVar(&o.YamlConfigFilePath, cmd, "Use this flag to delete resources using a YAML file")
	cmd.Flags().BoolVar(&o.Yes, ConfirmFlag, false, "Confirm delete resource(s) in the file")
}

func NewDeleteOnAllOptions(logger log.Logger) CmdOptions {
	options := new(DeleteOnAllOptions)
	options.logger = logger
	return options
}

// describe function deleteFileBytes
// deleteFileBytes deletes the resources in the YAML file
// fileBytes is the YAML file contents
// returns error
func deleteFileBytes(cmd *cobra.Command, o *DeleteOnAllOptions, fileBytes []byte) error {
	var err error

	// split the file and update individual resources
	m, _, err := utils.SplitYamlAndGetListByKind(fileBytes)
	if err != nil {
		return err
	}

	log.GetLogger().Debugf("deleteFileBytes map %s", m)

	// iterate over the clusters
	for _, resource := range m["Cluster"] {
		var cluster common.Cluster

		log.GetLogger().Debugf("cluster resoure %s", resource)
		err = yaml.Unmarshal(resource, &cluster)
		if err != nil {
			log.GetLogger().Infof("failed parse cluster %s", resource)
			return err
		}

		// get the istioctl binary
		istioCtlDir, err := istioctl.DownloadIsctioCtl(cluster.Spec.Version, "")
		if err != nil {
			log.GetLogger().Debugf("failed to get istioctl binary info for version %s error", cluster.Spec.Version, err)
			return err
		}

		// istioctl uninstall --purge
		// uninstall the istio control plane using above command
		output, err := istioctl.RunIstioCtlCmd(cluster.Spec.Version, istioCtlDir, "uninstall", cluster.Spec.KubeconfigFile, cluster.Spec.Context, "", "")
		if err != nil {
			log.GetLogger().Debugf("failed to uninstall using istioctl error %s", err)
			return err
		}

		if output.Len() > 0 {
			log.GetLogger().Debugf("istioctl output to deploy controlplane %s", output.String())
		}

		// delete the namespace istio-system, HELLO_WORLD_NAMESPACE
		deleteNamespace(&cluster)
	}

	return nil
}

// describe function deleteNamespace
// deleteNamespace deletes the namespace istio-system and HELLO_WORLD_NAMESPACE
// cluster is the cluster object
// returns error
func deleteNamespace(cluster *common.Cluster) error {
	// kubectl delete namespace istio-system
	deeteClusterNamespace := []string{
		"--kubeconfig",
		cluster.Spec.KubeconfigFile,
		"--context=" + cluster.Spec.Context,
		"delete",
		"namespace",
		"istio-system",
	}

	output, err := utils.KubeCtlCmd(deeteClusterNamespace)
	if err != nil {
		log.GetLogger().Debugf("failed to delete istio-system namespace error %s", err)
		return err
	}

	if output.Len() > 0 {
		log.GetLogger().Debugf("kubectl output to delete  istio-system namespace output %s", output.String())
	}

	if cluster.Spec.InstallHelloWorld == "true" {
		// kubectl delete namespace HELLO_WORLD_NAMESPACE
		deeteClusterNamespace = []string{
			"--kubeconfig",
			cluster.Spec.KubeconfigFile,
			"--context=" + cluster.Spec.Context,
			"delete",
			"namespace",
			constants.HELLO_WORLD_NAMESPACE,
		}

		output, err := utils.KubeCtlCmd(deeteClusterNamespace)
		if err != nil {
			log.GetLogger().Debugf("failed to delete istio-system namespace error %s", err)
			return err
		}

		if output.Len() > 0 {
			log.GetLogger().Debugf("kubectl output to delete  istio-system namespace output %s", output.String())
		}

	}

	return nil
}
