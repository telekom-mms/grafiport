package restore

import (
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	"grafiport/common"
	"path/filepath"
)

// AlertRules is a function that restores all alert rules to a Grafana instance from a given folder
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func AlertRules(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "alertRules"
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}

	log.Info("Starting to restore alert rules")
	// path is the based on a provided directory and the naming of sub-folder os our tool
	path := filepath.Join(directory, folderName)
	// get all objects from provided path
	AlertRuleSlice, err := common.ReadObjectsFromDisk[gapi.AlertRule](path)
	if err != nil {
		log.Error("Error reading AlertRules from Disk")
		return err
	}
	// looping over found files to restore
	for _, newAlertRule := range AlertRuleSlice {
		// interact with api
		// search for alertRule if exists to determine if update or create
		// if no error then object exists
		_, err = client.AlertRule(newAlertRule.UID)
		if err == nil {
			err = client.UpdateAlertRule(&newAlertRule)
			if err != nil {
				log.Error("Error updating AlertRule ", err)
				break
			}
			log.Info("Updated AlertRule " + newAlertRule.Title)

		} else {
			_, err = client.NewAlertRule(&newAlertRule)
			if err != nil {
				log.Error("Error creating AlertRule ", err)
				break
			}
			log.Info("Created AlertRule " + newAlertRule.Title)
		}
	}
	return err
}
