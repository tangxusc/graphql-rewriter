package rewriter

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/graphql-rewriter/package/plugin_manager"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/formatter"
	"github.com/vektah/gqlparser/parser"
	"gopkg.in/errgo.v2/fmt/errors"
	"io/ioutil"
	"net/http"
)

func RewriteGraphql(c *gin.Context, body []byte) error {
	var err error
	input := string(body)
	query, err := parser.ParseQuery(&ast.Source{Input: input, Name: "spec"})
	if err != nil && err.Error() != "" {
		return err
	}

	err = plugin_manager.ReWriteQl(query)
	if err != nil {
		return err
	}

	//TODO: io.pipe
	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatQueryDocument(query)

	logrus.Debugf("rewrite after:%s\n", buf.String())

	//redirect graphql to graphql server
	//http post
	client := &http.Client{Timeout: graphqlRequestTimeout}
	request, err := http.NewRequest(http.MethodPost, grahpqlUrl, bytes.NewBufferString(buf.String()))
	if err != nil {
		return err
	}
	//set header content-type application/graphql
	request.Header.Set("Content-Type", "application/graphql")
	do, err := client.Do(request)
	if err != nil {
		return err
	}
	defer do.Body.Close()
	//read body
	all, err := ioutil.ReadAll(do.Body)
	if err != nil {
		return err
	}
	//status code 200
	if do.StatusCode != http.StatusOK {
		return errors.Newf("graphql request failed with status code:%d , body:%s ", do.StatusCode, do.Body)
	}
	all, err = plugin_manager.ReWriteResult(all)
	if err != nil {
		return err
	}
	c.Writer.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	_, err = c.Writer.Write(all)
	return err
}
