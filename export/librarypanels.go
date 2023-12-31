package export

import (
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafiport/common"
)

// LibraryPanels is a function that exports all folders from a Grafana instance and stores them as JSON files in a directory.
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func LibraryPanels(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "libraryPanels"
	path := common.InitializeFolder(directory, folderName)          // initialize Sub-folder to export to it
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	// Get slice of all LibraryPanels in a short form
	libraryPanels, err := client.LibraryPanels()
	if err != nil {
		log.Error("Failed to get LibraryPanels ", err)
		return err
	}
	// iterate over LibraryPanel Slice
	for _, panel := range libraryPanels {
		// Get slice of the current panels by UID.
		// UID is in this case the Identifier
		p, _ := client.LibraryPanelByUID(panel.UID)
		if err != nil {
			log.Error("Error fetching LibraryPanel ", err)
		}
		err = common.WriteObjectToDisk(p, path, slug.Make(p.Name+" "+p.UID)+".json")
		if err != nil {
			log.Error("Couldn't write Dashboard to disk ", err)
		}
		log.Info("Exported Library Panels " + panel.Name)
	}
	return nil
}
