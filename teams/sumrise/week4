# -*- coding: utf-8 -*-

"""

@author: sumrise
@time: 2021-05-13 23:07
"""

from flow_py_sdk import flow_client
import json
import asyncio


async def get_nba_transaction_events():
    async with flow_client(
            host="access.mainnet.nodes.onflow.org", port=9000
    ) as client:
        block = await client.get_latest_block(is_sealed=True)
        print(f"block.height: {block.height}\n")
        startBlock, endBlock = block.height - 15, block.height
        event_type = "A.c1e4f4f4c4257510.Market.MomentPurchased"
        print("当前区块高度: {}".format(block.height))
        print("下面查询高度 {} - {} 的MomentPurchased".format(startBlock, endBlock))

        block_events = await client.get_events_for_height_range(
            type=event_type,
            start_height=startBlock,
            end_height=endBlock
        )
        for block_event in block_events:
            for event in block_event.events:
                print("type: {}".format(event.type))
                print("transaction_id: {}".format(event.transaction_id.hex()))

                fields = json.loads(event.payload)['value']['fields']
                id = fields[0]['value']['value']
                price = fields[1]['value']['value']
                seller = fields[2]['value']['value']['value']

                print("id: {}".format(int(id)))
                print("price: {}".format(float(price)))
                print("seller: {}".format(seller))
                print()


if __name__ == "__main__":
    asyncio.run(get_nba_transaction_events())
