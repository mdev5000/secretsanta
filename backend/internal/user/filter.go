package user

import (
	"github.com/mdev5000/secretsanta/internal/util/scim"
	"go.mongodb.org/mongo-driver/bson"
)

func ParseFilter(expr scim.Expression) (bson.M, error) {
	switch e := expr.(type) {
	case scim.BinaryExpression:
	case scim.AttrExpression:
		return parseFilterAttr(e)
	case scim.ValuePath:
	}
	return nil, scim.ErrInvalidExpr
}

//func parseValuePath(a scim.ValuePath) (bson.M, error) {
//	switch a.AttributeName {
//	case FieldFamilyIds:
//	}
//
//}

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
