package scim

import (
	"bytes"
	"fmt"
	scim2 "github.com/scim2/filter-parser"
	"github.com/stretchr/testify/require"
	"testing"
)

func parse(e Expression) {
	switch ex := e.(type) {
	case BinaryExpression:
		fmt.Println(ex.CompareOperator)
		fmt.Println("X:")
		parse(ex.X)
		fmt.Println("Y:")
		parse(ex.Y)
	case AttrExpression:
		fmt.Println(ex.AttributePath.AttributeName)
		fmt.Println(ex.AttributePath.SubAttribute)
		fmt.Println(ex.CompareOperator, scim2.EQ)
		fmt.Println(ex.CompareValue)
	case scim2.ValuePath:
		fmt.Println(ex.AttributeName)
		parse(ex.ValueExpression)
	}
}

func TestStuff(t *testing.T) {
	b := bytes.NewBufferString(`(stuff eq "stuff") or (stuff eq "another")`)
	p := scim2.NewParser(b)
	e, err := p.Parse()
	require.NoError(t, err)
	fmt.Printf("%#v\n", e)
	parse(e)
}

func TestStuff2(t *testing.T) {
	b := bytes.NewBufferString(`name.stuff eq "another"`)
	p := scim2.NewParser(b)
	e, err := p.Parse()
	require.NoError(t, err)
	fmt.Printf("%#v\n", e)
	parse(e)
}

func TestSub(t *testing.T) {
	b := bytes.NewBufferString(`families[id eq "1"]`)
	p := scim2.NewParser(b)
	e, err := p.Parse()
	require.NoError(t, err)
	fmt.Printf("%#v\n", e)
	parse(e)
}
