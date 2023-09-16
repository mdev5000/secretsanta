package scim

import (
	"bytes"
	"errors"
	"fmt"
	scim2 "github.com/scim2/filter-parser"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"strings"
)

// @todo test package

type Token = scim2.Token

var (
	CO = scim2.CO // contains

	EQ = scim2.EQ
	NE = scim2.NE
	GE = scim2.GE
	GT = scim2.GT
	LT = scim2.LT
	LE = scim2.LE

	AND = scim2.AND
	OR  = scim2.OR
)

var ErrInvalidExpr = errors.New("invalid expr")

type ErrBadOperator struct {
	Attr      string
	Op        string
	Supported []string
}

func (e ErrBadOperator) Error() string {
	return fmt.Sprintf(
		"unsupported SCIM operation '%s' for attribyte '%s' (supported %s)",
		e.Op,
		e.Attr,
		strings.Join(e.Supported, ","),
	)
}

type ErrBadAttr struct {
	Reason string
	Attr   string
}

func (e ErrBadAttr) Error() string {
	return fmt.Sprintf("bad SCIM attribute '%s': %s", e.Attr, e.Reason)
}

type (
	Expression       = scim2.Expression
	BinaryExpression = scim2.BinaryExpression
	UnaryExpression  = scim2.UnaryExpression
	AttrExpression   = scim2.AttributeExpression
	ValuePath        = scim2.ValuePath
	AttrPath         = scim2.AttributePath
)

func ToBSONNumericFilter(a AttrExpression, value interface{}) (bson.M, error) {
	var op string
	switch a.CompareOperator {
	case EQ:
		op = "$eq"
	case NE:
		op = "$ne"
	case GT:
		op = "$gt"
	case GE:
		op = "$gte"
	case LT:
		op = "$lt"
	case LE:
		op = "$le"
	default:
		return nil, fmt.Errorf("bad numeric filter operator '%s'", a.CompareOperator.String())
	}
	return bson.M{op: value}, nil
}

func ToBSONStringFilter(a AttrExpression, value string) (bson.M, error) {
	var op string
	switch a.CompareOperator {
	case EQ:
		op = "$eq"
	case NE:
		op = "$ne"
		// @todo add startsWith, etc.
	default:
		return nil, fmt.Errorf("bad string filter operator '%s'", a.CompareOperator.String())
	}
	return bson.M{op: value}, nil
}

func ParseString(s string) (Expression, error) {
	return Parse(bytes.NewBufferString(s))
}

func Parse(r io.Reader) (Expression, error) {
	p := scim2.NewParser(r)
	return p.Parse()
}

func EnsureOperatorSupported(a AttrExpression, supported ...Token) error {
	for _, op := range supported {
		if a.CompareOperator == op {
			return nil
		}
	}
	supportedS := make([]string, len(supported))
	for i, token := range supported {
		supportedS[i] = token.String()
	}
	return ErrBadOperator{
		Attr:      a.AttributePath.String(),
		Op:        a.CompareOperator.String(),
		Supported: supportedS,
	}
}
