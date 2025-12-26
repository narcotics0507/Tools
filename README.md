# 上校工具箱 (Colonel's Toolbox)

一个基于 Go (Gin) 和 Vue 3 的轻量级、模块化开发者工具箱。

## 🛠️ 功能列表
1.  **JSON 工具**: 格式化、压缩、语法高亮。
2.  **Base64 转换**: 文本编码与解码。
3.  **正则测试**: 实时正则表达式匹配与高亮。
4.  **Cron 表达式**: Crond 表达式生成器。
5.  **时间戳工具**: 时间戳转换。
6.  **URL 编解码**: URL Encode/Decode。
7.  **工具调用统计**: 查看工具使用热度与趋势。

## 在线体验地址：
https://tools.sonic.nyc.mn/


## 🚀 快速开始

### 本地运行
1.  **构建前端**:
    ```bash
    cd frontend
    npm install
    npm run build
    ```
2.  **运行后端**:
    ```bash
    go run main.go
    # 或构建二进制
    go build -o toolbox.exe
    ./toolbox.exe
    ```
3.  访问 `http://localhost:8080`

### 🐳 Docker 部署 (推荐)
这个项目已完全容器化。

#### 1. 构建镜像
```bash
docker build -t tools:v1.0.2 .
```

#### 2. 启动容器 (持久化数据)
为了防止重启容器后**统计数据丢失**，需要挂载 `/data` 目录：
####  可直接拉取最新镜像： isolatec/tools:latest
```bash
docker run -d \
  --name my-tools \
  -p 8080:8080 \
  -v tools_data:/data \
  tools:v1.0.2
```

*   `-d`: 后台运行
*   `-p 8080:8080`: 端口映射 (主机端口:容器端口)
*   `-v tools_data:/data`: **[重要]** 数据卷挂载。SQLite 数据库存储在容器内的 `/data` 目录，挂载卷可以确保数据持久化。
