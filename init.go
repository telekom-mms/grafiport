package main

import (
	"flag"
	"fmt"
	"os"
	"log"
)

var (
	username  string
	password  string
	url       string
	directory string
	help      bool
	credentials string
)

func init() {
	flag.StringVar(&username, "u", "", "Username for Grafana Instance")
	flag.StringVar(&password, "p", "", "Password for Grafana Instance")
	flag.StringVar(&url, "url", "", "Baseurl for your Grafana Instance.")
	flag.StringVar(&directory, "directory", "", "Directory where Output/Input is stored")
	flag.BoolVar(&help, "h", false, "Die Hilfe")

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
	info, err := os.Stat(directory)
	if os.IsNotExist(err) {
		log.Fatal("No directory provided. \n Please create target folder")
	}

	if !info.IsDir() {
		log.Fatal("Path is not a directory.")
	}
	
	credentials = username + ":" + password
	fmt.Println(username, directory, url, password)
}