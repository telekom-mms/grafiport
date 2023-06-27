package main

import (
	"github.com/charmbracelet/log"
	exports "grafana-exporter/export"
	restores "grafana-exporter/restore"
)

var (
	errors []error
	err    error
)

func main() {
	if restore {
		if datasources {
			err = restores.DataSources(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
		}
		if dashboards {
			err = restores.Folders(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
			err = restores.LibraryPanels(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
			err = restores.Dashboards(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
		}
		if alerting {
			err = restores.AlertRules(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
			err = restores.ContactPoints(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
			err = restores.NotificationPolicies(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
			err = restores.NotificationTemplates(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
		}

		if errors != nil {
			log.Error("Error in Export execution")
		}
	} else {
		if datasources {
			err = exports.DataSources(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
		}
		if dashboards {
			err = exports.Dashboards(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
			err = exports.Folders(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
			err = exports.LibraryPanels(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
		}
		if alerting {
			err = exports.AlertRules(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
			err = exports.ContactPoints(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
			err = exports.NotificationPolicies(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
			err = exports.NotificationTemplates(username, password, url, directory)
			if err != nil {
				errors = append(errors, err)
			}
		}
		if errors != nil {
			log.Error("Error in Restore execution")
		}
	}
}
