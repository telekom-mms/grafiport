package main

import (
	"flag"
	"github.com/charmbracelet/log"
	"os"
	"strconv"
)

var (
	username    string
	password    string
	url         string
	directory   string
	help        bool
	restore     bool
	alerting    bool
	dashboards  bool
	datasources bool
)

func init() {
	flag.StringVar(&username, "u", "", "Username for Grafana Instance")
	flag.StringVar(&password, "p", "", "Password for Grafana Instance")
	flag.StringVar(&url, "url", "", "Baseurl for your Grafana Instance.")
	flag.StringVar(&directory, "directory", "", "Directory where Output/Input is stored")
	flag.BoolVar(&help, "h", false, "The Help")
	flag.BoolVar(&restore, "r", false, "Restore of provided backup Directory")
	flag.BoolVar(&alerting, "alerting", false, "Export or Restore of the Alerting Objects including Alert Rules, Contact Point, Notification Policies and Notification Templates")
	flag.BoolVar(&dashboards, "dashboard", true, "Export or Restore of the Dashboards, Folders and Library Panels. Default is always true")
	flag.BoolVar(&datasources, "dashboard", true, "Export or Restore of the Datasources. Default is always true")

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if username == "" {
		username = os.Getenv("username")
	}
	if password == "" {
		password = os.Getenv("password")
	}
	if url == "" {
		url = os.Getenv("url")
	}
	if directory == "" {
		directory = os.Getenv("directory")
	}
	if !alerting {
		alerting = getEnvBool("alerting")
	}
	if dashboards != getEnvBool("dashboard") {
		dashboards = getEnvBool("dashboard")
	}
	if datasources != getEnvBool("datasources") {
		datasources = getEnvBool("datasources")
	}
	if !restore {
		restore = getEnvBool("restore")
	}
	info, err := os.Stat(directory)
	if os.IsNotExist(err) {
		log.Fatal("No directory provided. \n Please create target folder")
	}

	if !info.IsDir() {
		log.Fatal("Path is not a directory.")
	}
}

func getEnvBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		return false
	}
	binary, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatal("Not able to parse Boolean ", key)
	}
	return binary
}
