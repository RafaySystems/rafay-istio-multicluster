package cmd

import (
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/commands"
	"github.com/spf13/cobra"
)

func newApplyFileCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the pipeline command
	cmd := &cobra.Command{
		Use:     "apply -f config.yaml",
		Aliases: []string{"ppl"},
		Short:   "apply for creating or updating istio mesh across multiple Kubernetes clusters",
		Long:    `apply for creating or updating istio mesh across multiple Kubernetes clusters`,
		Example: `
  Using config file
    istio-mc apply -f config.yaml

---------------------------
# Sample config.yaml
---------------------------
apiVersion: ristioctl.k8smgmt.io/v3
kind: Certificate
metadata:
  name: ristioctl-certs
spec:
  validityHours: 87600
  password: false
  # Subject Alternative Name Suffix
  sanSuffix: istio.io
  meshID: usmesh
---
apiVersion: ristioctl.k8smgmt.io/v3
kind: Cluster
metadata:
  name: cluster1
spec:
  kubeconfigFile: kubeconfig.yaml
  context: cluster1-context
  meshID: usmesh
  version: "1.18.0"
---
apiVersion: ristioctl.k8smgmt.io/v3
kind: Cluster
metadata:
  name: cluster2
spec:
  kubeconfigFile: kubeconfig.yaml
  context: cluster2-context
  meshID: usmesh
  version: "1.18.0"
		`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
