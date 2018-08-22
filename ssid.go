package main

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
)

var (
	errCantGetSSID = errors.New("can't find SSID")
	osxSSIDre      = regexp.MustCompile(`\sSSID: (.+)\n`)
	lastSSID       string
)

func getSSID() (ssid string, hasChanged bool, err error) {
	switch runtime.GOOS {
	case "darwin":
		ssid, err = getSSIDDarwin()
	default:
		panic(fmt.Sprintf("unsupported OS: %v", runtime.GOOS))
	}
	return ssid, ssid != lastSSID, err
}

func getSSIDDarwin() (string, error) {
	cmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-I")
	b, err := cmd.Output()
	if err != nil {
		return "", err
	}

	submatch := osxSSIDre.FindStringSubmatch(string(b))
	// first index is whole string, second is the first capture group
	if len(submatch) != 2 {
		return "", errCantGetSSID
	}

	return submatch[1], nil
}
