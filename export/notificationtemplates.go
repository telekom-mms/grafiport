package export

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafana-exporter/common"
	"os"
	"path/filepath"
)

func NotificationTemplates(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "notificationTemplates"
	path := common.InitializeFolder(directory, folderName)          // initialize Subfolder to export to it
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	templates, err := client.MessageTemplates()
	if err != nil {
		log.Error("Failed to get NotificationTemplates ", err)
		return err
	}

	for _, template := range templates {
		t, err := client.MessageTemplate(template.Name)
		if err != nil {
			log.Error("Error fetching NotificationTemplate ", err)
		}
		jsonNotificationTemplates, err := json.Marshal(t)
		if err != nil {
			log.Error("Error unmarshalling json File ", err)
		}
		err = os.WriteFile(filepath.Join(path, slug.Make(template.Name)+".json"), jsonNotificationTemplates, os.FileMode(0666))
		if err != nil {
			log.Error("Couldn't write NotificationTemplates to disk ", err)
		} else {
			log.Info("Exported NotificationTemplate " + template.Name)
		}
	}
	return nil
}
