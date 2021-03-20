package depend

import "testing"

func TestNewRequest(t *testing.T) {
	root := createSimpleGraph(t, test_data_simple_graph)
	reqs := NewRequests(root)

	if err := reqs.AddRequest("foo", "branch", "main", "Test", "test@test.com"); err != nil {
		t.Fatal("Could not add a request: " + err.Error())
	}
}

func TestNewRequestFail(t *testing.T) {
	root := createSimpleGraph(t, test_data_simple_graph)
	reqs := NewRequests(root)

	if err := reqs.AddRequest("no-key", "branch", "main", "Test", "test@test.com"); err == nil {
		t.Fatal("Should not be able to add a new request.")
	}
}
