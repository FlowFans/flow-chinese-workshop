import sdk from '@onflow/sdk'

const contractAddress = 'c1e4f4f4c4257510';
const contractName = 'Market';
const eventName = 'MomentPurchased';

const getEvents = async () => {
  sdk.config().put('accessNode.api', 'http://access.mainnet.nodes.onflow.org:9000');

  let query = await sdk.build([
    sdk.getBlock(true)
  ]);
  const pipedQuery = await sdk.pipe(query)
  const latestBlockResponse = await sdk.send(pipedQuery)
  const latestBlock = await sdk.decode(latestBlockResponse);
  const toBlock = latestBlock.height;
  const fromBlock = toBlock - 20;
  const eventType = `A.${contractAddress}.${contractName}.${eventName}`;
  const eventsResponse = await sdk.send(await sdk.build([ sdk.getEventsAtBlockHeightRange(eventType, fromBlock, toBlock) ]));
  return await sdk.decode(eventsResponse);
}

const main = async () => {
  const event = await getEvents()
  console.log(event)
}

main()