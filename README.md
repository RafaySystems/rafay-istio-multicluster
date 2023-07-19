# rafay-istio-multi-cluster CLI (ristioctl) 
## Overview
The ristioctl CLI is a command line tool that allows you to manage your Istio Multi-Cluster deployments. The CLI is built using the [cobra]

## Installation
The CLI can be compiled using the following methods:
make build

## Usage
```
bin/ristioctl_darwin_amd64 -h
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
