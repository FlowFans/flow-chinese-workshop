## 測試網合约部署的交易 hash
3686f1a085281ea76701c440d46fd5c2cfdc4cca3222746b6b037fa94bd2da73

## 調用合约的交易 hash
8f4cde46061b8692bacaa0632fc4a8fb9eb961c097abe11b689dc3f5694a7202

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
