下面这份我按“从 0 到能做 Agent”的方式收敛成一条认知链路。重点不是把概念堆满，而是把 Tool 的边界、选择和错误处理放进同一套结构里理解。

# Tool Calling & Tool Design（Agent 核心笔记）

## 1. 核心认知

### Tool 的本质

```text id="a1k9xp"
Tool = 将模型的“语言能力”连接到“外部确定性能力”的接口
```

更工程一点：

```text id="b3m7vx"
Tool = 让 LLM 从“回答问题”变成“执行系统行为”的桥梁
```

### Tool System 的本质目标

```text id="o1v9md"
把一个不确定的语言模型，约束成一个“可执行系统”
```

### AI 在 Tool System 中的角色

```text id="p4k7qp"
AI：在边界内做概率推理
人：定义边界 + 系统结构
```

### 一句话总纲

```text id="r1v8md"
Tool system 的本质，是用人类定义的语义边界，去压缩模型的决策空间，从而换取系统行为的可控性与稳定性
```

## 2. Tool Calling 如何工作

### 选择与生成

LLM 做的不是“直接调用函数”，而是：

```text id="c6p2md"
在多个候选工具中，根据语义相关性选择一个，并生成参数
```

本质概率形式：

```text id="d8v1qx"
P(tool | query, tool_descriptions)
```

### 三层结构

```text id="e2k9mn"
① Tool Description（语义层）
② Tool Selection（决策层）
③ Tool Execution（执行层）
```

你现在主要在理解的是 ① 和 ②。

### Tool Description 的本质

```text id="f5q8wp"
Tool description = 对工具“能力语义”的定义，而不是执行说明
```

它的作用：

* 定义工具属于哪个“知识类别”
* 影响模型如何归类问题
* 影响 selection 概率分布

### Schema vs Description

```text id="g7m2vd"
Schema：约束输入/输出结构（硬约束）
Description：影响模型理解（软约束）
```

### Tool Boundary

```text id="h9k3xp"
Tool boundary = 人对系统能力的责任划分
```

它决定：

* 每个 tool 解决什么问题
* 不解决什么问题
* 出错归谁

### 为什么 Boundary 必须由人定义

```text id="i4v6mn"
AI只能做语义匹配，不能做系统设计决策
```

原因：

* 没有全局目标函数
* 只做局部概率判断
* 不理解业务 trade-off

### Tool Overlap

```text id="j6p1vq"
Tool overlap = 多个 tool 在语义上覆盖同一问题空间
```

后果：

* selection 不稳定
* 行为不可复现
* prompt 调优失效

本质问题：

```text id="k8m3wx"
破坏决策空间的“可分性”
```

### Description 不一致 vs Overlap

| 类型              | 影响     |
| --------------- | ------ |
| description 不一致 | 语义理解偏差 |
| tool overlap    | 决策结构崩坏 |

### Selection 的本质

```text id="n8k2vx"
Tool selection = 语义分类问题 + 概率决策问题
```

模型是在做：

* 这个问题属于哪个“能力类别”
* 哪个 tool 最接近

有些问题就是需要调用多个 tool 的。

Tool selection 的难点不在“能不能选对”，而在“有多个都能选的时候是否稳定”。

## 3. Tool Design 的目标与原则

### 设计目标

```text id="l2q9md"
Tool design 的目标不是“更聪明”，而是“更可控”
```

包括：

* 减少歧义
* 降低选择空间
* 提高行为可预测性

### 设计本质

你可以这样理解：

```text id="m5v7qp"
Prompt → 控制输出内容
Tool boundary → 控制行为路径
```

### 核心原则

```text id="q7m2xp"
① 单一职责（single responsibility）
② 最小重叠（minimal overlap）
③ 语义清晰（clear semantics）
④ 行为可预测（predictable behavior）
```

## 4. Tool 错误怎么处理

### 核心思想

在 Agent 里，Tool 的错误不应该一上来就由程序“硬处理完”，而是先把错误变成观测信息，再交给 LLM 决策下一步。

```text
Tool
↓
基础重试（程序处理）
↓
仍失败
↓
结构化错误返回
↓
LLM 决策是否继续、换 tool、降级或终止
```

### 为什么要这样做

* 程序层只负责确定性兜底，避免低级瞬时错误直接打断流程
* 结构化错误把失败原因变成可读信号，而不是纯文本噪音
* LLM 拿到错误上下文后，可以结合任务目标判断是否继续

### 这里的关键边界

```text id="s6k2qp"
Tool = 控制模型“能做什么”的结构，而不是“怎么回答”的技巧
```

错误处理也属于这个结构的一部分：程序负责观测与包装，模型负责基于观测做路径选择。

## 5. 往后怎么接

如果你后面继续往 Agent 学，这一套理解会直接延伸到：

* ReAct（在 tool 之上加循环）
* Workflow / Graph（控制执行流）
* MCP（工具标准化）
* Eval（评估 tool selection）

如果你下一步想继续，我可以帮你把这个体系往前推一层：

```text id="t4v9xn"
ReAct 为什么本质是在 tool system 上加“状态机”
```

那一层你就进入真正 Agent 结构了。
