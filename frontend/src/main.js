import { createApp } from 'vue'
import App from './App.vue'
import './style.css' // <--- 关键！引入刚才创建的全局样式

createApp(App).mount('#app')