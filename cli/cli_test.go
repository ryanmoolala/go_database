package cli

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			expectedType    TokenType
			expectedLiteral string
		}
	}{
		{
			name: "full SQL input",
			input: `SELECT * FROM employees WHERE key >= "alice" AND key <= "mia";
INSERT employees VALUES ("bob", "designer");
DELETE "carol";
PRINT employees;`,
			expected: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{SELECT,    "SELECT"},
				{ASTERISK,  "*"},
				{FROM,      "FROM"},
				{IDENT,     "EMPLOYEES"},
				{WHERE,     "WHERE"},
				{IDENT,     "KEY"},
				{GTE,       ">="},
				{STRING,    "ALICE"},
				{AND,       "AND"},
				{IDENT,     "KEY"},
				{LTE,       "<="},
				{STRING,    "MIA"},
				{SEMICOLON, ";"},
				{INSERT,    "INSERT"},
				{IDENT,     "EMPLOYEES"},
				{VALUES,    "VALUES"},
				{LPAREN,    "("},
				{STRING,    "BOB"},
				{COMMA,     ","},
				{STRING,    "DESIGNER"},
				{RPAREN,    ")"},
				{SEMICOLON, ";"},
				{DELETE,    "DELETE"},
				{STRING,    "CAROL"},
				{SEMICOLON, ";"},
				{PRINT,     "PRINT"},
				{IDENT,     "EMPLOYEES"},
				{SEMICOLON, ";"},
				{EOF,       ""},
			},
		},
		{
			name:  "illegal tokens",
			input: `$2 elect;`,
			expected: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{ILLEGAL,   "$"},
				{NUMBER,    "2"},
				{IDENT,     "ELECT"},
				{SEMICOLON, ";"},
				{EOF,       ""},
			},
		},
		{
			name:  "comparators",
			input: `== != < > <= >=`,
			expected: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{EQ,  "=="},
				{NEQ, "!="},
				{LT,  "<"},
				{GT,  ">"},
				{LTE, "<="},
				{GTE, ">="},
				{EOF, ""},
			},
		},
		{
			name:  "empty string literal",
			input: `""`,
			expected: []struct {
				expectedType    TokenType
				expectedLiteral string
			}{
				{STRING, ""},
				{EOF,    ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input) // fresh lexer per subtest — no shared state
			for i, exp := range tt.expected {
				tok := l.NextToken()
				if tok.Type != exp.expectedType {
					t.Fatalf("token[%d] wrong type: expected=%q got=%q (literal=%q)",
						i, exp.expectedType, tok.Type, tok.Literal)
				}
				if tok.Literal != exp.expectedLiteral {
					t.Fatalf("token[%d] wrong literal: expected=%q got=%q",
						i, exp.expectedLiteral, tok.Literal)
				}
			}
		})
	}
}