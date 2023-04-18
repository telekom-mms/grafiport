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

func AlertRules(username, password, url, directory string) {
	var (
		err error
	)
	foldername := "alertrules"
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
	// wir benötigen erst alle Folder mit den ID's bevor wir uns die AlertRules holen können, da diese an den folder hängen
	alertRules, err := client.AlertRules()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create search dashboards: %s\n", err)
		os.Exit(1)
	}
	for _, alertRule := range alertRules {
		jsonFolder, _ := json.Marshal(alertRule)
		_ = os.WriteFile(filepath.Join(path, slug.Make(alertRule.Title))+".json", jsonFolder, os.FileMode(int(0666)))
	}
}
