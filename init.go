package main

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	flag.StringVar(&username, "u", "", "Username for Grafana Instance")
	flag.StringVar(&password, "p", "", "Password for Grafana Instance")
	flag.StringVar(&url, "url", "", "Url for Grafana Instance")
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
	fmt.Println(username, directory, url, password)
}