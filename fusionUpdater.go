package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	latest := getLatestVersion()
	current := getCurrentVersion()

	// Compare and update if necessary
	if compareVersions(latest, current, &wg) {
		go update(latest, &wg)
	}

	wg.Wait()
}

func compareVersions(l string, c string, wg *sync.WaitGroup) bool {
	//fmt.Printf("old %s  :: new %s\n", c, l)
	l_vals := strings.Split(l, ".")
	c_vals := strings.Split(c, ".")

	for curr, lv := range l_vals {
		if lv != c_vals[curr] { // versions differ
			return true
		}
	}

	wg.Done()
	return false
}

func update(latest string, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		upgrade := exec.Command("C:\\FusionUpdater\\Fusion360AdminInstall.exe", "--process", "update", "quiet")
		err := upgrade.Start()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(upgrade.Args)
		upgrade.Wait()

		writeLatestVersion(latest)
	}()

}

func writeLatestVersion(latest string) {
	err := ioutil.WriteFile("C:\\FusionUpdater\\currentVersion.ini", []byte(latest), 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func getCurrentVersion() string {

	content, err := ioutil.ReadFile("C:\\FusionUpdater\\currentVersion.ini")
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func getLatestVersion() string {
	endpoint := "https://dl.appstreaming.autodesk.com/production/67316f5e79bc48318aa5f7b6bb58243d/73e72ada57b7480280f7a6f4a289729f/full.json"
	client := http.Client{
		Timeout: time.Second * 4,
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)

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
