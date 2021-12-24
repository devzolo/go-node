package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/devzolo/go-node/utils"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetNodeLtsDownloadUri() string {
	version := utils.FindNodeJsLtsVersion()
	soName := utils.GetSoName()
	arch := utils.GetArch()
	return fmt.Sprintf("https://nodejs.org/dist/%s/node-%s-%s-%s.zip", version, version, soName, arch)
}

func downloadLatestNodeLts() {
	os.RemoveAll("bin")
	fmt.Println("Downloading latest node lts")

	uri := GetNodeLtsDownloadUri()

	fmt.Println(uri)

	resp, err := http.Get(uri)
	check(err)

	if resp.StatusCode != 200 {
		panic(resp.Status)
	}

	fmt.Println("Downloaded latest node lts")

	os.Mkdir("temp", os.ModePerm)

	f, err := os.Create("temp/node.zip")
	check(err)

	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			break
		}
		f.Write(buf[:n])
	}

	f.Close()

	extractNodeLts("temp/node.zip")

	os.RemoveAll("temp")
}

func extractNodeLts(file string) {
	fmt.Println("Extracting node lts")
	version := utils.FindNodeJsLtsVersion()
	soName := utils.GetSoName()
	arch := utils.GetArch()

	firstPathToRemove := "node-" + version + "-" + soName + "-" + arch

	os.MkdirAll("bin/nodejs", 0777)

	r, err := zip.OpenReader(file)
	check(err)

	dirToExtract := "bin/nodejs/"

	for _, f := range r.File {
		rc, err := f.Open()
		check(err)
		defer rc.Close()

		entryName := dirToExtract + strings.Replace(f.Name, firstPathToRemove+"/", "", 1)
		if f.FileInfo().IsDir() {
			os.MkdirAll(entryName, 0777)
		} else {
			w, err := os.Create(entryName)
			check(err)
			defer w.Close()

			_, err = io.Copy(w, rc)
			check(err)
		}
	}
	fmt.Println("Extracted node lts")
	r.Close()
}

func main() {
	// FindNodeJsLtsVersion()
	// uri := GetNodeLtsDownloadUri()
	// fmt.Println(uri)

	downloadLatestNodeLts()
}
