package constants

import (
	"path/filepath"

	"github.com/RafaySystems/rafay-istio-multicluster/pkg/utils"
)

var RAFAY_DIR_DEFAULT_LOCATION = filepath.Join(utils.GetUserHome(), ".rafay")
var CLI_DIR_DEFAULT_LOCATION = filepath.Join(RAFAY_DIR_DEFAULT_LOCATION, "cli")
var CONFIG_FILE_DEFAULT_NAME = "config"
var LOG_FILE_DEFAULT_NAME = "cli.log"
var CONFIG_FILE_DEFAULT_LOCATION = filepath.Join(CLI_DIR_DEFAULT_LOCATION, CONFIG_FILE_DEFAULT_NAME)
var LOG_FILE_DEFAULT_LOCATION = filepath.Join(CLI_DIR_DEFAULT_LOCATION, LOG_FILE_DEFAULT_NAME)

var ENV_FLAG_NAME = "env"
var VERBOSE_FLAG_NAME = "verbose"
var DEBUG_FLAG_NAME = "debug"
var CONFIG_FLAG_NAME = "config"
var TEST_FLAG_NAME = "qa"
var PROFILE_FLAG_NAME = "profile"
var STRUCTURED_OUTPUT_FLAG_NAME = "structured"

var CONFIG_API_VERSION = "1.0"
var WP_API_VERSION = "1.0"
var CLI_VERSION = "1.0"
var CLI_BUILD_NUMBER = "NA"
var CLI_ARCH = "NA"
var CLI_BUILD_TIME = "NA"

var GENERIC_ERROR_MESSAGE = "CLI faced an issue while running the command %s. Please use -v flag to see debug logs."
var AUTH_FAILURE_ERROR_MESSAGE = "CLI faced an authentication failure. Please check your config file or use \"init\" command to configure cli with the appropriate credentials."
var CLI_MISMATCH_ERROR_MESSAGE = "Your CLI version is incompatible with the current API version. Please visit console to download the most uptodate CLI version."

var WORKLOAD_CREATED = 201
var REQUEST_SUCCESSFUL = 200
var AUTHENTICATION_FAILURE = 401

const DRIFT_ACTION_DETECT_NOTIFY = "DetectAndNotify"
const DRIFT_ACTION_BLOCK_NOTIFY = "BlockAndNotify"

// Value In seconds
const WAIT_OPERATION_POLL_INTERVAL = 30
const WAIT_OPERATION_TIMEOUT_INTERVAL = 7201
const DEFAULT_PROJECT_ID = "1"
const DEFAULT_PROJECT_NAME = "defaultproject"
const ErrIstioctlNotFound = "Unable to find Istioctl"

const CLI_NAME = "ristioctl"
const CLI_HOME_FOLDER = ".ristioctl"
const HELLO_WORLD_NAMESPACE = "sample"
