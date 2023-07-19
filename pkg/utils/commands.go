package utils

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/RafaySystems/rafay-istio-multicluster/pkg/log"
)

var (
	Expose string = "false" // Expose hidden commands
)

func Hidden() bool {
	return Expose == "false"
}

func ShowDeprecatedMessageOutput(msg, output string) {
	if output == "table" {
		fmt.Println(msg)
	}
}

func KubeCtlCmd(args []string) (bytes.Buffer, error) {
	var (
		out bytes.Buffer
		er  bytes.Buffer
	)

	log.GetLogger().Debugf("Run kubectl command %s ", args)
	// We need a variable executable here hence using nosec
	// #nosec
	command := exec.Command("kubectl", args...)
	command.Stdout = &out
	command.Stderr = &er
	err := command.Run()
	if err != nil {
		return out, err
	}

	return out, nil
}
