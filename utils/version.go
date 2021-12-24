package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CompareVersions(a, b string) int {

	aClean := strings.Replace(a, "v", "", -1)
	aClean = strings.Replace(aClean, ".", "", -1)

	bClean := strings.Replace(b, "v", "", -1)
	bClean = strings.Replace(bClean, ".", "", -1)

	aInt, err := strconv.Atoi(aClean)
	check(err)

	bInt, err := strconv.Atoi(bClean)
	check(err)

	if aInt < bInt {
		return -1
	} else if aInt > bInt {
		return 1
	}
	return 0
}

type FlexBool struct {
	Val bool
	Txt string
}

func (f *FlexBool) UnmarshalJSON(b []byte) error {
	if f.Txt = strings.Trim(string(b), `"`); f.Txt == "" {
		f.Txt = "false"
	}
	f.Val = f.Txt == "1" ||
		strings.EqualFold(f.Txt, "true") ||
		strings.EqualFold(f.Txt, "yes") ||
		strings.EqualFold(f.Txt, "t") ||
		strings.EqualFold(f.Txt, "armed") ||
		strings.EqualFold(f.Txt, "active") ||
		strings.EqualFold(f.Txt, "enabled") ||
		strings.EqualFold(f.Txt, "ready") ||
		strings.EqualFold(f.Txt, "up") ||
		strings.EqualFold(f.Txt, "ok")
	return nil
}

type VersionInfo struct {
	Version string   `json:"version"`
	Lts     FlexBool `json:"lts"`
}

func FindNodeJsLtsVersion() string {
	baseUrl := "https://nodejs.org/dist/index.json"
	resp, err := http.Get(baseUrl)
	check(err)

	if resp.StatusCode != 200 {
		panic(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	var versions []VersionInfo
	err = json.Unmarshal(body, &versions)
	check(err)

	var versionsSlice []string

	for _, archive := range versions {
		if archive.Lts.Txt != "false" {
			versionsSlice = append(versionsSlice, archive.Version)
		}
	}

	var latestVersion string

	for _, version := range versionsSlice {
		if latestVersion == "" {
			latestVersion = version
		} else {
			if CompareVersions(version, latestVersion) > 0 {
				latestVersion = version
			}
		}
	}

	fmt.Println("Latest version is " + latestVersion)

	return latestVersion
}
