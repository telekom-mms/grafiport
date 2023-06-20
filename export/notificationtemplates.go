package export

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafana-exporter/common"
	"os"
	"path/filepath"
)

// NotificationTemplates is a function that exports all templates for notifications from a Grafana instance and stores them as JSON files in a directory.
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func NotificationTemplates(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "notificationTemplates"
	path := common.InitializeFolder(directory, folderName)          // initialize Sub-folder to export to it
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	templates, err := client.MessageTemplates() //Get slice of all templates in a short version
	if err != nil {
		log.Error("Failed to get NotificationTemplates ", err)
		return err
	}
	// iterate over all templates
	for _, template := range templates {
		// Get slice of the current templates by Name.
		// Name is in this case the Identifier
		t, err := client.MessageTemplate(template.Name)
		if err != nil {
			log.Error("Error fetching NotificationTemplate ", err)
		}
		jsonNotificationTemplates, err := json.Marshal(t) // create JSON Object
		if err != nil {
			log.Error("Error unmarshalling json File ", err)
		}
		// write Templates as json to file
		err = os.WriteFile(filepath.Join(path, slug.Make(template.Name)+".json"), jsonNotificationTemplates, os.FileMode(0666)) // Make sure filename is unique, FileMode is irrelevant, but required for WriteFile
		if err != nil {
			log.Error("Couldn't write NotificationTemplates to disk ", err)
		}
		log.Info("Exported NotificationTemplate " + template.Name)
	}
	return nil
}
