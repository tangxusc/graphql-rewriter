package web

import (
	"fmt"
	"github.com/hasura/go-graphql-client"
	"log"
	"strings"
	"testing"
)

func TestRemoteSubscribe(t *testing.T) {
	client := graphql.NewSubscriptionClient(`ws://192.168.31.12:8091/graphql`).
		WithProtocol(graphql.SubscriptionsTransportWS).
		WithConnectionParams(map[string]interface{}{
			//"headers": map[string]string{
			//	xHasuraAdminSecret: adminSecret,
			//},
		}).WithLog(log.Println).
		OnError(func(sc *graphql.SubscriptionClient, err error) error {
			if strings.Contains(err.Error(), "invalid x-hasura-admin-secret/x-hasura-access-key") {
				return err
			}
			return nil
		})

	defer client.Close()

	client.Exec(`subscription queryProduct{
    queryProduct{
        productID
        name
    }
}`, nil, func(message []byte, err error) error {
		ms := string(message)
		fmt.Println(ms, err)
		return nil
	})

	client.Run()
}
