package export

import (
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafiport/common"
)

// ContactPoints is a function that exports all contact points from a Grafana instance and stores them as JSON files in a directory.
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func ContactPoints(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "contactPoints"
	path := common.InitializeFolder(directory, folderName)          // initialize Sub-folder to export to it
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	contactPoints, err := client.ContactPoints() //Get slice of all ContactPoints
	if err != nil {
		log.Error("Failed to get ContactPoints ", err)
		return err
	}
	// iterate over all ContactPoints to be consistent with all other objects
	for _, contactPoint := range contactPoints {
		if contactPoint.UID == "" {
			continue
		}
		err = common.WriteObjectToDisk(contactPoint, path, slug.Make(contactPoint.Name+" "+contactPoint.UID)+".json")
		if err != nil {
			log.Error("Couldn't write ContactPoint to disk ", err)
		}
		log.Info("Exported ContactPoint " + contactPoint.Name)
	}
	return nil
}
