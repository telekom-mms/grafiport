package export

import (
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
)

func LibaryPanels(username, password, url, directory string) {
	var (
		err             error
		rawLibaryPanels []gapi.LibraryPanel
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
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0760)
	}
	libarypanels, err := client.LibraryPanels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create search dashboards: %s\n", err)
		os.Exit(1)
	}

	for _, panel := range libarypanels {
		p, _ := client.LibraryPanelByUID(panel.UID)
		rawLibaryPanels = append(rawLibaryPanels, *p) //TODO entfernen - die Lib wird nie verwendet
		jsonLibaryPanels, _ := json.Marshal(p)
		_ = os.WriteFile(filepath.Join(path, slug.Make(panel.Name))+".json", jsonLibaryPanels, os.FileMode(int(0666)))
	}
}
