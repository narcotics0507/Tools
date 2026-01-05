<script setup>
import { ref, computed } from 'vue'
import { 
    camelCase, 
    kebabCase, 
    snakeCase, 
    startCase, 
    toLower, 
    toUpper, 
    upperFirst 
} from 'lodash-es'

const props = defineProps(['reportEvent'])

const inputText = ref('')

// 自定义一些通过组合 lodash 函数实现的格式
const pascalCase = (str) => upperFirst(camelCase(str))
const constantCase = (str) => toUpper(snakeCase(str))
const dotCase = (str) => kebabCase(str).replace(/-/g, '.')
const pathCase = (str) => kebabCase(str).replace(/-/g, '/')

const conversions = computed(() => {
    const txt = inputText.value
    if (!txt) return []
    
    return [
        { label: 'Lowercase (小写)', value: toLower(txt) },
        { label: 'Uppercase (大写)', value: toUpper(txt) },
        { label: 'CamelCase (驼峰)', value: camelCase(txt) },
        { label: 'PascalCase (大驼峰)', value: pascalCase(txt) },
        { label: 'SnakeCase (下划线)', value: snakeCase(txt) },
        { label: 'KebabCase (短横线)', value: kebabCase(txt) },
        { label: 'ConstantCase (常量)', value: constantCase(txt) },
        { label: 'DotCase (点连接)', value: dotCase(txt) },
        { label: 'PathCase (路径)', value: pathCase(txt) },
        { label: 'TitleCase (标题)', value: startCase(camelCase(txt)) }, // startCase 稍微粗糙，先这样
    ]
})

const copy = (text) => {
    if(!text) return
    navigator.clipboard.writeText(text)
    // 简单的 Toast 提示，这里用 alert 代替
    props.reportEvent('text', 'copy')
}

</script>

<template>
  <div class="card full-height">
    <div class="section-header">
        <span>文本大小写转换</span>
    </div>

    <div class="split-container">
        <!-- Input Area -->
        <div class="input-pane">
            <div class="pane-header">输入文本</div>
            <textarea 
                v-model="inputText" 
                class="text-input" 
                placeholder="Type something here..."
                autofocus
            ></textarea>
        </div>

        <!-- Output List -->
        <div class="output-pane">
            <div class="pane-header">转换结果</div>
            <div class="results-list">
                <div v-for="(item, index) in conversions" :key="index" class="result-item">
                     <div class="label">{{ item.label }}</div>
                     <div class="input-wrapper">
                        <input type="text" readonly :value="item.value" @click="copy(item.value)" />
                        <button @click="copy(item.value)">Copy</button>
                     </div>
                </div>
                <div v-if="conversions.length === 0" class="placeholder">
                    输入内容并实时预览所有格式
                </div>
            </div>
        </div>
    </div>
  </div>
</template>

<style scoped>
.full-height {
    display: flex;
    flex-direction: column;
}

.split-container {
    flex: 1;
    display: flex;
    gap: 24px;
    padding: 24px;
    overflow: hidden;
}

.input-pane, .output-pane {
    flex: 1;
    display: flex;
    flex-direction: column;
    background: white;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    overflow: hidden;
}

.pane-header {
    padding: 12px 16px;
    background: #f9fafb;
    border-bottom: 1px solid var(--border-color);
    font-size: 13px;
    font-weight: 600;
    color: #4b5563;
}

.text-input {
    flex: 1;
    resize: none;
    border: none;
    outline: none;
    padding: 16px;
    font-size: 14px;
    line-height: 1.6;
    font-family: var(--font-mono);
}

.output-pane {
    background: #f9fafb;
}

.results-list {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.result-item {
    background: white;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    padding: 12px;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

.result-item .label {
    font-size: 12px;
    color: #9ca3af;
    margin-bottom: 8px;
}

/* Input wrapper styles moved to global style.css */

.placeholder {
    margin-top: 40px;
    text-align: center;
    color: #9ca3af;
    font-size: 13px;
}
</style>
