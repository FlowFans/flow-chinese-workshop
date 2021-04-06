## 使用 Javascript 与 Flow 交互

本文假设读者是熟悉 JavaScript 和 [React](https://reactjs.org/) 的开发者，对 [Flow](https://www.onflow.org/) 有着一定的了解，或者熟悉 Flow 智能合约语言 [Cadence](https://docs.onflow.org/tutorial/cadence/00-introduction) 相关的概念。

我们将通过本文熟悉并搭建本地开发环境，使用 JavaScript 根据现有的 [JS-SDK](https://github.com/onflow/flow-js-sdk) 完成对链的调用与交互。

包含以下内容：

- 搭建本地开发模拟环境

- 部署开发版钱包服务

- 使用 Dev Wallet 创建本地账户

- 查询账户信息

- 执行 Cadence 脚本

- 部署 Cadence 合约并与之交互

> 教程内容参照了原文[flow-js-sdk quick start](https://github.com/onflow/flow-js-sdk/tree/master/packages/fcl#quick-start) 内容根据最新的代码和示例略有增补。

### 初始化仓库和开发环境

> 为了方便读者理解，我们直接使用 flow-js-sdk 官方提供的代码库作为基础，并针对原有的示例略有一些调整，请参照 fork 的仓库 [react-fcl-demo](https://github.com/caosbad/react-fcl-demo) 来完成部署和演示

```shell
git clone https://github.com/caosbad/react-fcl-demo.git 
cd react-fcl-demo
yarn 
```

首先将远程仓库克隆到本地，然后在实例项目中安装依赖，yarn 会将 `package.json` 文件中的项目依赖 `@onflow/fcl` `@onflow/sdk` `@onflow/six-set-code` `@onflow/dev-wallet`  等下载。

- [`@onflow/fcl`](https://github.com/caosbad/flow-js-sdk/blob/master/packages/fcl)  -- 基于 DApp 开发者的需求，对 @onflow/sdk 的一层封装。

- [`@onflow/sdk`](https://github.com/caosbad/flow-js-sdk/blob/master/packages/sdk)  -- 使用 JavaScript 进行 [build](https://github.com/caosbad/flow-js-sdk/blob/master/packages/sdk/src/build), [resolve](https://github.com/caosbad/flow-js-sdk/blob/master/packages/sdk/src/resolve), [send](https://github.com/caosbad/flow-js-sdk/blob/master/packages/send) 和 [decode](https://github.com/caosbad/flow-js-sdk/blob/master/packages/decode) 工具，与 Flow 链进行交互。

- [`@onflow/dev-wallet`](https://github.com/caosbad/flow-js-sdk/blob/master/packages/dev-wallet) -- 提供给测试与开发的本地钱包环境。

在这之前，我们还需要初始化 Flow 的本地模拟器启动 wallet-dev 服务

#### 安装 & 启动模拟器

模拟器是帮助我们在本机启动一个本地的 Flow 网络，类似于以太坊的 ganache，模拟器的下载安装步骤可以参考这里 [instructions](https://github.com/onflow/flow/blob/master/docs/cli.md#installation). 

```shell
// Linux and macOS
sh -ci "$(curl -fsSL https://storage.googleapis.com/flow-cli/install.sh)"

// Windows
iex "& { $(irm 'https://storage.googleapis.com/flow-cli/install.ps1') }"
```

```shell
// --init 参数是在第一次启动的时候添加，如果已经初始化过，就直接执行 start 命令
flow emulator start --init 
```

在示例项目的目录里执行 init 命令后，我们会发现目录下多出了一个 `flow.json` 文件，类似于以下的结构：

```json
{
    "accounts": {
        "service": {
            "address": "f8d6e0586b0a20c7",
            "privateKey": "84f82df6790f07b281adb5bbc848bd6298a2de67f94bdfac7a400d5a1b893de5",
            "sigAlgorithm": "ECDSA_P256",
            "hashAlgorithm": "SHA3_256"
        }
    }
}
```

模拟器启动之后你会看到启动成功的日志，模拟器提供了 gRPC 和 http 通信的接口

![模拟器](https://trello-attachments.s3.amazonaws.com/5fccc55f9c47787592af6b96/634x173/169f1b96e21ad2553331ead89e65c75a/image.png)

接下来在新的终端启动 Dev wallet 

#### 启动 Dev wallet 服务

在 `package.json` 文件中，我们会看到 `scripts` 的配置项中有名为 `dev-wallet` 和 `dev-wallet-win` 两个脚本，现在把我们上一步模拟初始化生成的 `privateKey` 覆盖现有的配置。

然后执行 `yarn run dev-wallet` 或 `yarn run dev-wallet-win`

成功之后，将会看到以下的日志：

![dev wallet](https://trello-attachments.s3.amazonaws.com/5fccc55f9c47787592af6b96/634x405/df970bc9b37de3fd16dd934a900e41a2/image.png)

> 这里启动了多个服务，同时注意 Service Address 和 Private Key 与模拟器生成的一致。

环境已经配置成功，接下来就是启动示例项目：

### 启动示例项目

```shell
yarn start
```

确保模拟器和 Dev wallet 也在启动的状态，我们可以看到页面上的一些示例操作，下面我们从代码层面了解一些交互的细节

#### 获取最新区块

```typescript
// src/demo/GetLatestBlock.tsx
import { decode, send, getLatestBlock } from "@onflow/fcl"

const GetLatestBlock = () => {
  const [data, setData] = useState(null)

  const runGetLatestBlock = async (event: any) => {
    event.preventDefault()

    const response = await send([
      getLatestBlock(),
    ])

    setData(await decode(response)) // 解码返回的数据，并更新 state
  }
```

```json
// 返回结果
{
  "id": "de37aabaf1ce314da4a6e2189d9584b71a7f302844b4ed5fb1ca3042afbad3d0", // 区块的 id
  "parentId": "1ae736bdea1065a98262348d5a7a2141d2b21a76ac2184b3e1181088de430255",  // 上一个区块的 id
  "height": 2,
  "timestamp": {
    "wrappers_": null,
    "arrayIndexOffset_": -1,
    "array": [
      1607256408,
      195959000
    ],
    "pivot_": 1.7976931348623157e+308,
    "convertedPrimitiveFields_": {}
  },
  "collectionGuarantees": [
    {
      "collectionId": "49e27fcf465075e6afd9009478788ba801fefa85a919d48df740e541cc514497",
      "signatures": [
        {}
      ]
    }
  ],
  "blockSeals": [],
  "signatures": [
    {}
  ]
}
```

#### 查询用户信息

这里需要我们输入用户地址来完成查询，

```typescript
// src/demo/GetAccount.tsx
  const runGetAccount = async (event: any) => {
    const response = await fcl.send([
      fcl.getAccount(addr),            // 通过地址获取用户信息
    ])

    setData(await fcl.decode(response))
  }
```

```json
{
  "address": "01cf0e2f2f715450",      // 地址
  "balance": 0,
  "code": {},
  "keys": [
    {
      "index": 0,
      "publicKey": "7b3f982ebf0e87073831aa47543d7c2a375f99156e3d0cff8c3638bb8d3f166fd0db7c858b4b77709bf25c07815cf15d7b2b7014f3f31c2efa9b5c7fdac5064d",  // 公钥
      "signAlgo": 2,
      "hashAlgo": 3,
      "weight": 1000,
      "sequenceNumber": 1
    }
  ]
}
```

#### 执行脚本

执行脚本我们可以理解为是一种无需用户授权的查询操作

```typescript
// src/demo/ScriptOne.tsx

const scriptOne = `\
pub fun main(): Int {
  return 42 + 6
}
`

const runScript = async (event: any) => {
    const response = await fcl.send([
      fcl.script(scriptOne),
    ])
    setData(await fcl.decode(response)) // 48
  }
```

#### 用定义的结构解析脚本运行的结果

这里我们可以看到在智能合约里可以定义复杂的数据结构， 并且通过 typescript 的类型进行数据的解构，能够将复杂的数据与前端的应用层友好的关联。

```typescript
// src/model/Point.ts 这里定义了结构数据的类型
class Point {
  public x: number;
  public y: number;

  constructor (p: Point) {
    this.x = p.x
    this.y = p.y
  }
}

export default Point;

// src/demo/ScriptTwo.tsx
const scriptTwo = `
pub struct SomeStruct {
  pub var x: Int
  pub var y: Int

  init(x: Int, y: Int) {
    self.x = x
    self.y = y
  }
}

pub fun main(): [SomeStruct] {
  return [SomeStruct(x: 1, y: 2), SomeStruct(x: 3, y: 4)]
}
`;

fcl.config()
  .put("decoder.SomeStruct", (data: Point) => new Point(data)) // 这里定义了 fcl 对数据的解构方式

  const runScript = async (event: any) => {
    event.preventDefault()

    const response = await fcl.send([   // 脚本的执行可以认为是一种读操作，不需要用户授权
      fcl.script(scriptTwo),
    ])

    setData(await fcl.decode(response))
  }

// class 中的 public 和 脚本中的 pub 替换 
```

这里需要注意几点：

- config 中 decoder.SomeStruct 名称要与脚本中的 SomeStruct 类型名称对应

- 回调函数中的 data 要指定对应的类型，也就是负责解构的 Point 类型

- 解构的类型，需要有自己的 constructor 函数

```json
// 输出结果
Point 0
{
  "x": 1,
  "y": 2
}
--
Point 1
{
  "x": 3,
  "y": 4
}
--
```

#### 登入（创建账户）登出

确保我们本地运行了 Dev wallet  服务

在 demo 的页面点击  Sign In/Up Dev wallet 将会弹出授权页面：

![sign up](https://trello-attachments.s3.amazonaws.com/5fccc55f9c47787592af6b96/738x276/8a33ec1d21827d3c3db1ffca2fed7923/image.png)

接着点击授权，会进入到更新 profile 的界面

![profile](https://trello-attachments.s3.amazonaws.com/5fccc55f9c47787592af6b96/735x356/4d2f927326b154ab27e07cdd8a5d4ce7/image.png)

保存并应用之后，Dev wallet 会将 profile 的信息存入数据库中，订阅函数将会执行回调，将 user 的信息作为参数传递回来

```ts
// src/demo/Authenticate.tsx
const signInOrOut = async (event: any) => {
    event.preventDefault()

    if (loggedIn) {
      fcl.unauthenticate() // logout 
    } else {
      fcl.authenticate() // sign in or sign up ，这里会呼出 Dev wallet 的窗口
    }
  }

// line:38
fcl.currentUser().subscribe((user: any) => setUser({...user})) // fcl.currentUser() 这里提供了监听方法，并动态获取 User 数据
```

对应用开发者来说，fcl 帮助我们管理用户的登录状态和所需要的授权操作，会在下文发送交易的章节详述。

```json
// user 返回值
{
  "VERSION": "0.2.0",
  "addr": "179b6b1cb6755e31",  // 用户的地址
  "cid": "did:fcl:179b6b1cb6755e31",
  "loggedIn": true,            // 登录状态
  "services": [                // 服务数据
    {
      "type": "authz",
      "keyId": 0,
      "id": "asdf8701#authz-http-post",
      "addr": "179b6b1cb6755e31",
      "method": "HTTP/POST",
      "endpoint": "http://localhost:8701/flow/authorize",
      "params": {
        "userId": "37b92714-2713-41b0-9749-fc08b3fdd827"
      }
    },
    {
      "type": "authn",
      "id": "wallet-provider#authn",
      "pid": "37b92714-2713-41b0-9749-fc08b3fdd827",
      "addr": "asdf8701",
      "name": "FCL Dev Wallet",
      "icon": "https://avatars.onflow/avatar/asdf8701.svg",
      "authn": "http://localhost:8701/flow/authenticate"
    }
  ]
}
```

Dev wallet 会将 profile 的数据存储到 GraphQL 的数据服务中，供第二次登陆时由 Dev wallet 调用展示

#### 发送交易

```ts
 // src/demo/SendTransaction.tsx

const simpleTransaction = `
transaction {
  execute {
    log("A transaction happened")
  }
}
`
const { transactionId } = await fcl.send([
   fcl.transaction(simpleTransaction),
   fcl.proposer(fcl.currentUser().authorization), // 交易触发者
   fcl.payer(fcl.currentUser().authorization),    // 费用支付者
  ])

    setStatus("Transaction sent, waiting for confirmation")

 const unsub = fcl
    .tx({ transactionId })
    .subscribe((transaction: any) => {
       setTransaction(transaction)  // 更新 state

       if (fcl.tx.isSealed(transaction)) {
          setStatus("Transaction is Sealed")
          unsub()
          }
        })
```

与执行脚本不同的是，这里的 `transaction` 调用需要发起交易，相当于执行链上的交易操作，虽然这里只进执行了 log ，但仍需要指定交易发起者和费用支付者。

```json
// 脚本运行成功的返回值
{
  "status": 4,
  "statusCode": 0,
  "errorMessage": "",
  "events": []     // 触发的事件列表
}
```

由于这里 Cadence 脚本只单纯的只执行了 log ，events 里没有数据返回

#### 部署合约

这里我们定义一个示例的合约，并声明一个公开可调用的函数来通过外部传入的参数触发合约的事件，这里并未对合约的变量进行改变

```ts
// src/demo/DeployContract.tsx
// 需要部署的合约脚本，这里为了测试方便我添加了 access(all) 的访问权限声明
const simpleContract = ` 
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

// 为账户部署合约的脚本
const deployScript = `
transaction(code: String) {
  prepare(acct: AuthAccount) {
      acct.contracts.add(name: "HelloWorld", code: code.decodeHex())
  }
}
`


const runTransaction = async (event: any) => {
    const result = await Send(simpleContract, deployScript);  // 这里的 send 是一个封装函数
    setTransaction(result);  // 更新 state
  }


// src/helper/fcl-deployer.ts 
export async function Send(code: string, deployScript: string) {
    const response = await fcl.send([
        setCode({
            proposer: fcl.currentUser().authorization,
            authorization: fcl.currentUser().authorization,     
            payer: fcl.currentUser().authorization,             
            code: code,
            deployScript: deployScript
        })
    ])

    try {
      return await fcl.tx(response).onceExecuted()  // 返回执行结果
    } catch (error) {
      return error;
    }
}
```

- Send 函数封装了当前用户的签名信息

- 部署合约的脚本和合约自身都是 Cadence 脚本

- `fcl.tx(res).onceExecuted` 可以用作交易执行的监听函数

- `acct.contracts.add(name: "HelloWorld", code: code.decodeHex())` 其中 `add` 函数的 `name` 参数需要与合约脚本中声明的名称一致

- 同样名字的合约在一个账户下只能有一个

```json
// 部署合约返回值
{
  "status": 4,
  "statusCode": 0,
  "errorMessage": "",
  "events": [
    {
      "type": "flow.AccountContractAdded",    // 类型
      "transactionId": "8ba62635f73f7f5d3e1a73d5fd860ea7369662109556e510b4af904761944e2a",  // trx id
      "transactionIndex": 1,
      "eventIndex": 0,
      "data": {
        "address": "0x179b6b1cb6755e31",  // 地址
        "codeHash": [...],               // 编码之后的合约 code ，此处有省略
        "contract": "HelloWorld"        // 合约名称
      }
    }
  ]
}
```

#### 与合约交互

在界面上我们需要输入之前部署合约账户的地址，才能够成功的导入合约并调用其公开的函数，注意调用的交易体中（由 transaction 包裹，execute 中执行合约代码的调用），传入 massage 作为合约方法的参数

```ts
// src/demo/ScriptOne.tsx
// 这里的 addr 是我们部署合约的地址
const simpleTransaction = (address: string | null) => `\
  import HelloWorld from 0x${address}

  transaction {
    execute {
      HelloWorld.hello(message: "Hello from visitor")
    }
  }
`
  const runTransaction = async (event: any) => {
    try {
      // 通过 transactionId 获得交易监听
      const { transactionId } = await fcl.send([
        fcl.transaction(simpleTransaction(addr)),
        fcl.proposer(fcl.currentUser().authorization),
        fcl.payer(fcl.currentUser().authorization),
      ])
      // 交易的监听函数定义，返回值是取消监听的函数
      const unsub = fcl
        .tx({
          transactionId,  // 解构出交易 id
        })
        .subscribe((transaction: any) => {
          setTransaction(transaction)    // 更新 state

          if (fcl.tx.isSealed(transaction)) {
            unsub()    // 取消监听
          }
        })
    } catch (error) {
      setStatus("Transaction failed")
    }
  }
```

```json
{
  "status": 4,
  "statusCode": 0,
  "errorMessage": "",
  "events": [
    {
      "type": "A.179b6b1cb6755e31.HelloWorld.HelloEvent",     // 调用的合约的事件类型
      "transactionId": "28ec7c9c0eecb4408dfc3b7b23720a6038a8379721eb7b532747cfc016a3b1cc",
      "transactionIndex": 1,
      "eventIndex": 0,
      "data": {                                         // 数据
        "message": "Hello from visitor"                 // 事件监听的参数
      }
    }
  ]
}
```

- `unsub = fcl.tx({transactionId}).subscribe(func)` 是交易结果监听的方式，等价于 `fcl.tx(response).onceExecuted()`

- 事件的 type 以  A.用户地址.合约名称.事件名称  规则来命名

- fcl 已经将获取到的数据进行了解码操作，可以直接看到返回的结果

最后合约的监听函数可以获得合约触发的事件，并将其通过回调函数返回给我们。

### 最后

现在我们已经熟悉了如何通过 fcl 与 flow 链的交互，我们已经具备了在 Flow 链上开发 DApp 的最小知识，接下来可以继续根据现有的 demo 做一些测试，或者深入探索有关 Dev wallet 、flow-sdk 或 Cadence 相关的代码与服务，相信会有更多的收获。
