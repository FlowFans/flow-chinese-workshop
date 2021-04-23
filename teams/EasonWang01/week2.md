## 測試網合约部署的交易 hash
ef7646a0bd894629ee2f12835db5e1258119e50377a7d5dcf543359450fa959a

## 調用合约的交易 hash
29b6b50ffbf1db78b4947f3b71156bafedfb05a8185bb2c5fdba453956d0dee0

## 合约的源代码

```ts
pub contract HelloWorld1 {
  pub resource HelloAsset {
    pub fun hello(): String {
      return "Hello, World!"
    }
  }

  init() {
    let newHello <- create HelloAsset()

    self.account.save(<-newHello, to: /storage/Hello)

    log("HelloAsset created and stored")
  }
}
```