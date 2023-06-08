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

// ContactPoints is a function that exports all contact points from a Grafana instance and stores them as JSON files in a directory.
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func ContactPoints(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "contactPoints"
	userInfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userInfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}
	log.Info("Starting to export ContactPoints")
	path := filepath.Join(directory, folderName)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0760)
		if err != nil {
			log.Fatal("Error creating directory ", err)
		}
	}
	contactPoints, err := client.ContactPoints()
	if err != nil {
		log.Error("Failed to get ContactPoints ", err)
		return err
	}
	for _, contactPoint := range contactPoints {
		if contactPoint.UID == "" {
			continue
		}
		jsonContactPoint, err := json.Marshal(contactPoint)
		if err != nil {
			log.Error("Error unmarshalling json File ", err)
		}
		err = os.WriteFile(filepath.Join(path, slug.Make(contactPoint.Name))+".json", jsonContactPoint, os.FileMode(0666))
		if err != nil {
			log.Error("Couldn't write ContactPoint to disk ", err)
		} else {
			log.Info("Exported ContactPoint " + contactPoint.Name)
		}
	}
	return nil
}
