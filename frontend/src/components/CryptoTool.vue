<script setup>
import { ref } from 'vue'
const props = defineProps(['reportEvent'])

const crypto = ref({ base64In: '', base64Out: '' })

const base64Encode = () => {
  try {
    crypto.value.base64Out = btoa(unescape(encodeURIComponent(crypto.value.base64In)))
    props.reportEvent('crypto', 'b64_encode')
  } catch(e) { props.reportEvent('crypto', 'b64_encode_error', e.message) }
}
const base64Decode = () => {
  try {
    crypto.value.base64Out = decodeURIComponent(escape(atob(crypto.value.base64In)))
    props.reportEvent('crypto', 'b64_decode_success')
  } catch (e) {
    crypto.value.base64Out = 'Error'
    props.reportEvent('crypto', 'b64_decode_error', e.message)
  }
}
</script>

<template>
  <div class="card">
    <div class="section-header">加解密工具 (BASE64)</div>
    <textarea v-model="crypto.base64In" placeholder="输入内容..." style="height:100px; margin-bottom:10px;"></textarea>
    <div style="display:flex; gap:8px;">
      <button class="btn btn-primary" @click="base64Encode">编码</button>
      <button class="btn btn-outline" @click="base64Decode">解码</button>
    </div>
    <textarea v-model="crypto.base64Out" readonly style="height:100px; background:#f9fafb; margin-top:10px;"></textarea>
  </div>
</template>