package bptree

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
)

/*

command to test everything

*/

/* Entry */
func Test_create_entry(t *testing.T) {
	var correctKey string = "key"
	var correctValue string = "value"
	var newEntry *Entry = create_entry("key", "value")
	if newEntry.key != "key" {
		t.Errorf("create_entry Want: %s, Got: %s", correctKey, newEntry.key)
	}
	if newEntry.value != "value" {
		t.Errorf("create_entry Want: %s, Got: %s", correctValue, newEntry.value)
	}
}

/* BPTree node  */
func Test_create_tree_node(t *testing.T) {

	t.Run("leaf node with entries", func(t *testing.T) {
		entries := []*Entry{
			{key: "a", value: "1"},
			{key: "b", value: "2"},
		}
		node := create_tree_node(entries, nil, nil, true, 3)

		if !node.isLeaf {
			t.Error("expected isLeaf=true")
		}
		if len(node.leafEntries) != 2 {
			t.Errorf("expected 2 leafEntries, got %d", len(node.leafEntries))
		}
		if len(node.childNodes) != 0 {
			t.Errorf("expected 0 childNodes, got %d", len(node.childNodes))
		}
		if len(node.keys) != 0 {
			t.Errorf("expected 0 keys, got %d", len(node.keys))
		}
		if node.next != nil || node.prev != nil {
			t.Error("expected next and prev to be nil")
		}
		if node.order != 3 {
			t.Errorf("expected order=3, got %d", node.order)
		}
	})

	t.Run("leaf node with no entries", func(t *testing.T) {
		node := create_tree_node([]*Entry{}, nil, nil, true, 4)

		if !node.isLeaf {
			t.Error("expected isLeaf=true")
		}
		if len(node.leafEntries) != 0 {
			t.Errorf("expected 0 leafEntries, got %d", len(node.leafEntries))
		}
	})

	t.Run("lowest internal node derives separator keys from children", func(t *testing.T) {
		child1 := &TreeNode{leafEntries: []*Entry{{key: "a"}}}
		child2 := &TreeNode{leafEntries: []*Entry{{key: "d"}}}
		child3 := &TreeNode{leafEntries: []*Entry{{key: "g"}}}

		node := create_tree_node(nil, []*TreeNode{child1, child2, child3}, nil, false, 3)

		if node.isLeaf {
			t.Error("expected isLeaf=false")
		}
		if len(node.childNodes) != 3 {
			t.Errorf("expected 3 childNodes, got %d", len(node.childNodes))
		}
		// separator keys come from child[1:] → ["d", "g"]
		if len(node.keys) != 2 {
			t.Errorf("expected 2 separator keys, got %d", len(node.keys))
		}
		if node.keys[0] != "d" || node.keys[1] != "g" {
			t.Errorf("expected keys [d g], got %v", node.keys)
		}
		if len(node.leafEntries) != 0 {
			t.Errorf("expected 0 leafEntries on internal node, got %d", len(node.leafEntries))
		}
	})

	t.Run("internal node with single child has no separator keys", func(t *testing.T) {
		child := &TreeNode{leafEntries: []*Entry{{key: "x"}}}
		node := create_tree_node(nil, []*TreeNode{child}, nil, false, 3)

		if len(node.keys) != 0 {
			t.Errorf("expected 0 keys for single child, got %d", len(node.keys))
		}
		if len(node.childNodes) != 1 {
			t.Errorf("expected 1 childNode, got %d", len(node.childNodes))
		}
	})

	t.Run("internal node skips children with no leafEntries for key derivation", func(t *testing.T) {
		child1 := &TreeNode{leafEntries: []*Entry{{key: "a"}}}
		child2 := &TreeNode{leafEntries: []*Entry{}} // empty — contributes no key
		child3 := &TreeNode{leafEntries: []*Entry{{key: "z"}}}

		node := create_tree_node(nil, []*TreeNode{child1, child2, child3}, nil, false, 3)

		// only child3 contributes → ["z"]
		if len(node.keys) != 1 {
			t.Errorf("expected 1 key, got %d: %v", len(node.keys), node.keys)
		}
		if node.keys[0] != "z" {
			t.Errorf("expected key 'z', got %q", node.keys[0])
		}
	})

	t.Run("internal node with no children", func(t *testing.T) {
		node := create_tree_node(nil, []*TreeNode{}, nil, false, 3)

		if len(node.keys) != 0 {
			t.Errorf("expected 0 keys, got %d", len(node.keys))
		}
		if len(node.childNodes) != 0 {
			t.Errorf("expected 0 childNodes, got %d", len(node.childNodes))
		}
	})
}

func Test_insert_leaf_entry(t *testing.T) {
	t.Run("insert entry into empty leaf node", func(t *testing.T) {
		node := create_empty_tree_node()
		entry := create_entry("key", "value")
		node.insert_leaf_entry(entry)
		if len(node.leafEntries) != 1 {
			t.Errorf("Expected 1 entry, got %d", len(node.leafEntries))
		}
		if len(node.keys) != 0 {
			t.Errorf("Expected 0 keys, got %d", len(node.keys))
		}
	}) 

	t.Run("insert entry into non-empty leaf node ", func(t *testing.T) {
		entries := []*Entry{
			{key: "a", value: "1"},
			{key: "b", value: "2"},
		}
		node := create_tree_node(entries, nil, nil, true, 3)

		entry := create_entry("c", "3")

		node.insert_leaf_entry(entry)
		if len(node.leafEntries) != 3 {
			t.Errorf("Expected 3 entry, got %d", len(node.leafEntries))
		}
		if len(node.keys) != 0 {
			t.Errorf("Expected 0 keys, got %d", len(node.keys))
		} 
		for i, e := range entries {
			if i > 0 && e.key < entries[i-1].key {
				t.Errorf("Expected sorted keys, got keys in the wrong order")
			}
		}
 	}) 
}

func Test_delete_leaf_entry(t *testing.T) {
    makeEntries := func() []*Entry {
        return []*Entry{
            {key: "a", value: "1"},
            {key: "b", value: "2"},
            {key: "c", value: "3"},
        }
    }

    tests := []struct {
        name          string
        deleteKey     string
        wantFound     bool
        wantLen       int
        wantRemaining []string
    }{
        {
            name:          "delete first entry",
            deleteKey:     "a",
            wantFound:     true,
            wantLen:       2,
            wantRemaining: []string{"b", "c"},
        },
        {
            name:          "delete middle entry",
            deleteKey:     "b",
            wantFound:     true,
            wantLen:       2,
            wantRemaining: []string{"a", "c"},
        },
        {
            name:          "delete last entry",
            deleteKey:     "c",
            wantFound:     true,
            wantLen:       2,
            wantRemaining: []string{"a", "b"},
        },
        {
            name:          "delete non-existent key",
            deleteKey:     "z",
            wantFound:     false,
            wantLen:       3,
            wantRemaining: []string{"a", "b", "c"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            node := create_tree_node(makeEntries(), nil, nil, true, 3)

            got := node.delete_leaf_entry(tt.deleteKey)

            if got != tt.wantFound {
                t.Errorf("delete_leaf_entry() = %v, want %v", got, tt.wantFound)
            }

            if len(node.leafEntries) != tt.wantLen {
                t.Errorf("len(leafEntries) = %d, want %d", len(node.leafEntries), tt.wantLen)
            }

            for _, e := range node.leafEntries {
                if e.key == tt.deleteKey && tt.wantFound {
                    t.Errorf("key %q still present after deletion", tt.deleteKey)
                }
            }

			remaining := make([]string, len(node.leafEntries))
            for i, e := range node.leafEntries {
                remaining[i] = e.key
            }
            if !slices.Equal(remaining, tt.wantRemaining) {
                t.Errorf("remaining keys = %v, want %v", remaining, tt.wantRemaining)
            }
        })
    }
}

func Test_handle_leaf_overflow(t *testing.T) {
    // helper to create leaf entries
    makeEntries := func(keys ...string) []*Entry {
        entries := make([]*Entry, len(keys))
        for i, k := range keys {
            entries[i] = &Entry{key: k}
        }
        return entries
    }

    t.Run("separator key is first key of right node", func(t *testing.T) {
        node := create_tree_node(makeEntries("a", "b", "c", "d", "e"), nil, nil, true, 2)

        separatorKey, _ := node.handle_leaf_overflow()

        if separatorKey != "c" {
            t.Errorf("expected separator key 'c', got '%s'", separatorKey)
        }
    })

    t.Run("left node gets first half of entries", func(t *testing.T) {
        node := create_tree_node(makeEntries("a", "b", "c", "d", "e"), nil, nil, true, 2)

        node.handle_leaf_overflow()

        if len(node.leafEntries) != 2 {
            t.Errorf("expected left node to have 2 entries, got %d", len(node.leafEntries))
        }
        if node.leafEntries[0].key != "a" || node.leafEntries[1].key != "b" {
            t.Errorf("expected left node keys [a, b], got %v", node.leafEntries)
        }
    })

    t.Run("right node gets second half of entries", func(t *testing.T) {
        node := create_tree_node(makeEntries("a", "b", "c", "d", "e"), nil, nil, true, 2)

        _, rightNode := node.handle_leaf_overflow()

        if len(rightNode.leafEntries) != 3 {
            t.Errorf("expected right node to have 2 entries, got %d", len(rightNode.leafEntries))
        }
        if rightNode.leafEntries[0].key != "c" || rightNode.leafEntries[1].key != "d" {
            t.Errorf("expected right node keys [c, d], got %v", rightNode.leafEntries)
        }
    })

    t.Run("linked list - right.prev points to left", func(t *testing.T) {
        node := create_tree_node(makeEntries("a", "b", "c", "d", "e"), nil, nil, true, 2)

        _, rightNode := node.handle_leaf_overflow()

        if rightNode.prev != node {
            t.Errorf("expected rightNode.prev to point to left node")
        }
    })

    t.Run("linked list - left.next points to right", func(t *testing.T) {
        node := create_tree_node(makeEntries("a", "b", "c", "d", "e"), nil, nil, true, 2)

        _, rightNode := node.handle_leaf_overflow()

        if node.next != rightNode {
            t.Errorf("expected left node.next to point to right node")
        }
    })

    t.Run("linked list - right.next points to old next", func(t *testing.T) {
        oldNext := create_tree_node(makeEntries("f", "g"), nil, nil, true, 2)
        node := create_tree_node(makeEntries("a", "b", "c", "d", "e"), nil, nil, true, 2)
        node.next = oldNext
        oldNext.prev = node

        _, rightNode := node.handle_leaf_overflow()

        if rightNode.next != oldNext {
            t.Errorf("expected rightNode.next to point to oldNext")
        }
    })

    t.Run("linked list - oldNext.prev points to right node", func(t *testing.T) {
        oldNext := create_tree_node(makeEntries("f", "g"), nil, nil, true, 2)
        node := create_tree_node(makeEntries("a", "b", "c", "d", "e"), nil, nil, true, 2)
        node.next = oldNext
        oldNext.prev = node

        _, rightNode := node.handle_leaf_overflow()

        if oldNext.prev != rightNode {
            t.Errorf("expected oldNext.prev to point to rightNode")
        }
    })

    t.Run("right node parent should be same as left node parent - BUG", func(t *testing.T) {
        parent := create_tree_node(nil, nil, []string{"e"}, false, 2)
        node := create_tree_node(makeEntries("a", "b", "c", "d", "e"), nil, nil, true, 2)
        node.parent = parent

        _, rightNode := node.handle_leaf_overflow()

        // this will FAIL - rightChildNode.parent = treeNode instead of treeNode.parent
        if rightNode.parent != parent {
            t.Errorf("expected rightNode.parent to be parent node, got %v", rightNode.parent)
        }
    })

    t.Run("right node is a leaf", func(t *testing.T) {
        node := create_tree_node(makeEntries("a", "b", "c", "d", "e"), nil, nil, true, 2)

        _, rightNode := node.handle_leaf_overflow()

        if !rightNode.isLeaf {
            t.Errorf("expected right node to be a leaf node")
        }
    })

    t.Run("odd number of entries", func(t *testing.T) {
        node := create_tree_node(makeEntries("a", "b", "c", "d", "e"), nil, nil, true, 2)

        _, rightNode := node.handle_leaf_overflow()

        if len(node.leafEntries) != 2 {
            t.Errorf("expected left node to have 2 entries, got %d", len(node.leafEntries))
        }
        if len(rightNode.leafEntries) != 3 {
            t.Errorf("expected right node to have 3 entries, got %d", len(rightNode.leafEntries))
        }
    })
}

func Test_insert_internal_key(t *testing.T) {
    makeLeafNode := func(key string) *TreeNode {
        return &TreeNode{leafEntries: []*Entry{{key: key}}}
    }

    t.Run("insert separator key at beginning", func(t *testing.T) {
        c0 := &TreeNode{leafEntries: []*Entry{{key: "a"}}}
        c1 := &TreeNode{leafEntries: []*Entry{{key: "c"}}}
        node := create_tree_node(nil, []*TreeNode{c0, c1}, []string{"c"}, false, 2)
        newRight := makeLeafNode("b")

        node.insert_internal_key("b", newRight)

        if len(node.keys) != 2 {
            t.Errorf("expected 2 keys, got %d", len(node.keys))
        }
        if node.keys[0] != "b" || node.keys[1] != "c" {
            t.Errorf("expected keys [b, c], got %v", node.keys)
        }
    })

    t.Run("insert separator key at end", func(t *testing.T) {
        c0 := &TreeNode{leafEntries: []*Entry{{key: "a"}}}
        c1 := &TreeNode{leafEntries: []*Entry{{key: "c"}}}
        node := create_tree_node(nil, []*TreeNode{c0, c1}, []string{"c"}, false, 2)
        newRight := makeLeafNode("e")

        node.insert_internal_key("e", newRight)

        if len(node.keys) != 2 {
            t.Errorf("expected 2 keys, got %d", len(node.keys))
        }
        if node.keys[0] != "c" || node.keys[1] != "e" {
            t.Errorf("expected keys [c, e], got %v", node.keys)
        }
    })

    t.Run("insert separator key in middle", func(t *testing.T) {
        c0 := &TreeNode{leafEntries: []*Entry{{key: "a"}}}
        c1 := &TreeNode{leafEntries: []*Entry{{key: "c"}}}
        c2 := &TreeNode{leafEntries: []*Entry{{key: "e"}}}
        node := create_tree_node(nil, []*TreeNode{c0, c1, c2}, []string{"c", "e"}, false, 2)
        newRight := makeLeafNode("d")

        node.insert_internal_key("d", newRight)

        if len(node.keys) != 3 {
            t.Errorf("expected 3 keys, got %d", len(node.keys))
        }
        if node.keys[0] != "c" || node.keys[1] != "d" || node.keys[2] != "e" {
            t.Errorf("expected keys [c, d, e], got %v", node.keys)
        }
    })

    t.Run("new right child inserted at correct child position", func(t *testing.T) {
        c0 := &TreeNode{leafEntries: []*Entry{{key: "a"}}}
        c1 := &TreeNode{leafEntries: []*Entry{{key: "c"}}}
        c2 := &TreeNode{leafEntries: []*Entry{{key: "e"}}}
        node := create_tree_node(nil, []*TreeNode{c0, c1, c2}, []string{"c", "e"}, false, 2)
        newRight := makeLeafNode("d")

        node.insert_internal_key("d", newRight)

        if len(node.childNodes) != 4 {
            t.Errorf("expected 4 children, got %d", len(node.childNodes))
        }
        if node.childNodes[2] != newRight {
            t.Errorf("expected newRight at childPos 2")
        }
    })

    t.Run("n keys always has n+1 children", func(t *testing.T) {
        c0 := &TreeNode{leafEntries: []*Entry{{key: "a"}}}
        c1 := &TreeNode{leafEntries: []*Entry{{key: "c"}}}
        node := create_tree_node(nil, []*TreeNode{c0, c1}, []string{"c"}, false, 2)
        newRight := makeLeafNode("b")

        node.insert_internal_key("b", newRight)

        if len(node.childNodes) != len(node.keys)+1 {
            t.Errorf("expected %d children, got %d", len(node.keys)+1, len(node.childNodes))
        }
    })

    t.Run("new right child parent pointer updated", func(t *testing.T) {
        c0 := &TreeNode{leafEntries: []*Entry{{key: "a"}}}
        c1 := &TreeNode{leafEntries: []*Entry{{key: "c"}}}
        node := create_tree_node(nil, []*TreeNode{c0, c1}, []string{"c"}, false, 2)
        newRight := makeLeafNode("e")

        node.insert_internal_key("e", newRight)

        if newRight.parent != node {
            t.Errorf("expected newRight.parent to point to internal node")
        }
    })

    t.Run("existing children are not displaced", func(t *testing.T) {
        c0 := &TreeNode{leafEntries: []*Entry{{key: "a"}}}
        c1 := &TreeNode{leafEntries: []*Entry{{key: "c"}}}
        c2 := &TreeNode{leafEntries: []*Entry{{key: "e"}}}
        node := create_tree_node(nil, []*TreeNode{c0, c1, c2}, []string{"c", "e"}, false, 2)
        newRight := makeLeafNode("d")

        node.insert_internal_key("d", newRight)

        if node.childNodes[0] != c0 {
            t.Errorf("expected c0 at position 0")
        }
        if node.childNodes[1] != c1 {
            t.Errorf("expected c1 at position 1")
        }
        if node.childNodes[3] != c2 {
            t.Errorf("expected c2 at position 3")
        }
    })

    t.Run("returns true", func(t *testing.T) {
        c0 := &TreeNode{leafEntries: []*Entry{{key: "a"}}}
        c1 := &TreeNode{leafEntries: []*Entry{{key: "c"}}}
        node := create_tree_node(nil, []*TreeNode{c0, c1}, []string{"c"}, false, 2)
        newRight := makeLeafNode("e")

        result := node.insert_internal_key("e", newRight)

        if !result {
            t.Errorf("expected insert_internal_key to return true")
        }
    })
}

func Test_handle_internal_overflow(t *testing.T) {
	c0 := &TreeNode{leafEntries: []*Entry{{key: "a"}}}
	c1 := &TreeNode{leafEntries: []*Entry{{key: "b"}}}
	c2 := &TreeNode{leafEntries: []*Entry{{key: "c"}}}
	c3 := &TreeNode{leafEntries: []*Entry{{key: "d"}}}
	c4 := &TreeNode{leafEntries: []*Entry{{key: "e"}}}
	c5 := &TreeNode{leafEntries: []*Entry{{key: "f"}}}
    
	t.Run("median key is returned", func(t *testing.T) {
		node := create_tree_node(nil, []*TreeNode{c0, c1, c2, c3, c4, c5}, []string{"a", "b", "c", "d", "e"}, false, 2)
        medianKey, _ := node.handle_internal_overflow()

        if medianKey != "c" {
            t.Errorf("expected median key 'c', got '%s'", medianKey)
        }
    })

    t.Run("left node keeps keys before median", func(t *testing.T) {
		node := create_tree_node(nil, []*TreeNode{c0, c1, c2, c3, c4, c5}, []string{"a", "b", "c", "d", "e"}, false, 2)
        node.handle_internal_overflow()
        if len(node.keys) != 2 {
            t.Errorf("expected left node to have 2 keys, got %d", len(node.keys))
        }
        if node.keys[0] != "a" || node.keys[1] != "b" {
            t.Errorf("expected left node keys [a, b], got %v", node.keys)
        }
    })

    t.Run("right node gets keys after median", func(t *testing.T) {
		node := create_tree_node(nil, []*TreeNode{c0, c1, c2, c3, c4, c5}, []string{"a", "b", "c", "d", "e"}, false, 2)
        _, rightNode := node.handle_internal_overflow()

        if len(rightNode.keys) != 2 {
            t.Errorf("expected right node to have 2 keys, got %d", len(rightNode.keys))
        }
        if rightNode.keys[0] != "d" || rightNode.keys[1] != "e" {
            t.Errorf("expected right node keys [d, e], got %v", rightNode.keys)
        }
    })

    t.Run("median key is not in left node", func(t *testing.T) {
		node := create_tree_node(nil, []*TreeNode{c0, c1, c2, c3, c4, c5}, []string{"a", "b", "c", "d", "e"}, false, 2)
        medianKey, _ := node.handle_internal_overflow()

        for _, k := range node.keys {
            if k == medianKey {
                t.Errorf("median key '%s' should not remain in left node", medianKey)
            }
        }
    })

    t.Run("median key is not in right node", func(t *testing.T) {
		node := create_tree_node(nil, []*TreeNode{c0, c1, c2, c3, c4, c5}, []string{"a", "b", "c", "d", "e"}, false, 2)
        medianKey, rightNode := node.handle_internal_overflow()

        for _, k := range rightNode.keys {
            if k == medianKey {
                t.Errorf("median key '%s' should not remain in right node", medianKey)
            }
        }
    })

    t.Run("left node keeps children up to and including n", func(t *testing.T) {
		node := create_tree_node(nil, []*TreeNode{c0, c1, c2, c3, c4, c5}, []string{"a", "b", "c", "d", "e"}, false, 2)
        node.handle_internal_overflow()

        // n = 5/2 = 2, left keeps childNodes[:n+1] = [:3] = [c0, c1, c2]
        if len(node.childNodes) != 3 {
            t.Errorf("expected left node to have 3 children, got %d", len(node.childNodes))
        }
        if node.childNodes[0] != c0 || node.childNodes[1] != c1 || node.childNodes[2] != c2 {
            t.Errorf("expected left children [c0, c1, c2], got %v", node.childNodes)
        }
    })

    t.Run("right node gets children after n", func(t *testing.T) {
		node := create_tree_node(nil, []*TreeNode{c0, c1, c2, c3, c4, c5}, []string{"a", "b", "c", "d", "e"}, false, 2)
        _, rightNode := node.handle_internal_overflow()

        // right gets childNodes[n+1:] = [3:] = [c3, c4, c5]
        if len(rightNode.childNodes) != 3 {
            t.Errorf("expected right node to have 3 children, got %d", len(rightNode.childNodes))
        }
        if rightNode.childNodes[0] != c3 || rightNode.childNodes[1] != c4 || rightNode.childNodes[2] != c5 {
            t.Errorf("expected right children [c3, c4, c5], got %v", rightNode.childNodes)
        }
    })

    t.Run("right node is not a leaf", func(t *testing.T) {
		node := create_tree_node(nil, []*TreeNode{c0, c1, c2, c3, c4, c5}, []string{"a", "b", "c", "d", "e"}, false, 2)
        _, rightNode := node.handle_internal_overflow()

        if rightNode.isLeaf {
            t.Errorf("expected right node to be an internal node")
        }
    })

    t.Run("n keys has n+1 children on both sides after split", func(t *testing.T) {
		node := create_tree_node(nil, []*TreeNode{c0, c1, c2, c3, c4, c5}, []string{"a", "b", "c", "d", "e"}, false, 2)
        _, rightNode := node.handle_internal_overflow()

        if len(node.childNodes) != len(node.keys)+1 {
            t.Errorf("left: expected %d children for %d keys, got %d", len(node.keys)+1, len(node.keys), len(node.childNodes))
        }
        if len(rightNode.childNodes) != len(rightNode.keys)+1 {
            t.Errorf("right: expected %d children for %d keys, got %d", len(rightNode.keys)+1, len(rightNode.keys), len(rightNode.childNodes))
        }
    })
}


/* BPTree */

func Test_create_new_tree(t *testing.T) {
	t.Run("returns non-nil tree", func(t *testing.T) {
		tree := create_new_tree(2)
		if tree == nil {
			t.Fatal("expected non-nil tree")
		}
	})

	t.Run("root is a leaf node", func(t *testing.T) {
		tree := create_new_tree(2)
		if !tree.root.isLeaf {
			t.Error("expected root to be a leaf")
		}
	})

	t.Run("root starts empty", func(t *testing.T) {
		tree := create_new_tree(2)
		if len(tree.root.leafEntries) != 0 {
			t.Errorf("expected 0 entries, got %d", len(tree.root.leafEntries))
		}
	})

	t.Run("root order is set correctly", func(t *testing.T) {
		tree := create_new_tree(3)
		if tree.root.order != 3 {
			t.Errorf("expected order 3, got %d", tree.root.order)
		}
	})

	t.Run("min reflects order", func(t *testing.T) {
		tree := create_new_tree(4)
		if tree.min != 4 {
			t.Errorf("expected min 4, got %d", tree.min)
		}
	})

	t.Run("panics on order less than 2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for order < 2")
			}
		}()
		create_new_tree(1)
	})

	t.Run("search on empty tree returns not found", func(t *testing.T) {
		tree := create_new_tree(2)
		_, found := tree.searchEntry("x")
		if found {
			t.Error("expected not found on empty tree")
		}
	})
}

func Test_insert_and_search(t *testing.T) {
	t.Run("insert single entry and find it", func(t *testing.T) {
		tree := create_new_tree(2)
		tree.insert_element(Entry{key: "a", value: "1"})
		e, found := tree.searchEntry("a")
		if !found {
			t.Fatal("expected to find inserted entry")
		}
		if e.value != "1" {
			t.Errorf("expected value '1', got '%s'", e.value)
		}
	})

	t.Run("insert many entries and find each", func(t *testing.T) {
		tree := create_new_tree(2)
		keys := []string{"d", "b", "f", "a", "c", "e", "g"}
		for _, k := range keys {
			tree.insert_element(Entry{key: k, value: k + "_val"})
		}
		for _, k := range keys {
			e, found := tree.searchEntry(k)
			if !found {
				t.Errorf("expected to find key '%s'", k)
			} else if e.value != k+"_val" {
				t.Errorf("key '%s': expected value '%s_val', got '%s'", k, k, e.value)
			}
		}
	})

	t.Run("search for non-existent key returns false", func(t *testing.T) {
		tree := create_new_tree(2)
		tree.insert_element(Entry{key: "a", value: "1"})
		_, found := tree.searchEntry("z")
		if found {
			t.Error("expected not found for missing key")
		}
	})

	t.Run("insert causes root split", func(t *testing.T) {
		tree := create_new_tree(2)
		for _, k := range []string{"a", "b", "c", "d", "e"} {
			tree.insert_element(Entry{key: k, value: k})
		}
		if tree.root.isLeaf {
			t.Error("expected root to have split into internal node")
		}
		for _, k := range []string{"a", "b", "c", "d", "e"} {
			_, found := tree.searchEntry(k)
			if !found {
				t.Errorf("key '%s' missing after root split", k)
			}
		}
	})
}

func Test_printTree(t *testing.T) {
	// Helper: capture stdout from printTree
	captureOutput := func(tree *BTree) string {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		printTree(tree.root)
		w.Close()
		os.Stdout = old
		var buf bytes.Buffer
		buf.ReadFrom(r)
		return strings.TrimRight(buf.String(), "\n")
	}

	t.Run("empty tree", func(t *testing.T) {
		tree := create_new_tree(2)
		out := captureOutput(tree)
		if out != "(empty tree)" {
			t.Errorf("expected '(empty tree)', got %q", out)
		}
	})

	t.Run("single entry — root is leaf", func(t *testing.T) {
		tree := create_new_tree(2)
		tree.insert_element(Entry{key: "a", value: "1"})
		out := captureOutput(tree)
		lines := strings.Split(out, "\n")
		if lines[0] != "[a]" {
			t.Errorf("expected '[a]', got %q", lines[0])
		}
		if lines[1] != "leaves: a" {
			t.Errorf("expected 'leaves: a', got %q", lines[1])
		}
	})

	t.Run("leaves line is sorted", func(t *testing.T) {
		tree := create_new_tree(2)
		for _, k := range []string{"d", "b", "f", "a", "c", "e", "g"} {
			tree.insert_element(Entry{key: k, value: k})
		}
		out := captureOutput(tree)
		lines := strings.Split(out, "\n")
		last := lines[len(lines)-1]
		if last != "leaves: a b c d e f g" {
			t.Errorf("expected sorted leaves, got %q", last)
		}
	})

	t.Run("internal node shows separator keys", func(t *testing.T) {
		tree := create_new_tree(2)
		for _, k := range []string{"a", "b", "c", "d", "e"} {
			tree.insert_element(Entry{key: k, value: k})
		}
		out := captureOutput(tree)
		lines := strings.Split(out, "\n")
		// Root must be an internal node after splitting
		if strings.HasPrefix(lines[0], "leaves:") {
			t.Error("expected an internal node at root after split")
		}
	})
}

func Test_printTree_visual(t *testing.T) {
	tree := create_new_tree(2)
	for _, k := range []string{"d", "b", "f", "a", "c", "e", "g"} {
		tree.insert_element(Entry{key: k, value: k})
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	printTree(tree.root)
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	t.Log("\n" + buf.String())
}

func Test_BTree_delete(t *testing.T) {
	setup := func(keys []string) *BTree {
		tree := create_new_tree(2)
		for _, k := range keys {
			tree.insert_element(Entry{key: k, value: k + "_val"})
		}
		return tree
	}

	t.Run("delete from empty tree returns false", func(t *testing.T) {
		tree := create_new_tree(2)
		if tree.delete("a") {
			t.Error("expected false deleting from empty tree")
		}
	})

	t.Run("delete non-existent key returns false", func(t *testing.T) {
		tree := setup([]string{"a", "b", "c"})
		if tree.delete("z") {
			t.Error("expected false for missing key")
		}
	})

	t.Run("delete existing key returns true", func(t *testing.T) {
		tree := setup([]string{"a", "b", "c"})
		if !tree.delete("b") {
			t.Error("expected true for existing key")
		}
	})

	t.Run("deleted key is no longer searchable", func(t *testing.T) {
		tree := setup([]string{"a", "b", "c"})
		tree.delete("b")
		_, found := tree.searchEntry("b")
		if found {
			t.Error("expected key 'b' to be gone after delete")
		}
	})

	t.Run("other keys survive delete", func(t *testing.T) {
		tree := setup([]string{"a", "b", "c"})
		tree.delete("b")
		for _, k := range []string{"a", "c"} {
			_, found := tree.searchEntry(k)
			if !found {
				t.Errorf("expected key '%s' to still exist after deleting 'b'", k)
			}
		}
	})

	t.Run("delete only entry in tree", func(t *testing.T) {
		tree := setup([]string{"a"})
		if !tree.delete("a") {
			t.Error("expected true deleting sole entry")
		}
		_, found := tree.searchEntry("a")
		if found {
			t.Error("expected tree to be empty after deleting sole entry")
		}
	})

	t.Run("delete all entries one by one", func(t *testing.T) {
		keys := []string{"a", "b", "c", "d", "e", "f", "g"}
		tree := setup(keys)
		for _, k := range keys {
			if !tree.delete(k) {
				t.Errorf("expected true deleting '%s'", k)
			}
			_, found := tree.searchEntry(k)
			if found {
				t.Errorf("key '%s' still present after delete", k)
			}
		}
	})

	t.Run("delete across a split — leaf linked list stays intact", func(t *testing.T) {
		// Insert enough to force a root split, then delete from each leaf
		tree := setup([]string{"a", "b", "c", "d", "e"})
		tree.delete("c")
		for _, k := range []string{"a", "b", "d", "e"} {
			_, found := tree.searchEntry(k)
			if !found {
				t.Errorf("key '%s' missing after deleting 'c'", k)
			}
		}
	})

	t.Run("delete same key twice — second returns false", func(t *testing.T) {
		tree := setup([]string{"a", "b", "c"})
		tree.delete("b")
		if tree.delete("b") {
			t.Error("expected false deleting an already-deleted key")
		}
	})
}

/* ------------------------------------------------------------------ */
/*  Helpers for underflow tests                                         */
/* ------------------------------------------------------------------ */

// buildLeafChain builds a parent with `count` leaf children each holding
// `entriesPerLeaf` entries. Returns (parent, leaves).
// Keys are auto-generated: leaf 0 gets "a0","a1"..., leaf 1 gets "b0","b1"... etc.
func buildLeafChain(order, count, entriesPerLeaf int) (*TreeNode, []*TreeNode) {
	leaves := make([]*TreeNode, count)
	for i := 0; i < count; i++ {
		entries := make([]*Entry, entriesPerLeaf)
		for j := 0; j < entriesPerLeaf; j++ {
			k := string(rune('a'+i)) + string(rune('0'+j))
			entries[j] = create_entry(k, k)
		}
		leaves[i] = create_tree_node(entries, nil, nil, true, order)
	}

	// Wire linked list
	for i := 0; i < count-1; i++ {
		leaves[i].next = leaves[i+1]
		leaves[i+1].prev = leaves[i]
	}

	// Build parent: separator keys are the first key of each non-first leaf
	keys := make([]string, count-1)
	for i := 1; i < count; i++ {
		keys[i-1] = leaves[i].leafEntries[0].key
	}
	parent := create_tree_node(nil, leaves, keys, false, order)
	for _, leaf := range leaves {
		leaf.parent = parent
	}
	return parent, leaves
}

// leafKeys returns the keys in a leaf node as a string slice.
func leafKeys(node *TreeNode) []string {
	out := make([]string, len(node.leafEntries))
	for i, e := range node.leafEntries {
		out[i] = e.key
	}
	return out
}

// parentKeys returns the separator keys of an internal node.
func parentKeys(node *TreeNode) []string {
	return slices.Clone(node.keys)
}

/* ------------------------------------------------------------------ */
/*  handle_leaf_underflow                                               */
/* ------------------------------------------------------------------ */

func Test_handle_leaf_underflow(t *testing.T) {

	t.Run("redistribute from right sibling when right has extra", func(t *testing.T) {
		// order=2: min entries = 2
		// left has 1 entry (underflow), right has 3 (can spare one)
		_, leaves := buildLeafChain(2, 2, 2)
		left, right := leaves[0], leaves[1]

		// Give right an extra entry so it has 3
		right.insert_leaf_entry(create_entry("b2", "b2"))

		// Force left into underflow by removing one entry directly (bypass delete)
		left.leafEntries = left.leafEntries[:1]

		left.handle_leaf_underflow()

		if len(left.leafEntries) != 2 {
			t.Errorf("left should have 2 entries after redistribute, got %d", len(left.leafEntries))
		}
		if len(right.leafEntries) != 2 {
			t.Errorf("right should have 2 entries after redistribute, got %d", len(right.leafEntries))
		}
	})

	t.Run("redistribute from right updates parent separator key", func(t *testing.T) {
		parent, leaves := buildLeafChain(2, 2, 2)
		left, right := leaves[0], leaves[1]

		right.insert_leaf_entry(create_entry("b2", "b2"))
		left.leafEntries = left.leafEntries[:1]
		left.handle_leaf_underflow()

		// Separator key must now equal the new first entry of right sibling
		wantSep := right.leafEntries[0].key
		if parent.keys[0] != wantSep {
			t.Errorf("parent separator = %q, want %q", parent.keys[0], wantSep)
		}
	})

	t.Run("redistribute from left sibling when left has extra", func(t *testing.T) {
		_, leaves := buildLeafChain(2, 2, 2)
		left, right := leaves[0], leaves[1]

		// Give left an extra entry
		left.insert_leaf_entry(create_entry("a2", "a2"))

		// Force right into underflow
		right.leafEntries = right.leafEntries[:1]

		right.handle_leaf_underflow()

		if len(right.leafEntries) != 2 {
			t.Errorf("right should have 2 entries after redistribute, got %d", len(right.leafEntries))
		}
		if len(left.leafEntries) != 2 {
			t.Errorf("left should have 2 entries after redistribute, got %d", len(left.leafEntries))
		}
	})

	t.Run("redistribute from left updates parent separator key", func(t *testing.T) {
		parent, leaves := buildLeafChain(2, 2, 2)
		left, right := leaves[0], leaves[1]

		left.insert_leaf_entry(create_entry("a2", "a2"))
		right.leafEntries = right.leafEntries[:1]

		right.handle_leaf_underflow()

		wantSep := right.leafEntries[0].key
		if parent.keys[0] != wantSep {
			t.Errorf("parent separator = %q, want %q", parent.keys[0], wantSep)
		}
	})

	t.Run("merge with right sibling when both at minimum", func(t *testing.T) {
		parent, leaves := buildLeafChain(2, 2, 2)
		left := leaves[0]

		// Drop left to underflow — right is already at minimum (2)
		left.leafEntries = left.leafEntries[:1]
		left.handle_leaf_underflow()

		// Left should have absorbed right's entries
		if len(left.leafEntries) != 3 {
			t.Errorf("merged leaf should have 3 entries, got %d", len(left.leafEntries))
		}
		// Parent should have lost the separator key
		if len(parent.keys) != 0 {
			t.Errorf("parent should have 0 keys after merge, got %d: %v", len(parent.keys), parent.keys)
		}
		// Parent should have lost the right child pointer
		if len(parent.childNodes) != 1 {
			t.Errorf("parent should have 1 child after merge, got %d", len(parent.childNodes))
		}
	})

	t.Run("merge fixes linked list — next pointer", func(t *testing.T) {
		_, leaves := buildLeafChain(2, 3, 2)
		left, mid, right := leaves[0], leaves[1], leaves[2]

		// Drop mid to underflow, right is at minimum — will merge mid+right
		mid.leafEntries = mid.leafEntries[:1]
		mid.handle_leaf_underflow()

		// mid absorbed right; mid.next should now skip right and point to right.next (nil)
		if mid.next != right.next {
			t.Errorf("linked list not fixed after merge: mid.next = %v, want %v", mid.next, right.next)
		}
		_ = left
	})

	t.Run("merge fixes linked list — prev pointer of successor", func(t *testing.T) {
		_, leaves := buildLeafChain(2, 3, 2)
		left, mid, right := leaves[0], leaves[1], leaves[2]

		// Add a fourth leaf after right so we can check its prev pointer
		extra := create_tree_node([]*Entry{create_entry("z0", "z0")}, nil, nil, true, 2)
		right.next = extra
		extra.prev = right

		mid.leafEntries = mid.leafEntries[:1]
		mid.handle_leaf_underflow()

		if extra.prev != mid {
			t.Errorf("extra.prev should point to mid after merge, got %v", extra.prev)
		}
		_ = left
	})

	t.Run("no-op when parent is nil (root leaf)", func(t *testing.T) {
		node := create_tree_node([]*Entry{create_entry("a", "a")}, nil, nil, true, 2)
		// Should not panic
		node.handle_leaf_underflow()
	})
}

/* ------------------------------------------------------------------ */
/*  handle_internal_underflow                                           */
/* ------------------------------------------------------------------ */

// buildInternalWithLeaves creates a 2-level tree:
//   parent (internal) → children (internal) → grandchildren (leaves)
// Each internal child has `keysPerChild` keys and `keysPerChild+1` leaf grandchildren.
func buildTwoLevelInternal(order, childCount, keysPerChild int) (*TreeNode, []*TreeNode) {
	children := make([]*TreeNode, childCount)
	for i := 0; i < childCount; i++ {
		// Build keysPerChild+1 leaf grandchildren for each child
		grandLeaves := make([]*TreeNode, keysPerChild+1)
		for j := 0; j <= keysPerChild; j++ {
			k := string(rune('a'+i)) + string(rune('0'+j))
			grandLeaves[j] = create_tree_node([]*Entry{create_entry(k, k)}, nil, nil, true, order)
		}
		keys := make([]string, keysPerChild)
		for j := 0; j < keysPerChild; j++ {
			keys[j] = grandLeaves[j+1].leafEntries[0].key
		}
		children[i] = create_tree_node(nil, grandLeaves, keys, false, order)
		for _, gl := range grandLeaves {
			gl.parent = children[i]
		}
	}

	// Build root
	rootKeys := make([]string, childCount-1)
	for i := 1; i < childCount; i++ {
		rootKeys[i-1] = children[i].childNodes[0].leafEntries[0].key
	}
	root := create_tree_node(nil, children, rootKeys, false, order)
	for _, child := range children {
		child.parent = root
	}
	return root, children
}

func Test_handle_internal_underflow(t *testing.T) {

	t.Run("redistribute from right internal sibling", func(t *testing.T) {
		// order=2: internal nodes need >= 2 keys
		// left child has 1 key (underflow), right child has 3 keys (can spare)
		root, children := buildTwoLevelInternal(2, 2, 2)
		left, right := children[0], children[1]

		// Give right an extra key by adding a grandchild leaf
		extraLeaf := create_tree_node([]*Entry{create_entry("extra", "extra")}, nil, nil, true, 2)
		extraLeaf.parent = right
		right.keys = append([]string{"a9"}, right.keys...)
		right.childNodes = append([]*TreeNode{extraLeaf}, right.childNodes...)

		// Drop left to underflow
		left.keys = left.keys[:1]
		left.childNodes = left.childNodes[:2]

		left.handle_internal_underflow()

		if len(left.keys) < root.order {
			t.Errorf("left still in underflow after redistribute: has %d keys, need >= %d", len(left.keys), root.order)
		}
		if len(right.keys) < root.order {
			t.Errorf("right went into underflow after redistribute: has %d keys", len(right.keys))
		}
	})

	t.Run("redistribute from left internal sibling", func(t *testing.T) {
		root, children := buildTwoLevelInternal(2, 2, 2)
		left, right := children[0], children[1]

		// Give left an extra key
		extraLeaf := create_tree_node([]*Entry{create_entry("extra", "extra")}, nil, nil, true, 2)
		extraLeaf.parent = left
		left.keys = append(left.keys, "a9")
		left.childNodes = append(left.childNodes, extraLeaf)

		// Drop right to underflow
		right.keys = right.keys[:1]
		right.childNodes = right.childNodes[:2]

		right.handle_internal_underflow()

		if len(right.keys) < root.order {
			t.Errorf("right still in underflow after redistribute: has %d keys, need >= %d", len(right.keys), root.order)
		}
		if len(left.keys) < root.order {
			t.Errorf("left went into underflow: has %d keys", len(left.keys))
		}
	})

	t.Run("merge internal nodes — parent loses a key", func(t *testing.T) {
		root, children := buildTwoLevelInternal(2, 2, 2)
		left := children[0]

		// Drop left to underflow — right sibling is at minimum so must merge
		left.keys = left.keys[:1]
		left.childNodes = left.childNodes[:2]

		left.handle_internal_underflow()

		if len(root.keys) != 0 {
			t.Errorf("root should have 0 keys after merge, got %d: %v", len(root.keys), root.keys)
		}
		if len(root.childNodes) != 1 {
			t.Errorf("root should have 1 child after merge, got %d", len(root.childNodes))
		}
	})

	t.Run("merge re-parents grandchildren to new owner", func(t *testing.T) {
		_, children := buildTwoLevelInternal(2, 2, 2)
		left, right := children[0], children[1]

		// Record right's grandchildren before merge
		grandchildren := slices.Clone(right.childNodes)

		left.keys = left.keys[:1]
		left.childNodes = left.childNodes[:2]
		left.handle_internal_underflow()

		// All of right's former grandchildren should now point to left as parent
		for i, gc := range grandchildren {
			if gc.parent != left {
				t.Errorf("grandchild[%d].parent = %p, want left (%p)", i, gc.parent, left)
			}
		}
	})

	t.Run("no-op when parent is nil (root internal node)", func(t *testing.T) {
		root, _ := buildTwoLevelInternal(2, 2, 2)
		root.parent = nil
		// Should not panic
		root.handle_internal_underflow()
	})
}

/* ------------------------------------------------------------------ */
/*  End-to-end: delete triggers underflow and tree stays valid          */
/* ------------------------------------------------------------------ */

func Test_underflow_end_to_end(t *testing.T) {
	// allKeysFound verifies every key in `want` is searchable in the tree
	allKeysFound := func(t *testing.T, tree *BTree, want []string) {
		t.Helper()
		for _, k := range want {
			_, found := tree.searchEntry(k)
			if !found {
				t.Errorf("key %q missing from tree", k)
			}
		}
	}

	t.Run("delete triggers leaf redistribute — remaining keys still searchable", func(t *testing.T) {
		// order=2: leaf holds max 4 entries, min 2
		// Insert 6 keys so we get at least 2 leaves, one with room to spare
		tree := create_new_tree(2)
		keys := []string{"a", "b", "c", "d", "e", "f"}
		for _, k := range keys {
			tree.insert_element(Entry{key: k, value: k})
		}
		tree.delete("a")
		allKeysFound(t, tree, []string{"b", "c", "d", "e", "f"})
	})

	t.Run("delete triggers leaf merge — remaining keys still searchable", func(t *testing.T) {
		tree := create_new_tree(2)
		keys := []string{"a", "b", "c", "d", "e"}
		for _, k := range keys {
			tree.insert_element(Entry{key: k, value: k})
		}
		// Delete enough to force a merge
		tree.delete("a")
		tree.delete("b")
		allKeysFound(t, tree, []string{"c", "d", "e"})
	})

	t.Run("leaf linked list intact after merge", func(t *testing.T) {
		tree := create_new_tree(2)
		for _, k := range []string{"a", "b", "c", "d", "e", "f"} {
			tree.insert_element(Entry{key: k, value: k})
		}
		tree.delete("a")
		tree.delete("b")

		// Walk the linked list and collect all keys
		leftmost := tree.root
		for !leftmost.isLeaf {
			leftmost = leftmost.childNodes[0]
		}
		var got []string
		for n := leftmost; n != nil; n = n.next {
			for _, e := range n.leafEntries {
				got = append(got, e.key)
			}
		}
		want := []string{"c", "d", "e", "f"}
		if !slices.Equal(got, want) {
			t.Errorf("linked list after merge = %v, want %v", got, want)
		}
	})

	t.Run("delete all but one key — tree still searchable", func(t *testing.T) {
		tree := create_new_tree(2)
		keys := []string{"a", "b", "c", "d", "e", "f", "g"}
		for _, k := range keys {
			tree.insert_element(Entry{key: k, value: k})
		}
		for _, k := range keys[:len(keys)-1] {
			tree.delete(k)
		}
		_, found := tree.searchEntry("g")
		if !found {
			t.Error("last remaining key 'g' not found")
		}
	})
}

func Test_create_new_tree_bulkload(t *testing.T) {
	makeEntrySlice := func(keys []string) []*Entry {
		out := make([]*Entry, len(keys))
		for i, k := range keys {
			out[i] = create_entry(k, k+"_val")
		}
		return out
	}

	t.Run("panics on order less than 2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for order < 2")
			}
		}()
		create_new_tree_bulkload(makeEntrySlice([]string{"a"}), 1)
	})

	t.Run("panics on empty entries", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty entries")
			}
		}()
		create_new_tree_bulkload([]*Entry{}, 2)
	})

	t.Run("single entry — root is a leaf", func(t *testing.T) {
		tree := create_new_tree_bulkload(makeEntrySlice([]string{"a"}), 2)
		if !tree.root.isLeaf {
			t.Error("expected root to be a leaf for single entry")
		}
	})

	t.Run("all entries searchable after bulkload", func(t *testing.T) {
		keys := []string{"a", "b", "c", "d", "e", "f", "g"}
		tree := create_new_tree_bulkload(makeEntrySlice(keys), 2)
		for _, k := range keys {
			e, found := tree.searchEntry(k)
			if !found {
				t.Errorf("key %q not found after bulkload", k)
			} else if e.value != k+"_val" {
				t.Errorf("key %q: expected value %q, got %q", k, k+"_val", e.value)
			}
		}
	})

	t.Run("root splits into internal node when leaves exceed max", func(t *testing.T) {
		// order=2: max 4 entries per leaf, so 5+ entries forces a split
		keys := []string{"a", "b", "c", "d", "e"}
		tree := create_new_tree_bulkload(makeEntrySlice(keys), 2)
		if tree.root.isLeaf {
			t.Error("expected internal root after bulkloading 5 entries with order 2")
		}
	})

	t.Run("leaf linked list is sorted and complete", func(t *testing.T) {
		keys := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
		tree := create_new_tree_bulkload(makeEntrySlice(keys), 2)

		leftmost := tree.root
		for !leftmost.isLeaf {
			leftmost = leftmost.childNodes[0]
		}
		var got []string
		for n := leftmost; n != nil; n = n.next {
			for _, e := range n.leafEntries {
				got = append(got, e.key)
			}
		}
		if !slices.Equal(got, keys) {
			t.Errorf("linked list = %v, want %v", got, keys)
		}
	})

	t.Run("prev pointers are wired correctly", func(t *testing.T) {
		keys := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
		tree := create_new_tree_bulkload(makeEntrySlice(keys), 2)

		leftmost := tree.root
		for !leftmost.isLeaf {
			leftmost = leftmost.childNodes[0]
		}
		var nodes []*TreeNode
		for n := leftmost; n != nil; n = n.next {
			nodes = append(nodes, n)
		}
		for i := 1; i < len(nodes); i++ {
			if nodes[i].prev != nodes[i-1] {
				t.Errorf("nodes[%d].prev is wrong", i)
			}
		}
	})

	t.Run("bulkload and insert produce same search results", func(t *testing.T) {
		keys := []string{"a", "b", "c", "d", "e", "f", "g"}
		entries := makeEntrySlice(keys)

		bulk := create_new_tree_bulkload(entries, 2)
		onebyone := create_new_tree(2)
		for _, e := range entries {
			onebyone.insert_element(*e)
		}

		for _, k := range keys {
			_, bulkFound := bulk.searchEntry(k)
			_, oneFound := onebyone.searchEntry(k)
			if bulkFound != oneFound {
				t.Errorf("key %q: bulkload found=%v, one-by-one found=%v", k, bulkFound, oneFound)
			}
		}
	})

	t.Run("large bulkload — all entries searchable", func(t *testing.T) {
		keys := make([]string, 100)
		for i := range keys {
			keys[i] = fmt.Sprintf("key%03d", i)
		}
		tree := create_new_tree_bulkload(makeEntrySlice(keys), 3)
		for _, k := range keys {
			_, found := tree.searchEntry(k)
			if !found {
				t.Errorf("key %q not found in large bulkload", k)
			}
		}
	})
}


func Test_searchRangeEntry(t *testing.T) {
	setup := func() *BTree {
		tree := create_new_tree(2)
		for _, k := range []string{"a", "b", "c", "d", "e", "f", "g"} {
			tree.insert_element(Entry{key: k, value: k + "_val"})
		}
		return tree
	}
 
	resultKeys := func(entries []*Entry) []string {
		keys := make([]string, len(entries))
		for i, e := range entries {
			keys[i] = e.key
		}
		return keys
	}
 
	t.Run("range with >= and <=", func(t *testing.T) {
		tree := setup()
		got := resultKeys(tree.searchRangeEntry(map[string]string{"e": "<=", "b": ">="}))
		want := []string{"b", "c", "d", "e"}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
 
	t.Run("strict range with > and <", func(t *testing.T) {
		tree := setup()
		got := resultKeys(tree.searchRangeEntry(map[string]string{"e": "<", "b": ">"}))
		want := []string{"c", "d"}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
 
	t.Run("upper bound only", func(t *testing.T) {
		tree := setup()
		got := resultKeys(tree.searchRangeEntry(map[string]string{"c": "<="}))
		want := []string{"a", "b", "c"}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
 
	t.Run("lower bound only", func(t *testing.T) {
		tree := setup()
		got := resultKeys(tree.searchRangeEntry(map[string]string{"e": ">="}))
		want := []string{"e", "f", "g"}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
 
	t.Run("exact match with ==", func(t *testing.T) {
		tree := setup()
		got := resultKeys(tree.searchRangeEntry(map[string]string{"c": "=="}))
		want := []string{"c"}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
 
	t.Run("no conditions returns all entries", func(t *testing.T) {
		tree := setup()
		got := resultKeys(tree.searchRangeEntry(map[string]string{}))
		want := []string{"a", "b", "c", "d", "e", "f", "g"}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
 
	t.Run("impossible conditions return empty", func(t *testing.T) {
		tree := setup()
		// key must be both > "e" and < "b" — impossible
		got := tree.searchRangeEntry(map[string]string{"e": ">", "b": "<"})
		if len(got) != 0 {
			t.Errorf("expected empty for impossible conditions, got %v", resultKeys(got))
		}
	})
 
	t.Run("empty tree returns nil", func(t *testing.T) {
		tree := create_new_tree(2)
		got := tree.searchRangeEntry(map[string]string{"e": "<=", "b": ">="})
		if got != nil {
			t.Errorf("expected nil for empty tree, got %v", got)
		}
	})
 
	t.Run("results are in sorted order", func(t *testing.T) {
		tree := setup()
		got := resultKeys(tree.searchRangeEntry(map[string]string{"f": "<=", "b": ">="}))
		if !slices.IsSorted(got) {
			t.Errorf("results not sorted: %v", got)
		}
	})
 
	t.Run("values on returned entries are correct", func(t *testing.T) {
		tree := setup()
		entries := tree.searchRangeEntry(map[string]string{"d": "<=", "b": ">="})
		for _, e := range entries {
			wantVal := e.key + "_val"
			if e.value != wantVal {
				t.Errorf("key %q: expected value %q, got %q", e.key, wantVal, e.value)
			}
		}
	})
 
	t.Run("search after delete reflects deletion", func(t *testing.T) {
		tree := setup()
		tree.delete("c")
		got := resultKeys(tree.searchRangeEntry(map[string]string{"e": "<=", "b": ">="}))
		want := []string{"b", "d", "e"}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}