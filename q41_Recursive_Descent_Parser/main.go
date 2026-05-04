package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Parser struct{ tokens []string; pos int }

func tokenize(expr string) []string {
	var tokens []string
	i := 0
	for i < len(expr) {
		if unicode.IsSpace(rune(expr[i])) { i++; continue }
		if unicode.IsDigit(rune(expr[i])) {
			j := i
			for j < len(expr) && unicode.IsDigit(rune(expr[j])) { j++ }
			tokens = append(tokens, expr[i:j]); i = j
		} else {
			tokens = append(tokens, string(expr[i])); i++
		}
	}
	return tokens
}

func (p *Parser) peek() string {
	if p.pos >= len(p.tokens) { return "" }
	return p.tokens[p.pos]
}

func (p *Parser) consume() string {
	t := p.peek(); p.pos++; return t
}

// expr = term (('+' | '-') term)*
func (p *Parser) parseExpr() int {
	result := p.parseTerm()
	for p.peek() == "+" || p.peek() == "-" {
		op := p.consume()
		right := p.parseTerm()
		if op == "+" { result += right } else { result -= right }
	}
	return result
}

// term = factor (('*' | '/') factor)*
func (p *Parser) parseTerm() int {
	result := p.parseFactor()
	for p.peek() == "*" || p.peek() == "/" {
		op := p.consume()
		right := p.parseFactor()
		if op == "*" { result *= right } else { result /= right }
	}
	return result
}

// factor = number | '(' expr ')'
func (p *Parser) parseFactor() int {
	t := p.peek()
	if t == "(" {
		p.consume()
		result := p.parseExpr()
		p.consume() // ')'
		return result
	}
	p.consume()
	n, _ := strconv.Atoi(t)
	return n
}

func evaluate(expr string) int {
	tokens := tokenize(expr)
	p := &Parser{tokens: tokens}
	return p.parseExpr()
}

func main() {
	exprs := []string{
		"3 + 5",
		"10 - 3 * 2",
		"3 + 5 * (2 - 1)",
		"(2 + 3) * (4 - 1)",
		"100 / (2 + 3) * 2",
	}
	for _, e := range exprs {
		fmt.Printf("%-30s = %d\n", strings.TrimSpace(e), evaluate(e))
	}
}
