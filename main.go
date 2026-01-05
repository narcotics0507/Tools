package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ==========================================
// 1. 静态资源嵌入配置
// ==========================================

// go:embed 指令会将 dist 目录下的所有文件打包进最终的二进制文件。
// 这样发布时只需要一个 .exe 或可执行文件，无需携带静态文件目录。
//
//go:embed dist/*
var staticFiles embed.FS

// ==========================================
// 2. IP 速率限制器 (Rate Limiter)
// ==========================================

// IPRateLimiter 用于针对每个 IP 地址进行限流，防止恶意刷接口。
type IPRateLimiter struct {
	ips map[string]*rate.Limiter // 存储 IP -> 限流器的映射
	mu  *sync.RWMutex            // 读写锁，保证并发安全
	r   rate.Limit               // 速率限制 (每秒允许多少个请求)
	b   int                      // 突发大小 (Burst)，允许瞬间突发的请求数量
}

// NewIPRateLimiter 创建一个新的 IP 限流器实例
// r: 每秒生成的令牌数 (request/s)
// b: 令牌桶容量 (burst)
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// GetLimiter 获取指定 IP 的限流器，如果不存在则自动创建一个新的
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// 全局限流器配置: 每秒允许 5 个请求，突发 10 个
var limiter = NewIPRateLimiter(5, 10)

// RateLimitMiddleware Gin 中间件：拦截请求并检查限流状态
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		// 如果该 IP 请求速率超标
		if !limiter.GetLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Too many requests (请求过于频繁，请稍后再试)",
			})
			c.Abort() // 终止后续处理
			return
		}
		c.Next() // 放行
	}
}

// ==========================================
// 3. 数据库模型定义 (GORM)
// ==========================================

// UserEvent 记录用户使用工具的行为事件
type UserEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"` // 自增主键
	ToolName  string    `json:"tool_name"`            // 工具名称 (如: json, hash, cron)
	Action    string    `json:"action"`               // 具体动作 (如: format, calculate)
	IsError   bool      `json:"is_error"`             // 是否发生错误
	ErrorMsg  string    `json:"error_msg"`            // 错误信息详情
	ClientIP  string    `json:"client_ip"`            // 用户 IP
	UserAgent string    `json:"user_agent"`           // 用户浏览器标识
	CreatedAt time.Time `json:"created_at"`           // 创建时间 (自动记录)
}

// ==========================================
// 4. 主程序入口
// ==========================================

func main() {
	// ---------------------------
	// 4.1 数据库初始化
	// ---------------------------

	// 确定数据库文件路径：
	// 优先检查 /data 目录是否存在 (通常用于 Docker 挂载卷)。
	// 如果 /data 存在，则将数据库存在 /data/toolbox.db 以便持久化。
	// 否则默认存在当前目录 (本地开发环境)。
	dbPath := "toolbox.db"
	if _, err := os.Stat("/data"); err == nil {
		dbPath = "/data/toolbox.db"
	}

	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if dbDir != "." {
		_ = os.MkdirAll(dbDir, 0755)
	}

	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("无法连接数据库:", err)
	}

	// 自动迁移模式：自动创建或更新 UserEvent 表结构
	db.AutoMigrate(&UserEvent{})

	// ---------------------------
	// 4.2 Web 服务配置 (Gin)
	// ---------------------------

	// 设置为发布模式 (减少控制台的调试日志输出)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 启用 Gzip 压缩 (DefaultCompression)
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// 准备静态文件系统：从嵌入的 staticFiles 中提取 "dist" 子目录
	distFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		log.Fatal("无法加载静态文件:", err)
	}

	// 🔥【核心优化】：将 index.html 预读到内存中
	// 这样做是为了实现 SPA (单页应用) 的兜底路由。
	// 当用户访问不存在的路径时，直接返回内存中的 index.html，避免 301 重定向，提升体验。
	indexData, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		log.Fatal("无法读取 dist/index.html，请检查 npm run build 是否成功:", err)
	}

	// ---------------------------
	// 4.3 路由定义
	// ---------------------------

	// (A) API 接口组 (前缀 /api)
	api := r.Group("/api")
	{
		// 1. 上报工具使用事件 (POST /api/event)
		// 启用了限流中间件 RateLimitMiddleware
		api.POST("/event", RateLimitMiddleware(), func(c *gin.Context) {
			var evt UserEvent
			// 绑定 JSON 请求体
			if err := c.ShouldBindJSON(&evt); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// 补充上下文信息
			evt.ClientIP = c.ClientIP()
			evt.UserAgent = c.Request.UserAgent()
			evt.CreatedAt = time.Now()

			// 异步写入数据库，不阻塞 HTTP 响应
			go func(e UserEvent) { db.Create(&e) }(evt)

			c.JSON(http.StatusOK, gin.H{"status": "success"})
		})

		// 2. 获取统计数据 (GET /api/stats)
		// 返回所有工具的使用次数排名
		api.GET("/stats", func(c *gin.Context) {
			var results []struct {
				ToolName string `json:"tool_name"`
				Count    int    `json:"count"`
			}
			// SQL: SELECT tool_name, count(*) FROM user_events GROUP BY tool_name ORDER BY count DESC
			db.Model(&UserEvent{}).
				Select("tool_name, count(*) as count").
				Group("tool_name").
				Order("count desc").
				Scan(&results)

			// 同时返回服务器时间 (CST/UTC+8)
			c.JSON(http.StatusOK, gin.H{
				"total_tools_usage": results,
				"server_time":       time.Now().In(time.FixedZone("CST", 8*3600)).Format(time.DateTime),
			})
		})

		// 3. 获取趋势数据 (GET /api/stats/trend)
		// 返回过去 7 天的每日工具使用量
		api.GET("/stats/trend", func(c *gin.Context) {
			var results []struct {
				Date     string `json:"date"`
				ToolName string `json:"tool_name"`
				Count    int    `json:"count"`
			}
			// SQLite 特有语法: strftime('%Y-%m-%d', created_at) 用于按天分组
			db.Model(&UserEvent{}).
				Select("strftime('%Y-%m-%d', created_at) as date, tool_name, count(*) as count").
				Where("created_at > ?", time.Now().AddDate(0, 0, -7)). // 仅查询最近 7 天
				Group("date, tool_name").
				Order("date asc").
				Scan(&results)

			c.JSON(http.StatusOK, gin.H{"trend": results})
		})
	}

	// (B) 静态资源托管 (/assets)
	// 这里的 assetsFS 对应 dist/assets 目录
	assetsFS, _ := fs.Sub(distFS, "assets")
	r.StaticFS("/assets", http.FS(assetsFS))

	// (C) 首页渲染函数
	// 直接向响应写入内存中的 index.html 内容
	renderIndex := func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", indexData)
	}

	// (D) 根路径路由
	r.GET("/", renderIndex)

	// (E) Favicon 支持 (避免浏览器控制台报 404)
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("favicon.ico", http.FS(distFS))
	})

	// (F) 兜底路由 (NoRoute) - 核心 SPA 支持
	// 对于任何未定义的路由 (且不是 /api 开头的)，都返回 index.html
	// 让前端 Vue Router 接管路由处理
	r.NoRoute(func(c *gin.Context) {
		// 如果是 API 请求但路径不对，返回 404 JSON
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
			return
		}
		// 否则返回前端页面
		renderIndex(c)
	})

	// ---------------------------
	// 4.4 启动服务
	// ---------------------------
	log.Println("✅ 服务已启动 (内存直出模式)")
	log.Println("👉 请访问: http://localhost:8080")
	r.Run(":8080")
}
