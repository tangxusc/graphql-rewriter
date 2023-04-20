package rewriter

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/graphql-rewriter/package/plugin_manager"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/formatter"
	"github.com/vektah/gqlparser/parser"
)

func RewriteGraphql(q string) (result string, err error) {
	query, err := parser.ParseQuery(&ast.Source{Input: q, Name: "spec"})
	if err != nil && err.Error() != "" {
		return
	}

	err = plugin_manager.ReWriteQl(query)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatQueryDocument(query)

	logrus.Debugf("[rewrite]rewrite after:\n%s\n", buf.String())

	return buf.String(), err
}

func RewriteResult(r []byte) (result []byte, err error) {
	result, err = plugin_manager.ReWriteResult(r)
	return
}
