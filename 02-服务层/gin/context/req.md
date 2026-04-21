req,  c中的req，有关请求信息的内容以及处理方法

已经拿到req了，我们的任务就是 解析 请求头和请求体，提取出我们需要的信息。
请求头：包含请求方法、路径、协议版本、请求头字段等。
请求体：包含实际的请求数据，如表单数据、JSON数据等。

请求头：主要就是请求头类型
c.Request.Header 获取所有请求头字段
核心方法：value:=c.GetHeader("key") 获取指定请求头字段的value

常用快捷获取（高频）
获取认证头：GetHeader ("Authorization")
获取内容类型：GetHeader ("Content-Type")
获取客户端标识：GetHeader ("User-Agent")

请求体：body就按照请求头的Content-Type来解析内容
实际就是 从流中读取数据     ，content-type就是文件流中的 .xxx 文件类型

流式读取？  按照什么格式读取到go内存中？  考虑好解析出来的数据是什么struct？

这里跟文件读取类似，我们只需要  数据绑定

数据绑定这里有一个最佳实践：
请求体复杂而且需要校验，使用ShouldBind系列  (需要写一个结构体)
请求体简单，不需要校验，使用对应的小工具：例如PostForm....（不需要写结构体，但是需要手动写出对应字段名）

1. 无敌通用的ShouldBindxxx(&data)
特点：自动+Tag+宽松
使用场景：需要解析的字段体比较复杂
使用方法： c.ShouldBindxxx(&data)

推荐使用ShouldBind系列，Bind系列错误直接阻断请求，太强硬了

Gin 可以从多种来源绑定数据：JSON、XML、YAML、TOML、表单数据（URL 编码和 multipart）、查询字符串、URI 参数和请求头。使用相应的结构体标签（json、xml、yaml、form、uri、header）来映射字段。验证规则放在 binding 标签中
原理其实就是 自动读取请求头+按照结构体的验证规则，将请求体中的数据解析到结构体中。

binding 标签最常用的校验参数
```
一、基础必传校验
required：字段必填，不能为空
二、数值范围校验
gte：大于等于指定值gt：大于指定值lte：小于等于指定值lt：小于指定值
三、字符串长度校验
len：字符串长度必须等于指定值min：字符串最小长度max：字符串最大长度
四、格式校验
email：必须是邮箱格式uri：必须是网址格式uuid：必须是 UUID 格式
五、枚举校验
oneof：值必须是指定枚举中的一个
```
使用示例代码：
```go
type LoginForm struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

r.POST("/login", func(c *gin.Context) {
    var form LoginForm
    
    // 绑定JSON数据到结构体
    if err := c.ShouldBindJSON(&form); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 处理登录逻辑...
    
    c.JSON(200, gin.H{
        "message": "登录成功",
        "username": form.Username,
    })
})
```


2. 精细化的小工具
特点：不需要struct，但是需要手动写出对应字段名
使用场景：需要解析的字段体比较简单
使用方法：
value:=c.PostForm("key")

username := c.PostForm("username")
password := c.DefaultPostForm("password", "default_password")

file, err := c.FormFile("file")
form, err := c.MultipartForm()
files := form.File["files"]
FormFile：专门获取单个文件，针对单文件上传场景
MultipartForm：解析完整的多部件表单，支持多个文件上传，同时可获取表单内的所有文本参数







