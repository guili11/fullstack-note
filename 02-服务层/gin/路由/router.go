package main

import (
	"github.com/gin-gonic/gin"
)

func init() {
	r:=gin.Default()
	// 1.路径参数（GET请求中传递参数的方式）
	// 确定路径参数，业务接口，例如登录，获取用户信息，就是只需要一个id参数
	r.GET("/student/:id", func(c *gin.Context) {
		// 前端在路径中可以输入字符串，也就是student/xxx，这个xxx会被读取到，赋值给id
		id := c.Param("id")
		c.JSON(200, gin.H{
			"id": id,
		})
	})
	// 2.通配符参数
	// 看起来，这里的path是未知的不好处理啊， 实际开发中，这里path其实就是对应文件实际路径的参数 文件 / 静态资源（图片、文档、视频、前端页面）：必须用 * 多级通配符
	r.GET("/student/*path", func(c *gin.Context) {
		// 前端在路径中可以输入字符串，也就是student/xxx/yyy，这个xxx会被读取到，赋值给path
		path := c.Param("path")
		c.JSON(200, gin.H{
			"path": path,
		})
	})
	// 例如：
	// 	// 访问 /file/b/c/2.pdf
	// → filepath = /b/c/2.pdf
	// → 后端读取 ./static/b/c/2.pdf 返回

	//3.查询参数  
	// 查询参数，GET请求中传递参数的方式，例如：/student?id=123&name=张三
	r.GET("/student", func(c *gin.Context) {
		// 前端在路径中可以输入字符串，也就是student?id=123&name=张三，这个123会被读取到，赋值给id
		id := c.Query("id")
		name := c.Query("name")
		c.JSON(200, gin.H{
			"id": id,
			"name": name,
		})
	})
	// 例如：
	// 	// 访问 /student?id=123&name=张三
	// → id = 123
	// → name = 张三
}
