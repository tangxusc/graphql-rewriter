package web

import (
	"context"
	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/graphql-rewriter/package/rewriter"
)

type subscribeService struct {
	subscribeId string
}

func (s *subscribeService) Subscribe(ctx context.Context, document string, operationName string,
	variableValues map[string]interface{}) (<-chan interface{}, error) {
	rewriteGraphql, err := rewriter.RewriteGraphql(document)
	if err != nil {
		return nil, err
	}
	c := make(chan interface{})
	onError := func(sc *graphql.SubscriptionClient, err error) error {
		logrus.Errorf("[remote-client]subscribe error:%v", err)
		//sc.Unsubscribe(s.subscribeId)
		sc.Close()
		close(c)
		return err
	}
	client := graphql.
		NewSubscriptionClient(grahpqlUrlWs).
		WithProtocol(graphql.SubscriptionProtocolType(subscriptionProtocolType)).
		WithConnectionParams(map[string]interface{}{
			"headers": ctx.Value("Header"),
		}).
		WithLog(func(args ...interface{}) {
			logrus.Debugf("[remote-client]:%v\n", args)
		}).
		OnError(onError)

	go func() {
		select {
		case <-ctx.Done():
			logrus.Debugf("[remote-client]Unsubscribe Id:%s", s.subscribeId)
			if c != nil {
				close(c)
			}
			client.Unsubscribe(s.subscribeId)
			client.Close()
		}
	}()
	//pull request : https://github.com/hasura/go-graphql-client/pull/88
	//reason: DEBU[0014] [remote-client]:[{"id":"70e0f2e6-abd2-4e0e-9e85-60ab21a6d768","type":"error","payload":{"message":"input:3: Cannot query field \"productID1\" on type \"Product\". Did you mean \"productID\"?\n"}} server]
	h := func(message []byte, err error) error {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("[remote-client]subscribe client recever panic at Error:%v", err)
			}
		}()
		if err != nil {
			logrus.Debugf("[remote-client]subscribe client recever error:%v", err)
			c <- err.Error()
			close(c)
			c = nil
			return err
		}
		logrus.Debugf("[remote-client]subscribe client recever message:%s", string(message))
		result, err := rewriter.RewriteResult(message)
		if err != nil {
			logrus.Errorf("[remote-client]🌶RewriteResult Error:%v", err)
			c <- formatError(err)
			return err
		}
		c <- formatData(result)
		return nil
	}
	s.subscribeId, err = client.Exec(rewriteGraphql, variableValues, h)
	if err != nil {
		logrus.Errorf("[remote-client]subscribe client Exec error:%v", err)
		return c, err
	}

	go func() {
		if err = client.Run(); err != nil {
			logrus.Errorf("[remote-client]subscribe client Run error:%v", err)
		}
	}()

	return c, err
}
