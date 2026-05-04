package bptree

import (
	"fmt"
	"strings"
)

type BTree struct {
	root *TreeNode //Each Btree structure contains a node and other essential info like min, max.
	min int 
}

/* basic functions for tree itself */

func printTree(root *TreeNode) {
	if root == nil || (root.isLeaf && len(root.leafEntries) == 0) {
		fmt.Println("(empty tree)")
		return
	}
 
	// BFS level-order traversal
	currentLevel := []*TreeNode{root}
	for len(currentLevel) > 0 {
		var nextLevel []*TreeNode
		var sb strings.Builder
 
		for _, node := range currentLevel {
			sb.WriteString("[")
			if node.isLeaf {
				for i, e := range node.leafEntries {
					if i > 0 {
						sb.WriteString(" ")
					}
					sb.WriteString(e.key)
				}
			} else {
				for i, k := range node.keys {
					if i > 0 {
						sb.WriteString(" ")
					}
					sb.WriteString(k)
				}
				nextLevel = append(nextLevel, node.childNodes...)
			}
			sb.WriteString("] ")
		}
 
		fmt.Println(strings.TrimRight(sb.String(), " "))
		currentLevel = nextLevel
	}
 
	// Walk the leaf linked list for sorted entries
	var sb strings.Builder
	sb.WriteString("leaves: ")
	leftmost := root
	for !leftmost.isLeaf {
		leftmost = leftmost.childNodes[0]
	}
	first := true
	for node := leftmost; node != nil; node = node.next {
		for _, e := range node.leafEntries {
			if !first {
				sb.WriteString(" ")
			}
			sb.WriteString(e.key)
			first = false
		}
	}
	fmt.Println(sb.String())
}


// create_new_tree creates a B+ tree with a given order.
// order defines the max keys per node (2*order) and min keys per
// non-root node (order). Must be at least 2.
func create_new_tree(order int) *BTree {
	if order < 2 {
		panic("bptree: order must be at least 2")
	}
	emptyLeaf := &TreeNode{
		isLeaf:      true,
		keys:        make([]string, 0),
		leafEntries: make([]*Entry, 0),
		childNodes:  make([]*TreeNode, 0),
		order:       order,
	}
    
	return &BTree{
		root: emptyLeaf,
		min:  order,
	}
}

func create_new_tree_bulkload(entries []*Entry, order int) *BTree {
	if order < 2 {
		panic("bptree: order must be at least 2")
	}
	if len(entries) == 0 {
		panic("bptree: bulkload requires at least one entry")
	}
 
	// --- Step 1: build leaf nodes ---
	maxLeaf := order * 2
	var leaves []*TreeNode
	for i := 0; i < len(entries); i += maxLeaf {
		end := i + maxLeaf
		if end > len(entries) {
			end = len(entries)
		}
		// Copy the slice so the leaf owns its own backing array
		batch := make([]*Entry, end-i)
		copy(batch, entries[i:end])
		leaf := &TreeNode{
			isLeaf:      true,
			keys:        make([]string, 0),
			leafEntries: batch,
			childNodes:  make([]*TreeNode, 0),
			order:       order,
		}
		leaves = append(leaves, leaf)
	}
 
	// --- Step 2: wire leaf linked list ---
	for i := 0; i < len(leaves)-1; i++ {
		leaves[i].next = leaves[i+1]
		leaves[i+1].prev = leaves[i]
	}
 
	// --- Step 3: build internal levels bottom-up ---
	// If there is only one leaf the loop is skipped and that leaf becomes root.
	currentLevel := make([]*TreeNode, len(leaves))
	copy(currentLevel, leaves)
 
	maxInternal := order * 2
	for len(currentLevel) > 1 {
		var nextLevel []*TreeNode
 
		for i := 0; i < len(currentLevel); i += maxInternal + 1 {
			end := i + maxInternal + 1
			if end > len(currentLevel) {
				end = len(currentLevel)
			}
			children := currentLevel[i:end]
 
			// Separator keys: first key of each child except the leftmost
			keys := make([]string, len(children)-1)
			for j := 1; j < len(children); j++ {
				// Walk down to the leftmost leaf of this child to get its smallest key
				node := children[j]
				for !node.isLeaf {
					node = node.childNodes[0]
				}
				keys[j-1] = node.leafEntries[0].key
			}
 
			internal := &TreeNode{
				isLeaf:     false,
				keys:       keys,
				childNodes: append([]*TreeNode{}, children...),
				leafEntries: make([]*Entry, 0),
				order:      order,
			}
			for _, child := range children {
				child.parent = internal
			}
			nextLevel = append(nextLevel, internal)
		}
		currentLevel = nextLevel
	}
 
	root := currentLevel[0]
    printTree(root)
	return &BTree{root: root, min: order}
}

func (tree *BTree) searchEntry(searchKey string) (*Entry, bool) {
	if tree.root == nil {
		return nil, false
	}
    //var key string = entry.key
	var current *TreeNode = tree.root

	for !current.isLeaf { // search the keys first	
        i := 0
        for i < len(current.keys) && searchKey >= current.keys[i] {
            i++
        }
        current = current.childNodes[i]
    }

    for _, entry := range current.leafEntries { //search leaf leafEntries
        if entry.key == searchKey {
            return entry, true
        }
    }
    return nil, false
}

func conditionSatisfied(key string, conditions map[string]string) bool {
	for limitKey, comparator := range conditions {
		switch comparator {
		case "<=":
			if !(key <= limitKey) {
				return false
			}
		case ">=":
			if !(key >= limitKey) {
				return false
			}
		case "<":
			if !(key < limitKey) {
				return false
			}
		case ">":
			if !(key > limitKey) {
				return false
			}
		case "==":
			if !(key == limitKey) {
				return false
			}
		default:
			return false // unknown comparator — fail safe
		}
	}
	return true
}
 
func (tree *BTree) searchRangeEntry(conditions map[string]string) []*Entry {
	if tree.root == nil {
		return nil
    }
	// Extract the tightest lower bound from conditions so we can start
	// the leaf scan as late as possible rather than from the very first leaf
	startKey := ""
	for limitKey, comparator := range conditions {
		if comparator == ">=" || comparator == ">" {
			if limitKey > startKey {
				startKey = limitKey
			}
		}
	}
 
	// Traverse down to the leaf that would contain startKey
	current := tree.root
	for !current.isLeaf {
		i := 0
		for i < len(current.keys) && startKey >= current.keys[i] {
			i++
		}
		current = current.childNodes[i]
	}
 
	// Walk the linked list, collecting entries that pass all conditions
	var results []*Entry
	for node := current; node != nil; node = node.next {
		for _, entry := range node.leafEntries {
			if conditionSatisfied(entry.key, conditions) {
				results = append(results, entry)
			}
		}
	}
	return results
}

func (tree *BTree) insert_element(entry Entry) bool {
    //fmt.Println("-------PRE INSERTION-------")
	var key string = entry.key
	var current *TreeNode = tree.root
	var parentStack []*TreeNode

	// Traverse down to the correct leaf, tracking parents
	for !current.isLeaf {
		i := 0
		parentStack = append(parentStack, current)
		for i < len(current.keys) && key >= current.keys[i] {
			i++
		}
		current = current.childNodes[i]
	}

	// Insert into leaf
	current.insert_leaf_entry(&entry)

	// No overflow — done
	if len(current.leafEntries) <= current.order*2 {
		return true
	}

	// Handle leaf overflow: split and propagate up
	separatorKey, rightNode := current.handle_leaf_overflow()

	// Walk up the parent stack, inserting separator keys and splitting as needed
	for len(parentStack) > 0 {
		parent := parentStack[len(parentStack)-1]
		parentStack = parentStack[:len(parentStack)-1]

		parent.insert_internal_key(separatorKey, rightNode)

		// No overflow at this internal node — done
		if len(parent.keys) <= parent.order*2 {
			return true
		}

		// Internal node overflows — split it and continue up
		separatorKey, rightNode = parent.handle_internal_overflow()
	}

	// If we reach here, the root itself overflowed — create a new root
	newRoot := create_tree_node(
		nil,
		[]*TreeNode{tree.root, rightNode},
		[]string{separatorKey},
		false,
		tree.root.order,
	)
	tree.root.parent = newRoot
	rightNode.parent = newRoot
	tree.root = newRoot
    
    //tree.printTree()
	return true
}

func (tree *BTree) delete(key string) bool {
    if tree.root == nil {
		return false
	}
	var current *TreeNode = tree.root

	for !current.isLeaf { // search the keys first	
        i := 0
        for i < len(current.keys) && key >= current.keys[i] {
            i++
        }
        current = current.childNodes[i]
    }

    return current.delete_leaf_entry(key)
}

/* ------------------------------------------------------------------ */
/* Public API                                                           */
/* ------------------------------------------------------------------ */
 
func CreateNewTree(order int) *BTree                                          { return create_new_tree(order) }
func CreateNewTreeBulkload(entries []*Entry, order int) *BTree                { return create_new_tree_bulkload(entries, order) }
func PrintTree(root *TreeNode)                                                { printTree(root) }
 
func (tree *BTree) Insert(entry *Entry) bool                                  { return tree.insert_element(*entry) }
func (tree *BTree) Search(key string) (*Entry, bool)                          { return tree.searchEntry(key) }
func (tree *BTree) Delete(key string) bool                                    { return tree.delete(key) }
func (tree *BTree) SearchRange(conditions map[string]string) []*Entry         { return tree.searchRangeEntry(conditions) }
func (tree *BTree) Root() *TreeNode                                           { return tree.root }