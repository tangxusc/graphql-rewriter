package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tangxusc/graphql-rewriter/package/config"
	"github.com/tangxusc/graphql-rewriter/package/plugin_manager"
	"github.com/tangxusc/graphql-rewriter/package/web"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func NewCommand() (*cobra.Command, context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	command := &cobra.Command{
		Use:   ``,
		Short: ``,
		Long:  ``,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Kill)
				<-c
				cancelFunc()
			}()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				logrus.Debugf("Closing application...")
				cancelFunc()
			}()
			logrus.SetLevel(logrus.TraceLevel)
			rand.Seed(time.Now().UnixNano())
			if err := plugin_manager.StartPlugins(ctx); err != nil {
				return err
			}
			if err := web.StartWeb(ctx); err != nil {
				return err
			}

			return nil
		},
	}

	plugin_manager.InitFlags()
	web.InitFlags()
	config.InitFlags()
	viper.AutomaticEnv()
	viper.AddConfigPath(`./config`)
	viper.SetConfigName("config")
	config.BuildFlags(command)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("[config]read config error:%v", err)
	}
	if err := viper.BindPFlags(command.PersistentFlags()); err != nil {
		logrus.Errorf("[config]BindPFlags config error:%v", err)
	}

	return command, ctx, cancelFunc
}

func main() {
	command, _, _ := NewCommand()
	if err := command.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
