package token

type Token struct {
	Type    string
	Literal string
}

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"
	STRING     = "STRING"
	ARRAY      = "ARRAY"

	ASSIGN       = "="
	EQUAL        = "=="
	NOT          = "!"
	NOT_EQUAL    = "!="
	LESS_THAN    = "<"
	GREATER_THAN = ">"

	PLUS     = "+"
	MINUS    = "-"
	DIVIDE   = "/"
	MULTIPLY = "*"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN     = "("
	RPAREN     = ")"
	LBRACE     = "{"
	RBRACE     = "}"
	LSQBRACKET = "["
	RSQBRACKET = "]"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]string{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func GetIdentType(ident string) string {
	if tokType, ok := keywords[ident]; ok {
		return tokType
	}
	return IDENTIFIER
}
