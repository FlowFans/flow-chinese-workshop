## 测试网合约部署的交易 hash

474ac76caf98fdeb631ccdb5be020c2ef3a7817e8df4ec8e821ca22c00bc191f

## 调用合约的交易 hash

0595b59b38fa6ba0051b9a40081f5702d77e49e81dd0769a682fe3d5c5d1d2db

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