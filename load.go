package tinymce

import (
	"io/ioutil"
)

func LoadFileOnDisk(file string) ([]byte, error) {
	fb, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return fb, nil
}
