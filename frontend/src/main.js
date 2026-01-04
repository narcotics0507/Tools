import { createApp } from 'vue'
import App from './App.vue'
import './style.css' // <--- 引入全局 CSS 样式 (Tailwind 或自定义样式)

// 创建 Vue 应用实例并挂载到 index.html 的 #app 节点上
createApp(App).mount('#app')