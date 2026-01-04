<script setup>
import { ref, shallowRef, watch } from 'vue'
import { VueMonacoDiffEditor } from '@guolao/vue-monaco-editor'

const props = defineProps(['reportEvent'])

const original = ref('')
const modified = ref('')

// Debounce logic could be added if performance is an issue, 
// but for text diffs, Monaco handles updates reasonably well.
// We will bind v-model directly.

const OPTIONS = {
  originalEditable: true, 
  readOnly: false,        
  renderSideBySide: true,
  minimap: { enabled: false }, // Disable minimap to save space in split view
  scrollBeyondLastLine: false,
  automaticLayout: true,
  scrollbar: {
    verticalScrollbarSize: 6,
    horizontalScrollbarSize: 6,
    verticalSliderSize: 6,
    horizontalSliderSize: 6,
    useShadows: false
  },
  wordWrap: 'on'
}

const editorRef = shallowRef()

const handleMount = (diffEditor) => {
  editorRef.value = diffEditor
  // Force update options to ensure scrollbar styles are applied
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

const copyResult = () => {
    navigator.clipboard.writeText(modified.value)
    alert('已复制修改后的文本')
    props.reportEvent('diff', 'copy')
}

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

    <!-- Toolbar -->
    <div class="toolbar">
        <div class="left-tip">
            ✨ 实时对比模式：请在下方输入文本，差异将自动显示。
        </div>

        <div class="btn-group">
            <button class="btn btn-sm btn-outline" @click="clear">清空</button>
            <button class="btn btn-sm btn-outline" @click="copyResult">复制右侧</button>
        </div>
    </div>

    <!-- Main Content: Split Vertical -->
    <div class="split-container">
        <!-- Top: Inputs -->
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

        <!-- Bottom: Diff Editor -->
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
.full-height {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 100px);
}

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

/* Update: Input take 35%, Editor take 65% */
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
    white-space: pre-wrap;       /* Preserve whitespace/returns */
    overflow-wrap: break-word;     /* Break long words if needed */
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

.btn-sm {
    padding: 4px 12px;
    font-size: 12px;
}

/* Custom Scrollbar */
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

/* Global Monaco Overrides - placed here to ensure they load with the component */
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
/* Hide the decorations ruler which can look like a thick border/shadow */
.monaco-editor .decorationsOverviewRuler {
    display: none !important;
}
/* Ensure the scrollbar container itself is not too wide/intrusive */
.monaco-editor .scrollbar.vertical {
    width: 6px !important;
    background: transparent !important;
    box-shadow: none !important;
}
</style>
