package restore

import (
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	"grafana-exporter/common"
	"path/filepath"
	"strings"
)

// LibraryPanels is a function that restores all libraryPanel to a Grafana instance from a given folder
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func LibraryPanels(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "libraryPanels"
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}

	log.Info("Starting to restore Library Panel")
	// path is the based on a provided directory and the naming of sub-folder os our tool
	path := filepath.Join(directory, folderName)
	// get all objects from provided path
	LibraryPanelsSlice, err := common.ReadObjectsFromDisk[gapi.LibraryPanel](path)
	if err != nil {
		log.Error("Error reading AlertRules from Disk")
		return err
	}
	// looping over found files to restore
	for _, newPanel := range LibraryPanelsSlice {
		// interact with api
		// search for alertRule if exists to determine if update or create
		// if no error then object exists
		status, err := client.LibraryPanelByUID(newPanel.UID)
		if err != nil && !(strings.Contains(err.Error(), "library element could not be found")) {
			log.Error("Error getting UID for Library Panel ", err)
		}
		// check if a corresponding folder for the LibraryPanel exists

		folder, err := client.FolderByUID(newPanel.Meta.FolderUID)
		if err != nil {
			log.Error("Error getting UID for Folder in Library Panel ", err)
		}
		// folder id has to be set in the libraryPanel Config
		// folder id is only instance unique
		newPanel.Folder = folder.ID
		if status != nil {
			// the Panel other FolderVariables and IDs have to be overwritten in order to restore correctly
			// this is to keep references in sync
			newPanel.Meta.FolderUID = folder.UID
			newPanel.Meta.FolderName = folder.Title
			newPanel.ID = status.ID
			newPanel.Version = 0
			_, err := client.PatchLibraryPanel(newPanel.UID, newPanel)
			if err != nil {
				log.Error("Error updating Library Panel ", err)
				break
			}
			log.Info("Updated  Library Panel " + newPanel.Name)

		} else {

			_, err = client.NewLibraryPanel(newPanel)
			if err != nil {
				log.Error("Error creating Library Panel ", err)
				break
			}
			log.Info("Created  Library Panel " + newPanel.Name)

		}
	}
	return nil
}
