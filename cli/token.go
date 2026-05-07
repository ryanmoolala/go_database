package cli

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

const (
	//query optr
	SELECT   TokenType = "SELECT"
	FROM     TokenType = "FROM"
	WHERE    TokenType = "WHERE"
	DELETE   TokenType = "DELETE"
	INSERT   TokenType = "INSERT"
	VALUES   TokenType = "VALUES"
	AND      TokenType = "AND"
	RANGE    TokenType = "RANGE"
	BETWEEN  TokenType = "BETWEEN"
	BULKLOAD TokenType = "BULKLOAD"
	PRINT    TokenType = "PRINT"
	HELP     TokenType = "HELP"
	CREATE   TokenType = "CREATE"
 
	// comparators 
	EQ  TokenType = "=="
	NEQ TokenType = "!="
	LT  TokenType = "<"
	GT  TokenType = ">"
	LTE TokenType = "<="
	GTE TokenType = ">="
 
	// literals
	IDENT  TokenType = "IDENT"  // tree name, key name e.g. employees
	STRING TokenType = "STRING" // quoted value e.g. "engineer"
	NUMBER TokenType = "NUMBER" // numeric value e.g. 42
 
	// punctuation
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	ASTERISK  TokenType = "*"
 
	// special
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"
)



var keywords = map[string]TokenType{
	"SELECT":   SELECT,
	"FROM":     FROM,
	"WHERE":    WHERE,
	"DELETE":   DELETE,
	"INSERT":   INSERT,
	"VALUES":   VALUES,
	"AND":      AND,
	"RANGE":    RANGE,
	"BETWEEN":  BETWEEN,
	"BULKLOAD": BULKLOAD,
	"PRINT":    PRINT,
	"HELP" : HELP,
	"CREATE": CREATE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
	return tok
	}
	return IDENT
}







