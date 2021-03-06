package depend

import (
	"fmt"
	"os"
	"path"
	"testing"
)

var test_data string = `
[
	{
		"Name": "foo",
		"Url": "foo.git.com",
		"Deps": ["bar", "baz"]
	},
	{
		"Name": "bar",
		"Url": "bar.git.com",
		"Deps": []
	},
	{
		"Name": "baz",
		"Url": "baz.git.com",
		"Deps": []
		}
]
`

var test_data_circular string = `
[
	{
		"Name": "foo",
		"Url": "foo.git.com",
		"Deps": ["bar"]
	},
	{
		"Name": "bar",
		"Url": "bar.git.com",
		"Deps": ["baz"]
	},
	{
		"Name": "baz",
		"Url": "baz.git.com",
		"Deps": ["qux"]
	},
	{
		"Name": "qux",
		"Url": "baz.git.com",
		"Deps": ["foo"]
	}
]
`

var test_data_multi_root string = `
[
	{
		"Name": "foo",
		"Url": "foo.git.com",
		"Deps": ["bar", "baz"]
	},
	{
		"Name": "bar",
		"Url": "bar.git.com",
		"Deps": []
	},
	{
		"Name": "baz",
		"Url": "baz.git.com",
		"Deps": []
	},
	{
		"Name": "qux",
		"Url": "qux.git.com",
		"Deps": ["baz"]
	}
]
`

func TestNewTableFromFile(t *testing.T) {
	temp_file := writeJson(t, test_data)

	root, err := NewGraphFromFile(temp_file)
	if err != nil {
		t.Fatal(err)
	}

	if len(root.Table) != 3 {
		t.Fatal("Incorrect number of entries: ", len(root.Table))
	}

	repo, ok := root.Table["foo"]
	if !ok {
		t.Fatal("foo does not exist.")
	}

	if repo.Name != "foo" {
		t.Fatal("Incorrect name: ", repo.Name)
	}

	if repo.URL != "foo.git.com" {
		t.Fatal("Incorrect URL: ", repo.URL)
	}

	if len(repo.Deps) != 2 {
		t.Fatal("Incorrect number of dependencies: ", len(repo.Deps))
	}
}

func TestCreateRoot(t *testing.T) {
	temp_file := writeJson(t, test_data)

	root, err := NewGraphFromFile(temp_file)
	if err != nil {
		t.Fatal(err)
	}

	if len(root.Deps) != 1 {
		t.Fatal("Incorrect number of dependencies: ", len(root.Deps))
	}

	if root.Deps[0].Name != "foo" {
		t.Fatal("Incorrect name: ", root.Deps[0].Name)
	}

	if len(root.Deps[0].Deps) != 2 {
		t.Fatal("Incorrect number of dependencies: ", len(root.Deps[0].Deps))
	}
}

func TestCiruclar(t *testing.T) {
	temp_file := writeJson(t, test_data_circular)
	_, err := NewGraphFromFile(temp_file)
	if err == nil {
		t.Fatal("No cycle detected.")
	}
}

func TestMultiRoot(t *testing.T) {
	temp_file := writeJson(t, test_data_multi_root)

	root, err := NewGraphFromFile(temp_file)
	if err != nil {
		t.Fatal(err)
	}

	if len(root.Deps) != 2 {
		t.Fatal("Incorrect number of dependencies: ", len(root.Deps))
	}

	if root.Deps[0].Name != "foo" && root.Deps[1].Name != "foo" {
		fmt.Println(root.Deps[0])
		fmt.Println(root.Deps[1])
		t.Fatal("foo not declared in root.")
	}
	if root.Deps[0].Name != "qux" && root.Deps[1].Name != "qux" {
		fmt.Println(root.Deps[0])
		fmt.Println(root.Deps[1])
		t.Fatal("qux not declared in root.")
	}
	for _, v := range root.Deps {
		if v.Name == "qux" {
			if v.Deps[0].Name != "baz" {
				t.Fatal("Expected qux to contain a baz dependency.")
			}
			if v.Deps[0].URL != "baz.git.com" {
				t.Fatal("Expected baz to contain a baz.git.com url.")
			}
		}
	}
}

// Writes json data to a file and returns the file path.
func writeJson(t *testing.T, data string) string {
	temp_dir := t.TempDir()
	temp_file := path.Join(temp_dir, "test.json")

	f, err := os.Create(temp_file)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		t.Fatal(err)
	}

	return temp_file
}
