<script setup>
import { ref } from 'vue'
import tinycolor from 'tinycolor2'

const props = defineProps(['reportEvent'])

const inputColor = ref('#2563eb')
const hex = ref('')
const rgb = ref('')
const hsl = ref('')
const previewStyle = ref({})

// 更新所有颜色格式
const updateColors = (val) => {
    const color = tinycolor(val)
    if (color.isValid()) {
        hex.value = color.toHexString()
        rgb.value = color.toRgbString()
        hsl.value = color.toHslString()
        previewStyle.value = {
            backgroundColor: hex.value,
            color: color.isLight() ? '#000' : '#fff' // 自动适配文字颜色
        }
    }
}

// 初始化
updateColors(inputColor.value)

const onInput = (e) => {
    const val = e.target.value
    inputColor.value = val
    updateColors(val)
}

const copy = (text) => {
    navigator.clipboard.writeText(text)
    alert('已复制: ' + text)
    props.reportEvent('color', 'copy')
}

// 快速预设颜色
const swatches = [
    '#ef4444', '#f97316', '#f59e0b', '#84cc16', '#22c55e', 
    '#06b6d4', '#2563eb', '#6366f1', '#a855f7', '#ec4899', 
    '#111827', '#6b7280', '#ffffff'
]

const selectSwatch = (c) => {
    inputColor.value = c
    updateColors(c)
}

</script>

<template>
  <div class="card full-height">
    <div class="section-header">
        <span>颜色转换</span>
    </div>

    <div class="color-tool-container">
        <!-- 颜色预览与选择器 -->
        <div class="preview-card" :style="previewStyle">
            <div class="preview-text">Color Preview</div>
            <input type="color" class="color-picker" :value="hex" @input="onInput" />
        </div>

        <!-- 数值展示区域 -->
        <div class="values-container">
            <div class="input-group">
                <label>HEX</label>
                <div class="input-wrapper">
                    <input type="text" v-model="hex" @input="onInput" spellcheck="false" />
                    <button @click="copy(hex)">复制</button>
                </div>
            </div>

            <div class="input-group">
                <label>RGB</label>
                 <div class="input-wrapper">
                    <input type="text" v-model="rgb" @input="onInput" spellcheck="false" />
                    <button @click="copy(rgb)">复制</button>
                </div>
            </div>

            <div class="input-group">
                <label>HSL</label>
                 <div class="input-wrapper">
                    <input type="text" v-model="hsl" @input="onInput" spellcheck="false" />
                    <button @click="copy(hsl)">复制</button>
                </div>
            </div>
        </div>
        
        <!-- 常用色板 -->
        <div class="swatches-area">
            <div class="label">常用色板</div>
            <div class="swatches-grid">
                <div 
                    v-for="c in swatches" 
                    :key="c" 
                    class="swatch" 
                    :style="{ background: c }"
                    @click="selectSwatch(c)"
                ></div>
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

.color-tool-container {
    max-width: 600px;
    margin: 0 auto;
    width: 100%;
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 32px;
}

.preview-card {
    height: 120px;
    border-radius: var(--radius-md);
    box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1);
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    border: 1px solid rgba(0,0,0,0.1);
    transition: background 0.2s;
}

.preview-text {
    font-weight: 600;
    font-size: 18px;
    opacity: 0.9;
}

.color-picker {
    position: absolute;
    bottom: 12px;
    right: 12px;
    opacity: 0; 
    /* 实际上我们希望点击预览区域能触发选择，或者仅仅用一个小图标。
       这里简单处理：让 Input type=color 覆盖在右下角或者整个区域？
       为了美观，通常隐藏原本的 input，用自定义按钮触发。
       但原生 input type=color 最简单。让我们放在右下角可见。
    */
    width: 32px;
    height: 32px;
    padding: 0;
    border: 2px solid white;
    border-radius: 50%;
    cursor: pointer;
    opacity: 1;
    overflow: hidden;
    box-shadow: 0 2px 4px rgba(0,0,0,0.2);
}
.color-picker::-webkit-color-swatch-wrapper {
	padding: 0;
}
.color-picker::-webkit-color-swatch {
	border: none;
    border-radius: 50%;
}

.values-container {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.input-group label {
    display: block;
    margin-bottom: 6px;
    font-size: 13px;
    color: #6b7280;
    font-weight: 500;
}

/* Input wrapper styles moved to global style.css */

.swatches-area .label {
    font-size: 13px;
    color: #6b7280;
    margin-bottom: 12px;
}

.swatches-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
}

.swatch {
    width: 36px;
    height: 36px;
    border-radius: 8px;
    cursor: pointer;
    box-shadow: 0 1px 2px rgba(0,0,0,0.1);
    border: 2px solid transparent;
    transition: transform 0.1s;
}

.swatch:hover {
    transform: scale(1.1);
    border-color: white;
    box-shadow: 0 4px 6px rgba(0,0,0,0.2);
}
</style>
