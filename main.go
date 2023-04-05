package main

import (
	exports "grafana-exporter/export"
	restores "grafana-exporter/restore"
)

func main() {
	if restore {
		//		restores.Datasources(credentials, url, directory)
		//		time.Sleep(5 * time.Second)
		//		restores.Folders(credentials, url, directory)
		//		time.Sleep(5 * time.Second)
		restores.Dashboards(username, password, url, directory)
	} else {
		exports.Datasources(username, password, url, directory)
		exports.Dashboards(username, password, url, directory)
		exports.Folders(username, password, url, directory)
	}
}
