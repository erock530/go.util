package util

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/erock530/go.util/fileutil"
	"github.com/erock530/go.util/osutil"
)

//ReadSeekCloser interface for holding file io handles
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

//GetFile Retrieves a file pointer and the files mtime
func GetFile(fileName string) (ReadSeekCloser, time.Time, error) {
	var modTime time.Time
	fi, err := os.Open(fileName)
	if err == nil {
		fInfo, _ := fi.Stat()
		if fInfo != nil {
			modTime = fInfo.ModTime()
		}
	}
	return fi, modTime, err
}

// GetModTime makes sure the file exists
func GetModTime(filePath string) (time.Time, error) {
	return fileutil.GetModTime(filePath)
}

//CreateZipArchive takes in a map of file names and their file paths, and a map of file names and the metadata
//it then zips these files and returns a reader
func CreateZipArchive(archiveName string, files map[string]string, metas map[string]string) (io.ReadSeeker, time.Time, error) {
	return fileutil.CreateZipArchive(archiveName, files, metas)
}

//ReadFile will read a file from wherever it is stored and return the bytes
func ReadFile(filePath string) ([]byte, error) {
	return fileutil.ReadFile(filePath)
}

//TrimExt removes the suffix from a file path
func TrimExt(path string) string {
	return fileutil.TrimExt(path)
}

//IsSameTimestamp tells whether the files are the same
//compares the timestamps of the two files.
func IsSameTimestamp(sourcePath, destPath string) (bool, error) {
	return fileutil.IsSameTimestamp(sourcePath, destPath)
}

//RunCmdSafely runs the command and logs errors
func RunCmdSafely(cmdStr string, args []string, workdir string, env []string) error {
	return osutil.RunCmdSafely(cmdStr, args, workdir, env)
}

//RunCmdSafelyStdOut runs the command and logs errors and captures std out into buffer
func RunCmdSafelyStdOut(cmdStr string, out *bytes.Buffer, args []string, workdir string, env []string) error {
	return osutil.RunCmdSafelyStdOut(cmdStr, out, args, workdir, env)
}

// TempFile creates a new temporary file in the directory dir
// with a name beginning with prefix, and ending with suffix (ea file extension .html)
// opens the file for reading
// and writing, and returns the resulting *os.File.
// If dir is the empty string, TempFile uses the default directory
// for temporary files (see os.TempDir).
// Multiple programs calling TempFile simultaneously
// will not choose the same file.  The caller can use f.Name()
// to find the pathname of the file.  It is the caller's responsibility
// to remove the file when no longer needed.
func TempFile(dir, prefix string, suffix string) (f *os.File, err error) {
	return fileutil.TempFile(dir, prefix, suffix)
}

//ReadYamlFile unmarshal's a yaml file into given object
func ReadYamlFile(fileName string, i interface{}) error {
	return fileutil.ReadYamlFile(fileName, i)
}

//ReadJSONFile unmarshal's a yaml file into given object
func ReadJSONFile(fileName string, i interface{}) error {
	return fileutil.ReadJSONFile(fileName, i)
}

//UnmarshalFile A generic unmarshaler from a file
func UnmarshalFile(fileName string, i interface{}, f func([]byte, interface{}) error) error {
	return fileutil.UnmarshalFile(fileName, i, f)
}
