package tree

import (
	"fmt"
)

func Build(records []Record) (*Node, error) {

	if len(records) == 0 {
		return nil, nil
	}

	matrix := make([][]Record, len(records))
	var root *Node

	// 1. fill up the matrix
	for _, r := range records {
		node := newNode(r.ID)
		isRoot := node.IsRoot(r.Parent)

		if isRoot && root == nil {
			root = node
		} else if isRoot && root != nil {
			return nil, fmt.Errorf("More than 1 root found")
		} else if node.ID < len(matrix) {
			// we can do this since we know node.ID is always going to be between 0 and len(records)
			matrix[node.ID] = append(matrix[node.ID], r)
		} else {
			return nil, fmt.Errorf("Error. Likely index out of bounds")
		}
	}

	// still no root? error!
	if root == nil {
		return nil, fmt.Errorf("No Root Found")
	}

	failed := 0
	for _, recs := range matrix {
		for _, r := range recs {
			if !root.Insert(r) {
				failed++
			}
		}
	}

	if failed > 0 {
		return nil, fmt.Errorf("Not All Nodes inserted. Current Tree \n %v", matrix)
	}

	return root, nil
}

// A record
type Record struct {
	ID     int
	Parent int
}

// A Node
type Node struct {
	ID       int
	Children []*Node
}

func newNode(id int) *Node {
	return &Node{
		ID: id,
		// Children: make([]*Node, 0),
	}
}

// Insert inserts a record in the node/ node's children
// I first attempt to find the correct parent node using breadth first traversal and then inserting the node if not existing
// Insert returns true if the record was inserted and returns false if the record was not inserted
func (n *Node) Insert(record Record) bool {

	parent := record.Parent

	if parent == n.ID && n.idDoesNotYetExist(record.ID) {
		node := newNode(record.ID)
		n.Children = append(n.Children, node)
		return true
	}

	nodes := make([]*Node, 0)
	nodes = append(nodes, n.Children...)
	var next *Node
	for {

		if len(nodes) == 0 {
			break
		}

		next, nodes = unshift(nodes)
		nodes = append(nodes, next.Children...)
		if next.ID == parent {
			break
		}

	}

	if next != nil {
		return next.Insert(record)
	}

	return false

}

func (n *Node) idDoesNotYetExist(id int) bool {
	for _, child := range n.Children {
		if child.ID == id {
			return false
		}
	}

	return true
}

func (n *Node) IsRoot(parent int) bool {
	if n.ID == parent {
		return true
	}

	return false
}

func unshift(nodes []*Node) (*Node, []*Node) {
	node := nodes[0]
	nodes = nodes[1:]
	return node, nodes
}
