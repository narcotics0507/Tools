<script setup>
import { ref } from 'vue'
const props = defineProps(['reportEvent'])

// ==========================================
// 状态定义
// ==========================================
const urlInput = ref('')
const urlOutput = ref('')

// ==========================================
// 逻辑处理
// ==========================================

// URL 编码 (保留特殊字符的含义，转义非安全字符)
const urlEncode = () => {
  // 使用 encodeURIComponent 对 URI 组件进行编码
  urlOutput.value = encodeURIComponent(urlInput.value)
  props.reportEvent('url', 'encode')
}

// URL 解码
const urlDecode = () => {
  try {
    urlOutput.value = decodeURIComponent(urlInput.value)
    props.reportEvent('url', 'decode_success')
  } catch (e) {
    urlOutput.value = 'Error'
    props.reportEvent('url', 'decode_error', e.message)
  }
}

// 复制结果
const copy = () => {
  navigator.clipboard.writeText(urlOutput.value)
  alert('已复制')
}
</script>

<template>
  <div class="card full-height">
    <div class="section-header">URL 编解码</div>
    
    <div class="dual-pane">
      <!-- 左侧输入 -->
      <div class="pane input-pane">
        <div class="pane-content">
            <textarea class="code-area" v-model="urlInput" placeholder="输入原始内容..."></textarea>
        </div>
        <div class="pane-footer">
          <button class="btn btn-primary" @click="urlEncode">编码 (Encode)</button>
          <button class="btn btn-outline" @click="urlDecode">解码 (Decode)</button>
        </div>
      </div>
      
      <!-- 右侧结果 -->
      <div class="pane output-pane">
        <div class="pane-content">
            <textarea class="code-area" v-model="urlOutput" readonly placeholder="转换结果..."></textarea>
        </div>
        <div class="pane-footer right-footer">
            <button class="btn btn-outline" @click="copy">复制结果</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 局部微调: 让 Pane 的内部布局更精致 */
.pane {
    background: white; /* 统一白底 */
    display: flex;
    flex-direction: column;
}

.pane-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 0; /* Remove padding to let textarea fill */
}

.code-area {
    width: 100%;
    height: 100%;
    border: none; /* Remove border from textarea, let pane handle it */
    box-shadow: none !important;
    padding: 16px;
    resize: none;
    background: transparent;
}

.pane-footer {
    padding: 12px 16px;
    border-top: 1px solid #f2f2f7;
    background: #fafafa; /* Very light gray footer */
    display: flex;
    gap: 12px;
    align-items: center;
}

.right-footer {
    justify-content: flex-end;
}
</style>