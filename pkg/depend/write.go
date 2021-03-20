package depend

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/git-depend/git-depend/pkg/git"
)

type Status string

const (
	Locked Status = "Locked"
)

var lock_name string = ""

// Lock allows us to safely write to a note.
type Lock struct {
	ID        string    `json:"Id"`
	Timestamp time.Time `json:"Timestamp"`
	Status    string    `json:"Status"`
}

// WriteLock will attempt to
func (lock *Lock) WriteLock(cache *git.Cache, node *Node) error {
	node.Lock()
	defer node.Unlock()

	data, err := json.Marshal(lock)
	if err != nil {
		return err
	}
	if err = cache.AddNotes(node.URL, lock_name, string(data)); err != nil {
		return err
	}
	return nil
}

func (lock *Lock) WriteUnlock(id string) {

}

func (requests *Requests) WriteRequests(cache *git.Cache) error {
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
		lock.WriteLock(cache, node)
	}
	return nil
}
