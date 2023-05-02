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

func ContactPoints(username, password, url, directory string) {
	var (
		filesInDir      []os.DirEntry
		rawContactPoint []byte
	)
	fmt.Println("restoring contact points")
	foldername := "contactpoints"
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
			if rawContactPoint, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}

			var newContactPoint gapi.ContactPoint
			if err = json.Unmarshal(rawContactPoint, &newContactPoint); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
			// TODO create special case for default email grafana point
			_, err := client.ContactPoint(newContactPoint.UID)
			fmt.Println(err)
			// TODO explore error handling with http returncodes to correctly handle this case
			if err == nil {
				client.UpdateContactPoint(&newContactPoint)
				println("update contact point")
			} else {
				client.NewContactPoint(&newContactPoint)
			}

		}
	}
}
