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

func LibraryPanels(username, password, url, directory string) error {
	var (
		filesInDir []os.DirEntry
		rawPanel   []byte
	)
	folderName := "libraryPanels"
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}

	log.Info("Starting to restore Libary Panel")
	path := filepath.Join(directory, folderName)

	filesInDir, err = os.ReadDir(path)
	if err != nil {
		log.Error("Failed to read folder%s\n", err)
		return err
	}

	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawPanel, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				log.Error(err)
				continue
			}

			var newPanel gapi.LibraryPanel
			if err = json.Unmarshal(rawPanel, &newPanel); err != nil {
				log.Error(err)
				continue
			}
			status, err := client.LibraryPanelByUID(newPanel.UID)
			if err != nil && !(strings.Contains(err.Error(), "library element could not be found")) {
				log.Error("Error getting UID for Library Panel ", err)
			}
			folder, err := client.FolderByUID(newPanel.Meta.FolderUID)
			if err != nil {
				log.Error("Error getting UID for Folder in Library Panel ", err)
			}
			newPanel.Folder = folder.ID
			if status != nil {
				newPanel.Folder = folder.ID
				newPanel.Meta.FolderUID = folder.UID
				newPanel.Meta.FolderName = folder.Title
				newPanel.ID = status.ID
				newPanel.Version = 0
				_, err := client.PatchLibraryPanel(newPanel.UID, newPanel)
				if err != nil {
					log.Error("Error updating Library Panel ", err)
				} else {
					log.Info("Updated  Library Panel " + newPanel.Name)
				}

			} else {

				_, err = client.NewLibraryPanel(newPanel)
				if err != nil {
					log.Error("Error creating Library Panel ", err)
				} else {
					log.Info("Created  Library Panel " + newPanel.Name)
				}

			}
		}
	}
	return nil
}
