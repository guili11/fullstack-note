package main

import (
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

// r gin中最关键的实例对象
func main() {
	// 1.创建路由管理器
	r := gin.Default()

	// 这个r
	// 1.储存路由对应信息
	// 2.当请求来的时候执行对应的handler
	// 3.监听端口
	// 4.启动服务
	// 5.挂载中间件（因为这个r控制一整个请求处理流程，所以中间件的挂载也需要他）

	// 其实呢，要理解为什么很简单，在代码世界中，要储存信息，或者运行函数，是需要一个machine的，这里就是r
	// 别的地方，例如axios，router，都是有machine的，  除非很简单的死操作，例如json翻译，直接调用函数就可以，不需要machine

	// 2.创建路由,以及对应的handler
	r.GET("/student", func(c *gin.Context) {
		// 内部处理逻辑
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
	r.POST("/student", func(c *gin.Context) {})
	r.PUT("/student", func(c *gin.Context) {})
	r.DELETE("/student", func(c *gin.Context) {})

	var testHandler func(c *gin.Context)
	testHandler = testDownload
	r.POST("/test", testHandler)

	// 3.启动服务
	r.Run(":8080")
}

func testForm(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	defaultLove := "unknown"
	love := c.DefaultPostForm("love", defaultLove)
	c.JSON(200, gin.H{
		"username": username,
		"password": password,
		"love":     love,
	})
}

func testFile(c *gin.Context) {
	// 一般来说，服务端获取到文件后，需要保存到服务器端
	// 所以主要操作就是读取 和 写入
	file, err := c.FormFile("file")
	// 自动识别格式：multipart/form-data
	if err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	// 保存文件
	dst := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":  "文件上传成功",
		"filename": file.Filename,
	})
}

func testMultiFile(c *gin.Context) {
	// 多文件上传
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	files := form.File["files"]
	filenames := []string{}

	for _, file := range files {
		dst := "./uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		filenames = append(filenames, file.Filename)
	}

	c.JSON(200, gin.H{
		"message":   "多文件上传成功",
		"filenames": filenames,
	})
}

func testDownload(c *gin.Context) {
	// 下载文件
	c.File("./uploads/OIP.webp")
	// 实际开发，要带一个身份信息，id，  路由参数  download/:id 
}

