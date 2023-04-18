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

func Folders(username, password, url, directory string) {

	var (
		filesInDir []os.DirEntry
		rawFolder  []byte
	)
	foldername := "folders"
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
			if rawFolder, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}

			var newFolder gapi.Folder
			if err = json.Unmarshal(rawFolder, &newFolder); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
			status, _ := client.FolderByUID(newFolder.UID)
			if status.UID != "" {
				client.UpdateFolder(newFolder.UID, newFolder.Title)

			} else {
				client.NewFolder(newFolder.Title, newFolder.UID)
			}
		}
	}
}
