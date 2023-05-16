package export

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
)

func Folders(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "folders"
	userInfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userInfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}
	log.Info("Starting to export Folders")
	path := filepath.Join(directory, folderName)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0760)
		if err != nil {
			log.Fatal("Error creating directory", err)
		}
	}
	folders, err := client.Folders()
	if err != nil {
		log.Error("Failed to get Folders", err)
		return err
	}
	for _, folder := range folders {
		f, err := client.FolderByUID(folder.UID)
		if err != nil {
			log.Error("Error fetching Folder", err)
		}
		jsonFolder, err := json.Marshal(f)
		if err != nil {
			log.Error("Error unmarshalling json File", err)
		}
		err = os.WriteFile(filepath.Join(path, slug.Make(folder.Title))+".json", jsonFolder, os.FileMode(0666))
		if err != nil {
			log.Error("Couldn't write Folder to disk", err)
		} else {
			log.Info("Exported Folder", folder.Title)
		}
	}
	return nil
}
