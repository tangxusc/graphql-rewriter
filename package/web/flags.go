package web

import (
	"github.com/spf13/cobra"
	"github.com/tangxusc/graphql-rewriter/package/config"
)

var webAddr string

func InitFlags() {
	config.RegisterFlags(func(c *cobra.Command) {
		c.PersistentFlags().StringVar(&webAddr, "webAddr", ":8080", "web addr")
	})
}
