/*
Naive Minimal Spanning Arborescence (spanning arborescence of minimum weight) algorithm using Chu–Liu/Edmonds' algorithm
See on wikipedia: https://en.wikipedia.org/wiki/Edmonds'_algorithm

WARNING: Work In Progress
TODO:
	- Implement efficient version
	- Add basic multithreading
	- Use general graph data structure
*/
package msa

type Node uint

type Edge struct {
	From   Node
	To     Node
	Weight float64
}

// Returns true if empty
func (e Edge) Empty() bool {
	return (e.Weight == 0 && e.From == 0 && e.To == 0)
}

type Graph struct {
	N []Node
	E []Edge
}

func (g *Graph) TotalWeight() (total float64) {
	for _, e := range g.E {
		total += e.Weight
	}
	return
}

// Remove
func (g *Graph) removeRootIncoming(root Node) *Graph {
	// Create a new graph, reinitialize edge slice
	ng := *g
	ng.E = make([]Edge, 0)

	// Sift through all edges
	for _, e := range g.E {
		// If an edge goes to root, don't add it
		if e.To != root {
			ng.E = append(ng.E, e)
		}
	}

	// Return
	return &ng
}

// Find the lightest incoming edge
func (g *Graph) findLightestIncomingEdge(to Node) Edge {
	var lie Edge

	// Look at each edge
	for _, e := range g.E {
		// If that edge goes to the node
		if e.To == to {
			// Then if it is lower, replace the lowest weight with it
			if e.Weight < lie.Weight || lie.Empty() {
				lie = e
			}
		}
	}

	return lie
}

// Create a graph made only with the lightest incoming edge
func (g *Graph) removeHeavyEdges(root Node) *Graph {
	// Create a new graph, reinitialize edge slice
	ng := *g
	ng.E = make([]Edge, 0)

	// Sift through all nodes
	for _, n := range g.N {
		// Optimise by not doing it for root
		if n != root {
			// Take the lightest edge for each node and append it
			ng.E = append(ng.E, g.findLightestIncomingEdge(n))
		}
	}

	// Return
	return &ng
}

// Contract all cycles and return a new graph
func (g *Graph) contract() *Graph {
	//cycles := g.cycles()
	return nil
}

func (g *Graph) cyclesPresent() bool {
	//temp
	return false
}

// Given a graph, calculate the MSA
func (g *Graph) MSA(root Node) (*Graph, error) {
	// First remove incoming to root
	g = g.removeRootIncoming(root)

	// Create a graph with only the lightest edges
	g = g.removeHeavyEdges(root)

	// If there are no cycle, we found the right graph
	if !g.cyclesPresent() {
		return g, nil
	}

	// If there are
	g = g.contract()

	return g, nil
}
