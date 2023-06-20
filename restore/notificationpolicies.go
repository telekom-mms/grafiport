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

// NotificationPolicies is a function that restores all notification policies to a Grafana instance from a given folder
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func NotificationPolicies(username, password, url, directory string) error {

	var (
		filesInDir              []os.DirEntry
		rawNotificationPolicies []byte
	)
	folderName := "notificationPolicies"
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}

	path := filepath.Join(directory, folderName)
	// path is the based on a provided directory and the naming of sub-folder os our tool
	filesInDir, err = os.ReadDir(path)
	// get all files in provided path
	if err != nil {
		log.Error("Failed to read folder%s\n", err)
		return err
	}
	// looping over found files to restore
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			// read in files to json and Unmarshall them to be Object
			if rawNotificationPolicies, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				log.Error(err)
				continue
			}

			var newNotificationPolicies gapi.NotificationPolicyTree
			if err = json.Unmarshal(rawNotificationPolicies, &newNotificationPolicies); err != nil {
				log.Error(err)
				continue
			}
			// interact with api
			// search for alertRule if exists to determine if update or create
			// if no error then object exists
			_, er := client.NotificationPolicyTree()
			if er == nil {
				// library misses Update Function, so implemented by deleting and overwriting the new Config
				err = client.ResetNotificationPolicyTree()
				if err != nil {
					log.Error("Error updating Notification Policy Tree - resetting (1/2) ", err)
					continue
				}
				err = client.SetNotificationPolicyTree(&newNotificationPolicies)
				if err != nil {
					log.Error("Error updating Policy Tree - creating (1/2) ", err)
				} else {
					log.Info("Updated Policy Tree " + newNotificationPolicies.Receiver)
				}
			} else {
				err = client.SetNotificationPolicyTree(&newNotificationPolicies)
				if err != nil {
					log.Error("Error creating  Notification Policy Tree ", err)
				} else {
					log.Info("Created  Policy Tree " + newNotificationPolicies.Receiver)
				}
			}

		}
	}
	return nil
}
