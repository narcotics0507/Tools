<script setup>
import { ref, computed } from 'vue'

const props = defineProps(['reportEvent'])

const pattern = ref('\\d+')
const flags = ref('g')
const testString = ref('Hello 123 World 456')

const resultHtml = computed(() => {
    if (!pattern.value) return testString.value
    try {
        const regex = new RegExp(pattern.value, flags.value)
        const content = testString.value
        
        // Highlight logic
        // We wrap matches in specific spans.
        // A simple way is to use replace callback, but we need to handle non-global regex nicely too.
        
        let html = ''
        let lastIndex = 0
        
        // Safe html escape
        const escapeHtml = (unsafe) => {
             return unsafe
                 .replace(/&/g, "&amp;")
                 .replace(/</g, "&lt;")
                 .replace(/>/g, "&gt;")
                 .replace(/"/g, "&quot;")
                 .replace(/'/g, "&#039;");
        }

        content.replace(regex, (match, ...args) => {
            const offset = args[args.length - 2] // offset is second to last arg in replace
            // Append text before match
            html += escapeHtml(content.slice(lastIndex, offset))
            // Append highlighted match
            html += `<span class="match">${escapeHtml(match)}</span>`
            lastIndex = offset + match.length
            
            // If not global, replace only runs once naturally
            return match 
        })
        
        // Append remaining text
        html += escapeHtml(content.slice(lastIndex))
        
        // Trigger report (optional: consider debouncing)
        // props.reportEvent('regex', 'test') 
        
        return html
    } catch (e) {
        return `<span class="error">正则错误: ${e.message}</span>`
    }
})
</script>

<template>
  <div class="card">
    <div class="section-header">正则表达式测试</div>
    
    <div class="input-row">
        <div style="flex: 3;">
            <label>正则表达式 Pattern</label>
            <input type="text" v-model="pattern" class="code-area text-input" placeholder="\d+">
        </div>
        <div style="flex: 1;">
            <label>修饰符 Flags</label>
            <input type="text" v-model="flags" class="code-area text-input" placeholder="g">
        </div>
    </div>
    
    <div class="dual-pane">
      <div class="pane">
        <label>测试字符串</label>
        <textarea class="code-area" v-model="testString"></textarea>
      </div>
      <div class="pane">
        <label>匹配结果</label>
        <div class="code-area result-area" v-html="resultHtml"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.input-row { display: flex; gap: 16px; margin-bottom: 16px; }
.text-input { height: 38px; }
.dual-pane { display: flex; gap: 16px; height: 300px; }
.pane { flex: 1; display: flex; flex-direction: column; gap: 8px; }
.code-area { flex: 1; resize: none; font-size: 14px; padding: 12px; }
.result-area { background: white; overflow: auto; white-space: pre-wrap; word-break: break-all; border: 1px solid #e5e7eb; }

label { font-size: 12px; font-weight: 500; color: #6b7280; margin-bottom: 4px; display: block; }

:deep(.match) { background-color: #fef08a; border-radius: 2px; box-shadow: 0 0 0 1px #facc15; }
:deep(.error) { color: #ef4444; font-weight: bold; }
</style>
