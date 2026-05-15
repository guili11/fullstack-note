好问题！站在面试官角度，针对你这个 workflow 项目，我推测会问这些：

***

## 第一层：基础理解（看你是否自己写的）

> **"你介绍一下你这个工作流的整体设计吧"**

要能说清楚：

- 用了 Eino 的 `compose.NewGraph` 构建 DAG
- 6 个节点 + 1 个分支条件
- N1 纯数据查询，N2/N3/N4 用 ChatModel，N5 持久化
- 异步触发，不阻塞投递请求

***

## 第二层：技术细节（看深度）

### 1. 分支条件    --确定性高

> **"为什么用** **`GraphBranch`** **而不是** **`ChatModelAgent`？"**

你的回答思路：

- `ChatModelAgent` 是黑盒，LLM 自主决定调什么工具
- 但\*\*"简历是否完整 → 走哪条路"\*\*是确定的业务规则，不应该让 LLM 决定
- `GraphBranch` 是白盒分支，逻辑可控、可测试

### 2. Graph 的 `Invoke` 模式

> **"Graph 的执行模式是什么？AllPredecessor 和 AnyPredecessor 有什么区别？"**
AddEdge 方法默认是 `AllPredecessor`模式
- 你的 DAG 是串行的，默认 `AllPredecessor`（所有前驱完成才触发当前节点）
- `AnyPredecessor` 用于并行场景（比如多个来源汇聚到一个节点）

### 3. 异步触发的方式

> **"为什么用 goroutine 异步触发？有什么问题？"**

你的思路：

- 投递请求不应该等 AI 审核完才返回，用户体验差
- goroutine 的问题是：**服务重启后未完成的工作流会丢失**
- 改进方向：用消息队列（如 Redis Stream / RabbitMQ）保证可靠投递

### 4. 状态共享

> **"为什么用** **`ApplyState`** **共享状态？多个节点并发会有什么问题？"**

- 你的 DAG 是串行的，没有并发，所以单指针传递没问题
- 如果以后要并行（N3 match 和 N4 同时跑），就需要用 `sync.Mutex` 或不可变数据结构

***

## 第三层：架构设计（看全局视野）

> **"这个工作流和你们的 ChatModelAgent（AI 对话）有什么区别？"**

对比：

- **ChatModelAgent**：LLM 主动思考+调工具，适用开放场景
- **Graph DAG**：固定流程 + 条件分支，适用确定性场景

> **"如果 HR 想修改审核规则（比如改评分权重），怎么改？"**

- 当前是写死在 Prompt 里的，修改需要改代码
- 改进方向：Prompt 模板化 + 配置文件 + 规则引擎

> **"为什么不直接用 LangChain？"**

- 项目是 Go 写的，Eino 是 Go 生态
- Eino 设计理念和 LangChain 对齐（Graph/Chain/Tool 等概念）

***

## 第四层：极限场景（看应变能力）

> **"AI 调用超时了怎么办？"**

- 当前没有超时控制，会一直卡住
- 改进：context.WithTimeout 加超时，超时后走降级路径

> **"用户重复投递怎么办？"**

- 可以通过 `UNIQUE KEY(position_id, candidate_id)` 防重
- workflow 前加一层幂等校验

> **"workflow 跑到一半服务崩了，恢复后怎么处理？"**

- 当前 goroutine 方式会丢失
- 改进：Eino 的 `CheckPointStore` + 消息队列可以恢复

***

## 总结：简历上的加分表述

```
基于 Eino 框架实现异步投递审核工作流（DAG），核心能力：
1. AI 评估简历完整性 + 条件分支分流
2. AI 匹配岗位评分（Prompt Engineering）
3. 异步 goroutine 触发，不阻塞主流程
```

面试时把这个逻辑讲清楚，再加一句 **"不足之处是目前没有做中断恢复，后续可以用消息队列优化"**，就是满分答案。
