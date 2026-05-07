package cli

import (
	"fmt"
	"strconv"
)

type Parser struct {
	l *Lexer
	tokens []Token

	curToken  Token
	peekToken Token

	curTokenIndex int 
	peekTokenIndex int 

	errors    []string
}

func NewParser(l *Lexer, tokens []Token) *Parser {
	p := &Parser{l: l}

	p.tokens = tokens
	p.curTokenIndex = 0
	p.peekTokenIndex = 1
	
	p.nextToken()
	
	return p
}

func (p *Parser) nextToken() {
	if p.curTokenIndex < len(p.tokens) {
		p.curToken = p.tokens[p.curTokenIndex]
	} else {
		p.curToken = Token{Type: EOF, Literal: ""}
	}

	if p.peekTokenIndex < len(p.tokens) {
		p.peekToken = p.tokens[p.peekTokenIndex]
	} else {
		p.peekToken = Token{Type: EOF, Literal: ""}
	}
	//fmt.Printf("nextToken: curToken=%s peekToken=%s\n", p.curToken, p.peekToken) // debug
	p.curTokenIndex++
	p.peekTokenIndex++
}
	
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) curTokenIs(t TokenType) bool  { return p.curToken.Type == t }
func (p *Parser) peekTokenIs(t TokenType) bool { return p.peekToken.Type == t }
 
func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors,
		fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
	return false
}

//main entry point for parser
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
 
	for !p.curTokenIs(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
 
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case SELECT:
		return p.parseSelectStatement() 
	case INSERT:
		return p.parseInsertStatement() 
	case DELETE:
		return p.parseDeleteStatement()
	case PRINT:
		return p.parsePrintStatement() 
	case CREATE:
		return p.parseCreateStatement() 
	case BULKLOAD:
		return p.parseBulkloadStatement() 
	default:
		p.errors = append(p.errors,
			fmt.Sprintf("unknown statement starting with %s", p.curToken.Type))
		return nil
	}
}

//select statements like 
// SELECT * FROM EMPLOYEES WHERE key >= "alice" ;
// SELECT * FROM EMPLOYEES WHERE key >= "alice" AND key <= "mia" ;

func (p *Parser) parseSelectStatement() *SelectStatement {
	stmt := &SelectStatement{Token: p.curToken}
	 
	if !p.expectPeek(ASTERISK) {
		return nil
	}

	if !p.expectPeek(FROM) {
		return nil
	}

	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.Table = p.curToken.Literal

	if p.peekTokenIs(WHERE) {
		p.nextToken() // consume WHERE	
		if p.peekTokenIs(IDENT) {
			p.nextToken() // consume key identifier
		}
		stmt.Conditions = p.parseConditions()
	}
 
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConditions() []Condition {
	var conditions []Condition

	for {
		if !isComparator(p.peekToken.Type) {
			break
		}
		p.nextToken()
		comparator := p.curToken.Literal
 
		if !p.expectPeek(STRING) {
			break
		}
		limitKey := p.curToken.Literal
 
		conditions = append(conditions, Condition{
			LimitKey:   limitKey,
			Comparator: comparator,
		})
 
		if p.peekTokenIs(AND) {
			p.nextToken() // consume AND
			if p.peekTokenIs(IDENT) {
				p.nextToken() // consume next key identifier
			}
		} else {
			break
		}
	}
 
	return conditions
}
 
func isComparator(t TokenType) bool {
	switch t {
	case EQ, NEQ, LT, GT, LTE, GTE:
		return true
	}
	return false
}


// insert operations
//. INSERT <table> VALUES (<key>, <value>) ;        
func (p *Parser) parseInsertStatement() *InsertStatement {
	stmt := &InsertStatement{Token: p.curToken}
 
	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.Table = p.curToken.Literal
 
	if !p.expectPeek(VALUES) {
		return nil
	}
	if !p.expectPeek(LPAREN) {
		return nil
	}
	if !p.expectPeek(STRING) {
		return nil
	}
	stmt.Key = p.curToken.Literal
 
	if !p.expectPeek(COMMA) {
		return nil
	}
	if !p.expectPeek(STRING) {
		return nil
	}
	stmt.Value = p.curToken.Literal
 
	if !p.expectPeek(RPAREN) {
		return nil
	}
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}
 
	return stmt
}

//delete operations
// DELETE FROM <table> <key> ;   
func (p *Parser) parseDeleteStatement() *DeleteStatement {
	stmt := &DeleteStatement{Token: p.curToken}
 
	if !p.expectPeek(FROM) {
		return nil
	}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Table = p.curToken.Literal

	if !p.expectPeek(STRING) {
		return nil
	}

	stmt.Key = p.curToken.Literal
 
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}
 
	return stmt
}

//print operations
// PRINT <table> ;     
func (p *Parser) parsePrintStatement() *PrintStatement {
	stmt := &PrintStatement{Token: p.curToken}
 
	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.Table = p.curToken.Literal
 
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}
 
	return stmt
}

//create tree
// CREATE <table> <int>;
func (p *Parser) parseCreateStatement() *CreateStatement {
	stmt := &CreateStatement{Token: p.curToken}
 
	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.Table = p.curToken.Literal
 
	if !p.expectPeek(NUMBER) {
		return nil
	}
	var err error
	stmt.Order, err = strconv.Atoi(p.curToken.Literal)

	if err != nil {
		p.errors = append(p.errors, "invalid order value")
		return nil
	}
	 
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}	
 
	return stmt
}

//bulkload operations
//  BULKLOAD <table> VALUES (<key>,<val>), (<key>,<val>) ... ;
func (p *Parser) parseBulkloadStatement() *BulkloadStatement {
	stmt := &BulkloadStatement{Token: p.curToken}
 
	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.Table = p.curToken.Literal
 
	if !p.expectPeek(VALUES) {
		return nil
	}
 
	for p.peekTokenIs(LPAREN) {
		p.nextToken() // consume (
		if !p.expectPeek(STRING) {
			return nil
		}
		key := p.curToken.Literal
		if !p.expectPeek(COMMA) {
			return nil
		}
		if !p.expectPeek(STRING) {
			return nil
		}
		value := p.curToken.Literal
		if !p.expectPeek(RPAREN) {
			return nil
		}
		stmt.Entries = append(stmt.Entries, BulkEntry{Key: key, Value: value})
 
		if p.peekTokenIs(COMMA) {
			p.nextToken()
		}
	}
 
	// parser enforces at least one entry — evaluator handles the 2*order minimum
	if len(stmt.Entries) == 0 {
		p.errors = append(p.errors, "BULKLOAD requires at least one (key, value) pair")
		return nil
	}
 
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}
 
	return stmt
}