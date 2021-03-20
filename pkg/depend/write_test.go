package depend

import (
	"os/exec"
	"testing"

	"github.com/git-depend/git-depend/pkg/git"
)

func TestWriteRequests(t *testing.T) {
	createLocalGitCache(t)
}

// Creates a new local git cache in a temporary directory.
func createLocalGitCache(t *testing.T) *git.Cache {
	dir := t.TempDir()
	cmd := exec.Command("git", "-C", dir, "init")
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	path := "file://" + dir
	cache, _ := git.NewCache(path)
	return cache
}
