测试网合约部署的交易 hash
"transactionId": "1a56f39aefa3bc67e48276c6874bc2abd29468408bc65c8fb0537c3182d69101"

调用合约的交易 hash
"transactionId": "29c021579a4449fe9d930b72f108764968eb04f23a919aac8b9a7e94f3e017e3"

合约的源代码
```
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
```