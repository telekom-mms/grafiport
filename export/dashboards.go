package export

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
)

// Dashboards is a function that exports all dashboards from a Grafana instance and stores them as JSON files in a directory.
// The function takes four parameters: username, password, url and directory.
// username and password are the credentials for the Grafana instance.
// url is the base URL of the Grafana instance.
// directory is the path of the directory where the dashboards will be stored.
func Dashboards(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "dashboards"                        //Name of sub-folder
	userInfo := url2.UserPassword(username, password) // User Password Combination for Login to Grafana
	config := gapi.Config{BasicAuth: userInfo}        //Config object passed to client
	client, err := gapi.New(url, config)              // Client used to query Grafana Api
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}
	log.Info("Starting to export Dashboard")
	path := filepath.Join(directory, folderName)
	// create sub-folder to put Dashboards in
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0760) //create complete path
		if err != nil {
			log.Fatal("Error creating directory ", err)
		}
	}
	dashboards, err := client.Dashboards() //Get slice of all dashboards in a short version
	if err != nil {
		log.Error("Failed to get Dashboards ", err)
		return err
	}
	// iterate over all Dashboards
	for _, dashboard := range dashboards {
		// Get slice of the current dashboard by UID.
		// UID is in this case the Identifier
		ds, err := client.DashboardByUID(dashboard.UID)
		if err != nil {
			log.Error("Error fetching Dashboard ", err)
		}
		jsonDashboard, err := json.Marshal(ds) // create JSON Object from Dashboard
		if err != nil {
			log.Error("Error unmarshalling json File ", err)
		}
		// write Dashboards as json to a File
		// TODO create a unique and descriptive namingscheme for file
		err = os.WriteFile(filepath.Join(path, slug.Make(dashboard.Title))+".json", jsonDashboard, os.FileMode(0666)) // Name of File is not necessarily unique. So overwrite is possible, but unlikely
		if err != nil {
			log.Error("Couldn't write Dashboard to disk ", err)
		}
		log.Info("Exported Dashboard " + dashboard.Title)
	}
	return nil
}
