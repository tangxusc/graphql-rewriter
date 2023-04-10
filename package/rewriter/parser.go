package rewriter

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/graphql-rewriter/package/plugin_manager"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/formatter"
	"github.com/vektah/gqlparser/parser"
	"net/http"
)

func RewriteGraphql(c *gin.Context, body []byte) error {
	var err error
	input := string(body)
	query, err := parser.ParseQuery(&ast.Source{Input: input, Name: "spec"})
	if err != nil && err.Error() != "" {
		return err
	}

	err = plugin_manager.ReWrite(query)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatQueryDocument(query)

	logrus.Debugf("rewrite after:%s\n", buf.String())
	//TODO:转发请求至graphql服务器
	c.String(http.StatusOK, string(""))
	return nil
}
