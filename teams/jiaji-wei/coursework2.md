## 测试网合约部署的交易 hash

abdf9cce811a0ad599894d8eaf474ab9adc8785b9383516422edf95bc1bd4d61

## 调用合约的交易 hash

726c8250420dc18381b9584e47467360c9c229212ada219f5e538b74d2c6a543

## 合约的源代码

```ts
access(all) contract HelloWorld {
  pub let greeting: String
  pub var counter: Int
  pub event HelloEvent(message: String, number: Int)

  init() {
    self.greeting = "Hello, Flow!"
    self.counter = 0
  }

  pub fun hello(message: String): String {
    self.counter = self.counter + 1
    emit HelloEvent(message: message, number: self.counter)
    return self.greeting
  }
}

```