package rewriter

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tangxusc/graphql-rewriter/package/plugins"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
	"net/http"
)

func RewriteGraphql(c *gin.Context, body []byte) error {
	var err error
	query, err := parser.ParseQuery(&ast.Source{Input: string(body), Name: "spec"})
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(query)
	if err != nil {
		return err
	}

	patchs, err := plugins.Patchs(marshal)
	if err != nil {
		return err
	}
	out := marshal
	for _, patch := range patchs {
		out, err = patch.Apply(out)
		if err != nil {
			return err
		}
	}
	var doc ast.QueryDocument
	err = json.Unmarshal(out, &doc)
	if err != nil {
		return err
	}
	//TODO:转发请求至graphql服务器
	c.String(http.StatusOK, string(out))
	return nil
}
