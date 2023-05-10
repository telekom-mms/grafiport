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

func ContactPoints(username, password, url, directory string) error {
	var (
		filesInDir      []os.DirEntry
		rawContactPoint []byte
		err             error
	)
	folderName := "contactPoints"
	userInfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userInfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}

	path := filepath.Join(directory, folderName)

	filesInDir, err = os.ReadDir(path)
	if err != nil {
		log.Error("Failed to read folder%s\n", err)
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawContactPoint, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				log.Error(err)
				continue
			}

			var newContactPoint gapi.ContactPoint
			if err = json.Unmarshal(rawContactPoint, &newContactPoint); err != nil {
				log.Error(err)
				continue
			}

			_, err = client.ContactPoint(newContactPoint.UID)
			if err == nil {
				err = client.UpdateContactPoint(&newContactPoint)
				log.Error("Error updating ContactPoint ", err)
				println("update contact point")

			} else {
				_, err = client.NewContactPoint(&newContactPoint)
				log.Error("Error creating ContactPoint ", err)
			}
		}
	}
	return err
}
