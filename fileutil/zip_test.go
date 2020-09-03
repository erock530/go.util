package fileutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//Test_Unzip test unzipping file using function
func Test_Unzip(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	testFile := filepath.Join(goPath, "src", "github.com", "erock530", "go.util", "fileutil", "test.zip")
	tmp := os.TempDir()
	output, err := Unzip(testFile, filepath.Join(tmp, "unzip"))
	if err != nil {
		t.Errorf("problem reading %s %v", testFile, err)
		return
	}

	foundA := false
	foundB := false
	foundzipSlip := false
	foundLink := false

	for _, file := range output {
		if file == filepath.Join(tmp, "unzip", "a.txt") {
			foundA = true
		}
		if file == filepath.Join(tmp, "unzip", "sub_dir", "b.txt") {
			foundB = true
		}
		if file == filepath.Join(tmp, "unzip", "testLink") {
			foundLink = true
		}
		if strings.Contains(file, "zipSlip") && strings.Contains(file, "etc_systemd_system_aide-search.service") {
			foundzipSlip = true
		}
	}
	if !foundA || !foundB || !foundLink || !foundzipSlip {
		t.Errorf("Unable to find expected output foundA %v foundB %v foundLink %v foundzipSlip %v", foundA, foundB, foundLink, foundzipSlip)
	}
}
