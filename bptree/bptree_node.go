package bptree

import (
	"slices"
	"strings"
)


type TreeNode struct {
	isLeaf bool
	keys []string // separator keys 
	leafEntries []*Entry //for leaf nodes only 
	childNodes []*TreeNode //for internal nodes
	order int
	next *TreeNode //leaf level linked list
	prev *TreeNode //leaf level linked list
	parent *TreeNode //parent of the current node
}

func create_tree_node(leafEntries []*Entry, childNodes []*TreeNode, keys []string, isLeaf bool, order int) *TreeNode {
	node := &TreeNode{
		isLeaf:     isLeaf,
		keys:       make([]string, 0),
		leafEntries:    make([]*Entry, 0),
		childNodes: make([]*TreeNode, 0),
		order: order,
		next:       nil,
		prev:       nil,
	}
	if isLeaf {
		node.leafEntries = leafEntries
	} else if len(keys) == 0 {
		node.childNodes = childNodes
		if len(childNodes) > 1 {
			for _, child := range childNodes[1:] {
				if len(child.leafEntries) > 0 {
					node.keys = append(node.keys, child.leafEntries[0].key)
				}
			}
		}
	} else {
		node.keys = keys
		node.childNodes = childNodes
	}
	return node
}

func create_empty_tree_node() *TreeNode {
	return &TreeNode{
		isLeaf:     true,
		keys:       make([]string, 0),
		leafEntries:    make([]*Entry, 0),
		childNodes: make([]*TreeNode, 0),
		next:       nil,
		prev:       nil,
		order: 0,
	}
}

/*LEAF NODE*/

func (treeNode *TreeNode) insert_leaf_entry(entry *Entry) bool {
	treeNode.leafEntries = append(treeNode.leafEntries, entry)
	slices.SortFunc(treeNode.leafEntries, func(a, b *Entry) int { //maintain sorted order 
		return strings.Compare(a.key, b.key)
	})
	return true
}  

func (treeNode *TreeNode) handle_leaf_overflow() (separatorKey string, node2 *TreeNode) {
	var n int = len(treeNode.leafEntries)/2

	var rightChildNode *TreeNode = create_tree_node(
		treeNode.leafEntries[n:],
		nil,
		nil,
		true, 
		treeNode.order,
	)

	treeNode.leafEntries = treeNode.leafEntries[:n]

	/* add linked list feature here */
	rightChildNode.next = treeNode.next     // right's next = left's old next
    rightChildNode.prev = treeNode          // right's prev = left
    if treeNode.next != nil {
        treeNode.next.prev = rightChildNode // old next's prev = right
    }
    treeNode.next = rightChildNode 
	
	rightChildNode.parent = treeNode.parent

	return rightChildNode.leafEntries[0].key, rightChildNode
}

func (treeNode *TreeNode) delete_leaf_entry(key string) bool {
    entries := treeNode.leafEntries
    for i, e := range entries {
        if e.key == key {
            treeNode.leafEntries = slices.Delete(entries, i, i+1)
			if len(treeNode.leafEntries) < treeNode.order {
				treeNode.handle_leaf_underflow()
			}
            return true
        }
    }
    return false
}

func (treeNode *TreeNode) handle_leaf_underflow() {
	parent := treeNode.parent
	if parent == nil {
		return // root leaf — nothing to do
	}
 
	// Find which child index this node is in the parent
	myIndex := 0
	for i, child := range parent.childNodes {
		if child == treeNode {
			myIndex = i
			break
		}
	}
 
	// Prefer right sibling, fall back to left
	rightSibling := treeNode.next
	if rightSibling != nil && rightSibling.parent == parent && len(rightSibling.leafEntries) > treeNode.order {
		// --- Redistribute from right sibling ---
		// Steal the leftmost entry from right sibling
		stolen := rightSibling.leafEntries[0]
		rightSibling.leafEntries = rightSibling.leafEntries[1:]
		treeNode.leafEntries = append(treeNode.leafEntries, stolen)
 
		// Update the separator key in parent: the key separating us from
		// the right sibling is now the new first entry of the right sibling
		parent.keys[myIndex] = rightSibling.leafEntries[0].key
		return
	}
 
	leftSibling := treeNode.prev
	if leftSibling != nil && leftSibling.parent == parent && len(leftSibling.leafEntries) > treeNode.order {
		// --- Redistribute from left sibling ---
		// Steal the rightmost entry from left sibling
		stolen := leftSibling.leafEntries[len(leftSibling.leafEntries)-1]
		leftSibling.leafEntries = leftSibling.leafEntries[:len(leftSibling.leafEntries)-1]
		treeNode.leafEntries = append([]*Entry{stolen}, treeNode.leafEntries...)
 
		// Update separator key in parent: key at myIndex-1 separates left from us
		parent.keys[myIndex-1] = treeNode.leafEntries[0].key
		return
	}
 
	// --- Merge ---
	// Prefer merging with right sibling (we absorb it), else merge into left
	if rightSibling != nil && rightSibling.parent == parent {
		// Absorb right sibling into treeNode
		treeNode.leafEntries = append(treeNode.leafEntries, rightSibling.leafEntries...)
 
		// Fix linked list: skip over right sibling
		treeNode.next = rightSibling.next
		if rightSibling.next != nil {
			rightSibling.next.prev = treeNode
		}
 
		// Remove the separator key and right sibling pointer from parent
		parent.keys = slices.Delete(parent.keys, myIndex, myIndex+1)
		parent.childNodes = slices.Delete(parent.childNodes, myIndex+1, myIndex+2)
 
	} else if leftSibling != nil && leftSibling.parent == parent {
		// Absorb treeNode into left sibling
		leftSibling.leafEntries = append(leftSibling.leafEntries, treeNode.leafEntries...)
 
		// Fix linked list
		leftSibling.next = treeNode.next
		if treeNode.next != nil {
			treeNode.next.prev = leftSibling
		}
 
		// Remove the separator key and treeNode pointer from parent
		parent.keys = slices.Delete(parent.keys, myIndex-1, myIndex)
		parent.childNodes = slices.Delete(parent.childNodes, myIndex, myIndex+1)
	}
 
	// Propagate underflow up to the parent internal node
	if len(parent.keys) < parent.order && parent.parent != nil {
		parent.handle_internal_underflow()
	}
}



/*INTERNAL NODES*/

func (treeNode *TreeNode) insert_internal_key(separatorKey string, newRightChild *TreeNode) bool {
	insertPos := len(treeNode.keys)
    for i, k := range treeNode.keys {
        if separatorKey < k {
            insertPos = i
            break
        }
    }

	treeNode.keys = append(treeNode.keys, "")
	copy(treeNode.keys[insertPos+1:], treeNode.keys[insertPos:])
    treeNode.keys[insertPos] = separatorKey

	childPos := insertPos + 1
	treeNode.childNodes = append(treeNode.childNodes, nil)
	copy(treeNode.childNodes[childPos+1:], treeNode.childNodes[childPos:])
	treeNode.childNodes[childPos] = newRightChild
    newRightChild.parent = treeNode

	return true
}

func (treeNode *TreeNode) handle_internal_overflow() (key string, node2 *TreeNode) {
	var n int = len(treeNode.keys)/2 //always the median of the array keys

	var medianKey string = treeNode.keys[n]

	var rightChildNode *TreeNode = create_tree_node(
		nil,
		treeNode.childNodes[n+1:],
		treeNode.keys[n+1:],
		false,
		treeNode.order,
	)

	treeNode.childNodes = treeNode.childNodes[:n+1]
	treeNode.keys = treeNode.keys[:n]

	return medianKey, rightChildNode 
}

//new challenge: handle underflow
func (treeNode *TreeNode) handle_internal_underflow() {
	parent := treeNode.parent
	if parent == nil {
		return // root — underflow here just means the tree shrinks, handled at tree level
	}
 
	// Find our index in parent
	myIndex := 0
	for i, child := range parent.childNodes {
		if child == treeNode {
			myIndex = i
			break
		}
	}
 
	// Try right sibling first
	if myIndex < len(parent.childNodes)-1 {
		rightSibling := parent.childNodes[myIndex+1]
		if len(rightSibling.keys) > treeNode.order {
			// --- Redistribute from right sibling ---
			// Pull separator key down from parent into our keys
			treeNode.keys = append(treeNode.keys, parent.keys[myIndex])
			// Adopt right sibling's leftmost child
			treeNode.childNodes = append(treeNode.childNodes, rightSibling.childNodes[0])
			rightSibling.childNodes[0].parent = treeNode
			// Promote right sibling's first key up to parent
			parent.keys[myIndex] = rightSibling.keys[0]
			// Remove the borrowed key and child from right sibling
			rightSibling.keys = rightSibling.keys[1:]
			rightSibling.childNodes = rightSibling.childNodes[1:]
			return
		}
	}
 
	// Try left sibling
	if myIndex > 0 {
		leftSibling := parent.childNodes[myIndex-1]
		if len(leftSibling.keys) > treeNode.order {
			// --- Redistribute from left sibling ---
			// Pull separator key down from parent, prepend to our keys
			treeNode.keys = append([]string{parent.keys[myIndex-1]}, treeNode.keys...)
			// Adopt left sibling's rightmost child
			treeNode.childNodes = append([]*TreeNode{leftSibling.childNodes[len(leftSibling.childNodes)-1]}, treeNode.childNodes...)
			leftSibling.childNodes[len(leftSibling.childNodes)-1].parent = treeNode
			// Promote left sibling's last key up to parent
			parent.keys[myIndex-1] = leftSibling.keys[len(leftSibling.keys)-1]
			// Remove borrowed key and child from left sibling
			leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
			leftSibling.childNodes = leftSibling.childNodes[:len(leftSibling.childNodes)-1]
			return
		}
	}
 
	// --- Merge ---
	if myIndex < len(parent.childNodes)-1 {
		// Merge with right sibling: absorb right into treeNode
		rightSibling := parent.childNodes[myIndex+1]
		// Pull separator key down from parent to join the two halves
		treeNode.keys = append(treeNode.keys, parent.keys[myIndex])
		treeNode.keys = append(treeNode.keys, rightSibling.keys...)
		treeNode.childNodes = append(treeNode.childNodes, rightSibling.childNodes...)
		for _, child := range rightSibling.childNodes {
			child.parent = treeNode
		}
		// Remove separator key and right sibling from parent
		parent.keys = slices.Delete(parent.keys, myIndex, myIndex+1)
		parent.childNodes = slices.Delete(parent.childNodes, myIndex+1, myIndex+2)
 
	} else {
		// Merge with left sibling: absorb treeNode into left
		leftSibling := parent.childNodes[myIndex-1]
		leftSibling.keys = append(leftSibling.keys, parent.keys[myIndex-1])
		leftSibling.keys = append(leftSibling.keys, treeNode.keys...)
		leftSibling.childNodes = append(leftSibling.childNodes, treeNode.childNodes...)
		for _, child := range treeNode.childNodes {
			child.parent = leftSibling
		}
		// Remove separator key and treeNode from parent
		parent.keys = slices.Delete(parent.keys, myIndex-1, myIndex)
		parent.childNodes = slices.Delete(parent.childNodes, myIndex, myIndex+1)
	}
 
	// Propagate up if parent is now under minimum
	if len(parent.keys) < parent.order && parent.parent != nil {
		parent.handle_internal_underflow()
	}
}
 