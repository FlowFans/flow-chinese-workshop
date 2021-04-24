### 创建合约Hash
  > ed1fdde887774866f6860df86088860c20367d01ca972995c9dfed772f4582d4

### 调用合约Hash
> 7ff19a03806c016bd2d6b27ba488543cbf130b55f1f3fb01e3d54d065c5b4643

### 源代码

  ```ts
  access(all) contract HelloWorld3 {
  pub let greeting: String
  pub var greetingCount: Int

  pub event HelloEvent(message: String)

  init() {
    self.greeting = "Hello, World!"
    self.greetingCount = 0
  }

  pub fun hello2(message: String): String {
    self.greetingCount = self.greetingCount + 1
    emit HelloEvent(message: message)
    return message
  }
}

```