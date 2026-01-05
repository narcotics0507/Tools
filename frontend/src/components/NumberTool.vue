<script setup>
import { ref } from 'vue'

const props = defineProps(['reportEvent'])

const bin = ref('')
const oct = ref('')
const dec = ref('')
const hex = ref('')

// 更新所有字段
const updateAll = (value, radix) => {
    if (!value) {
        bin.value = ''
        oct.value = ''
        dec.value = ''
        hex.value = ''
        return
    }

    try {
        // 使用 BigInt 处理大数
        let num
        if (radix === 10) {
            num = BigInt(value)
        } else if (radix === 16) {
             // BigInt 构造函数接受 0x 前缀或无前缀如果是十进制，但这里是十六进制字符串
             // BigInt('0x' + value) 可以
             // 去掉可能的 0x 前缀
             const cleanHex = value.toLowerCase().replace(/^0x/, '')
             num = BigInt('0x' + cleanHex)
        } else if (radix === 2) {
             num = BigInt('0b' + value)
        } else if (radix === 8) {
             num = BigInt('0o' + value)
        }

        // 更新其他字段
        if (radix !== 2) bin.value = num.toString(2)
        if (radix !== 8) oct.value = num.toString(8)
        if (radix !== 10) dec.value = num.toString(10)
        if (radix !== 16) hex.value = num.toString(16).toUpperCase()

    } catch (e) {
        // 输入不合法时暂时忽略，或是显示错误状态 (这里简化处理，不做更新)
        // console.warn(e)
    }
}

const onInput = (type, e) => {
    const val = e.target.value.trim()
    
    // 简单的字符清理 (只允许对应进制的字符)
    let cleanVal = val
    if (type === 'bin') {
        cleanVal = val.replace(/[^01]/g, '')
        bin.value = cleanVal // 立即回显清理后的值
        updateAll(cleanVal, 2)
    } else if (type === 'oct') {
         cleanVal = val.replace(/[^0-7]/g, '')
         oct.value = cleanVal
         updateAll(cleanVal, 8)
    } else if (type === 'dec') {
         cleanVal = val.replace(/[^0-9]/g, '')
         dec.value = cleanVal
         updateAll(cleanVal, 10)
    } else if (type === 'hex') {
         cleanVal = val.replace(/[^0-9a-fA-F]/g, '')
         hex.value = cleanVal
         updateAll(cleanVal, 16)
    }
}

const copy = (text) => {
    if(!text) return
    navigator.clipboard.writeText(text)
    alert('已复制')
}

</script>

<template>
  <div class="card full-height">
    <div class="section-header">
        <span>进制转换</span>
    </div>

    <div class="converter-container">
        
        <!-- 10 进制 -->
        <div class="input-group">
            <label>十进制 (Decimal)</label>
            <div class="input-wrapper">
                <input class="input-box" type="text" :value="dec" @input="e => onInput('dec', e)" placeholder="0-9" />
                <button class="btn-copy" @click="copy(dec)">复制</button>
            </div>
        </div>
        
        <!-- 16 进制 -->
        <div class="input-group">
            <label>十六进制 (Hexadecimal)</label>
            <div class="input-wrapper">
                <input class="input-box" type="text" :value="hex" @input="e => onInput('hex', e)" placeholder="0-9, A-F" />
                <button class="btn-copy" @click="copy(hex)">复制</button>
            </div>
        </div>

        <!-- 2 进制 -->
        <div class="input-group">
            <label>二进制 (Binary)</label>
            <div class="input-wrapper">
                 <input class="input-box" type="text" :value="bin" @input="e => onInput('bin', e)" placeholder="0-1" />
                 <button class="btn-copy" @click="copy(bin)">复制</button>
            </div>
        </div>

        <!-- 8 进制 -->
        <div class="input-group">
            <label>八进制 (Octal)</label>
             <div class="input-wrapper">
                <input class="input-box" type="text" :value="oct" @input="e => onInput('oct', e)" placeholder="0-7" />
                <button class="btn-copy" @click="copy(oct)">复制</button>
             </div>
        </div>

    </div>
  </div>
</template>

<style scoped>
.converter-container {
    padding: 24px;
    max-width: 600px;
    margin: 0 auto;
    width: 100%;
}
.input-group { margin-bottom: 24px; }
label {
    display: block; margin-bottom: 8px; font-size: 14px; font-weight: 500; color: #374151;
}
/* Input wrapper styles moved to global style.css */
.input-box {
    flex: 1; border: none; outline: none; padding: 0 12px; font-size: 14px; font-family: var(--font-mono); color: var(--text-primary);
}

</style>
