package user

import (
	"github.com/mdev5000/secretsanta/internal/types"
	"github.com/mdev5000/secretsanta/internal/util/scim"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func parseFilter(t *testing.T, e string) scim.Expression {
	ex, err := scim.ParseString(e)
	require.NoError(t, err)
	return ex
}

func parseExpr(t *testing.T, filter string) (bson.M, error) {
	f := NewFinder(nil)
	return f.ParseFilter(parseFilter(t, filter))
}

type testCase struct {
	expr     string
	expected bson.M
}

func TestParseFilter_simple(t *testing.T) {
	cases := []testCase{
		{
			`username eq "john01"`,
			bson.M{
				FieldUsername: bson.M{"$eq": "john01"},
			},
		},
		{
			`firstname ne "greg"`,
			bson.M{
				FieldFirstname: bson.M{"$ne": "greg"},
			},
		},
		{
			`families[id eq "family1"]"`,
			bson.M{
				FieldFamilyIds: bson.M{"$eq": "family1"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			b, err := parseExpr(t, tc.expr)
			require.NoError(t, err)
			require.Equal(t, tc.expected, b)
		})
	}
}

func TestParseFilter_binaryExpr(t *testing.T) {
	cases := []testCase{
		{
			`(username eq "john01") and ((firstname eq "fred") or (lastname ne "john"))"`,
			bson.M{
				"$and": bson.A{
					bson.M{FieldUsername: bson.M{"$eq": "john01"}},
					bson.M{
						"$or": bson.A{
							bson.M{FieldFirstname: bson.M{"$eq": "fred"}},
							bson.M{FieldLastname: bson.M{"$ne": "john"}},
						},
					},
				},
			},
		},
		{
			`(username eq "user1") and families[id eq "1" or id eq "3"]`,
			bson.M{
				"$and": bson.A{
					bson.M{FieldUsername: bson.M{"$eq": "user1"}},
					bson.M{
						"$or": bson.A{
							bson.M{FieldFamilyIds: bson.M{"$eq": "1"}},
							bson.M{FieldFamilyIds: bson.M{"$eq": "3"}},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			b, err := parseExpr(t, tc.expr)
			require.NoError(t, err)
			require.Equal(t, tc.expected, b)
		})
	}
}

type familiesFinderMock struct {
	mock.Mock
}

func (f *familiesFinderMock) FindFamilies(expr scim.Expression) ([]types.Family, error) {
	ret := f.Called(expr)
	return ret.Get(0).([]types.Family), ret.Error(1)
}

func TestParseFilter_familiesSubExpr(t *testing.T) {
	cases := []struct {
		name               string
		expr               string
		expectedFamilyExpr string
		familyRs           []types.Family
		expected           bson.M
	}{
		{
			name:               "sub expr",
			expr:               `families[name eq "family1" or name eq "family2"]`,
			expectedFamilyExpr: `name eq "family1" or name eq "family2"`,
			familyRs:           []types.Family{{ID: "family1"}, {ID: "family2"}},
			expected:           bson.M{FieldFamilyIds: bson.M{"$in": bson.A{"family1", "family2"}}},
		},
		{
			name:               "half-failed expr",
			expr:               `families[id eq "family-1" or name eq "family2"]`,
			expectedFamilyExpr: `id eq "family-1" or name eq "family2"`,
			familyRs:           []types.Family{{ID: "family1"}, {ID: "family2"}},
			expected:           bson.M{FieldFamilyIds: bson.M{"$in": bson.A{"family1", "family2"}}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expr := parseFilter(t, tc.expr)
			familiesExpr := parseFilter(t, tc.expectedFamilyExpr)

			familiesMock := &familiesFinderMock{}
			familiesMock.On("FindFamilies", familiesExpr).Return(tc.familyRs, nil)

			finder := NewFinder(familiesMock)
			b, err := finder.ParseFilter(expr)
			require.NoError(t, err)
			require.Equal(t, tc.expected, b)

			familiesMock.AssertExpectations(t)
		})
	}
}
