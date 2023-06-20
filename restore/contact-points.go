package restore

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	"grafana-exporter/common"
	"os"
	"path/filepath"
	"strings"
)

// ContactPoints is a function that restores all contactPoints to a Grafana instance from a given folder
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func ContactPoints(username, password, url, directory string) error {
	var (
		filesInDir      []os.DirEntry
		rawContactPoint []byte
		err             error
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
	// get all files in provided path
	filesInDir, err = os.ReadDir(path)
	if err != nil {
		log.Error("Failed to read folder%s\n", err)
	}
	// looping over found files to restore
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			// read in files to json and Unmarshall them to be Object
			if rawContactPoint, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				log.Error(err)
				continue
			}

			var newContactPoint gapi.ContactPoint
			if err = json.Unmarshal(rawContactPoint, &newContactPoint); err != nil {
				log.Error(err)
				continue
			}
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
	}
	return err
}
