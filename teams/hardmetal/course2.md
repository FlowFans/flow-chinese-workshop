# 测试网部署合约哈希

0a251b4b51a64d842f71c11f19ea10d20b4586c7fd123b70e7a3f548710ba9a7

# 调用合约的哈希

2b28a06b0c71f9437deb067cd0b8924e9d382e9fd0be302b5b1a0c82cd92c070

# 合约源代码 
```
access(all) contract HelloWorld {
  pub let greeting: String
  pub event HelloEvent(message: String)
  init() {
    self.greeting = "Hello, Flow!"
  }
  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}
```