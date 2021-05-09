
package main
import (
	"context"
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func main() {
	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	blockEvents, err := flowClient.GetEventsForHeightRange(context.Background(), client.EventRangeQuery{
		Type:        "A.7e60df042a9c0868.FlowToken.TokensDeposited",
		StartHeight: 14305951,
		EndHeight:   14305961,
	})

	for _, blockEvent := range blockEvents {
		for _, event := range blockEvent.Events {
			amount := event.Value.Fields[0].(cadence.UFix64).ToGoValue().(uint64)
			address := event.Value.Fields[1].(cadence.Optional).Value.(cadence.Address).String()
			fmt.Println("amount--- ", amount)
			fmt.Println("address -- ", address)
		}
	}
}