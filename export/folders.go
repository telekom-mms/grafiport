package export

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafana-exporter/common"
	"os"
	"path/filepath"
)

// Folders is a function that exports all folders from a Grafana instance and stores them as JSON files in a directory.
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func Folders(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "folders"
	path := common.InitializeFolder(directory, folderName)          // initialize Subfolder to export to it
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	//Get slice of all Folders in a short version
	folders, err := client.Folders()
	if err != nil {
		log.Error("Failed to get Folders ", err)
		return err
	}
	// iterate over Folder Slice
	for _, folder := range folders {
		// Get slice of the current folder by UID.
		// UID is in this case the Identifier
		f, err := client.FolderByUID(folder.UID)
		if err != nil {
			log.Error("Error fetching Folder ", err)
		}
		jsonFolder, err := json.Marshal(f) // create JSON Object from Folder
		if err != nil {
			log.Error("Error unmarshalling json File ", err)
		}
		err = os.WriteFile(filepath.Join(path, slug.Make(folder.Title+" "+folder.UID)+".json"), jsonFolder, os.FileMode(0666)) // Make sure Name of File is unique, Filemode is irrelevant, but required for Writefile
		if err != nil {
			log.Error("Couldn't write Folder to disk ", err)
		}
		log.Info("Exported Folder " + folder.Title)
	}
	return nil
}
