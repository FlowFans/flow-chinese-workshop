package main
import (
	"context"
	"fmt"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func main() {

	ctx := context.Background()

	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	
	blockHeader, err := flowClient.GetLatestBlockHeader(ctx, true)
	if err != nil{
		panic(err)
	}

	//get last block height
	lastBlockHeight := blockHeader.Height
	fmt.Printf("Current block height: ", blockHeader.Height)

	startBlock, endBlock := lastBlockHeight-1000, lastBlockHeight

	//get tokendeposited event from blocks
	blockEvents, err := flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
		Type:        "A.7e60df042a9c0868.FlowToken.TokensDeposited",
		StartHeight: startBlock,
		EndHeight:   endBlock,
	})

	//print blockevent
	for _, blockEvent := range blockEvents {
		for _, event := range blockEvent.Events {
			amount := event.Value.Fields[0].(cadence.UFix64).ToGoValue().(uint64)
			address := event.Value.Fields[1].(cadence.Optional).Value.(cadence.Address).String()
			fmt.Println("amount: ", amount)
			fmt.Println("address: ", address)
		}
	}
}