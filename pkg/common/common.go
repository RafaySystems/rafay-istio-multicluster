package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/RafaySystems/rafay-istio-multicluster/pkg/log"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/yaml"
)

type Meta struct {
	Name        string            `yaml:"name"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type ResourceType struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   Meta   `yaml:"metadata"`
}

type CertificateSpec struct {
	ValidityHours string `yaml:"validityHours,omitempty"`
	SanSuffix     string `yaml:"sanSuffix,omitempty"`
	MeshID        string `yaml:"meshID,omitempty"`
	Folder        string `yaml:"folder,omitempty"`
	Password      bool   `yaml:"password,omitempty"`
}
type Certificate struct {
	ApiVersion string          `yaml:"apiVersion,omitempty"`
	Kind       string          `yaml:"kind"`
	Metadata   Meta            `yaml:"metadata"`
	Spec       CertificateSpec `yaml:"spec"`
}

type ClusterSpec struct {
	KubeconfigFile    string `yaml:"kubeconfigFile,omitempty"`
	Context           string `yaml:"context"`
	MeshID            string `yaml:"meshID"`
	Version           string `yaml:"version"`
	InstallHelloWorld string `yaml:"installHelloWorld"`
}

type Cluster struct {
	ApiVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   Meta        `yaml:"metadata"`
	Spec       ClusterSpec `yaml:"spec"`
}

type TempllateValues struct {
	MeshID      string
	ClusterName string
}

type RafayRemoteSecreteTempllateValues struct {
	ClusterName string
	ServerCA    string
	Server      string
	User        string
	UserCert    string
	UserKey     string
}

// describe function CheckForRafayKubeConfig
// CheckForRafayKubeConfig checks if the Rafay ZTKA kubeconfig file
// cluster is the cluster object
// returns error
func CheckForRafayKubeConfig(cluster *Cluster) (*clientcmdapiv1.Config, error) {
	var rafayConfig clientcmdapiv1.Config

	kubeconfig, err := parseKubeConfig(cluster.Spec.KubeconfigFile)
	if err != nil {
		return nil, err
	}

	log.GetLogger().Debugf("CheckForRafayKubeConfig kubeconfig %s", kubeconfig)

	for _, cntx := range kubeconfig.Contexts {
		log.GetLogger().Debugf("CheckForRafayKubeConfig cntx %s", cntx)
		if cntx.Name != cluster.Spec.Context {
			continue
		}

		log.GetLogger().Debugf("CheckForRafayKubeConfig matched cntx %s", cntx)

		rafayConfig.Contexts = append(rafayConfig.Contexts, cntx)
		rafayConfig.CurrentContext = cntx.Name

		//find the cluster
		for _, cst := range kubeconfig.Clusters {
			if cst.Name == cntx.Context.Cluster {
				if strings.Contains(cst.Cluster.Server, "rafay-edge.net") || strings.Contains(cst.Cluster.Server, "rafay.dev") {
					rafayConfig.Clusters = append(rafayConfig.Clusters, cst)
					break
				} else {
					return nil, fmt.Errorf("server does not contain rafay-edge.net")
				}
			}
		}

		log.GetLogger().Debugf("CheckForRafayKubeConfig cntx.Context.AuthInfo %s kubeconfig.AuthInfos %s", cntx.Context.AuthInfo, kubeconfig.AuthInfos)
		// find the user
		for _, ati := range kubeconfig.AuthInfos {
			log.GetLogger().Debugf("CheckForRafayKubeConfig ati.Name %s", ati.Name)
			if ati.Name == cntx.Context.AuthInfo {
				rafayConfig.AuthInfos = append(rafayConfig.AuthInfos, ati)
				break
			}
		}

	}

	return &rafayConfig, nil
}

// describe function parseKubeConfig
// parseKubeConfig parses the kubeconfig file
// fileName is the kubeconfig file name
// returns *clientcmdapiv1.Config, error
func parseKubeConfig(fileName string) (*clientcmdapiv1.Config, error) {
	var kubeconfig clientcmdapiv1.Config

	yb, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	jb, err := yaml.YAMLToJSON(yb)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jb, &kubeconfig)
	if err != nil {
		return nil, err
	}

	return &kubeconfig, nil

}
