# Flow JS-SDK 源码分析 —— 部署 Cadence 合约

在上一篇 「查询的组装与数据解码」的源码分析文章里，我们一起熟悉了 Flow JS-SDK 与 Flow 链相关的查询流程，其主要执行的都是查询相关的操作，并不涉及到编码和签名相关的流程。

本篇文章，我们将从合约的部署和交互的例子开始，通过阅读源码来学习如何进行写入的操作流程。

我们将从两个部分开始：

- 部署 Cadence 合约

- 接着与 Cadence 交互

## 部署 Cadence 合约

我们在 FCL demo 中也看到了有部署合约的例子，这个文件里定义了用户需要部署的合约脚本，使用 `@onflow/six-set-code` 组装部署脚本的交易体：

```ts
// examples/react-fcl-demo/src/demo/DeployContract.tsx

const simpleContract = `
pub contract HelloWorld {
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

const result = await Send(simpleContract) // fcl-deployer.ts


// examples/react-fcl-demo/src/helper/fcl-deployer.ts

import * as fcl from "@onflow/fcl"
import { template as setCode } from "@onflow/six-set-code"

export async function Send(code: string) {
    const response = await fcl.send([
        setCode({
            proposer: fcl.currentUser().authorization,
            authorization: fcl.currentUser().authorization,     
            payer: fcl.currentUser().authorization,             
            code: code,
        })
    ])

    try {
      return await fcl.tx(response).onceExecuted()
    } catch (error) {
      return error;
    }
}
```

这里我们看到 `setCode` 函数将当前授权用户的授权信息与 code 一起传入

- Proposer —— 交易发起人

- Authorization —— 授权人

- Payer —— 费用支付人

- code —— 需要部署的合约代码

这里将三种角色都设置为当前登录的用户，我们来看看 `setCode` 内部的实现：

```ts
// packages/six-set-code/src/six-set-code.js
import * as sdk from "@onflow/sdk"
import * as t from "@onflow/types"
// 定义默认的常量参数
export const TITLE = "Set Account Code"
export const DESCRIPTION = "Set an Account Code on Flow with given code."
export const VERSION = "0.0.0"
export const HASH = "7375dc3feb96e2f8061eff548220a96bf77ceb17affd1ac113f10d15411a92c4"

// 部署账户代码的脚本，也是 cadence 语法，当前的版本已经旧了，现在使用的 acct.setCode 已经被替换为 
// acct.contracts.add(name: "HelloWorld", code: code.decodeHex())
export const CODE = 
`transaction(code: String) {
    prepare(acct: AuthAccount) {
        acct.contracts.add(name: "HelloWorld", code: code.decodeHex())
    }
}`

// 这里的函数我添加了一个折中方案，将部署脚本也作为参数传入
export const template = ({ proposer, authorization, payer, code = "", deployScript = CODE }) => sdk.pipe([
    sdk.transaction(deployScript),
    sdk.args([sdk.arg(Buffer.from(code, "utf8").toString("hex"), t.String)]), // 这里对 code 进行了编码与类型指定操作
    sdk.proposer(proposer),
    sdk.authorizations([authorization]),
    sdk.payer(payer),
    sdk.validator((ix, {Ok, Bad}) => {   // 添加了验证的逻辑，要求 template 类型的交易只能有一个授权
        if (ix.authorizations.length > 1) return Bad(ix, "template only requires one authorization.")
        return Ok(ix)
    })
])
```

`@onflow/six-set-code` 里定义了一些常量和默认的合约代码部署脚本，同样使用了 `pipe` 来组装交易体数据，拼装处理 IX 所需要的数据，最后进行校验。

我们要注意，`template` 方法的参数这里我稍微做了一些改动，增加了 `deployScript` 并将原有的 `sdk.transaction(CODE)`替换成了 `sdk.transaction(deployScript)`，实现的详情参考[这里](https://github.com/caosbad/react-fcl-demo/blob/master/src/demo/DeployContract.tsx#L24) ,为了能够兼容新版本的合约代码部署，我们需要自定义 `deployScript` 将其默认值覆盖掉，传入与脚本匹配的部署脚本。

新的脚本如下所示：

```ts
// 兼容新版的部署脚本
const deployScript = `
transaction(code: String) {
  prepare(acct: AuthAccount) {
      acct.contracts.add(name: "HelloWorld", code: code.decodeHex())
  }
}
`
```

### Build transaction

我们再回头看上面的交易构建函数 `sdk.transaction`，依然使用了之前提到到 `pipe` 函数，把 `deployScript` 传入调用 `interaction` 构建 IX 交易体，将需要部署到账户中的智能合约代码作为参数传入，通过 `@onflow/sdk-build-transaction` 中的代码将 `ix.cadence` 赋值，并填充其余授权的字段，最终返回 `ix`。

```ts
// packages/sdk-build-transaction/src/index.js
/* ... */
export function transaction(...args) {
  return pipe([
    makeTransaction,     // 调用 interaction 的 IX 初始化字段
    put("ix.cadence", template(...args)),   // template 这里直接返回了 string 类型的 code 内容
    ix => {     // 继续补全其他的字段
      ix.message.computeLimit = ix.message.computeLimit || DEFAULT_COMPUTE_LIMIT
      ix.message.refBlock = ix.message.refBlock || DEFUALT_REF
      ix.authorizations = ix.authorizations || DEFAULT_SCRIPT_ACCOUNTS
      return Ok(ix)
    },
  ])
}
```

接着就是不熟脚本参数的设置，我们可以看到部署代码的脚本，其实是一个 Cadence 函数，并定义了其参数的名称和类型。在初始化完成`ix.cadence`  之后，我们仍然需要将部署脚本中使用到的合约代码通过编码和参数类型的定义与部署脚本函数的参数进行关联，这样在处理部署脚本的时候，可以将需要部署的代码作为参数传入到脚本中。

```ts
// sdk export {args, arg} from "@onflow/sdk-build-arguments"

sdk.args([sdk.arg(Buffer.from(code, "utf8").toString("hex"), t.String)]) // 初始化 arg 然后拼接成 args 
```

### makeArgument

这里其实使用到了 `@onflow/sdk-build-arguments` 中的两个函数，而其中 `arg` 定义参数类型，`args` 则将参数结构化赋值给 `interaction` 中的 `ix` 结构

```ts
// packages/sdk-build-arguments/src/index.js
import {pipe, makeArgument} from "@onflow/interaction"

export function args(ax = []) {
  return pipe(ax.map(makeArgument))  // 将参数赋值给 ix 数据
}

export function arg(value, xform) {  // 返回参数和参数类型定义
  return {value, xform}
}

/* interaction */
const ARG = `{
  "kind":${ARGUMENT},
  "tempId":null,
  "value":null,
  "asArgument":null,
  "xform":null,
  "resolve": null
}`
// interaction makeArgument
export const makeArgument = (arg) => (ix) => {  // 基于 ix 数据构建参数
  let tempId = uuid()    // 生成唯一id
  ix.message.arguments.push(tempId) 

  ix.arguments[tempId] = JSON.parse(ARG) // 按照定义的模板初始化，之后基于参数数据赋值
  ix.arguments[tempId].tempId = tempId 
  ix.arguments[tempId].value = arg.value
  ix.arguments[tempId].asArgument = arg.asArgument
  ix.arguments[tempId].xform = arg.xform
  ix.arguments[tempId].resolve = arg.resolve
  return Ok(ix)   // 返回 ix
}
```

到这里，部署合约的参数都已经构建完成了，接下来就是授权相关的操作

### Build proposer

`sdk.proposer(proposer)` 这里是来自 `@onflow/sdk-build-proposer` 中的构建函数，其中也依赖了 `interaction` 定义的  `makeAccount` 

```ts
// packages/sdk-build-proposer/src/index.js
import {pipe, makeProposer} from "@onflow/interaction"

const roles = {
  proposer: true,
}

// 
export async function proposer(authz) {
  return typeof authz === "function"
    ? makeProposer({resolve: authz, role: roles, roles})  // fcl.currentUser().authorization 作为函数传递，设置给 resolve
    : makeProposer({...authz, role: roles, roles})
}


// interaction

// account 的初始化结构
const ACCT = `{
  "kind":${ACCOUNT},
  "tempId":null,
  "addr":null,
  "keyId":null,
  "sequenceNum":null,
  "signature":null,
  "signingFunction":null,
  "resolve":null,
  "role": {
    "proposer":false,
    "authorizer":false,
    "payer":false,
    "param":false
  }
}`

export const makeProposer = (acct) => (ix) => {
  let tempId = uuid()
  ix.proposer = tempId
  return Ok(pipe(ix, [makeAccount(acct, tempId)]))  // pipe 处理 ix 结构，这里我们先跳过授权的逻辑，后面的文章再详细拆解
}


// 将获得的用户信息继续拼接给 ix 并返回
const makeAccount = (acct, tempId) => (ix) => {
  ix.accounts[tempId] = JSON.parse(ACCT)
  ix.accounts[tempId].tempId = tempId
  ix.accounts[tempId].addr = acct.addr
  ix.accounts[tempId].keyId = acct.keyId
  ix.accounts[tempId].sequenceNum = acct.sequenceNum
  ix.accounts[tempId].signature = acct.signature
  ix.accounts[tempId].signingFunction = acct.signingFunction
  ix.accounts[tempId].resolve = acct.resolve
  ix.accounts[tempId].role = {         // 设置角色信息
    ...ix.accounts[tempId].role,
    ...acct.role,
  }
  return Ok(ix)
}
```

将用户信息初始化给 IX 结构中，并设置好角色信息。

### Build Authorizations

这里的流程与 Proposer 的处理方式相同，只不过将 `rose`  设置为 `{authorizer: true }`， 之后将 `ix.authorizations` 中推入传入用户的 ID, 返回组装好的 ix 对象。

### Build Payer

则是将 `ix.payer` 设置为账户对应的 uuid, `rose` 设置为 `{ payer: true }`

### Build Validator

```ts
sdk.validator((ix, {Ok, Bad}) => { // 添加了验证的逻辑，要求 template 类型的交易只能有一个授权
 if (ix.authorizations.length > 1) return Bad(ix, "template only requires one authorization.")
 return Ok(ix)
 }) 

// @onflow/sdk-build-validator 
import {update} from "@onflow/interaction" 

export function validator(cb) { // 将验证的回调函数组装至 ix 
 return update('ix.validators', validators => 
Array.isArray(validators) ? validators.push(cb) : [cb]
 )
}

// interaction 中更新 ix 指定 key 数据的函数
export const update = (key, fn = identity) => (ix) => {
 ix.assigns[key] = fn(ix.assigns[key], ix)
 return Ok(ix)
}
```

这里只是将 `validator` 的验证函数赋值给了 ix 数据结构，完成了最后一步的交易数据组装。

### send 函数

这里又进入了 `@onflow/send` 的函数中，上一篇中我们提到，send 会根据不同的 ix 结构和 `tag` 路由交易处理的逻辑，我们直接进入到 `sendTransaction` 函数查看代码，这里主要是把之前的步骤组装起来的 ix 数据转化成具体的交易数据，同时完成了数据格式的编码，具体的内容请参考注释。

```ts
// packages/send/src/send-transaction.js



export async function sendTransaction(ix, opts = {}) {
  ix = await ix

  const tx = new Transaction()    // 初始化交易体，适配 gRPC 的数据类型
  tx.setScript(scriptBuffer(ix.message.cadence))  // 设置需要执行的 Cadence 脚本
  tx.setGasLimit(ix.message.computeLimit)   // 设置 gas 上限
  tx.setReferenceBlockId(         // 设置最新的区块信息
    ix.message.refBlock ? hexBuffer(ix.message.refBlock) : null
  )
  tx.setPayer(addressBuffer(sansPrefix(ix.accounts[ix.payer].addr)))  // 设置支付人
  ix.message.arguments.forEach(arg =>
    tx.addArguments(argumentBuffer(ix.arguments[arg].asArgument))   // 转换交易的参数
  )
  // 设置权限相关信息
  ix.authorizations 
    .map(tempId => ix.accounts[tempId].addr)
    .reduce((prev, current) => {
      return prev.find(item => item === current) ? prev : [...prev, current]
    }, [])
    .forEach(addr => tx.addAuthorizers(addressBuffer(sansPrefix(addr))))
  // 交易发起人
  const proposalKey = new Transaction.ProposalKey()
  proposalKey.setAddress(
    addressBuffer(sansPrefix(ix.accounts[ix.proposer].addr))
  )
  proposalKey.setKeyId(ix.accounts[ix.proposer].keyId)
  proposalKey.setSequenceNumber(ix.accounts[ix.proposer].sequenceNum)

  tx.setProposalKey(proposalKey)
  // 如果没有设置付款人，则默认使用交易签名人作为付款人
  // Apply Non Payer Signatures to Payload Signatures
  for (let acct of Object.values(ix.accounts)) {
    try {
      if (!acct.role.payer && acct.signature != null) {
        const sig = new Transaction.Signature()
        sig.setAddress(addressBuffer(sansPrefix(acct.addr)))
        sig.setKeyId(acct.keyId)
        sig.setSignature(hexBuffer(acct.signature))
        tx.addPayloadSignatures(sig)
      }
    } catch (error) {
      console.error("Trouble applying payload signature", {acct, ix})
      throw error
    }
  }

  // 如果没有设置签名人，则默认使用付款人作为签名人
  // Apply Payer Signatures to Envelope Signatures
  for (let acct of Object.values(ix.accounts)) {
    try {
      if (acct.role.payer && acct.signature != null) {
        const sig = new Transaction.Signature()
        sig.setAddress(addressBuffer(sansPrefix(acct.addr)))
        sig.setKeyId(acct.keyId)
        sig.setSignature(hexBuffer(acct.signature))
        tx.addEnvelopeSignatures(sig)
      }
    } catch (error) {
      console.error("Trouble applying envelope signature", {acct, ix})
      throw error
    }
  }

  // 初始化 gRPC 的请求体
  const req = new SendTransactionRequest()
  req.setTransaction(tx)
  // 时间记录
  var t1 = Date.now()
  const res = await unary(opts.node, AccessAPI.SendTransaction, req)  // 广播交易
  var t2 = Date.now()

  let ret = response()  // 初始化响应对象
  ret.tag = ix.tag      // 交易类型设置
  ret.transactionId = u8ToHex(res.getId_asU8())     // trx id 赋值
  // 浏览器环境的事件广播
  if (typeof window !== "undefined") {
    window.dispatchEvent(
      new CustomEvent("FLOW::TX", {
        detail: {txId: ret.transactionId, delta: t2 - t1},
      })
    )
  }

  return ret  // 返回响应对象
}
```

这里返回的是一个具备默认结构的 `response` 对象，接着通过 `fcl.tx` 将其构造成具备监听交易结果函数的对象，将构造的 `response` 传入其中，生成响应的 `transaction` 并用其定义的 `onceExecuted` 函数完成监听。

`return await fcl.tx(response).onceExecuted()` ，至此合约的部署交易已经完成。

### 最后

这里我们从 Flow JS-SDK 构建写入交易的流程入手，了解了在部署合约这样的交易中，ix 是如何构建的，包括最终返回交易结果的监听。

这里因为篇幅所限，没有涵盖授权和与合约交互的流程，将会在后面几篇详述。本文中所涉及到的代码片段和注释均在[github](https://github.com/caosbad/flow-js-sdk/tree/comm) 中，供大家查阅。

`2021-01-10`
