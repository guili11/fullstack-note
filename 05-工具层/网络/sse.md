# 一、SSE 基础概念

## 1. 定义

SSE 是浏览器原生支持的一种 **服务器推送技术**

**本质**：浏览器通过 HTTP 长连接 接收后端连续推送的事件流

## 2. 特点

| 特性      | 描述                             |
| ------- | ------------------------------ |
| 单向      | 后端推送 → 前端接收，前端不能通过 SSE 发数据     |
| 长连接     | 建立一次连接，可以持续推送多条事件              |
| HTTP 原生 | 基于 GET 请求，EventSource API 原生支持 |
| 事件驱动    | 后端发送事件触发前端回调                   |

## 3. 典型场景

- 实时更新工作流节点状态（你的项目）
- 实时监控 / 日志推送
- 新闻或行情实时推送

***

# 二、SSE 如何使用

## 1️⃣ 前端

```javascript
const evtSource = new EventSource('/subscribeWorkflowProgress?applicationId=123')

evtSource.onmessage = (event) => {
    const data = JSON.parse(event.data)
    console.log('节点完成事件', data)
}

evtSource.addEventListener('workflow_done', (event) => {
    console.log('Workflow 完成', JSON.parse(event.data))
})

evtSource.onerror = (err) => {
    console.error('SSE 错误', err)
}
```

**API 说明**：

- `EventSource(url)` → 建立 SSE 连接
- `.onmessage` → 默认接收所有事件
- `.addEventListener('事件名')` → 指定事件类型
- `.close()` → 关闭连接

***

## 2️⃣ 后端

**设置响应头**：

```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
```

**发送事件**：

```
event: node_complete
data: {"node":"preprocess","status":"done"}

```

每条事件以 `event:` + `data:` + 两个换行结尾

***

# 三、SSE 核心原理

## 1. 长连接

- 建立一次 TCP 连接，后端持续写入数据流
- 浏览器 EventSource 自动保持连接

## 2. 事件流

- 后端每发一条事件，浏览器触发对应回调
- 支持事件类型和数据分离 → 可区分 `node_complete`、`workflow_done`、`workflow_error`

## 3. 自动重连

- EventSource 会在连接断开后自动重连
- 可以通过 `retry:` 指定重连间隔

***

# 四、SSE vs HTTP Stream vs WebSocket

| 特性    | SSE               | HTTP Stream | WebSocket          |
| ----- | ----------------- | ----------- | ------------------ |
| 方向    | 单向                | 单向          | 双向                 |
| 建立方式  | GET + EventSource | fetch/axios | `ws://` 或 `wss://` |
| 数据格式  | 文本事件              | 任意 chunk    | 任意                 |
| 浏览器支持 | 原生                | 原生          | 原生/库               |
| 适用场景  | 实时事件驱动            | 纯数据流        | 双向交互               |
| 复杂度   | 低                 | 中           | 高                  |

> **你的项目适合 SSE**：Workflow 事件流 单向、事件驱动、浏览器直接订阅

***

# 五、设计原因 & 好处

## 1. 事件驱动契合业务

- Workflow 每个节点完成 → 发事件
- 前端按事件更新状态 → 无需轮询 → 高效

## 2. 单向流降低复杂度

- WebSocket 双向没必要
- SSE 自带重连机制 → 简单、可靠

## 3. 浏览器原生支持

- 不依赖第三方库
- 易于集成 Vue3 Composition Hook

***

# 六、进阶点 / 面试亮点

- **SSE 与 gRPC 流结合**：后端异步执行节点 → gRPC Stream → SSE → 前端
- **支持多事件类型**：`node_complete`、`workflow_done`、`workflow_error`
- **前端状态建模**：`steps` + `currentNode` + `workflowStatus` → 可动态渲染分支
- **生命周期管理**：页面卸载 → `evtSource.close()`，避免资源泄漏
- **可扩展性**：未来多实例 → 可换 Redis Pub/Sub + SSE 代理

***

## 💡 一句话总结 SSE

> SSE 是 **事件驱动、单向、浏览器原生支持**的长连接，非常适合 Workflow 节点异步执行状态的实时推送，结合 gRPC 流可以形成端到端可靠的事件链路。

***

