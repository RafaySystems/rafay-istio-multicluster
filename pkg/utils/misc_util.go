package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type ErrorDetails struct {
	ErrorCode string `json:"error_code,omitempty" yaml:"error_code,omitempty"`
	Detail    string `json:"detail,omitempty" yaml:"detail,omitempty"`
	Info      string `json:"info,omitempty" yaml:"info,omitempty"`
}

type RafayErrorMessage struct {
	StatusCode int            `json:"status_code,omitempty" yaml:"status_code,omitempty"`
	Details    []ErrorDetails `json:"details,omitempty" yaml:"details,omitempty"`
}

type LabelStruct struct {
	Items map[string]string `yaml:"LABELS"`
}

func GetUserHome() string {
	homeEnvVariable := "HOME"
	if runtime.GOOS == "windows" {
		homeEnvVariable = "USERPROFILE"
	}
	return os.Getenv(homeEnvVariable)
}

func FormatYamlMessage(data interface{}) (string, error) {
	var ret string
	bArr, err := yaml.Marshal(data)
	if err == nil {
		ret = string(bArr)
	}
	return ret, err
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func FullPath(parentFile, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(filepath.Dir(parentFile), path)
}

func FullPaths(parentFile, path string) string {
	allPaths := strings.Split(path, ",")
	if len(allPaths) <= 1 {
		return FullPath(parentFile, path)
	}
	allFullPaths := make([]string, len(allPaths))
	for i, aPath := range allPaths {
		allFullPaths[i] = FullPath(parentFile, aPath)
	}
	return strings.Join(allFullPaths, ",")
}

func GetAsString(i interface{}) string {
	if i == nil {
		return ""
	}
	return i.(string)
}

func GetAsMap(array []string) map[string]string {
	asMap := make(map[string]string)
	for _, entry := range array {
		asMap[entry] = ""
	}
	return asMap
}

func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func PrettyPrint(responseStr string) {
	if !gjson.Valid(responseStr) {
		fmt.Println(responseStr)
		return
	}
	result := gjson.Get(responseStr, "@pretty")
	fmt.Println(result.String())
}

func ExpandFile(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		abs, err := filepath.Abs(path)
		return abs, err
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

func ProcessValueFiles(valuesFilePath string) (string, error) {
	allValuesFileNames := ""
	allValuesFiles := strings.Split(valuesFilePath, ",")
	for _, valuesFile := range allValuesFiles {
		if valuesFile != "" {
			absFile, err := ExpandFile(valuesFile)
			if err != nil {
				return "", fmt.Errorf("values file %s does not exist error %s", valuesFile, err.Error())
			}
			if !FileExists(absFile) {
				return "", fmt.Errorf("values file %s does not exist", valuesFile)
			}
			allValuesFileNames = allValuesFileNames + absFile + ","
		}
	}
	return allValuesFileNames, nil
}

func DownloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func CreateFile(srcfilename, buffer string) error {
	var file, err = os.Create(srcfilename)
	if err != nil {
		return fmt.Errorf("error while creating file: %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(buffer)
	if err != nil {
		return fmt.Errorf("error while writing data to file ")
	}
	return nil
}
func SCPFileToRemote(sshIpaddress, sshPort, sshPrivateKeyPath, sshUser, srcFilePath string) error {

	if !FileExists(sshPrivateKeyPath) {
		return fmt.Errorf("file %s not exist", sshPrivateKeyPath)
	}

	if !FileExists(srcFilePath) {
		return fmt.Errorf("file %s not exist", srcFilePath)
	}

	key, err := ioutil.ReadFile(sshPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %v", err)
	}

	IpAddressAndPort := sshIpaddress + ":" + sshPort

	client, err := ssh.Dial("tcp", IpAddressAndPort, &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}

	src, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("failed to open file :%v", err)
	}
	src.Close()

	errs := scp.CopyPath(src.Name(), src.Name(), session)
	if errs != nil {
		return fmt.Errorf("failed to Copyfile to Dest: %v", errs)
	}
	return nil
}

func Contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func FilterRafayLabels(l map[string]string) string {
	for key := range l {
		if strings.HasPrefix(key, "rafay.dev/") {
			delete(l, key)
		}
	}
	mJson, _ := json.Marshal(l)
	if string(mJson) == "null" {
		return ""
	} else {
		return string(mJson)
	}
}

func CustomLabelPrinter(labels map[string]string) string {
	for key := range labels {
		if strings.HasPrefix(key, "rafay.dev/") {
			delete(labels, key)
		}
	}
	items := &LabelStruct{Items: labels}
	yamlBytes, err := yaml.Marshal(items)
	if err != nil {
		return ""
	}
	return string(yamlBytes)

}

func ValidateLabel(label string) (map[string]string, error) {

	// label : [key1:value1 key2:value2]

	arrayMaps := strings.Split(label, ",")
	mapper := make(map[string]string, len(arrayMaps))

	for _, dict := range arrayMaps {
		if strings.Contains(dict, ":") {
			splitter := strings.Split(dict, ":")
			if strings.HasPrefix(splitter[0], "kubernetes.io/") || strings.HasPrefix(splitter[0], "k8s.io/") || strings.HasPrefix(splitter[0], "rafay.dev/") {
				return nil, fmt.Errorf("cannot use %s, reserved for kubernetes and Rafay labels", splitter[0])
			} else {
				mapper[splitter[0]] = splitter[1]
			}
		} else {
			return nil, fmt.Errorf("invalid label(s) format, please look the labels schema in exmaples")
		}
	}

	return mapper, nil

}
