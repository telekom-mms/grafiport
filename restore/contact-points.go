package restore

import (
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	"grafiport/common"
	"path/filepath"
)

// ContactPoints is a function that restores all contactPoints to a Grafana instance from a given folder
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func ContactPoints(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "contactPoints"
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}

	log.Info("Starting to restore ContactPoints")
	// path is the based on a provided directory and the naming of sub-folder os our tool
	path := filepath.Join(directory, folderName)
	// get all objects from provided path
	ContactPointsSlice, err := common.ReadObjectsFromDisk[gapi.ContactPoint](path)
	if err != nil {
		log.Error("Error reading AlertRules from Disk")
		return err
	}
	// looping over found files to restore
	for _, newContactPoint := range ContactPointsSlice {
		// interact with api
		// search for ContactPoints if exists to determine if update or create
		// if no error then object exists
		_, err = client.ContactPoint(newContactPoint.UID)
		if err == nil {
			err = client.UpdateContactPoint(&newContactPoint)
			if err != nil {
				log.Error("Error updating ContactPoint ", err)
			} else {
				log.Info("Updated ContactPoint " + newContactPoint.Name)
			}

		} else {
			_, err = client.NewContactPoint(&newContactPoint)
			if err != nil {
				log.Error("Error creating ContactPoint ", err)
			} else {
				log.Info("Created ContactPoint " + newContactPoint.Name)
			}
		}
	}
	return err
}
