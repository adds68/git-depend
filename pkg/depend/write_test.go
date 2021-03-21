package depend

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/git-depend/git-depend/pkg/git"
)

// Format this date with local URLs.
var test_data_simple_local_graph string = `
[
	{
		"Name": "foo",
		"Url": "%s",
		"Deps": ["bar", "baz"]
	},
	{
		"Name": "bar",
		"Url": "%s",
		"Deps": []
	},
	{
		"Name": "baz",
		"Url": "%s",
		"Deps": []
	}
]
`

// TestWriteRequests clones temporary git repositories into a cache and then writes a lock to them.
func TestWriteRequests(t *testing.T) {
	cache := createLocalGitCache(t)
	urls := []string{createLocalGitRepo(t), createLocalGitRepo(t), createLocalGitRepo(t)}
	data := fmt.Sprintf(test_data_simple_local_graph, urls[0], urls[1], urls[2])

	root, err := NewGraph([]byte(data))
	if err != nil {
		t.Fatal("Could not create graph: " + err.Error())
	}

	// cache_other := createLocalGitCache(t)
	// cache_other.AddNotes(urls[0], ref_lock_name, "Some note.")
	// cache_other.PushNotes(urls[0], ref_lock_name)

	// out, err := cache_other.ShowNotes(urls[0], ref_lock_name)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(string(out))

	// dir := strings.TrimPrefix(urls[0], "file://")
	// git.AddNotes(dir, ref_lock_name, "Some Note")
	// out, err = git.ShowNotes(dir, ref_lock_name)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(string(out))

	requests := NewRequests(root)
	requests.AddRequest("foo", "branch", "main", "Eric", "eric@email.com")
	if err = requests.WriteRequests(cache); err != nil {
		t.Fatal("Could not write requests: " + err.Error())
	}

	out, err := cache.ShowNotes(urls[0], ref_lock_name)
	if err != nil {
		t.Fatal(err)
	}
	lock := &Lock{}
	json.Unmarshal(out, &lock)

	if lock.ID != "foo" {
		t.Fatal("Lock not created: " + lock.ID)
	}

	if lock.Status != string(Locked) {
		t.Fatal("Not locked: " + string(Locked))
	}
}

// Creates a new local git cache in a temporary directory.
func createLocalGitCache(t *testing.T) *git.Cache {
	cache, err := git.NewCache(t.TempDir())
	if err != nil {
		t.Fatal("Failed to create cache: " + err.Error())
	}
	return cache
}

// Creates a git repo in a temp directory.
// Returns file://{dir}
func createLocalGitRepo(t *testing.T) string {
	dir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		t.Fatal("Failed to create local git repo: " + err.Error())
	}

	cmd = exec.Command("git", "config", "user.email", "you@example.com")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		t.Fatal("Failed to create local git repo: " + err.Error())
	}

	cmd = exec.Command("git", "config", "user.name", "Your Name")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		t.Fatal("Failed to create local git repo: " + err.Error())
	}

	emptyFile, err := os.Create(path.Join(dir, "emptyFile.txt"))
	if err != nil {
		t.Fatal("Failed to create file: " + err.Error())
	}
	emptyFile.Close()

	cmd = exec.Command("git", "add", "-A")
	cmd.Dir = dir
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		t.Fatal("Failed to add files: " + err.Error())
	}

	cmd = exec.Command("git", "commit", "-m", "Init.")
	cmd.Dir = dir
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		t.Fatal("Failed to commit files: " + err.Error())
	}

	return "file://" + dir
}
