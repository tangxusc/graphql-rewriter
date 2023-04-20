package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	"github.com/sirupsen/logrus"
	"net/http"
)

func StartWeb(ctx context.Context) error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	graphQLHandler := graphqlws.NewHandlerFunc(&subscribeService{}, &handler{}, graphqlws.WithWriteTimeout(graphqlRequestTimeout))
	r.POST("/graphql", gin.WrapF(graphQLHandler))
	r.GET("/graphql", gin.WrapF(graphQLHandler))

	srv := &http.Server{
		Addr:    webAddr,
		Handler: r,
	}
	go func() {
		select {
		case <-ctx.Done():
			if err := srv.Shutdown(ctx); err != nil {
				logrus.Errorf("[web]Shutdown server:%+v", err)
			}
		}
	}()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Errorf("[web]ListenAndServe error:%v", err)
		return err
	}
	return nil
}
