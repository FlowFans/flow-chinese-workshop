package main

import (
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func main() {
	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure()) // 主網
	if err != nil {
		panic(err)
	}

	latestBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	if err != nil {
		panic(err)
	}
	fmt.Println("current height: ", latestBlock.Height)

	// NBA TOP SHOT EVENT
	blockEvents, err := flowClient.GetEventsForHeightRange(context.Background(), client.EventRangeQuery{
		Type:        "A.c1e4f4f4c4257510.Market.MomentListed",
		StartHeight: latestBlock.Height - 249, //主網最多只能輪詢 250 個區塊
		EndHeight:   latestBlock.Height,
	})
	if err != nil {
		panic(err)
	}

	for _, blockEvent := range blockEvents {
		for _, event := range blockEvent.Events {
			fmt.Println("BlockID: ", blockEvent.Height)
			fmt.Println("TransactionID: ", event.TransactionID)
			fmt.Println("event.Value: ", event.Value)
			fmt.Println("id: ", event.Value.Fields[0])
			fmt.Println("price: ", event.Value.Fields[1])
			fmt.Println("seller: ", event.Value.Fields[2])
		}
	}
}
