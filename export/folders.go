package export

import (
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafiport/common"
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
	path := common.InitializeFolder(directory, folderName)          // initialize Sub-folder to export to it
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	// Get slice of all Folders in a short version
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
		err = common.WriteObjectToDisk(f, path, slug.Make(f.Title+" "+f.UID)+".json")
		if err != nil {
			log.Error("Couldn't write Folder to disk ", err)
		}
		log.Info("Exported Folder " + folder.Title)
	}
	return nil
}
