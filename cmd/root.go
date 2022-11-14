package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	exports "grafana-exporter/cmd/export"
	"log"
	"os"
)

var (
	password    string
	username    string
	credentials string
	directory   string
	url         string
	userLicense string
	rootCmd     = &cobra.Command{
		Use:   "grafana-exporter",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				username = os.Getenv("USERNAME")
				password = os.Getenv("PASSWORD")
				url = os.Getenv("URL")
				directory = os.Getenv("DIRECTORY")
				log.Println(username)
			}
			credentials = username + ":" + password
			fileinfo, err := os.Stat(directory)
			if os.IsNotExist(err) {
				log.Println(fileinfo)
				log.Fatal("No directory provided. \n Please create target folder")

			}
			if !fileinfo.IsDir() {
				log.Fatal("Path is not a directory.")
			}
			log.Println(fileinfo)
			exports.Datasources(credentials, url, directory)
			exports.Dashboards(credentials, url, directory)
			exports.Folders(credentials, url, directory)
			exports.NotificationChannels(credentials, url, directory)

		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Password for the Api")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Username for the Api")
	rootCmd.PersistentFlags().StringVar(&url, "url", "", "Url for the Api")
	rootCmd.PersistentFlags().StringVar(&directory, "directory", "", "name of")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")

}
