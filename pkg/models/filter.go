package models

import (
	stderrors "errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Filter struct {
	ListFilter
	CondFilter
}

type ListFilter struct {
	Page     int `query:"page"`
	PageRows int `query:"page_rows"`
}

func (lf *ListFilter) GetLimitOffset() (int, int) {
	return lf.PageRows, (lf.Page - 1) * lf.PageRows
}

//
// Advance filtering
//
type CondFilter struct {
	QueryStr        string `query:"query" json:"query"`
	EntityTableName string
}

func (cf *CondFilter) BuildQuery() (string, error) {
	if cf.QueryStr == "" {
		return "true", nil
	}

	if err := checkParens(cf.QueryStr); err != nil {
		return "", err
	}

	es, err := cf.parse()
	if err != nil {
		return "", err
	}

	if err := cf.validate(es); err != nil {
		return "", err
	}

	return cf.generateQuery(es[0])
}

func checkParens(queryStr string) error {
	stack := []string{}
	for _, c := range queryStr {
		if c == '(' {
			stack = append(stack, string(c))
		} else if c == ')' {
			if len(stack) < 1 {
				return errors.New("syntax error")
			}
			stack = stack[:len(stack)-1]
		} else {
			// chars
		}
	}
	if len(stack) > 0 {
		return errors.New("syntax error")
	}
	return nil
}

type Expression struct {
	isElement bool
	Val       string

	Op       string
	Lhs, Rhs *Expression
}

func NewElement(s string) Expression {
	return Expression{
		isElement: true,
		Val:       s,
	}
}

func NewExpression(op string, l, r Expression) Expression {
	e := Expression{
		Op:  op,
		Lhs: &l,
		Rhs: &r,
	}
	// fmt.Println("NewExpression:", e)
	return e
}

func (e *Expression) String() string {
	if e.isElement {
		return "E[" + e.Val + "]"
	}
	return "Ex[" + e.Lhs.String() + e.Op + e.Rhs.String() + "]"
}

func (cf *CondFilter) parse() ([]Expression, error) {
	stack := []Expression{}

	open := 0
	includeSpace := false
	buff := []rune{}
	for _, t := range cf.QueryStr {
		// fmt.Println(string(t), "stack:", stack)
		if t == '(' {
			open += 1
			continue
		} else if t == ')' {
			if open <= 0 {
				return nil, errors.New("syntax error")
			}
			open -= 1

			if len(buff) > 0 {
				stack = append(stack, NewElement(string(buff)))
				// fmt.Println("new stack:", stack)
			}

			// reset buffer
			buff = buff[:0]

			// pop 3 elements from stack
			if len(stack) < 3 {
				return nil, errors.New("syntax error")
			}
			es := stack[len(stack)-3:]

			// process and put that element on stack
			// fmt.Println("popped:", es)
			ee := NewExpression(es[1].Val, es[0], es[2])
			stack = stack[:len(stack)-3]
			stack = append(stack, ee)
			// fmt.Println("new stack:", stack)
			continue
		} else if t == '\'' {
			includeSpace = !includeSpace
			buff = append(buff, t)
			continue
		}

		if open <= 0 && len(stack) == 0 {
			return nil, errors.New("syntax error")
		}

		if t == ' ' && includeSpace {
			buff = append(buff, t)
		} else if t == ' ' && !includeSpace {
			// element done
			// put it on stack
			if len(buff) > 0 {
				stack = append(stack, NewElement(string(buff)))
				// fmt.Println("new stack:", stack)
			}
			buff = buff[:0]
		} else {
			buff = append(buff, t)
		}
	}

	if len(buff) > 0 {
		return nil, errors.New("syntax error")
	}
	// fmt.Printf("len stack: %d, %+v\n", len(stack), stack)

	if len(stack) != 1 && len(stack)%3 != 0 {
		return nil, errors.New("can't combine condition clauses")
	}

	// fmt.Println("combining to build a single expression")
	if len(stack) > 1 {
		for {
			if len(stack) == 1 {
				break
			}
			es := stack[len(stack)-3:]

			// process and put that element on stack
			// fmt.Println("popped:", es)
			ee := NewExpression(es[1].Val, es[0], es[2])
			stack = stack[:len(stack)-3]
			stack = append(stack, ee)
			// fmt.Println("new stack:", stack)
		}
	}

	return stack, nil
}

var ErrUnknownOperator = stderrors.New("unknown operator used")
var ErrInvorrectElement = stderrors.New("element is neither of [symbol, int, str]")

var operators = map[string]string{
	"AND": "AND",
	"OR":  "OR",
	"LT":  "<",
	"GT":  ">",
	"EQ":  "=",
	"NE":  "!=",
}

var symbols = map[string]map[string]struct{}{
	"activities": {
		"ts":       {},
		"distance": {},
		"seconds":  {},
		"loc":      {},
		"w_cond":   {},
	},
	"users": {
		"name": {},
		"role": {},
	},
}

func (cf *CondFilter) validateElement(e Expression) error {
	// fmt.Println("validating:", e)

	// is symbol
	_, ok := symbols[cf.EntityTableName][e.Val]
	if ok {
		// fmt.Println(e.Val, "is symbol")
		return nil
	}

	// is number
	_, err := strconv.Atoi(e.Val)
	if err == nil {
		// fmt.Println(e.Val, "is num")
		return nil
	}

	// is string
	if strings.HasPrefix(e.Val, "'") && strings.HasSuffix(e.Val, "'") {
		// fmt.Println(e.Val, "is str")
		return nil
	}

	return errors.Wrapf(ErrInvorrectElement, "element: %v", e.Val)
}

func (cf *CondFilter) validateExpression(ex Expression) error {
	if ex.isElement {
		return cf.validateElement(ex)
	}

	if _, ok := operators[ex.Op]; !ok {
		return errors.Wrapf(ErrUnknownOperator, "op: %v", ex.Op)
	}

	if ex.Lhs != nil {
		if err := cf.validateExpression(*ex.Lhs); err != nil {
			return err
		}
	}
	if ex.Rhs != nil {
		if err := cf.validateExpression(*ex.Rhs); err != nil {
			return err
		}
	}
	return nil
}

func (cf *CondFilter) validate(es []Expression) error {
	for _, e := range es {
		if e.isElement {
			if err := cf.validateElement(e); err != nil {
				return err
			}
		} else {
			if err := cf.validateExpression(e); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cf *CondFilter) generateExpressionQuery(ex Expression) (string, error) {
	if ex.isElement {
		return ex.Val, nil
	}

	lq, err := cf.generateExpressionQuery(*ex.Lhs)
	if err != nil {
		return "", err
	}

	rq, err := cf.generateExpressionQuery(*ex.Rhs)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s %s %s)", lq, operators[ex.Op], rq), nil
}

func (cf *CondFilter) generateQuery(ex Expression) (string, error) {
	if ex.isElement {
		return "", errors.New("can't generate query with just the element")
	}

	q, err := cf.generateExpressionQuery(ex)
	if err != nil {
		return "", err
	}

	return q, nil
}
