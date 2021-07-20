// antivirus.go
// Provides methods and routines to interact with a virus scanner

package antivirus

import (
	"encoding/json"
	"fmt"
	"strconv"

	logging "github.com/erock530/go.logging/v2"
	"github.com/erock530/go.util/osutil/v2"
)

const (
	useVirusScannerEnv             = "VIRUS_SCAN_ENABLED"
	scannerLocEnv                  = "VIRUS_SCANNER_PATH"
	scannerDatLocEnv               = "VIRUS_SCANNER_DAT_PATH"
	maxFileSizeEnv                 = "MAXFILESIZE"
	memSizeEnv                     = "MEMSIZE"
	ignoreSelfExtractingExeEnv     = "NOCOMP"
	avTimeoutEnv                   = "VIRUS_SCANNER_TIMEOUT"
	argListEnv                     = "VIRUS_SCANNER_ARGS"
	defaultUseVirusScanner         = "false"
	defaultScannerLoc              = "/usr/local/bin/uvscan"
	defaultDatLoc                  = "/usr/local/uvscan/"
	defaultMaxFileSize             = "10"
	defaultMemSize                 = "1000"
	defaultIgnoreSelfExtractingExe = "false"
	defaultAvTimeout               = "10m"
	defaultArgList                 = "[]"
)

var (
	useVirusScanner         bool
	scannerLoc              string
	scannerDatLoc           string
	maxFileSize             int //MB
	memSize                 int //KB
	ignoreSelfExtractingExe bool
	argList                 []string
	avTimeout               string //timeout format, e.g. 1m or 20s
)

func errorMsg(env string, err error) error {
	return fmt.Errorf("During antivirus init, Unable to parse %s: %s", env, err.Error())
}

func init() {
	var err error
	useVirusScanner, err = strconv.ParseBool(osutil.GetSetEnv(useVirusScannerEnv, defaultUseVirusScanner))
	if err != nil {
		logging.Fatal(errorMsg(useVirusScannerEnv, err))
	}

	if !useVirusScanner {
		logging.Info("Virus scanning is currently disabled for this service")
	} else {
		logging.Info("Virus scanning is enabled for this service")
	}
	scannerLoc = osutil.GetSetEnv(scannerLocEnv, defaultScannerLoc)
	scannerDatLoc = osutil.GetSetEnv(scannerDatLocEnv, defaultDatLoc)

	maxFileSize, err = strconv.Atoi(osutil.GetSetEnv(maxFileSizeEnv, defaultMaxFileSize))
	if err != nil {
		logging.Fatal(errorMsg(maxFileSizeEnv, err))
	}

	memSize, err = strconv.Atoi(osutil.GetSetEnv(memSizeEnv, defaultMemSize))
	if err != nil {
		logging.Fatal(errorMsg(memSizeEnv, err))
	}

	ignoreSelfExtractingExe, err = strconv.ParseBool(osutil.GetSetEnv(ignoreSelfExtractingExeEnv, defaultIgnoreSelfExtractingExe))
	if err != nil {
		logging.Fatal(errorMsg(ignoreSelfExtractingExeEnv, err))
	}

	avTimeout = osutil.GetSetEnv(avTimeoutEnv, defaultAvTimeout)

	err = json.Unmarshal([]byte(osutil.GetSetEnv(argListEnv, defaultArgList)), &argList)
	if err != nil {
		logging.Fatal(errorMsg(argListEnv, err))
	}

	if useVirusScanner {
		logging.Info("AV timeout is set to: %s", avTimeout)
		logging.Info("Scanner location is set to: %s", scannerLoc)
		logging.Info("Dat file location is set to: %s", scannerDatLoc)
		logging.Info("Mem size env is set to %d", memSize)
		logging.Info("Ignore self extracting exe - %t", ignoreSelfExtractingExe)
		logging.Info("Additional arg list is set to: %v", argList)
	}
}

//GetAVScanner returns an initialized anti-virus scanner interface
func GetAVScanner() (AVScanner, error) {
	uvs := &Uvscanner{}
	err := uvs.Initialize()
	return uvs, err
}
