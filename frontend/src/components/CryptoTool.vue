<script setup>
import { ref } from 'vue'
import { debounce } from '../utils/debounce'

const props = defineProps(['reportEvent'])

// 防抖上报函数
const debouncedReport = debounce((type, action, msg) => {
    props.reportEvent(type, action, msg)
}, 2000)

// ==========================================
// 1. 状态定义
// ==========================================
const input = ref('')
const output = ref('')
const mode = ref('encode') // 当前模式: 'encode' (编码) | 'decode' (解码)

// ==========================================
// 2. 逻辑处理
// ==========================================

// 执行转换处理
const process = () => {
    try {
        if (mode.value === 'encode') {
            // Base64 编码
            // 注意: btoa 不能直接处理 UTF-8 字符 (如中文)，需要先转义
            // 流程: 原始字符串 -> encodeURIComponent (UTF-8转义) -> unescape (转为 Latin1) -> btoa
            output.value = btoa(unescape(encodeURIComponent(input.value)))
            debouncedReport('crypto', 'b64_encode_success')
        } else {
            // Base64 解码
            // 流程: atob -> escape -> decodeURIComponent
            output.value = decodeURIComponent(escape(atob(input.value)))
            debouncedReport('crypto', 'b64_decode_success')
        }
    } catch (e) {
        // 解码失败通常是因为输入不是合法的 Base64 字符串
        output.value = "错误: " + e.message
        debouncedReport('crypto', 'b64_error', e.message)
    }
}

// 复制结果
const copy = () => {
    navigator.clipboard.writeText(output.value)
}
</script>

<template>
  <div class="card">
    <div class="section-header">加解密 & Base64</div>
    
    <!-- 模式切换按钮 -->
    <div class="tab-bar">
       <div class="tab-btn" :class="{ active: mode === 'encode' }" @click="mode = 'encode'">编码 (Encode)</div>
       <div class="tab-btn" :class="{ active: mode === 'decode' }" @click="mode = 'decode'">解码 (Decode)</div>
    </div>

    <!-- 输入输出区域 -->
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