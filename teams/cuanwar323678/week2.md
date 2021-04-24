# 测试网合约部署的交易 hash

tx hash: d1e4ea2e73bcaae4292404b743ec576b4938649d5c6c53ae830432a80d8ab19d
contract address: 0x33b75fa9399ff5b9

# 调用合约的交易 hash

tx hash: 204f72d0bc52c06789cfa1be1333af53725b86153d8089aad8f689f3a0fce784

# 合约的源代码

access(all) contract helloworld{
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "Hello, world"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}