package common

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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

func WriteObjectToDisk(object any, path, filename string) error {
	jsonObject, err := json.Marshal(object) // Create JSON Object of AlertRule from received Bytes
	if err != nil {
		log.Error("Error unmarshalling json File", err)
		return err
	}
	// write Object as json to a File
	err = os.WriteFile(filepath.Join(path, filename), jsonObject, os.FileMode(0666))
	if err != nil {
		log.Error("Couldn't write Object to disk of type", reflect.TypeOf(object), err)
		return err
	} // Make sure filename is unique, FileMode is irrelevant, but required for WriteFile
	return nil
}

// read in files to json and Unmarshall them to be Object
func ReadObjectsFromDisk[T any](path string) ([]T, error) {
	var r []T
	var Object T
	var rawObject []byte
	filesInDir, err := os.ReadDir(path)
	if err != nil {
		log.Error("Failed to read folder%s\n", err)
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			rawObject, err = os.ReadFile(filepath.Join(path, file.Name()))
			if err != nil {
				log.Error("Error reading Objects from Disk")
				return nil, err
			}
			err = json.Unmarshal(rawObject, &Object)
			if err != nil {
				log.Error("Error reading Objects from Disk")
				return nil, err
			}
			r = append(r, Object)
		}
	}
	return r, nil
}
