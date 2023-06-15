package export

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafana-exporter/common"
	"os"
	"path/filepath"
)

func NotificationPolicies(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "notificationPolicies"
	path := common.InitializeFolder(directory, folderName)          // initialize Subfolder to export to it
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	notificationPolicies, err := client.NotificationPolicyTree()

	if err != nil {
		log.Error("Failed to get NotificationPolicies ", err)
		return err
	}

	jsonFolder, err := json.Marshal(notificationPolicies)
	if err != nil {
		log.Error("Error unmarshalling json File ", err)
	}
	err = os.WriteFile(filepath.Join(path, slug.Make(notificationPolicies.Receiver)+".json"), jsonFolder, os.FileMode(0666))
	if err != nil {
		log.Error("Couldn't write NotificationPolicies to disk ", err)
	} else {
		log.Info("Exported NotificationPolicies " + notificationPolicies.Receiver)
	}
	return nil
}
