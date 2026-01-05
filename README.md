<div align="center">
  <h1>🚀 上校工具箱 (Colonel's Toolbox)</h1>
  <p>
    <strong>现代 · 轻量 · 全能</strong>
  </p>
  <p>
    基于 <strong>Go (Gin)</strong> 和 <strong>Vue 3</strong> 构建的开发者效能神器
  </p>
  <p>
    <a href="https://github.com/narcotics0507/Tools/stargazers"><img src="https://img.shields.io/github/stars/narcotics0507/Tools?color=ffcb47&labelColor=black&style=for-the-badge" alt="Stars"></a>
    <a href="https://github.com/narcotics0507/Tools/network/members"><img src="https://img.shields.io/github/forks/narcotics0507/Tools?color=blue&labelColor=black&style=for-the-badge" alt="Forks"></a>
    <a href="https://github.com/narcotics0507/Tools/blob/master/LICENSE"><img src="https://img.shields.io/github/license/narcotics0507/Tools?color=green&labelColor=black&style=for-the-badge" alt="License"></a>
    <a href="https://hub.docker.com/r/isolatec/tools"><img src="https://img.shields.io/docker/v/isolatec/tools?color=blue&labelColor=black&style=for-the-badge&logo=docker&logoColor=white" alt="Docker Version"></a>
    <a href="https://hub.docker.com/r/isolatec/tools"><img src="https://img.shields.io/docker/pulls/isolatec/tools?color=blue&labelColor=black&style=for-the-badge" alt="Docker Pulls"></a>
  </p>
  <p>
    <a href="https://tools.sonic.nyc.mn/">🌐 在线体验</a> •
    <a href="#-快速部署">🚀 快速部署</a> •
    <a href="#-功能概览">✨ 功能概览</a>
  </p>
</div>

---

## 📖 项目简介

**上校工具箱** 是一款专为开发者打造的现代化工具集合。它摆脱了传统且臃肿的工具站模式，采用前后端分离架构，以极简的 Docker 镜像交付。
无论是本地开发调试，还是私有化部署在内网环境，它都能为您提供 **安全**、**快速**、**纯净** 的工具服务。

> 🛡️ **安全承诺**: 所有 API 均内置 IP 速率限制 (Rate Limiting)，有效防御恶意扫描与滥用。

---

## ✨ 功能概览

| 工具模块 | 功能描述 | 核心亮点 |
| :--- | :--- | :--- |
| **🎨 JSON 工具箱** | JSON 格式化、压缩、转义 | 支持左右分栏实时编辑与对比 |
| **⚖️ 差异比对 (Diff)** | 文本/代码差异检测 | Monaco Editor 驱动，支持多种语言高亮 |
| **🛢️ SQL 格式化** | SQL 美化、压缩 | 支持多种 SQL 方言，一键美化复杂查询 |
| **🔐 加解密 & Base64** | AES / DES / RSA / Base64 | 常用算法一站式解决 |
| **🧩 正则测试** | 正则表达式实时匹配 | 实时高亮显示，调试更直观 |
| **#️⃣ 哈希计算** | MD5, SHA1, SHA256, SHA512 | 快速生成文件或文本指纹 |
| **🔢 进制转换** | Binary / Octal / Decimal / Hex | 实时双向转换，支持大数 |
| **📱 二维码生成** | 文本 / URL 转 QR Code | 支持自定义纠错级别和尺寸 |
| **🎨 颜色转换** | HEX / RGB / HSL / CMYK | 开发者必备的色彩转换器 |
| **🔠 文本工具** | 大小写 / 驼峰 / 蛇形互转 | 变量命名风格快速切换 |
| **⏲️ Cron 表达式** | Cron 生成与反解析 | 可视化操作，不仅是工具也是教程 |
| **📅 时间戳转换** | 秒/毫秒 ↔ 日期时间 | 自动获取当前时间，双向转换 |
| **🔗 URL 编解码** | Encode / Decode | 快速处理 URL 参数 |
| **🆔 UUID 生成** | 批量生成 UUID (v1, v4) | 支持自定义分隔符和大小写 |
| **📊 调用统计** | 数据看板 | 可视化展示工具使用热度与趋势 |

---

## 🚀 快速部署

本项目已完全容器化，支持 **Docker** 一键启动。

### 🐳 方案一：Docker CLI (推荐)

```bash
docker run -d \
  --name my-tools \
  -p 8080:8080 \
  -v tools_data:/data \
  isolatec/tools:latest
```

*   **挂载说明**: `-v tools_data:/data` 用于持久化 SQLite 数据库，确保统计数据在容器重启后不丢失。

### 🐳 方案二：Docker Compose

```yaml
version: '3'
services:
  tools:
    image: isolatec/tools:latest
    container_name: my-tools
    ports:
      - "8080:8080"
    volumes:
      - tools_data:/data
    restart: unless-stopped

volumes:
  tools_data:
```

### 🛠️ 方案三：手动构建

如果您需要二次开发或手动构建镜像：

```bash
# 自动使用 goproxy.cn 加速构建
docker build -t tools:latest .
```

---

## 🛠️ 本地开发

### 环境要求
*   **Go** 1.24+
*   **Node.js** 20+

### 启动流程

1.  **构建前端**:
    ```bash
    cd frontend
    npm install
    npm run build
    ```

2.  **启动后端**:
    ```bash
    # 回到项目根目录
    go mod download
    go run main.go
    ```

3.  **访问**: 打开浏览器访问 `http://localhost:8080`

---

## 📝 技术栈

<p>
  <img src="https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Gin-008ECF?style=flat-square&logo=go&logoColor=white" alt="Gin" />
  <img src="https://img.shields.io/badge/Vue.js-4FC08D?style=flat-square&logo=vue.js&logoColor=white" alt="Vue" />
  <img src="https://img.shields.io/badge/Vite-646CFF?style=flat-square&logo=vite&logoColor=white" alt="Vite" />
  <img src="https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white" alt="Docker" />
  <img src="https://img.shields.io/badge/SQLite-003B57?style=flat-square&logo=sqlite&logoColor=white" alt="SQLite" />
</p>

---

## 📄 开源协议

本项目基于 **MIT 协议** 开源。
Copyright © 2025 narcotics0507.

<div align="center">
  <p>Made with ❤️ by narcotics0507</p>
</div>
