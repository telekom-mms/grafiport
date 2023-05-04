package restore

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
	"strings"
)

func Dashboards(username, password, url, directory string) error {

	var (
		filesInDir []os.DirEntry
		rawDB      []byte
	)
	folderName := "dashboards"
	userInfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userInfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}

	path := filepath.Join(directory, folderName)

	filesInDir, err = os.ReadDir(path)
	if err != nil {
		log.Error("Failed to read folder%s\n", err)
		return err
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawDB, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				log.Error(err)
				continue
			}

			var newDB gapi.Dashboard
			if err = json.Unmarshal(rawDB, &newDB); err != nil {
				log.Error(err)
				continue
			}
			newDB.Model["id"] = ""
			newDB.FolderID = 0
			uid := fmt.Sprint(newDB.Model["uid"])
			exists, _ := client.DashboardByUID(uid)
			if exists != nil {
				err = client.DeleteDashboardByUID(uid)
				if err != nil {
					log.Error("Error updating Dashboard - delete (1/2)", err)
				}
				_, err = client.NewDashboard(newDB)
				if err != nil {
					log.Error("Error updating Dashboard - create (2/2)", err)
				}
			} else {
				_, err = client.NewDashboard(newDB)
				if err != nil {
					log.Error("Error creating Dashboard", err)
				}
			}

		}
	}
	return nil
}
