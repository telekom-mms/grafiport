package export

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafana-exporter/common"
	"os"
	"path/filepath"
)

// Dashboards is a function that exports all dashboards from a Grafana instance and stores them as JSON files in a directory.
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func Dashboards(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "dashboards"                                      //Name of sub-folder
	path := common.InitializeFolder(directory, folderName)          // initialize Sub-folder to export to it
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	dashboards, err := client.Dashboards() //Get slice of all dashboards in a short version
	if err != nil {
		log.Error("Failed to get Dashboards ", err)
		return err
	}
	// iterate over all Dashboards
	for _, dashboard := range dashboards {
		// Get slice of the current dashboard by UID.
		// UID is in this case the Identifier
		ds, err := client.DashboardByUID(dashboard.UID)
		if err != nil {
			log.Error("Error fetching Dashboard ", err)
		}
		jsonDashboard, err := json.Marshal(ds) // create JSON Object from Dashboard
		if err != nil {
			log.Error("Error unmarshalling json File ", err)
		}
		// write Dashboards as json to a file
		err = os.WriteFile(filepath.Join(path, slug.Make(dashboard.Title+" "+dashboard.UID)+".json"), jsonDashboard, os.FileMode(0666)) // Make sure filename is unique, FileMode is irrelevant, but required for WriteFile
		if err != nil {
			log.Error("Couldn't write Dashboard to disk ", err)
		}
		log.Info("Exported Dashboard " + dashboard.Title)
	}
	return nil
}
