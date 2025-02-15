package ast

import (
	"errors"
	"fmt"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/parse/lexer"
	"github.com/Tom5521/xgotext/pkg/po/parse/token"
)

type Parser struct {
	input    []byte
	tokens   []token.Token
	position int

	nodes   []Node
	content []byte
	name    string
}

func (p *Parser) collectTokens(l *lexer.Lexer) {
	tok := l.NextToken()
	for tok.Type != token.EOF {
		p.tokens = append(p.tokens, tok)
		tok = l.NextToken()
	}
}

func NewParser(input []byte, filename string) *Parser {
	p := &Parser{
		input:   input,
		content: input,
		name:    filename,
	}

	p.collectTokens(lexer.New(input))

	return p
}

func NewParserFromString(input, filename string) *Parser {
	return NewParser([]byte(input), filename)
}

type parserFunc = func() (Node, error)

func (p *Parser) genParseMap() map[token.Type]parserFunc {
	return map[token.Type]parserFunc{
		token.COMMENT:      p.comment,
		token.MSGID:        p.msgid,
		token.MSGSTR:       p.msgstr,
		token.MSGCTXT:      p.msgctxt,
		token.PluralMsgid:  p.pluralMsgid,
		token.PluralMsgstr: p.pluralMsgstr,
	}
}

func (p *Parser) Normalizer() (*Normalizer, []error) {
	errs := p.Parse()
	return NewNormalizer(p.name, p.content, p.nodes), errs
}

func (p *Parser) Parse() []error {
	var errs []error

	parseMap := p.genParseMap()

	var tok token.Token
	for p.position, tok = range p.tokens {
		if len(errs) > 3 {
			errs = append(errs, errors.New("too many errors"))
			break
		}
		var node Node
		var err error
		switch tok.Type {
		case token.ILLEGAL:
			err = fmt.Errorf(
				"token at %s:%d is illegal",
				p.name,
				util.FindLine(p.input, tok.Pos),
			)
		case token.MSGID,
			token.MSGSTR,
			token.MSGCTXT,
			token.PluralMsgid,
			token.PluralMsgstr,
			token.COMMENT:
			parser := parseMap[tok.Type]
			node, err = parser()
		case token.STRING:
			continue
		default:
			err = fmt.Errorf(
				"unknown token type at %s:%d",
				p.name,
				util.FindLine(p.input, tok.Pos),
			)
		}

		if err != nil {
			errs = append(errs, err)
			continue
		}

		p.nodes = append(p.nodes, node)
	}

	return errs
}

func (p *Parser) Nodes() []Node {
	return p.nodes
}
