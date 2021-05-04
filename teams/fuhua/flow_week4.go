package main

import (
	"context"
	"fmt"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func HandleErr(err error) {
	if err != nil {
		fmt.Println("err:", err.Error())
		panic(err)
	}
}

func main() {
	ctx := context.Background()

	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	HandleErr(err)

	blockHeader, err := flowClient.GetLatestBlockHeader(ctx, true)
	HandleErr(err)

	blockTip := blockHeader.Height
	fmt.Printf("Current block tip: %d\n", blockHeader.Height)

	// Query by type in latest 10 blocks
	startBlock, endBlock := blockTip-10, blockTip
	const eventType = "A.c1e4f4f4c4257510.Market.MomentPurchased"
	results, err := flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
		Type:        eventType,
		StartHeight: startBlock,
		EndHeight:   endBlock,
	})
	HandleErr(err)

	fmt.Printf("\nQuery result of NBA Topshot event: %s from block %d to %d\n", eventType, startBlock, endBlock)
	for _, block := range results {
		for i, event := range block.Events {
			fmt.Printf("Found event #%d in block #%d\n", i+1, block.Height)
			fmt.Printf("Transaction ID: %s\n", event.TransactionID)
			fmt.Printf("Event ID: %s\n", event.ID())
			fmt.Printf("Moment ID: %d\n", event.Value.Fields[0])
			fmt.Printf("Price(int): %d\n", event.Value.Fields[1].ToGoValue())
			fmt.Printf("Price(str): %s\n", event.Value.Fields[1].String())
			fmt.Printf("Seller: %s\n", event.Value.Fields[2])
			fmt.Printf("%s\n\n",event.Value)
		}
	}
}
