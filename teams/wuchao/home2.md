
1. deploy contract   transactionId:7cdd5e4a5c5f3da2ab9e672868244ad52858bd4cbaf18f160f801e9f2b7bb61a
2. Interact with contract transactionId:db2e66020ebd0d2224f13fcfce5a33336f2a0a5ee50ff2166782f5a7e8de253a
3. code(unchange):

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