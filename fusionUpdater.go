package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var workingPath string = "C:\\FusionUpdater\\"
var fusionInstallerURL string = "https://dl.appstreaming.autodesk.com/production/installers/Fusion%20360%20Admin%20Install.exe"
var releaseVersionEndpoint string = "https://dl.appstreaming.autodesk.com/production/67316f5e79bc48318aa5f7b6bb58243d/73e72ada57b7480280f7a6f4a289729f/full.json"

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	latest := getReleaseVersionNumber(&wg)
	wg.Wait()
	current := getInstalledVersionNumber()
	if current == "" { // not installed
		fmt.Println("Fusion360 not installed.")
		wg.Add(1)
		go installFusion(&wg)
		wg.Wait()
		writeLatestVersion(latest)
	} else {
		fmt.Println("You are on version " + current + " and the latest version is " + latest)

		// Compare and update if necessary
		if compareVersions(latest, current) {
			wg.Add(1)
			go update(latest, &wg)
		} else {
			fmt.Println("Fusion is already up to date... exiting")
		}

		wg.Wait()
	}

	fmt.Println("Operation completed succcessfully. Enjoy Fusion!")
}

func installFusion(wg *sync.WaitGroup) {
	// create directory if it doesn't exists
	if _, err := os.Stat(workingPath); os.IsNotExist(err) {
		err := os.Mkdir(workingPath, 0777)

		if err != nil {
			log.Panic(err)
		}
	}

	var downloadWG sync.WaitGroup
	downloadWG.Add(1)
	go downloadNewestVersion(&downloadWG)
	downloadWG.Wait()
	fmt.Println("Download Complete. Installing")

	defer wg.Done()

	upgrade := exec.Command("C:\\FusionUpdater\\Fusion360AdminInstall.exe") //, "quiet")
	err := upgrade.Run()
	if err != nil {
		log.Fatal(err)
	}
}

/// returns true if versions differ
func compareVersions(l string, c string) bool {
	l_vals := strings.Split(l, ".")
	c_vals := strings.Split(c, ".")

	for curr, lv := range l_vals {
		if lv != c_vals[curr] { // versions differ
			return true
		}
	}

	return false
}

func update(latest string, wg *sync.WaitGroup) {

	var downloadWG sync.WaitGroup
	downloadWG.Add(1)
	go downloadNewestVersion(&downloadWG)
	downloadWG.Wait()
	fmt.Println("Download Complete. Applying Update")

	defer wg.Done()

	upgrade := exec.Command("C:\\FusionUpdater\\Fusion360AdminInstall.exe", "--process", "update", "quiet")
	err := upgrade.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Update Complete. Cleaning up")
	upgrade = exec.Command("C:\\FusionUpdater\\Fusion360AdminInstall.exe", "--process uninstall", "--purge-incomplete.exe")
	err = upgrade.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Working diretory clean.")

	writeLatestVersion(latest)
	fmt.Println("Version information updated")
}

func downloadNewestVersion(wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	fmt.Println("Downloading latest version.")

	// Create the file
	out, err := os.Create(workingPath + "Fusion360AdminInstall.exe")
	defer func() {
		if err != nil {
			log.Fatal("Error occured while downloading the latest version.")
		}
	}()
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(fusionInstallerURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func writeLatestVersion(latest string) {
	err := ioutil.WriteFile("C:\\FusionUpdater\\currentVersion.ini", []byte(latest), 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func getInstalledVersionNumber() string {

	content, err := ioutil.ReadFile("C:\\FusionUpdater\\currentVersion.ini")
	if err != nil { // file not found, make it
		return ""
	}

	return string(content)
}

func getReleaseVersionNumber(wg *sync.WaitGroup) string {
	defer wg.Done()

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, releaseVersionEndpoint, nil)

	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	res_ob := Response{}
	jsonErr := json.Unmarshal(body, &res_ob)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return res_ob.BuildVersion
}

type Response struct {
	AbortInstallUrls          []string      `json:"abort-install-urls"`
	BuildVersion              string        `json:"build-version"`
	InitialInstallDeletePaths []string      `json:"initial-install-delete-paths"`
	LauncherPath              string        `json:"launcher-path"`
	MajorUpdateVersion        string        `json:"major-update-version"`
	MajorUpdateVersionList    []interface{} `json:"major-update-version-list"`
	Packages                  []struct {
		Checksum       string `json:"checksum"`
		CompressedSize int    `json:"compressed-size"`
		Size           int    `json:"size"`
	} `json:"packages"`
	Patches             []string `json:"patches"`
	PatchesBuildVersion []string `json:"patches_build_version"`
	PreInstallTasks     []string `json:"pre-install-tasks"`
	Properties          struct {
		AutoLaunch struct {
			ID string `json:"id"`
		} `json:"auto-launch"`
		DisplayName string `json:"display-name"`
		ExecName    string `json:"exec-name"`
		RequiredOs  struct {
			FriendlyVersion string `json:"friendly-version"`
			Version         []int  `json:"version"`
		} `json:"required-os"`
		SessionFileName string   `json:"session-file-name"`
		SubApplications []string `json:"sub-applications"`
		UninstallIcon   string   `json:"uninstall-icon"`
	} `json:"properties"`
	ReleaseVersion string `json:"release-version"`
	Streamer       struct {
		Checksum       string `json:"checksum"`
		FeatureVersion string `json:"feature-version"`
	} `json:"streamer"`
	UpdateGroup string `json:"update-group"`
}
