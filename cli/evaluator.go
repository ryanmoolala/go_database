package cli

import (
	"fmt"
	"godatabase/bptree"
	"sort"
)

func Eval(program *Program) {
    for _, stmt := range program.Statements {
        switch s := stmt.(type) {
        case *InsertStatement:
            evalInsert(s)
        case *SelectStatement:
            evalSelect(s)
        case *DeleteStatement:
            evalDelete(s)
        case *PrintStatement:
            evalPrint(s)
        case *BulkloadStatement:
            evalBulkload(s)
		case *CreateStatement:
			evalCreate(s)
        }
    }
}

func evalCreate(stmt *CreateStatement) {
	table := stmt.Table
	// check for existing table
	// if not exists, create new B+ tree and add to store
	if _, exists := storageReference.Storage[table]; exists {
		fmt.Printf("Table %s already exists\n", table)
	} else {
		storageReference.Storage[table] = bptree.CreateNewTree(stmt.Order)
		fmt.Printf("Table %s created\n", table)
	}
}

func evalPrint(stmt *PrintStatement) {
	table := stmt.Table
	if tree, exists := storageReference.Storage[table]; exists {
		fmt.Printf("Contents of table %s:\n", table)
		bptree.PrintTree(tree.Root())
	} else {
		fmt.Printf("Table %s does not exist\n", table)
	}
}

func tableExists(table string) (*bptree.BTree, bool) {
	tree, exists := storageReference.Storage[table]
	if !exists {
		fmt.Printf("Error: table %q does not exist\n", table)
	}
	return tree, exists
}
 
func tableFree(table string) bool {
	if _, exists := storageReference.Storage[table]; exists {
		fmt.Printf("Error: table %q already exists\n", table)
		return false
	}
	return true
}

/* ------------------------------------------------------------------ */
/* INSERT <table> VALUES (<key>, <value>) ;                            */
/* ------------------------------------------------------------------ */
 
func evalInsert(stmt *InsertStatement) {
	if stmt.Key == "" {
		fmt.Println("Error: key cannot be empty")
		return
	}
	tree, ok := tableExists(stmt.Table)
	if !ok {
		return
	}
	if _, found := tree.Search(stmt.Key); found {
		fmt.Printf("Error: key %q already exists in %q\n", stmt.Key, stmt.Table)
		return
	}
	tree.Insert(bptree.CreateEntry(stmt.Key, stmt.Value))
	fmt.Printf("Inserted (%q, %q) into %q\n", stmt.Key, stmt.Value, stmt.Table)
}
 
/* ------------------------------------------------------------------ */
/* SELECT * FROM <table> WHERE key <op> <val> AND ... ;               */
/* ------------------------------------------------------------------ */
 
func evalSelect(stmt *SelectStatement) {
	tree, ok := tableExists(stmt.Table)
	if !ok {
		return
	}
	if len(stmt.Conditions) == 0 {
		fmt.Println("Error: SELECT requires at least one WHERE condition")
		return
	}
	conditions := make(map[string]string)
	for _, c := range stmt.Conditions {
		if c.LimitKey == "" {
			fmt.Println("Error: condition key cannot be empty")
			return
		}
		conditions[c.LimitKey] = c.Comparator
	}
	results := tree.SearchRange(conditions)
	if len(results) == 0 {
		fmt.Println("No results found")
		return
	}
	fmt.Printf("%d result(s):\n", len(results))
	for _, e := range results {
		fmt.Printf("  %s → %s\n", e.GetKey(), e.GetValue())
	}
}
 
/* ------------------------------------------------------------------ */
/* DELETE FROM <table> <key> ;                                         */
/* ------------------------------------------------------------------ */
 
func evalDelete(stmt *DeleteStatement) {
	if stmt.Key == "" {
		fmt.Println("Error: key cannot be empty")
		return
	}
	tree, ok := tableExists(stmt.Table)
	if !ok {
		return
	}
	if tree.Delete(stmt.Key) {
		fmt.Printf("Deleted %q from %q\n", stmt.Key, stmt.Table)
	} else {
		fmt.Printf("Error: key %q not found in %q\n", stmt.Key, stmt.Table)
	}
}
 
/* ------------------------------------------------------------------ */
/* BULKLOAD <table> VALUES (<key>,<val>), ... ;                       */
/* ------------------------------------------------------------------ */
 
func evalBulkload(stmt *BulkloadStatement) {
	if stmt.Table == "" {
		fmt.Println("Error: table name cannot be empty")
		return
	}
	if len(stmt.Entries) == 0 {
		fmt.Println("Error: BULKLOAD requires at least one entry")
		return
	}
	seen := make(map[string]bool)
	for _, e := range stmt.Entries {
		if e.Key == "" {
			fmt.Println("Error: entry key cannot be empty")
			return
		}
		if seen[e.Key] {
			fmt.Printf("Error: duplicate key %q in BULKLOAD entries\n", e.Key)
			return
		}
		seen[e.Key] = true
	}
	if !tableFree(stmt.Table) {
		return
	}
	entries := make([]*bptree.Entry, len(stmt.Entries))
	for i, e := range stmt.Entries {
		entries[i] = bptree.CreateEntry(e.Key, e.Value)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].GetKey() < entries[j].GetKey()
	})
	order := 2
	storageReference.Storage[stmt.Table] = bptree.CreateNewTreeBulkload(entries, order)
	fmt.Printf("Bulkloaded %d entries into table %q\n", len(entries), stmt.Table)
}