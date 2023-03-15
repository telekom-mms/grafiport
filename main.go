package main

import (
	"flag"
	"os"
)

var (
	username  string
	password  string
	url       string
	directory string
	help      bool
)

func main() {
	//	cmd.Execute()
	if help {
		flag.Usage()
		os.Exit(0)
	}
}

func init() {
	flag.StringVar(&username, "u", "", "Username for Grafana Instance")
	flag.StringVar(&password, "p", "", "Password for Grafana Instance")
	flag.StringVar(&url, "url", "", "Url for Grafana Instance")
	flag.StringVar(&directory, "directory", "", "Directory where Output/Input is stored")
	flag.BoolVar(&help, "h", false, "Die Hilfe")

	flag.Parse()
}
