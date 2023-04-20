package web

import (
	"context"
	"encoding/json"
	"github.com/hasura/go-graphql-client"
	"github.com/tangxusc/graphql-rewriter/package/rewriter"
	"net/http"
)

type handler struct {
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		//TODO: "message":"error"
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rewriteGraphql, err := rewriter.RewriteGraphql(params.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpClient := &http.Client{Timeout: graphqlRequestTimeout}
	client := graphql.NewClient(grahpqlUrl, httpClient)
	raw, err := client.ExecRaw(context.TODO(), rewriteGraphql, params.Variables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := rewriter.RewriteResult(raw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
