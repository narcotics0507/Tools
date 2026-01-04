<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
  ArcElement,
  PointElement,
  LineElement
} from 'chart.js'
import { Pie, Line } from 'vue-chartjs'

// 注册 Chart.js 组件，必须在使用前注册
ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement, PointElement, LineElement)

// ==========================================
// 1. 响应式状态
// ==========================================
const stats = ref([])         // 工具总调用统计数据
const trendData = ref([])     // 7日趋势数据
const serverTime = ref('')    // 显示给用户的服务器时间字符串
const loading = ref(true)     // 加载状态

// 定时器引用，用于组件销毁时清除
let timer = null
// 本地时间与服务器时间的偏移量 (毫秒)
// 用于在本地准确显示服务器时间，不受本地时钟误差影响
let timeOffset = 0

// ==========================================
// 2. 时间同步逻辑
// ==========================================

// updateTime 每秒执行一次，利用 timeOffset 计算准确的服务器时间
const updateTime = () => {
    const now = Date.now()
    // 计算当前的服务器时间 = 本地时间 + 偏移量
    const serverNow = new Date(now + timeOffset)
    
    // 格式化为 YYYY-MM-DD HH:mm:ss
    const pad = (n) => n.toString().padStart(2, '0')
    const year = serverNow.getFullYear()
    const month = pad(serverNow.getMonth() + 1)
    const day = pad(serverNow.getDate())
    const hour = pad(serverNow.getHours())
    const minute = pad(serverNow.getMinutes())
    const second = pad(serverNow.getSeconds())
    
    serverTime.value = `${year}-${month}-${day} ${hour}:${minute}:${second}`
}

// fetchData 从后端获取统计数据和初始服务器时间
const fetchData = async () => {
    loading.value = true
    try {
        // 请求总调用统计
        const resStats = await fetch('/api/stats').then(r => r.json())
        stats.value = resStats.total_tools_usage || []
        
        // --- 核心时间校准逻辑 ---
        // 后端返回的是格式化的服务器时间字符串 (如 "2023-10-01 12:00:00")
        // 我们需要解析它，并计算它与本地时间的差值
        
        // 1. 解析服务器时间字符串
        const [dStr, tStr] = resStats.server_time.split(' ')
        const [y, m, d] = dStr.split('-').map(Number)
        const [hh, mm, ss] = tStr.split(':').map(Number)
        // 注意：Month 需要减 1 (0-11)
        const serverDate = new Date(y, m - 1, d, hh, mm, ss)
        
        // 2. 计算偏移量 (服务器时间 - 本地接收到响应的时间)
        timeOffset = serverDate.getTime() - Date.now()
        
        // 3. 立即更新一次 UI，并启动定时器
        updateTime() 
        if (timer) clearInterval(timer)
        timer = setInterval(updateTime, 1000)

        // 请求趋势数据
        const resTrend = await fetch('/api/stats/trend').then(r => r.json())
        trendData.value = resTrend.trend || []
    } catch (e) {
        console.error("无法获取统计数据:", e)
    } finally {
        loading.value = false
    }
}

// 组件卸载时清理定时器，防止内存泄漏
onUnmounted(() => {
    if (timer) clearInterval(timer)
})

// ==========================================
// 3. 图表配置
// ==========================================

// 预定义一组鲜艳的颜色，用于图表展示
const chartColors = ['#41B883', '#E46651', '#00D8FF', '#DD1B16', '#F7DF1E', '#3178C6', '#9B59B6', '#34495E']

// 计算属性：生成饼图所需的数据结构
const pieChartData = computed(() => {
    if (!stats.value.length) return { labels: [], datasets: [] }
    return {
        labels: stats.value.map(s => s.tool_name),
        datasets: [{
            backgroundColor: chartColors,
            data: stats.value.map(s => s.count),
            borderWidth: 0
        }]
    }
})

// 计算属性：生成折线图所需的数据结构
const lineChartData = computed(() => {
    if (!trendData.value.length) return { labels: [], datasets: [] }
    
    // 1. 提取所有唯一的日期和工具名称
    const dates = [...new Set(trendData.value.map(t => t.date))].sort()
    const tools = [...new Set(trendData.value.map(t => t.tool_name))]
    
    // 2. 构建每个工具的数据集
    const datasets = tools.map((tool, index) => {
        const color = chartColors[index % chartColors.length]
        // 映射每天的数据，如果某天没有数据则补 0
        const data = dates.map(date => {
            const entry = trendData.value.find(t => t.date === date && t.tool_name === tool)
            return entry ? entry.count : 0
        })
        return {
            label: tool,
            backgroundColor: color,
            borderColor: color,
            data: data,
            fill: false,
            tension: 0.4 // 曲线平滑度
        }
    })
    
    return { labels: dates, datasets }
})

// 饼图配置选项
const pieOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: {
            position: 'right',
            labels: { color: '#374151' } // 适配浅色主题文字颜色
        }
    }
}

// 折线图配置选项
const lineOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
        y: {
            beginAtZero: true,
            ticks: { color: '#6b7280', stepSize: 1 },
            grid: { color: '#e5e7eb' }
        },
        x: {
            ticks: { color: '#6b7280' },
            grid: { color: '#e5e7eb' }
        }
    },
    plugins: {
        legend: {
            labels: { color: '#374151' }
        }
    },
    elements: {
        point: {
            radius: 5,     // 数据点大小
            hoverRadius: 7 // 悬停时放大
        },
        line: {
            tension: 0.4   // 贝塞尔曲线
        }
    }
}

// 挂载时自动获取数据
onMounted(fetchData)
</script>

<template>
  <div class="stats-container">
    <!-- 顶部标题栏 -->
    <div class="header-row">
        <div>
            <h2>工具调用统计</h2>
            <div class="server-time">Server Time: {{ serverTime }}</div>
        </div>
        <button @click="fetchData" class="refresh-btn">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M23 4v6h-6"/><path d="M1 20v-6h6"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>
            刷新
        </button>
    </div>
    
    <!-- 加载中状态 -->
    <div v-if="loading" class="loading">
        <div class="spinner"></div> 加载数据中...
    </div>
    
    <!-- 数据内容区域 -->
    <div v-else class="content-wrapper">
        <div v-if="stats.length === 0" class="empty-state">
            暂无统计数据
        </div>
        
        <div v-else class="charts-grid">
            <!-- 饼图卡片 -->
            <div class="card">
                <h3>工具使用占比</h3>
                <div class="chart-wrapper">
                    <Pie :data="pieChartData" :options="pieOptions" />
                </div>
            </div>
            
            <!-- 折线图卡片 -->
            <div class="card full-width">
                <h3>7日使用趋势</h3>
                <div class="chart-wrapper">
                    <Line :data="lineChartData" :options="lineOptions" />
                </div>
            </div>
        </div>
        
        <!-- 底部详细数据表格 -->
        <div v-if="stats.length > 0" class="card mt-4">
             <h3>详细数据</h3>
             <table class="data-table">
                <thead>
                    <tr>
                        <th>工具名称</th>
                        <th>总调用次数</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="item in stats" :key="item.tool_name">
                        <td>{{ item.tool_name }}</td>
                        <td>{{ item.count }}</td>
                    </tr>
                </tbody>
             </table>
        </div>
    </div>
  </div>
</template>

<style scoped>
.stats-container {
    padding: 0;
    max-width: 1200px;
    margin: 0 auto;
    color: var(--text-primary);
    animation: fadeIn 0.4s ease-out;
}

.header-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
}

h2 { margin: 0; font-size: 1.25rem; font-weight: 600; color: #111827; }
h3 { margin: 0 0 16px 0; font-size: 1rem; font-weight: 600; color: #111827; }

.server-time { font-size: 0.85rem; color: var(--text-secondary); margin-top: 4px; }

.refresh-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    background: white;
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    padding: 8px 16px;
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: 0.85rem;
    transition: all 0.2s;
    font-weight: 500;
}
.refresh-btn:hover { background: #f9fafb; border-color: #d1d5db; }
.refresh-btn:active { transform: translateY(1px); }

/* Grid 布局：大屏显示两列，小屏单列 */
.charts-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 20px;
}

@media (min-width: 1024px) {
    .charts-grid { grid-template-columns: 1fr 2fr; } /* 饼图占 1/3，折线图占 2/3 */
}

.card {
    background: var(--bg-card);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    padding: 24px;
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.05);
}

.chart-wrapper {
    height: 300px;
    position: relative;
}

.mt-4 { margin-top: 20px; }

.data-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 8px;
}
.data-table th, .data-table td {
    text-align: left;
    padding: 12px;
    border-bottom: 1px solid var(--border-color);
    font-size: 0.9rem;
}
.data-table th { color: var(--text-secondary); font-weight: 500; background: #f9fafb; }
.data-table td { color: var(--text-primary); }
.data-table tr:last-child td { border-bottom: none; }

.loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 300px;
    color: var(--text-secondary);
}

.spinner {
    width: 32px;
    height: 32px;
    border: 3px solid var(--border-color);
    border-radius: 50%;
    border-top-color: var(--accent);
    animation: spin 1s linear infinite;
    margin-bottom: 16px;
}

@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>
