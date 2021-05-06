from flow_py_sdk import Script, flow_client
import json
import asyncio


async def get_nba_transaction_events():
    async with flow_client(
        host="access.mainnet.nodes.onflow.org", port=9000
    ) as client:
        block = await client.get_latest_block(is_sealed=True)
        print(f"block.height: {block.height}\n")
        startBlock, endBlock = block.height-10, block.height
        event_type = "A.c1e4f4f4c4257510.Market.MomentPurchased"
        block_events = await client.get_events_for_height_range(
            type=event_type,
            start_height=startBlock,
            end_height=endBlock
        )
        for block_event in block_events:
            for event in block_event.events:
                print(f"type: {event.type}")
                print(f"transaction_id: {event.transaction_id.hex()}")
                fields = json.loads(event.payload)['value']['fields']
                id = fields[0]['value']['value']
                print(f"id: {int(id)}")
                price = fields[1]['value']['value']
                print(f"price: {float(price)}")
                seller = fields[2]['value']['value']['value']
                print(f"seller: {seller}")
                print('\n')

if __name__ == "__main__":
    asyncio.run(get_nba_transaction_events())
