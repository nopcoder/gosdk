package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

const goDownloadInformationURL = "https://go.dev/dl/?mode=json"

type GoDownloadInformation []struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []struct {
		Filename string `json:"filename"`
		Os       string `json:"os"`
		Arch     string `json:"arch"`
		Version  string `json:"version"`
		Sha256   string `json:"sha256"`
		Size     int    `json:"size"`
		Kind     string `json:"kind"`
	} `json:"files"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gosdk [remote|install]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "remote":
		latestVersion, err := getLatestGoVersion()
		if err != nil {
			fmt.Printf("Error fetching latest version information: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Remote Go version: %s\n", latestVersion)
	case "install":
		var version string
		if len(os.Args) < 3 {
			if latest, err := getLatestGoVersion(); err != nil {
				fmt.Printf("Error fetching latest version information: %s\n", err)
				os.Exit(1)
			} else {
				version = latest
			}
		} else {
			version = os.Args[2]
		}
		err := installGoVersion(version)
		if err != nil {
			fmt.Printf("Error installing Go version %s: %v\n", version, err)
			os.Exit(1)
		}
	default:
		fmt.Println("Invalid command. Use 'remote' or 'install'.")
		os.Exit(1)
	}
}

func getLatestGoVersion() (string, error) {
	resp, err := http.Get(goDownloadInformationURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var downloadInformation GoDownloadInformation
	if err := json.NewDecoder(resp.Body).Decode(&downloadInformation); err != nil {
		return "", err
	}
	if len(downloadInformation) == 0 {
		return "", fmt.Errorf("no versions found")
	}
	return downloadInformation[0].Version, nil
}

func installGoVersion(version string) error {
	fmt.Printf("Download %s installer...\n", version)
	installCmd := exec.Command("go", "install", fmt.Sprintf("golang.org/dl/%s@latest", version))
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Download %s SDK...\n", version)
	execCmd := exec.Command(version, "download")
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	return execCmd.Run()
}
