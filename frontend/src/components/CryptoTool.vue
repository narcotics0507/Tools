<script setup>
import { ref } from 'vue'

const props = defineProps(['reportEvent'])

const input = ref('')
const output = ref('')
const mode = ref('encode') // encode | decode

const process = () => {
    try {
        if (mode.value === 'encode') {
            output.value = btoa(unescape(encodeURIComponent(input.value)))
            props.reportEvent('crypto', 'b64_encode_success')
        } else {
            output.value = decodeURIComponent(escape(atob(input.value)))
            props.reportEvent('crypto', 'b64_decode_success')
        }
    } catch (e) {
        output.value = "错误: " + e.message
        props.reportEvent('crypto', 'b64_error', e.message)
    }
}

const copy = () => {
    navigator.clipboard.writeText(output.value)
}
</script>

<template>
  <div class="card">
    <div class="section-header">加解密 & Base64</div>
    
    <div class="tab-bar">
       <div class="tab-btn" :class="{ active: mode === 'encode' }" @click="mode = 'encode'">编码 (Encode)</div>
       <div class="tab-btn" :class="{ active: mode === 'decode' }" @click="mode = 'decode'">解码 (Decode)</div>
    </div>

    <div class="dual-pane">
      <div class="pane">
        <textarea class="code-area" v-model="input" placeholder="输入内容..." @input="process"></textarea>
      </div>
      <div class="pane">
        <textarea class="code-area" v-model="output" readonly placeholder="结果..."></textarea>
      </div>
    </div>
    
    <div style="margin-top:16px; text-align:right;">
         <button class="btn btn-primary" @click="copy">复制结果</button>
    </div>
  </div>
</template>

<style scoped>
.dual-pane { display: flex; gap: 16px; height: 300px; margin-top: 16px; }
.pane { flex: 1; display: flex; flex-direction: column; }
.code-area { flex: 1; resize: none; }
.controls { display: flex; gap: 16px; font-size: 14px; }
.controls label { cursor: pointer; display: flex; align-items: center; gap: 4px; }
</style>