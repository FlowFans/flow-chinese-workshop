from flow_py_sdk import Script, flow_client
import asyncio
import json



async def get_nbatopshot_transaction_events():
    async with flow_client(
        host="access.mainnet.nodes.onflow.org", port=9000
    ) as client:
        block = await client.get_latest_block(is_sealed=True)
        startBlock = block.height-20
        endBlock = block.height
        event_type = "A.c1e4f4f4c4257510.Market.MomentPurchased"
        print(f"当前区块高度: {block.height}\n")
        print(f"下面查询高度{startBlock}至{endBlock}的MomentPurchased交易: \n")
        block_events = await client.get_events_for_height_range(
            type=event_type,
            start_height=startBlock,
            end_height=endBlock
        )
        for block_event in block_events:
            for event in block_event.events:
                fields = json.loads(event.payload)['value']['fields']
                id = fields[0]['value']['value']
                price = fields[1]['value']['value']
                seller = fields[2]['value']['value']['value']
                print(f"type: {event.type}")
                print(f"transaction_id: {event.transaction_id.hex()}")
                print(f"id: {int(id)}")
                print(f"price: {float(price)}")
                print(f"seller: {seller}")
                print('\n')

if __name__ == "__main__":
    asyncio.run(get_nbatopshot_transaction_events())
