测试网合约部署的交易tx d131683acfd950205a0d3a54ee3038a5b9bbaab8853d579c7d34afc6a7d575f9
调用合约的交易tx      a96aaa38875eb9525b4a50f5bf7ef883b7624d8d587cc2ebf671c6bf07a3025d
helloworld 代码
```
pub contract HelloWorld {
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "Hello, World!"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}
```