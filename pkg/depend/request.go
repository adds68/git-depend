package depend

// Request holds all the information about the merge which is taking place.
// Dependencies will only contain the URL and Branch.
type Request struct {
	Name   string `json:"Name"`
	From   string `json:"From"`
	To     string `json:"To"`
	Author string `json:"Author,omitempty"`
	Email  string `json:"Email,omitempty"`
}

// NewRequest creates the Request struct from the given fields
func NewRequest(name string, from string, to string, author string, email string) *Request {
	return &Request{
		Name:   name,
		From:   from,
		To:     to,
		Author: author,
		Email:  email,
	}
}
