package istioctl

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/RafaySystems/rafay-istio-multicluster/pkg/constants"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/log"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/utils"
	"github.com/layer5io/meshkit/errors"
	"github.com/mitchellh/go-homedir"
)

const (
	platform = runtime.GOOS
	arch     = runtime.GOARCH
)

var (
	// ErrRunIstioCtlCmdCode represents the errors which are generated
	// during fetch manifest process
	ErrRunIstioCtlCmdCode = "1000"

	// ErrIstioctlNotFoundCode implies istioctl couldn't be found anywhere
	// on the fs
	ErrIstioctlNotFoundCode = "1001"

	// ErrGettingIstioReleaseCode implies failure while failing istio release
	// bundle
	ErrGettingIstioReleaseCode = "1002"

	// during unzip process
	ErrUnzipFileCode = "1003"

	// ErrTarXZFCode represents the errors which are generated
	// during decompressing and extracting tar.gz file
	ErrTarXZFCode = "1004"

	// ErrUnsupportedPlatformCode implies unavailbility of Istio on the
	// requested platform
	ErrUnsupportedPlatformCode = "1005"

	// ErrDownloadingTarCode implies error while downloading istio tar
	ErrDownloadingTarCode = "1006"

	// ErrUnpackingTarCode implies error while unpacking istio release
	// bundle tar
	ErrUnpackingTarCode = "1007"

	// ErrMakingBinExecutableCode implies error while makng istioctl executable
	ErrMakingBinExecutableCode = "1008"

	// ErrLoadNamespaceCode implies error while finding namespace
	ErrLoadNamespaceCode = "1009"

	// ErrIstioctlNotFound implies istioctl was not found locally
	ErrIstioctlNotFound = errors.New(ErrIstioctlNotFoundCode, errors.Alert, []string{"Unable to find Istioctl"}, []string{}, []string{}, []string{})

	// ErrUnsupportedPlatform represents runtime platform is
	// unsupported
	ErrUnsupportedPlatform = errors.New(ErrUnsupportedPlatformCode, errors.Alert, []string{"requested platform is not supported by Istio"}, []string{"Istio only supports Windows, Linux and Darwin"}, []string{}, []string{""})

	//downloadLocation = os.TempDir()
)

// ErrIstioctlNotFound implies istioctl was not found locally

// ErrRunIstioCtlCmd is the error for mesh port forward
func ErrRunIstioCtlCmd(err error, des string) error {
	return errors.New(ErrRunIstioCtlCmdCode, errors.Alert, []string{"Error running istioctl command"}, []string{err.Error()}, []string{"Corrupted istioctl binary", "Command might be invalid"}, []string{})
}

// ErrGettingIstioRelease is the error when the yaml unmarshal fails
func ErrGettingIstioRelease(err error) error {
	return errors.New(ErrGettingIstioReleaseCode, errors.Alert, []string{"Error occured while fetching Istio release artifacts"}, []string{err.Error()}, []string{}, []string{})
}

// ErrUnzipFile is the error for unzipping the file
func ErrUnzipFile(err error) error {
	return errors.New(ErrUnzipFileCode, errors.Alert, []string{"Error while unzipping"}, []string{err.Error()}, []string{"File might be corrupt"}, []string{})
}

// ErrTarXZF is the error for unzipping the file
func ErrTarXZF(err error) error {
	return errors.New(ErrTarXZFCode, errors.Alert, []string{"Error while extracting file"}, []string{err.Error()}, []string{"/The gzip might be corrupt"}, []string{})
}

// ErrDownloadingTar is the error when tar download fails
func ErrDownloadingTar(err error) error {
	return errors.New(ErrDownloadingTarCode, errors.Alert, []string{"Error occured while downloading Istio tar"}, []string{err.Error()}, []string{}, []string{})
}

// ErrUnpackingTar is the error when tar unpack fails
func ErrUnpackingTar(err error) error {
	return errors.New(ErrUnpackingTarCode, errors.Alert, []string{"Error occured while unpacking tar"}, []string{err.Error()}, []string{}, []string{})
}

// ErrMakingBinExecutable occurs when istioctl binary couldn't be made
// executable
func ErrMakingBinExecutable(err error) error {
	return errors.New(ErrMakingBinExecutableCode, errors.Alert, []string{"Error while making istioctl an executable"}, []string{err.Error()}, []string{}, []string{})
}

// ErrLoadNamespace implies error while finding namespace
func ErrLoadNamespace(err error, str string) error {
	return errors.New(ErrLoadNamespaceCode, errors.Alert, []string{"Error while labeling namespace:", str}, []string{err.Error()}, []string{}, []string{})
}

// describe function KubeCtlGetAllIstioSystem
// Runs kubectl get all -n istio-system
func KubeCtlGetAllIstioSystem(kubeconfig, context string) error {
	getCmd := []string{
		"--kubeconfig",
		kubeconfig,
		"--context=" + context,
		"-n",
		"istio-system",
		"get",
		"all",
	}

	output, err := utils.KubeCtlCmd(getCmd)
	if err != nil {
		log.GetLogger().Debugf("failed to KubeCtlGetAllIstioSystem %s", err)
		return err
	}

	if output.Len() > 0 {
		log.GetLogger().Debugf("kubectl get all -n istio-system \n %s", output.String())
	}

	return nil
}

// describe function RunIstioCtlCmd
// Runs istioctl install/uninstall/create-remote-secret commands
// version is the istio version
// dirName is the directory where istioctl is downloaded
// istioCmd is the istioctl command to run
// kubeconfig is the kubeconfig file
// context is the context to use
// file is the file to use
// name is the name to use
func RunIstioCtlCmd(version, dirName, istioCmd, kubeconfig, context, file, name string) (bytes.Buffer, error) {
	var (
		out     bytes.Buffer
		er      bytes.Buffer
		execCmd []string
	)

	log.GetLogger().Info("Running istioctl...")

	// get the istioctl executable path
	Executable, err := getExecutable(version, dirName)
	if err != nil {
		return out, ErrRunIstioCtlCmd(err, err.Error())
	}

	switch istioCmd {
	case "install":
		execCmd = []string{"--kubeconfig=" + kubeconfig, "--context=" + context, istioCmd, "-y", "-f", file, "--readiness-timeout", "15s"}
	case "uninstall":
		execCmd = []string{"--kubeconfig=" + kubeconfig, "--context=" + context, "x", istioCmd, "--purge", "-y"}
	case "create-remote-secret":
		execCmd = []string{"--kubeconfig=" + kubeconfig, "--context=" + context, "x", istioCmd, "--name=" + name}
	}

	log.GetLogger().Debugf("RunIstioCtlCmd execCmd %s %s ", Executable, execCmd)

	command := exec.Command(Executable, execCmd...)
	command.Stdout = &out
	command.Stderr = &er
	if istioCmd == "create-remote-secret" {
		log.GetLogger().Debugf("running create-remote-secret %s", command.String())
		err = command.Run()
		if err != nil {
			return out, ErrRunIstioCtlCmd(err, er.String())
		}
	} else {
		// start the command and watch istion-system namespace for resources
		command.Start()
		closeCh := make(chan struct{})

		ticker := time.NewTicker(5 * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					KubeCtlGetAllIstioSystem(kubeconfig, context)
				case <-closeCh:
					ticker.Stop()
					return
				}
			}
		}()
		err = command.Wait()
		closeCh <- struct{}{}
		if err != nil {
			return out, ErrRunIstioCtlCmd(err, er.String())
		}
	}

	return out, nil
}

// describe function getExecutable
// getExecutable returns the path to istioctl executable
// release is the istio version
// dirName is the directory where istioctl is downloaded
// returns the path to istioctl executable
// returns error
func getExecutable(release, dirName string) (string, error) {
	binaryName := generatePlatformSpecificBinaryName("istioctl", platform)

	log.GetLogger().Info("Using istioctl from the downloaded release bundle...")
	executable := path.Join(dirName, "bin", binaryName)
	if _, err := os.Stat(executable); err == nil {
		return executable, nil
	}

	log.GetLogger().Info("Done")
	return "", ErrIstioctlNotFound
}

func generatePlatformSpecificBinaryName(binName, platform string) string {
	if platform == "windows" && !strings.HasSuffix(binName, ".exe") {
		return binName + ".exe"
	}

	return binName
}

func DownloadIsctioCtl(version string, folder string) (string, error) {
	// Fetch and/or return the path to downloaded and extracted release bundle
	dirName, err := getIstioRelease(version, folder)
	if err != nil {
		// ErrGettingIstioRelease
		return "", ErrGettingIstioRelease(err)
	}

	return dirName, nil
}

// describe function getIstioRelease
// getIstioRelease returns the path to istioctl executable
// release is the istio version
// folder is the directory where istioctl is downloaded
// returns the path to istioctl executable
// returns error
func getIstioRelease(release, folder string) (string, error) {
	var dPath string
	var artifactPath string
	releaseName := fmt.Sprintf("istio-%s", release)

	if folder == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dPath = path.Join(home, constants.CLI_HOME_FOLDER, "istio")
		err = os.MkdirAll(dPath, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else {
		dPath = folder
	}

	artifactPath = path.Join(dPath, releaseName)
	log.GetLogger().Debugf("Looking for artifacts of requested version %s ...", artifactPath)

	_, err := os.Stat(artifactPath)
	if err == nil {
		log.GetLogger().Debugf("return artifactPath %s", artifactPath)
		return artifactPath, nil
	}
	log.GetLogger().Debugf("Artifacts not found...")

	log.GetLogger().Debugf("Downloading requested istio version artifacts... releaseName %s, release %s", releaseName, release)
	res, err := downloadTar(releaseName, release)
	if err != nil {
		return "", ErrGettingIstioRelease(err)
	}

	err = extractTar(dPath, res)
	if err != nil {
		return "", ErrGettingIstioRelease(err)
	}

	return artifactPath, nil
}

// describe function downloadTar
// downloadTar downloads the istio release bundle
// releaseName is the istio release name
// release is the istio version
// returns the http response
// returns error
func downloadTar(releaseName, release string) (*http.Response, error) {
	log.GetLogger().Debugf("downloadTar releaseName %s release %s", releaseName, release)
	url := "https://github.com/istio/istio/releases/download"

	switch platform {
	case "darwin":
		url = fmt.Sprintf("%s/%s/%s-osx.tar.gz", url, release, releaseName)
	case "windows":
		url = fmt.Sprintf("%s/%s/%s-win.zip", url, release, releaseName)
	case "linux":
		url = fmt.Sprintf("%s/%s/%s-%s-%s.tar.gz", url, release, releaseName, platform, arch)
	default:
		return nil, ErrUnsupportedPlatform
	}

	log.GetLogger().Debugf("isto download url %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrDownloadingTar(err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, ErrDownloadingTar(fmt.Errorf("status is not http.StatusOK"))
	}

	return resp, nil
}

// describe function extractTar
// extractTar extracts the istio release bundle
// dPath is the directory where istioctl is downloaded
// res is the http response
// returns error
func extractTar(dPath string, res *http.Response) error {
	log.GetLogger().Debugf("extractTar to path %s", dPath)
	// Close the response body
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	switch platform {
	case "darwin":
		fallthrough
	case "linux":
		if err := tarxzf(dPath, res.Body); err != nil {
			//ErrExtracingFromTar
			return ErrUnpackingTar(err)
		}
	case "windows":
		if err := unzip(dPath, res.Body); err != nil {
			return ErrUnpackingTar(err)
		}
	}

	return nil
}

// describe function tarxzf
// tarxzf extracts the tar.gz file
// location is the directory where istioctl is downloaded
// stream is the http response body
// returns error
func tarxzf(location string, stream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(stream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return ErrTarXZF(err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// File traversal is required to store the extracted manifests at the right place
			if err := os.MkdirAll(path.Join(location, header.Name), 0750); err != nil {
				return ErrTarXZF(err)
			}
		case tar.TypeReg:
			// File traversal is required to store the extracted manifests at the right place
			outFile, err := os.Create(path.Join(location, header.Name))
			if err != nil {
				return ErrTarXZF(err)
			}
			// Trust istioctl tar
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return ErrTarXZF(err)
			}
			if err = outFile.Close(); err != nil {
				return ErrTarXZF(err)
			}

			if header.FileInfo().Name() == "istioctl" {
				// istioctl binary needs to be executable
				if err = os.Chmod(outFile.Name(), 0750); err != nil {
					return ErrMakingBinExecutable(err)
				}
			}

		default:
			return ErrTarXZF(err)
		}
	}

	return nil
}

// describe function unzip
// unzip extracts the zip file
// location is the directory where istioctl is downloaded
// zippedContent is the http response body
// returns error
func unzip(location string, zippedContent io.Reader) error {
	// Keep file in memory: Approx size ~ 50MB
	// TODO: Find a better approach
	zipped, err := ioutil.ReadAll(zippedContent)
	if err != nil {
		return ErrUnzipFile(err)
	}

	zReader, err := zip.NewReader(bytes.NewReader(zipped), int64(len(zipped)))
	if err != nil {
		return ErrUnzipFile(err)
	}

	for _, file := range zReader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return ErrUnzipFile(err)
		}
		defer func() {
			if err := zippedFile.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		// need file traversal to place the extracted files at the right place, hence
		extractedFilePath := path.Join(location, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(extractedFilePath, file.Mode()); err != nil {
				return ErrUnzipFile(err)
			}
		} else {
			// we need a variable path hence,
			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				return ErrUnzipFile(err)
			}
			defer func() {
				if err := outputFile.Close(); err != nil {
					fmt.Println(err)
				}
			}()

			// Trust istio zip hence,
			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return ErrUnzipFile(err)
			}
		}
	}

	return nil
}
