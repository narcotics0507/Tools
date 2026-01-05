<script setup>
import { ref, computed, watch } from 'vue'
// ✅ 引用 cron-schedule 库进行 Cron 表达式解析和预测
import { parseCronExpression } from 'cron-schedule'

const props = defineProps(['reportEvent'])

// ==========================================
// 1. 状态定义
// ==========================================
// cron 对象存储生成器的当前状态: 模式(每天/每周/每月), 周几, 日期, 小时, 分钟
const cron = ref({ mode: 'daily', week: '1', day: '1', hour: '00', minute: '00' })

// 反解析输入框的值
const cronInputStr = ref('')
// 预测的未来执行时间列表
const nextRuns = ref([])
// 解析错误信息
const parseError = ref('')

// 星期几的映射 (0=周日, 1=周一)
const weekDays = [{ val: '1', label: '一' }, { val: '2', label: '二' }, { val: '3', label: '三' }, { val: '4', label: '四' }, { val: '5', label: '五' }, { val: '6', label: '六' }, { val: '0', label: '日' }]

// ==========================================
// 2. 计算与逻辑
// ==========================================

// 根据当前 cron 对象的状态生成的 Crond 表达式字符串
const cronResultString = computed(() => {
  let d = '*', w = '*'
  if (cron.value.mode === 'weekly') w = cron.value.week
  if (cron.value.mode === 'monthly') d = cron.value.day
  // 格式: 分 时 日 月 周
  return `${parseInt(cron.value.minute)} ${parseInt(cron.value.hour)} ${d} * ${w}`
})

// 🔥 核心功能：预测未来 5 次运行时间
const predictNextRuns = (expression) => {
  try {
    parseError.value = ''
    // 1. 使用库解析表达式
    const cronJob = parseCronExpression(expression)

    // 2. 循环计算未来 5 次
    const times = []
    let lastDate = new Date() // 起始时间为当前

    for (let i = 0; i < 5; i++) {
      const next = cronJob.getNextDate(lastDate)
      times.push(next.toString())
      lastDate = next // 下一次基于这一次继续往后算
    }
    nextRuns.value = times
  } catch (err) {
    parseError.value = '格式错误: ' + err.message
    nextRuns.value = []
  }
}

// 监听输入框变化，自动触发预测
watch(cronInputStr, (newVal) => { if(newVal) predictNextRuns(newVal) })

// 按钮点击事件：解析用户输入的 Cron 字符串并回填到 UI
const parseCronString = () => {
  try {
    const str = cronInputStr.value.trim()
    if (!str) return
    predictNextRuns(str)

    // 如果表达式本身非法，就不尝试回填 UI 了
    if (parseError.value) return

    // 简单的空格分割解析 (注意：这只能处理简单的 5 段式标准 Cron)
    const p = str.split(/\s+/)
    if (p.length < 5) throw new Error("Length < 5")

    const getFirst = (s) => s.includes(',') ? s.split(',')[0] : (s.includes('/') ? s.split('/')[0] : s)
    cron.value.minute = getFirst(p[0]).padStart(2, '0')
    cron.value.hour = getFirst(p[1]).padStart(2, '0')
    const d = p[2]
    const w = p[4]

    // 根据字段判断当前是哪种模式
    if (w !== '*' && w !== '?') { cron.value.mode = 'weekly'; cron.value.week = getFirst(w) }
    else if (d !== '*' && d !== '?') { cron.value.mode = 'monthly'; cron.value.day = getFirst(d) }
    else { cron.value.mode = 'daily' }
    
    props.reportEvent('cron', 'parse_success')
  } catch (e) {
    // 解析失败通常是因为表达式太复杂 (如区间、列表等)，UI无法完全还原，仅报错日志
    console.log(e)
  }
}

// 复制结果到剪贴板
const copyResult = () => {
  navigator.clipboard.writeText(cronResultString.value)
  props.reportEvent('cron', 'copy_result')
  alert('已复制')
}
</script>

<template>
  <div class="card">
    <div class="section-header">Crond 表达式助手</div>

    <!-- 顶部：反解析面板 -->
    <div class="panel">
      <label class="label">输入表达式 (反解析)</label>
      <div style="display: flex; gap: 8px;">
        <input v-model="cronInputStr" placeholder="例: 30 2 * * 1" class="code-area">
        <button class="btn btn-primary" @click="parseCronString">解析并回填</button>
      </div>

      <!-- 预测结果显示区域 -->
      <div style="margin-top: 12px; min-height: 20px;">
        <div v-if="parseError" class="error-msg">{{ parseError }}</div>
        <div v-else-if="nextRuns.length > 0">
          <div style="font-size:11px; color:#9ca3af; margin-bottom:4px;">未来 5 次执行:</div>
          <div style="display:flex; flex-wrap:wrap; gap:6px;">
                    <span v-for="time in nextRuns" :key="time" class="time-badge">
                        {{ new Date(time).toLocaleString() }}
                    </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 中部：可视化生成器 -->
    <div class="section-header" style="margin-top: 24px; border: none; margin-bottom: 12px;">可视化生成</div>
    
    <!-- 模式切换 Tab -->
    <div class="tab-bar">
      <div class="tab-btn" :class="{active: cron.mode==='daily'}" @click="cron.mode='daily'">每天</div>
      <div class="tab-btn" :class="{active: cron.mode==='weekly'}" @click="cron.mode='weekly'">每周</div>
      <div class="tab-btn" :class="{active: cron.mode==='monthly'}" @click="cron.mode='monthly'">每月</div>
    </div>

    <div class="generator-box">
      <!-- 左侧：日期选择控件 -->
      <div class="left-controls">
        <div v-if="cron.mode==='daily'" class="desc-text">任务将在每天指定时间执行。</div>
        
        <div v-if="cron.mode==='weekly'">
          <div class="sub-label">选择星期</div>
          <div style="display:flex; gap:6px;">
            <button v-for="d in weekDays" :key="d.val" class="btn"
                    :class="cron.week === d.val ? 'btn-primary':'btn-outline'"
                    @click="cron.week = d.val" style="flex:1;">{{ d.label }}</button>
          </div>
        </div>
        
        <div v-if="cron.mode==='monthly'">
          <div class="sub-label">选择日期</div>
          <div class="month-grid">
            <div v-for="n in 31" :key="n" class="grid-item"
                 :class="{active: parseInt(cron.day) === n}"
                 @click="cron.day = n.toString()">{{ n }}</div>
          </div>
        </div>
      </div>

      <!-- 右侧：时间选择控件 -->
      <div class="time-picker">
        <div class="sub-label">执行时间</div>
        <div style="display:flex; align-items:center; gap:4px;">
          <select v-model="cron.hour" style="flex:1;"><option v-for="h in 24" :value="(h-1).toString().padStart(2,'0')">{{(h-1).toString().padStart(2,'0')}}</option></select>
          <span style="font-weight:bold; color:#d1d5db;">:</span>
          <select v-model="cron.minute" style="flex:1;"><option v-for="m in 60" :value="(m-1).toString().padStart(2,'0')">{{(m-1).toString().padStart(2,'0')}}</option></select>
        </div>
      </div>
    </div>

    <!-- 底部：结果展示条 -->
    <div @click="copyResult" class="result-bar">
      {{ cronResultString }}
    </div>
  </div>
</template>

<style scoped>
.panel { background: #f9fafb; padding: 16px; border-radius: var(--radius-sm); border: 1px solid rgba(0,0,0,0.05); }
.label { display:block; font-size:11px; font-weight:600; color:#6b7280; margin-bottom:6px; text-transform:uppercase; }
.sub-label { font-size:12px; font-weight:600; color:#374151; margin-bottom:8px; }
.time-badge { 
    background:#eff6ff; 
    color:#2563eb; 
    padding:6px 12px; /* Larger padding */
    border-radius:99px; /* Pill shape */
    font-size:13px; /* Larger text */
    font-family:var(--font-mono); 
    border: 1px solid #bfdbfe; /* Visible border */
    font-weight: 600;
}
.error-msg { color: #ef4444; font-size: 12px; }
.desc-text { color:#6b7280; font-size:13px; }

/* 布局修复: 使左右控件在小屏下堆叠，大屏下并排 */
.generator-box {
  display: flex;
  gap: 24px;
  padding: 20px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  margin-bottom: 16px;
  align-items: flex-start;
  flex-wrap: wrap;
}

.left-controls {
  flex: 1;
  min-width: 280px;
}

/* 强制不压缩右侧时间选择器 */
.time-picker {
  flex: 0 0 240px;
  border-left: 1px dashed #e5e7eb;
  padding-left: 24px;
}

@media (max-width: 640px) {
  .generator-box { flex-direction: column; gap: 16px; }
  .time-picker {
    border-left: none;
    padding-left: 0;
    border-top: 1px dashed #e5e7eb;
    padding-top: 16px;
    width: 100%;
    flex: auto;
  }
}

.result-bar { background:#1f2937; color:#4ade80; padding:16px; text-align:center; font-family:var(--font-mono); font-size:18px; font-weight:600; border-radius:8px; cursor:pointer; transition:transform 0.1s; }
.result-bar:active { transform: scale(0.99); }
</style>