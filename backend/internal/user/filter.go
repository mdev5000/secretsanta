package user

import (
	"github.com/mdev5000/secretsanta/internal/types"
	"github.com/mdev5000/secretsanta/internal/util/scim"
	"github.com/mdev5000/secretsanta/internal/util/validator"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	FilterFieldFamilies = "families"
)

type FamilyFinder interface {
	FindFamilies(expr scim.Expression) ([]types.Family, error)
}

type Finder struct {
	families FamilyFinder
}

func NewFinder(families FamilyFinder) *Finder {
	return &Finder{families: families}
}

func (f *Finder) ParseFilter(expr scim.Expression) (bson.M, error) {
	switch e := expr.(type) {

	case scim.BinaryExpression:
		{
			op, err := scim.OpToBSON(e.CompareOperator)
			if err != nil {
				return nil, err
			}
			x, err := f.ParseFilter(e.X)
			if err != nil {
				return nil, err
			}

			y, err := f.ParseFilter(e.Y)
			if err != nil {
				return nil, err
			}
			return bson.M{
				op: bson.A{x, y},
			}, nil
		}

	case scim.AttrExpression, scim.ValuePath:
		return f.parseValue(expr)
	}
	return nil, scim.ErrInvalidExpr
}

func (f *Finder) parseValue(expr scim.Expression) (bson.M, error) {
	switch e := expr.(type) {
	case scim.AttrExpression:
		return parseFilterAttr(e)
	case scim.ValuePath:
		return f.parseValuePath(e)
	}
	return nil, scim.ErrInvalidExpr
}

func (f *Finder) parseValuePath(a scim.ValuePath) (bson.M, error) {
	switch a.AttributeName {
	case FilterFieldFamilies:
		b, err, ok := f.parseInlineFamiliesExpr(a.ValueExpression)
		if ok {
			return b, err
		}
		return f.parseFamiliesExpr(a.ValueExpression)
	}
	return nil, scim.ErrInvalidExpr
}

// parseInlineFamiliesExpr tries to avoid a second mongo call to families. If only the id is required
// the sub call can be avoided.
func (f *Finder) parseInlineFamiliesExpr(expr scim.Expression) (bson.M, error, bool) {
	switch e := expr.(type) {

	case scim.BinaryExpression:
		{
			op, err := scim.OpToBSON(e.CompareOperator)
			if err != nil {
				return nil, err, true
			}
			x, err, ok := f.parseInlineFamiliesExpr(e.X)
			if err != nil || !ok {
				return nil, err, ok
			}

			y, err, ok := f.parseInlineFamiliesExpr(e.Y)
			if err != nil || !ok {
				return nil, err, ok
			}
			return bson.M{
				op: bson.A{x, y},
			}, nil, true
		}

	case scim.AttrExpression:
		{
			if e.AttributePath.AttributeName != "id" {
				return nil, nil, false
			}
			if err := scim.EnsureOperatorSupported(e, scim.EQ, scim.NE); err != nil {
				return nil, err, true
			}
			filter, err := scim.ToBSONStringFilter(e, e.CompareValue)
			return bson.M{
				FieldFamilyIds: filter,
			}, err, true
		}
	}

	return nil, nil, false
}

func (f *Finder) parseFamiliesExpr(expr scim.Expression) (bson.M, error) {
	families, err := f.families.FindFamilies(expr)
	if err != nil {
		return nil, err
	}
	filter := make(bson.A, len(families))
	for i, family := range families {
		filter[i] = family.ID
	}
	return bson.M{
		FieldFamilyIds: bson.M{"$in": filter},
	}, err
}

func parseFilterAttr(a scim.AttrExpression) (bson.M, error) {
	switch a.AttributePath.AttributeName {
	case FieldUsername, FieldFirstname, FieldLastname:
		// @todo improve the error
		value := a.CompareValue
		err := validator.ValidateSingle(a.AttributePath.AttributeName, value, validator.StringSearchField)
		if err != nil {
			return nil, err
		}
		filter, err := scim.ToBSONStringFilter(a, value)
		return bson.M{
			a.AttributePath.AttributeName: filter,
		}, err
	}

	return nil, scim.ErrBadAttr{
		Reason: "unsupported user attribute",
		Attr:   a.AttributePath.String(),
	}
}
