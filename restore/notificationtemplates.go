package restore

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	"grafana-exporter/common"
	"os"
	"path/filepath"
	"strings"
)

func NotificationTemplates(username, password, url, directory string) error {
	var (
		filesInDir  []os.DirEntry
		rawTemplate []byte
	)
	folderName := "notificationTemplates"
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	log.Info("Starting to restore NotificationTemplates")
	path := filepath.Join(directory, folderName)

	filesInDir, err = os.ReadDir(path)
	if err != nil {
		log.Error("Failed to read folder%s\n", err)
		return err
	}

	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawTemplate, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				log.Error(err)
				continue
			}

			var newTemplate gapi.AlertingMessageTemplate
			if err = json.Unmarshal(rawTemplate, &newTemplate); err != nil {
				log.Error(err)
				continue
			}
			status, _ := client.MessageTemplate(newTemplate.Name)
			if status.Name != "" {
				err = client.DeleteMessageTemplate(newTemplate.Name)
				if err != nil {
					log.Error("Error updating NotificationTemplate - delete (1/2) ", err)
					continue
				}
				err = client.SetMessageTemplate(newTemplate.Name, newTemplate.Template)
				if err != nil {
					log.Error("Error updating NotificationTemplate - set (2/2) ", err)
				} else {
					log.Info("Updated NotificationTemplate " + newTemplate.Name)
				}

			} else {
				err = client.SetMessageTemplate(newTemplate.Name, newTemplate.Template)
				if err != nil {
					log.Error("Error creating NotificationTemplate ", err)
				} else {
					log.Info("Created NotificationTemplate " + newTemplate.Name)
				}
			}

		}
	}
	return nil
}
