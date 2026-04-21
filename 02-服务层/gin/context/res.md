响应res的处理，
三点：
响应头，响应体，文件下载

# 响应头： 这里比较关键，因为大部分响应头的作用就是服务端来控制客户端的行为，比如缓存、压缩、跨域等。

核心方法：c.Header(key, value)：设置响应头的key-value对

例如：（具体细节见跨域.md）

1. 跨域： 服务端允许跨域行为：
   ```go
   c.Header("Access-Control-Allow-Origin", "*")
   ```
2. 缓存： 服务端设置缓存时间，客户端根据缓存时间来判断是否需要重新请求。
   [http缓存](file:///c:/Users/Administrator/Desktop/note/fullstack-note/05-工具层/网络/http缓存.md)
   ```go
   c.Header("Cache-Control", "max-age=300")
   ```

# 响应体：

这里就是我们通常说的 JSON 数据，或者 HTML 页面等。服务端根据客户端的请求，返回对应的响应体。

核心方法：status+data
c.JSON(status, data)：返回JSON格式的响应体
c.HTML(status, template, data)：返回HTML格式的响应体

# 文件下载

没有 Gin 的时候，这是读取的反操作：大致就是读取文件流，然后传入 body，再返回给客户端。

Gin 直接提供了 `c.File()` 方法，直接返回文件流，将 **读取 + 返回** 二合一了。

> 为什么二合一？因为不需要中间操作啊。你要这个文件，我直接给你。我需要操作文件的，在存储文件的时候已经做了。

上传文件也很方便，本质就是 **读取 + 存储**，两个函数完事：

```go
file := c.FormFile("file")
// 中间操作嘛，这也是为什么要分开的原因
c.SaveUploadedFile(file, "./uploads/"+file.Filename)
```

 **基本类似的操作，为什么有的地方是n合一，有的地方是分开的呢。**\
  **其实就是看你需不需要插入中间操作。**



