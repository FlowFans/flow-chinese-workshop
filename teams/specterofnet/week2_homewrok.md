## 测试网合约部署的交易 hash

ef703e8e8e1e3c5e85ea8fd9273000569165fd6d8e2f5ec50fe6a2e34faa1450

## 调用合约的交易 hash

44c8da86e3d7f1a43684bc4de57e63b5c129a0ed5534557e7258b6903faae310

## 合约的源代码

```ts

access(all) contract FirstContractOnFlow {
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "Hello, specter"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}
```
