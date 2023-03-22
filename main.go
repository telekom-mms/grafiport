package main

import(
	exports "grafana-exporter/export"
)

func main() {
	exports.Datasources(credentials, url, directory)
	exports.Dashboards(credentials, url, directory)
	exports.Folders(credentials, url, directory)
	exports.NotificationChannels(credentials, url, directory)
}
