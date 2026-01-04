<script setup>
import { ref, watch } from 'vue'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'

const props = defineProps(['reportEvent'])

const jsonInput = ref('')
const jsonData = ref(null)
const parseError = ref('')
const showLine = ref(true)
const showLineNumber = ref(true)
const showIcon = ref(true)
const deep = ref(3)

const handleJsonChange = () => {
  try {
    if (!jsonInput.value.trim()) {
      jsonData.value = null
      parseError.value = ''
      return
    }
    const obj = JSON.parse(jsonInput.value)
    jsonData.value = obj
    parseError.value = ''
    props.reportEvent('json', 'format_success')
  } catch (e) {
    // Keep showing old data if parse fails? Or clear it? 
    // Usually json.cn clears or shows error. Let's show error only.
    jsonData.value = null
    parseError.value = 'JSON 解析失败: ' + e.message
    props.reportEvent('json', 'format_error', e.message)
  }
}

const copyResult = () => {
  if (jsonData.value) {
    navigator.clipboard.writeText(JSON.stringify(jsonData.value, null, 4))
    alert('已复制格式化后的 JSON')
  }
}

const compress = () => {
    try {
        const obj = JSON.parse(jsonInput.value)
        jsonInput.value = JSON.stringify(obj)
        handleJsonChange()
        props.reportEvent('json', 'minify')
    } catch (e) {
        alert('JSON 无效，无法压缩')
    }
}
</script>

<template>
  <div class="card full-height">
    <div class="section-header">
        <span>JSON 编辑器 (类似 json.cn)</span>
    </div>

    <!-- Enhanced Controls -->
    <div class="controls-bar">
       <div class="toggles-group">
          <label class="toggle-label">
            <div class="toggle-switch">
                <input type="checkbox" v-model="showLine">
                <span class="slider"></span>
            </div>
            <span>连线</span>
          </label>
          
          <label class="toggle-label">
            <div class="toggle-switch">
                <input type="checkbox" v-model="showLineNumber">
                <span class="slider"></span>
            </div>
            <span>行号</span>
          </label>

          <label class="toggle-label">
            <div class="toggle-switch">
                <input type="checkbox" v-model="showIcon">
                <span class="slider"></span>
            </div>
            <span>图标</span>
          </label>
       </div>

       <div class="center-group">
          <span class="label-text">默认展开:</span>
          <select v-model="deep" class="select-box">
              <option :value="1">1层</option>
              <option :value="2">2层</option>
              <option :value="3">3层</option>
              <option :value="4">4层</option>
              <option :value="10">全部</option>
          </select>
       </div>

       <div class="actions-group">
          <button class="btn btn-sm btn-outline" @click="compress">压缩</button>
          <button class="btn btn-sm btn-primary" @click="handleJsonChange">格式化</button>
          <button class="btn btn-sm btn-outline" @click="copyResult">复制结果</button>
       </div>
    </div>

    <div class="dual-pane">
      <div class="pane left-pane">
        <textarea 
            class="code-area" 
            v-model="jsonInput" 
            @input="handleJsonChange"
            placeholder="在此粘贴 JSON 字符串..."
        ></textarea>
        <div v-if="parseError" class="error-msg">{{ parseError }}</div>
      </div>
      
      <div class="pane right-pane">
        <div class="json-viewer-container" v-if="jsonData">
            <vue-json-pretty
                :data="jsonData"
                :showLine="showLine"
                :showLineNumber="showLineNumber"
                :showIcon="showIcon"
                :deep="deep"
                :editable="false" 
            />
        </div>
        <div v-else class="placeholder-text">
            等待输入有效的 JSON...
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.full-height {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 100px);
}

.controls-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    margin-bottom: 12px;
    border-bottom: 1px solid var(--border-color);
}

.toggles-group, .actions-group, .center-group {
    display: flex;
    align-items: center;
    gap: 16px;
}

.label-text {
    font-size: 13px;
    color: #4b5563;
    white-space: nowrap;
}

.select-box {
    padding: 4px 8px;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    font-size: 13px;
    outline: none;
    cursor: pointer;
}

/* Toggle Switch CSS */
.toggle-label {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    font-size: 13px;
    color: #374151;
    user-select: none;
}

.toggle-switch {
    position: relative;
    width: 36px;
    height: 20px;
}

.toggle-switch input {
    opacity: 0;
    width: 0;
    height: 0;
}

.slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: #e5e7eb;
    transition: .4s;
    border-radius: 20px;
}

.slider:before {
    position: absolute;
    content: "";
    height: 16px;
    width: 16px;
    left: 2px;
    bottom: 2px;
    background-color: white;
    transition: .4s;
    border-radius: 50%;
    box-shadow: 0 1px 2px rgba(0,0,0,0.1);
}

input:checked + .slider {
    background-color: var(--primary);
}

input:checked + .slider:before {
    transform: translateX(16px);
}

/* Common Layout */
.dual-pane {
  flex: 1;
  display: flex;
  gap: 16px;
  overflow: hidden; 
}

.pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  background: white;
}

.left-pane {
    position: relative;
}

.code-area { 
    flex: 1; 
    resize: none; 
    border: none;
    outline: none;
    padding: 12px;
    font-family: var(--font-mono); 
    font-size: 13px; 
    line-height: 1.5;
}

.error-msg {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    background: #fef2f2;
    color: #b91c1c;
    padding: 8px 12px;
    font-size: 12px;
    border-top: 1px solid #fecaca;
}

.right-pane {
    background: #f9fafb;
    overflow: auto;
}

.json-viewer-container {
    padding: 12px;
}

.placeholder-text {
    padding: 20px;
    color: #9ca3af;
    text-align: center;
    margin-top: 40px;
}

.btn-sm {
    padding: 4px 12px;
    font-size: 12px;
}
</style>