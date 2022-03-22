//go:build generate
// +build generate

package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetNPMPackageJson() ([]byte, error) {
	// Query the NPM registry API for the package.json file.
	url := "https://registry.npmjs.org/tinymce/latest"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func GetNPMPackageVersionFromJSON() (string, error) {
	// Get the JSON
	byt, err := GetNPMPackageJson()
	if err != nil {
		return "", err
	}
	// Marshal the JSON into an object
	var data map[string]interface{}
	err = json.Unmarshal(byt, &data)
	if err != nil {
		return "", err
	}
	// Get the version
	for i, v := range data {
		if i == "version" {
			fmt.Println(i, ":", v)
			return v.(string), nil
		}
	}
	version := "6.0.0"
	return version, nil
}

func GetTinyMCEReleaseFromVersion(version string) (string, error) {
	url := fmt.Sprintf("https://download.tiny.cloud/tinymce/community/tinymce_%s.zip", version)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile("tinymce.zip", ret, 0644)
	if err != nil {
		return "", err
	}
	return "tinymce.zip", nil
}

func UnzipTinyMCEDotZip(zipFile string) error {
	// Unzip the file
	dst := "www"
	os.RemoveAll(dst)
	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path")
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}

func CopyIndexHTML() error {
	src := "index.html"
	dst := "www/index.html"
	os.Remove(dst)
	byt, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, byt, 0644)
	if err != nil {
		return err
	}
	return nil
}

func CopyStyleCSS() error {
	src := "style.css"
	dst := "www/style.css"
	os.Remove(dst)
	byt, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, byt, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	version, err := GetNPMPackageVersionFromJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(version)
	release, err := GetTinyMCEReleaseFromVersion(version)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(release)
	err = UnzipTinyMCEDotZip(release)
	if err != nil {
		log.Fatal(err)
	}
	err = CopyIndexHTML()
	if err != nil {
		log.Fatal(err)
	}
	err = CopyStyleCSS()
	if err != nil {
		log.Fatal(err)
	}
}
