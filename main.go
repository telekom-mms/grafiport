package main

import (
	exports "grafana-exporter/export"
	restores "grafana-exporter/restore"
)

func main() {
	if restore {
		restores.Datasources(credentials, url, directory)
	} else {
		exports.Datasources(credentials, url, directory)
		exports.Dashboards(credentials, url, directory)
		exports.Folders(credentials, url, directory)
		exports.NotificationChannels(credentials, url, directory)
	}
}
