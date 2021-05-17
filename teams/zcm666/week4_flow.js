import * as sdk from "@onflow/sdk"

const node = "https://access-mainnet-beta.onflow.org"
const responseForGetLatestBlock = await sdk.send(await sdk.build([
    sdk.getLatestBlock()
]), { node: node })

let lastheight = responseForGetLatestBlock.block.height
let eventType = "A.c1e4f4f4c4257510.Market.MomentPurchased"

const responseForGetEvents = await sdk.send(await sdk.build([
    sdk.getEvents(eventType, lastheight, lastheight),
]), { node: node })


let txEvents = responseForGetEvents.events
txEvents.forEach(tx => {

    let payload = tx.payload
    let id_fields = payload.value.fields[0]
    let price_fields = payload.value.fields[1]
    let seller_fields = payload.value.fields[2]
    console.log("transactionId:", tx.transactionId)
    console.log("id:", id_fields.value.value)
    console.log("price:", price_fields.value.value)
    console.log("seller:", seller_fields.value.value.value)
    console.log('\n')
}); 