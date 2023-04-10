package main

import (
	"fmt"
	"github.com/vektah/gqlparser/ast"
)

func ReWrite(gql *ast.QueryDocument) error {
	fmt.Printf("Before rewrite:\n%+v\n\n", gql)
	gql.Operations[0].Name = "test"
	return nil
}

func main() {
	println("print plugin main")
}
