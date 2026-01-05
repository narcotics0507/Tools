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
  LineElement,
  Filler
} from 'chart.js'
import { Pie, Line } from 'vue-chartjs'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement, PointElement, LineElement, Filler)

// 1. 响应式状态
const stats = ref([])         
const trendData = ref([])     
const timerTime = ref('')     
const timerDate = ref('')     
const loading = ref(true)     

let timer = null
let timeOffset = 0

// 2. 时间逻辑
const updateTime = () => {
    const now = Date.now()
    const serverNow = new Date(now + timeOffset)
    
    // Apple 风格：简洁
    const pad = (n) => n.toString().padStart(2, '0')
    const hour = pad(serverNow.getHours())
    const minute = pad(serverNow.getMinutes())
    const second = pad(serverNow.getSeconds())
    timerTime.value = `${hour}:${minute}:${second}`

    const year = serverNow.getFullYear()
    const month = pad(serverNow.getMonth() + 1)
    const day = pad(serverNow.getDate())
    const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
    const week = weekdays[serverNow.getDay()]
    timerDate.value = `${year}年${month}月${day}日 ${week}`
}

const fetchData = async () => {
    loading.value = true
    try {
        const resStats = await fetch('/api/stats').then(r => r.json())
        stats.value = resStats.total_tools_usage || []
        
        const [dStr, tStr] = resStats.server_time.split(' ')
        const [y, m, d] = dStr.split('-').map(Number)
        const [hh, mm, ss] = tStr.split(':').map(Number)
        const serverDate = new Date(y, m - 1, d, hh, mm, ss)
        timeOffset = serverDate.getTime() - Date.now()
        
        updateTime() 
        if (timer) clearInterval(timer)
        timer = setInterval(updateTime, 1000)

        const resTrend = await fetch('/api/stats/trend').then(r => r.json())
        trendData.value = resTrend.trend || []
    } catch (e) {
        console.error(e)
    } finally {
        loading.value = false
    }
}

onUnmounted(() => {
    if (timer) clearInterval(timer)
})

// 3. 计算属性 (中文)
const totalCalls = computed(() => stats.value.reduce((sum, item) => sum + item.count, 0))

// 缩写映射逻辑
const getToolAbbr = (name) => {
    const map = {
        'json': 'JSON',
        'sql': 'SQL',
        'number': 'BASE',
        'qrcode': 'QR',
        'color': 'COLOR',
        'text': 'TEXT',
        'timestamp': 'TIME'
    }
    // 默认转大写
    return map[name.toLowerCase()] || name.toUpperCase()
}

const topTool = computed(() => { 
    if (stats.value.length > 0) {
        return getToolAbbr(stats.value[0].tool_name)
    }
    return 'N/A'
})
const activeToolsCount = computed(() => stats.value.length)

// Apple Color Palette
const colors = [
    '#007AFF', // System Blue
    '#34C759', // System Green
    '#FF9500', // System Orange
    '#FF2D55', // System Pink
    '#AF52DE', // System Purple
    '#5856D6', // System Indigo
    '#5AC8FA', // System Teal
    '#8E8E93'  // System Gray
]

const pieChartData = computed(() => {
    if (!stats.value.length) return { labels: [], datasets: [] }
    // 使用缩写显示在图表中? 用户只说 "最常使用TOPtool中显示的工具都大写显示缩写"
    // 图表标签也可以跟随缩写，更简洁
    return {
        labels: stats.value.map(s => getToolAbbr(s.tool_name)),
        datasets: [{
            backgroundColor: colors,
            data: stats.value.map(s => s.count),
            borderWidth: 2,
            borderColor: '#ffffff',
            hoverOffset: 4
        }]
    }
})

const lineChartData = computed(() => {
    if (!trendData.value.length) return { labels: [], datasets: [] }
    const dates = [...new Set(trendData.value.map(t => t.date))].sort()
    const top3Tools = stats.value.slice(0, 3).map(s => s.tool_name)
    
    const datasets = top3Tools.map((tool, index) => {
        const color = colors[index]
        const data = dates.map(date => {
            const entry = trendData.value.find(t => t.date === date && t.tool_name === tool)
            return entry ? entry.count : 0
        })
        return {
            label: getToolAbbr(tool), // 缩写
            backgroundColor: (ctx) => {
                const canvas = ctx.chart.ctx;
                const gradient = canvas.createLinearGradient(0, 0, 0, 400);
                gradient.addColorStop(0, color + '66'); // 40% opacity
                gradient.addColorStop(1, color + '00'); // 0% opacity
                return gradient;
            },
            borderColor: color,
            data: data,
            fill: true, // 开启面积图填充
            tension: 0.4, // 平滑曲线
            pointRadius: 0, // 默认不显示点，更干净
            pointHoverRadius: 6,
            pointBackgroundColor: '#fff',
            pointBorderWidth: 2,
            borderWidth: 2
        }
    })
    
    return { labels: dates, datasets }
})

const commonOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: {
            position: 'bottom',
            labels: { 
                usePointStyle: true, 
                padding: 20,
                color: '#8E8E93',
                font: { size: 12, family: "-apple-system, BlinkMacSystemFont, sans-serif" }
            }
        },
        tooltip: {
            backgroundColor: 'rgba(255, 255, 255, 0.9)',
            titleColor: '#000',
            bodyColor: '#333',
            borderColor: 'rgba(0,0,0,0.05)',
            borderWidth: 1,
            padding: 10,
            cornerRadius: 10, // 更圆润
            titleFont: { size: 13, weight: 600, family: "-apple-system, BlinkMacSystemFont, sans-serif" },
            bodyFont: { size: 12, family: "-apple-system, BlinkMacSystemFont, sans-serif" },
            displayColors: true,
            boxWidth: 8,
            boxHeight: 8,
            usePointStyle: true,
            // 阴影
            callbacks: {
                labelTextColor: () => '#333'
            }
        }
    }
}

const lineOptions = {
    ...commonOptions,
    scales: {
        y: {
            beginAtZero: true,
            border: { display: false },
            grid: { color: '#F2F2F7', drawBorder: false }, // 极淡网格
            ticks: { color: '#C7C7CC', padding: 10, font: { family: "-apple-system, BlinkMacSystemFont, sans-serif" } }
        },
        x: {
            grid: { display: false, drawBorder: false },
            ticks: { color: '#C7C7CC', padding: 10, font: { family: "-apple-system, BlinkMacSystemFont, sans-serif" } }
        }
    },
    interaction: {
        mode: 'index',
        intersect: false,
    },
}

onMounted(fetchData)
</script>

<template>
  <div class="dashboard-container">
    <!-- Header -->
    <div class="header-section">
        <div class="title-group">
            <h1>System Status</h1>
        </div>
        <button @click="fetchData" class="refresh-btn">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M23 4v6h-6"/><path d="M1 20v-6h6"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>
            <span>刷新</span>
        </button>
    </div>

    <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
    </div>

    <div v-else class="grid-layout">
        <!-- Metrics - Apple Style: Glass / Clean White -->
        <div class="card metric-card">
            <div class="metric-label">服务器时间</div>
            <div class="metric-value time-val">{{ timerTime }}</div>
            <div class="metric-sub">{{ timerDate }}</div>
        </div>

        <div class="card metric-card">
            <div class="metric-label">Total Calls</div>
            <div class="metric-value">{{ totalCalls }}</div>
             <div class="metric-sub">累计请求量</div>
        </div>

        <div class="card metric-card">
            <div class="metric-label">Top Tool</div>
            <div class="metric-value text-ellipsis" style="color:#007AFF">{{ topTool }}</div>
             <div class="metric-sub">最常使用</div>
        </div>

        <div class="card metric-card">
            <div class="metric-label">Active Tools</div>
            <div class="metric-value">{{ activeToolsCount }}</div>
             <div class="metric-sub">在线服务</div>
        </div>

        <!-- Charts -->
        <div class="card chart-card line-chart-area">
             <div class="card-header">
                <h3>Last 7 Days Trend</h3>
             </div>
             <div class="chart-box">
                <Line :data="lineChartData" :options="lineOptions" />
             </div>
        </div>

        <div class="card chart-card pie-chart-area">
            <div class="card-header">
                <h3>工具分布</h3>
            </div>
            <div class="chart-box">
                <Pie :data="pieChartData" :options="commonOptions" />
            </div>
        </div>
        
        <!-- Table -->
        <div class="card table-card">
            <div class="card-header">
                <h3>详细报表</h3>
            </div>
            <div class="table-responsive">
                <table>
                    <thead>
                        <tr>
                            <th>TOOL NAME</th>
                            <th class="text-right">CALLS</th>
                            <th class="text-right">SHARE</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in stats" :key="item.tool_name">
                            <td class="font-medium">{{ getToolAbbr(item.tool_name) }}</td>
                            <td class="text-right num-font" style="font-weight:600">{{ item.count }}</td>
                            <td class="text-right text-gray num-font">
                                {{ ((item.count / totalCalls) * 100).toFixed(1) }}%
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
  </div>
</template>

<style scoped>
/* Apple System Fonts Stack - Clean & Native */
.dashboard-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 32px;
    font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue", Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    color: #1d1d1f; 
    /* Apple 浅灰背景 */
    background: #F5F5F7; 
    min-height: calc(100vh - 80px);
}

/* Header */
.header-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 32px;
}
.title-group h1 {
    font-size: 28px;
    font-weight: 700;
    margin: 0;
    color: #1d1d1f;
    letter-spacing: -0.01em;
}

.refresh-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    background: white;
    border: none;
    border-radius: 99px; /* Pill shape */
    padding: 8px 16px;
    font-size: 13px;
    font-weight: 500;
    color: #007AFF;
    cursor: pointer;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1); 
    transition: all 0.2s;
}
.refresh-btn:hover { 
    background: #f2f2f7;
    transform: scale(1.02);
}
.refresh-btn:active { transform: scale(0.98); }

/* Grid */
.grid-layout {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    grid-template-rows: auto auto auto; 
    gap: 20px;
}

/* Card - Apple Style: Soft, White, Minimal Shadow */
.card {
    background: white;
    border-radius: 18px; /* 更大的圆角 */
    padding: 24px;
    box-shadow: 0 4px 24px -1px rgba(0,0,0,0.04); 
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}
.card:hover {
    box-shadow: 0 8px 30px -5px rgba(0,0,0,0.08); 
    transform: translateY(-2px);
}

/* Metric Cards */
.metric-card {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    min-height: 140px;
}

.metric-label {
    font-size: 13px;
    color: #86868b;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}
.metric-value {
    font-size: 38px;
    font-weight: 700;
    color: #1d1d1f;
    line-height: 1.1;
    letter-spacing: -0.02em;
    margin: 12px 0;
}
.metric-sub {
    font-size: 13px;
    color: #86868b;
}
.time-val {
    font-family: "SF Pro Display", -apple-system, BlinkMacSystemFont, sans-serif;
    font-variant-numeric: tabular-nums; 
}

/* Chart Areas */
.line-chart-area {
    grid-column: span 3;
    height: 440px;
    display: flex;
    flex-direction: column;
}
.pie-chart-area {
    grid-column: span 1;
    height: 440px;
    display: flex;
    flex-direction: column;
}
.table-card {
    grid-column: span 4;
}

.card-header {
    margin-bottom: 24px;
    border-bottom: 1px solid #f5f5f7;
    padding-bottom: 12px;
}
.card-header h3 {
    margin: 0;
    font-size: 17px;
    font-weight: 600;
    color: #1d1d1f;
}

.chart-box {
    flex: 1;
    position: relative;
    min-height: 0; 
}

/* Table - Clean List */
.table-responsive { width: 100%; overflow-x: auto; }
table { width: 100%; border-collapse: collapse; }
th { 
    text-align: left; 
    font-size: 11px; 
    color: #86868b; 
    font-weight: 600; 
    padding: 12px 0; 
    border-bottom: 1px solid #e5e5ea;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}
td { 
    font-size: 14px; 
    padding: 16px 0; 
    color: #1d1d1f; 
    border-bottom: 1px solid #f5f5f7;
}
tr:last-child td { border-bottom: none; }
.text-right { text-align: right; }
.font-medium { font-weight: 500; }
.text-gray { color: #86868b; }
.num-font { font-variant-numeric: tabular-nums; }

/* Utilities */
.text-ellipsis { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

/* Loading */
.loading-state {
    height: 400px;
    display: flex;
    align-items: center;
    justify-content: center;
}
.spinner {
    width: 32px; height: 32px;
    border: 3px solid #f5f5f7;
    border-top-color: #007AFF; 
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* Responsive */
@media (max-width: 1024px) {
    .grid-layout { grid-template-columns: repeat(2, 1fr); }
    .line-chart-area, .pie-chart-area, .table-card { grid-column: span 2; }
}
@media (max-width: 640px) {
    .grid-layout { grid-template-columns: 1fr; }
    .line-chart-area, .pie-chart-area, .table-card, .metric-card { grid-column: span 1; }
}
</style>
