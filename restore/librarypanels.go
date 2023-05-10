package restore

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
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
	userInfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userInfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}

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
			if err != nil {
				log.Error("Error getting UID for Library Panel", err)
			}
			folder, err := client.FolderByUID(newPanel.Meta.FolderUID)
			if err != nil {
				log.Error("Error getting UID for Folder in Library Panel", err)
			}
			newPanel.Folder = folder.ID
			if status != nil {
				_, err = client.PatchLibraryPanel(newPanel.UID, newPanel)
				if err != nil {
					log.Error("Error updating Library Panel", err)
				}

			} else {
				_, err = client.NewLibraryPanel(newPanel)
				if err != nil {
					log.Error("Error creating Library Panel", err)
				}
			}
		}
	}
	return nil
}
