package msa

import (
	"github.com/gyuho/goraph"
	"os"
	"testing"
)

func TestGraph_MSA_13(t *testing.T) {

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
