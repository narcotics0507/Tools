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

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement, PointElement, LineElement)

const stats = ref([])
const trendData = ref([])
const serverTime = ref('')
const loading = ref(true)

let timer = null
let timeOffset = 0

const updateTime = () => {
    const now = Date.now()
    const serverNow = new Date(now + timeOffset)
    // Format: YYYY-MM-DD HH:mm:ss
    const pad = (n) => n.toString().padStart(2, '0')
    const year = serverNow.getFullYear()
    const month = pad(serverNow.getMonth() + 1)
    const day = pad(serverNow.getDate())
    const hour = pad(serverNow.getHours())
    const minute = pad(serverNow.getMinutes())
    const second = pad(serverNow.getSeconds())
    serverTime.value = `${year}-${month}-${day} ${hour}:${minute}:${second}`
}

const fetchData = async () => {
    loading.value = true
    try {
        const resStats = await fetch('/api/stats').then(r => r.json())
        stats.value = resStats.total_tools_usage || []
        
        // Calculate offset
        // Parse server time (YYYY-MM-DD HH:mm:ss) to timestamp
        // Note: Replacing space with T might be needed for reliable parsing across browsers if ISO format isn't strictly followed, 
        // but given the Go format time.DateTime ("2006-01-02 15:04:05"), 
        // we can parse it carefully or trust new Date() in modern browsers.
        // Let's manually parse to be safe and avoid timezone issues (assuming server sends "local" time which we want to display as-is)
        // Actually, the server sends a string representing "CST". Ideally we display THAT time.
        // If we just do new Date(resStats.server_time), browser might interpret it as local time or UTC.
        // Let's stick to the server string as the anchor.
        
        const serverStr = resStats.server_time.replace(' ', 'T') // Make it ISO-like
        // We know server is +08:00. 
        // If we treat it as a string, we can just parse the components.
        
        // Better approach for visual "tick":
        // 1. Parse server string components.
        // 2. Construct a Date object representing that moment.
        // 3. Compare with Date.now() to get offset.
        
        const [dStr, tStr] = resStats.server_time.split(' ')
        const [y, m, d] = dStr.split('-').map(Number)
        const [hh, mm, ss] = tStr.split(':').map(Number)
        const serverDate = new Date(y, m - 1, d, hh, mm, ss)
        
        timeOffset = serverDate.getTime() - Date.now()
        
        updateTime() // Initial update
        if (timer) clearInterval(timer)
        timer = setInterval(updateTime, 1000)

        const resTrend = await fetch('/api/stats/trend').then(r => r.json())
        trendData.value = resTrend.trend || []
    } catch (e) {
        console.error("Failed to fetch stats:", e)
    } finally {
        loading.value = false
    }
}

onUnmounted(() => {
    if (timer) clearInterval(timer)
})

const chartColors = ['#41B883', '#E46651', '#00D8FF', '#DD1B16', '#F7DF1E', '#3178C6', '#9B59B6', '#34495E']

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

const lineChartData = computed(() => {
    if (!trendData.value.length) return { labels: [], datasets: [] }
    
    // Sort logic to ensure we cover a reasonable range or just existing data
    // For specific requirement "last 7 days", server returns data.
    // We just aggregate unique dates and tool names.
    const dates = [...new Set(trendData.value.map(t => t.date))].sort()
    const tools = [...new Set(trendData.value.map(t => t.tool_name))]
    
    const datasets = tools.map((tool, index) => {
        const color = chartColors[index % chartColors.length]
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
            tension: 0.4
        }
    })
    
    return { labels: dates, datasets }
})

const pieOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: {
            position: 'right',
            labels: { color: '#374151' } // Dark text for light theme
        }
    }
}

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
            radius: 5,
            hoverRadius: 7
        },
        line: {
            tension: 0.4
        }
    }
}

onMounted(fetchData)
</script>

<template>
  <div class="stats-container">
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
    
    <div v-if="loading" class="loading">
        <div class="spinner"></div> 加载数据中...
    </div>
    
    <div v-else class="content-wrapper">
        <div v-if="stats.length === 0" class="empty-state">
            暂无统计数据
        </div>
        
        <div v-else class="charts-grid">
            <div class="card">
                <h3>工具使用占比</h3>
                <div class="chart-wrapper">
                    <Pie :data="pieChartData" :options="pieOptions" />
                </div>
            </div>
            
            <div class="card full-width">
                <h3>7日使用趋势</h3>
                <div class="chart-wrapper">
                    <Line :data="lineChartData" :options="lineOptions" />
                </div>
            </div>
        </div>
        
        <!-- Detailed Table -->
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
    padding: 0; /* Let main-content handle padding */
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

.charts-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 20px;
}

@media (min-width: 1024px) {
    .charts-grid { grid-template-columns: 1fr 2fr; }
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
