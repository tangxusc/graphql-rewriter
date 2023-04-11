package rewriter

import (
	"github.com/spf13/cobra"
	"github.com/tangxusc/graphql-rewriter/package/config"
	"time"
)

var grahpqlUrl = ":8080"

var graphqlRequestTimeout = time.Second * 5

func InitFlags() {
	config.RegisterFlags(func(c *cobra.Command) {
		c.PersistentFlags().StringVar(&grahpqlUrl, "graphql-url", "http://192.168.31.12:8091/graphql", "GraphQL server address")
		//http request timeout duration
		c.PersistentFlags().DurationVar(&graphqlRequestTimeout, "graphql-request-timeout", time.Second*5, "GraphQL request Timeout")
	})
}
