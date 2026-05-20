从需求  一步步 引入到  为什么要使用ctx，ctx本质是什么。



当我们想在 在函数之间传递信息的时候，原本是使用 参数传递，  当信息过多，这个参数还需要额外维护，

所以不如直接维护一个ctx 全局变量，各个函数，协程访问

由于ctx本质就是一块内存 全局，

- **Go 原生 `ctx` 的唯一性**：只作用于 **「单次请求内部」**
  
  不跨服务器、不跨服务、不跨集群！

- **分布式全局唯一性**：不是 `ctx` 本身，是 **存在 ctx 里的 `TraceID`**
  
  这个 ID 才是 **全集群、全服务、全链路唯一** 的！



ctx的概念？  go context.Contect()
context.Context 是 Go 提供的一种跨 API 边界、跨 goroutine 传递控制信息 和共享数据 的机制。
它主要解决三个问题：
取消控制（cancellation）：上游请求取消时，下游可以及时停止耗时操作。
超时/截止时间（timeout / deadline）：给请求链路统一的超时限制，防止无限等待。
跨 API 传递请求范围信息（key/value）：比如 traceID、用户信息等，可在整个请求链路中透传递

实践范例（按照业务粒度 划分    后续可以升级为按照基础操作粒度的）
handler 入口：创建带超时的 ctx（例：3-10 秒，按接口粒度）  cancel是用于在操作没超时，对于ctx的销毁处理
service 层：接收并透传 ctx，不自己新建
repo / utils / 外部服务调用：在耗时操作内部介入ctx ，在耗时操作后捕获err返回给handler
顶层 handler：统一捕获 context.DeadlineExceeded 并返回超时响应

ctx的核心点：
ctx超时创建（一定要在 handler 入口 使用 context.WithTimeout 创建 目的：保证整条请求链都能被取消或超时控制）
ctx超时/取消控制（耗时操作介入ctx）
ctx错误捕获（操作后 捕获）
ctx json返回  

不要在耗时操作后添加 ctx检查，耗时操作本身就应该支持ctx。  例如DB.WithContext()

基础操作粒度的 超时控制 ，跟上述流程基本一致，只是 ctx的创建，从 业务入口创建 修改为了在耗时操作前分别创建对应的ctx

为什么ctx创建要放在 handler呢（保障用户 取消请求的时候，能够捕获到    区别于服务本身的超时控制）
