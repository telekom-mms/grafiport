package export

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"grafana-exporter/common"
	"os"
	"path/filepath"
)

// DataSources is a function that exports all data sources from a Grafana instance and stores them as JSON files in a directory.
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func DataSources(username, password, url, directory string) error {
	folderName := "dataSources" //Name of sub-folder

	log.Info("Starting to export DataSource")
	path := common.InitializeFolder(directory, folderName)          // initialize Sub-folder to export to it
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}
	dataSources, err := client.DataSources() //Get slice of all DataSources in a short version
	if err != nil {
		log.Error("Failed to create search dataSources ", err)
		return err
	}
	// iterate over all DataSources
	for _, datasource := range dataSources {
		ds, _ := client.DataSourceByUID(datasource.UID) // request the Datasource by UID to get the complete Config
		if err != nil {
			log.Error("Error fetching DataSource from Grafana ", err)
		}
		jsonDatasource, err := json.Marshal(ds) // create JSON Object from Datasource
		if err != nil {
			log.Error("Error unmarshalling json File ", err)
		}
		err = os.WriteFile(filepath.Join(path, slug.Make(datasource.Name+" "+datasource.UID)+".json"), jsonDatasource, os.FileMode(0666)) // Make sure filename is unique, FileMode is irrelevant, but required for WriteFile
		if err != nil {
			log.Error("Couldn't write DataSource to disk ", err)
		}
		log.Info("Exported DataSource " + datasource.Name)
	}
	return nil
}
