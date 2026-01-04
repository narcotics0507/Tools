<script setup>
import { ref } from 'vue'

const props = defineProps(['reportEvent'])

// ==========================================
// 状态定义
// ==========================================
const count = ref(1)    // 生成数量
const result = ref('') // 结果字符串

// ==========================================
// 逻辑处理
// ==========================================

// 生成 UUID (v4)
const generate = () => {
    try {
        const uuids = []
        for (let i = 0; i < count.value; i++) {
            // 使用浏览器原生的 crypto.randomUUID()
            // 相比旧的 Math.random() 实现，这个更加安全且唯一性更有保障
            uuids.push(crypto.randomUUID())
        }
        // 多行显示
        result.value = uuids.join('\n')
        
        // 上报生成数量作为元数据
        props.reportEvent('uuid', 'generate', count.value.toString())
    } catch (e) {
        // 部分旧浏览器可能不支持 crypto.randomUUID (需要 localhost 或 HTTPS 上下文)
        result.value = "Error: " + e.message
        props.reportEvent('uuid', 'error', e.message)
    }
}

// 复制结果
const copy = () => {
    if (!result.value) return
    navigator.clipboard.writeText(result.value)
}
</script>

<template>
  <div class="card">
    <div class="section-header">UUID 生成器</div>
    
    <div class="controls">
        <label>
            生成数量:
            <input type="number" v-model="count" min="1" max="100" class="input-count">
        </label>
        <button class="btn btn-primary" @click="generate">生成 UUID</button>
        <button class="btn btn-outline" @click="copy">复制结果</button>
    </div>

    <textarea class="result-area" v-model="result" readonly placeholder="生成的 UUID 将显示在这里..."></textarea>
  </div>
</template>

<style scoped>
.controls {
    display: flex;
    gap: 16px;
    align-items: center;
    margin-bottom: 16px;
    background: #f9fafb;
    padding: 16px;
    border-radius: 8px;
    border: 1px solid #e5e7eb;
}

.input-count {
    padding: 6px 10px;
    border: 1px solid #d1d5db;
    border-radius: 4px;
    width: 60px;
}

.result-area {
    width: 100%;
    height: 300px;
    padding: 12px;
    border: 1px solid #d1d5db;
    border-radius: 8px;
    font-family: monospace;
    resize: vertical;
    outline: none;
}
.result-area:focus { border-color: var(--primary); }
</style>
