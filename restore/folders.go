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

func Folders(username, password, url, directory string) error {

	var (
		filesInDir []os.DirEntry
		rawFolder  []byte
		err        error
	)
	folderName := "folders"
	userInfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userInfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}

	log.Info("Starting to restore Folders")
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
			status, _ := client.FolderByUID(newFolder.UID)
			if status.UID != "" {
				err = client.UpdateFolder(newFolder.UID, newFolder.Title)
				if err != nil {
					log.Error("Error updating Folder", err)
				} else {
					log.Info("Updated Folder" + newFolder.Title)
				}
			} else {
				_, err = client.NewFolder(newFolder.Title, newFolder.UID)
				if err != nil {
					log.Error("Error creating Folder", err)
				} else {
					log.Info("Created Folder" + newFolder.Title)
				}
			}
		}
	}
	return err
}
