package fileutil

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/erock530/go.logging"
)

const zipSlip = "zipSlip"

//ZipData returns a zip io.ReadSeeker
func ZipData(files map[string]io.Reader) io.ReadSeeker {

	buf := new(bytes.Buffer)
	zipper := zip.NewWriter(buf)

	for k, v := range files {
		z, err := zipper.Create(k)
		if err != nil {
			logging.Errorf("Error creating zip file: %+v", err)
			return nil
		}
		_, err = io.Copy(z, v)
		if err != nil {
			logging.Errorf("Error writing zip file: %+v", err)
			return nil
		}
	}
	zipper.Close()
	return bytes.NewReader(buf.Bytes())
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file to an output directory.
//Any zipslip files will be put into directory called zipslip  {dest}/zipsip/_etc_systemd_system_possible_bad_file.conf
//Modified from https://golangcode.com/unzip-files-in-go/ to work with soft-links and perserve zipSlip files
func Unzip(src string, dest string) ([]string, error) {
	filenames := []string{}

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: https://snyk.io/research/zip-slip-vulnerability#go
		if !strings.HasPrefix(fpath, filepath.Clean(dest)) {
			base := filepath.Base(fpath)
			if base == "." || base == ".." || base == string(os.PathSeparator) {
				continue
			}
			//perserve problem path within filename
			name := strings.Replace(fpath, string(os.PathSeparator), "_", -1)
			fpath = filepath.Join(dest, zipSlip, name)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		if f.Mode()&os.ModeSymlink != 0 {
			os.Remove(fpath)
			buff, _ := ioutil.ReadAll(rc)
			oldname := string(buff)
			err = os.Symlink(oldname, fpath)
			rc.Close()
			if err != nil {
				return filenames, fmt.Errorf("Create Symlink failed, needs Admin rights on windows %v", err)
			}
			continue
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
