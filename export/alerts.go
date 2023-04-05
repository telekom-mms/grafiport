package export

import (
	"encoding/json"
	"fmt"
	"log"
	url2 "net/url"
	"os"
	"path/filepath"

	"github.com/gosimple/slug"
	gapi "github.com/grafana/grafana-api-golang-client"
)

func Alerts(username, password, url, directory string) {
	var (
		err        error
	)
	foldername := "alerts"
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

// wir benötigen erst alle Folder mit den ID's bevor wir uns die AlertRules holen können, da diese an den folder hängen
	folders, err := client.Folders()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create search dashboards: %s\n", err)
		os.Exit(1)
	}
	for _, folder := range folders {
		f, _ := client.FolderByUID(folder.UID)
		title := f.Title
		fuid := f.UID
		fmt.Println(fuid, title)
		arg, err := client.AlertRuleGroup(fuid, title)
		if err != nil {
			log.Printf("Konnte AlertRuleGroup nicht beziehen %v %v", fuid, title)
			continue
		}
		ar, err := client.AlertRule(arg.FolderUID)
		if err != nil {
			log.Fatalf("Konnte die AlertRule nicht beziehen %v", err)
		}
		jsonFolder, _ := json.Marshal(ar)
		_ = os.WriteFile(filepath.Join(path, slug.Make(ar.Title))+".json", jsonFolder, os.FileMode(int(0666)))
	}
}
