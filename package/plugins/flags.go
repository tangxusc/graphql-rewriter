package plugins

import (
	"github.com/spf13/cobra"
	"github.com/tangxusc/graphql-rewriter/package/config"
)

var pluginDir string

func InitFlags() {
	config.RegisterFlags(func(c *cobra.Command) {
		c.PersistentFlags().StringVar(&pluginDir, "plugins", "./plugins", "wasm plugins")
	})
}
