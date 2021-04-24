- Contract Deployment TX id: 1b31a7d037c759d4e4824e9491692765487929a3a2fc6fa51742911766d0dedc

- Interact with contract TX id: f3aff7b0b9cdf9206d6c91c15ce4ff30027419aee75e1f19e662163dc03ba996

- Contract code:

```
access(all) contract HelloWorld {
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "Xing here!"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}
```
