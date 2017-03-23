package msa

import (
	"github.com/gyuho/goraph"
	"os"
	"strconv"
	"testing"
)

func TestGraph_MSA_17_D(t *testing.T) {

	// Get graph
	f, err := os.Open("testdata/graph.json")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	g, err := goraph.NewGraphFromJSON(f, "graph_17")
	if err != nil {
		t.Error(err)
	}
	startgstr := g.String()

	// Process graph
	feasible, err := MSA(g, goraph.StringID("D"))
	t.Logf("For MSA test with root %s:\n\tInput: \n%s\n\tFeasible: %v\n\tValid: \n%s\n\tGot: \n%s\n", "D", startgstr, feasible, "NONE", g.String())
	if err != nil {
		t.Errorf("Error while calculating MSA (%v)", err)
	}
}

func TestGraph_MSA_17_C(t *testing.T) {

	// Get graph
	f, err := os.Open("testdata/graph.json")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	g, err := goraph.NewGraphFromJSON(f, "graph_17")
	if err != nil {
		t.Error(err)
	}
	startgstr := g.String()

	// Process graph
	feasible, err := MSA(g, goraph.StringID("C"))
	t.Logf("For MSA test with root %s:\n\tInput: \n%s\n\tFeasible: %v\n\tValid: \n%s\n\tGot: \n%s\n", "C", startgstr, feasible, "NONE", g.String())
	if err != nil {
		t.Errorf("Error while calculating MSA (%v)", err)
	}
}

func TestGraph_MSAAllRoots_17(t *testing.T) {

	// Get graph
	f, err := os.Open("testdata/graph.json")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	g, err := goraph.NewGraphFromJSON(f, "graph_17")
	if err != nil {
		t.Error(err)
	}
	startgstr := g.String()
	t.Logf("Using starting graph:\n%s\n", startgstr)

	// Process graph
	feasible, graph, rootID, err := MSAAllRoots(g)
	t.Logf("Got results:\n\tFeasible: %v\n\tGraph\n%s\n\tRoot: %s\n", feasible, graph.String(), rootID.String())
	if err != nil {
		t.Errorf("Error while calculating MSAAllRoots: %v", err)
	}
}

// Test with all the graphs
func TestGraph_MSAAllRoots_All(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Log("TESTING ALL GRAPHS...")
	for i := 0; i <= 17; i++ {
		t.Logf("Testing graph %d...", i)
		// Get graph
		f, err := os.Open("testdata/graph.json")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		graphNumberStr := strconv.FormatInt(int64(i), 10)
		if i < 10 {
			graphNumberStr = "0" + graphNumberStr
		}

		graphID := "graph_" + graphNumberStr
		g, err := goraph.NewGraphFromJSON(f, graphID)
		if err != nil {
			t.Error(err)
		}
		//startgstr := g.String()
		//t.Logf("Using starting graph:\n%s\n", startgstr)

		// Process graph
		feasible, graph, rootID, err := MSAAllRoots(g)
		if graph != nil && rootID != nil {
			//t.Logf("Got results:\n\t Graph\n%s\n\tRoot: %s\n", graph.String(), rootID.String())
		}
		if err != nil {
			t.Errorf("Error while calculating MSAAllRoots: %v", err)
		}
		t.Logf("DONE, feasability: %v, root: %s", feasible, rootID)
	}
}
