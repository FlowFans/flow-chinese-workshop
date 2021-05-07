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

 startBlock, endBlock := blockTip-20, blockTip
 const eventType = "A.c1e4f4f4c4257510.Market.MomentPurchased"
 results, err := flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
  Type:        eventType,
  StartHeight: startBlock,
  EndHeight:   endBlock,
 })
 if err != nil{
  panic(err)
 }

 fmt.Println("Query result of NBA Topshot event: ",eventType" from block ",startBlock, " to ",endBlock)
 for _, block := range results {
  for i, event := range block.Events {
   fmt.Println("Found event ",i+1," in block ", block.Height)
   fmt.Println("Transaction ID: ", event.TransactionID)
   fmt.Println("Event ID: ", event.ID())
   fmt.Println("Moment ID: ", event.Value.Fields[0])
   fmt.Println("Price(int): ", event.Value.Fields[1].ToGoValue())
   fmt.Println("Price(str): ", event.Value.Fields[1].String())
   fmt.Println("Seller: ", event.Value.Fields[2])
   fmt.Println(event.Value)
  }
 }
}
