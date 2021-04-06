# Flow JS-SDK 源码分析 —— 查询的组装与数据解码

在之前的文章 ——「使用 Javascript 与 Flow 交互」里，我们熟悉了开发环境的初始化和 FCL 与 Flow 链进行交互的一些功能，这篇文章我们将从 [onflow/flow-js-sdk](https://github.com/onflow/flow-js-sdk)源代码的层面去分析，sdk 是如何完成交易的组装和验证。

### 交互举例

以查询账户信息`getAccount`这个基本的账户查询举例， 我们根据具体实现的代码来分析：

```js
const response = await sdk.send(await sdk.build([
      sdk.getAccount(addr)
    ]), { node: "http://localhost:8080" })

    setResult(await sdk.decodeResponse(response))
```

这里面涉及到几个 package ：

- send ——  初始化配置，根据构造的交易类型，调用不同的 `send` 方法
  
  `packages/send/src/send.js`

- build —— 初始化 `interaction` 并传递交易函数
  `packages/sdk/src/build/index.js`

- gerAccount —— 初始化查询交易，并校验地址的合法性
  `packages/sdk-build-get-account/src/index.js`

- decodeResponse —— 根据返回结果解码
  `packages/decode/src/decode.js`

#### getAccount

```js
import {pipe, makeGetAccount, Ok} from "@onflow/interaction"
import {sansPrefix} from "@onflow/util-address"

export function getAccount(addr) {
  return pipe([
    makeGetAccount, // 构造具体查询交易
    ix => {
      ix.accountAddr = sansPrefix(addr) // 截取地址前缀
      return Ok(ix) // 返回状态
    }
  ])
}
```

##### makeGetAccount

这里我们需要着重分析一下交易体构造的工具 [interaction](https://github.com/onflow/flow-js-sdk/blob/master/packages/interaction/src/interaction.js)  ，interaction 是及不同交易类型构造，处理与验证的整体，包括了参数处理，类型验证和状态设置，这里我们先按照交易构造流程完成流程的熟悉。

```js
export const GET_ACCOUNT /*             */ = 0b0000010000
export const OK /*  */ = 0b10

/*    ...      */

export const Ok = (ix) => {
  ix.status = OK
  return ix
}

/*    ...      */

const makeIx = (wat) => (ix) => {
  ix.tag = wat
  return Ok(ix)
}

/*    ...      */

export const makeGetAccount /*            */ = makeIx(GET_ACCOUNT)
```

在这里，通过将 IX 数据类型的替换，把具体的交易查询的 `tag` 设置为二进制的数据类型，与具体的查询或交易体相对应。

接着将交易体的状态设置为 `OK` 类型，并将其返回，我们可以通过下面的代码了解 IX 交易体的数据结构。

##### 交易体 IX 的初始结构

```js
export const UNKNOWN /*                 */ = 0b0000000001
export const OK /*  */ = 0b10

const IX = `{
  "tag":${UNKNOWN},
  "assigns":{},
  "status":${OK},
  "reason":null,
  "accounts":{},
  "params":{},
  "arguments":{},
  "message": {
    "cadence":null,
    "refBlock":null,
    "computLimit":null,
    "proposer":null,
    "payer":null,
    "authorizations":[],
    "params":[],
    "arguments":[]
  },
  "proposer":null,
  "authorizations":[],
  "payer":null,
  "events": {
    "eventType":null,
    "start":null,
    "end":null
  },
  "latestBlock": {
    "isSealed":null
  },
  "block": {
    "isSealed":null,
    "id":null,
    "height":null
  },
  "accountAddr":null,
  "transactionId":null
}`
```

最终我们通过 `getAccount` 获取了一个具备类型的交易体数据，并准备进行下一步的处理，我们也应该注意到了，在 `build` 函数中，嵌套调用了 `pipe` 函数，接下来是了解 `pipe` 函数发挥的作用

##### Pipe

顾名思义，Pipe 提供了将 IX 中所需要的多种处理函数或数据组装需求而实现的管道调用逻辑，最终形成一个链式处理的结果，最终返回组装与验证完成的交易结构体。

```js
const recPipe = async (ix, fns = []) => {
  ix = hardMode(await ix) // 严格校验与 IX 结构比对
  if (isBad(ix) || !fns.length) return ix // 判断处理的错误状态或结束递归的条件
  const [hd, ...rest] = fns // 处理函数拆分
  const cur = await hd
  if (isFn(cur)) return recPipe(cur(ix), rest) // 不同逻辑的递归处理
  if (isNull(cur) || !cur) return recPipe(ix, rest)
  if (isInteraction(cur)) return recPipe(cur, rest)
  throw new Error("Invalid Interaction Composition")
}

export const pipe = (...args) => {
  const [arg1, arg2] = args
  if (isArray(arg1) && arg2 == null) return (d) => pipe(d, arg1) // 拆分链式调用的数组函数
  return recPipe(arg1, arg2) // 启动递归处理
}
```

- fns 中的处理函数会返回 IX 结构给下一个函数

- IX 通过 `isBad` 来判断上一函数交易结构的状态

- `recPipe` 会递归处理所有链式调用的函数

#### build

```js
import {pipe, interaction} from "@onflow/interaction"


export function build(fns = []) {
  return pipe(interaction(), fns)
}


// interaction

export const interaction = () => JSON.parse(IX)
```

这里我们看到 `build` 函数也是调用了 `pipe` 方法先完成了 IX 基本结构的初始化，然后将 `getAccount` 计算的 IX 结果作为数组元素传入到 `pipe` 方法中，其实这里并没有用到最先初始化的 IX，而是通过 `recPipe` 函数中的判断，使用 `getAccount` 函数计算出来的 IX 替换初始化的 IX 结构。

```js
 if (isInteraction(cur)) return recPipe(cur, rest) // cur 作为 getAccount 输出的结构，替换了 interaction() 
```

#### send

我们查看 [send](https://github.com/onflow/flow-js-sdk/blob/master/packages/send/src/send.js#L23) 源代码可以看到 `send` 函数起到了校验和路由的功能，代码有删减，只保留了 sendGetAccount 

```js
export const send = async (ix, opts = {}) => {
  opts.node = opts.node || (await config().get("accessNode.api")) // 初始化自定义节点配置
  ix = await ix // 

  // 根据交易类型，决定返回具体的交易方法
  switch (true) {
   /*    ...      */
    case isGetAccount(ix): // 路由到相对的查询
      return sendGetAccount(ix, opts)
   /*    ...      */
    default:
      return ix
  }
}
```

##### sendGetAccount

这里 Flow SDK 使用了 gRPC  `protoc` 的工具定义交互的数据类型，较为易于研发与维护的数据交互方式。详情请见 `packages/protobuf`[protobuf ](https://github.com/onflow/flow-js-sdk/tree/master/packages/protobuf) 在此不做详述。

```js
export async function sendGetAccount(ix, opts = {}) {
  ix = await ix // 获得具体的 ix 结构

  const req = new GetAccountRequest() // 定义 gRPC message 结构
  req.setAddress(addressBuffer(sansPrefix(ix.accountAddr))) // 设置请求参数

  const res = await unary(opts.node, AccessAPI.GetAccount, req) // 获得实例对象

  let ret = response() // 初始化响应数据
  ret.tag = ix.tag // 赋值请求类型

  const account = res.getAccount() // 调用查询 api 获得数据示例
  ret.account = {
    address: withPrefix(u8ToHex(account.getAddress_asU8())), // 获得地址信息并添加前缀
    balance: account.getBalance(), // 获得余额
    code: account.getCode_asU8(), // 获得地址中部署的合约代码
    keys: account.getKeysList().map(publicKey => ({ // 遍历地址下绑定的 key 信息
      index: publicKey.getIndex(),
      publicKey: u8ToHex(publicKey.getPublicKey_asU8()),
      signAlgo: publicKey.getSignAlgo(),
      hashAlgo: publicKey.getHashAlgo(),
      weight: publicKey.getWeight(),
      sequenceNumber: publicKey.getSequenceNumber(),
    })),
  }

  return ret
}
```

- `GetAccountRequest` 是 gRPC 定义的信息交互类型，需要设置参数才能与 gRPC 交互

- unary 同样是获得 gRPC 数据的封装接口，接受 node 和查询所需参数类型与实际参数数据，返回 message 数据

- 调用定义的查询类型 `getAccount` 获得具体 account 的值

- 组装并填充至初始化后的 `response` 数据结构中

```js
// packages/response/src/response.js
const DEFAULT_RESPONSE =
'{"tag": 0, "transaction":null, "transactionId":null, "encodedData":null, "events": null, "account": null}'

export const response = () => JSON.parse(DEFAULT_RESPONSE)
```

现在数据已经获取到，接下来就是将返回的数据解码处理，为项目所用。

#### decode

最后一步是将获取到的数据进行解码，在 [decode.js](https://github.com/onflow/flow-js-sdk/blob/master/packages/decode/src/decode.js#L128) 中定义了请求数据响应的解码函数 `decodeResponse` (代码略有删减)

```js
// 解码响应数据
export const decodeResponse = async (response, customDecoders = {}) => {
  let decoders = { ...defaultDecoders, ...customDecoders }

  if (response.encodedData) { // 返回的查询数据解码
    return await decode(response.encodedData, decoders)
    /* ... */
  } else if (response.account) { // 返回账户响应的数据解码，这里只对 account 中 code 的数据进行解码处理
    const acct = response.account 
    acct.code = new TextDecoder("utf-8").decode(acct.code || new UInt8Array()) // 解码账户的合约代码
    return acct
    /* ... */
  } else if (response.transactionId) {
    return response.transactionId
  }

  return null
}
```

我们可以看到，在 account 返回值的 decode 逻辑中，只对账户的代码进行了解码，其余的都按照查询出来的结果直接返回，从 [protobuf](https://github.com/onflow/flow-js-sdk/blob/master/packages/protobuf/src/proto/flow/entities/account.proto#L7) 的定义中我们也知道 request 查询出的结果已经被自动转为相对应的类型。

## 最后

我们从 getAccount 的查询操作开始，从交易体构建，验证，交易类型的设置与分发，再到 interaction 中链式处理的逻辑，经过 gRPC 的查询和类型定义，最后到 Decode 的解码，还原了整个 Flow JS-SDK 的数据查询与获取流程。

在这个过程中我们也对 Flow JS-SDK 的模块化设计有了一个比较清晰的认识，也总结出其设计的特点：

- 所有的查询或交易都遵循同样的处理原则

- 交易查询的封装和处理用交易体数据直接体现，更加直观

- 交易处理根据不同的需求增加了灵活的链式调用逻辑

- 在查询交易打包之前，不同模块完成自己的校验逻辑

- gRPC 作为数据类型的定义，可以提高应用层查询的便利性

这次源码分析了较为简单的查询流程，让我们对应用层与区块链交互的流程有一个大致的思路，感兴趣的同学可以举一反三的查看其它查询流程的细节逻辑，相关流程的代码和注释维护在[Github](https://github.com/caosbad/flow-js-sdk/tree/comm)，供大家参考。
