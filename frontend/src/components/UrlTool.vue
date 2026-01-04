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
  <div class="card">
    <div class="section-header">URL 编解码</div>
    
    <div class="dual-pane">
      <!-- 左侧输入 -->
      <div class="pane">
        <textarea class="code-area" v-model="urlInput" placeholder="输入原始内容..."></textarea>
        <div style="display:flex; gap:8px;">
          <button class="btn btn-primary" @click="urlEncode">编码</button>
          <button class="btn btn-outline" @click="urlDecode">解码</button>
        </div>
      </div>
      
      <!-- 右侧结果 -->
      <div class="pane">
        <textarea class="code-area" v-model="urlOutput" readonly style="background:#f9fafb;"></textarea>
        <button class="btn btn-outline" @click="copy">复制结果</button>
      </div>
    </div>
  </div>
</template>