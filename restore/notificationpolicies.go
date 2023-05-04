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

func NotificationPolicies(username, password, url, directory string) error {

	var (
		filesInDir              []os.DirEntry
		rawNotificationPolicies []byte
	)
	folderName := "notificationPolicies"
	userinfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userinfo}
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
			if rawNotificationPolicies, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				log.Error(err)
				continue
			}

			var newNotificationPolicies gapi.NotificationPolicyTree
			if err = json.Unmarshal(rawNotificationPolicies, &newNotificationPolicies); err != nil {
				log.Error(err)
				continue
			}
			_, er := client.NotificationPolicyTree()
			if er == nil {
				err = client.ResetNotificationPolicyTree()
				if err != nil {
					log.Error("Error updating Notification Policy Tree - resetting (1/2)", err)
					continue
				}
				err = client.SetNotificationPolicyTree(&newNotificationPolicies)
				if err != nil {
					log.Error("Error updating Policy Tree - creating (1/2)", err)
				}
			} else {
				err = client.SetNotificationPolicyTree(&newNotificationPolicies)
				if err != nil {
					log.Error("Error creating  Notification Policy Tree", err)
				}
			}

		}
	}
	return nil
}
