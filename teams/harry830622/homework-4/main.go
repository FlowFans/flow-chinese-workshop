package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/onflow/cadence"
	// "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	// "github.com/onflow/flow-go-sdk/crypto"
	// "github.com/onflow/flow-go-sdk/templates"
)

func main() {
	ctx := context.Background()
	// flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	block, err := flowClient.GetLatestBlock(ctx, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(block.Height)

	// tx, err := flowClient.GetTransaction(ctx, flow.HexToID("f9707a1a98a156fd9dc7d0cac84c95db6b013f3999789577a6b476a8d83a97f2"))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(tx.Authorizers)

	momentListedEvtType := "A.c1e4f4f4c4257510.Market.MomentListed"
	// momentPriceChangedEvtType := "A.c1e4f4f4c4257510.Market.MomentPriceChanged"
	// momentMintedEvtType := "A.0b2a3299cc857e29.TopShot.MomentMinted"

	q := client.EventRangeQuery{Type: momentListedEvtType, StartHeight: block.Height - 0*60*60 - 1, EndHeight: block.Height - 0*60*60}
	blockEvts, err := flowClient.GetEventsForHeightRange(ctx, q)
	if err != nil {
		panic(err)
	}

	for _, blockEvt := range blockEvts {
		for _, evt := range blockEvt.Events {
			fmt.Println(evt.Type)
			fmt.Println(evt.Value.Fields)
			script := `
import TopShot from 0x0b2a3299cc857e29
import Market from 0xc1e4f4f4c4257510
pub struct Moment {
  pub var id: UInt64
  pub var playId: UInt32
  pub var setId: UInt32
  pub var serialNumber: UInt32
  init(moment: &TopShot.NFT) {
    self.id = moment.id
    self.playId = moment.data.playID
    self.setId = moment.data.setID
    self.serialNumber = moment.data.serialNumber
  }
}
pub fun main(owner: Address, momentId: UInt64): Moment {
  let acct = getAccount(owner)
  let collectionRef = acct.getCapability(/public/topshotSaleCollection)!.borrow<&{Market.SalePublic}>() ?? panic("Could not borrow capability from public collection")
  let moment = collectionRef.borrowMoment(id: momentId) ?? panic("Could not borrow moment from public collection")
  return Moment(moment: moment)
}
`
			res, err := flowClient.ExecuteScriptAtBlockHeight(ctx, block.Height, []byte(script), []cadence.Value{evt.Value.Fields[2].(cadence.Optional).Value, evt.Value.Fields[0]})
			if err != nil {
				panic(err)
			}
			fmt.Println(res)

			script = `
import TopShot from 0x0b2a3299cc857e29
pub struct MomentData {
  pub var playMetaData: {String: String}?
  pub var setName: String?
  init(playId: UInt32, setId: UInt32) {
    self.playMetaData = TopShot.getPlayMetaData(playID: playId)
    self.setName = TopShot.getSetName(setID: setId)
  }
}
pub fun main(playId: UInt32, setId: UInt32): MomentData {
  return MomentData(playId: playId, setId: setId)
}
`
			res, err = flowClient.ExecuteScriptAtBlockHeight(ctx, block.Height, []byte(script), []cadence.Value{res.(cadence.Struct).Fields[1], res.(cadence.Struct).Fields[2]})
			if err != nil {
				panic(err)
			}
			fmt.Println(res)
		}
	}
}
