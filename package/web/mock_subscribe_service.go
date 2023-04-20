package web

import (
	"context"
	"time"
)

type mockSubscribeService struct {
}

func (m *mockSubscribeService) Subscribe(ctx context.Context, document string, operationName string,
	variableValues map[string]interface{}) (<-chan interface{}, error) {
	c := make(chan interface{})
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(c)
				return
			case c <- time.Now().String():
				time.Sleep(time.Second)
			}
		}
	}()
	return c, nil
}
