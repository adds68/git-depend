package depend

import (
	"sync"
	"time"
)

// Lock allows us to safely write to a note.
type Lock struct {
	sync.Mutex
	ID        string    `json:"Id"`
	Timestamp time.Time `json:"Timestamp"`
	Status    string    `json:"Status"`
}

// WriteLock will attempt to
func (lock *Lock) WriteLock(id string, timestamp time.Time) {
	lock.Lock()
	lock.ID = id
	lock.Timestamp = timestamp
	lock.Status = "LOCKING"
	lock.Unlock()
}

func (lock *Lock) WriteUnlock(id string) {

}
