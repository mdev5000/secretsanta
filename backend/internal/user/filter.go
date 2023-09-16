package user

import (
	"github.com/mdev5000/secretsanta/internal/util/scim"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	FilterFieldFamilies = "families"
)

func ParseFilter(expr scim.Expression) (bson.M, error) {
	switch e := expr.(type) {

	case scim.BinaryExpression:
		{
			op, err := scim.OpToBSON(e.CompareOperator)
			if err != nil {
				return nil, err
			}
			x, err := ParseFilter(e.X)
			if err != nil {
				return nil, err
			}

			y, err := ParseFilter(e.Y)
			if err != nil {
				return nil, err
			}
			return bson.M{
				op: bson.A{x, y},
			}, nil
		}

	case scim.AttrExpression, scim.ValuePath:
		return parseValue(expr)
	}
	return nil, scim.ErrInvalidExpr
}

func parseValue(expr scim.Expression) (bson.M, error) {
	switch e := expr.(type) {
	case scim.AttrExpression:
		return parseFilterAttr(e)
	case scim.ValuePath:
		return parseValuePath(e)
	}
	return nil, scim.ErrInvalidExpr
}

func parseValuePath(a scim.ValuePath) (bson.M, error) {
	switch a.AttributeName {
	case FilterFieldFamilies:
		return parseFamiliesExpr(a.ValueExpression)
	}
	return nil, scim.ErrInvalidExpr
}

func parseFamiliesExpr(expr scim.Expression) (bson.M, error) {
	switch e := expr.(type) {

	case scim.BinaryExpression:
		{
			op, err := scim.OpToBSON(e.CompareOperator)
			if err != nil {
				return nil, err
			}
			x, err := parseFamiliesExpr(e.X)
			if err != nil {
				return nil, err
			}

			y, err := parseFamiliesExpr(e.Y)
			if err != nil {
				return nil, err
			}
			return bson.M{
				op: bson.A{x, y},
			}, nil
		}

	case scim.AttrExpression:
		{
			if e.AttributePath.AttributeName != "id" {
				return nil, scim.ErrBadAttr{
					Reason: "unsupported user attribute",
					Attr:   e.AttributePath.String(),
				}
			}
			if err := scim.EnsureOperatorSupported(e, scim.EQ, scim.NE); err != nil {
				return nil, err
			}
			filter, err := scim.ToBSONStringFilter(e, e.CompareValue)
			return bson.M{
				FieldFamilyIds: filter,
			}, err
		}
	}
	return nil, scim.ErrInvalidExpr
}

func parseFilterAttr(a scim.AttrExpression) (bson.M, error) {
	switch a.AttributePath.AttributeName {
	case FieldUsername, FieldFirstname, FieldLastname:
		// @todo improve the error
		filter, err := scim.ToBSONStringFilter(a, a.CompareValue)
		return bson.M{
			a.AttributePath.AttributeName: filter,
		}, err
	}

	return nil, scim.ErrBadAttr{
		Reason: "unsupported user attribute",
		Attr:   a.AttributePath.String(),
	}
}
