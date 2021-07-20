package antivirus

import (
	"fmt"
	"time"
)

//AVScanner is an antivirus scanner
type AVScanner interface {
	//Initialize gathers anti-virus scanner information
	Initialize() error
	//Scan takes in a path and returns an error object
	Scan(string) *ScanResult
}

//ScanResult is the results of a scan
type ScanResult struct {
	Path      string
	ScanError error
	Positive  bool
	ScanTime  time.Time
	Metadata  *AVMetadata
}

//AVMetadata is the antivirus metadata
type AVMetadata struct {
	AVSoftware       string
	AVVersion        string
	AVEngineVersion  string
	AVDatVersion     int
	AVDatVersionDate string
}

func (sr *ScanResult) Error() string {
	if sr.ScanError == nil {
		return ""
	}
	if sr.IsPositive() {
		return fmt.Sprintf("Antivirus software has tagged %s as containing a virus", sr.Path)
	}
	return sr.ScanError.Error()
}

//IsPositive returns true if scan detected something, false otherwise
func (sr *ScanResult) IsPositive() bool {
	return sr.Positive
}

//GetBlurb gets the scanner blurb
func (sr *ScanResult) GetBlurb() string {
	blurb := fmt.Sprintf("%s: %s", sr.Metadata.AVSoftware, sr.Metadata.AVVersion)
	if sr.Metadata.AVEngineVersion != "" {
		blurb = fmt.Sprintf("%s\nAV Engine version: %s", blurb, sr.Metadata.AVEngineVersion)
	}
	if sr.Metadata.AVDatVersion > 0 {
		blurb = fmt.Sprintf("%s\nDat set version: %d", blurb, sr.Metadata.AVDatVersion)
		if sr.Metadata.AVDatVersionDate != "" {
			blurb = fmt.Sprintf("%s created %s", blurb, sr.Metadata.AVDatVersionDate)
		}
	}
	return blurb
}
