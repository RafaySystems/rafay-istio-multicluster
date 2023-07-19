package commands

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/RafaySystems/rafay-istio-multicluster/internal/fixtures"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/audit"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/common"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/constants"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/istioctl"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/log"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/stepcerts"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

// describe ApplyOnAllOptions
// ApplyOnAllOptions is the struct to hold the options for apply command
// YamlConfigFilePath is the path to the YAML file containing the resources to apply
// DryRun is the flag to indicate if the apply operation should be a dry run
// logger is the logger instance
type ApplyOnAllOptions struct {
	YamlConfigFilePath string
	DryRun             bool
	logger             log.Logger
}

const (
	DryRunFlag = "dry-run"
)

var DryRun bool

// describe Validate
// Validate is the function to validate the apply command options
func (o *ApplyOnAllOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

var dryRunUnsupportedError = `Dry run unsupported for resource type "%s". Skipping apply operation.\n`

// describe Run
// Run is the function to run the apply command
// cmd is the cobra command instance
// args is the list of arguments
// returns error
func (o *ApplyOnAllOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Debugf("apply start %s", o.YamlConfigFilePath)
	// retrieve the flags
	flagSet := cmd.Flags()

	DryRun = flagSet.Changed(DryRunFlag) && o.DryRun
	waitForOperationToComplete, err := flagSet.GetBool("wait")
	if err != nil {
		log.GetLogger().Errorf("failed to get global wait flag %s", err)
		return err
	}

	if DryRun && waitForOperationToComplete {
		return fmt.Errorf("wait flag cannot be set during dryrun")
	}

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

	return applyFileBytes(cmd, o, fileBytes, waitForOperationToComplete)
}

// describe function getCertFromMeshID
// getCertFromMeshID is the function to get the certificate for the given meshID
// meshID is the meshID to get the certificate for
// m is the map of resources from the YAML file
// returns *common.Certificate, string, error
// *common.Certificate is the certificate for the given meshID
// string is the path to the certificate
// error is the error encountered
func getCertFromMeshID(meshID string, m map[string][][]byte) (*common.Certificate, string, error) {
	for _, resource := range m["Certificate"] {
		var cert common.Certificate
		var cPath string
		log.GetLogger().Debugf("certificate resoure %s", resource)
		err := yaml.Unmarshal(resource, &cert)
		if err != nil {
			log.GetLogger().Infof("failed parse certificate %s", resource)
			return nil, "", err
		}

		if cert.Spec.Folder == "" {
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cPath = path.Join(home, constants.CLI_HOME_FOLDER, "istio", "certs")
			err = os.MkdirAll(cPath, 0755)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			cPath = cert.Spec.Folder
		}

		if cert.Spec.MeshID == meshID {
			log.GetLogger().Debugf("found cert matching meshid %s", meshID)
			return &cert, cPath, nil
		}
	}

	log.GetLogger().Debugf("did not find cert matching meshid %s", meshID)
	return nil, "", fmt.Errorf("did not find cert matching meshid %s", meshID)
}

// describe function combineFiles
// combineFiles is the function to combine two files into a destination file
// file1 is the first file to combine
// file2 is the second file to combine
// destination is the destination file path
// returns error
func combineFiles(file1, file2, destination string) error {
	f1, err := os.Open(file1)
	if err != nil {
		log.GetLogger().Debugf("failed to open %s", file1)
		return err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.GetLogger().Debugf("failed to open %s", file1)
		return err
	}
	defer f2.Close()

	out, err := os.Create(destination)
	if err != nil {
		log.GetLogger().Debugf("failed to open %s", file1)
		return err
	}
	defer out.Close()

	n1, err := io.Copy(out, f1)
	if err != nil {
		log.GetLogger().Debugf("failed to read %s", file1)
		return err
	}

	n2, err := io.Copy(out, f2)
	if err != nil {
		log.GetLogger().Debugf("failed to read %s", file2)
		return err
	}
	log.GetLogger().Debugf("wrote %d bytes to %s", n1+n2, out)

	return nil
}

// describe function createClusterYaml
// createClusterYaml is the function to create a YAML file for the given cluster
// cluster is the cluster to create the YAML file for
// fileName is the name of the file to create
// bb is the bytes buffer to write to the file
// returns string, error
// string is the path to the created YAML file
// error is the error encountered
func createClusterYaml(cluster *common.Cluster, fileName string, bb *bytes.Buffer) (string, error) {
	clusterDir := path.Join(audit.AuditFolder, cluster.Metadata.Name)
	err := os.MkdirAll(clusterDir, 0755)
	if err != nil {
		log.GetLogger().Debugf("failed to mkdir %s", clusterDir)
		return "", err
	}
	fileYaml := path.Join(clusterDir, fileName)
	yamlFile, err := os.Create(fileYaml)
	if err != nil {
		log.GetLogger().Debugf("failed to open %s", fileYaml)
		return "", err
	}
	defer yamlFile.Close()

	_, err = yamlFile.Write(bb.Bytes())
	if err != nil {
		log.GetLogger().Debugf("failed to write to %s", fileYaml)
		return "", err
	}

	return fileYaml, nil
}

// describe function deployEastWestGateway
// deployEastWestGateway is the function to deploy the eastwest gateway for the given cluster
// cluster is the cluster to deploy the eastwest gateway for
// returns error
// error is the error encountered
func deployEastWestGateway(cluster *common.Cluster) error {
	ew := new(bytes.Buffer)
	// prepare the template values
	valData := common.TempllateValues{
		MeshID:      cluster.Spec.MeshID,
		ClusterName: cluster.Metadata.Name,
	}
	// perpare the eastwest gateway YAML from template
	err := fixtures.EastWestGateway.Execute(ew,
		struct {
			Values common.TempllateValues
		}{
			valData,
		})
	if err != nil {
		log.GetLogger().Debugf("fixtures.EastWestGateway failed to execute")
		return err
	}

	// create the YAML file for the eastwest gateway
	ewFile, err := createClusterYaml(cluster, "eastwest-gateway.yaml", ew)
	if err != nil {
		return err
	}

	// if dry run, return
	if DryRun {
		return nil
	}

	// download the istioctl binary if not already downloaded
	istioCtlDir, err := istioctl.DownloadIsctioCtl(cluster.Spec.Version, "")
	if err != nil {
		log.GetLogger().Debugf("failed to get istioctl binary info for version %s error", cluster.Spec.Version, err)
		return err
	}

	// istioctl install -y -f eastwest-gateway.yaml
	// deploy the eastwest gateway using istioctl uses above istioctl syntax
	output, err := istioctl.RunIstioCtlCmd(cluster.Spec.Version, istioCtlDir, "install", cluster.Spec.KubeconfigFile, cluster.Spec.Context, ewFile, "")
	if err != nil {
		log.GetLogger().Debugf("failed to deploy eastwest-gateway using istioctl error %s", err)
		return err
	}

	if output.Len() > 0 {
		log.GetLogger().Debugf("istioctl output to deploy  eastwest-gateway %s %s", ewFile, output.String())
	}

	return nil
}

// describe function deployControlPlane
// deployControlPlane is the function to deploy the control plane for the given cluster
// cluster is the cluster to deploy the control plane for
// returns error
// error is the error encountered
func deployControlPlane(cluster *common.Cluster) error {
	cp := new(bytes.Buffer)
	// prepare the template values
	valData := common.TempllateValues{
		MeshID:      cluster.Spec.MeshID,
		ClusterName: cluster.Metadata.Name,
	}
	// perpare the control plane YAML from template
	err := fixtures.ControlPlane.Execute(cp, struct {
		Values common.TempllateValues
	}{
		valData,
	})
	if err != nil {
		log.GetLogger().Debugf("fixtures.ControlPlnae failed to execute")
		return err
	}

	// create the YAML file for the control plane
	cpFile, err := createClusterYaml(cluster, "controlplane.yaml", cp)
	if err != nil {
		return err
	}

	// download the istioctl binary if not already downloaded
	istioCtlDir, err := istioctl.DownloadIsctioCtl(cluster.Spec.Version, "")
	if err != nil {
		log.GetLogger().Debugf("failed to get istioctl binary info for version %s error", cluster.Spec.Version, err)
		return err
	}

	// if dry run, return
	if DryRun {
		return nil
	}

	// istioctl install -y -f controlplane.yaml
	// deploy the control plane using istioctl uses above istioctl syntax
	output, err := istioctl.RunIstioCtlCmd(cluster.Spec.Version, istioCtlDir, "install", cluster.Spec.KubeconfigFile, cluster.Spec.Context, cpFile, "")
	if err != nil {
		log.GetLogger().Debugf("failed to deploy contrpolplane using istioctl error %s", err)
		return err
	}

	if output.Len() > 0 {
		log.GetLogger().Debugf("istioctl output to deploy controlplane %s %s", cpFile, output.String())
	}

	return nil
}

// describe function deployExposeService
// deployExposeService is the function to deploy the expose service for the given cluster
// cluster is the cluster to deploy the expose service for
// returns error
// error is the error encountered
func deployExposeService(cluster *common.Cluster) error {
	es := new(bytes.Buffer)
	// prepare the template values
	valData := common.TempllateValues{
		MeshID:      cluster.Spec.MeshID,
		ClusterName: cluster.Metadata.Name,
	}
	// perpare the expose service YAML from template
	err := fixtures.ExposeService.Execute(es, struct {
		Values common.TempllateValues
	}{
		valData,
	})
	if err != nil {
		log.GetLogger().Debugf("fixtures.ExposeService failed to execute")
		return err
	}

	// create the YAML file for the expose service
	esFile, err := createClusterYaml(cluster, "expose-services.yaml", es)
	if err != nil {
		return err
	}

	// if dry run, return
	if DryRun {
		return nil
	}

	// kubectl apply -n istio-system -f expose-services.yaml
	// deploy the expose service using kubectl uses above kubectl syntax
	applyClusterNamespace := []string{
		"--kubeconfig",
		cluster.Spec.KubeconfigFile,
		"--context=" + cluster.Spec.Context,
		"-n",
		"istio-system",
		"apply",
		"-f",
		esFile,
	}

	// apply the expose service
	output, err := utils.KubeCtlCmd(applyClusterNamespace)
	if err != nil {
		log.GetLogger().Debugf("failed to apply cluster namespace from %s error %s", esFile, err)
		return err
	}

	if output.Len() > 0 {
		log.GetLogger().Debugf("kubectl output of apply expose-services.yaml %s output %s", esFile, output.String())
	}

	return nil
}

func deployNamespace(cluster *common.Cluster) error {

	ns := new(bytes.Buffer)
	valData := common.TempllateValues{
		MeshID:      cluster.Spec.MeshID,
		ClusterName: cluster.Metadata.Name,
	}
	err := fixtures.NameSpace.Execute(ns, struct {
		Values common.TempllateValues
	}{
		valData,
	})
	if err != nil {
		log.GetLogger().Debugf("fixtures.NameSpace failed to execute")
		return err
	}

	nsFile, err := createClusterYaml(cluster, "namespace.yaml", ns)
	if err != nil {
		return err
	}

	if DryRun {
		return nil
	}

	// kubectl --context="ctx-${i}" apply -f tmp/namespace-${i}.yaml
	applyClusterNamespace := []string{
		"--kubeconfig",
		cluster.Spec.KubeconfigFile,
		"--context=" + cluster.Spec.Context,
		"-n",
		"istio-system",
		"apply",
		"-f",
		nsFile,
	}

	output, err := utils.KubeCtlCmd(applyClusterNamespace)
	if err != nil {
		log.GetLogger().Debugf("failed to apply cluster cert from %s error %s", nsFile, err)
		return err
	}

	if output.Len() > 0 {
		log.GetLogger().Debugf("kubectl output to apply cluster cert from %s output %s", nsFile, output.String())
	}

	return nil
}

// describe function installClusterCerts
// installClusterCerts is the function to install the cluster certs for the given cluster
// cluster is the cluster to install the cluster certs for
// rootCertPem is the path to the root cert
// rootCakey is the path to the root key
// clusterChainPem is the path to the cluster chain cert
// clusterCertPem is the path to the cluster cert
// clusterKeyPem is the path to the cluster key
// returns error
// error is the error encountered
func installClusterCerts(cluster *common.Cluster, rootCertPem, rootCakey, clusterChainPem, clusterCertPem, clusterKeyPem string) error {
	// kubectl create secret generic cacerts -n istio-system \
	//     --from-file=ca-cert.pem \
	//     --from-file=ca-key.pem \
	//     --from-file=root-cert.pem \
	//     --from-file=cert-chain.pem --dry-run -o yaml
	// create the secret for the cluster certs using abiove kubectl syntax
	createSecret := []string{
		"create",
		"secret",
		"generic",
		"cacerts",
		"-n",
		"istio-system",
		"--from-file=" + clusterCertPem,
		"--from-file=" + clusterKeyPem,
		"--from-file=" + rootCertPem,
		"--from-file=" + clusterChainPem,
		"--dry-run",
		"-o",
		"yaml",
	}

	// execute the above kubectl command
	output, err := utils.KubeCtlCmd(createSecret)
	if err != nil {
		log.GetLogger().Debugf("installClusterCerts failed to cretae certs.yaml err %s", err)
		return err
	}

	// create the folder for the cluster certs
	certsDir := path.Join(audit.AuditFolder, cluster.Metadata.Name)
	err = os.MkdirAll(certsDir, 0755)
	if err != nil {
		log.GetLogger().Debugf("failed to mkdir %s", certsDir)
		return err
	}
	certsYaml := path.Join(certsDir, "certs.yaml")
	certsYamlFile, err := os.Create(certsYaml)
	if err != nil {
		log.GetLogger().Debugf("failed to open %s", certsYaml)
		return err
	}
	defer certsYamlFile.Close()

	// write the output of the above kubectl command to the certs.yaml file
	_, err = certsYamlFile.Write(output.Bytes())
	if err != nil {
		log.GetLogger().Debugf("failed to write to %s", certsYaml)
		return err
	}

	// if dry run, return
	if DryRun {
		return nil
	}

	//kubectl -n istio-system apply -f certs.yaml
	// apply the certs.yaml using above kubectl syntax
	applyClusterSecret := []string{
		"--kubeconfig",
		cluster.Spec.KubeconfigFile,
		"--context=" + cluster.Spec.Context,
		"-n",
		"istio-system",
		"apply",
		"-f",
		certsYaml,
	}

	// execute the above kubectl command
	output, err = utils.KubeCtlCmd(applyClusterSecret)
	if err != nil {
		log.GetLogger().Debugf("failed to apply cluster cert from %s error %s", certsYaml, err)
		return err
	}

	if output.Len() > 0 {
		log.GetLogger().Debugf("kubectl output to apply cluster cert from %s output %s", certsYaml, output.String())
	}

	return nil
}

// describe function applyFileBytes
// applyFileBytes is the function to apply the resources in the given YAML file
// cmd is the cobra command instance
// o is the apply command options
// fileBytes is the bytes of the YAML file
// waitForOperationToComplete is the flag to indicate if the apply operation should wait for completion
// returns error
// error is the error encountered
func applyFileBytes(cmd *cobra.Command, o *ApplyOnAllOptions, fileBytes []byte, waitForOperationToComplete bool) error {
	// split the file and update individual resources
	m, _, err := utils.SplitYamlAndGetListByKind(fileBytes)
	if err != nil {
		return err
	}

	log.GetLogger().Debugf("applyFileBytes map %s", m)

	// process certificates first
	for _, resource := range m["Certificate"] {
		var cert common.Certificate
		var cPath string

		log.GetLogger().Debugf("certificate resoure %s", resource)
		err := yaml.Unmarshal(resource, &cert)
		if err != nil {
			log.GetLogger().Infof("failed parse certificate %s", resource)
			return err
		}

		if cert.Spec.Folder == "" {
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cPath = path.Join(home, constants.CLI_HOME_FOLDER, "istio", "certs")
			err = os.MkdirAll(cPath, 0755)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			cPath = cert.Spec.Folder
		}

		_, err = os.Stat(cPath)
		if err != nil {
			log.GetLogger().Debugf("cannot stat certificate path %s", cPath)
			return err
		}

		// check root-ca.key exist or not
		_, err = os.Stat(path.Join(cPath, "root-ca.key"))
		if err != nil {
			pemPath := path.Join(cPath, "root-cert.pem")
			keyPath := path.Join(cPath, "root-ca.key")
			san := "root." + cert.Spec.SanSuffix
			// Generate rootca
			args := []string{
				"step",
				"certificate",
				"create",
				san,
				pemPath,
				keyPath,
				"--profile",
				"root-ca",
				"--no-password",
				"--insecure",
				"--san",
				san,
				"--not-after",
				cert.Spec.ValidityHours + "h",
				"--kty",
				"RSA",
			}
			log.GetLogger().Debugf("generate rootca in path %s args:%s", cPath, args)
			stepcerts.GenStepCerts(args)
		}
	}

	// deploy namespaces
	for _, resource := range m["Cluster"] {
		var cluster common.Cluster

		err = yaml.Unmarshal(resource, &cluster)
		if err != nil {
			log.GetLogger().Infof("failed parse cluster %s", resource)
			return err
		}

		// create namespace
		err = deployNamespace(&cluster)
		if err != nil {
			log.GetLogger().Infof("failed in deployNamespace error %s", err)
			return err
		}
	}

	// process clusters
	for _, resource := range m["Cluster"] {
		var cert *common.Certificate
		var cluster common.Cluster
		var err error
		var cPath string
		var rootCertPem, rootCakey, remoteSecret string
		var clusterChainPem, clusterCertPem, clusterKeyPem string

		log.GetLogger().Debugf("cluster resoure %s", resource)
		err = yaml.Unmarshal(resource, &cluster)
		if err != nil {
			log.GetLogger().Infof("failed parse cluster %s", resource)
			return err
		}

		//check for cluster certificate config
		cert, cPath, err = getCertFromMeshID(cluster.Spec.MeshID, m)
		if err != nil {
			log.GetLogger().Infof("failed get cert config for cluster %s", resource)
			return err
		}

		clusterCertPath := path.Join(cPath, cluster.Metadata.Name)
		err = os.MkdirAll(clusterCertPath, 0755)
		if err != nil {
			log.GetLogger().Infof("failed mkdir cluser cert path %s for cluster %s", clusterCertPath, resource)
			return err
		}

		rootCertPem = path.Join(cPath, "root-cert.pem")
		rootCakey = path.Join(cPath, "root-ca.key")
		clusterChainPem = path.Join(clusterCertPath, "cert-chain.pem")
		clusterKeyPem = path.Join(clusterCertPath, "ca-key.pem")
		clusterCertPem = path.Join(clusterCertPath, "ca-cert.pem")

		log.GetLogger().Infof("rootCertPem %s rootCakey %s  clusterChainPem %s clusterKeyPem %s clusterCertPem %s", rootCertPem, rootCakey, clusterChainPem, clusterKeyPem, clusterCertPem)
		// check for cluster certificate is generated
		_, err = os.Stat(path.Join(clusterCertPath, "cert-chain.pem"))
		if err != nil {
			// generate cert for cluster
			intermediate := cluster.Metadata.Name + "." + cert.Spec.SanSuffix
			args := []string{
				"step",
				"certificate",
				"create",
				intermediate,
				clusterCertPem,
				clusterKeyPem,
				"--ca",
				rootCertPem,
				"--ca-key",
				rootCakey,
				"--profile",
				"intermediate-ca",
				"--no-password",
				"--insecure",
				"--san",
				intermediate,
				"--not-after",
				cert.Spec.ValidityHours + "h",
				"--kty",
				"RSA",
			}
			log.GetLogger().Debugf("generate cluster certs in path %s args:%s", clusterCertPath, args)
			stepcerts.GenStepCerts(args)
			err = combineFiles(clusterCertPem, rootCertPem, clusterChainPem)
			if err != nil {
				log.GetLogger().Infof("failed to combine %s %s to chained cert %s", clusterCertPem, rootCertPem, clusterChainPem)
				return err
			}
		}

		// install cluster certs
		err = installClusterCerts(&cluster, rootCertPem, rootCakey, clusterChainPem, clusterCertPem, clusterKeyPem)
		if err != nil {
			log.GetLogger().Infof("failed installClusterCerts error %s", err)
			return err
		}

		// deploy controlplane for cluster
		deployControlPlane(&cluster)

		// deploy eastwest gateway for cluster
		deployEastWestGateway(&cluster)

		// deploy expose service for cluster
		err = deployExposeService(&cluster)
		if err != nil {
			log.GetLogger().Infof("failed in deployExposeService error %s", err)
			return err
		}

		// Check for Rafay kubeconfig - server dns conatins rafay.dev
		rafayCfg, err := common.CheckForRafayKubeConfig(&cluster)
		if err == nil {
			remoteSecret, err = cretaeRafayRemoteSecret(&cluster, rafayCfg)
			if err != nil {
				log.GetLogger().Infof("failed in create rafay remote secret for %s error %s", cluster.Metadata.Name, err)
				return err
			}
		} else {
			// use istioctl to genereate secrets
			remoteSecret, err = createRemoteSecret(&cluster)
			if err != nil {
				log.GetLogger().Infof("failed in create remote secret for %s error %s", cluster.Metadata.Name, err)
				return err
			}
		}

		// Enable Endpoint Discovery. Deploy remoteSecret to all other clusters
		err = enableEndpontDiscovery(m, cluster.Metadata.Name, remoteSecret)
		if err != nil {
			log.GetLogger().Infof("failed in enableEndpontDiscovery error %s", err)
			return err
		}

		// deploy helloworld
		if cluster.Spec.InstallHelloWorld == "true" {
			err = deployHelloWorld(&cluster)
			if err != nil {
				log.GetLogger().Infof("failure in deployHelloWorld error %s", err)
			}
		}
	}

	return nil
}

// describe function createRemoteSecret
// createRemoteSecret is the function to create the remote secret for the given cluster
// cluster is the cluster to create the remote secret for
// returns string, error
// string is the path to the created remote secret YAML file
// error is the error encountered
func createRemoteSecret(cluster *common.Cluster) (string, error) {
	istioCtlDir, err := istioctl.DownloadIsctioCtl(cluster.Spec.Version, "")
	if err != nil {
		log.GetLogger().Debugf("failed to get istioctl binary info for version %s error", cluster.Spec.Version, err)
		return "", err
	}

	// istioctl x create-remote-secret --context="cluster1" --name=cluster1
	output, err := istioctl.RunIstioCtlCmd(cluster.Spec.Version, istioCtlDir, "create-remote-secret", cluster.Spec.KubeconfigFile, cluster.Spec.Context, "", cluster.Metadata.Name)
	if err != nil {
		log.GetLogger().Debugf("failed to deploy eastwest-gateway using istioctl error %s", err)
		return "", err
	}

	if output.Len() > 0 {
		output.Bytes()
		remoteSecretFile, err := createClusterYaml(cluster, "remote-secret.yaml", &output)
		if err != nil {
			return "", err
		}
		log.GetLogger().Debugf("istioctl remote-secret %s ", remoteSecretFile)
		return remoteSecretFile, nil
	}

	return "", fmt.Errorf("failed to create remote secret for %s", cluster.Metadata.Name)
}

// describe function enableEndpontDiscovery
// enableEndpontDiscovery is the function to enable endpoint discovery for the given cluster
// m is the map of resources
// srcClusterName is the name of the source cluster whose remote secret is to be applied to other clusters
// remoteSecret is the path to the remote secret YAML file
// returns error
// error is the error encountered
func enableEndpontDiscovery(m map[string][][]byte, srcClusterName, remoteSecret string) error {
	// loop through all clusters
	for _, resource := range m["Cluster"] {
		var cluster common.Cluster
		var err error

		log.GetLogger().Debugf("cluster resoure %s", resource)
		err = yaml.Unmarshal(resource, &cluster)
		if err != nil {
			log.GetLogger().Infof("failed parse cluster %s", resource)
			return err
		}

		// skip source cluster
		if srcClusterName == cluster.Metadata.Name {
			continue
		}

		// kubectl apply -f remoteSecret.yaml --context="remoteCluster"
		// apply the remote secret to the cluster using above kubectl syntax
		applyRemoteSecret := []string{
			"--kubeconfig",
			cluster.Spec.KubeconfigFile,
			"--context=" + cluster.Spec.Context,
			"-n",
			"istio-system",
			"apply",
			"-f",
			remoteSecret,
		}

		// if dry run, return
		if DryRun {
			return nil
		}

		// execute the above kubectl command
		output, err := utils.KubeCtlCmd(applyRemoteSecret)
		if err != nil {
			log.GetLogger().Debugf("failed to apply remote secret of %s to %s error %s", srcClusterName, cluster.Metadata.Name, err)
			return err
		}
		if output.Len() > 0 {
			log.GetLogger().Debugf("apply remote secret of %s to %s output %s", srcClusterName, cluster.Metadata.Name, output.String())
		}
	}

	return nil
}

// describe function cretaeRafayRemoteSecret
// cretaeRafayRemoteSecret is the function to create the remote secret for the given cluster
// cluster is the cluster to create the remote secret for
// rafayCfg is the rafay kubeconfig
// returns string, error
// string is the path to the created remote secret YAML file
// error is the error encountered
func cretaeRafayRemoteSecret(cluster *common.Cluster, rafayCfg *clientcmdapiv1.Config) (string, error) {
	es := new(bytes.Buffer)
	// prepare the template values
	valData := common.RafayRemoteSecreteTempllateValues{
		ClusterName: cluster.Metadata.Name,
		ServerCA:    b64.StdEncoding.EncodeToString(rafayCfg.Clusters[0].Cluster.CertificateAuthorityData),
		Server:      rafayCfg.Clusters[0].Cluster.Server,
		User:        rafayCfg.AuthInfos[0].Name,
		UserCert:    b64.StdEncoding.EncodeToString(rafayCfg.AuthInfos[0].AuthInfo.ClientCertificateData),
		UserKey:     b64.StdEncoding.EncodeToString(rafayCfg.AuthInfos[0].AuthInfo.ClientKeyData),
	}
	// perpare the remote secret YAML from template
	err := fixtures.RafayRemoteSecrete.Execute(es, struct {
		Values common.RafayRemoteSecreteTempllateValues
	}{
		valData,
	})
	if err != nil {
		log.GetLogger().Debugf("fixtures.RafayRemoteSecrete failed to execute")
		return "", err
	}

	// create the YAML file for the remote secret
	remoteSecret, err := createClusterYaml(cluster, "remote-secret.yaml", es)
	if err != nil {
		return "", err
	}

	return remoteSecret, nil
}

func (o *ApplyOnAllOptions) AddFlags(cmd *cobra.Command) {
	AddYamlConfigFlagVar(&o.YamlConfigFilePath, cmd, "Use this flag to create/update all-pipeline using a YAML file")
	cmd.Flags().BoolVar(&o.DryRun, DryRunFlag, false, "Dry Run")
}

func NewApplyOnAllOptions(logger log.Logger) CmdOptions {
	options := new(ApplyOnAllOptions)
	options.logger = logger
	return options
}

// describe function deployHelloWorld
// deployHelloWorld is the function to deploy the helloworld for the given cluster
// cluster is the cluster to deploy the helloworld for
// returns error
// error is the error encountered
func deployHelloWorld(cluster *common.Cluster) error {
	hw := new(bytes.Buffer)
	// prepare the template values
	valData := common.TempllateValues{
		MeshID:      cluster.Spec.MeshID,
		ClusterName: cluster.Metadata.Name,
	}
	// perpare the helloworld YAML from template
	err := fixtures.HelloWorld.Execute(hw, struct {
		Values common.TempllateValues
	}{
		valData,
	})
	if err != nil {
		log.GetLogger().Debugf("fixtures.HelloWorld failed to execute")
		return err
	}

	// create the YAML file for the helloworld
	hwFile, err := createClusterYaml(cluster, "helloworld.yaml", hw)
	if err != nil {
		return err
	}

	// if dry run, return
	if DryRun {
		return nil
	}

	// kubectl apply -f helloworld.yaml
	// deploy the helloworld using kubectl uses above kubectl syntax
	applyClusterNamespace := []string{
		"--kubeconfig",
		cluster.Spec.KubeconfigFile,
		"--context=" + cluster.Spec.Context,
		"-n",
		constants.HELLO_WORLD_NAMESPACE,
		"apply",
		"-f",
		hwFile,
	}

	output, err := utils.KubeCtlCmd(applyClusterNamespace)
	if err != nil {
		log.GetLogger().Debugf("failed to apply cluster HelloWorld from %s error %s", hwFile, err)
		return err
	}

	if output.Len() > 0 {
		log.GetLogger().Debugf("kubectl output to apply HelloWorld from %s output %s", hwFile, output.String())
	}

	return nil
}
