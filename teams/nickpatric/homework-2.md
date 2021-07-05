### 创建合约Hash
> 1c9f50a110bf89b4778368ba85dc1268174d5871876ce9c31f54dd9d8d312bf4
### 调用合约Hash
> 08eba71dc6594045234c723d5d70cb50c8cdfd5db1007ecc2110c6b6acbff600
### 源代码

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