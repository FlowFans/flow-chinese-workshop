#### 部署合約交易Hash
0c7b9abec9e54a6336f47ea65996d1709929642496d87d17f219cca11d8d8147

#### 調用合約交易Hash
fe7fbb8dc6cac9f2fc1bd12b3cf00281886e53d7b49dc4c1538e0fadd3b90e50

#### 測試網合約部署交易hash
```
pub contract Workshop3 {
  pub let greeting: String
  pub event HelloEvent(message: String)

  pub resource NFT {
    pub let id:UInt64
    pub var metadata: {String: String}
    init(initID: UInt64) {
      self.id = initID
      self.metadata = {}
    }
    pub fun nftHello(): UInt64 {
      return self.id
    }
  }

  init() {
    self.greeting = "Hello, Workshop3!"
    self.account.save<@NFT>(<-create NFT(initID: 1), to: /storage/NFTWorkshop3)
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}
```