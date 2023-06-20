package export

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafana-exporter/common"
	"os"
	"path/filepath"
)

// AlertRules is a function that exports all alert rules from a Grafana instance and stores them as JSON files in a directory.
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func AlertRules(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "alertRules"
	path := common.InitializeFolder(directory, folderName)          // initialize Sub-folder to export to it
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	alertRules, err := client.AlertRules() //Get slice of all AlertRules
	if err != nil {
		log.Error("Failed to get AlertRules", err)
		return err
	}
	// iterate over all AlertRules to be consistent with all other objects
	for _, alertRule := range alertRules {
		jsonAlertRule, err := json.Marshal(alertRule) // Create JSON Object of AlertRule from received Bytes
		if err != nil {
			log.Error("Error unmarshalling json File", err)
		}
		// write Dashboards as json to a File
		err = os.WriteFile(filepath.Join(path, slug.Make(alertRule.Title+" "+alertRule.UID)+".json"), jsonAlertRule, os.FileMode(0666)) // Make sure filename is unique, FileMode is irrelevant, but required for WriteFile
		if err != nil {
			log.Error("Couldn't write AlertRule to disk ", err)
		}
		log.Info("Exported AlertRule " + alertRule.Title)
	}
	return nil
}
