package depend

import "errors"

// Request holds all the information about the merge which is taking place.
// Dependencies will only contain the URL and Branch.
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

func NewRequests(root *Root) *Requests {
	return &Requests{
		NodesRoot: root,
	}
}

func (requests *Requests) AddRequest(name string, from string, to string, author string, email string) error {
	if _, ok := requests.Table[name]; !ok {
		return errors.New("Request already exists")
	}

	if node := requests.NodesRoot.GetNode(name); node != nil {
		return errors.New("Node does not exist")
	}

	requests.Table[name] = newRequest(name, from, to, author, email)
	return nil
}

// NewRequest creates the Request struct from the given fields
func newRequest(name string, from string, to string, author string, email string) *Request {
	return &Request{
		Name:   name,
		From:   from,
		To:     to,
		Author: author,
		Email:  email,
	}
}
