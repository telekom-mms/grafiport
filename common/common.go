package common

import (
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
)

func InitializeFolder(directory, folderName string) string {
	var err error
	path := filepath.Join(directory, folderName)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0760) //create complete path
		if err != nil {
			log.Fatal("Error creating directory ", err)
		}
	}
	return path
}

func InitializeClient(username, password, url string) (*gapi.Client, error) {
	userInfo := url2.UserPassword(username, password) // User Password Combination for Login to Grafana
	config := gapi.Config{BasicAuth: userInfo}        //Config object passed to client
	client, err := gapi.New(url, config)              // Client used to query Grafana Api
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return nil, err
	}
	return client, nil
}
