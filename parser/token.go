package parser

import "fmt"

type Token struct {
	kind TokenKind
	from int
	to   int

	fnum float64
	inum int64
	text string
}

func (t Token) Is(kind TokenKind) bool {
	return t.kind == kind
}

func (t Token) IsChar(c rune) bool {
	return t.Is(tkChar) && t.text == string(c)
}

func (t Token) IsEOF() bool {
	return t.Is(tkEOF)
}

func (t Token) String() string {
	s := fmt.Sprintf("%-15s", t.kind)
	switch t.kind {
	case tkInvalid:
	case tkChar:
		s += fmt.Sprintf(": %s", t.text)
	case tkComment:
		s += fmt.Sprintf(": %q", t.text)
	case tkEOF:
	case tkError:
		s += fmt.Sprintf(": %q", t.text)
	case tkFloat:
		s += fmt.Sprintf(": %v", t.fnum)
	case tkInt:
		s += fmt.Sprintf(": %v", t.inum)
	case tkName:
		s += fmt.Sprintf(": %q", t.text)
	case tkSpace:
		s += fmt.Sprintf(": %q", t.text)
	case tkString:
		s += fmt.Sprintf(": %q", t.text)
	default:
		s += fmt.Sprintf(": unknown (%d)", t.kind)
	}
	return s
}

type TokenKind int

const (
	tkInvalid TokenKind = iota
	tkChar
	tkComment
	tkEOF
	tkError
	tkFloat
	tkInt
	tkName
	tkSpace
	tkString
)

func (k TokenKind) String() string {
	switch k {
	case tkInvalid:
		return "invalid"
	case tkChar:
		return "char"
	case tkComment:
		return "comment"
	case tkEOF:
		return "EOF"
	case tkError:
		return "error"
	case tkFloat:
		return "float"
	case tkInt:
		return "int"
	case tkName:
		return "name"
	case tkSpace:
		return "space"
	case tkString:
		return "string"
	default:
		return fmt.Sprintf("unknown (%d)", k)
	}
}
