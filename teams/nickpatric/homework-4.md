# Flow Go-SDK

- Go module file 
- 
```go
module flow-grants

go 1.16

require (
	github.com/onflow/cadence v0.15.0 // indirect
	github.com/onflow/flow-go-sdk v0.19.0 // indirect
	google.golang.org/grpc v1.37.1 // indirect
)


```

- Go-sdk demo codes

```go

package main

import (
	"context"
	"fmt"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func main(){
	flowClient, err := client.New("access.devnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	blockEvents, err := flowClient.GetEventsForHeightRange(context.Background(), client.EventRangeQuery{
		Type:        "A.7e60df042a9c0868.FlowToken.TokensDeposited",
		StartHeight: 31267819,
		EndHeight:   31267829,
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

```