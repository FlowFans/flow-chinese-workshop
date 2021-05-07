package main

import (
	"context"
	"fmt"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil{
		panic(err)
	}

	blockHeader, err := flowClient.GetLatestBlockHeader(ctx, true)
	if err != nil{
		panic(err)
	}

	blockTip := blockHeader.Height
	fmt.Printf("Current block tip: ", blockHeader.Height)

	startBlock, endBlock := blockTip-40, blockTip
	const eventType = "A.c1e4f4f4c4257510.Market.MomentPurchased"
	results, err := flowClient.GetEventsForHeightRange(ctx, 
		client.EventRangeQuery{
		Type:        eventType,
		StartHeight: startBlock,
		EndHeight:   endBlock,
	})
	if err != nil{
		panic(err)
	}

	fmt.Println("Query result of NBA Topshot event: ",eventType," from block ",startBlock, " to ",endBlock)
	fmt.Println("-----------------------------")
	for _, block := range results {
		for _, event := range block.Events {
			fmt.Println("block ", block.Height)
			fmt.Println("transaction_id: ", event.TransactionID)
			fmt.Println("evetn_id: ", event.ID())
			fmt.Println("moment_id: ", event.Value.Fields[0])
			fmt.Println("price: ", event.Value.Fields[1].String())
			fmt.Println("seller: ", event.Value.Fields[2])
			fmt.Println(event.Value)
			fmt.Println("---------------------------")
		}
	}
}
