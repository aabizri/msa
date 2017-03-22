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

import (
	"fmt"
	"github.com/gyuho/goraph"
)

// cleanParallelEdges removes all parallel edges except the lightest one

// removeRootIncoming removes all incoming edges to the root Node
func removeRootIncoming(g goraph.Graph, root goraph.ID) error {
	// Get all sources
	sources, err := g.GetSources(root)
	if err != nil {
		return fmt.Errorf("removeRootIncoming: error while retrieving sources going to root (%s): %v", root.String(), err)
	}

	// Delete these edges
	for sourceID, _ := range sources {
		// Delete the edge going from s to root
		err = g.DeleteEdge(sourceID, root)
		if err != nil {
			return fmt.Errorf("removeRootIncoming: error while deleting edge going from %s to %s (root) : %v", sourceID.String(), root.String(), err)
		}
	}

	// Return
	return err
}

// removeHeavyEdges removes all the incoming edges to target except the lightest one
// DESTRUCTIVE
func removeHeavyEdges(g goraph.Graph, target goraph.ID) error {
	// Get all sources
	sources, err := g.GetSources(target)

	// Find the lightest one
	var (
		lightestEdgeSource goraph.ID
		lightestWeight     float64
	)
	for sourceID, _ := range sources {
		// Retrieve the weight of that specific edge
		weight, err := g.GetWeight(sourceID, target)
		if err != nil {
			return fmt.Errorf("removeHeavyEdges: error while getting weight of edge going from %s to %s : %v", sourceID.String(), target.String(), err)
		}

		// If that weight is lighter than the lightest, or if the lightest weight hasn't yet been set
		if weight <= lightestWeight || lightestWeight == 0 {
			lightestEdgeSource = sourceID
			lightestWeight = weight
		}
	}

	// Delete all the edges going from source indicated in sources and terminating in target
	for sourceID, _ := range sources {
		// Don't delete the edge incoming from the lightest edge source
		if sourceID.String() != lightestEdgeSource.String() {
			// Delete the edge going from sourceID to the target id
			err = g.DeleteEdge(sourceID, target)
			if err != nil {
				return fmt.Errorf("removeHeavyEdges: error removing edge from %s to %s : %v", sourceID.String(), target.String(), err)
			}
		}
	}

	// Return
	return err

}

// removeAllHeavyEdges removes all the edges going to every node that aren't root except the lightest one, and returns a graph made of that
func removeAllHeavyEdges(g goraph.Graph, root goraph.ID) error {
	// Get all nodes
	nodes := g.GetNodes()

	// For every node in the nodes map, call removeHeavyEdges, removing all incoming edges except the lightest
	var err error
	for nodeID, _ := range nodes {
		// Obviously don't remove the heaviest incoming edges coming to root , they are already removed
		if nodeID.String() != root.String() {
			err := removeHeavyEdges(g, nodeID)
			if err != nil {
				return fmt.Errorf("removeAllHeavyEdges: error while removing the heaviest edges going to %s: %v", nodeID.String(), err)
			}
		}
	}

	return err
}

func copyGraph(g goraph.Graph) (goraph.Graph, error) {
	// Create a graph copy
	tmpg := goraph.NewGraph()
	tmpg.Init()
	// Add each node
	fmt.Println("copyGraph: Adding nodes...")
	for _, oldNode := range g.GetNodes() {
		ok := tmpg.AddNode(goraph.NewNode(oldNode.ID().String()))
		fmt.Printf("\tadded %s\n", oldNode.String())
		if ok != true {
			return nil, fmt.Errorf("copyGraph: Error while adding node %s to new graph", oldNode.String())
		}
	}
	// Add every edge of original graph
	oldEdges, err := GetEdges(g)
	if err != nil {
		return nil, fmt.Errorf("copyGraph: Error while retrieving edges: %v", err)
	}
	fmt.Printf("copyGraph: Adding %d edges...\n", len(oldEdges))
	for _, edge := range oldEdges {
		tmpg.AddEdge(edge.Source().ID(), edge.Target().ID(), edge.Weight()) // ID is workaround for badly coded goraph library
		fmt.Printf("\tadded %s", edge.String())
	}

	newEdges, _ := GetEdges(tmpg)
	fmt.Printf("copyGraph: Resulting graph:\n\tNodes: %v\n\tEdges: %v\n", tmpg.GetNodes(), newEdges)
	return tmpg, nil
}

// Given a graph, calculate the MSA
func MSA(g goraph.Graph, root goraph.ID) error {

	// First remove every edge coming into root
	fmt.Printf("Calling removeRootIncoming with parameters:\n\tRoot: %s\n\tGraph: \n%s\n...", root.String(), g.String())
	err := removeRootIncoming(g, root)
	if err != nil {
		return fmt.Errorf("MSA: removeRootIncoming returned erro r: %s", err.Error())
	}
	fmt.Printf("DONE - Results:\n\tRoot: %s\n\tGraph: \n%s\n", root.String(), g.String())

	// Copy the graph
	/*
		fmt.Printf("Calling copy with parameters:\n\tGraph: \n%s\n...",g.String())
		tmpg,err := copyGraph(g)
		if err != nil{
			return fmt.Errorf("MSA: copy returned an error: %v",err)
		}
		fmt.Printf("DONE - Resulting graph: \n%s\n",tmpg.String())

		// Pass this graph to removeAllHeavyEdges
		err = removeAllHeavyEdges(g,root)
		if err != nil{
			fmt.Errorf("MSA: Error in removeAllHeavyEdges: %v",err)
		}*/

	// Now let's check if there are any cycles in that temporary graph
	// First let's retrieve all strongly connected components
	stronglyConnectedComponents := goraph.Tarjan(g)

	// Now let's iterate through the list to check if there are any sublists longer than one, and add them to the list of cycles
	cycles := make([][]goraph.ID, 0)
	for _, l := range stronglyConnectedComponents {
		// If the quantity of strongly connected components is higher than one, then we found a cycle
		if len(l) > 1 {
			cycles = append(cycles, l)
		}
	}

	// If there are no cycles, then we found the minimal spanning arborescence
	if len(cycles) == 0 {
		fmt.Printf("There are no cycles in graph:\n\tRoot: %s\n\tGraph: %s\n", root.String(), g.String())
		return nil
	}
	fmt.Printf("There are cycles in graph:\n\tRoot: %s\n\tGraph: %s\n", root.String(), g.String())

	// If there are, let's contract them
	err = contract(g, root, cycles)
	return err
}
