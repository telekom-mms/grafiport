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

func Datasources(username, password, url, directory string) {
	foldername := "datasources"
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
	datasources, err := client.DataSources()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create search datasources: %s\n", err)
		os.Exit(1)
	}
	for _, datasource := range datasources {
		ds, _ := client.DataSourceByUID(datasource.UID)
		jsonDatasource, _ := json.Marshal(ds)
		_ = os.WriteFile(filepath.Join(path, slug.Make(datasource.Name))+".json", jsonDatasource, os.FileMode(int(0666)))
	}
}
