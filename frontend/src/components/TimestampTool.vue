<script setup>
import { ref, onUnmounted } from 'vue'

const props = defineProps(['reportEvent'])

// --- 实时时钟逻辑 ---
const now = ref({ bj: '', ts: 0 })
const updateClock = () => {
  const d = new Date()
  // 强制转换为北京时间 (UTC+8)
  const utc = d.getTime() + (d.getTimezoneOffset() * 60000)
  const bjDate = new Date(utc + (3600000 * 8))

  const z = n => n < 10 ? '0' + n : n
  now.value.bj = `${bjDate.getFullYear()}-${z(bjDate.getMonth() + 1)}-${z(bjDate.getDate())} ${z(bjDate.getHours())}:${z(bjDate.getMinutes())}:${z(bjDate.getSeconds())}`
  now.value.ts = Math.floor(d.getTime() / 1000)
}

// 启动定时器
const timer = setInterval(updateClock, 1000)
updateClock()

// 组件销毁时清除定时器，防止内存泄漏
onUnmounted(() => clearInterval(timer))

// --- 批量转换逻辑 ---
const batchTsInput = ref('')
const batchOutput = ref('')

const batchConvertTs = () => {
  props.reportEvent('timestamp', 'batch_convert')
  const lines = batchTsInput.value.split(/\r?\n/)
  batchOutput.value = lines.map((line, i) => {
    const raw = line.trim()
    const idx = (i + 1).toString().padStart(3, ' ')
    if (!raw) return `${idx} |                  | `
    const n = Number(raw)
    // 简单判断是秒还是毫秒
    const d = new Date(raw.replace('-', '').length <= 10 ? n * 1000 : n)
    return `${idx} | ${raw.padEnd(15)} | ${isNaN(d) ? '无效' : d.toLocaleString()}`
  }).join('\n')
}

const exportBatchCsv = () => {
  props.reportEvent('timestamp', 'export_csv')
  const lines = batchTsInput.value.split(/\n/)
  let csv = '\ufeff序号,时间戳,结果\r\n'
  lines.forEach((line, i) => {
    const raw = line.trim()
    if (!raw) { csv += `${i + 1},,""\r\n`; return }
    const n = Number(raw)
    const d = new Date(raw.replace('-', '').length <= 10 ? n * 1000 : n)
    csv += `${i + 1},${raw},"${isNaN(d) ? '无效' : d.toLocaleString()}"\r\n`
  })
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
    <textarea v-model="batchTsInput" placeholder="输入多行时间戳，支持空行..." style="height:120px; margin-bottom:12px;"></textarea>
    <div style="display:flex; gap:8px;">
      <button class="btn btn-primary" @click="batchConvertTs">批量转换</button>
      <button class="btn btn-outline" @click="exportBatchCsv">导出 CSV</button>
      <button class="btn btn-outline" @click="batchTsInput=''; batchOutput=''">清空</button>
    </div>
    <textarea v-model="batchOutput" readonly style="height:180px; background:#f9fafb; margin-top:12px;"></textarea>
  </div>
</template>