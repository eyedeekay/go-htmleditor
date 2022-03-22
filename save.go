package tinymce

import (
	"io/ioutil"
	"os"
)

func SaveFileOnDisk(file string, body []byte) error {
	// get the original file permissions
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	// write the new file
	err = ioutil.WriteFile(file, body, fi.Mode())
	if err != nil {
		return err
	}
	return nil
}
