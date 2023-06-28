package export

import (
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafiport/common"
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
		err = common.WriteObjectToDisk(ds, path, slug.Make(dashboard.Title+" "+dashboard.UID)+".json")
		if err != nil {
			log.Error("Couldn't write Dashboard to disk ", err)
		}
		log.Info("Exported Dashboard " + dashboard.Title)
	}
	return nil
}
