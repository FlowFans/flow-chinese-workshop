# 部署合约Hash

> 71e31f87e47b29378a743e356e0b473e19119843381c7192e1969f8150f1037e

# 调用合约Hash

> bd28606f030e9a4ebe35cabfe5439cd17a280822fda39e3ab2d1f93470241479

# 代码

```
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
