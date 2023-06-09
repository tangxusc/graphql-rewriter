package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/ast"
)

func ReWriteQl(gql *ast.QueryDocument) error {
	logrus.Debugf("[plugin][print]ReWriteQl result:\n%+v\n", gql)
	return nil
}

func ReWriteResult(data []byte) ([]byte, error) {
	logrus.Debugf("[plugin][print]ReWriteResult :\n%+v\n", string(data))
	return data, nil
}

func main() {
	println("print plugin main")
}
