package antivirus

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/erock530/go.util/osutil/v2"
)

//Uvscanner is the Mcafee antivirus CLI
type Uvscanner struct {
	Metadata *AVMetadata
}

//Initialize gathers anti-virus scanner information
func (uv *Uvscanner) Initialize() error {
	const initialScanTimeout = 10
	//Get antivirus scanner output metadata
	scannerArgs := []string{"-d", scannerDatLoc, scannerLoc}
	_, out, stderr := osutil.RunCmd(scannerLoc, scannerArgs, initialScanTimeout, false)
	if stderr.String() != "" {
		return fmt.Errorf("Unable to initialize uvscanner, encountered error when scanner at %s tried to scan itself: %s", scannerLoc, stderr.String())
	}
	blurb := out.String()

	//Parse blurb
	uv.Metadata = getMetadataFromAVOutput(blurb)
	return nil
}

//Scan can perform the antivirus scanning
func (uv *Uvscanner) Scan(path string) *ScanResult {
	var err error
	sr := &ScanResult{
		Path:     path,
		Metadata: uv.Metadata,
	}
	if !useVirusScanner {
		return sr
	}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		sr.ScanError = err
		return sr
	}

	env := os.Environ()
	basedir := filepath.Dir(path)
	env = append(env, fmt.Sprintf("Home=%s", basedir))

	scannerArgs := []string{avTimeout, scannerLoc, "-d", scannerDatLoc, "--maxfilesize", strconv.Itoa(maxFileSize), "--nocomp", ignoreSelfExtractingExeEnv}
	scannerArgs = append(scannerArgs, argList...)
	scannerArgs = append(scannerArgs, path)
	cmd := exec.Command(scannerLoc, scannerArgs...)
	sr.ScanError = cmd.Run()

	if sr.ScanError != nil {
		if exitError, ok := sr.ScanError.(*exec.ExitError); ok {
			if exitError.ExitCode() == 13 {
				sr.Positive = true
			}
		}
	}
	sr.ScanTime = time.Now().UTC()
	return sr
}

func getMetadataFromAVOutput(blurb string) *AVMetadata {
	meta := &AVMetadata{}
	meta.AVSoftware = "McAfee VirusScan Command Line for Linux64"
	lines := strings.Split(blurb, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		words := strings.Split(l, " ")

		if strings.HasPrefix(l, "McAfee Virus") {
			meta.AVVersion = words[len(words)-1]
		}
		if strings.HasPrefix(l, "AV Engine version: ") {
			if len(words) >= 4 {
				meta.AVEngineVersion = words[3]
			}
		}
		if strings.HasPrefix(l, "Dat set version: ") {
			if len(words) >= 4 {
				readVersion, err := strconv.Atoi(words[3])
				if err == nil {
					meta.AVDatVersion = readVersion
				}
			}
			if len(words) >= 6 && words[4] == "created" {
				meta.AVDatVersionDate = strings.Join(words[5:], " ")
			}
		}
	}
	return meta
}
