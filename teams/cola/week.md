## 测试网合约部署的交易 hash

561f0a6204337d12dceea85c0956d0fe9be54ce417fd913448a896d07b0e9610

## 调用合约的交易 hash

f7a8bd01fff839f632a9657544905716857d13acddefa378cb5dde2cc6c485f2

## 合约的源代码

```ts
access(all) contract HelloWorld {
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
