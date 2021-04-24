### Transaction hash deployed by testnet contract: 
da6734817f012b24da47b6c6718c3fe816d7a2bc0d22f5c7fc5cba0a92d20a71

### Call the transaction hash of the contract: 
ef881e73cf14f5b5576515ec11cb1646a52f9eb4bb1a6a160449be19ac93f747

### Contract code:
```ts
access(all) contract HelloWorld {
  pub let greeting: String
  pub event HelloEvent(message: String)

  init() {
    self.greeting = "Hello, World! -- TLT"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }
}
```