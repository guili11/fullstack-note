使用air热更新服务器

**一个服务器要做的核心工作**

1. 请求解析

读取客户端请求

解析 HTTP Method、URL、Header、Body 等信息

处理跨域、Cookie、Session 等

2. 路由匹配

根据请求 URL + Method 找到对应的 Handler / 控制器

支持参数提取、REST 风格接口

3. 执行业务逻辑（Handler 服务）

调用 Service 层 / 业务逻辑

对请求数据进行处理，生成结果

4. 响应创建与返回

封装 HTTP 状态码、Header、Body

支持 JSON / HTML / 文件等多种格式

发送给客户端

---

Gin 的优势体现

封装 Context：统一请求解析 + 响应管理，Handler 可以直接访问 c.Query()、c.Param()、c.JSON() 等

路由灵活：支持分组、RESTful 路由、动态参数

中间件机制：在 Handler 前后可统一处理日志、鉴权、限流等

高性能：底层使用 httprouter，路由匹配速度快

简化开发：很多重复操作（绑定、JSON 响应、错误处理）封装好，减少样板代码

---

💡 面试亮点表述：

> “一个 HTTP 服务器核心工作是请求解析、路由匹配、业务逻辑执行和响应返回。Gin 框架在这些环节提供了封装和中间件机制，极大提高了开发效率和可维护性，同时保证性能和灵活性。”

---

```go
package main

import (
    "fmt"
    "net/http"
)

// 模拟一个业务处理函数
func helloHandler(w http.ResponseWriter, r *http.Request) {
    // 原始请求解析
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "world"
    }

    // 原始响应写入
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, `{"message": "Hello %s"}`, name)
}

func main() {
    // 路由注册
    http.HandleFunc("/hello", helloHandler)

    // 启动服务器
    fmt.Println("Server running at http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Server error:", err)
    }
}
```

![](C:\Users\Administrator\AppData\Roaming\marktext\images\2026-05-09-17-25-52-image.png)