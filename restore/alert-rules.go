package restore

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
	"strings"
)

func AlertRules(username, password, url, directory string) error {
	var (
		filesInDir   []os.DirEntry
		rawAlertRule []byte
		err          error
	)
	folderName := "alertRules"
	userInfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userInfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}

	log.Info("Starting to restore alert rules")
	path := filepath.Join(directory, folderName)

	filesInDir, err = os.ReadDir(path)
	if err != nil {
		log.Error("Failed to read alert rules%s\n", err)
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawAlertRule, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				log.Error(err)
				continue
			}

			var newAlertRule gapi.AlertRule
			if err = json.Unmarshal(rawAlertRule, &newAlertRule); err != nil {
				log.Error(err)
				continue
			}

			_, err = client.AlertRule(newAlertRule.UID)
			if err == nil {
				err = client.UpdateAlertRule(&newAlertRule)
				if err != nil {
					log.Error("Error updating AlertRule ", err)
				} else {
					log.Info("Updated AlertRule", newAlertRule.Title)
				}

			} else {
				_, err = client.NewAlertRule(&newAlertRule)
				if err != nil {
					log.Error("Error creating AlertRule ", err)
				} else {
					log.Info("Created AlertRule", newAlertRule.Title)
				}
			}
		}
	}
	return err
}
