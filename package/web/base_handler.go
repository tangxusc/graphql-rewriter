package web

import (
	"context"
	"encoding/json"
	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/graphql-rewriter/package/rewriter"
	"github.com/tidwall/sjson"
	"io"
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
	w.Header().Set("Content-Type", "application/json")
	var params paramPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logrus.Errorf("[web]decode request param error:%v", err)
		writeError(w, err)
		return
	}

	rewriteGraphql, err := rewriter.RewriteGraphql(params.Query)
	if err != nil {
		logrus.Errorf("[web]Rewrite Graphql error:%v", err)
		writeError(w, err)
		return
	}

	httpClient := &http.Client{Timeout: graphqlRequestTimeout}
	client := graphql.NewClient(grahpqlUrl, httpClient)
	raw, err := client.ExecRaw(context.TODO(), rewriteGraphql, params.Variables)
	if err != nil {
		logrus.Errorf("[web]exec raw graphql for server: %v , error:%v", grahpqlUrl, err)
		writeError(w, err)
		return
	}

	result, err := rewriter.RewriteResult(raw)
	if err != nil {
		logrus.Errorf("[web]Rewrite Result error:%v", err)
		writeError(w, err)
		return
	}

	writeData(w, result)
}

func writeError(w io.Writer, err error) {
	_, _ = w.Write(formatError(err))
}

func formatError(err error) []byte {
	value, err := sjson.SetBytes([]byte(`{"errors":[{"message":""}]}`), "errors.0.message", err.Error())
	if err != nil {
		logrus.Errorf("[web]SetBytes error:%v", err)
	}
	return value
}

type resultPayload struct {
	Data json.RawMessage `json:"data"`
}

func writeData(w io.Writer, data []byte) {
	_, _ = w.Write(formatData(data))
}

func formatData(data []byte) json.RawMessage {
	payload := resultPayload{Data: data}
	marshal, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("[web] marshal result payload error:%v", err)
		return nil
	}
	return marshal
}
