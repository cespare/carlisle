package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// This file defines a simple AST for arithmetic. It's represented by the user
// using an s-expression.

type Arith interface {
	Eval(*MoveResizeState) (float64, error)
	String() string
}

func ParseArith(b []byte) (Arith, error) {
	arith, rest, err := parseArith(b)
	if err != nil {
		return nil, fmt.Errorf("error parsing %q: %s", b, err)
	}
	if len(bytes.TrimSpace(rest)) != 0 {
		return nil, fmt.Errorf("trailing junk after arithmetic expression: %q", rest)
	}
	return arith, nil
}

func parseArith(b []byte) (arith Arith, rest []byte, err error) {
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return nil, nil, errUnexpectedEnd
	}
	switch {
	case b[0] == '(':
		return parseArithFn(b)
	case b[0] >= '0' && b[0] <= '9':
		return parseArithConst(b)
	}
	return parseArithAtom(b)
}

func scanTok(b []byte) (first, rest []byte, err error) {
	var tok []byte
	for i, c := range b {
		switch c {
		case '\r', '\n', ' ', '\t', '(', ')':
			tok = bytes.TrimSpace(tok)
			if len(tok) == 0 {
				return nil, nil, errUnexpectedEnd
			}
			return tok, bytes.TrimSpace(b[i:]), nil
		}
		tok = append(tok, c)
	}
	return tok, nil, nil
}

type FnType int

const (
	FnAdd FnType = iota
	FnSub
	FnMul
	FnDiv
	FnMin
	FnMax
)

var FnToSymbol = map[FnType]string{
	FnAdd: "+",
	FnSub: "-",
	FnMul: "*",
	FnDiv: "/",
	FnMin: "min",
	FnMax: "max",
}

var SymbolToFn = make(map[string]FnType)

type ArithFn struct {
	Type FnType
	Args []Arith
}

var errUnexpectedEnd = errors.New("reached end of expression unexpectedly")

func parseArithFn(b []byte) (*ArithFn, []byte, error) {
	b = b[1:]
	first, rest, err := scanTok(b)
	if err != nil {
		return nil, nil, err
	}
	typ, ok := SymbolToFn[string(first)]
	if !ok {
		return nil, nil, fmt.Errorf("unrecognized function: %s", string(first))
	}
	b = rest
	args := []Arith{}
	var arith Arith
	for {
		if len(b) == 0 {
			return nil, nil, errUnexpectedEnd
		}
		if b[0] == ')' {
			rest = b[1:]
			break
		}
		arith, rest, err = parseArith(b)
		if err != nil {
			return nil, nil, err
		}
		args = append(args, arith)
		b = rest
	}
	return &ArithFn{typ, args}, rest, nil
}

func (a *ArithFn) Eval(state *MoveResizeState) (float64, error) {
	eval := make([]float64, len(a.Args))
	for i, arg := range a.Args {
		result, err := arg.Eval(state)
		if err != nil {
			return 0, err
		}
		eval[i] = result
	}
	switch a.Type {
	case FnAdd:
		if len(eval) != 2 {
			return 0, fmt.Errorf("+ expects 2 arguments")
		}
		return eval[0] + eval[1], nil
	case FnSub:
		if len(eval) != 2 {
			return 0, fmt.Errorf("- expects 2 arguments")
		}
		return eval[0] - eval[1], nil
	case FnMul:
		if len(eval) != 2 {
			return 0, fmt.Errorf("* expects 2 arguments")
		}
		return eval[0] * eval[1], nil
	case FnDiv:
		if len(eval) != 2 {
			return 0, fmt.Errorf("/ expects 2 arguments")
		}
		return eval[0] / eval[1], nil
	case FnMin:
		if len(eval) < 2 {
			return 0, fmt.Errorf("min expects at least 2 arguments")
		}
		min := eval[0]
		for _, e := range eval[1:] {
			if e < min {
				min = e
			}
		}
		return min, nil
	case FnMax:
		if len(eval) < 2 {
			return 0, fmt.Errorf("max expects at least 2 arguments")
		}
		max := eval[0]
		for _, e := range eval[1:] {
			if e > max {
				max = e
			}
		}
		return max, nil
	}
	panic("unreached")
}

func (a *ArithFn) String() string {
	args := []string{}
	for _, arg := range a.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("fn[%s](%s)", FnToSymbol[a.Type], strings.Join(args, ", "))
}

type ArithAtom int

const (
	AtomX ArithAtom = iota
	AtomY
	AtomWidth
	AtomHeight
	AtomScreenWidth
	AtomScreenHeight
)

var AtomToSymbol = map[ArithAtom]string{
	AtomX:            "x",
	AtomY:            "y",
	AtomWidth:        "w",
	AtomHeight:       "h",
	AtomScreenWidth:  "sw",
	AtomScreenHeight: "sh",
}

var SymbolToAtom = make(map[string]ArithAtom)

func parseArithAtom(b []byte) (ArithAtom, []byte, error) {
	first, rest, err := scanTok(b)
	if err != nil {
		return 0, nil, err
	}
	atom, ok := SymbolToAtom[string(first)]
	if !ok {
		return 0, nil, fmt.Errorf("bad atom: %s", string(first))
	}
	return atom, rest, nil
}

func (a ArithAtom) Eval(state *MoveResizeState) (float64, error) {
	switch a {
	case AtomX:
		return float64(state.X), nil
	case AtomY:
		return float64(state.Y), nil
	case AtomWidth:
		return float64(state.W), nil
	case AtomHeight:
		return float64(state.H), nil
	case AtomScreenWidth:
		return float64(state.ScreenW), nil
	case AtomScreenHeight:
		return float64(state.ScreenH), nil
	}
	panic("unreached")
}

func (a ArithAtom) String() string {
	return fmt.Sprintf("atom(%s)", AtomToSymbol[a])
}

type ArithConst float64

func parseArithConst(b []byte) (ArithConst, []byte, error) {
	first, rest, err := scanTok(b)
	if err != nil {
		return 0, nil, err
	}
	f, err := strconv.ParseFloat(string(first), 64)
	if err != nil {
		return 0, nil, err
	}
	return ArithConst(f), rest, nil
}

func (c ArithConst) Eval(*MoveResizeState) (float64, error) {
	return float64(c), nil
}

func (c ArithConst) String() string {
	return fmt.Sprintf("const(%.1f)", c)
}

func init() {
	for k, v := range FnToSymbol {
		SymbolToFn[v] = k
	}
	for k, v := range AtomToSymbol {
		SymbolToAtom[v] = k
	}
}
