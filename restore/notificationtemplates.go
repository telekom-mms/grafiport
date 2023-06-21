package restore

import (
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	"grafana-exporter/common"
	"path/filepath"
)

// NotificationTemplates is a function that restores all notification templates to a Grafana instance from a given folder
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func NotificationTemplates(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "notificationTemplates"
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	log.Info("Starting to restore NotificationTemplates")
	// path is the based on a provided directory and the naming of sub-folder os our tool
	path := filepath.Join(directory, folderName)
	// get all objects from provided path
	AlertingMessageTemplatesSlice, err := common.ReadObjectsFromDisk[gapi.AlertingMessageTemplate](path)
	if err != nil {
		log.Error("Error reading AlertRules from Disk")
		return err
	}
	// looping over found files to restore
	for _, newTemplate := range AlertingMessageTemplatesSlice {
		// interact with api
		// search for alertRule if exists to determine if update or create
		// if no error then object exists
		status, _ := client.MessageTemplate(newTemplate.Name)
		if status.Name != "" {
			// library misses Update Function, so implemented by deleting and overwriting the new Config
			err = client.DeleteMessageTemplate(newTemplate.Name)
			if err != nil {
				log.Error("Error updating NotificationTemplate - delete (1/2) ", err)
				break
			}
			err = client.SetMessageTemplate(newTemplate.Name, newTemplate.Template)
			if err != nil {
				log.Error("Error updating NotificationTemplate - set (2/2) ", err)
				break
			}
			log.Info("Updated NotificationTemplate " + newTemplate.Name)

		} else {
			err = client.SetMessageTemplate(newTemplate.Name, newTemplate.Template)
			if err != nil {
				log.Error("Error creating NotificationTemplate ", err)
				break
			}
			log.Info("Created NotificationTemplate " + newTemplate.Name)

		}

	}
	return nil
}
