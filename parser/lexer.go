package parser

import (
	"fmt"

	"github.com/ab36245/go-source/buffer"
)

func newLexer(b buffer.Buffer) func() Token {
	i := 0
	return func() Token {
		t := Token{}
		t.from = i
		this := b.RuneAt(i)
		next := b.RuneAt(i + 1)
		switch {
		case this.IsEOF():
			t.kind = tkEOF

		case this.IsSpace():
			i = lexSpace(&t, b, i)

		case this.Is('/') && next.Is('/'):
			i = lexCommentSingle(&t, b, i+2)

		case this.Is('/') && next.Is('*'):
			i = lexCommentMulti(&t, b, i+2)

		case this.IsLetter() || this.Is('_'):
			i = lexName(&t, b, i)

		case this.IsDigit(10):
			i = lexNumber(&t, b, i)

		case this.IsAny("-+") && next.IsDigit(10):
			i = lexNumber(&t, b, i+1)
			if this.Is('-') {
				switch t.kind {
				case tkFloat:
					t.fnum = -t.fnum
				case tkInt:
					t.fnum = -t.fnum
				default:
					// should be tkError
				}
			}

		case this.Is('"'):
			i = lexString(&t, b, i+1)

		default:
			i = lexChar(&t, b, i)
		}
		t.to = i
		return t
	}
}

func lexChar(t *Token, b buffer.Buffer, i int) int {
	t.kind = tkChar
	t.text += string(b.RuneAt(i))
	return i + 1
}

func lexCommentMulti(t *Token, b buffer.Buffer, i int) int {
	t.kind = tkComment
	nest := 0
	for {
		char := b.RuneAt(i)
		if char.IsEOF() {
			return i
		}
		next := b.RuneAt(i + 1)
		if char.Is('/') && next.Is('*') {
			t.text += string(char) + string(next)
			nest += 1
			i += 2
		} else if char.Is('*') && next.Is('/') {
			if nest == 0 {
				return i + 2
			}
			t.text += string(char) + string(next)
			nest -= 1
			i += 2
		} else {
			t.text += string(char)
			i += 1
		}
	}
}

func lexCommentSingle(t *Token, b buffer.Buffer, i int) int {
	t.kind = tkComment
	for {
		char := b.RuneAt(i)
		if char.IsEOF() {
			return i
		}
		if char.Is('\n') {
			return i + 1
		}
		t.text += string(char)
		i += 1
	}
}

func lexName(t *Token, b buffer.Buffer, i int) int {
	t.kind = tkName
	r := b.RuneAt(i)
	for r.IsLetter() || r.Is('_') || r.IsDigit(10) {
		t.text += string(r)
		i += 1
		r = b.RuneAt(i)
	}
	return i
}

func lexNumber(t *Token, b buffer.Buffer, i int) int {
	char := b.RuneAt(i)
	next := b.RuneAt(i + 1)
	base := 10
	if !char.Is('0') {
		base = 10
	} else if next.Is('b') {
		base = 2
		i += 2
	} else if next.Is('o') {
		base = 8
		i += 2
	} else if next.Is('d') {
		base = 10
		i += 2
	} else if next.Is('x') {
		base = 16
		i += 2
	} else if next.IsDigit(8) {
		base = 8
		i += 1
	} else {
		t.kind = tkError
		t.text = "invalid explicit number base"
		return i + 1
	}

	nextDigit := func() int64 {
		for b.RuneAt(i).Is('_') {
			i += 1
		}
		return int64(b.RuneAt(i).AsDigit(base))
	}

	t.kind = tkInt
	digit := nextDigit()
	if digit < 0 {
		t.kind = tkError
		t.text = fmt.Sprintf("invalid base %d number", base)
		return i + 1
	}
	for digit >= 0 {
		t.inum = t.inum*int64(base) + digit
		i += 1
		digit = nextDigit()
	}
	if !b.RuneAt(i).IsAny(".-+") {
		return i
	}

	t.kind = tkFloat
	t.fnum = float64(t.inum)
	if b.RuneAt(i).Is('.') {
		i += 1
		fdiv := float64(base)
		digit = nextDigit()
		if digit < 0 {
			t.kind = tkError
			t.text = fmt.Sprintf("invalid base %d float", base)
			return i + 1
		}
		for digit >= 0 {
			t.fnum += float64(digit) / fdiv
			fdiv *= 10
			i += 1
			digit = nextDigit()
		}
	}
	if b.RuneAt(i).IsAny("-+") {
		// TODO exponent
	}

	return i
}

func lexSpace(t *Token, b buffer.Buffer, i int) int {
	t.kind = tkSpace
	r := b.RuneAt(i)
	for r.IsSpace() {
		t.text += string(r)
		i += 1
		r = b.RuneAt(i)
	}
	return i
}

func lexString(t *Token, b buffer.Buffer, i int) int {
	t.kind = tkString
	for {
		r := b.RuneAt(i)
		if r.IsEOF() {
			t.kind = tkError
			t.text = "end-of-file in string"
			return i
		}
		if r.Is('"') {
			return i + 1
		}
		if !r.Is('\\') {
			t.text += string(r)
			i += 1
		} else {
			// TODO process escape sequence
		}
	}
}
