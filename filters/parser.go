/*
 Copyright 2021 Linka Cloud  All rights reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package filters

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenWord
	tokenString
	tokenLParen
	tokenRParen
	tokenComma
)

type token struct {
	typ   tokenType
	value string
	pos   int
}

type parser struct {
	tokens []token
	idx    int
}

var caseInsensitiveOps = map[string]struct{}{
	"eq":         {},
	"has_prefix": {},
	"has_suffix": {},
	"matches":    {},
	"in":         {},
	"inf":        {},
	"sup":        {},
}

// ParseExpression builds an Expression from its formatted representation.
// An empty string returns (nil, nil) to mirror Expression.Format().
func ParseExpression(input string) (*Expression, error) {
	if strings.TrimSpace(input) == "" {
		return nil, nil
	}
	p, err := newParser(input)
	if err != nil {
		return nil, err
	}
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if err := p.expectEOF(); err != nil {
		return nil, err
	}
	return expr, nil
}

// ParseFieldFilter builds a FieldFilter from its formatted representation.
// An empty string returns (nil, nil).
func ParseFieldFilter(input string) (*FieldFilter, error) {
	if strings.TrimSpace(input) == "" {
		return nil, nil
	}
	p, err := newParser(input)
	if err != nil {
		return nil, err
	}
	ff, err := p.parseFieldFilter()
	if err != nil {
		return nil, err
	}
	if err := p.expectEOF(); err != nil {
		return nil, err
	}
	return ff, nil
}

// ParseFilter builds a Filter from its formatted representation.
// An empty string returns (nil, nil).
func ParseFilter(input string) (*Filter, error) {
	if strings.TrimSpace(input) == "" {
		return nil, nil
	}
	p, err := newParser(input)
	if err != nil {
		return nil, err
	}
	f, err := p.parseFilter()
	if err != nil {
		return nil, err
	}
	if err := p.expectEOF(); err != nil {
		return nil, err
	}
	return f, nil
}

func newParser(input string) (*parser, error) {
	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}
	return &parser{tokens: tokens}, nil
}

func tokenize(input string) ([]token, error) {
	var tokens []token
	for idx := 0; idx < len(input); {
		r, w := utf8.DecodeRuneInString(input[idx:])
		switch {
		case unicode.IsSpace(r):
			idx += w
		case r == '\'':
			start := idx
			idx += w
			var sb strings.Builder
			closed := false
			for idx < len(input) {
				r, w = utf8.DecodeRuneInString(input[idx:])
				if r == '\\' {
					idx += w
					if idx >= len(input) {
						return nil, fmt.Errorf("filters: unterminated escape at %d", start)
					}
					r, w = utf8.DecodeRuneInString(input[idx:])
					sb.WriteRune(r)
					idx += w
					continue
				}
				if r == '\'' {
					idx += w
					tokens = append(tokens, token{typ: tokenString, value: sb.String(), pos: start})
					closed = true
					break
				}
				sb.WriteRune(r)
				idx += w
			}
			if closed {
				continue
			}
			return nil, fmt.Errorf("filters: unterminated string literal at %d", start)
		case r == '(':
			tokens = append(tokens, token{typ: tokenLParen, value: "(", pos: idx})
			idx += w
		case r == ')':
			tokens = append(tokens, token{typ: tokenRParen, value: ")", pos: idx})
			idx += w
		case r == ',':
			tokens = append(tokens, token{typ: tokenComma, value: ",", pos: idx})
			idx += w
		default:
			start := idx
			for idx < len(input) {
				r, w = utf8.DecodeRuneInString(input[idx:])
				if unicode.IsSpace(r) || r == '(' || r == ')' || r == ',' || r == '\'' {
					break
				}
				idx += w
			}
			tokens = append(tokens, token{typ: tokenWord, value: input[start:idx], pos: start})
		}
	}
	tokens = append(tokens, token{typ: tokenEOF, pos: len(input)})
	return tokens, nil
}

func (p *parser) parseExpression() (*Expression, error) {
	expr, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.error(p.peek(), "expected expression")
	}
	return expr, nil
}

func (p *parser) parseOr() (*Expression, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.peekWord("or") {
		p.next()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left.OrExprs = append(left.OrExprs, right)
	}
	return left, nil
}

func (p *parser) parseAnd() (*Expression, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for p.peekWord("and") {
		p.next()
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		left.AndExprs = append(left.AndExprs, right)
	}
	return left, nil
}

func (p *parser) parsePrimary() (*Expression, error) {
	tok := p.peek()
	if tok.typ == tokenLParen {
		p.next()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expectToken(tokenRParen); err != nil {
			return nil, err
		}
		return expr, nil
	}
	ff, err := p.parseFieldFilter()
	if err != nil {
		return nil, err
	}
	return &Expression{Condition: ff}, nil
}

func (p *parser) parseFieldFilter() (*FieldFilter, error) {
	tok := p.next()
	if tok.typ != tokenWord {
		return nil, p.error(tok, "expected field name")
	}
	filter, err := p.parseFilter()
	if err != nil {
		return nil, err
	}
	return &FieldFilter{Field: tok.value, Filter: filter}, nil
}

func (p *parser) parseFilter() (*Filter, error) {
	negated := false
	for {
		tok := p.peek()
		if tok.typ == tokenWord && strings.EqualFold(tok.value, "not") {
			negated = !negated
			p.next()
			continue
		}
		break
	}
	tok := p.next()
	if tok.typ != tokenWord {
		return nil, p.error(tok, "expected filter operator")
	}
	lower := strings.ToLower(tok.value)
	if lower == "is" {
		return p.parseIsFilter(negated)
	}
	op, ci := normalizeOperator(lower)
	switch op {
	case "eq":
		return p.parseEq(ci, negated)
	case "has_prefix":
		return p.parseStringFunc(ci, negated, func(val string) isStringFilter_Condition { return &StringFilter_HasPrefix{HasPrefix: val} })
	case "has_suffix":
		return p.parseStringFunc(ci, negated, func(val string) isStringFilter_Condition { return &StringFilter_HasSuffix{HasSuffix: val} })
	case "matches":
		return p.parseStringFunc(ci, negated, func(val string) isStringFilter_Condition { return &StringFilter_Regex{Regex: val} })
	case "in":
		return p.parseIn(ci, negated)
	case "inf":
		return p.parseOrder(op, ci, negated)
	case "sup":
		return p.parseOrder(op, ci, negated)
	case "before", "after":
		if ci {
			return nil, p.error(tok, "case insensitive modifier is invalid for %s", op)
		}
		return p.parseTimeComparison(op, negated)
	default:
		return nil, p.error(tok, "unexpected operator %q", tok.value)
	}
}

func (p *parser) parseIsFilter(negated bool) (*Filter, error) {
	tok := p.next()
	if tok.typ != tokenWord {
		return nil, p.error(tok, "expected value after 'is'")
	}
	switch strings.ToLower(tok.value) {
	case "null":
		return makeNullFilter(negated), nil
	case "true":
		return makeBoolFilter(negated, true), nil
	case "false":
		return makeBoolFilter(negated, false), nil
	default:
		return nil, p.error(tok, "unexpected value %q after 'is'", tok.value)
	}
}

func (p *parser) parseStringFunc(ci, negated bool, builder func(string) isStringFilter_Condition) (*Filter, error) {
	tok := p.next()
	if tok.typ != tokenString {
		return nil, p.error(tok, "expected quoted string value")
	}
	return makeStringFilter(ci, negated, builder(tok.value)), nil
}

func (p *parser) parseEq(ci, negated bool) (*Filter, error) {
	tok := p.next()
	lit, err := p.classifyLiteral(tok)
	if err != nil {
		return nil, err
	}
	switch lit.kind {
	case literalString:
		return makeStringFilter(ci, negated, &StringFilter_Equals{Equals: lit.str}), nil
	case literalTime:
		if ci {
			return nil, p.error(tok, "case insensitive modifier is invalid for time comparisons")
		}
		return makeTimeFilter(negated, &TimeFilter_Equals{Equals: timestamppb.New(lit.ts)}), nil
	case literalDuration:
		if ci {
			return nil, p.error(tok, "case insensitive modifier is invalid for durations")
		}
		return makeDurationFilter(negated, &DurationFilter_Equals{Equals: durationpb.New(lit.dur)}), nil
	case literalNumber:
		if ci {
			return nil, p.error(tok, "case insensitive modifier is invalid for numbers")
		}
		return makeNumberFilter(negated, &NumberFilter_Equals{Equals: lit.num}), nil
	default:
		return nil, p.error(tok, "unsupported literal for eq")
	}
}

func (p *parser) parseOrder(op string, ci, negated bool) (*Filter, error) {
	tok := p.next()
	lit, err := p.classifyLiteral(tok)
	if err != nil {
		return nil, err
	}
	switch lit.kind {
	case literalString:
		if op == "inf" {
			return makeStringFilter(ci, negated, &StringFilter_Inf{Inf: lit.str}), nil
		}
		return makeStringFilter(ci, negated, &StringFilter_Sup{Sup: lit.str}), nil
	case literalDuration:
		if ci {
			return nil, p.error(tok, "case insensitive modifier is invalid for durations")
		}
		if op == "inf" {
			return makeDurationFilter(negated, &DurationFilter_Inf{Inf: durationpb.New(lit.dur)}), nil
		}
		return makeDurationFilter(negated, &DurationFilter_Sup{Sup: durationpb.New(lit.dur)}), nil
	case literalNumber:
		if ci {
			return nil, p.error(tok, "case insensitive modifier is invalid for numbers")
		}
		if op == "inf" {
			return makeNumberFilter(negated, &NumberFilter_Inf{Inf: lit.num}), nil
		}
		return makeNumberFilter(negated, &NumberFilter_Sup{Sup: lit.num}), nil
	default:
		return nil, p.error(tok, "unsupported literal for %s", op)
	}
}

func (p *parser) parseIn(ci, negated bool) (*Filter, error) {
	if _, err := p.expectToken(tokenLParen); err != nil {
		return nil, err
	}
	peek := p.peek()
	if peek.typ == tokenRParen {
		return nil, p.error(peek, "expected at least one value in 'in' clause")
	}
	if peek.typ == tokenString {
		var values []string
		for {
			tok := p.next()
			if tok.typ != tokenString {
				return nil, p.error(tok, "expected quoted string value in 'in' clause")
			}
			values = append(values, tok.value)
			if p.peek().typ != tokenComma {
				break
			}
			p.next()
		}
		if _, err := p.expectToken(tokenRParen); err != nil {
			return nil, err
		}
		return makeStringFilter(ci, negated, &StringFilter_In_{In: &StringFilter_In{Values: values}}), nil
	}
	if ci {
		return nil, p.error(peek, "case insensitive modifier is invalid for numeric 'in'")
	}
	var numbers []float64
	for {
		tok := p.next()
		if tok.typ != tokenWord {
			return nil, p.error(tok, "expected number in 'in' clause")
		}
		val, err := strconv.ParseFloat(tok.value, 64)
		if err != nil {
			return nil, p.error(tok, "invalid number %q", tok.value)
		}
		numbers = append(numbers, val)
		if p.peek().typ != tokenComma {
			break
		}
		p.next()
	}
	if _, err := p.expectToken(tokenRParen); err != nil {
		return nil, err
	}
	return makeNumberFilter(negated, &NumberFilter_In_{In: &NumberFilter_In{Values: numbers}}), nil
}

func (p *parser) parseTimeComparison(op string, negated bool) (*Filter, error) {
	tok := p.next()
	if tok.typ != tokenWord {
		return nil, p.error(tok, "expected RFC3339 timestamp")
	}
	ts, err := time.Parse(time.RFC3339, tok.value)
	if err != nil {
		return nil, p.error(tok, "invalid RFC3339 timestamp %q", tok.value)
	}
	if op == "before" {
		return makeTimeFilter(negated, &TimeFilter_Before{Before: timestamppb.New(ts)}), nil
	}
	return makeTimeFilter(negated, &TimeFilter_After{After: timestamppb.New(ts)}), nil
}

type literalKind int

const (
	literalString literalKind = iota
	literalNumber
	literalDuration
	literalTime
)

type literalValue struct {
	kind literalKind
	str  string
	num  float64
	dur  time.Duration
	ts   time.Time
}

func (p *parser) classifyLiteral(tok token) (literalValue, error) {
	switch tok.typ {
	case tokenString:
		return literalValue{kind: literalString, str: tok.value}, nil
	case tokenWord:
		if ts, err := time.Parse(time.RFC3339, tok.value); err == nil {
			return literalValue{kind: literalTime, ts: ts}, nil
		}
		if dur, err := time.ParseDuration(tok.value); err == nil {
			return literalValue{kind: literalDuration, dur: dur}, nil
		}
		if num, err := strconv.ParseFloat(tok.value, 64); err == nil {
			return literalValue{kind: literalNumber, num: num}, nil
		}
	}
	return literalValue{}, p.error(tok, "invalid literal %q", tok.value)
}

func makeStringFilter(ci, negated bool, cond isStringFilter_Condition) *Filter {
	return &Filter{
		Match: &Filter_String_{
			String_: &StringFilter{
				Condition:       cond,
				CaseInsensitive: ci,
			},
		},
		Not: negated,
	}
}

func makeNumberFilter(negated bool, cond isNumberFilter_Condition) *Filter {
	return &Filter{
		Match: &Filter_Number{
			Number: &NumberFilter{Condition: cond},
		},
		Not: negated,
	}
}

func makeBoolFilter(negated bool, val bool) *Filter {
	return &Filter{
		Match: &Filter_Bool{Bool: &BoolFilter{Equals: val}},
		Not:   negated,
	}
}

func makeNullFilter(negated bool) *Filter {
	return &Filter{
		Match: &Filter_Null{Null: &NullFilter{}},
		Not:   negated,
	}
}

func makeTimeFilter(negated bool, cond isTimeFilter_Condition) *Filter {
	return &Filter{
		Match: &Filter_Time{Time: &TimeFilter{Condition: cond}},
		Not:   negated,
	}
}

func makeDurationFilter(negated bool, cond isDurationFilter_Condition) *Filter {
	return &Filter{
		Match: &Filter_Duration{Duration: &DurationFilter{Condition: cond}},
		Not:   negated,
	}
}

func normalizeOperator(word string) (string, bool) {
	lower := strings.ToLower(word)
	if strings.HasPrefix(lower, "i") {
		base := lower[1:]
		if _, ok := caseInsensitiveOps[base]; ok {
			return base, true
		}
	}
	return lower, false
}

func (p *parser) peek() token {
	if p.idx >= len(p.tokens) {
		return token{typ: tokenEOF, pos: p.tokens[len(p.tokens)-1].pos}
	}
	return p.tokens[p.idx]
}

func (p *parser) next() token {
	tok := p.peek()
	if p.idx < len(p.tokens) {
		p.idx++
	}
	return tok
}

func (p *parser) peekWord(word string) bool {
	tok := p.peek()
	return tok.typ == tokenWord && strings.EqualFold(tok.value, word)
}

func (p *parser) expectToken(tt tokenType) (token, error) {
	tok := p.next()
	if tok.typ != tt {
		return token{}, p.error(tok, "expected %s", tokenTypeName(tt))
	}
	return tok, nil
}

func (p *parser) expectEOF() error {
	tok := p.peek()
	if tok.typ != tokenEOF {
		return p.error(tok, "unexpected token %q", tok.value)
	}
	return nil
}

func (p *parser) error(tok token, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("filters: %s (at position %d)", msg, tok.pos)
}

func tokenTypeName(tt tokenType) string {
	switch tt {
	case tokenWord:
		return "word"
	case tokenString:
		return "string"
	case tokenLParen:
		return "'('"
	case tokenRParen:
		return "')'"
	case tokenComma:
		return "','"
	default:
		return "token"
	}
}
