package export

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/grafana-tools/sdk"
	"os"
	"path/filepath"
)

func Datasources(credentials string, url string, directory string) {
	var (
		datasources []sdk.Datasource
		dsPacked    []byte
		meta        sdk.BoardProperties
		err         error
	)
	ctx := context.Background()
	c, err := sdk.NewClient(url, credentials, sdk.DefaultHTTPClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}
	if datasources, err = c.GetAllDatasources(ctx); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	path := filepath.Join(directory, "datasources")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0660)
	}
	for _, ds := range datasources {
		if dsPacked, err = json.Marshal(ds); err != nil {
			fmt.Fprintf(os.Stderr, "%s for %s\n", err, ds.Name)
			continue
		}
		if err = os.WriteFile(filepath.Join(path, fmt.Sprintf("%s.json", slug.Make(ds.Name))), dsPacked, os.FileMode(int(0666))); err != nil {
			fmt.Fprintf(os.Stderr, "%s for %s\n", err, meta.Slug)
		}
	}
}
