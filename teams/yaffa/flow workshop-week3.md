flow workshop-week3



通过Flow CLI 本地开发测试合约， 部署 `Fungible Tokens` 或 `Non-Fungible Tokens` 合约部署至测试网，可以添加新的代码和交互逻辑，将自部署的地址提交到 github 仓库里



https://flow-view-source.com/testnet/account/0xc34e97fc542fbbad





flow accounts add-contract ExampleToken D:\Develop\flow\ExampleToken.cdc --signer testnet-account --network testnet
Contract 'ExampleToken' deployed to the account 'c34e97fc542fbbad'.

Address  0xc34e97fc542fbbad
Balance  1999.99950000     
Keys     1

Key 0   Public Key               48933d54d4462e232bbd100adcd7dc3ab2da2be5e5c2978e9f3926534d09d7f920247c9a2721c43a6da8213b98b75c9a0d8c058b582f1660e8ed490f4c7268dd
        Weight                   1000
        Hash Algorithm           SHA3_256
        Revoked                  false
        Sequence Number          15
        Index                    0

Contracts Deployed: 5
Contract: 'HelloWorld'
Contract: 'Kibble'
Contract: 'KittyItems'
Contract: 'KittyItemsMarket'
Contract: 'ExampleToken'




PS D:\Develop\flow> flow accounts add-contract ExampleNFT D:\Develop\flow\ExampleNFT.cdc --signer testnet-account --network testnet
Contract 'ExampleNFT' deployed to the account 'c34e97fc542fbbad'.

Address  0xc34e97fc542fbbad
Balance  1999.99940000
Keys     1

Key 0   Public Key               48933d54d4462e232bbd100adcd7dc3ab2da2be5e5c2978e9f3926534d09d7f920247c9a2721c43a6da8213b98b75c9a0d8c058b582f1660e8ed490f4c7268dd
        Weight                   1000
        Signature Algorithm      ECDSA_P256
        Hash Algorithm           SHA3_256
        Revoked                  false
        Sequence Number          16
        Index                    0

Contracts Deployed: 6
Contract: 'ExampleNFT'
Contract: 'ExampleToken'
Contract: 'HelloWorld'
Contract: 'Kibble'
Contract: 'KittyItems'
Contract: 'KittyItemsMarket'