package msa

import "github.com/gyuho/goraph"

// GetEdges returns all edges from the given Graph
// It is not destructive
// Exported because goraph.Graph doesn't provide it
func GetEdges(g goraph.Graph) ([]goraph.Edge, error) {
	edges := []goraph.Edge{}
	foundEdge := make(map[string]struct{})
	for id1, nd1 := range g.GetNodes() {
		tm, err := g.GetTargets(id1)
		if err != nil {
			return nil, err
		}
		var weight float64
		for id2, nd2 := range tm {
			weight, err = g.GetWeight(id1, id2)
			if err != nil {
				return nil, err
			}
			edge := goraph.NewEdge(nd1, nd2, weight)
			if _, ok := foundEdge[edge.String()]; !ok {
				edges = append(edges, edge)
				foundEdge[edge.String()] = struct{}{}
			}
		}

		sm, err := g.GetSources(id1)
		if err != nil {
			return nil, err
		}
		for id3, nd3 := range sm {
			weight, err := g.GetWeight(id3, id1)
			if err != nil {
				return nil, err
			}
			edge := goraph.NewEdge(nd3, nd1, weight)
			if _, ok := foundEdge[edge.String()]; !ok {
				edges = append(edges, edge)
				foundEdge[edge.String()] = struct{}{}
			}
		}
	}
	return edges, nil
}

// TotalWeight calculates the total weight of a graph
func TotalWeight(g goraph.Graph) (total float64, err error) {
	// Get all edges
	edges, err := GetEdges(g)
	if err != nil {
		return 0, err
	}

	// Sum it all
	for _, e := range edges {
		total += e.Weight()
	}

	return
}
