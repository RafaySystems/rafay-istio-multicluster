# rafay-istio-multi-cluster CLI (ristioctl) 
## Overview
The ristioctl CLI is a command line tool that allows you to manage your Istio Multi-Cluster deployments. The CLI is built using the [cobra]

## Build
```
make build
```

## Usage
```
>> bin/ristioctl_darwin_amd64 -h
A CLI tool to manage istio mesh across multiple Kubernetes clusters.

Usage:
  ristioctl [command]

Available Commands:
  apply       apply for creating or updating istio mesh across multiple Kubernetes clusters
  completion  Generate the autocompletion script for the specified shell
  delete      Delete various resources in Console
  download    Download various resources in Console
  help        Help about any command
  version     Displays version of the CLI utility

Flags:
  -d, --debug     Enable debug logs
  -h, --help      help for ristioctl
  -v, --verbose   Verbose mode. A lot more information output.
      --wait      Wait for the operation to complete

Use "ristioctl [command] --help" for more information about a command.
```

## Apply
```
>> bin/ristioctl_darwin_amd64 apply -f example/mesh1.yaml
```

## Delete
```
>> bin/ristioctl_darwin_amd64 delete -f example/mesh1.yaml
```

## Goal
The goal of writing this CLI tool is to automate the setup of a Multi-primary Istio Mesh using the Gitops pipeline. Configuring a multi-cluster service mesh can present several challenges. Here are some common challenges:

- Complexity of Configuration: Managing a service mesh configuration across multiple clusters can be complex. 
- Consistent Configuration: Ensuring consistent configuration across multiple clusters is crucial for the proper functioning of the service mesh.
- Network Connectivity: Establishing network connectivity between clusters is vital for a multi-cluster service mesh. It requires setting up secure communication channels, often across public or hybrid cloud environments.
- Service Discovery: Service discovery becomes more challenging in a multi-cluster environment. Ensuring that services in one cluster can discover and communicate with services in other clusters requires careful configuration and coordination.
- Monitoring and Troubleshooting: Monitoring and troubleshooting a multi-cluster service mesh configuration can be complex due to the increased number of components and the distributed nature of the infrastructure.s.

To address these challenges, adopting infrastructure-as-code (IaC) approaches for configuration management and leveraging automation tools for consistent deployments is recommended. At Rafay, we developed this CLI tool to deploy multi-cluster Istio service mesh in our internal environments. This tool uses the “Multi-Primary on different networks”, topology as described in the [Istio documentation](https://istio.io/latest/docs/setup/install/multicluster/multi-primary_multi-network/).



The CLI consumes a simple configuration as shown below to set up a multi-cluster service mesh.
```
>> cat examples/mesh.yaml
apiVersion: istiomc.k8smgmt.io/v3
kind: Certificate
metadata:
  name: istiomc-certs
spec:
  validityHours: 2190
  password: false
  sanSuffix: istio.io # Subject Alternative Name Suffix
  meshID: uswestmesh
---
apiVersion: istiomc.k8smgmt.io/v3
kind: Cluster
metadata:
  name: cluster1
spec:
  kubeconfigFile: "kubeconfig-istio-demo.yaml"
  context: cluster1
  meshID: uswestmesh
  version: "1.18.0"
  installHelloWorld: true #deploy sample HelloWorld application
---
apiVersion: istiomc.k8smgmt.io/v3
kind: Cluster
metadata:
  name: cluster2
spec:
  kubeconfigFile: "kubeconfig-istio-demo.yaml"
  context: cluster2
  meshID: uswestmesh
  version: "1.18.0"
  installHelloWorld: true   #deploy sample HelloWorld application
```

## Certificate:
```kind: Certificate:``` The CLI establishes trust between all clusters in the mesh using this configuration. It will generate and deploy distinct certificates for each cluster. All cluster certificates are issued by the same root certificate authority (CA). Internally the CLI uses the [step-ca](https://smallstep.com/docs/step-ca/) tool.

## Cluster: 
```kind: Cluster:``` The CLI identifies the clusters in the mesh based on this configuration. Specify a unique name, the istio version, the kubeconfig file, and the context for each cluster. The CLI tool will generate all the required configurations and deploy them to each cluster. The CLI internally takes care of the following steps.
- Configure Trust across all clusters in the mesh.
- Deploy Istio into the clusters
- Deploy east-west gateway into the clusters
- Expose services in the clusters
- Enable cross-cluster service discovery using Rafay ZTKA-based secure channel.

