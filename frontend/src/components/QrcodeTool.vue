<script setup>
import { ref, watch } from 'vue'
import QRCode from 'qrcode'
import { debounce } from '../utils/debounce'

const props = defineProps(['reportEvent'])

const textInput = ref('')
const qrDataUrl = ref('')
const errorLevel = ref('M') // L, M, Q, H
const margin = ref(4)
const scale = ref(4)

const generate = debounce(async () => {
    if (!textInput.value.trim()) {
        qrDataUrl.value = ''
        return
    }
    
    try {
        qrDataUrl.value = await QRCode.toDataURL(textInput.value, {
            errorCorrectionLevel: errorLevel.value,
            margin: Number(margin.value),
            scale: Number(scale.value),
            width: 300 // 固定宽度方便展示，实际保存时可能需要按 scale
        })
        props.reportEvent('qrcode', 'generate')
    } catch (err) {
        console.error(err)
    }
}, 300)

watch([textInput, errorLevel, margin], generate)

const download = () => {
    if (!qrDataUrl.value) return
    const a = document.createElement('a')
    a.href = qrDataUrl.value
    a.download = 'qrcode.png'
    a.click()
    props.reportEvent('qrcode', 'download')
}
</script>

<template>
  <div class="card full-height">
    <div class="section-header">
        <span>二维码生成</span>
    </div>

    <div class="split-container">
        <!-- 左侧设置与输入 -->
        <div class="left-pane">
            <div class="form-group">
                <label>文本内容 / URL</label>
                <textarea 
                    v-model="textInput" 
                    class="text-input qr-input-fix" 
                    placeholder="在此输入文本..."
                ></textarea>
            </div>

            <div class="options-group">
                <div class="option-item">
                    <label>容错率</label>
                    <select v-model="errorLevel" class="select-box">
                        <option value="L">Low (7%)</option>
                        <option value="M">Medium (15%)</option>
                        <option value="Q">Quartile (25%)</option>
                        <option value="H">High (30%)</option>
                    </select>
                </div>
                
                 <div class="option-item">
                    <label>边距 (Margin)</label>
                    <select v-model="margin" class="select-box">
                        <option value="0">0</option>
                        <option value="1">1</option>
                        <option value="2">2</option>
                        <option value="4">4</option>
                    </select>
                </div>
            </div>
        </div>

        <!-- 右侧展示 -->
        <div class="right-pane">
            <div class="preview-box">
                <img v-if="qrDataUrl" :src="qrDataUrl" alt="QR Code" class="qr-img" />
                <div v-else class="placeholder">
                    输入文本以生成二维码
                </div>
            </div>
            
            <button v-if="qrDataUrl" class="btn btn-primary btn-block" @click="download">
                下载图片
            </button>
        </div>
    </div>
  </div>
</template>

<style scoped>
/* Global .full-height used */

.split-container {
    flex: 1;
    display: flex;
    gap: 32px;
    padding: 24px;
}

.left-pane {
    flex: 2;
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.right-pane {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
    background: #f9fafb;
    padding: 24px;
    border-radius: var(--radius-md);
    border: 1px solid var(--border-color);
}

.form-group {
    display: flex;
    flex-direction: column;
    flex: 1;
}

.form-group label, .option-item label {
    font-size: 14px;
    font-weight: 500;
    color: #374151;
    margin-bottom: 8px;
}

.text-input {
    flex: 1;
    resize: none;
    /* Visuals inherited from global textarea */
}

.text-input:focus {
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.options-group {
    display: flex;
    gap: 24px;
}

.option-item {
    display: flex;
    flex-direction: column;
}

.select-box {
    padding: 8px 12px;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    outline: none;
}

.preview-box {
    width: 200px;
    height: 200px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: white;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    border-radius: 4px;
}

.qr-img {
    max-width: 100%;
    max-height: 100%;
}

.placeholder {
    color: #9ca3af;
    font-size: 13px;
    text-align: center;
    padding: 20px;
}

.btn-block {
    width: 100%;
    max-width: 200px;
    padding: 10px;
}
</style>
