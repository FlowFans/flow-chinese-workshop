1 测试网合约部署的交易 hash
9949250d8c2ae905b97af13ec95e08b515999bbed45e979551bb7faad5b3d0a9

2 调用合约的交易 hash
9bcd1f5acbbadd615b95d4bcbd9f9861bccdd6d00646b3ad9647cfd0c4f9c825

3 合约的源代码
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
