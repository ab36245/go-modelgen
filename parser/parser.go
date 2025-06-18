package parser

import (
	"fmt"
	"path/filepath"

	"github.com/ab36245/go-modelgen/defx"
	"github.com/ab36245/go-source/buffer"
)

func NewParser(path string) (*Parser, error) {
	buffer, err := buffer.File(path)
	if err != nil {
		return nil, err
	}
	lexer := newLexer(buffer)
	parser := &Parser{
		buffer: buffer,
		lexer:  lexer,
		path:   path,
	}
	return parser, nil
}

type Parser struct {
	buffer buffer.Buffer
	errors []parseError
	lexer  func() Token
	path   string
	token  Token
}

func (p *Parser) Parse() ([]defx.Model, error) {
	p.next()
	return p.parseModels()
}

func (p *Parser) parseModels() ([]defx.Model, error) {
	var models []defx.Model
	for !p.token.IsEOF() {
		model, err := p.parseModel()
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}

func (p *Parser) parseModel() (defx.Model, error) {
	model := defx.Model{}

	if !p.token.Is(tkInt) {
		return model, p.expected("model id")
	}
	model.Id = int(p.token.inum)
	p.next()

	if !p.token.Is(tkName) {
		return model, p.expected("model name")
	}
	model.Name = p.token.text
	p.next()

	if !p.token.IsChar('{') {
		return model, p.expected("'{'")
	}
	p.next()

	fields, err := p.parseFields()
	if err != nil {
		return model, err
	}
	model.Fields = fields

	if !p.token.IsChar('}') {
		return model, p.expected("'}'")
	}
	p.next()

	return model, nil
}

func (p *Parser) parseFields() ([]defx.Field, error) {
	var fields []defx.Field
	for !p.token.IsChar('}') {
		field, err := p.parseField()
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func (p *Parser) parseField() (defx.Field, error) {
	field := defx.Field{}

	if !p.token.Is(tkName) {
		return field, p.expected("field name")
	}
	field.Name = p.token.text
	p.next()

	if !p.token.IsChar(':') {
		return field, p.expected("':'")
	}
	p.next()

	typ, err := p.parseType()
	if err != nil {
		return field, err
	}
	field.Type = typ

	return field, nil
}

func (p *Parser) parseType() (*defx.Type, error) {
	if !p.token.IsChar('[') {
		return p.parseSimpleType()
	}
	p.next()
	if p.token.IsChar(']') {
		p.next()
		return p.parseArrayType()
	}
	return p.parseMapType()
}

func (p *Parser) parseSimpleType() (*defx.Type, error) {
	typ := &defx.Type{}
	if !p.token.Is(tkName) {
		return typ, p.expected("type name")
	}
	switch p.token.text {
	case "bool":
		typ.Kind = defx.BoolType
	case "bytes":
		typ.Kind = defx.BytesType
	case "float":
		typ.Kind = defx.FloatType
	case "int":
		typ.Kind = defx.IntType
	case "ref":
		typ.Kind = defx.RefType
	case "string":
		typ.Kind = defx.StringType
	case "time":
		typ.Kind = defx.TimeType
	default:
		typ.Kind = defx.ModelType
		typ.Name = p.token.text
	}
	p.next()
	return typ, nil
}

func (p *Parser) parseArrayType() (*defx.Type, error) {
	typ := &defx.Type{
		Kind: defx.ArrayType,
	}
	sub, err := p.parseType()
	if err != nil {
		return nil, err
	}
	typ.Sub = sub
	return typ, nil
}

func (p *Parser) parseMapType() (*defx.Type, error) {
	typ := &defx.Type{
		Kind: defx.MapType,
	}

	key, err := p.parseSimpleType()
	if err != nil {
		return nil, err
	}
	typ.Key = key

	if !p.token.IsChar(']') {
		return nil, p.expected("']'")
	}
	p.next()

	sub, err := p.parseType()
	if err != nil {
		return nil, err
	}
	typ.Sub = sub

	return typ, nil
}

func (p *Parser) expected(want any) error {
	var mesg string
	if p.token.Is(tkError) {
		mesg = p.token.text
	} else {
		mesg = fmt.Sprintf("expected %v, found %s", want, p.token.kind)
	}
	err := parseError{
		path:   p.path,
		buffer: p.buffer,
		mesg:   mesg,
		from:   p.token.from,
		to:     p.token.to,
	}
	return err
}

func (p *Parser) next() {
	for {
		token := p.lexer()
		if !token.Is(tkComment) && !token.Is(tkSpace) {
			p.token = token
			return
		}
	}
}

type parseError struct {
	path   string
	buffer buffer.Buffer
	mesg   string
	from   int
	to     int
}

func (e parseError) Error() string {
	path, _ := filepath.Abs(e.path)
	r := e.buffer.Range(e.from, e.to-1)
	s := fmt.Sprintf("%s:%d : %s\n", path, r.From.Line.Number, e.mesg)
	s += r.Show()
	return s
}
