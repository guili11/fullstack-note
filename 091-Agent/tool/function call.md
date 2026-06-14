# 1. 定义（是什么 / 边界）

Function Calling 是一种机制，用于让模型输出结构化的，函数调用请求（函数名 + 参数）。

# 2. 解决问题（为什么存在）

关键：决策和执行的分离。

标准解析：

Function Calling 解决的问题是：

LLM能够理解用户意图并做出工具调用决策，
但无法直接与外部系统交互并执行操作，
导致“决策能力”和“执行能力”脱节。

人话： 实际就是大脑和手脚的关系
llm能判断做什么，但是不能做。
系统能做，但不知道要做什么。     

# 3. 思想（用什么方式解决）      
注意，思想要聚焦的是 问题的关键点，另外的补丁，我们不需要管，那个是实现处理的。

这里问题关键是 决策和执行的断裂，
所以我们的思想应该是 聚焦如何修复这种断裂：
用结构化输出作为AI与系统之间的执行协议

> 至于 如何解析，如何执行，这是建模需要关注的问题。

# 4. 实现 / 建模（如何落地成系统）

建模：
把解决方案拆解成系统结构
（流程、角色、数据、职责）

流程：
Tool Schema（工具定义）
↓
用户请求
↓
LLM

（参考 Tool Schema）

↓
Tool Call
{
    name,
    arguments
}
↓判断是否需要调用工具，并路由到对应工具
Executor
↓
真实函数
↓
执行结果
↓append到 context
LLM
↓
最终回答

接下来分 角色来划分职责：
系统：注册tool信息，让llm识别到工具的存在。
LLM：理解用户意图，生成 Tool Call。根据执行结果，生成最终回答。   
Executor：识别到 Tool Call，然后 路由到真实函数，执行操作。
真实函数：执行实际的操作，返回结果。

距离最终的实现，要解决三个问题：

**A. 如何让 LLM 产生正确 Tool Call**
- Registry — 收集管理所有工具定义，供 LLM 知道有哪些工具可用
- Tool Schema — 用 JSON Schema 描述工具的签名（名称、参数、类型），让 LLM 理解调用规则
- Prompt 注入 — 将工具定义写入 System Prompt，通过上下文告知 LLM 何时调用
- Function Calling API — 利用 LLM 原生 `tools` 参数，让模型自动生成结构化 Tool Call

**B. 如何识别 Tool Call**
- 决策解析 — 判断 LLM 是否要调工具、调哪个、参数是什么。早期需要手动从文本中解析 JSON，现在 Function Calling API 直接返回结构化 `tool_calls`
- JSON 解析 — 从 LLM 返回值中提取结构化的函数名和参数
- 参数校验 — 校验参数完整性、类型正确性、值有效性，防止 LLM 幻觉

**C. 如何执行 Tool Call**
- Dispatcher — 根据工具名称路由到对应的注册工具，类似 Web 路由
- Executor — 封装执行流程：超时控制、错误处理、重试、结果格式化

# 5. 未解决的问题（遗留问题 / 下一层问题）

Function Calling 解决了什么，但又没解决什么？

或者说，function call 实践中需要注意的核心问题是？


Function Calling 的遗留问题 = “从自然语言到结构化指令的可靠性问题”
  
1. 结构化输出不稳定   检测 重试机制   （已解决）
2. 工具选择依赖上下文，不是强约束    需要 指令约束来弥补 llm的意图概率识别
3. 模型意图与系统执行结果可能不一致

---

# 6. 实践（记录最新的开发方式）

> 记录当前工程中推荐使用的工具和模式，持续更新。

## 6.1 决策解析 —— 直接使用 Function Calling API

当前主流 SDK（OpenAI / Claude / Gemini）已内置 Tool Call 解析能力，无需手写 JSON 解析。

```
传统方式（Prompt 注入）:
  LLM 返回文本
    → 你手写正则/JSON 解析  ← 你自己干了解析的活
    → 拿到 {tool, arguments}
    → 这时的"决策解析" = JSON 解析

现在的 Function Calling API:
  LLM 返回 response
    → OpenAI SDK 内部封装了解析逻辑  ← 工具帮你干了
    → 你直接读 message.tool_calls
    → 这时的"决策解析" = 判断 tool_calls 是否为 null
```

演进路径对比：

| 阶段 | 方式 | 谁做解析 | 可靠性 | 适用范围 |
|------|------|---------|--------|---------|
| 1. 手工 Prompt | System Prompt 描述工具 → 解析文本 | 开发者手写正则/json.loads | 低，LLM 格式不稳定 | 早期、不支持 FC 的开源模型 |
| 2. 半结构化 | JSON Mode + Prompt 约束 | 开发者 + 平台辅助 | 中 | 部分支持约束的模型 |
| 3. Function Calling API | 传入 `tools` 参数，直接拿 `tool_calls` | OpenAI SDK 内置 | 高，结构化输出稳定 | OpenAI、Claude、Gemini 等主流模型 |

**核心洞察：不是理解错了，是工具进步了。** JSON 解析的能力没有消失，而是被下沉到 SDK 内部。遇到不支持 Function Calling 的模型时，依然需要回到传统方式。

## 6.2 职责划分 —— 开发者应该关注什么

在成熟的 Function Calling 体系中，职责划分非常清晰：

```
LLM          → 决策调用什么工具（这是 LLM 唯一做的事）
框架 / SDK   → 工具注册、工具注入、Tool Call 解析、工具路由、参数校验、结果回注
开发者       → 关注工具本身的能力定义和业务逻辑实现
```

开发者不需要重复搭建整套工具调用基础设施：

```python
# ❌ 开发者不需要关心这些（框架/SDK 帮你做了）
# - 如何把 tools 列表传给 LLM
# - 如何从 response 中提取 tool_calls
# - 如何校验参数类型
# - 如何把执行结果回传给 LLM

# ✅ 开发者只需要关心这个
def get_weather(city: str, date: str = None) -> str:
    """查询天气的实际逻辑"""
    return f"{city} 的天气是晴天，25°C"

def send_email(to: str, subject: str, body: str) -> str:
    """发送邮件的实际逻辑"""
    return f"邮件已发送至 {to}"
```

**实践建议：**
- 选择支持 Function Calling 的 SDK 和模型（OpenAI / Claude / Gemini），让框架承担调用链路
- 把精力花在工具的 Schema 设计（参数描述是否清晰）和业务逻辑的正确性上
- 只在遇到不支持 FC 的开源模型时，才自行搭建 JSON 解析和注入链路
