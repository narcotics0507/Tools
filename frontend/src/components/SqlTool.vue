<script setup>
import { ref, shallowRef } from 'vue'
import { format } from 'sql-formatter'
import { VueMonacoEditor } from '@guolao/vue-monaco-editor'

const props = defineProps(['reportEvent'])

// ==========================================
// 1. 状态定义
// ==========================================
const sqlInput = ref('')
const formattedSql = ref('')
const dialect = ref('mysql') // 默认 SQL 方言 (改为 MySQL 以支持更多常见语法)

const editorRef = shallowRef()

// Monaco Editor 配置
const EDITOR_OPTIONS = {
    automaticLayout: true,
    formatOnType: true,
    formatOnPaste: true,
    minimap: { enabled: false },
    readOnly: true, // 只读，用于展示格式化结果
    scrollBeyondLastLine: false,
    wordWrap: 'on',
    scrollbar: {
        verticalScrollbarSize: 6,
        horizontalScrollbarSize: 6,
        useShadows: false
    }
}

const handleMount = (editor) => {
    editorRef.value = editor
}

// ==========================================
// 2. 逻辑处理
// ==========================================

const formatSql = () => {
    if (!sqlInput.value.trim()) return
    try {
        formattedSql.value = format(sqlInput.value, {
            language: dialect.value,
            tabWidth: 4,
            keywordCase: 'upper',
        })
        props.reportEvent('sql', 'format')
    } catch (e) {
        alert('格式化失败: ' + e.message)
    }
}

const compressSql = () => {
     if (!sqlInput.value.trim()) return
     // 简单的压缩逻辑：去除多余空白
     formattedSql.value = sqlInput.value.replace(/\s+/g, ' ').trim()
     props.reportEvent('sql', 'compress')
}

const copyResult = () => {
    if (formattedSql.value) {
        navigator.clipboard.writeText(formattedSql.value)
        alert('已复制结果')
        props.reportEvent('sql', 'copy')
    }
}

// 示例 SQL
const inputExample = () => {
    sqlInput.value = "SELECT * FROM users WHERE id = 1 AND name = 'test' ORDER BY created_at DESC"
    formatSql()
}

</script>

<template>
  <div class="card full-height">
    <div class="section-header">
        <span>SQL 格式化</span>
    </div>

    <!-- 控制栏 -->
    <div class="controls-bar">
       <div class="left-group">
          <span class="label-text">方言:</span>
          <select v-model="dialect">
              <option value="sql">Standard SQL</option>
              <option value="mysql">MySQL</option>
              <option value="postgresql">PostgreSQL</option>
              <option value="hive">Hive</option>
          </select>
          <button class="btn btn-text" @click="inputExample">填入示例</button>
       </div>

       <div class="actions-group">
          <button class="btn btn-outline" @click="compressSql">压缩</button>
          <button class="btn btn-primary" @click="formatSql">格式化</button>
          <button class="btn btn-outline" @click="copyResult">复制结果</button>
       </div>
    </div>

    <div class="dual-pane">
      <!-- 左侧：输入 -->
      <div class="pane left-pane">
        <textarea 
            class="code-area" 
            v-model="sqlInput" 
            placeholder="在此粘贴 SQL 语句..."
        ></textarea>
      </div>
      
      <!-- 右侧：Monaco Output -->
      <div class="pane right-pane">
         <VueMonacoEditor
            v-if="formattedSql"
            v-model:value="formattedSql"
            :language="dialect"
            theme="vs-light"
            :options="EDITOR_OPTIONS"
            @mount="handleMount"
            class="sql-viewer"
         />
         <div v-else class="placeholder-text">
            等待格式化...
         </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Global .full-height used */

.controls-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    margin-bottom: 12px;
    border-bottom: 1px solid var(--border-color);
}

.left-group, .actions-group {
    display: flex;
    align-items: center;
    gap: 12px;
}

.label-text {
    font-size: 13px;
    color: #4b5563;
    white-space: nowrap; /* 强制不换行 */
}



.dual-pane {
  flex: 1;
  display: flex;
  gap: 16px;
  overflow: hidden; 
}

.pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  background: white;
}

.code-area { 
    flex: 1; 
    resize: none; 
    border: none;
    outline: none;
    padding: 12px;
    font-family: var(--font-mono); 
    font-size: 13px; 
    line-height: 1.5;
}

.right-pane {
    background: #f9fafb;
}

.sql-viewer {
    width: 100%;
    height: 100%;
}

.placeholder-text {
    padding: 40px;
    color: #9ca3af;
    text-align: center;
}

/* Removed .select-box, .btn-sm, .btn-text to use global styles */
</style>
