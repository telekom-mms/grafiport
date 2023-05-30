package restore

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
	"strings"
)

func DataSources(username, password, url, directory string) error {
	var (
		filesInDir    []os.DirEntry
		rawDatasource []byte
	)
	folderName := "dataSources"
	userInfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userInfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}

	log.Info("Starting to restore dataSources")
	path := filepath.Join(directory, folderName)

	filesInDir, err = os.ReadDir(path)
	if err != nil {
		log.Error("Failed to read folder%s\n", err)
		return err
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawDatasource, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				log.Error(err)
				continue
			}

			var newDatasource gapi.DataSource
			if err = json.Unmarshal(rawDatasource, &newDatasource); err != nil {
				log.Error(err)
				continue
			}
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
				} else {
					log.Info("Updated Datasource " + newDatasource.Name)
				}

			} else {
				_, err = client.NewDataSource(&newDatasource)
				if err != nil {
					log.Error("Error creating Datasource ", err)
				} else {
					log.Info("Created Datasource " + newDatasource.Name)
				}

			}
		}
	}
	return nil
}
