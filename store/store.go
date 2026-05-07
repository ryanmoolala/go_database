//public kv interface
package store

import (
	"godatabase/bptree"
)

type Store struct {
	Storage map[string]*bptree.BTree
}

func NewStore() *Store {
	return &Store{
		Storage: make(map[string]*bptree.BTree),
	}
}
var GlobalSettings Store 
