package restore

import (
	"encoding/json"
	"fmt"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
	"strings"
)

func Dashboards(username, password, url, directory string) {

	var (
		filesInDir []os.DirEntry
		rawDB      []byte
	)
	foldername := "dashboards"
	userinfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userinfo}
	client, err := gapi.New(url, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}

	path := filepath.Join(directory, foldername)

	filesInDir, err = os.ReadDir(path)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawDB, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}

			var newDB gapi.Dashboard
			if err = json.Unmarshal(rawDB, &newDB); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
			newDB.Model["id"] = ""
			newDB.FolderID = 0
			uid := fmt.Sprint(newDB.Model["uid"])
			exists, _ := client.DashboardByUID(uid)
			if exists != nil {
				client.DeleteDashboardByUID(uid)

				client.NewDashboard(newDB)
			} else {
				client.NewDashboard(newDB)
			}

		}
	}
}
