package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/RafaySystems/rafay-istio-multicluster/cmd"
	"github.com/RafaySystems/rafay-istio-multicluster/internal/fixtures"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/audit"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/constants"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/versioninfo"
	"github.com/mitchellh/go-homedir"
)

// cliversion is the version of the CLI
// clitime is the time the CLI was built
// cliarch is the architecture the CLI was built for
var (
	cliversion string
	clitime    string
	cliarch    string
	buildNum   string
)

func main() {
	// load the yaml templates from internal/fixtures/data folder
	err := fixtures.Load()
	if err != nil {
		fmt.Println("failed to load fixtures")
		fmt.Println(err)
		os.Exit(1)
	}

	// create the audit folder, default to users home directory
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("failed to get home dir ")
		fmt.Println(err)
		os.Exit(1)
	}

	// create the audit folder with a timestamp. Each cli execution creates a new folder
	timestamp := "audit-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	audit.AuditFolder = path.Join(home, constants.CLI_HOME_FOLDER, "istio", timestamp)
	err = os.MkdirAll(audit.AuditFolder, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cliarch == "" {
		cliarch = runtime.GOOS + "/" + runtime.GOARCH
	}
	versioninfo.Init(cliversion, buildNum, clitime, cliarch)

	// Execute the root command
	cmd.Execute()
}
