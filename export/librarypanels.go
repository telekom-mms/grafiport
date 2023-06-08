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
	userinfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userinfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}
	log.Info("Starting to export LibraryPanels")
	path := filepath.Join(directory, folderName)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0760)
		if err != nil {
			log.Fatal("Error creating directory ", err)
		}
	}
	libraryPanels, err := client.LibraryPanels()
	if err != nil {
		log.Error("Failed to get LibraryPanels ", err)
		return err
	}

	for _, panel := range libraryPanels {
		p, _ := client.LibraryPanelByUID(panel.UID)
		if err != nil {
			log.Error("Error fetching LibraryPanel ", err)
		}
		jsonLibraryPanels, err := json.Marshal(p)
		if err != nil {
			log.Error("Error unmarshalling json File ", err)
		}
		err = os.WriteFile(filepath.Join(path, slug.Make(panel.Name))+".json", jsonLibraryPanels, os.FileMode(0666))
		if err != nil {
			log.Error("Couldn't write Dashboard to disk ", err)
		} else {
			log.Info("Exported Library Panels " + panel.Name)
		}
	}
	return nil
}
