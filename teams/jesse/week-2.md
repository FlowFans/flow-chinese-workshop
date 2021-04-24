# contract deployment transaction hash

9888f87379beebd0b53dbf09a68dff09f04be04039e32920e27f74530c67fcca

# contract interaction transaction hash

354c7cc8d2febc699092f4ce9cbec4433ca4e2ba67b76d0bc7a32b2a966dd50f

# contract source code

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
