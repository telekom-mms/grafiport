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

func Folders(credentials string, url string, directory string) {
	var (
		filesInDir []os.DirEntry
		folder     sdk.Folder
		rawFolder  []byte
	)
	ctx := context.Background()
	c, err := sdk.NewClient(url, credentials, sdk.DefaultHTTPClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}
	path := filepath.Join(directory, "folders")
	filesInDir, err = os.ReadDir(path)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawFolder, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
		}
		if err = json.Unmarshal(rawFolder, &folder); err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}
		if folder, err = c.CreateFolder(ctx, folder); err != nil {
			fmt.Fprintf(os.Stderr, "error on importing folder %s with %s", folder.Title, err)
		}

	}
}
