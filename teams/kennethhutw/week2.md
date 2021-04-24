//部署合约的Hash: 1b1ecfe567a9d029e6d7383a6892034be5773bee0fa79c8c978738d26ea19d7d

//调用合约的Hash: 0b899af717d3b13530df90196f6c24a49c99851e0fb83e4fb1d77db09fff5025


//合约代码(未修改):
 
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