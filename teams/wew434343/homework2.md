# 测试网合约部署的交易 hash

contract address: 0xce5eeec1a233e912
tx hash: ebff575afbf0e29dc0fe15603a2a1f2e18be4a65f39f11e878e134b1e7215b22

# 调用合约的交易 hash

tx hash: 12b4787399d66588ab1d076039b385e7eef2c6834fb0519125dd220ccbead321

# 合约的源代码

access(all) contract contract1{
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "first contract"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}