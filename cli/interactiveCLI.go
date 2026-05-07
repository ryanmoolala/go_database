package cli

import (
	"bufio"
	"fmt"
	"os"
	"godatabase/store"
)

func printIntro() {
	fmt.Println("enter help to get started")
}

const PROMPT = ">> "

var commandHistory = make([]string, 0)
var storageReference *store.Store = store.NewStore()
 
func Start() {

	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print(PROMPT)
 
		if !scanner.Scan() { 
			return
		}
 
		line := scanner.Text() //This waits for a line of input
 
		if line == "exit" {
			return
		}

		if line == "help" {
			fmt.Println()
			fmt.Println("┌──────────────────────────────────────────────────────────────────┐")
			fmt.Println("│                        Available Commands                        │")
			fmt.Println("├──────────────────────────────────────────────────────────────────┤")
			fmt.Println("│  CREATE <table> <order> ;                                        │")
			fmt.Println("│    Create a new B+ tree table                                    │")
			fmt.Println("│    e.g.  CREATE employees 2;                                     │")
			fmt.Println("├──────────────────────────────────────────────────────────────────┤")
			fmt.Println("│  INSERT <table> VALUES (<key>, <value>) ;                        │")
			fmt.Println("│    Insert a key-value pair into a table                          │")
			fmt.Println("│    e.g.  INSERT employees VALUES (\"alice\", \"engineer\");          │")
			fmt.Println("├──────────────────────────────────────────────────────────────────┤")
			fmt.Println("│  SELECT * FROM <table> WHERE key <op> <val> [AND key <op> <val>] │")
			fmt.Println("│    Search entries matching conditions                            │")
			fmt.Println("│    Operators: ==  !=  <  >  <=  >=                              │")
			fmt.Println("│    e.g.  SELECT * FROM employees WHERE key >= \"alice\";           │")
			fmt.Println("├──────────────────────────────────────────────────────────────────┤")
			fmt.Println("│  DELETE FROM <table> <key> ;                                     │")
			fmt.Println("│    Delete an entry by key                                        │")
			fmt.Println("│    e.g.  DELETE FROM employees \"alice\";                          │")
			fmt.Println("├──────────────────────────────────────────────────────────────────┤")
			fmt.Println("│  BULKLOAD <table> VALUES (<k>,<v>), (<k>,<v>), ... ;             │")
			fmt.Println("│    Create a new table loaded from a list of entries              │")
			fmt.Println("│    e.g.  BULKLOAD employees VALUES (\"a\",\"1\"),(\"b\",\"2\");          │")
			fmt.Println("├──────────────────────────────────────────────────────────────────┤")
			fmt.Println("│  PRINT <table> ;                                                 │")
			fmt.Println("│    Print the tree structure of a table                           │")
			fmt.Println("│    e.g.  PRINT employees;                                        │")
			fmt.Println("├──────────────────────────────────────────────────────────────────┤")
			fmt.Println("│  HELP ;    Show this menu                                        │")
			fmt.Println("│  exit      Quit the program                                      │")
			fmt.Println("└──────────────────────────────────────────────────────────────────┘")
			fmt.Println()
			continue
		}

		// if line == "prev" {
		// 	if len(commandHistory) == 0 {
		// 		fmt.Println("No previous commands.")
		// 		continue
		// 	}
		// 	line = commandHistory[len(commandHistory)-1]
		// 	fmt.Printf("Re-running previous command: %s\n", line)
		// }

		commandHistory = append(commandHistory, line) //save history
 
		tokens, ok, lexer := tokenize(line)
		if !ok {
			continue // illegal token found — skip to next input
		}
		
		parser := NewParser(lexer, tokens)
		program := parser.ParseProgram()
		if len(parser.Errors()) != 0 {
			fmt.Println("Parser errors:")
			for _, err := range parser.Errors() {
				fmt.Printf("\t%s\n", err)
			}
			continue
		}
		 
		Eval(program)
	}
}

func tokenize(line string) ([]Token, bool, *Lexer) {
	var tokens []Token
	l := New(line)
 
	for {
		tok := l.NextToken()
		if tok.Type == EOF {
			break
		}
		if tok.Type == ILLEGAL {
			fmt.Printf("Illegal token: %q — check your syntax\n", tok.Literal)
			return nil, false, nil
		}
		tokens = append(tokens, tok)
	}
	return tokens, true, l
}

