package main

import (
	"bytes"
	"fmt"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/formatter"
	"github.com/vektah/gqlparser/parser"
	"testing"
)

func TestReWriteQl(t *testing.T) {
	input := `query MyQuery {
  queryProduct {
    name
    productID
  }
}`
	query, err := parser.ParseQuery(&ast.Source{Input: input, Name: "spec"})
	if err != nil && err.Error() != "" {
		t.Fatal(err.Error())
	}
	field := query.Operations[0].SelectionSet[0].(*ast.Field)
	//field.Name = "query_test_Product_name"
	field.Alias = "query_test_Product_Alias"
	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatQueryDocument(query)
	fmt.Println(query)
	fmt.Println(buf.String())
}
