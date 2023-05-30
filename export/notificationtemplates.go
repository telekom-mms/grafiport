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

func NotificationTemplates(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "notificationTemplates"
	userinfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userinfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}
	log.Info("Starting to export NotificationTemplates")
	path := filepath.Join(directory, folderName)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0760)
		if err != nil {
			log.Fatal("Error creating directory", err)
		}
	}
	templates, err := client.MessageTemplates()
	if err != nil {
		log.Error("Failed to get NotificationTemplates ", err)
		return err
	}

	for _, template := range templates {
		t, err := client.MessageTemplate(template.Name)
		if err != nil {
			log.Error("Error fetching NotificationTemplate ", err)
		}
		jsonNotificationTemplates, err := json.Marshal(t)
		if err != nil {
			log.Error("Error unmarshalling json File ", err)
		}
		err = os.WriteFile(filepath.Join(path, slug.Make(template.Name))+".json", jsonNotificationTemplates, os.FileMode(0666))
		if err != nil {
			log.Error("Couldn't write NotificationTemplates to disk ", err)
		} else {
			log.Info("Exported NotificationTemplate " + template.Name)
		}
	}
	return nil
}
