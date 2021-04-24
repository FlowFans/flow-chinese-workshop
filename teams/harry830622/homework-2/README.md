# Homework 2

## Contract Deployment Transaction ID
8c36b9fbefbd230484632b4ccee855710877f302b8d67501cccc471814e19ac4

## Contract Call Transaction ID
a51b0225f554e80cf2eb02b4453959fdbb3a3920d4b664a99ec3a3cbdd9eb362

## Contract Code
Deploy at `0x8522c973359f8bf7`:
```cdc
access(all) contract HelloWorldHarry {
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "Hello, World! I am Harry!"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}
```
