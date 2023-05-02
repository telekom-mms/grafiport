package restore

import (
	"encoding/json"
	"fmt"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
	"strings"
)

func NotificationPolicies(username, password, url, directory string) {

	var (
		filesInDir              []os.DirEntry
		rawNotificationPolicies []byte
	)
	foldername := "notificationpolicies"
	userinfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userinfo}
	client, err := gapi.New(url, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}

	path := filepath.Join(directory, foldername)

	filesInDir, err = os.ReadDir(path)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawNotificationPolicies, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}

			var newNotificationPolicies gapi.NotificationPolicyTree
			if err = json.Unmarshal(rawNotificationPolicies, &newNotificationPolicies); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
			_, err := client.NotificationPolicyTree()
			if err == nil {
				client.ResetNotificationPolicyTree()
				fmt.Println("update NotificationPolicyTree")

				err := client.SetNotificationPolicyTree(&newNotificationPolicies)
				fmt.Println(err)
			} else {
				client.SetNotificationPolicyTree(&newNotificationPolicies)
			}

		}
	}
}
