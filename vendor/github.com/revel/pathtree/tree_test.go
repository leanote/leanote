package pathtree

import (
	"fmt"
	"reflect"
	"testing"
)

func TestColon(t *testing.T) {
	n := New()

	n.Add("/:first/:second/", 1)
	n.Add("/:first", 2)
	n.Add("/", 3)

	found(t, n, "/", nil, 3)
	found(t, n, "/a", []string{"a"}, 2)
	found(t, n, "/a/", []string{"a"}, 2)
	found(t, n, "/a/b", []string{"a", "b"}, 1)
	found(t, n, "/a/b/", []string{"a", "b"}, 1)

	notfound(t, n, "/a/b/c")
}

func TestStar(t *testing.T) {
	n := New()

	n.Add("/first/second/*star", 1)
	n.Add("/:first/*star/", 2)
	n.Add("/*star", 3)
	n.Add("/", 4)

	found(t, n, "/", nil, 4)
	found(t, n, "/a", []string{"a"}, 3)
	found(t, n, "/a/", []string{"a"}, 3)
	found(t, n, "/a/b", []string{"a", "b"}, 2)
	found(t, n, "/a/b/", []string{"a", "b"}, 2)
	found(t, n, "/a/b/c", []string{"a", "b/c"}, 2)
	found(t, n, "/a/b/c/", []string{"a", "b/c"}, 2)
	found(t, n, "/a/b/c/d", []string{"a", "b/c/d"}, 2)
	found(t, n, "/first/second", []string{"first", "second"}, 2)
	found(t, n, "/first/second/", []string{"first", "second"}, 2)
	found(t, n, "/first/second/third", []string{"third"}, 1)
}

func TestMixedTree(t *testing.T) {
	n := New()

	n.Add("/", 0)
	n.Add("/path/to/nowhere", 1)
	n.Add("/path/:i/nowhere", 2)
	n.Add("/:id/to/nowhere", 3)
	n.Add("/:a/:b", 4)
	n.Add("/not/found", 5)

	found(t, n, "/", nil, 0)
	found(t, n, "/path/to/nowhere", nil, 1)
	found(t, n, "/path/to/nowhere/", nil, 1)
	found(t, n, "/path/from/nowhere", []string{"from"}, 2)
	found(t, n, "/walk/to/nowhere", []string{"walk"}, 3)
	found(t, n, "/path/to/", []string{"path", "to"}, 4)
	found(t, n, "/path/to", []string{"path", "to"}, 4)
	found(t, n, "/not/found", []string{"not", "found"}, 4)
	notfound(t, n, "/path/to/somewhere")
	notfound(t, n, "/path/to/nowhere/else")
	notfound(t, n, "/path")
	notfound(t, n, "/path/")

	notfound(t, n, "")
	notfound(t, n, "xyz")
	notfound(t, n, "/path//to/nowhere")
}

func TestExtensions(t *testing.T) {
	n := New()

	n.Add("/:first/:second.json", 1)
	n.Add("/a/:second.xml", 2)
	n.Add("/:first/:second", 3)

	found(t, n, "/a/b", []string{"a", "b"}, 3)
	found(t, n, "/a/b.json", []string{"a", "b"}, 1)
	found(t, n, "/a/b.xml", []string{"b"}, 2)
	found(t, n, "/a/b.c.xml", []string{"b.c"}, 2)
	found(t, n, "/other/b.xml", []string{"other", "b.xml"}, 3)
}

func TestErrors(t *testing.T) {
	n := New()
	fails(t, n.Add("//", 1), "empty path elements not allowed")
}

func BenchmarkTree100(b *testing.B) {
	n := New()
	n.Add("/", "root")

	// Exact matches
	for i := 0; i < 100; i++ {
		depth := i%5 + 1
		key := ""
		for j := 0; j < depth-1; j++ {
			key += fmt.Sprintf("/dir%d", j)
		}
		key += fmt.Sprintf("/resource%d", i)
		n.Add(key, "literal")
		// b.Logf("Adding %s", key)
	}

	// Wildcards at each level if no exact matches work.
	for i := 0; i < 5; i++ {
		var key string
		for j := 0; j < i; j++ {
			key += fmt.Sprintf("/dir%d", j)
		}
		key += "/:var"
		n.Add(key, "var")
		// b.Logf("Adding %s", key)
	}

	n.Add("/public/*filepath", "static")
	// b.Logf("Adding /public/*filepath")

	queries := map[string]string{
		"/": "root",
		"/dir0/dir1/dir2/dir3/resource4":    "literal",
		"/dir0/dir1/resource97":             "literal",
		"/dir0/variable":                    "var",
		"/dir0/dir1/dir2/dir3/variable":     "var",
		"/public/stylesheets/main.css":      "static",
		"/public/images/icons/an-image.png": "static",
	}

	for query, answer := range queries {
		leaf, _ := n.Find(query)
		if leaf == nil {
			b.Errorf("Failed to find leaf for querY %s", query)
			return
		}
		if leaf.Value.(string) != answer {
			b.Errorf("Incorrect answer for querY %s: expected: %s, actual: %s",
				query, answer, leaf.Value.(string))
			return
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N/len(queries); i++ {
		for k, _ := range queries {
			n.Find(k)
		}
	}
}

func notfound(t *testing.T, n *Node, p string) {
	if leaf, _ := n.Find(p); leaf != nil {
		t.Errorf("Should not have found: %s", p)
	}
}

func found(t *testing.T, n *Node, p string, expectedExpansions []string, val interface{}) {
	leaf, expansions := n.Find(p)
	if leaf == nil {
		t.Errorf("Didn't find: %s", p)
		return
	}
	if !reflect.DeepEqual(expansions, expectedExpansions) {
		t.Errorf("%s: Wildcard expansions (actual) %v != %v (expected)", p, expansions, expectedExpansions)
	}
	if leaf.Value != val {
		t.Errorf("%s: Value (actual) %v != %v (expected)", p, leaf.Value, val)
	}
}

func fails(t *testing.T, err error, msg string) {
	if err == nil {
		t.Errorf("expected an error. %s", msg)
	}
}
