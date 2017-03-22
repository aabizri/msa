package msa

import (
	"fmt"
	"github.com/gyuho/goraph"
	"log"
)

// contract all cycles
func contract(g goraph.Graph, root goraph.ID, cycles [][]goraph.ID) error {
	// Choose an arbitrary cycle
	if len(cycles) == 0 {
		return fmt.Errorf("contract: WTF, no cycles here")
	}
	c := cycles[0]

	log.Printf("contract: Contracting cycle of IDs: %v", c)

	// Create a new graph
	ng := goraph.NewGraph()
	ng.Init()
	// Create the contracted node
	vc := goraph.NewNode("vc")

	// Add the non-cycle nodes and the contracted node to the graph
	// First add the contracted one
	ok := ng.AddNode(vc)
	if !ok {
		return fmt.Errorf("contract: couldn't add contracted node (id: %s) to graph", vc.String())
	}
	// Now add the non-cycle nodes
	log.Printf("Adding non-cycle nodes...\n")
	for id, node := range g.GetNodes() {
		// If a node isn't in the cycle, add it
		if !idInCycle(c, id) {
			ok := ng.AddNode(node)
			log.Printf("\tadded node %s\n", node.ID().String())
			if !ok {
				return fmt.Errorf("contract: couldn't add node (id: %s) to new graph", id.String())
			}
		}
	}
	log.Printf("added nodes, result: %v\n", ng.GetNodes())

	// Now process the edges
	// Get the list of all edges
	edges, err := GetEdges(g)
	log.Printf("All edges of g: %v\n", edges)
	if err != nil {
		return fmt.Errorf("contract: Error in call to GetEdges: %v", err)
	}

	// Make a memory for pairs of old graph edges - new graph edges
	var ep []edgePair

	// Three cases: (pi(v) is the source of the lowest incoming edge to v
	// Case 1: If (u,v) is an edge in E with u not in C and v in C (an edge coming into the cycle), then include in E' a new edge e =(u,vc), and define w'(e) = w(u,v) - w(pi(v),v).
	// Case 2: If (u,v) is an edge in E with u in C and v not in C (an edge going away from the cycle), then include in E' a new edge e = (vc, v), and define w'(e) = w(u,v) .
	// Case 3: If (u,v) is an edge in E with u not in C and v not in C (an edge unrelated to the cycle), then include in E' a new edge e=(u,v), and define w'(e) = w(u,v) .
	for _, e := range edges {
		sourceID := e.Source().ID()
		targetID := e.Target().ID()
		sourceInCycle := idInCycle(c, e.Source())
		targetInCycle := idInCycle(c, e.Target())
		log.Printf("Processing edge from %s to %s...\n", sourceID.String(), targetID.String())
		switch {
		case !sourceInCycle && targetInCycle:
			log.Printf("CASE 1: For edge %s to %s, as %s isn't in cycle but %s is, add a new edge from %s to %s\n", sourceID.String(), targetID.String(), sourceID.String(), targetID.String(), sourceID.String(), vc.ID().String())
			lowestWeight, err := findLightestIncomingEdgeWeight(g, e.Target().ID())
			if err != nil {
				return fmt.Errorf("contract: error in findLightestIncomingEdgeWeight: %v", err)
			}
			err = ng.AddEdge(e.Source().ID(), vc.ID(), e.Weight()-lowestWeight)
			ep = append(ep, newEdgePair(e, goraph.NewEdge(e.Source(), vc, e.Weight()-lowestWeight)))
		case sourceInCycle && !targetInCycle:
			log.Printf("CASE 2: For edge %s to %s, as %s is in cycle but %s isn't, add a new edge from %s to %s\n", sourceID.String(), targetID.String(), sourceID.String(), targetID.String(), vc.ID().String(), targetID.String())
			err = ng.AddEdge(vc.ID(), e.Target().ID(), e.Weight())
			ep = append(ep, newEdgePair(e, goraph.NewEdge(vc, e.Target(), e.Weight())))
		case !sourceInCycle && !targetInCycle:
			log.Printf("CASE 3: For edge %s to %s, as %s and %s aren't in the cycle, add a new edge from %s to %s\n", sourceID.String(), targetID.String(), sourceID.String(), targetID.String(), sourceID.String(), targetID.String())
			err = ng.AddEdge(e.Source().ID(), e.Target().ID(), e.Weight())
			ep = append(ep, newEdgePair(e, goraph.NewEdge(e.Source(), e.Target(), e.Weight())))
		}
		if err != nil {
			return fmt.Errorf("contract: Error while doing the three case contraction process for %s: %v\n", e.String(), err)
		}
	}
	newedges, _ := GetEdges(ng)
	log.Printf("contract: Created %d new edges", len(newedges))

	// The fun begins, let's GO RECURSIVE WOOHOO
	// And enjoy the ride
	log.Printf("contract: Calling MSA on contracted graph...")
	err = MSA(ng, root)
	if err != nil {
		return fmt.Errorf("contract: Call to MSA (recursion) failed with error: %v", err)
	}
	log.Printf("contract: MSA Call finished")

	// Now, delete the lightest edge going to the corresponding destination of (u,vc)
	// First get that edge
	log.Printf("contract: last step: now last part: find the edge corresponding to (u,vc)\n")
	for _, pair := range ep {
		// If it goes to vc
		log.Printf("Trying pair \n\tOldest: %s\tNewest: %s", pair.oldest.String(), pair.newest.String())
		if pair.newest.Target().ID().String() == vc.ID().String() {
			log.Printf("contract: last step: Correct pair !")
			log.Printf("contract: last step: Recovering lightest incoming edge source to that target...")
			// Get the lightest incoming edge source
			source, err := findLightestIncomingEdgeSource(g, pair.oldest.Target().ID())
			log.Printf("contract: last step: Lightest incoming edge source is %s\n", source.String())
			if err != nil {
				return fmt.Errorf("contract: error in lightestIncomingEdgeSource while recovering for target %s: %v", pair.oldest.Target().ID().String(), err)
			}
			// Remove it
			log.Printf("contract: last step: Deleting edge from %s to %s\n", source.String(), pair.oldest.Target().ID().String())
			err = g.DeleteEdge(source, pair.oldest.Target().ID())
			if err != nil {
				return fmt.Errorf("contract: error while deleting lightest edge to %s, that is %s --> %s: %v", pair.oldest.Target().ID().String(), pair.oldest.Source().ID().String(), pair.newest.Target().ID().String(), err)
			}
		}
	}

	return err
}

// idInCycle returns true if the node is countained in the given cycles
func idInCycle(cycleNodesIDs []goraph.ID, nodeID goraph.ID) bool {
	for _, cycleNodeID := range cycleNodesIDs {
		if cycleNodeID.String() == nodeID.String() {
			return true
		}
	}
	return false
}

// findLightestIncomingEdge finds the edge with the lowest edge incoming to target
func findLightestIncomingEdgeWeight(g goraph.Graph, target goraph.ID) (float64, error) {
	// Get all sources
	sources, err := g.GetSources(target)
	if err != nil {
		return 0, fmt.Errorf("findLightestIncomingEdgeWeight: error while retrieving sources of target %s: %v", target.String(), err)
	}

	// Find the lightest one
	var lightestWeight float64

	for sourceID, _ := range sources {
		// Retrieve the weight of that specific edge
		weight, err := g.GetWeight(sourceID, target)
		if err != nil {
			return 0, fmt.Errorf("findLightestIncomingEdgeWeight: error while getting weight of edge going from %s to %s : %v", sourceID.String(), target.String(), err)
		}

		// If that weight is lighter than the lightest, or if the lightest weight hasn't yet been set
		if weight <= lightestWeight || lightestWeight == 0 {
			lightestWeight = weight
		}
	}

	return lightestWeight, nil
}

// findLightestIncomingEdge finds the edge with the lowest edge incoming to target
func findLightestIncomingEdgeSource(g goraph.Graph, target goraph.ID) (goraph.ID, error) {
	// Get all sources
	sources, err := g.GetSources(target)
	if err != nil {
		return nil, fmt.Errorf("findLightestIncomingEdgeWeight: error while retrieving sources of target %s: %v", target.String(), err)
	}

	// Find the lightest one
	var (
		lightestWeight float64
		lightestSource goraph.ID
	)

	for sourceID, _ := range sources {
		// Retrieve the weight of that specific edge
		weight, err := g.GetWeight(sourceID, target)
		if err != nil {
			return nil, fmt.Errorf("findLightestIncomingEdgeWeight: error while getting weight of edge going from %s to %s : %v", sourceID.String(), target.String(), err)
		}

		// If that weight is lighter than the lightest, or if the lightest weight hasn't yet been set
		if weight <= lightestWeight || lightestWeight == 0 {
			lightestWeight = weight
			lightestSource = sourceID
		}
	}

	return lightestSource, nil
}

type edgePair struct {
	oldest goraph.Edge
	newest goraph.Edge
}

func newEdgePair(oldest goraph.Edge, newest goraph.Edge) edgePair {
	return edgePair{oldest, newest}
}

/*
	Code for future debug
	// Create the new list of nodes
	n := g.GetNodeCount()-len(c)+1 // n is the new amount of nodes (all the nodes of the graph with the nodes of the cycle contracted into a single one)
	v := make([]goraph.Node,n)
*/
