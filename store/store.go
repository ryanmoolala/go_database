//public kv interface
package store

type Store interface {
	Put(key, value []byte) error
	Get(key []byte) error
	Delete(key []byte) error
    Close() error
}