package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// types of tokens
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// identifiers + literals
	IDFIER = "IDENTIFIER"
	INT    = "INT"
	FLOAT  = "FLOAT"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COMMENT   = "$$"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// string
	STRING = "STRING"

	// hash
	COLON = ":"

	// array support
	LBRACKET = "["
	RBRACKET = "]"

	// Language Keywords
	THOOS_MUJI    = "THOOS_MUJI"
	KAAM_GAR_MUJI = "KAAM_GAR_MUJI"

	YEDI_MUJI     = "YEDI_MUJI"
	NABHAE_CHIKNE = "NABHAE_CHIKNE"
	PATHA_MUJI    = "PATHA_MUJI"

	SACHO_MUJI = "SACHO_MUJI"
	JHUT_MUJI  = "JHUT_MUJI"

	JABA_SAMMA_MUJI = "JABA_SAMMA_MUJI"
	GHUMA_MUJI      = "GHUMA_MUJI"
)

func NewToken(t TokenType, r rune) Token {
	return Token{Type: t, Literal: string(r)}
}

func NewTokenFromStr(t TokenType, s string) Token {
	return Token{Type: t, Literal: s}
}

var keywords = map[string]TokenType{
	"thoos_muji":      THOOS_MUJI,
	"kaam_gar_muji":   KAAM_GAR_MUJI,
	"yedi_muji":       YEDI_MUJI,
	"nabhae_chikne":   NABHAE_CHIKNE,
	"sacho_muji":      SACHO_MUJI,
	"jhut_muji":       JHUT_MUJI,
	"patha_muji":      PATHA_MUJI,
	"jaba_samma_muji": JABA_SAMMA_MUJI,
	"ghuma_muji":      GHUMA_MUJI,
}

// this distinguishes reserved keywords from variable names
func LookupIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDFIER
}
