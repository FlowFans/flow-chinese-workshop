# Flow JS-SDK 源码解析 —— 账号授权与签名

在之前的文章里，我们熟悉了如何使用 Flow JS-SDK 完成一些针对 FLOW 链的简单用例，也熟悉了如何使用 [fcl](https://github.com/onflow/flow-js-sdk)完成数据的查询、合约部署，交易的发送。

这篇文章是帮助我们熟悉 JS-SDK 的授权与签名操作，从之前的例子我们知道，在构建交易的过程中，`fcl.authenticate` 是用户登录和授权的入口，同时 `fcl.currentUser()` 函数可以帮助开发者获得当前用户的信息，`currentUser().subscribe` 函数会帮助我们监听用户信息的变化，那么这些是怎么实现的呢，下面我们就从源码开始学习。

## authenticate

我们先从授权的操作开始，提供授权操作的源码位置在 current-user 的包中，我们看到 index 文件中包含了几个导出的函数：

```ts
// packages/fcl/src/current-user/index.js
export const currentUser = () => {
  return {
    authenticate,  // 授权
    unauthenticate, // 取消授权
    authorization,  // 用户交易签名认证
    subscribe,  // 事件监听
    snapshot,  // 账户信息快照
  }
}
```

其中 `authenticate` 函数中调用上初始化上下文的函数，同时定义了 `iframe` 的授权窗口交互的逻辑，允许第三方的开发者提供外部的授权服务界面，帮助用户方便的使用 FLOW 

```ts
// packages/fcl/src/current-user/index.js

async function authenticate() {
  return new Promise(async resolve => {
    // 初始化全局用户信息与上下文
    spawnCurrentUser()
    // 获取当前用户的快照信息
    const user = await snapshot()
    // 如果用户有登录的信息则直接返回
    if (user.loggedIn && notExpired(user)) return resolve(user)
    // 通过 iframe 的方式呼出授权界面
    const [$frame, unrender] = renderAuthnFrame({
      handshake: await config().get("challenge.handshake"), // 在 fcl 配置中配置的第三方授权服务地址，localhost 是依赖 dev-wallet 测试网依赖 blocto
      l6n: window.location.origin,  // 当前页面的 url
    })

    // 定义响应函数
    const replyFn = async ({data}) => {
      if (data.type === CHALLENGE_CANCEL_EVENT || data.type === CANCEL_EVENT) { // 取消授权，关闭窗口，取消事件监听
        unrender()
        window.removeEventListener("message", replyFn)
        return
      }
      if (data.type !== CHALLENGE_RESPONSE_EVENT) return // 非登录响应的数据都返回
      // 正常的数据响应流程
      unrender()
      window.removeEventListener("message", replyFn)

      send(NAME, SET_CURRENT_USER, await buildUser(data)) // 根据返回的数据初始化用户信息，并设置给 ctx
      resolve(await snapshot()) // 返回用户信息
    }

    window.addEventListener("message", replyFn) // 添加消息的响应函数
  })
}
```

`authenticate` 函数做了三件事情：

- 用 `spawnCurrentUser` 初始化全局的用户结构与上下文
- 初始化第三方钱包界面渲染
- 定义全局响应函数 `message` 

我们分别从这三部分开始讲起：

#### 初始化全局用户结构上下文

@onflow/util-actor 作为初始化并提供上下文工具的库起到了很关键的作用，我们先看上下文是如何初始化的

```ts
// packages/fcl/src/current-user/index.js
import {spawn, send, INIT, SUBSCRIBE, UNSUBSCRIBE} from "@onflow/util-actor"

// 作为全局的缓存上下文的标志
const NAME = "CURRENT_USER"
// 初始化调用函数，传入定义好的操作接口
const spawnCurrentUser = () => spawn(HANDLERS, NAME)
// 用户信息的默认数据结构
const DATA = `{
  "f_type": "User",
  "f_vsn": "1.0.0",
  "addr":null,
  "cid":null,
  "loggedIn":null,
  "expiresAt":null,
  "services":[]
}`

// 操作接口
const HANDLERS = {
  [INIT]: async ctx => {
    ctx.merge(JSON.parse(DATA)) // 初始化用户信息数据结构
    if (await canColdStorage()) { // 如果支持 sessionStorage，则加载信息
      const user = await coldStorage.get()
      if (notExpired(user)) ctx.merge(user)
    }
  },
  [SUBSCRIBE]: (ctx, letter) => { // 订阅
    ctx.subscribe(letter.from)
    ctx.send(letter.from, UPDATED, {...ctx.all()})
  },
  [UNSUBSCRIBE]: (ctx, letter) => { // 取消订阅
    ctx.unsubscribe(letter.from)
  },
  [SNAPSHOT]: async (ctx, letter) => { // 获取用户信息
    letter.reply({...ctx.all()})
  },
  [SET_CURRENT_USER]: async (ctx, letter, data) => {  // 更新用户信息
    ctx.merge(data)
    if (await canColdStorage()) coldStorage.put(ctx.all())
    ctx.broadcast(UPDATED, {...ctx.all()})
  },
  [DEL_CURRENT_USER]: async (ctx, letter) => { // 重置用户信息
    ctx.merge(JSON.parse(DATA))
    if (await canColdStorage()) coldStorage.put(ctx.all())
    ctx.broadcast(UPDATED, {...ctx.all()})
  },
}
```

这里我们看到 `spawn` 函数将 `current-user` 中定义的 HANDLERS 接口传入了初始化的函数，并制定了唯一的标志 `NAME` ，操作接口中定义了一系列消息类型的处理函数，这些函数会获得上下文和消息体的内容，方便进一步的对上下文中存储的数据进行维护。

那么上下文里存储的内容是什么呢？我们接着看 @onflow/util-actor 源码是怎么做的：

```ts
// packages/util-actor/src/index.js

const root =
  (typeof self === "object" && self.self === self && self) ||
  (typeof global === "object" && global.global === global && global) ||
  (typeof window === "object" && window.window === window && window)

// 全局对象
root.FCL_REGISTRY = root.FCL_REGISTRY == null ? {} : root.FCL_REGISTRY

// @onflow/util-actor 中构造用户信息与响应上下文函数
export const spawn = (fn, addr = null) => {
  if (addr == null) addr = ++pid
  // 判断是否已经有全局的对象
  if (root.FCL_REGISTRY[addr] != null) return addr
  // 定义全局对象，设置所需要的变量
  root.FCL_REGISTRY[addr] = {
    addr, // 服务标志
    mailbox: createMailbox(), // 通信组件
    subs: new Set(), // 订阅对象
    kvs: {}, // 用户数据
  }
  // 定义上下文
  const ctx = {
    self: () => addr,
    receive: () => root.FCL_REGISTRY[addr].mailbox.receive(),
    send: (to, tag, data, opts = {}) => { // 发送消息
      opts.from = addr
      return send(to, tag, data, opts)
    },

    broadcast: (tag, data, opts = {}) => { // 广播
      opts.from = addr
      for (let to of root.FCL_REGISTRY[addr].subs) send(to, tag, data, opts)
    },
    subscribe: sub => sub != null && root.FCL_REGISTRY[addr].subs.add(sub), // 订阅
    unsubscribe: sub => sub != null && root.FCL_REGISTRY[addr].subs.delete(sub),
    update: (key, fn) => {
      if (key != null)
        root.FCL_REGISTRY[addr].kvs[key] = fn(root.FCL_REGISTRY[addr].kvs[key])
    },
    all: () => {
      // 返回存储的所有的用户数据
      return root.FCL_REGISTRY[addr].kvs
    },
    merge: (data = {}) => {
      // 合并用户信息
      Object.keys(data).forEach(
        key => (root.FCL_REGISTRY[addr].kvs[key] = data[key])
      )
    },
  }
  // 封装在 current-user 中定义的 handlers
  if (typeof fn === "object") fn = fromHandlers(fn)

  queueMicrotask(async () => {
    await fn(ctx)
    kill(addr)
  })

  return addr
}
```

这里的代码有点长，略做删减保留了 `ctx` 后面会用到的接口，这里我们看到代码用 `root.FCL_REGISTRY[addr]` 注册了一个全局唯一的对象，并初始化了默认字段 `ctx` 作为后续操作的上下文，传递给 `handler` 的接口

![全局结构在浏览器中](https://trello-attachments.s3.amazonaws.com/5aceaf1164c86a15f5956cda/5fccc55f9c47787592af6b96/5a63afa55f06f4cc3e5231fa9caba920/image.png)

这里除了初始化了全局的 `ctx` 数据结构之外，还通过 `fromHandlers` 函数，将 `handlers` 的处理函数启动了一个循环的监听，用 mailbox 完成消息的接收与发送。后面我们会看到 `letter` 数据，就是 `mailbox` 消息通信的数据封装，也在之前定义的 `handlers` 用到过。

```ts
// 将定义的处理函数进行封装，并注册监听
const fromHandlers = (handlers = {}) => async ctx => {
  if (typeof handlers[INIT] === "function") await handlers[INIT](ctx)
  __loop: while (1) {
    const letter = await ctx.receive() // 接收订阅的 send 消息
    try {
      if (letter.tag === EXIT) {
        // 退出处理函数
        if (typeof handlers[TERMINATE] === "function") {
          await handlers[TERMINATE](ctx, letter, letter.data || {})
        }
        break __loop
      }
      await handlers[letter.tag](ctx, letter, letter.data || {}) // 根据接受的消息做不同的函数处理
    } catch (error) {
      console.error(`${ctx.self()} Error`, letter, error)
    } finally {
      continue __loop
    }
  }
}
```

其实这里维护了一个全局的可以通过事件监听与触发的状态，用来保存维护用户的信息，方便第三方的应用或是钱包服务商将服务整合进来。

#### 初始化第三方钱包界面渲染

接下来是初始化第三方钱包的服务，这里用到了 `iframe` 和 `fcl` 之前的配置信息，这也是为什么我们需要在开发环境启动一个 local 的 dev-wallet 服务的原因，交易的签名和确认会通过 `iframe` 的形式渲染出钱包服务的界面，完成用户的登录或授权。

```ts
// 通过 iframe 的方式呼出授权界面
 const [$frame, unrender] = renderAuthnFrame({
 handshake: await config().get("challenge.handshake"), // 在 fcl 配置中配置的第三方授权服务地址，localhost 是依赖 dev-wallet 测试网依赖 blocto
 l6n: window.location.origin,  // 当前页面的 url
 })


// packages/fcl/src/current-user/render-authn-frame.js
import {renderFrame} from "./render-frame"

export function renderAuthnFrame({handshake, l6n}) {
  var url = new URL(handshake)
  url.searchParams.append("l6n", l6n)
  return renderFrame(url.href)
}

// packages/fcl/src/current-user/render-frame.js

const FRAME_ID = "FCL_IFRAME"

export function renderFrame(src) {
  if (document.getElementById(FRAME_ID)) return

  const $frame = document.createElement("iframe")
  $frame.src = src
  $frame.id = FRAME_ID
  $frame.allow = "usb *"
  $frame.frameBorder = "0"
  $frame.style.cssText = `
    position:fixed;
    top: 0px;
    right: 0px;
    bottom: 0px;
    left: 0px;
    height: 100vh;
    width: 100vw;
    display:block;
    background:rgba(0,0,0,0.25);
    z-index: 2147483647;
    box-sizing: border-box;
  `
  document.body.append($frame)

  const unmount = () => {
    if (document.getElementById(FRAME_ID)) {
      document.getElementById(FRAME_ID).remove()
    }
  }

  return [$frame, unmount]
}
```

这里是将在 fcl 全局配置的链接放入 iframe 中，并呼出其界面，加载第三方钱包提供的登录授权界面。

#### 定义全局响应函数

这个是为第三方钱包服务所准备的通信机制，通过 `message` 的方式让 iframe 之间互相通信，加载的第三方的钱包服务，可以在 iframe 提供的页面中实现自己的授权逻辑，将授权的获取的用户信息，通过事件的方式传递到上下文中。

```ts
 // 定义响应函数
const replyFn = async ({data}) => {
  if (data.type === CHALLENGE_CANCEL_EVENT || data.type === CANCEL_EVENT) { // 取消授权，关闭窗口，取消事件监听
    unrender()
    window.removeEventListener("message", replyFn)
    return
  }
  if (data.type !== CHALLENGE_RESPONSE_EVENT) return // 非登录响应的数据都返回
  // 正常的数据响应流程
  unrender()
  window.removeEventListener("message", replyFn)

  send(NAME, SET_CURRENT_USER, await buildUser(data)) // 根据返回的数据初始化用户信息，并设置给 ctx
  resolve(await snapshot()) // 返回用户信息
}
```

这里使用了 `buildUser` 将返回的 data 信息构建成 current-user 所需要的用户信息，和授权相关的服务信息，我们用 local 环境的 dev-wallet 举例，在 config 中我们配置了授权的信息，同样也会在 currentUser 的信息中展示：

```json
{
  "VERSION": "0.2.0",
  "addr": "01cf0e2f2f715450",
  "cid": "did:fcl:01cf0e2f2f715450",
  "loggedIn": true,
  "services": [
    { // 签名相关的服务接口
      "type": "authz",
      "keyId": 0,
      "id": "asdf8701#authz-http-post",
      "addr": "01cf0e2f2f715450",
      "method": "HTTP/POST",
      "endpoint": "http://localhost:8701/flow/authorize",
      "params": {
        "userId": "b4eebc63-8e2e-4166-9ca2-5ce97c5d078d"
      }
    },
    { // 登录授权相关的服务接口
      "type": "authn",
      "id": "wallet-provider#authn",
      "pid": "b4eebc63-8e2e-4166-9ca2-5ce97c5d078d",
      "addr": "asdf8701",
      "name": "FCL Dev Wallet",
      "icon": "https://avatars.onflow/avatar/asdf8701.svg",
      "authn": "http://localhost:8701/flow/authenticate"
    }
  ]
}
```

在之前与 flow 交互的文章中，我们其实略过了这部分的讲述，这里我们看到 services 里有两个配置：

- `authn` —— 定义了账户登录认证的接口与参数
- `authz` —— 定义了账户签名交易授权的参数

通过这两个配置，fcl 才能够在用户登录和发起签名操作的时候请求对应的授权接口，`buildUser` 中 通过 `fetchServices` 获得

```ts
// packages/fcl/src/current-user/build-user.js
// 构建用户信息
export async function buildUser(data) {
  data = normalizeData(data)
  // 合并服务信息，拉取授权接口信息
  var services = mergeServices(
    data.services || [],
    await fetchServices(data.hks, data.code) // 获取签名授权的信息 authz
  ).map((service) => normalizeService(service, data))

  // console.log("BUILD USER", services)

  const authn = findService("authn", services)
  // 返回填充后的数据
  return {
    ...USER_PRAGMA,
    addr: withPrefix(data.addr),
    cid: deriveCompositeId(authn),
    loggedIn: true,
    services: services,
    expiresAt: data.exp,
  }
}
```

到这里，授权登录整个流程就结束了，我们还顺带的熟悉了一下通过第三方服务授权的例子，下面我们看看 current-user 中其他的接口：

- `unauthenticate` // 用户注销
- `authorization`  *// 用户交易签名认证*
- `subscribe`  *// 事件监听*
- `snapshot`  *// 账户信息快照*

## nauthenticate

这里就是直接清除掉客户端全局缓存的 `root.FCL_REGISTRY` 信息，同时调用上下文对象中的订阅方法，通知订阅者用户注销

```ts
function unauthenticate() {
  spawnCurrentUser()
  send(NAME, DEL_CURRENT_USER)
}

// handler 
[DEL_CURRENT_USER]: async (ctx, letter) => {
    ctx.merge(JSON.parse(DATA))  // 使用默认空数据格式覆盖原有的缓存
    if (await canColdStorage()) coldStorage.put(ctx.all())
    ctx.broadcast(UPDATED, {...ctx.all()}) // 广播
  },

// @onflow/util-actor
broadcast: (tag, data, opts = {}) => {
  opts.from = addr
  for (let to of root.FCL_REGISTRY[addr].subs) send(to, tag, data, opts)
},
```

## authorization

这里实现了第三方托管类型的用户签名授权服务，将用户信息中定义的签名服务 `authz` 的信息转化为请求，需要签名的交易体作为 `signable` 参数传入，第三方授权服务返回签名后的交易信息和签名，完成托管服务的授权操作

```ts
async function authorization(account) {
  spawnCurrentUser()
  const user = await authenticate() // 获得当前登录用户信息
  const authz = serviceOfType(user.services, "authz") // 遍历匹配 authz 签名相关信息

  const preAuthz = serviceOfType(user.services, "pre-authz") // 适配预签名
  if (preAuthz) {
    return {
      ...account,
      tempId: "CURRENT_USER",
      async resolve(account, preSignable) {
        return rawr(await execService(preAuthz, preSignable))
      },
    }
  }

  return {
    ...account,
    tempId: "CURRENT_USER",
    resolve: null,
    addr: sansPrefix(authz.identity.address),
    keyId: authz.identity.keyId,
    sequenceNum: null,
    signature: null,
    async signingFunction(signable) {
      return execService(authz, signable) // 请求三方服务器签名
    },
  }
}

// packages/fcl/src/current-user/exec-service/index.js
import {execHttpPost} from "./strategies/http-post"
import {execIframeRPC} from "./strategies/iframe-rpc"

const STRATEGIES = {
  "HTTP/RPC": execHttpPost,
  "HTTP/POST": execHttpPost,
  "IFRAME/RPC": execIframeRPC,
}

export async function execService(service, msg) {
  try {
    return STRATEGIES[service.method](service, msg) // 根据 service 中 authz 的配置执行签名授权的请求
  } catch (error) {
    console.error("execService(service, msg)", error, {service, msg})
    throw error
  }
}
```

这里我们可以把第三方的服务看做是一个签名授权的钱包服务，通过用户的 `session` 信息认证用户的身份，同时将客户端提供的所需签名的数据进行签名，并返回签名信息，这个过程中不暴露用户的私钥，也不会增加用户的认知负担，降低使用区块链的门槛。当然我们也可以用同上的 fcl 接口实现本地私钥的签名授权，这个将会放到后面的文章详述。

## subscribe

通过 subscribe 完成用户信息变动的监听

```ts
function subscribe(callback) {
  spawnCurrentUser()  // 初始化
  const EXIT = "@EXIT" // 定义退出条件
  const self = spawn(async ctx => {  // 初始化一个新的上下文,作为消息传递的 from 方
    ctx.send(NAME, SUBSCRIBE) // 发送订阅消息
    while (1) { // 定义监听循环
      const letter = await ctx.receive() // 接收更新
      if (letter.tag === EXIT) {
        ctx.send(NAME, UNSUBSCRIBE) // 注意这里使用的是 current-user 的上下文订阅
        return
      }
      callback(letter.data)
    }
  })
  return () => send(self, EXIT) // 结束监听函数，这里使用的是新的上下文
}

// 这里订阅 handler
 [SUBSCRIBE]: (ctx, letter) => {
    ctx.subscribe(letter.from)
    ctx.send(letter.from, UPDATED, {...ctx.all()})
  },
```

$subscribe$ 会为每个调用他的函数生成一个新的上下文，并添加到 `current_user` 的订阅中，完成订阅的隔离，之前我们看到的截图中的 subs 列表中有 1 和 2 两个订阅 ，其实是因为应用注册了两个 `subscribe` 的函数，除了 currentUser 外对应也有两个相同 `addr` 的上下文对象。

## snapshot

snapshot 获取当前上下文对象中存储的用户信息，并返回

```ts
// current-user
function snapshot() {
  spawnCurrentUser()
  return send(NAME, SNAPSHOT, null, {expectReply: true, timeout: 0})
}

// handler
[SNAPSHOT]: async (ctx, letter) => {
  letter.reply({...ctx.all()})
},

// ctx
all: () => {
  // 返回存储的所有的用户数据
  return root.FCL_REGISTRY[addr].kvs
},
```

## 最后

我们可以看到，fcl 在全局维护了一个共享的上下文，同时也提供了上下文之间通信的工具与监听模式，借助 iframe 嵌套第三方页面的形式完成三方可控的网页授权操作，使用 service 定义授权和签名所需要的服务接口描述，第三方的托管钱包只需要关注服务本身，将授权获得用户的信息通过订阅消息的形式传递给 fcl ，同时在签名的过程中也无需暴露用户的私钥，将签名作为返回值回传给 fcl 即可完成用户的授权。

这么设计的好处显而易见，降低普通用户的使用门槛，同时帮助第三方更好的提供链上账户的授权和与 DApp 功能业务整合的支持。

当然还有非托管形式的签名方法，将会在后面的源码解析文章中详述，本文涉及到的源码和注释在 [github](https://github.com/caosbad/flow-js-sdk/tree/comm) 中你也可以使用 [github1s](https://github1s.com/caosbad/flow-js-sdk/tree/comm) 来查看。

`2021-02-15@Caos`