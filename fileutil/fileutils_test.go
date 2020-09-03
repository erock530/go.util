package fileutil

import (
	"os"
	"path/filepath"
	"testing"
)

//TestProxyConfig is the type stored in config files
type TestProxyConfig struct {
	DisplayName string   `yaml:"displayName,omitempty"`
	Targets     []string `yaml:"target"`
}

//Test_ReadYamlFile test reading a .yaml file and converting it into an object
func Test_ReadYamlFile(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	testFile := filepath.Join(goPath, "src", "github.com", "erock530", "go.util", "fileutil", "test.yaml")
	testObject := &TestProxyConfig{}
	err := ReadYamlFile(testFile, testObject)
	if err != nil {
		t.Error("problem reading test.yaml, ", err)
	}
	if testObject.DisplayName != "Dashboard" {
		t.Error("Unable to parse Dashborad")
	}
	if len(testObject.Targets) != 2 {
		t.Error("Unable to parse Targets")
	}
}
