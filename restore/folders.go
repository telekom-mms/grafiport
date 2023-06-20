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

// Folders is a function that restores all folders to a Grafana instance from a given folder
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func Folders(username, password, url, directory string) error {

	var (
		filesInDir []os.DirEntry
		rawFolder  []byte
		err        error
	)
	folderName := "folders"
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	// path is the based on a provided directory and the naming of sub-folder os our tool
	log.Info("Starting to restore Folders")
	// get all files in provided path
	path := filepath.Join(directory, folderName)

	filesInDir, err = os.ReadDir(path)
	if err != nil {
		log.Error("Failed to read folder%s\n", err)
		return err
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawFolder, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				log.Error(err)
				continue
			}

			var newFolder gapi.Folder
			if err = json.Unmarshal(rawFolder, &newFolder); err != nil {
				log.Error(err)
				continue
			}
			// interact with api
			// search for alertRule if exists to determine if update or create
			// if no error then object exists
			status, _ := client.FolderByUID(newFolder.UID)
			if status.UID != "" {
				err = client.UpdateFolder(newFolder.UID, newFolder.Title)
				if err != nil {
					log.Error("Error updating Folder ", err)
					break
				}
				log.Info("Updated Folder " + newFolder.Title)

			} else {
				_, err = client.NewFolder(newFolder.Title, newFolder.UID)
				if err != nil {
					log.Error("Error creating Folder ", err)
					break
				}
				log.Info("Created Folder " + newFolder.Title)

			}
		}
	}
	return err
}
