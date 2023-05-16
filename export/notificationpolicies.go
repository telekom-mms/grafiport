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

func NotificationPolicies(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "notificationPolicies"
	userInfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userInfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}
	log.Info("Starting to export NotificationPolicies")
	path := filepath.Join(directory, folderName)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0760)
		if err != nil {
			log.Fatal("Error creating directory", err)
		}
	}
	notificationPolicies, err := client.NotificationPolicyTree()

	if err != nil {
		log.Error("Failed to get NotificationPolicies", err)
		return err
	}

	jsonFolder, err := json.Marshal(notificationPolicies)
	if err != nil {
		log.Error("Error unmarshalling json File", err)
	}
	err = os.WriteFile(filepath.Join(path, slug.Make(notificationPolicies.Receiver))+".json", jsonFolder, os.FileMode(0666))
	if err != nil {
		log.Error("Couldn't write NotificationPolicies to disk", err)
	} else {
		log.Info("Exported NotificationPolicies", notificationPolicies.Receiver)
	}
	return nil
}
