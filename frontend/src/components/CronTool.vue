<script setup>
import { ref, computed, watch } from 'vue'
// ✅ 使用新库：cron-schedule (原生支持 Vite，零报错)
import { parseCronExpression } from 'cron-schedule'

const props = defineProps(['reportEvent'])

const cron = ref({ mode: 'daily', week: '1', day: '1', hour: '00', minute: '00' })
const cronInputStr = ref('')
const nextRuns = ref([])
const parseError = ref('')

const weekDays = [{ val: '1', label: '一' }, { val: '2', label: '二' }, { val: '3', label: '三' }, { val: '4', label: '四' }, { val: '5', label: '五' }, { val: '6', label: '六' }, { val: '0', label: '日' }]

const cronResultString = computed(() => {
  let d = '*', w = '*'
  if (cron.value.mode === 'weekly') w = cron.value.week
  if (cron.value.mode === 'monthly') d = cron.value.day
  return `${parseInt(cron.value.minute)} ${parseInt(cron.value.hour)} ${d} * ${w}`
})

// 🔥 核心重写：使用 cron-schedule 预测时间
const predictNextRuns = (expression) => {
  try {
    parseError.value = ''
    // 1. 解析表达式
    const cronJob = parseCronExpression(expression)

    // 2. 计算未来 5 次运行时间
    const times = []
    let lastDate = new Date() // 从当前时间开始

    for (let i = 0; i < 5; i++) {
      // 获取下一次时间
      const next = cronJob.getNextDate(lastDate)
      times.push(next.toString())
      lastDate = next // 下一次基于这一次继续算
    }
    nextRuns.value = times
  } catch (err) {
    parseError.value = '格式错误: ' + err.message
    nextRuns.value = []
  }
}

watch(cronInputStr, (newVal) => { if(newVal) predictNextRuns(newVal) })

const parseCronString = () => {
  try {
    const str = cronInputStr.value.trim()
    if (!str) return
    predictNextRuns(str)

    // 如果解析失败，不回填 UI
    if (parseError.value) return

    const p = str.split(/\s+/)
    if (p.length < 5) throw new Error("Length < 5")

    const getFirst = (s) => s.includes(',') ? s.split(',')[0] : (s.includes('/') ? s.split('/')[0] : s)
    cron.value.minute = getFirst(p[0]).padStart(2, '0')
    cron.value.hour = getFirst(p[1]).padStart(2, '0')
    const d = p[2]
    const w = p[4]

    if (w !== '*' && w !== '?') { cron.value.mode = 'weekly'; cron.value.week = getFirst(w) }
    else if (d !== '*' && d !== '?') { cron.value.mode = 'monthly'; cron.value.day = getFirst(d) }
    else { cron.value.mode = 'daily' }
    props.reportEvent('cron', 'parse_success')
  } catch (e) {
    // 这里的错误通常是 split 分割导致的，非 cron 解析错误
    console.log(e)
  }
}

const copyResult = () => {
  navigator.clipboard.writeText(cronResultString.value)
  props.reportEvent('cron', 'copy_result')
  alert('已复制')
}
</script>

<template>
  <div class="card">
    <div class="section-header">Crond 表达式助手</div>

    <div class="panel">
      <label class="label">输入表达式 (反解析)</label>
      <div style="display: flex; gap: 8px;">
        <input v-model="cronInputStr" placeholder="例: 30 2 * * 1" class="mono-input">
        <button class="btn btn-blue" @click="parseCronString">解析并回填</button>
      </div>

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

    <div class="section-header" style="margin-top: 24px; border: none; margin-bottom: 12px;">可视化生成</div>
    <div class="tab-bar">
      <div class="tab-btn" :class="{active: cron.mode==='daily'}" @click="cron.mode='daily'">每天</div>
      <div class="tab-btn" :class="{active: cron.mode==='weekly'}" @click="cron.mode='weekly'">每周</div>
      <div class="tab-btn" :class="{active: cron.mode==='monthly'}" @click="cron.mode='monthly'">每月</div>
    </div>

    <div class="generator-box">
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

      <div class="time-picker">
        <div class="sub-label">执行时间</div>
        <div style="display:flex; align-items:center; gap:4px;">
          <select v-model="cron.hour" style="flex:1;"><option v-for="h in 24" :value="(h-1).toString().padStart(2,'0')">{{(h-1).toString().padStart(2,'0')}}</option></select>
          <span style="font-weight:bold; color:#d1d5db;">:</span>
          <select v-model="cron.minute" style="flex:1;"><option v-for="m in 60" :value="(m-1).toString().padStart(2,'0')">{{(m-1).toString().padStart(2,'0')}}</option></select>
        </div>
      </div>
    </div>

    <div @click="copyResult" class="result-bar">
      {{ cronResultString }}
    </div>
  </div>
</template>

<style scoped>
.panel { background: #f9fafb; padding: 16px; border-radius: 8px; border: 1px solid #e5e7eb; }
.label { display:block; font-size:11px; font-weight:600; color:#6b7280; margin-bottom:6px; text-transform:uppercase; }
.sub-label { font-size:12px; font-weight:600; color:#374151; margin-bottom:8px; }
.mono-input { font-family: var(--font-mono); font-size: 13px; }
.time-badge { background:#eef2ff; color:#4338ca; padding:2px 8px; border-radius:4px; font-size:11px; font-family:var(--font-mono); border: 1px solid #e0e7ff; }
.error-msg { color: #ef4444; font-size: 12px; }
.desc-text { color:#6b7280; font-size:13px; }

/* 布局修复 */
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

/* 强制不压缩右侧 */
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