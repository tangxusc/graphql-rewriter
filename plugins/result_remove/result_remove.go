package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tidwall/sjson"
	"github.com/vektah/gqlparser/ast"
)

func ReWriteQl(gql *ast.QueryDocument) error {
	return nil
}

func ReWriteResult(data []byte) ([]byte, error) {
	bytes, err := sjson.DeleteBytes(data, "extensions")
	if err != nil {
		return nil, err
	}
	logrus.Debugf("[plugin][print]ReWriteResult :\n%+v\n", string(bytes))
	return bytes, nil
}

func main() {
	println("print plugin main")
}
