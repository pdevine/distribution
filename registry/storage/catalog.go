package storage

import (
	"path"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/distribution/context"
	storageDriver "github.com/docker/distribution/registry/storage/driver"
)

// GetRepositories returns a list, or partial list, of repositories in the registry
// Because it's a quite expensive operation, it should only be used when building up
// an initial set of repositories.
func GetRepositories(ctx context.Context, driver storageDriver.StorageDriver, maxEntries int, lastEntry string) ([]string, error) {
	log.Infof("Retrieving up to %d entries of the catalog starting with '%s'", maxEntries, lastEntry)
	var repos []string

	root, err := defaultPathMapper.path(repositoriesRootPathSpec{})
	if err != nil {
		return repos, err
	}

	Walk(ctx, driver, root, func(fileInfo storageDriver.FileInfo) error {
		filePath := fileInfo.Path()

		// lop the base path off
		repoPath := filePath[len(root)+1:]

		_, file := path.Split(repoPath)
		if file == "_layers" {
			repoPath = strings.TrimSuffix(repoPath, "/_layers")
			if repoPath > lastEntry {
				repos = append(repos, repoPath)
			}
			return ErrSkipDir
		} else if strings.HasPrefix(file, "_") {
			return ErrSkipDir
		}

		return nil
	})

	sort.Strings(repos)
	repos = repos[0:min(maxEntries, len(repos))]

	return repos, nil
}

func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}
