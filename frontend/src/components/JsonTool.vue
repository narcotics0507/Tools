<script setup>
import { ref } from 'vue'

const props = defineProps(['reportEvent'])

const jsonInput = ref('')
const jsonResult = ref('')
const jsonRaw = ref('')

const jsonFormat = () => {
  try {
    const obj = JSON.parse(jsonInput.value)
    jsonRaw.value = JSON.stringify(obj, null, 4)
    // 简单的语法高亮实现
    jsonResult.value = jsonRaw.value.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, (match) => {
      let cls = 'json-number'
      if (/^"/.test(match)) { cls = /:$/.test(match) ? 'json-key' : 'json-string' }
      else if (/true|false/.test(match)) { cls = 'json-boolean' }
      else if (/null/.test(match)) { cls = 'json-null' }
      return `<span class="${cls}">${match}</span>`
    })
    props.reportEvent('json', 'format_success')
  } catch (e) {
    jsonResult.value = '<span style="color:red">JSON 解析失败: ' + e.message + '</span>'
    props.reportEvent('json', 'format_error', e.message)
  }
}

const jsonMinify = () => {
  try {
    jsonRaw.value = JSON.stringify(JSON.parse(jsonInput.value))
    jsonResult.value = jsonRaw.value
    props.reportEvent('json', 'minify')
  } catch { }
}

const copy = (t) => {
  if(!t) return
  navigator.clipboard.writeText(String(t))
  alert('已复制')
}
</script>

<template>
  <div class="card">
    <div class="section-header">
      <span>JSON 格式化 & 高亮</span>
    </div>
    
    <div class="toolbar">
        <button class="btn btn-primary" @click="jsonFormat">格式化 (Format)</button>
        <button class="btn btn-outline" @click="jsonMinify">压缩 (Minify)</button>
        <button class="btn btn-outline" @click="copy(jsonRaw)">复制 (Copy)</button>
    </div>

    <div class="dual-pane">
      <div class="pane">
        <textarea class="code-area" v-model="jsonInput" placeholder="输入 JSON..."></textarea>
      </div>
      <div class="pane">
        <div class="code-area result-area" v-html="jsonResult || '转换结果...'"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.toolbar { display: flex; gap: 8px; margin-bottom: 16px; }
.dual-pane {
  display: flex; gap: 16px; height: 600px;
}
.pane {
  flex: 1; display: flex; flex-direction: column; overflow: hidden;
}
.code-area { 
    flex: 1; 
    resize: none; 
    font-family: var(--font-mono); 
    font-size: 13px; 
    padding: 12px;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    outline: none;
}
.result-area { 
    background: #f9fafb; 
    overflow: auto; 
    white-space: pre-wrap; 
    word-break: break-all;
}

/* Syntax Highlighting Colors */
/* Syntax Highlighting Colors - High Contrast for Sharpness */
:deep(.json-key) { color: #7e22ce; }           /* Darker Purple */
:deep(.json-string) { color: #15803d; }        /* Darker Green */
:deep(.json-number) { color: #b45309; }        /* Darker Amber */
:deep(.json-boolean) { color: #1d4ed8; font-weight: bold; }   /* Blue */
:deep(.json-null) { color: #b91c1c; font-weight: bold; }      /* Dark Red */
</style>