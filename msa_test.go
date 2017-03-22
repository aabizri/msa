package msa

import (
	"github.com/gyuho/goraph"
	"os"
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
	err = MSA(g, goraph.StringID("D"))
	t.Logf("For MSA test with root %s:\n\tInput: \n%s\n\tValid: \n%s\n\tGot: \n%s\n", "D", startgstr, "NONE", g.String())
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
	err = MSA(g, goraph.StringID("C"))
	t.Logf("For MSA test with root %s:\n\tInput: \n%s\n\tValid: \n%s\n\tGot: \n%s\n", "C", startgstr, "NONE", g.String())
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
	graph, rootID, err := MSAAllRoots(g)
	t.Logf("Got results:\n\t Graph\n%s\n\tRoot: %s\n", graph.String(), rootID.String())
	if err != nil {
		t.Errorf("Error while calculating MSAAllRoots: %v", err)
	}
}
