package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/formatter"
	"github.com/vektah/gqlparser/parser"
	"testing"
)

func TestReWriteQl(t *testing.T) {
	viper.Set("rewriteql_regexps", map[string]string{
		`^query(\w+)$`: `queryAbcd_$1`,
	})
	var err error
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
	if err = ReWriteQl(query); err != nil {
		t.Fatal(err.Error())
	}

	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatQueryDocument(query)
	fmt.Println(buf.String())
}
