<script setup>
import { ref } from 'vue'
const props = defineProps(['reportEvent'])

const urlInput = ref('')
const urlOutput = ref('')

const urlEncode = () => {
  urlOutput.value = encodeURIComponent(urlInput.value)
  props.reportEvent('url', 'encode')
}
const urlDecode = () => {
  try {
    urlOutput.value = decodeURIComponent(urlInput.value)
    props.reportEvent('url', 'decode_success')
  } catch (e) {
    urlOutput.value = 'Error'
    props.reportEvent('url', 'decode_error', e.message)
  }
}
const copy = () => {
  navigator.clipboard.writeText(urlOutput.value)
  alert('已复制')
}
</script>

<template>
  <div class="card">
    <div class="section-header">URL 编解码</div>
    <div class="dual-pane">
      <div class="pane">
        <textarea class="code-area" v-model="urlInput" placeholder="输入原始内容..."></textarea>
        <div style="display:flex; gap:8px;">
          <button class="btn btn-primary" @click="urlEncode">编码</button>
          <button class="btn btn-outline" @click="urlDecode">解码</button>
        </div>
      </div>
      <div class="pane">
        <textarea class="code-area" v-model="urlOutput" readonly style="background:#f9fafb;"></textarea>
        <button class="btn btn-outline" @click="copy">复制结果</button>
      </div>
    </div>
  </div>
</template>