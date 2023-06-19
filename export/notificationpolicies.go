package export

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafana-exporter/common"
	"os"
	"path/filepath"
)

// NotificationPolicies is a function that exports all policies for notifications from a Grafana instance and stores them as JSON files in a directory.
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func NotificationPolicies(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "notificationPolicies"
	path := common.InitializeFolder(directory, folderName)          // initialize Sub-folder to export to it
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	// Get the complete NotificationPolicyTree
	notificationPolicies, err := client.NotificationPolicyTree()
	// NotificationPolicies are stored as a tree structure
	if err != nil {
		log.Error("Failed to get NotificationPolicies ", err)
		return err
	}

	jsonFolder, err := json.Marshal(notificationPolicies) // create JSON Object from Policy tree
	if err != nil {
		log.Error("Error unmarshalling json File ", err)
	}
	// store the complete tree as one object
	err = os.WriteFile(filepath.Join(path, slug.Make(notificationPolicies.Receiver)+".json"), jsonFolder, os.FileMode(0666)) // Make sure filename is unique, FileMode is irrelevant, but required for WriteFile
	if err != nil {
		log.Error("Couldn't write NotificationPolicies to disk ", err)
	}
	log.Info("Exported NotificationPolicies " + notificationPolicies.Receiver)
	return nil
}
