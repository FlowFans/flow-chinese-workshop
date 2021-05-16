package main

import (
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func main() {
	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if error != nil {
		panic(error)
	}

	latestBlock, error := flowClient.GetLatestBlock(context.Background(), false)
	if error != nil {
		panic(error)
	}
	fmt.Println("current height block: ", latestBlock.Height)

	// NBA TOP SHOT EVENT
	blockEvents, error := flowClient.GetEventsForHeightRange(context.Background(), client.EventRangeQuery{
		Type:        "A.c1e4f4f4c4257510.Market.MomentListed",
		StartHeight: latestBlock.Height - 10,
		EndHeight:   latestBlock.Height,
	})
	if error != nil {
		panic(error)
	}

	for _, blockEvent := range blockEvents {
		for _, event := range blockEvent.Events {
			fmt.Println("BlockID: ", blockEvent.Height)
			fmt.Println("TransactionID: ", event.TransactionID)
			fmt.Println("id: ", event.Value.Fields[0])
			fmt.Println("price: ", event.Value.Fields[1])
			fmt.Println("seller: ", event.Value.Fields[2])
		}
	}
}