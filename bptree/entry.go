package bptree

type Entry struct {
	key string
	value string
}

func create_entry(key string, value string) *Entry {
	return &Entry {
		key: key,
		value: value, 
	}
}

func CreateEntry(key string, value string) *Entry {
	return &Entry{
		key:   key,
		value: value,
	}
}
func (e *Entry) GetKey() string   { return e.key }
func (e *Entry) GetValue() string { return e.value }
 