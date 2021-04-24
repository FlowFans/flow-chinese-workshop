## 测试网合约部署的交易 hash

f0c10e029e30ea260bfa69c88036b5354660f2141ee5f84242bd5d37612fe136

## 调用合约的交易 hash

76b8c4ff88ac580336aa8ecaa4def64a83ffce9e8ea9c634ce33307091046ee5

## 合约的源代码

```ts
access(all) contract HelloWorld {
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "Hello, Caos!"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}
```
