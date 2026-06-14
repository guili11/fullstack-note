最终的实现，要解决三个问题：

**A. 如何让 LLM 产生正确 Tool Call**
- Registry
- Tool Schema
- Prompt 注入
- Function Calling API

**B. 如何识别 Tool Call**
- JSON 解析
- 参数校验

**C. 如何执行 Tool Call**
- Dispatcher
- Executor

---

## A. 如何让 LLM 产生正确 Tool Call

### Registry（注册中心）

**概念：** 收集和管理所有可用的工具定义，在调用 LLM 之前将工具列表传递给模型，让 LLM 知道有哪些工具可用。

```python
from openai import OpenAI

client = OpenAI()

# ===== 1. 先定义好两个真实函数（业务逻辑） =====

def get_weather(city: str, date: str = None) -> str:
    """查询天气的实际逻辑"""
    print(f"   [执行] get_weather(city={city}, date={date})")
    return f"{city} 的天气是晴天，25°C"

def send_email(to: str, subject: str, body: str) -> str:
    """发送邮件的实际逻辑"""
    print(f"   [执行] send_email(to={to}, subject={subject})")
    return f"邮件已发送至 {to}"

# ===== 2. 把函数注册到一个字典（Registry） =====

tool_functions: dict[str, callable] = {
    "get_weather": get_weather,
    "send_email": send_email,
}

# ===== 3. 把函数描述作为 tools 参数传给 OpenAI =====

response = client.chat.completions.create(
    model="gpt-4o",
    messages=[
        {"role": "system", "content": "你是一个有用的助手"},
        {"role": "user", "content": "北京今天天气怎么样？"},
    ],
    tools=[  # ← 这就是 Registry 的最终产物：tools 列表
        {
            "type": "function",
            "function": {
                "name": "get_weather",
                "description": "获取指定城市的天气信息",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "city": {"type": "string", "description": "城市名称"},
                        "date": {"type": "string", "description": "日期，格式 YYYY-MM-DD"},
                    },
                    "required": ["city"],
                },
            },
        },
        {
            "type": "function",
            "function": {
                "name": "send_email",
                "description": "发送邮件",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "to": {"type": "string", "description": "收件人地址"},
                        "subject": {"type": "string", "description": "邮件主题"},
                        "body": {"type": "string", "description": "邮件正文"},
                    },
                    "required": ["to", "subject", "body"],
                },
            },
        },
    ],
)

# Registry 的价值：把上面的 tools 列表动态组装出来，不用手写
```

如果想把 tools 列表动态组装，只需要一个简单函数：

```python
def build_tools_list() -> list[dict]:
    """动态构建 tools 列表"""
    tools = []
    for name, func in tool_functions.items():
        # 从函数签名和文档里自动提取 schema（这里简化，实际可用 pydantic 或 inspect）
        schema = TOOL_SCHEMAS[name]  # 预先定义好的 schema
        tools.append({
            "type": "function",
            "function": {
                "name": name,
                "description": func.__doc__ or "",
                "parameters": schema,
            },
        })
    return tools

# 调用时
# response = client.chat.completions.create(tools=build_tools_list(), ...)
```

### Tool Schema（工具描述）

**概念：** 用 JSON Schema 格式描述每个工具的调用方式——函数名、功能描述、参数列表和类型，LLM 根据这个 Schema 来决定调哪个工具、传什么参数。

```python
# Tool Schema 就是传给 OpenAI 的 tools 参数里的那一段 JSON
# 它定义了 LLM 可以调用什么、怎么调

get_weather_schema = {
    "type": "object",
    "properties": {
        "city": {
            "type": "string",
            "description": "城市名称，如 北京、上海",
        },
        "date": {
            "type": "string",
            "description": "日期，格式 YYYY-MM-DD，不传则查当天",
        },
    },
    "required": ["city"],  # city 是必填，date 选填
}

send_email_schema = {
    "type": "object",
    "properties": {
        "to": {
            "type": "string",
            "description": "收件人邮箱地址",
        },
        "subject": {
            "type": "string",
            "description": "邮件主题",
        },
        "body": {
            "type": "string",
            "description": "邮件正文内容",
        },
    },
    "required": ["to", "subject", "body"],
}

# 实际使用时，每个函数对应一个 dict 放入 tools 列表
# LLM 返回时，会按这里的定义来填充参数
```

Schema 中的关键字段：

| 字段 | 作用 | 说明 |
|------|------|------|
| `type: "object"` | 参数是一个 JSON 对象 | 固定写法 |
| `properties` | 定义每个参数 | 名称、类型、描述 |
| `required` | 标记必填参数 | 不在这个列表里就是选填 |
| `description` | 告诉 LLM 这个参数是什么 | 写清楚 LLM 才能正确填值 |

### Prompt 注入（将工具描述写到 Prompt 里）

**概念：** 不依赖 OpenAI 的 `tools` 参数，而是在 System Prompt 里用文字描述工具，让 LLM 按约定格式返回 Tool Call。

```python
# 适用场景：使用的 LLM 不支持 Function Calling API（如某些开源模型）

system_prompt = """
你是一个智能助手，你可以通过以下工具帮助用户：

工具列表：
1. get_weather(city: str, date?: str) — 获取指定城市的天气信息
2. send_email(to: str, subject: str, body: str) — 发送邮件

调用规则：
- 当用户的需求匹配某个工具时，返回 JSON 格式的工具调用
- 格式：{"tool": "工具名", "arguments": {"参数名": "值"}}
- 一次只调用一个工具
- 参数缺少时向用户追问
"""

# 调用（传入普通 messages，不传 tools 参数）
response = client.chat.completions.create(
    model="gpt-4o",
    messages=[
        {"role": "system", "content": system_prompt},
        {"role": "user", "content": "给 admin@example.com 发一封测试邮件"},
    ],
)

# LLM 会返回这样的文本（而不是原生 tool_calls）：
# {"tool": "send_email", "arguments": {"to": "admin@example.com", "subject": "测试邮件", "body": "这是一封测试邮件"}}
```

### Function Calling API（原生工具调用）

**概念：** OpenAI 等模型提供的原生功能——API 请求中传入 `tools` 参数，模型自动判断是否调用工具，返回结构化 `tool_calls`，不需要手写 Prompt 规则。

```python
from openai import OpenAI
import json

client = OpenAI()

# 1. 定义工具描述
tools = [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "获取指定城市的天气信息",
            "parameters": {
                "type": "object",
                "properties": {
                    "city": {"type": "string", "description": "城市名称"},
                },
                "required": ["city"],
            },
        },
    }
]

# 2. 发请求（传入 tools，不用在 Prompt 里写规则）
response = client.chat.completions.create(
    model="gpt-4o",
    messages=[
        {"role": "system", "content": "你是一个有用的助手"},
        {"role": "user", "content": "北京今天天气怎么样？"},
    ],
    tools=tools,
    tool_choice="auto",  # auto=LLM 自主决定；"required"=强制调用；{"type":"function","function":{"name":"get_weather"}}=指定工具
)

# 3. 检查 LLM 是否返回了 Tool Call
message = response.choices[0].message

if message.tool_calls:
    # 直接拿到结构化数据，不用自己解析文本
    tool_call = message.tool_calls[0]
    func_name = tool_call.function.name
    func_args = json.loads(tool_call.function.arguments)
    print(f"LLM 选择了工具: {func_name}")
    print(f"参数: {func_args}")
else:
    print(f"LLM 直接回复: {message.content}")
```

---

## B. 如何识别 Tool Call

### JSON 解析

**概念：** 从 LLM 返回的消息中提取函数名和参数——对于原生 `tool_calls` 直接读取即可；对于 Prompt 注入模式，需要从返回文本中解析 JSON。

```python
import json
import re

# 场景 1：使用 Function Calling API（推荐）
# response 中直接有 message.tool_calls，无需额外解析

response = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "北京天气"}],
    tools=[...],
)
message = response.choices[0].message

if message.tool_calls:
    tc = message.tool_calls[0]
    name = tc.function.name           # "get_weather"
    args = json.loads(tc.function.arguments)  # {"city": "北京"}
    print(f"工具: {name}, 参数: {args}")


# 场景 2：Prompt 注入模式（LLM 返回纯文本）
# 需要从文本中提取 JSON

def parse_tool_call(text: str) -> dict | None:
    """从 LLM 返回文本中解析 Tool Call"""
    # 尝试直接 JSON 解析
    try:
        data = json.loads(text)
        if "tool" in data and "arguments" in data:
            return data
    except json.JSONDecodeError:
        pass

    # 尝试从 ```json 代码块中提取
    match = re.search(r'```(?:json)?\s*(\{.*?"tool".*?\})\s*```', text, re.DOTALL)
    if match:
        return json.loads(match.group(1))

    return None

# 使用
result = parse_tool_call('{"tool": "get_weather", "arguments": {"city": "北京"}}')
print(result)  # {'tool': 'get_weather', 'arguments': {'city': '北京'}}
```

### 参数校验

**概念：** 校验 LLM 传的参数是否合法——必填参数有没有缺失、参数类型对不对、枚举值是否在范围内。因为 LLM 可能产生幻觉，传了不存在的参数或错误的值。

```python
def validate_arguments(func_name: str, arguments: dict) -> tuple[bool, str]:
    """
    校验 LLM 生成的参数是否合法
    返回 (是否合法, 错误信息)
    """
    # 每个工具预定义的参数规则
    schemas = {
        "get_weather": {
            "properties": {
                "city": {"type": "string"},
                "date": {"type": "string"},
            },
            "required": ["city"],
        },
        "send_email": {
            "properties": {
                "to": {"type": "string"},
                "subject": {"type": "string"},
                "body": {"type": "string"},
            },
            "required": ["to", "subject", "body"],
        },
    }

    schema = schemas.get(func_name)
    if not schema:
        return False, f"未知工具: {func_name}"

    props = schema["properties"]
    required = schema.get("required", [])

    # 1. 检查必填参数
    for field in required:
        if field not in arguments:
            return False, f"缺少必填参数: {field}"

    # 2. 检查是否有未知参数（LLM 可能多传）
    for field in arguments:
        if field not in props:
            return False, f"未知参数: {field}"

    # 3. 类型检查
    type_map = {"string": str, "number": (int, float), "integer": int, "boolean": bool}
    for field, value in arguments.items():
        expected = props[field].get("type")
        if expected and not isinstance(value, type_map.get(expected, object)):
            return False, f"参数 {field} 应为 {expected} 类型"

    return True, ""


# 使用
valid, err = validate_arguments("get_weather", {"city": "北京"})
print(valid, err)  # True ""

valid, err = validate_arguments("get_weather", {})  # 缺 city
print(valid, err)  # False "缺少必填参数: city"

valid, err = validate_arguments("get_weather", {"city": "北京", "foo": "bar"})  # 多传了未知参数
print(valid, err)  # False "未知参数: foo"
```

---

## C. 如何执行 Tool Call

### Dispatcher（调度器）

**概念：** 根据解析出的工具名称，找到对应的真实函数来执行，相当于一个 `工具名 → 函数` 的路由表。

```python
# ===== 真实函数 =====

def get_weather(city: str, date: str = None) -> str:
    return f"{city} 的天气是晴天，25°C"

def send_email(to: str, subject: str, body: str) -> str:
    return f"邮件已发送至 {to}"

# ===== 路由表（Dispatcher 的核心） =====

tool_functions = {
    "get_weather": get_weather,
    "send_email": send_email,
}

# ===== 派发执行 =====

def dispatch(func_name: str, arguments: dict) -> str:
    """
    根据工具名找到对应的函数并执行
    """
    func = tool_functions.get(func_name)
    if not func:
        raise ValueError(f"未找到工具: {func_name}")
    return func(**arguments)  # 解包参数，直接调用


# ===== 完整使用 =====

# 假设 LLM 返回了这个 Tool Call
tool_call = {
    "name": "get_weather",
    "arguments": {"city": "北京"},
}

# 校验 → 派发 → 执行
valid, err = validate_arguments(tool_call["name"], tool_call["arguments"])
if valid:
    result = dispatch(tool_call["name"], tool_call["arguments"])
    print(result)  # 输出: 北京的天气是晴天，25°C
```

### Executor（执行器）

**概念：** 在 Dispatcher 基础上再封装一层，处理超时、错误捕获、重试等边界情况，让调用方不用关心这些细节。

```python
import time
import threading

def execute_tool(func_name: str, arguments: dict, timeout: float = 10.0) -> dict:
    """
    安全的工具执行器
    
    返回统一格式：
    {"success": true/false, "result": "...", "error": "..."}
    """
    try:
        func = tool_functions.get(func_name)
        if not func:
            return {"success": False, "error": f"未找到工具: {func_name}"}

        # 带超时的执行
        result_container = []
        error_container = []

        def target():
            try:
                result_container.append(func(**arguments))
            except Exception as e:
                error_container.append(e)

        thread = threading.Thread(target=target)
        thread.start()
        thread.join(timeout=timeout)

        if thread.is_alive():
            return {"success": False, "error": f"执行超时（{timeout}s）"}

        if error_container:
            return {"success": False, "error": str(error_container[0])}

        return {"success": True, "result": result_container[0]}

    except Exception as e:
        return {"success": False, "error": str(e)}


# ===== 使用示例 =====

result = execute_tool("get_weather", {"city": "上海"})
if result["success"]:
    print(f"执行成功: {result['result']}")
else:
    print(f"执行失败: {result['error']}")


# 带重试的执行器
def execute_with_retry(func_name: str, arguments: dict, max_retries: int = 2) -> dict:
    """执行失败时自动重试"""
    for attempt in range(max_retries + 1):
        result = execute_tool(func_name, arguments)
        if result["success"]:
            return result
        print(f"第 {attempt + 1} 次尝试失败，{'重试中...' if attempt < max_retries else '放弃'}")
        time.sleep(1)
    return result
```

---

## 完整流程串联（A → B → C 闭环）

```python
from openai import OpenAI
import json

client = OpenAI()

# ===== 1. 真实业务函数 =====
def get_weather(city: str, date: str = None) -> str:
    return f"{city} 的天气是晴天，25°C"

def send_email(to: str, subject: str, body: str) -> str:
    return f"邮件已发送至 {to}"

# ===== 2. Registry + Schema =====
tool_functions = {"get_weather": get_weather, "send_email": send_email}
tools_schema = [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "获取天气信息",
            "parameters": {
                "type": "object",
                "properties": {
                    "city": {"type": "string", "description": "城市名称"},
                },
                "required": ["city"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "send_email",
            "description": "发送邮件",
            "parameters": {
                "type": "object",
                "properties": {
                    "to": {"type": "string"},
                    "subject": {"type": "string"},
                    "body": {"type": "string"},
                },
                "required": ["to", "subject", "body"],
            },
        },
    },
]

# ===== 3. 一次完整的 Tool Call 生命周期 =====

def agent(user_input: str) -> str:
    """A → B → C 完整流程"""

    # ===== A. 让 LLM 产生 Tool Call =====
    response = client.chat.completions.create(
        model="gpt-4o",
        messages=[
            {"role": "system", "content": "你是一个有用的助手"},
            {"role": "user", "content": user_input},
        ],
        tools=tools_schema,
        tool_choice="auto",
    )

    message = response.choices[0].message

    # LLM 没调工具，直接返回文本
    if not message.tool_calls:
        return message.content

    # ===== B. 识别 Tool Call =====
    tc = message.tool_calls[0]
    func_name = tc.function.name
    func_args = json.loads(tc.function.arguments)

    print(f"[LLM 选择了] {func_name}({func_args})")

    # ===== C. 执行 Tool Call =====
    func = tool_functions.get(func_name)
    if not func:
        return f"错误：未知工具 {func_name}"

    try:
        result = func(**func_args)
        print(f"[执行结果] {result}")
    except Exception as e:
        result = f"执行错误: {e}"

    # ===== 将结果返回给 LLM 生成最终回答 =====
    response = client.chat.completions.create(
        model="gpt-4o",
        messages=[
            {"role": "system", "content": "你是一个有用的助手"},
            {"role": "user", "content": user_input},
            message,  # assistant 的 tool_calls 消息
            {
                "role": "tool",
                "tool_call_id": tc.id,
                "content": result,
            },
        ],
        tools=tools_schema,
    )

    return response.choices[0].message.content


# ===== 运行 =====
print(agent("北京今天天气怎么样？"))
# 输出: 北京的天气是晴天，25°C

print(agent("给 admin@test.com 发一封测试邮件，主题是 Hello"))
# 输出: 邮件已发送至 admin@test.com
```
