package depend

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/git-depend/git-depend/pkg/git"
	"github.com/git-depend/git-depend/pkg/utils"
)

type Status string

const (
	Locked Status = "Locked"
)

var ref_lock_name string = "git-depend-lock"

// Lock allows us to safely write to a note.
type Lock struct {
	ID        string    `json:"Id"`
	Timestamp time.Time `json:"Timestamp"`
	Status    string    `json:"Status"`
}

// WriteLock will attempt to
func (lock *Lock) WriteLock(cache *git.Cache, node *Node) error {
	data, err := json.Marshal(lock)
	if err != nil {
		return err
	}

	if err = cache.AddNotes(node.URL, ref_lock_name, string(data)); err != nil {
		return err
	}

	return cache.PushNotes(node.URL, ref_lock_name)
}

func (lock *Lock) WriteUnlock(id string) {

}

func (requests *Requests) WriteRequests(cache *git.Cache) error {
	visited := utils.NewSet()
	lock := &Lock{
		ID:        "foo",
		Timestamp: time.Now(),
		Status:    string(Locked),
	}
	for k := range requests.Table {
		node, ok := requests.NodesRoot.GetNode(k)
		if !ok {
			return errors.New("Key does not exist: " + k)
		}
		if !visited.Exists(k) {
			lock.WriteLock(cache, node)
			visited.Add(k)
			for _, d := range node.GetChildren() {
				if !visited.Exists(d.Name) {
					lock.WriteLock(cache, d)
					visited.Add(d.Name)
				}
			}
		}
	}
	return nil
}
