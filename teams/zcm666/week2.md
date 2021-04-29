## 测试网合约部署的交易 hash

1e40d2ed7577cd2ad296af3cc9085d72382dd4126521a02fb9f8bdf5d11da3af

## 调用合约的交易 hash

c0e17fcf5bfb6b1185e553913d6bea135c9d43b323a65a1d3ed33400a888c0e5

## 合约的源代码

```ts

access(all) contract flowContractTest {
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "Hello, zcm"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}

```
