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
	"log"
)

// removeRootIncoming removes all incoming edges to the root Node
func removeRootIncoming(g goraph.Graph, root goraph.ID) error {
	log.Printf("removeRootIncoming: called with root %s", root.String())
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

	log.Printf("removeRootIncoming: removed %v", sources)

	// Return
	return err
}

// removeHeavyEdges removes all the incoming edges to target except the lightest one
// DESTRUCTIVE
func removeHeavyEdges(g goraph.Graph, root goraph.ID, target goraph.ID) error {
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

// removeAllHeavyEdges removes all the edges going to every node that aren't root except the lightest one
func removeAllHeavyEdges(g goraph.Graph, root goraph.ID) error {
	// Get all nodes
	nodes := g.GetNodes()

	// For every node in the nodes map, call removeHeavyEdges, removing all incoming edges except the lightest
	var err error
	for nodeID, _ := range nodes {
		// Obviously don't remove the heaviest incoming edges coming to root , they are already removed
		if nodeID.String() != root.String() {
			log.Printf("removeAllHeavyEdges: Calling removeHeavyEdges for %s", nodeID.String())
			err := removeHeavyEdges(g, root, nodeID)
			if err != nil {
				return fmt.Errorf("removeAllHeavyEdges: error while removing the heaviest edges going to %s: %v", nodeID.String(), err)
			}
		}
	}

	return err
}

// copyInPlace copies a graph into one that already exists
func copyInPlace(source goraph.Graph, target goraph.Graph) error {
	target.Init()
	// Add each node
	log.Println("copyInPlace: Adding nodes...")
	for _, oldNode := range source.GetNodes() {
		ok := target.AddNode(goraph.NewNode(oldNode.ID().String()))
		log.Printf("\tadded %s\n", oldNode.String())
		if ok != true {
			return fmt.Errorf("copyInPlace: Error while adding node %s to new graph", oldNode.String())
		}
	}

	// Add every edge of original graph
	oldEdges, err := GetEdges(source)
	if err != nil {
		return fmt.Errorf("copyGraph: Error while retrieving edges: %v", err)
	}
	log.Printf("copyGraph: Adding %d edges...\n", len(oldEdges))
	for _, edge := range oldEdges {
		target.AddEdge(edge.Source().ID(), edge.Target().ID(), edge.Weight()) // ID is workaround for badly coded goraph library
		log.Printf("\tadded %s", edge.String())
	}

	newEdges, _ := GetEdges(target)
	log.Printf("copyGraph: created %d edges", len(newEdges))
	//log.Printf("copyGraph: Resulting graph:\n\tNodes: %v\n\tEdges: %v\n", tmpg.GetNodes(), newEdges)
	return nil
}

func copyGraph(g goraph.Graph) (goraph.Graph, error) {
	// Create a graph copy
	tmpg := goraph.NewGraph()

	// Call copyInPlace
	err := copyInPlace(g, tmpg)
	return tmpg, err
}

// Given a graph, calculate the MSA
func MSA(g goraph.Graph, root goraph.ID) error {

	// First remove every edge coming into root
	//log.Printf("Calling removeRootIncoming with parameters:\n\tRoot: %s\n\tGraph: \n%s\n...", root.String(), g.String())
	log.Print("MSA: Calling removeRootIncoming...")
	err := removeRootIncoming(g, root)
	if err != nil {
		return fmt.Errorf("MSA: removeRootIncoming returned error: %s", err.Error())
	}
	log.Print("MSA: removeRootIncoming DONE")
	log.Printf("MSA: current graph:\n%s", g.String())

	// Create a dummy graph
	ng, err := copyGraph(g)
	if err != nil {
		return fmt.Errorf("MSA: error while copying graph: %v", err)
	}

	// Now remove all but the heaviest incoming edges on a dummy graph
	log.Print("MSA: Calling removeAllHeavyEdges")
	err = removeAllHeavyEdges(ng, root)
	if err != nil {
		return fmt.Errorf("MSA: removeAllEdges returned error: %s", err.Error())
	}
	log.Print("MSA: removeAllHeavyEdges DONE")
	log.Printf("MSA: current graph:\n%s", g.String())

	// Now let's check if there are any cycles in that graph
	// First let's retrieve all strongly connected components
	log.Print("MSA: calling Tarjan")
	stronglyConnectedComponents := goraph.Tarjan(ng)
	log.Print("MSA: tarjan returned")

	// Now let's iterate through the list to check if there are any sublists longer than one, and add them to the list of cycles
	cycles := make([][]goraph.ID, 0)
	for _, l := range stronglyConnectedComponents {
		// If the quantity of strongly connected components is higher than one, then we found a cycle
		if len(l) > 1 {
			cycles = append(cycles, l)
		}
	}
	log.Printf("MSA: found %d cycles: %v", len(cycles), cycles)

	// If there are no cycles, then we found the minimal spanning arborescence
	if len(cycles) == 0 {
		log.Print("MSA: No cycles found, returning...")
		return copyInPlace(ng, g)
	}

	// If there are, let's contract them
	log.Print("MSA: Calling contract...")
	err = contract(g, root, cycles)
	return err
}

// MSAAllRoots calls MSA with every possible root to find the lightest one
func MSAAllRoots(g goraph.Graph) (lightestGraph goraph.Graph, rootID goraph.ID, err error) {
	// Retrieve list of all nodes
	nodes := g.GetNodes()

	// Create a list of graphs
	var (
		lowestWeight float64
	)
	for id, _ := range nodes {
		ng, err := copyGraph(g)
		if err != nil {
			return nil, nil, err
		}
		err = MSA(ng, id)
		if err != nil {
			return nil, nil, err
		}

		totalWeight, err := TotalWeight(ng)
		if err != nil {
			return nil, nil, err
		}
		if totalWeight <= lowestWeight || lowestWeight == 0 {
			lowestWeight = totalWeight
			lightestGraph = ng
			rootID = id
		}
	}

	return
}
