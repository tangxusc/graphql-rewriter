package web

import (
	"github.com/hasura/go-graphql-client"
	"github.com/spf13/cobra"
	"github.com/tangxusc/graphql-rewriter/package/config"
	"time"
)

var webAddr string
var grahpqlUrl = ":8080"
var grahpqlUrlWs = ":8080"
var graphqlRequestTimeout = time.Second * 5
var subscriptionProtocolType = string(graphql.SubscriptionsTransportWS)

func InitFlags() {
	config.RegisterFlags(func(c *cobra.Command) {
		c.PersistentFlags().StringVar(&webAddr, "web-addr", ":8080", "web addr")
		c.PersistentFlags().StringVar(&grahpqlUrl, "proxy-graphql-url", "http://192.168.31.12:8091/graphql", "GraphQL server address")
		c.PersistentFlags().StringVar(&grahpqlUrlWs, "proxy-graphql-url-ws", "ws://192.168.31.12:8091/graphql", "GraphQL server address")
		c.PersistentFlags().StringVar(&subscriptionProtocolType, "subscription-protocol-type", "subscriptions-transport-ws",
			"default subscriptions-transport-ws,option: subscriptions-transport-ws,graphql-ws")
		//http request timeout duration
		c.PersistentFlags().DurationVar(&graphqlRequestTimeout, "proxy-graphql-request-timeout", time.Second*5, "GraphQL request Timeout")

	})
}
