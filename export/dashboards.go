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

func Dashboards(username, password, url, directory string) {
	var (
		err           error
	)
	foldername := "dashboards"
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
	dashboards, err := client.Dashboards()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create search dashboards: %s\n", err)
		os.Exit(1)
	}
	for _, dashboard := range dashboards {
		ds, _ := client.DashboardByUID(dashboard.UID)
		jsonDashboard, _ := json.Marshal(ds)
		_ = os.WriteFile(filepath.Join(path, slug.Make(dashboard.Title))+".json", jsonDashboard, os.FileMode(int(0666)))
	}

}
