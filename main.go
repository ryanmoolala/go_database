package main

import (
    "fmt"
    "godatabase/bptree"
)
func main() {
	tree := bptree.CreateNewTree(2)
 
	tree.Insert(bptree.CreateEntry("alice", "engineer"))
	tree.Insert(bptree.CreateEntry("bob", "designer"))
	tree.Insert(bptree.CreateEntry("carol", "manager"))
 
	// Search
	entry, found := tree.Search("bob")
	if found {
		fmt.Println("Found it ", entry.GetKey(), entry.GetValue()) // bob designer
	}

	// Conditional search
	hits := tree.SearchRange(map[string]string{"carol": "<", "alice": "<="})
	for _, e := range hits {
		fmt.Println("conditional search ", e.GetKey(), e.GetValue())
	}
 
	// Delete
	tree.Delete("bob")
 
	// Print tree structure
	bptree.PrintTree(tree.Root())
}
 