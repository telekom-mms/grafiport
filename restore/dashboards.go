package restore

import (
	"fmt"
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	"grafana-exporter/common"
	"path/filepath"
)

// Dashboards is a function that restores all dashboards to a Grafana instance from a given folder
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func Dashboards(username, password, url, directory string) error {

	folderName := "dashboards"
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}

	log.Info("Starting to restore Dashboards")
	// path is the based on a provided directory and the naming of sub-folder os our tool
	path := filepath.Join(directory, folderName)
	// get all objects from provided path
	DashboardsSlice, err := common.ReadObjectsFromDisk[gapi.Dashboard](path)
	if err != nil {
		log.Error("Error reading AlertRules from Disk")
		return err
	}
	// looping over found files to restore
	for _, newDB := range DashboardsSlice {

		// setting some fields in order to correctly restore the object
		// FolderID and ID of the Dashboard are controlled by the Grafana Instance so not unique enough
		newDB.Model["id"] = ""
		newDB.FolderID = 0
		// uid is sometimes missing in the Export Object, so set here
		uid := fmt.Sprint(newDB.Model["uid"])
		// interact with api
		// search for alertRule if exists to determine if update or create
		// if no error then object exists
		exists, _ := client.DashboardByUID(uid)
		if exists != nil {
			// library misses Update Function, so implemented by deleting and creating the new Config
			err = client.DeleteDashboardByUID(uid)
			if err != nil {
				log.Error("Error updating Dashboard - delete (1/2) ", err)
				continue
			}
			_, err = client.NewDashboard(newDB)
			if err != nil {
				log.Error("Error updating Dashboard - create (2/2) ", err)
				continue
			}
			log.Info("Updated Dashboard " + fmt.Sprint(newDB.Model["title"]))

		} else {
			_, err = client.NewDashboard(newDB)
			if err != nil {
				log.Error("Error creating Dashboard ", err)
			} else {
				log.Info("Created Dashboard " + fmt.Sprint(newDB.Model["title"]))
			}

		}

	}
	return nil
}
