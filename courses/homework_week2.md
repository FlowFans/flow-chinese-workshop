##  第二次作业提交

+ 部署账号: https://flow-view-source.com/testnet/account/0x9d75a22c9bcbfd40

+ 合约部署的交易hash: 6f45005b29ada4bad2b46cfa2a493b2f1f86b4ab4f40b161dc7ecc72e10aeb29

+ 合约调用的交易hash: 055d66409c67cca3599639080a6278908b828ac468940bc985e9dd86555f3eec

+ 合约代码:

`
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

`
