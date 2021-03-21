package depend

import "errors"

// Request holds all the information about the merge which is taking place.
type Request struct {
	Name   string `json:"Name"`
	From   string `json:"From"`
	To     string `json:"To"`
	Author string `json:"Author,omitempty"`
	Email  string `json:"Email,omitempty"`
}

// Requests contains a map of Node names to Requests.
type Requests struct {
	Table     map[string]*Request
	NodesRoot *Root
}

type MergeRequestsError struct {
}

func (e *MergeRequestsError) Error() string {
	return ""
}

// NewRequests returns a new Requests struct.
func NewRequests(root *Root) *Requests {
	return &Requests{
		Table:     make(map[string]*Request),
		NodesRoot: root,
	}
}

// AddRequest for merging.
func (requests *Requests) AddRequest(name string, from string, to string, author string, email string) error {
	if _, ok := requests.Table[name]; ok {
		return errors.New("Request already exists")
	}

	if _, ok := requests.NodesRoot.GetNode(name); !ok {
		return errors.New("Node does not exist")
	}

	requests.Table[name] = &Request{
		Name:   name,
		From:   from,
		To:     to,
		Author: author,
		Email:  email,
	}
	return nil
}
