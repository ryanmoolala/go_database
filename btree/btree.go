package main

//This is a Go implementation of a B+ Tree
type Entry struct {
	key string
	value interface{}
}

type TreeNode struct {
	isLeaf bool
	keys []string // separator keys 
	entries []*Entry //for leaf nodes only 
	childNodes []*TreeNode //for internal nodes
	next *TreeNode //leaf level linked list
}

type BTree struct {
	//Think of this as a wrapper for the root node of a B Tree
	root *TreeNode //Each Btree structure contains a node and other essential info like min, max.
	min int 
}

// basic functions eg. create
func createEntry(key string, value interface{}) *Entry {
	return &Entry {
		key: key,
		value: value, 
	}
}

func createTreeNode(bucket *BTree, entries []*Entry, childNodes []*TreeNode, isLeaf bool) *TreeNode {
	return &TreeNode{
		entries: entries,
		childNodes: childNodes,
		isLeaf: isLeaf ,
	}
}

func createEmptyTreeNode() *TreeNode {
	return &TreeNode{
		entries: []*Entry{},
		childNodes: []*TreeNode{},
	}
}

func createTreeWithRoot(root *TreeNode, min int) *BTree {
	var curr_root *BTree = &BTree{
		root: root,
	}
	curr_root.min = min //Root nodes can have a minimum of 1 
	return curr_root
}

func createTree(beta int) *BTree {
	return createTreeWithRoot(createEmptyTreeNode(),beta)
}

//Functions that interact with the tree itself
// Implement Search
// Implement Insert without splitting
// Add splitChild
// Handle root splitting
// Add leaf links
// Print + validate
// (Optional) Delete
// handle duplicate keys

func (tree *BTree) searchElement(key string) (*Entry, bool) {
	//Recursively search for the key, root has isLeaf = False hence don't need to check
	var current *TreeNode = tree.root

	//Search the keys first	
	for !current.isLeaf {
        i := 0
        for i < len(current.keys) && key >= current.keys[i] {
            i++
        }
        current = current.childNodes[i]
    }

    // Search leaf entries
    for _, entry := range current.entries {
        if entry.key == key {
            return entry, true
        }
    }

    return nil, false
}