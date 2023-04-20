package rewriter

import (
	"github.com/spf13/cobra"
	"github.com/tangxusc/graphql-rewriter/package/config"
)

func InitFlags() {
	config.RegisterFlags(func(c *cobra.Command) {
	})
}
