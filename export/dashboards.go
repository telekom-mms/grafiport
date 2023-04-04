package export

import (
	"encoding/json"
	"fmt"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
)

func Dashboards(username, password, url, directory string) {
	var (
		err           error
		rawDashboards []gapi.Dashboard
	)
	userinfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userinfo}
	client, err := gapi.New(url, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}

	dashboards, err := client.Dashboards()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create search dashboards: %s\n", err)
		os.Exit(1)
	}
	for _, dashboard := range dashboards {
		ds, _ := client.DashboardByUID(dashboard.UID)
		rawDashboards = append(rawDashboards, *ds)
		test, _ := json.Marshal(ds.Model)
		err = os.WriteFile(filepath.Join(directory, dashboard.Title)+".json", test, os.FileMode(int(0666)))
	}

}
