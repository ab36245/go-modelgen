package parser

import (
	"fmt"
	"path/filepath"

	"github.com/ab36245/go-modelgen/defs"
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

func (p *Parser) Parse() ([]defs.Model, error) {
	p.next()
	return p.parseModels()
}

func (p *Parser) parseModels() ([]defs.Model, error) {
	var models []defs.Model
	for !p.token.IsEOF() {
		model, err := p.parseModel()
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}

func (p *Parser) parseModel() (defs.Model, error) {
	model := defs.Model{}

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

func (p *Parser) parseFields() ([]defs.Field, error) {
	var fields []defs.Field
	for !p.token.IsChar('}') {
		field, err := p.parseField()
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func (p *Parser) parseField() (defs.Field, error) {
	field := defs.Field{}

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

func (p *Parser) parseType() (*defs.Type, error) {
	if p.token.IsChar('[') {
		p.next()
		if p.token.IsChar(']') {
			p.next()
			return p.parseArrayType()
		}
		return p.parseMapType()
	}
	if p.token.IsChar('?') {
		return p.parseOptionType()
	}
	return p.parseSimpleType()
}

func (p *Parser) parseSimpleType() (*defs.Type, error) {
	typ := &defs.Type{}
	if !p.token.Is(tkName) {
		return typ, p.expected("type name")
	}
	switch p.token.text {
	case "bool":
		typ.Kind = defs.BoolType
	case "bytes":
		typ.Kind = defs.BytesType
	case "float":
		typ.Kind = defs.FloatType
	case "int":
		typ.Kind = defs.IntType
	case "ref":
		typ.Kind = defs.RefType
	case "string":
		typ.Kind = defs.StringType
	case "time":
		typ.Kind = defs.TimeType
	default:
		typ.Kind = defs.ModelType
		typ.Name = p.token.text
	}
	p.next()
	return typ, nil
}

func (p *Parser) parseArrayType() (*defs.Type, error) {
	typ := &defs.Type{
		Kind: defs.ArrayType,
	}
	sub, err := p.parseType()
	if err != nil {
		return nil, err
	}
	typ.Sub = sub
	return typ, nil
}

func (p *Parser) parseMapType() (*defs.Type, error) {
	typ := &defs.Type{
		Kind: defs.MapType,
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

func (p *Parser) parseOptionType() (*defs.Type, error) {
	typ := &defs.Type{
		Kind: defs.OptionType,
	}
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
