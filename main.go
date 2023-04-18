package main

import (
	exports "grafana-exporter/export"
	restores "grafana-exporter/restore"
)

func main() {
	if restore {
		restores.Datasources(username, password, url, directory)
		restores.Folders(username, password, url, directory)
		restores.LibaryPanels(username, password, url, directory)
		restores.Dashboards(username, password, url, directory)
	} else {
		exports.Datasources(username, password, url, directory)
		exports.Dashboards(username, password, url, directory)
		exports.Folders(username, password, url, directory)
		exports.LibaryPanels(username, password, url, directory)
		exports.AlertRules(username, password, url, directory)
	}
}
