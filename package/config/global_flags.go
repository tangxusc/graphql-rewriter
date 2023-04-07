package config

import (
	"github.com/spf13/cobra"
)

var WebPort string

var Debug = false

func InitFlags() {
	RegisterFlags(func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringVar(&WebPort, "port", "8080", "web server port")
		//enable Debug
		cmd.PersistentFlags().BoolVar(&Debug, "Debug", false, "enable Debug")
	})
}
