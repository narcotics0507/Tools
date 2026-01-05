<script setup>
import { ref, onUnmounted, computed } from 'vue'
import FlatPickr from 'vue-flatpickr-component';
import 'flatpickr/dist/flatpickr.css';
import 'flatpickr/dist/themes/airbnb.css'; // Modern Theme
import { Mandarin } from 'flatpickr/dist/l10n/zh.js';

const props = defineProps(['reportEvent'])

// ==========================================
// 1. 实时时钟
// ==========================================
const now = ref({ bj: '', ts: 0 })

const updateClock = () => {
    const d = new Date()
    const ts = Math.floor(d.getTime() / 1000)
    
    // 简单格式化
    const z = n => n < 10 ? '0' + n : n
    const year = d.getFullYear()
    const month = z(d.getMonth() + 1)
    const day = z(d.getDate())
    const hour = z(d.getHours())
    const min = z(d.getMinutes())
    const sec = z(d.getSeconds())
    
    now.value.bj = `${year}-${month}-${day} ${hour}:${min}:${sec}`
    now.value.ts = ts
}

const timer = setInterval(updateClock, 1000)
updateClock() // Init

onUnmounted(() => clearInterval(timer))

// ==========================================
// 2. 日期 -> 时间戳 (Flatpickr)
// ==========================================
const dateVal = ref('') // YYYY-MM-DD
const timeVal = ref('') // HH:mm:ss

const dateToTsResult = ref({ s: '', ms: '' })

// Flatpickr Config
const dateConfig = {
    locale: Mandarin,
    dateFormat: 'Y-m-d',
    allowInput: true,
    disableMobile: "true" // Force flatpickr on mobile
}

const timeConfig = {
    locale: Mandarin,
    enableTime: true,
    noCalendar: true,
    dateFormat: "H:i:S",
    time_24hr: true,
    enableSeconds: true,
    allowInput: true,
    disableMobile: "true"
}

// 初始化为当前时间
const initDateInput = () => {
    const d = new Date()
    const z = n => n < 10 ? '0' + n : n
    dateVal.value = `${d.getFullYear()}-${z(d.getMonth() + 1)}-${z(d.getDate())}`
    timeVal.value = `${z(d.getHours())}:${z(d.getMinutes())}:${z(d.getSeconds())}`
    
    // Defer calculation slightly to ensure models update
    setTimeout(convertToTs, 50)
}

const convertToTs = () => {
    props.reportEvent('timestamp', 'date_to_ts')
    if (!dateVal.value) return

    let t = timeVal.value || '00:00:00'
    // Ensure HH:mm:ss format
    if (t.length === 5) t += ':00'

    const d = new Date(`${dateVal.value} ${t}`)
    if (isNaN(d.getTime())) {
        dateToTsResult.value = { s: '无效日期', ms: '无效日期' }
        return
    }
    dateToTsResult.value = {
        s: Math.floor(d.getTime() / 1000),
        ms: d.getTime()
    }
}

// 挂载时初始化
initDateInput()

const copy = (t) => {
  if(!t || t==='无效日期') return
  navigator.clipboard.writeText(String(t))
  alert('已复制: ' + t)
}

// ==========================================
// 3. 批量处理
// ==========================================
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
</script>

<template>
  <div class="ts-container">
    
    <!-- 卡片 1: 实时时钟 -->
    <div class="card clock-card">
      <div class="row">
         <div class="clock-item" @click="copy(now.bj)">
            <label>当前时间</label>
            <div class="val">{{ now.bj }}</div>
         </div>
         <div class="divider"></div>
         <div class="clock-item" @click="copy(now.ts)">
            <label>当前时间戳 (s)</label>
            <div class="val highlight">{{ now.ts }}</div>
         </div>
      </div>
    </div>

    <!-- 卡片 2: 日期转时间戳 (Flatpickr) -->
    <div class="card">
        <div class="section-header">日期 → 时间戳</div>
        <div class="picker-row">
            
            <div class="flatpickr-wrapper">
                 <span class="icon">📅</span>
                 <FlatPickr 
                    v-model="dateVal"
                    :config="dateConfig"
                    class="fp-input"
                    placeholder="选择日期"
                 />
            </div>

            <div class="flatpickr-wrapper time-wrapper">
                 <span class="icon">🕒</span>
                 <FlatPickr 
                    v-model="timeVal"
                    :config="timeConfig"
                    class="fp-input"
                    placeholder="选择时间"
                 />
            </div>

            <button class="btn btn-primary" @click="convertToTs">转换</button>
            <button class="btn btn-outline" @click="initDateInput">重置</button>
        </div>
        
        <div class="row-result" v-if="dateToTsResult.s">
            <div class="res-item" @click="copy(dateToTsResult.s)">
                <span class="label">秒 (s):</span>
                <span class="val">{{ dateToTsResult.s }}</span>
            </div>
            <div class="res-item" @click="copy(dateToTsResult.ms)">
                <span class="label">毫秒 (ms):</span>
                <span class="val">{{ dateToTsResult.ms }}</span>
            </div>
        </div>
    </div>

    <!-- 卡片 3: 批量 时间戳 -> 日期 -->
    <div class="card full-height">
        <div class="section-header">批量转换 (时间戳 → 日期)</div>
        <div class="toolbar">
             <button class="btn btn-sm btn-primary" @click="batchConvertTs">开始转换</button>
             <button class="btn btn-sm btn-outline" @click="exportBatchCsv">导出 CSV</button>
             <button class="btn btn-sm btn-text" @click="batchTsInput=''; batchOutput=''">清空</button>
        </div>
        <div class="dual-area">
            <textarea v-model="batchTsInput" placeholder="在此粘贴多行时间戳..." class="area-left"></textarea>
            <textarea v-model="batchOutput" readonly placeholder="结果区域..." class="area-right"></textarea>
        </div>
    </div>

  </div>
</template>

<style scoped>
.ts-container {
    display: flex;
    flex-direction: column;
    gap: 16px;
    height: calc(100vh - 100px);
    font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue", Helvetica, Arial, sans-serif;
}

.card {
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05); 
}

.section-header {
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 16px;
    color: #111827;
}

/* Clock */
.clock-card .row { display: flex; align-items: center; justify-content: space-around; }
.clock-item { text-align: center; cursor: pointer; padding: 12px 32px; border-radius: 12px; transition: background 0.2s; }
.clock-item:hover { background: #f9fafb; }
.clock-item label { display: block; font-size: 13px; color: #6b7280; margin-bottom: 4px; }
.clock-item .val { font-family: "SF Mono", Menlo, monospace; font-size: 20px; font-weight: 600; color: #1f2937; }
.clock-item .highlight { color: #2563eb; }
.divider { width: 1px; height: 40px; background: #e5e7eb; }

/* Picker Row */
.picker-row { display: flex; gap: 12px; margin-bottom: 20px; align-items: center; }

/* Flatpickr Wrapper Customization */
.flatpickr-wrapper {
    position: relative;
    display: flex;
    align-items: center;
    background: white;
    border: 1px solid #d1d5db;
    border-radius: 8px;
    padding: 0 12px;
    height: 38px;
    width: 200px;
    transition: all 0.2s;
}
.flatpickr-wrapper:focus-within {
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}
.time-wrapper { width: 160px; }

.flatpickr-wrapper .icon {
    font-size: 16px;
    margin-right: 8px;
    filter: grayscale(1);
    opacity: 0.6;
}
:deep(.fp-input) {
    border: none;
    outline: none;
    font-size: 14px;
    color: #374151;
    font-family: inherit;
    width: 100%;
    background: transparent;
    font-weight: 500;
}

/* Result Row */
.row-result { 
    display: flex; gap: 32px; 
    background: #f8fafc; 
    padding: 16px; 
    border-radius: 10px; 
    border: 1px solid #f1f5f9;
}
.res-item { display: flex; align-items: baseline; gap: 8px; cursor: pointer; }
.res-item .label { font-size: 13px; color: #64748b; font-weight: 500; }
.res-item .val { font-family: "SF Mono", Menlo, monospace; font-size: 18px; font-weight: 600; color: #2563eb; }

/* Batch */
.full-height { flex: 1; display: flex; flex-direction: column; }
.toolbar { display: flex; gap: 12px; margin-bottom: 12px; }
.dual-area { flex: 1; display: flex; gap: 16px; min-height: 0; }
textarea { 
    flex: 1; padding: 16px; 
    border: 1px solid #e2e8f0; 
    border-radius: 8px; 
    resize: none; 
    font-family: "SF Mono", Menlo, monospace; 
    font-size: 12px; 
    outline: none;
    line-height: 1.5;
    transition: all 0.2s;
}
textarea:focus { border-color: #3b82f6; box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1); }
.area-right { background: #f8fafc; color: #475569; }

/* Buttons */
/* Buttons use global styles */
.btn-sm { height: 32px; padding: 0 12px; font-size: 12px; }
</style>