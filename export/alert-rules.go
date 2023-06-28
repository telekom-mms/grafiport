package export

import (
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafiport/common"
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
		err = common.WriteObjectToDisk(alertRule, path, slug.Make(alertRule.Title+" "+alertRule.UID)+".json")
		if err != nil {
			log.Error("Couldn't write AlertRule to disk ", err)
		}
		log.Info("Exported AlertRule " + alertRule.Title)
	}
	return nil
}
