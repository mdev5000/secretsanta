package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
)

var IgnoreFields = cmpopts.IgnoreFields

func Equal(t *testing.T, x, y interface{}, option ...cmp.Option) {
	if !cmp.Equal(x, y, option...) {
		t.Errorf("failed structs were not equal\n%s", cmp.Diff(x, y, option...))
	}
}
