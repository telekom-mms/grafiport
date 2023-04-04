package restore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grafana-tools/sdk"
	"os"
	"path/filepath"
	"strings"
)

func Dashboards(credentials string, url string, directory string) {

	var (
		filesInDir []os.DirEntry
		rawDB      []byte
	)
	ctx := context.Background()
	c, err := sdk.NewClient(url, credentials, sdk.DefaultHTTPClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}

	path := filepath.Join(directory, "dashboard")

	filesInDir, err = os.ReadDir(path)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawDB, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
			var newDB sdk.Board
			if err = json.Unmarshal(rawDB, &newDB); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
			if status, err := c.SetRawDashboard(ctx, rawDB); err != nil {
				fmt.Fprintf(os.Stderr, "error on importing dashboard %s with %s (%s)", newDB.Title, err, *status.Message)
			}

		}
	}
}
