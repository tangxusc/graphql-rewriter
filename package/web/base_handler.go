package web

import (
	"context"
	"encoding/json"
	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/graphql-rewriter/package/rewriter"
	"net/http"
)

type paramPayload struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type handler struct {
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params paramPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logrus.Errorf("[web]decode request param error:%v", err)
		//TODO: "message":"error"
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rewriteGraphql, err := rewriter.RewriteGraphql(params.Query)
	if err != nil {
		logrus.Errorf("[web]Rewrite Graphql error:%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpClient := &http.Client{Timeout: graphqlRequestTimeout}
	client := graphql.NewClient(grahpqlUrl, httpClient)
	raw, err := client.ExecRaw(context.TODO(), rewriteGraphql, params.Variables)
	if err != nil {
		logrus.Errorf("[web]exec raw graphql for server: %v , error:%v", grahpqlUrl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := rewriter.RewriteResult(raw)
	if err != nil {
		logrus.Errorf("[web]Rewrite Result error:%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
