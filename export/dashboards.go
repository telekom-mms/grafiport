package export

import (
	"context"
	"fmt"
	"github.com/grafana-tools/sdk"
	"os"
	"path/filepath"
)

func Dashboards(credentials string, url string, directory string) {
	var (
		boardLinks []sdk.FoundBoard
		rawBoard   []byte
		meta       sdk.BoardProperties
		err        error
	)
	ctx := context.Background()
	c, err := sdk.NewClient(url, credentials, sdk.DefaultHTTPClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}
	if boardLinks, err = c.Search(ctx, sdk.SearchType(sdk.SearchTypeDashboard)); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	path := filepath.Join(directory, "dashboard")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0760)
	}
	for _, link := range boardLinks {
		if rawBoard, meta, err = c.GetRawDashboardByUID(ctx, link.UID); err != nil {
			fmt.Fprintf(os.Stderr, "%s for %s\n", err, link.URI)
			continue
		}
		if err = os.WriteFile(filepath.Join(path, fmt.Sprintf("%s.json", meta.Slug)), rawBoard, os.FileMode(int(0666))); err != nil {
			fmt.Fprintf(os.Stderr, "%s for %s\n", err, meta.Slug)
		}
	}
}
