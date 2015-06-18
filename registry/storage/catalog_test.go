package storage

import (
	"testing"

	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/storage/driver"
	"github.com/docker/distribution/registry/storage/driver/inmemory"
)

func setupFS(t *testing.T) (driver.StorageDriver, []string, context.Context) {
	d := inmemory.New()
	c := []byte("")
	ctx := context.Background()
	rootpath, _ := defaultPathMapper.path(repositoriesRootPathSpec{})

	repos := []string{
		"/foo/a/_layers/1",
		"/foo/b/_layers/2",
		"/bar/c/_layers/3",
		"/bar/d/_layers/4",
		"/foo/d/in/_layers/5",
		"/an/invalid/repo",
		"/bar/d/_layers/ignored/dir/6",
	}

	for _, repo := range repos {
		if err := d.PutContent(ctx, rootpath+repo, c); err != nil {
			t.Fatalf("Unable to put to inmemory fs")
		}
	}

	expected := []string{
		"bar/c",
		"bar/d",
		"foo/a",
		"foo/b",
		"foo/d/in",
	}

	return d, expected, ctx
}

func TestCatalog(t *testing.T) {
	d, expected, ctx := setupFS(t)

	repos, _ := GetRepositories(ctx, d, 100, "")

	if !testEq(repos, expected) {
		t.Errorf("Expected catalog repos err")
	}
}

func TestCatalogInParts(t *testing.T) {
	d, expected, ctx := setupFS(t)

	chunkLen := 2

	repos, _ := GetRepositories(ctx, d, chunkLen, "")
	if !testEq(repos, expected[0:chunkLen]) {
		t.Errorf("Expected catalog first chunk err")
	}

	lastRepo := repos[len(repos)-1]
	repos, _ = GetRepositories(ctx, d, chunkLen, lastRepo)

	if !testEq(repos, expected[chunkLen:chunkLen*2]) {
		t.Errorf("Expected catalog second chunk err")
	}

	lastRepo = repos[len(repos)-1]
	repos, _ = GetRepositories(ctx, d, chunkLen, lastRepo)

	if !testEq(repos, expected[chunkLen*2:chunkLen*3-1]) {
		t.Errorf("Expected catalog third chunk err")
	}

}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for count := range a {
		if a[count] != b[count] {
			return false
		}
	}

	return true
}
