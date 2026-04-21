package main

import (
	"log"
	"time"
	"github.com/gin-gonic/gin"
)
// 1. 全局日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("【请求】%s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

// 2. 请求耗时统计中间件
func CostMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		cost := time.Since(start)
		log.Printf("【耗时】%s 耗时: %v", c.Request.URL.Path, cost)
	}
}

// 3. 跨域中间件（生产必备）
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		c.Next()
	}
}

// 4. 全局异常捕获中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(500, gin.H{"msg": "服务器异常", "error": err})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// 5. 登录校验中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "admin123" {
			c.JSON(401, gin.H{"msg": "请先登录"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 6. 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("role")
		if role != "admin" {
			c.JSON(403, gin.H{"msg": "无管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 7. 接口限流中间件（模拟）
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 模拟：直接放行，仅做标记
		log.Println("【限流】接口限流校验通过")
		c.Next()
	}
}

// 8. 敏感接口校验中间件（单独给删除接口用）
func SensitiveMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("safe-key")
		if key != "123456" {
			c.JSON(403, gin.H{"msg": "敏感操作，校验失败"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func initi(){
	r:=gin.Default()
	
	// 全局中间件
	r.Use(LoggerMiddleware())
	r.Use(CostMiddleware())
	r.Use(CorsMiddleware())
	r.Use(RecoveryMiddleware())

	// 登录注册接口组
	auth:=r.Group("/auth")
	{
		auth.POST("/login", )
		auth.POST("/register", )
	}
	// 公开接口组：不用登录,但是限流
	open:=r.Group("/open")
	{
		open.Use(RateLimitMiddleware())
		open.GET("/notice", )
		open.GET("/captcha", )
	}

	// 用户业务组
	api:=r.Group("/api")
	{
		api.Use(AuthMiddleware()) // 🔥 全组统一加JWT鉴权
		api.GET("/info", )  
		api.GET("/order", )   
	}

	admin:=r.Group("/admin")
	{
		admin.Use(AuthMiddleware()) 
		admin.Use(AdminMiddleware()) 
		user:=admin.Group("/user")
		{
			user.GET("/list", )
      user.DELETE("/:id",SensitiveMiddleware())
		}
		stat:=admin.Group("/stat")
		{
			stat.GET("/total", )
		}
	}
	
}