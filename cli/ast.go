package cli

type Node interface {
	TokenLiteral() string
	StatementType() TokenType
}

type Statement interface {
	Node
	statementNode() // marker method to distinguish statements from expressions (if we had any)}
}

type SelectStatement struct {
	Token      Token  // the SELECT token
	Table      string // tree name
	Conditions []Condition // set if WHERE with comparators
}

func (s *SelectStatement) statementNode()       {} //no need to implement anything, just a marker method
func (s *SelectStatement) TokenLiteral() string { return s.Token.Literal }
func (s *SelectStatement) StatementType() TokenType { return s.Token.Type }

type InsertStatement struct {
	Token Token
	Table string
	Key   string
	Value string
}

func (s *InsertStatement) statementNode()       {} //no need to implement anything, just a marker method
func (s *InsertStatement) TokenLiteral() string { return s.Token.Literal }
func (s *InsertStatement) StatementType() TokenType { return s.Token.Type }

type DeleteStatement struct {
	Token Token
	Table string
	Key   string
}

func (s *DeleteStatement) statementNode()       {}
func (s *DeleteStatement) TokenLiteral() string { return s.Token.Literal }
func (s *DeleteStatement) StatementType() TokenType { return s.Token.Type }

type PrintStatement struct {
	Token Token
	Table string
}

func (s *PrintStatement) statementNode()       {}
func (s *PrintStatement) TokenLiteral() string { return s.Token.Literal }
func (s *PrintStatement) StatementType() TokenType { return s.Token.Type }

type CreateStatement struct {
	Token Token
	Table string
	Order int
}

func (s *CreateStatement) statementNode()       {}
func (s *CreateStatement) TokenLiteral() string { return s.Token.Literal }
func (s *CreateStatement) StatementType() TokenType { return s.Token.Type }

type BulkloadStatement struct {
	Token   Token
	Table   string
	Entries []BulkEntry
}

type BulkEntry struct {
	Key   string
	Value string
}

func (s *BulkloadStatement) statementNode()       {}
func (s *BulkloadStatement) TokenLiteral() string { return s.Token.Literal }
func (s *BulkloadStatement) StatementType() TokenType { return s.Token.Type }

type Condition struct {
	LimitKey   string // the value being compared against e.g. "alice"
	Comparator string // ==  !=  <  >  <=  >=
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}