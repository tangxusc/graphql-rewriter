package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/graphql-rewriter/package/rewriter"
	"io/ioutil"
	"net/http"
)

func StartWeb(ctx context.Context) error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/graphql", func(c *gin.Context) {
		//get graphql header
		header := c.GetHeader("Content-Type")
		logrus.Infof("header content-type:%s", header)
		//validate header is graphql
		if header != "application/graphql" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "invalid header",
			})
			return
		}
		//read body
		all, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "invalid body",
			})
			return
		}
		defer c.Request.Body.Close()

		if err := rewriter.RewriteGraphql(c, all); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		}
	})

	srv := &http.Server{
		Addr:    ":8080",
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
