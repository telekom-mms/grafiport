package restore

import (
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	"grafiport/common"
	"path/filepath"
	"strings"
)

// DataSources is a function that restores all dataSources to a Grafana instance from a given folder
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func DataSources(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "dataSources"
	client, err := common.InitializeClient(username, password, url) // initialize gapi Client
	if err != nil {
		log.Error("Failed to create gapi client", err)
		return err
	}

	log.Info("Starting to restore dataSources")
	// path is the based on a provided directory and the naming of sub-folder os our tool
	path := filepath.Join(directory, folderName)
	// get all objects from provided path
	DataSourcesSlice, err := common.ReadObjectsFromDisk[gapi.DataSource](path)
	if err != nil {
		log.Error("Error reading AlertRules from Disk")
		return err
	}
	// looping over found files to restore
	for _, newDatasource := range DataSourcesSlice {
		// interact with api
		// search for alertRule if exists to determine if update or create
		// if no error then object exists
		status, err := client.DataSourceByUID(newDatasource.UID)
		if (err != nil) && !(strings.Contains(err.Error(), "Data source not found")) {
			log.Error("Failed Status Check if Datasource already exists")
			continue
		}
		if status != nil {
			newDatasource.ID = status.ID
			err = client.UpdateDataSourceByUID(&newDatasource)
			if err != nil {
				log.Error("Error updating Datasource", err, "Datasource: ", newDatasource)
				break
			}
			log.Info("Updated Datasource " + newDatasource.Name)

		} else {
			_, err = client.NewDataSource(&newDatasource)
			if err != nil {
				log.Error("Error creating Datasource ", err)
				break
			}
			log.Info("Created Datasource " + newDatasource.Name)

		}
	}

	return nil
}
