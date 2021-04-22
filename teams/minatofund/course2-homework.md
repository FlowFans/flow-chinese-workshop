### 创建合约Hash
> 89497d6ab59a65168a8e728d16a46095c65845123808955241c27a413fb64d3a
### 调用合约Hash
> 8fc557e932d98af822a223955e19a3be39b746e488b3ffe9c00695146599a80b
### 源代码

  ```ts
  access(all) contract HelloWorld {
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "Hello, "
  }

  pub fun hello(message: String, name: String): String {

    emit HelloEvent(message: message)
    return self.greeting.concat(name).concat("!")
  }
}
```