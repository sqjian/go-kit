package consistent_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sqjian/go-kit/consistent"
	"strconv"
	"testing"
)

func TestBalance(t *testing.T) {
	x := consistent.New()
	if x == nil {
		t.Errorf("expected obj")
	}
	var reqs []string
	var nodes = make(map[string]int)

	for i := 0; i < 1e3; i++ {
		reqs = append(reqs, strconv.Itoa(i))
	}
	nodes = map[string]int{
		"node1": 0,
		"node2": 0,
		"node3": 0,
	}

	for node, _ := range nodes {
		x.Add(node)
	}

	for _, req := range reqs {
		node, nodeErr := x.Get(req)
		if nodeErr != nil {
			t.Fatal(nodeErr)
		}
		nodes[node]++
	}

	t.Logf(spew.Sdump(nodes))
}
