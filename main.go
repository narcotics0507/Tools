package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- 1. 嵌入静态文件 ---
//
//go:embed dist/*
var staticFiles embed.FS

// --- 2. 数据库结构定义 ---
type UserEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ToolName  string    `json:"tool_name"`
	Action    string    `json:"action"`
	IsError   bool      `json:"is_error"`
	ErrorMsg  string    `json:"error_msg"`
	ClientIP  string    `json:"client_ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	// ===========================
	//    第一部分：数据库初始化
	// ===========================
	dbPath := "toolbox.db"
	if _, err := os.Stat("/data"); err == nil {
		dbPath = "/data/toolbox.db"
	}
	dbDir := filepath.Dir(dbPath)
	if dbDir != "." {
		_ = os.MkdirAll(dbDir, 0755)
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("无法连接数据库:", err)
	}
	db.AutoMigrate(&UserEvent{})

	// ===========================
	//    第二部分：Web 服务配置
	// ===========================
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 准备静态文件系统
	distFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		log.Fatal("无法加载静态文件:", err)
	}

	// 🔥【核心大招】：启动时直接把 index.html 读入内存
	// 这样我们就不用 FileFromFS 了，绝对不会产生 301 重定向
	indexData, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		log.Fatal("无法读取 dist/index.html，请检查 npm run build 是否成功:", err)
	}

	// ===========================
	//    核心路由逻辑
	// ===========================

	// 1. API 接口
	api := r.Group("/api")
	{
		api.POST("/event", func(c *gin.Context) {
			var evt UserEvent
			if err := c.ShouldBindJSON(&evt); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			evt.ClientIP = c.ClientIP()
			evt.UserAgent = c.Request.UserAgent()
			evt.CreatedAt = time.Now()
			go func(e UserEvent) { db.Create(&e) }(evt)
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		})

		api.GET("/stats", func(c *gin.Context) {
			var results []struct {
				ToolName string `json:"tool_name"`
				Count    int    `json:"count"`
			}
			db.Model(&UserEvent{}).
				Select("tool_name, count(*) as count").
				Group("tool_name").
				Order("count desc").
				Scan(&results)
			c.JSON(http.StatusOK, gin.H{
				"total_tools_usage": results,
				"server_time":       time.Now().In(time.FixedZone("CST", 8*3600)).Format(time.DateTime),
			})
		})

		api.GET("/stats/trend", func(c *gin.Context) {
			var results []struct {
				Date     string `json:"date"`
				ToolName string `json:"tool_name"`
				Count    int    `json:"count"`
			}
			// SQLite specific syntax for date grouping
			db.Model(&UserEvent{}).
				Select("strftime('%Y-%m-%d', created_at) as date, tool_name, count(*) as count").
				Where("created_at > ?", time.Now().AddDate(0, 0, -7)).
				Group("date, tool_name").
				Order("date asc").
				Scan(&results)

			c.JSON(http.StatusOK, gin.H{"trend": results})
		})
	}

	// 2. 静态资源 (/assets) - JS/CSS 文件
	assetsFS, _ := fs.Sub(distFS, "assets")
	r.StaticFS("/assets", http.FS(assetsFS))

	// 3. 首页和兜底路由 - 直接返回内存中的 HTML 数据
	renderIndex := func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", indexData)
	}

	// 明确注册根路径 (防止 fallback 逻辑复杂化)
	r.GET("/", renderIndex)

	// 处理 Favicon (避免 404)
	r.GET("/favicon.ico", func(c *gin.Context) {
		// 如果 favicon 存在则返回，不存在也不报错
		c.FileFromFS("favicon.ico", http.FS(distFS))
	})

	// 兜底路由 (SPA 刷新支持)
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
			return
		}
		renderIndex(c) // 任何不认识的页面都返回 index.html
	})

	log.Println("✅ 服务已启动 (内存直出模式)")
	log.Println("👉 请访问: http://localhost:8080")
	r.Run(":8080")
}
