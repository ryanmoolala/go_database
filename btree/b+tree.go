package main

import (
    "encoding/binary"
    "fmt" 
)

const (
    BNODE_NODE = 1 // internal nodes with pointers
    BNODE_LEAF = 2 // leaf nodes with values
)
const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000


//data structures
type Node struct { //Btree leaf nodes store values, internal nodes storepointers to nodes
	keys [][]byte

	vals [][]byte
	children []*Node
}

type BNode []byte  //byte arrays can be variable-length
// | 2B type | 2B numKeys |
// | 8B pointer |
// | 2B keyLen | key bytes (≤1000) |
// | 4B valLen | val bytes (≤3000) |

//using a node is all about loading and getting sliced data, and its less memory intensive than using structs in go


//functions 
func assert(cond bool) {
	if (!cond) {
		panic("Exceeding page size limits!")
	}
}

func init() {
    node1max := 4 + 1*8 + 1*2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
    assert(node1max<=BTREE_PAGE_SIZE)
}

// getters
func (node BNode) btype() uint16 {
    return binary.LittleEndian.Uint16(node[0:2])
}
func (node BNode) nkeys() uint16 {
    return binary.LittleEndian.Uint16(node[2:4])
}

// set header (btype, n keys)
func (node BNode) setHeader(btype uint16, nkeys uint16) {
    binary.LittleEndian.PutUint16(node[0:2], btype)
    binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

// get pointer
func (node BNode) getPtr(idx uint16) uint64 {
	assert(idx < node.nkeys())
	pos := 4+8*idx
	return binary.LittleEndian.Uint64(node[pos:])
}
// set pointer
func (node BNode) setPtr(idx uint16, val uint64) {
	assert(idx < node.nkeys())
	pos := 4+8*idx
	binary.LittleEndian.PutUint64(node[pos:], val)
}

// read the `offsets` array
// offsets[i] = cumulative size up to KV[i-1]
func (node BNode) getOffset(idx uint16) uint16 {
    if idx == 0 {
        return 0
    }
    pos := 4 + 8*node.nkeys() + 2*(idx-1)
    return binary.LittleEndian.Uint16(node[pos:])
}

//set the new offset 
func (node BNode) setOffset(idx uint16, offset uint16) {
    assert(idx > 0)
    pos := 4 + 8*node.nkeys() + 2*(idx-1)
    binary.LittleEndian.PutUint16(node[pos:], offset)
}


//reading the KV entry; relative
func (node BNode) kvPos(idx uint16) uint16 {
    assert(idx <= node.nkeys())
    return 4 + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

// getting key
func (node BNode) getKey(idx uint16) []byte {
    assert(idx < node.nkeys())
    pos := node.kvPos(idx) 
    klen := binary.LittleEndian.Uint16(node[pos:])
    return node[pos+4:][:klen]
}

// getting value
func (node BNode) getVal(idx uint16) []byte {
    assert(idx < node.nkeys())
    pos := node.kvPos(idx)
    klen := binary.LittleEndian.Uint16(node[pos+0:])
    vlen := binary.LittleEndian.Uint16(node[pos+2:])
    return node[pos+4+klen:][:vlen]
}

func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) {
    new.setPtr(idx, ptr) //set pointers

    pos := new.kvPos(idx) //offset value of the previous key

    binary.LittleEndian.PutUint16(new[pos+0:], uint16(len(key)))
    binary.LittleEndian.PutUint16(new[pos+2:], uint16(len(val)))

    copy(new[pos+4:], key)
    copy(new[pos+4+uint16(len(key)):], val)

    new.setOffset(idx+1, new.getOffset(idx) + 4 + uint16((len(key) + len(val))))
}

//calculate node size in bytes
func (node BNode) getSizeBytes() uint16 {
    return node.kvPos(node.nkeys())
}


//Print out its contents
func printNode(node BNode) {
    fmt.Println("----- BNode -----")
    fmt.Printf("Type: %d\n", node.btype())
    fmt.Printf("NumKeys: %d\n", node.nkeys())

    fmt.Println("Pointers:")
    for i := uint16(0); i < node.nkeys(); i++ {
        fmt.Printf("  [%d] -> %d\n", i, node.getPtr(i))
    }

    fmt.Println("Offsets:")
    for i := uint16(0); i <= node.nkeys(); i++ {
        fmt.Printf("  offset[%d] = %d\n", i, node.getOffset(i))
    }

    fmt.Println("Key-Value entries:")
    for i := uint16(0); i < node.nkeys(); i++ {
        key := node.getKey(i)
        val := node.getVal(i)

        fmt.Printf("  [%d] keyLen=%d valLen=%d\n", i, len(key), len(val))
        fmt.Printf("      key=%q\n", key)
        fmt.Printf("      val=%q\n", val)
    }

    fmt.Printf("Total node size: %d bytes\n", node.getSizeBytes())
    fmt.Println("-----------------")
}






