package main

import (
	"github.com/charmbracelet/log"
	exports "grafana-exporter/export"
	restores "grafana-exporter/restore"
)

var err error

func main() {
	if restore {
		err = restores.DataSources(username, password, url, directory)
		err = restores.Folders(username, password, url, directory)
		err = restores.LibraryPanels(username, password, url, directory)
		err = restores.Dashboards(username, password, url, directory)
		if alerting {
			err = restores.ContactPoints(username, password, url, directory)
			err = restores.NotificationPolicies(username, password, url, directory)
		}
		if err != nil {
			log.Error("Error in Export execution")
		}
	} else {
		err = exports.DataSources(username, password, url, directory)
		err = exports.Dashboards(username, password, url, directory)
		err = exports.Folders(username, password, url, directory)
		err = exports.LibraryPanels(username, password, url, directory)
		if alerting {
			err = exports.AlertRules(username, password, url, directory)
			err = exports.ContactPoints(username, password, url, directory)
			err = exports.NotificationPolicies(username, password, url, directory)
		}
		if err != nil {
			log.Error("Error in Restore execution")
		}
	}
}
