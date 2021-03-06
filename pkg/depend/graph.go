package depend

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/git-depend/git-depend/pkg/utils"
)

// Root node only contains its direct dependencies.
type Root struct {
	Table      map[string]*Node
	Deps       []*Node
	reposTable map[string]*repo
}

// Node of the tree.
type Node struct {
	Name string
	URL  string
	Deps []*Node
}

type NodeCycleError struct {
	nodeName    string
	visitedName string
}

func (e *NodeCycleError) Error() string {
	msg := fmt.Sprintf("Cycle detected.\n%s ---> %s ---> %s", e.visitedName, e.nodeName, e.visitedName)
	return msg
}

// repo contains information about the repository
// The direct dependencies in this struct are the names of other repos.
type repo struct {
	Name string   `json:"Name"`
	URL  string   `json:"Url"`
	Deps []string `json:"Deps,omitempty"`
}

// NewGraphFromFile unmarshalls JSON from a file into a graph.
func NewGraphFromFile(path string) (*Root, error) {
	// Read and parse the data.
	reps, err := newReposFromFile(path)
	if err != nil {
		return nil, err
	}

	// Collect all the entries and check for duplications.
	repos := make(map[string]*repo)
	for _, repo := range reps {
		// Populate the table
		if _, ok := repos[repo.Name]; !ok {
			repos[repo.Name] = repo
		} else {
			msg := fmt.Sprintf("Deuplicate key: %s\nPath: %s\n", repo.Name, path)
			return nil, errors.New(msg)
		}
	}

	return newGraphFromRepos(repos)
}

func newGraphFromRepos(repos map[string]*repo) (*Root, error) {
	root := &Root{
		Table:      make(map[string]*Node, len(repos)),
		reposTable: repos,
	}

	_, err := root.topologicalSort()
	if err != nil {
		return nil, err
	}
	root.createGraph()
	return root, nil
}

// Uses Depth First
func (root *Root) topologicalSort() ([]*Node, error) {
	//TODO: Refactor this...
	return root.resolve(root.reposTable, nil)
}

// resolveFromRepos will take the parsed JSON and resolve it into the Root.Table
func (root *Root) resolveFromRepos(repos map[string]*repo, visited *utils.StringSet) ([]*Node, error) {
	for k, v := range repos {
		for _, d := range v.Deps {
			if visited.Exists(d) {
				return nil, &NodeCycleError{k, d}
			}
		}
	}
	return root.resolve(repos, visited)
}

func (root *Root) resolve(repos map[string]*repo, visited *utils.StringSet) ([]*Node, error) {
	var top bool
	if visited == nil {
		top = true
	}

	var nodes []*Node
	for k, v := range repos {
		if top {
			visited = utils.NewSet()
		}
		visited.Add(k)
		node, ok := root.Table[k]
		// Check to see if we have populated this node already.
		if !ok {
			deps := make(map[string]*repo)
			for _, d := range v.Deps {
				deps[d] = root.reposTable[d]
			}
			nod, err := root.resolveFromRepos(deps, visited)
			if err != nil {
				return nil, err
			}
			node = &Node{
				Name: v.Name,
				URL:  v.URL,
				Deps: nod,
			}
			root.Table[k] = node
			nodes = append(nodes, node)
		} else {
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

// createGraph will find the edges and populate the root.
func (root *Root) createGraph() {
	edges := make(map[string]*Node, len(root.Table))
	for k, v := range root.Table {
		edges[k] = v
	}

	// Find edges.
	for _, v := range root.Table {
		for _, d := range v.Deps {
			delete(edges, d.Name)
		}
	}

	// Populate root dependencies.
	root.Deps = make([]*Node, len(edges))
	i := 0
	for _, v := range edges {
		root.Deps[i] = v
		i++
	}
}

// newReposFromFile reads json from a file and returns a list of repos.
func newReposFromFile(path string) ([]*repo, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var repos []*repo
	if err := json.Unmarshal(data, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}
