## 测试网合约部署的交易 hash
1b0e44b44e16c4568b0f3cdfedf91e56d09b18c85a61095ea75af1967a088069

## 调用合约的交易 hash
c45eca8132f71a24797f76d3e2550d9505e3b5886fd9ef430540a8783e405190

## 合约的源代码

```ts
access(all) contract HelloWorld {
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "Hello, 423!"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}
```