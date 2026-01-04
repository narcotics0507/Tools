<script setup>
import { ref, onUnmounted } from 'vue'

const props = defineProps(['reportEvent'])

// ==========================================
// 1. 实时时钟逻辑
// ==========================================
const now = ref({ bj: '', ts: 0 })

// 更新时钟：获取当前时间并手动校准为北京时间 (UTC+8)
const updateClock = () => {
  const d = new Date()
  // 计算 UTC 时间戳 (本地时间 + 时区偏移)
  const utc = d.getTime() + (d.getTimezoneOffset() * 60000)
  // 加上 8 小时得到北京时间
  const bjDate = new Date(utc + (3600000 * 8))

  // 补零函数: 9 -> 09
  const z = n => n < 10 ? '0' + n : n
  
  // 格式化输出: YYYY-MM-DD HH:mm:ss
  now.value.bj = `${bjDate.getFullYear()}-${z(bjDate.getMonth() + 1)}-${z(bjDate.getDate())} ${z(bjDate.getHours())}:${z(bjDate.getMinutes())}:${z(bjDate.getSeconds())}`
  // 秒级时间戳
  now.value.ts = Math.floor(d.getTime() / 1000)
}

// 启动每秒定时器
const timer = setInterval(updateClock, 1000)
updateClock() // 立即执行一次

// 组件销毁前清除定时器
onUnmounted(() => clearInterval(timer))

// ==========================================
// 2. 批量时间戳转换逻辑
// ==========================================
const batchTsInput = ref('')
const batchOutput = ref('')

// 批量转换: 当用户输入多行时间戳时，逐行转换为可读时间
const batchConvertTs = () => {
  props.reportEvent('timestamp', 'batch_convert')
  const lines = batchTsInput.value.split(/\r?\n/)
  
  batchOutput.value = lines.map((line, i) => {
    const raw = line.trim()
    const idx = (i + 1).toString().padStart(3, ' ')
    // 空行处理
    if (!raw) return `${idx} |                  | `
    
    const n = Number(raw)
    // 智能判断: 如果数值较小 (10位以内) 视为秒，否则视为毫秒
    // 注意: 这只是一个简单的启发式判断
    const d = new Date(raw.replace('-', '').length <= 10 ? n * 1000 : n)
    
    // 生成格式化的行: 序号 | 原始输入 | 转换结果
    return `${idx} | ${raw.padEnd(15)} | ${isNaN(d) ? '无效' : d.toLocaleString()}`
  }).join('\n')
}

// 导出 CSV 文件
const exportBatchCsv = () => {
  props.reportEvent('timestamp', 'export_csv')
  const lines = batchTsInput.value.split(/\n/)
  // 添加 UTF-8 BOM (\ufeff) 防止 Excel 打开乱码
  let csv = '\ufeff序号,时间戳,结果\r\n'
  
  lines.forEach((line, i) => {
    const raw = line.trim()
    if (!raw) { csv += `${i + 1},,""\r\n`; return }
    
    const n = Number(raw)
    const d = new Date(raw.replace('-', '').length <= 10 ? n * 1000 : n)
    csv += `${i + 1},${raw},"${isNaN(d) ? '无效' : d.toLocaleString()}"\r\n`
  })
  
  // 创建 Blob 并触发下载
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'Toolbox_Export.csv'
  a.click()
}

// 复制辅助函数
const copy = (t) => {
  if(!t) return
  navigator.clipboard.writeText(String(t))
  alert('已复制')
}
</script>

<template>
  <div class="card">
    <div class="section-header">实时钟 & 批量转换</div>
    
    <!-- 两个时钟卡片 (北京时间 / UTC时间戳) -->
    <div style="display:grid; grid-template-columns:1fr 1fr; gap:16px; margin-bottom:16px;">
      <div @click="copy(now.bj)" style="cursor:pointer; padding:12px; border:1px solid var(--border-color); border-radius:8px;">
        <div style="font-size:10px; color:#999;">北京时间 (已校准)</div>
        <div style="font-family:var(--font-mono); font-weight:700;">{{ now.bj }}</div>
      </div>
      <div @click="copy(now.ts)" style="cursor:pointer; padding:12px; border:1px solid var(--border-color); border-radius:8px; border-left:3px solid var(--accent);">
        <div style="font-size:10px; color:#999;">UNIX 时间戳</div>
        <div style="font-family:var(--font-mono); font-weight:700; color:var(--accent);">{{ now.ts }}</div>
      </div>
    </div>
    
    <!-- 批量处理区域 -->
    <textarea v-model="batchTsInput" placeholder="输入多行时间戳，支持空行..." style="height:120px; margin-bottom:12px;"></textarea>
    <div style="display:flex; gap:8px;">
      <button class="btn btn-primary" @click="batchConvertTs">批量转换</button>
      <button class="btn btn-outline" @click="exportBatchCsv">导出 CSV</button>
      <button class="btn btn-outline" @click="batchTsInput=''; batchOutput=''">清空</button>
    </div>
    
    <textarea v-model="batchOutput" readonly style="height:180px; background:#f9fafb; margin-top:12px;"></textarea>
  </div>
</template>