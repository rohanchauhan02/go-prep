//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"strconv"
	"unicode"
)

// Tokenizer: splits "3 + 5 * (2 - 1)" into ["3","+","5","*","(","2","-","1",")"]
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

type Parser struct{ tokens []string; pos int }

func (p *Parser) peek() string {
	if p.pos >= len(p.tokens) { return "" }
	return p.tokens[p.pos]
}
func (p *Parser) consume() string { t := p.peek(); p.pos++; return t }

// TODO: parseExpr handles + and - (lowest precedence)
// expr = term (('+' | '-') term)*
func (p *Parser) parseExpr() int {
	// TODO: implement
	return 0
}

// TODO: parseTerm handles * and / (higher precedence)
// term = factor (('*' | '/') factor)*
func (p *Parser) parseTerm() int {
	// TODO: implement
	return 0
}

// TODO: parseFactor handles numbers and parenthesized expressions
// factor = NUMBER | '(' expr ')'
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
	p := &Parser{tokens: tokenize(expr)}
	return p.parseExpr()
}

func main() {
	cases := []struct{ expr string; want int }{
		{"3 + 5", 8},
		{"10 - 3 * 2", 4},
		{"3 + 5 * (2 - 1)", 8},
		{"(2 + 3) * (4 - 1)", 15},
		{"100 / (2 + 3) * 2", 40},
	}
	for _, c := range cases {
		got := evaluate(c.expr)
		status := "✅"
		if got != c.want { status = "❌" }
		fmt.Printf("%s  %-30s = %d (want %d)
", status, c.expr, got, c.want)
	}
}
