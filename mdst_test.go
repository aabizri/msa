package mdst

import (
	"fmt"
	"testing"
)

type testPair struct {
	root   Node
	input  *Graph
	output *Graph
}

var pairs []testPair = []testPair{
	{
		root: 3,
		input: &Graph{
			[]Node{0, 1, 2, 3},
			[]Edge{
				Edge{0, 1, 6},
				Edge{3, 0, 1},
				Edge{3, 2, 8},
				Edge{2, 1, 10},
				Edge{1, 2, 10},
				Edge{1, 3, 12},
			},
		},
		output: &Graph{
			[]Node{0, 1, 2, 3},
			[]Edge{
				Edge{3, 0, 1},
				Edge{3, 2, 8},
				Edge{0, 1, 6},
			},
		},
	},
	{
		root: 1,
		input: &Graph{
			[]Node{0, 1, 2, 3},
			[]Edge{
				Edge{0, 1, 6},
				Edge{3, 0, 1},
				Edge{3, 2, 8},
				Edge{2, 1, 10},
				Edge{1, 2, 10},
				Edge{1, 3, 12},
			},
		},
		output: &Graph{
			[]Node{0, 1, 2, 3},
			[]Edge{
				Edge{1, 3, 12},
				Edge{3, 0, 1},
				Edge{3, 2, 8},
			},
		},
	},
}

func sameGraph(a, b *Graph) error {
	if len(a.E) != len(b.E) {
		return fmt.Errorf("Not same amount of edges in both graphs: for a: %d, for b: %d", len(a.E), len(b.E))
	}
	return nil
}

func testMDST(p testPair) (*Graph, error) {
	output, err := MDST(p.input, p.root)
	if err != nil {
		return output, err
	}
	if err := sameGraph(output, p.output); err != nil {
		err = fmt.Errorf("Correct graph for input is invalid:\n\tInput: %v\n\tAlgorithm output: %v\n\tValid output: %v\n\tDiff: %v\n", p.input, output, p.output, err)
	}
	return output, err
}

func Test_MDSTComplete(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	for i, p := range pairs {

		output, err := testMDST(p)
		t.Logf("Pair %d:\n\tRoot:\t%v\n\tInput:\t%v\n\tValid:\t%v\n\tGot:\t%v", i, p.root, p.input, p.output, output)
		if err != nil {
			t.Error(err)
		}
	}
}
