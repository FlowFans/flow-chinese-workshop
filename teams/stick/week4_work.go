package main

import (
	"context"
	"fmt"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)
func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	handleErr(err)
	err = flowClient.Ping(context.Background())
	handleErr(err)
	latestBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	handleErr(err)
	fmt.Println("current height: ", latestBlock.Height)
	blockEvents, err := flowClient.GetEventsForHeightRange(context.Background(), client.EventRangeQuery{
		Type:        "A.c1e4f4f4c4257510.Market.MomentPurchased",
		StartHeight: latestBlock.Height-3,
		EndHeight:   latestBlock.Height,
	})
	handleErr(err)

	for _, blockEvent := range blockEvents {
		for _, purchaseEvent := range blockEvent.Events {
			fmt.Println("=====================================")
			fmt.Println(purchaseEvent)
			fmt.Println(purchaseEvent.Value)
			fmt.Printf("Transaction ID: %s\n", purchaseEvent.TransactionID)
			fmt.Printf("Event ID: %s\n", purchaseEvent.ID())
			fmt.Printf("Moment ID: %d\n", purchaseEvent.Value.Fields[0])
			fmt.Printf("Price(str): %s\n", purchaseEvent.Value.Fields[1].String())
			fmt.Printf("Seller: %s\n", purchaseEvent.Value.Fields[2])
		}
	}
}
