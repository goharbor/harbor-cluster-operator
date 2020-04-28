package parse

import "errors"

type Parser struct {
	Type  string
	Value interface{}
}

func Fortype(typ string) *Parser {
	return &Parser{
		Type: typ,
	}
}

func (p *Parser) Parse(value interface{}) *Parser {
	p.Value = value
	return p
}

func (p *Parser) ToInt() (int, error) {
	var result int

	if p.Type == "int" {
		result = p.Value.(int)
		return result, nil
	}
	return result, errors.New("value type is not int")
}

func (p *Parser) ToString() (string, error) {
	var result string
	if p.Type == "string" {
		result = p.Value.(string)
		return result, nil
	}
	return result, errors.New("value type is not string")
}
