package user

import (
	"github.com/mdev5000/secretsanta/internal/util/scim"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestParseFilter(t *testing.T) {
	parseFilter := func(e string) (bson.M, error) {
		ex, err := scim.ParseString(e)
		if err != nil {
			return nil, err
		}
		return ParseFilter(ex)
	}

	t.Run("simple expressions", func(t *testing.T) {
		cases := []struct {
			expr     string
			expected bson.M
		}{
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
		}

		for _, tc := range cases {
			t.Run(tc.expr, func(t *testing.T) {
				b, err := parseFilter(tc.expr)
				require.NoError(t, err)
				require.Equal(t, tc.expected, b)
			})
		}
	})
}
