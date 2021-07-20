//Unit tests for mcafee scanner uvscan
//
package antivirus

import (
	"path/filepath"
	"testing"
	"time"
)

var negativeTestFile = filepath.Join("testFiles", "Constitution.txt")
var positiveTestFile = filepath.Join("testFiles", "EICAR.COM")

const dirLoc = "testFiles"

func TestTrueNegativeScan(t *testing.T) {
	useVirusScanner = true
	uvs := Uvscanner{}
	err := uvs.Initialize()
	if err != nil {
		t.Errorf("Encountered intialization error: %s", err.Error())
	}

	scanResults := uvs.Scan(negativeTestFile)
	if scanResults.ScanError != nil {
		t.Errorf("Expected no error when scanning %s and instead found error: %s", negativeTestFile, scanResults.ScanError.Error())
	}
	if scanResults.IsPositive() {
		t.Errorf("Expected negative scan to return false on IsPositive check")
	}

	//Test scan metadata
	if scanResults.ScanTime.Add(1 * time.Minute).Before(time.Now().UTC()) {
		t.Errorf("Expected san time to be close to now, instead its more than a minute earlier")
	}
	if time.Now().UTC().Add(1 * time.Minute).Before(scanResults.ScanTime) {
		t.Errorf("Expected san time to be close to now, instead its more than a minute later")
	}
	if scanResults.Metadata.AVSoftware != "McAfee VirusScan Command Line for Linux64" {
		t.Errorf("Unexpected virus scanning software %s", scanResults.Metadata.AVSoftware)
	}
	if scanResults.Metadata.AVVersion != "6.1.4.305" {
		t.Errorf("Unexpected virus scanner version %s", scanResults.Metadata.AVVersion)
	}
	if scanResults.Metadata.AVEngineVersion != "6100.8979" {
		t.Errorf("Unexpected virus scanner engine version %s", scanResults.Metadata.AVEngineVersion)
	}
	if scanResults.Metadata.AVDatVersion != 9844 {
		t.Errorf("Unexpected dat file version %d", scanResults.Metadata.AVDatVersion)
	}
	if scanResults.Metadata.AVDatVersionDate != "Dec 23 2020" {
		t.Errorf("Unexpected dat file creation date %s", scanResults.Metadata.AVDatVersionDate)
	}

}

func TestTruePositiveScan(t *testing.T) {
	useVirusScanner = true
	uvs := Uvscanner{}
	scanResults := uvs.Scan(positiveTestFile)
	if scanResults.ScanError == nil {
		t.Fatalf("Expected a postive scan error when scanning %s and instead found no error", positiveTestFile)
	}

	if !scanResults.IsPositive() {
		t.Errorf("Expected positive scan to return true on IsPositive check")
	}

}

func TestTurningAntivirusOff(t *testing.T) {
	useVirusScanner = false
	uvs := Uvscanner{}
	scanResults := uvs.Scan(positiveTestFile)
	if scanResults.ScanError != nil {
		t.Errorf("Expected no scan error when scanning %s with the scanner turned off and instead found an error: %s", positiveTestFile, scanResults.ScanError.Error())
	}
}

func TestAdditionalParemterSupport(t *testing.T) {
	useVirusScanner = true
	uvs := Uvscanner{}
	//Test arg list option works by seeing if we can scan a folder
	argList = []string{"-r", "--secure"}

	scanResults := uvs.Scan(".")
	if scanResults.ScanError == nil {
		t.Fatalf("Expected a postive scan error when scanning %s and instead found no error", positiveTestFile)
	}
	if !scanResults.IsPositive() {
		t.Errorf("Expected positive scan to return true on IsPositive check")
	}

}
