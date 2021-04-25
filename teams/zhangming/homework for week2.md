## 测试网合约部署的交易 hash

d0461ccbff9f538d6e607ecfff83fcbed6a0224f7378033223d874ea3f3d5877

## 调用合约的交易 hash

de67af8740b620400c4b1ff4d0bff08788230e4b825dda72d98af8a39b6b0b71

## 合约的源代码

```ts
access(all) contract test1{
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "test for week2"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}
```