<script setup>
import { ref, computed, watch } from 'vue'
import { debounce } from '../utils/debounce'

const props = defineProps(['reportEvent'])

// ==========================================
// 1. 状态定义
// ==========================================
const pattern = ref('\\d+')         // 正则表达式模式
const flags = ref('g')              // 正则修饰符 (如 g, i, m)
const testString = ref('Hello 123 World 456') // 待测试字符串

// ==========================================
// 2. 核心逻辑: 实时匹配与高亮
// ==========================================

// 计算属性：根据正则匹配结果生成带有 HTML 高亮的字符串
const resultHtml = computed(() => {
    if (!pattern.value) return testString.value
    try {
        const regex = new RegExp(pattern.value, flags.value)
        const content = testString.value
        
        let html = ''
        let lastIndex = 0
        
        // 安全的 HTML 转义函数，防止 XSS
        const escapeHtml = (unsafe) => {
             return unsafe
                 .replace(/&/g, "&amp;")
                 .replace(/</g, "&lt;")
                 .replace(/>/g, "&gt;")
                 .replace(/"/g, "&quot;")
                 .replace(/'/g, "&#039;");
        }

        // 使用 replace 的回调函数来逐个处理匹配项
        content.replace(regex, (match, ...args) => {
            //倒数第二个参数在 replace 回调中通常是匹配项在原字符串中的偏移量 offset
            const offset = args[args.length - 2] 
            
            // 1. 拼接匹配项之前的内容 (普通文本)
            html += escapeHtml(content.slice(lastIndex, offset))
            // 2. 拼接高亮处理后的匹配项
            html += `<span class="match">${escapeHtml(match)}</span>`
            
            // 更新下次截取的起始位置
            lastIndex = offset + match.length
            
            // 如果不是全局匹配 (g)，replace 只会执行一次，这里返回 match 不影响原逻辑
            return match 
        })
        
        // 3. 拼接剩余未匹配的文本
        html += escapeHtml(content.slice(lastIndex))
        
        return html
    } catch (e) {
        // 如果正则语法错误，直接忽略，不crash页面，也可以显示错误提示
        return `<span class="error">正则错误: ${e.message}</span>`
    }
})

// ==========================================
// 3. 事件上报 (防抖)
// ==========================================

// 创建防抖上报函数
const debouncedReport = debounce(() => {
    props.reportEvent('regex', 'test')
}, 2000)

// 监听正则或测试文本的变化，触发上报
watch([pattern, flags, testString], () => {
    if (pattern.value) {
        debouncedReport()
    }
})
</script>

<template>
  <div class="card">
    <div class="section-header">正则表达式测试</div>
    
    <!-- 输入行: Pattern 和 Flags -->
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
    
    <!-- 双栏布局: 测试文本 vs 匹配结果 -->
    <div class="dual-pane">
      <div class="pane">
        <label>测试字符串</label>
        <textarea class="code-area" v-model="testString"></textarea>
      </div>
      <div class="pane">
        <label>匹配结果</label>
        <!-- 使用 v-html 渲染包含高亮 span 的结果 -->
        <div class="code-area result-area" v-html="resultHtml"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.input-row { display: flex; gap: 16px; margin-bottom: 16px; }
.code-area { flex: 1; resize: none; font-size: 14px; padding: 12px; }
.result-area { background: white; overflow: auto; white-space: pre-wrap; word-break: break-all; border: 1px solid #e5e7eb; }

label { font-size: 12px; font-weight: 500; color: #6b7280; margin-bottom: 4px; display: block; }

/* 深度选择器，用于匹配 v-html 动态生成的 .match 类 */
:deep(.match) { background-color: #fef08a; border-radius: 2px; box-shadow: 0 0 0 1px #facc15; }
:deep(.error) { color: #ef4444; font-weight: bold; }
</style>
