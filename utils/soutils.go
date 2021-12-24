package utils

import "runtime"

func GetSoName() string {
	var soName string
	if runtime.GOOS == "windows" {
		soName = "win"
	} else if runtime.GOOS == "darwin" {
		soName = "darwin"
	} else if runtime.GOOS == "linux" {
		soName = "linux"
	}
	return soName
}

func GetArch() string {
	var arch string
	if runtime.GOARCH == "amd64" {
		arch = "x64"
	} else if runtime.GOARCH == "386" {
		arch = "x86"
	}
	return arch
}
