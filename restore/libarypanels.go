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

func LibaryPanels(username, password, url, directory string) {
	var (
		filesInDir []os.DirEntry
		rawPanel   []byte
	)
	foldername := "libarypanels"
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
			if rawPanel, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}

			var newPanel gapi.LibraryPanel
			if err = json.Unmarshal(rawPanel, &newPanel); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
			status, _ := client.LibraryPanelByUID(newPanel.UID)
			folder, _ := client.FolderByUID(newPanel.Meta.FolderUID)
			newPanel.Folder = folder.ID
			if status != nil {
				client.PatchLibraryPanel(newPanel.UID, newPanel)

			} else {
				client.NewLibraryPanel(newPanel)
			}
		}
	}
}
