
package main

import (
	"context"
	"fmt"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func main() {

	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	processErr(err)
	err = flowClient.Ping(context.Background())
	processErr(err)

	//latestBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	//processErr(err)
	//fmt.Println("current block height: ", latestBlock.Height)

	blockEvents, err := flowClient.GetEventsForHeightRange(context.Background(), client.EventRangeQuery{
		Type:        "A.c1e4f4f4c4257510.Market.MomentListed",
		StartHeight: 14481080,
		EndHeight:   14481089,
	})
	processErr(err)

	for _, blockEvent := range blockEvents {
		for _, purchaseEvent := range blockEvent.Events {
			//fmt.Printf("%s\n\n",purchaseEvent.Value)
			fmt.Printf("Transaction ID: %s\n", purchaseEvent.TransactionID)
			fmt.Printf("Event ID: %s\n", purchaseEvent.ID())
			fmt.Printf("Moment ID: %d\n", purchaseEvent.Value.Fields[0])
			fmt.Printf("Price: %d\n", purchaseEvent.Value.Fields[1])
			fmt.Printf("Seller: %s\n", purchaseEvent.Value.Fields[2])

		}
	}
}

func processErr(err error) {
	if err != nil {
		fmt.Println("Err:" + err.Error())
		panic(err)
	}
}
