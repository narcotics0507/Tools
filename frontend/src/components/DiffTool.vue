<script setup>
import { ref, shallowRef } from 'vue'
// 使用 Vue 封装的 Monaco Diff Editor 组件
import { VueMonacoDiffEditor } from '@guolao/vue-monaco-editor'

const props = defineProps(['reportEvent'])

// ==========================================
// 1. 状态定义
// ==========================================
const original = ref('')  // 左侧原始文本
const modified = ref('')  // 右侧修改后文本

// 编辑器实例引用 (shallowRef 用于性能优化，避免深度响应式)
const editorRef = shallowRef()

// ==========================================
// 2. Monaco Editor 配置
// ==========================================
const OPTIONS = {
  originalEditable: true, // 允许编辑左侧原始文本 (默认通常是只读的，这里开启方便双向修改)
  readOnly: false,        // 允许编辑右侧
  renderSideBySide: true, // 并排显示模式 (Inline 为 false)
  minimap: { enabled: false }, // 禁用缩略图，节省空间
  scrollBeyondLastLine: false, // 禁止滚动超过最后一行
  automaticLayout: true,  // 自动响应容器大小变化
  // 自定义滚动条样式，使其更窄且不显眼
  scrollbar: {
    verticalScrollbarSize: 6,
    horizontalScrollbarSize: 6,
    verticalSliderSize: 6,
    horizontalSliderSize: 6,
    useShadows: false
  },
  wordWrap: 'on' // 开启自动换行
}

// ==========================================
// 3. 事件处理
// ==========================================

// 编辑器挂载完成后的回调
const handleMount = (diffEditor) => {
  editorRef.value = diffEditor
  
  // 强制更新子编辑器的选项，确保滚动条样式生效
  // Monaco Diff Editor 由两个独立的 Code Editor 组成
  const modifiedEditor = diffEditor.getModifiedEditor();
  const originalEditor = diffEditor.getOriginalEditor();
  
  const scrollbarOptions = {
    verticalScrollbarSize: 6,
    horizontalScrollbarSize: 6,
    verticalSliderSize: 6,
    horizontalSliderSize: 6,
    useShadows: false
  };

  modifiedEditor.updateOptions({ scrollbar: scrollbarOptions, wordWrap: 'on' });
  originalEditor.updateOptions({ scrollbar: scrollbarOptions, wordWrap: 'on' });
}

// 复制右侧结果
const copyResult = () => {
    navigator.clipboard.writeText(modified.value)
    alert('已复制修改后的文本')
    props.reportEvent('diff', 'copy')
}

// 清空两侧内容
const clear = () => {
    original.value = ''
    modified.value = ''
}
</script>

<template>
  <div class="card full-height">
    <div class="section-header">
        <span>文本差异比对</span>
    </div>

    <!-- 顶部工具栏 -->
    <div class="toolbar">
        <div class="left-tip">
            ✨ 实时对比模式：请在下方输入文本，差异将自动显示。
        </div>

        <div class="btn-group">
            <button class="btn btn-outline" @click="clear">清空</button>
            <button class="btn btn-outline" @click="copyResult">复制右侧</button>
        </div>
    </div>

    <!-- 主内容区：垂直分割布局 -->
    <div class="split-container">
        <!-- 上半部分：文本输入框 (各自占 35% 高度) -->
        <div class="input-container">
            <div class="input-box">
                <div class="label">原始文本</div>
                <textarea class="text-area" v-model="original" placeholder="在此粘贴..."></textarea>
            </div>
            <div class="input-box">
                <div class="label">需对比文本</div>
                <textarea class="text-area" v-model="modified" placeholder="在此粘贴..."></textarea>
            </div>
        </div>

        <!-- 下半部分：Diff 编辑器 (占剩余空间) -->
        <div class="editor-wrapper">
             <div class="label-bar">对比结果 (Diff View)</div>
             <VueMonacoDiffEditor
                v-model:original="original"
                v-model:modified="modified"
                language="json"
                theme="vs-light"
                :options="OPTIONS"
                @mount="handleMount"
                class="diff-editor"
            />
        </div>
    </div>
  </div>
</template>

<style scoped>
/* 撑满除了顶部导航外的高度 */
/* Global .full-height used */

.section-header {
    margin-bottom: 0;
    padding-bottom: 12px;
}

.toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: #f3f4f6;
    border: 1px solid var(--border-color);
    border-bottom: none;
    border-top-left-radius: var(--radius-sm);
    border-top-right-radius: var(--radius-sm);
}

.left-tip {
    font-size: 13px;
    color: #4b5563;
}

.btn-group {
    display: flex;
    gap: 8px;
}

.split-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    border: 1px solid var(--border-color);
    border-bottom-left-radius: var(--radius-sm);
    border-bottom-right-radius: var(--radius-sm);
    overflow: hidden;
}

/* 输入区域高度占比 35% */
.input-container {
    height: 35%; 
    display: flex;
    gap: 12px;
    padding: 12px;
    background: white;
    border-bottom: 1px solid var(--border-color);
}

.input-box {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.label {
    font-size: 12px;
    font-weight: 600;
    margin-bottom: 4px;
    color: #374151;
}

.text-area {
    flex: 1;
    resize: none;
    padding: 8px;
    font-family: var(--font-mono);
    font-size: 12px;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    outline: none;
    line-height: 1.4;
    white-space: pre-wrap;
    overflow-wrap: break-word;
}

.text-area:focus {
    border-color: var(--primary);
}

.editor-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    background: #fafafa;
}

.label-bar {
    padding: 4px 12px;
    font-size: 11px;
    color: #9ca3af;
    background: #f9fafb;
    border-bottom: 1px solid #f3f4f6;
}

.diff-editor {
    flex: 1;
    width: 100%;
}



/* 自定义滚动条 (Webkit) */
::-webkit-scrollbar {
    width: 6px;
    height: 6px;
}
::-webkit-scrollbar-track {
    background: transparent;
}
::-webkit-scrollbar-thumb {
    background: #e5e7eb;
    border-radius: 3px;
}
::-webkit-scrollbar-thumb:hover {
    background: #d1d5db;
}
</style>

<!-- 全局 Monaco 样式覆盖 (必须是非 scoped 才能生效) -->
<style>
.monaco-editor .scrollbar .slider {
    background: #e5e7eb !important;
    border-radius: 3px !important;
    width: 6px !important;
    left: 0 !important;
}
.monaco-editor .scrollbar .slider:hover {
    background: #d1d5db !important;
}
/* 隐藏 Monaco 编辑器左侧的装饰标尺，减少视觉干扰 */
.monaco-editor .decorationsOverviewRuler {
    display: none !important;
}
.monaco-editor .scrollbar.vertical {
    width: 6px !important;
    background: transparent !important;
    box-shadow: none !important;
}
</style>
