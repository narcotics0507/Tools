<script setup>
import { ref } from 'vue'

// 引入各个功能组件
import Announcement from './components/Announcement.vue'
import TimestampTool from './components/TimestampTool.vue'
import JsonTool from './components/JsonTool.vue'
import CronTool from './components/CronTool.vue'
import UrlTool from './components/UrlTool.vue'
import CryptoTool from './components/CryptoTool.vue'
import Stats from './components/Stats.vue'
import RegexTool from './components/RegexTool.vue'
import UuidTool from './components/UuidTool.vue'
import HashTool from './components/HashTool.vue'

import DiffTool from './components/DiffTool.vue'
import SqlTool from './components/SqlTool.vue'
import NumberTool from './components/NumberTool.vue'
import QrcodeTool from './components/QrcodeTool.vue'
import ColorTool from './components/ColorTool.vue'
import TextTool from './components/TextTool.vue'

// ==========================================
// 1. 状态管理
// ==========================================

// 当前选中的 Tab (默认为公告页)
const currentTab = ref('announcement')

// 移动端菜单的开关状态
const isMobileMenuOpen = ref(false)

// ==========================================
// 2. 菜单配置
// ==========================================

// 定义侧边栏菜单项 (ID, 名称, 图标 SVG)
const tabs = [
  { id: 'announcement', name: '公告页', icon: '<path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>' },
  { id: 'diff', name: '文本对比', icon: '<path d="M16 3h5v5M4 20L21 3M21 16v5h-5M3 21l8.5-8.5"></path>'},
  { id: 'sql', name: 'SQL 格式化', icon: '<path d="M4 6c0 1.66 3.58 3 8 3s8-1.34 8-3-3.58-3-8-3-8 1.34-8 3"/><path d="M4 6v6c0 1.66 3.58 3 8 3s8-1.34 8-3V6"/><path d="M4 12v6c0 1.66 3.58 3 8 3s8-1.34 8-3v-6"/>' },
  { id: 'number', name: '进制转换', icon: '<path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/>' },
  { id: 'qrcode', name: '二维码生成', icon: '<rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect>' },
  { id: 'color', name: '颜色转换', icon: '<circle cx="12" cy="12" r="10"/><path d="M12 2a9.96 9.96 0 0 1 6.3 3.4L12 12V2z"/><path d="M12 12l6.3 6.6A9.96 9.96 0 0 1 12 22V12z"/><path d="M12 12V2a9.96 9.96 0 0 0-6.3 3.4L12 12z"/><path d="M12 22a9.96 9.96 0 0 1-6.3-3.4L12 12v10z"/>' },
  { id: 'text', name: '文本大小写', icon: '<path d="M4 7V4h16v3"/><path d="M9 20h6"/><path d="M12 4v16"/>' },
  { id: 'timestamp', name: '时间戳工具', icon: '<circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>' },
  { id: 'url', name: 'URL 编解码', icon: '<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>' },
  { id: 'json', name: 'JSON 工具', icon: '<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>' },
  { id: 'cron', name: 'Crond 生成器', icon: '<path d="M5 22h14"/><path d="M5 2h14"/><path d="M17 22v-4a2 2 0 0 0-2-2H9a2 2 0 0 0-2 2v4"/>' },
  { id: 'regex', name: '正则测试', icon: '<path d="M21 12H3M7 8l-4 4 4 4M17 8l4 4-4 4"/>' },
  { id: 'uuid', name: 'UUID 生成器', icon: '<rect x="3" y="3" width="18" height="18" rx="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><circle cx="15.5" cy="8.5" r="1.5"></circle><circle cx="15.5" cy="15.5" r="1.5"></circle><circle cx="8.5" cy="15.5" r="1.5"></circle>' },
  { id: 'hash', name: '哈希计算器', icon: '<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>' },
  { id: 'crypto', name: '加解密 & Base64', icon: '<rect x="3" y="11" width="18" height="11" rx="2"></rect>' },
  { id: 'stats', name: '工具调用统计', icon: '<circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>' }
]

// ==========================================
// 3. 通用方法
// ==========================================

/**
 * 上报用户行为事件到后端
 * @param {string} tool - 工具名称 (如 json, diff)
 * @param {string} action - 动作名称 (如 operate, click)
 * @param {string} errorMsg - 错误信息 (可选)
 */
const reportEvent = (tool, action, errorMsg = '') => {
  fetch('/api/event', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      tool_name: tool,
      action: action,
      is_error: !!errorMsg,
      error_msg: errorMsg
    })
  }).catch(err => console.log('上报忽略'))
}

// 切换移动端菜单显示/隐藏
const toggleMobileMenu = () => {
    isMobileMenuOpen.value = !isMobileMenuOpen.value
}

// 选中 Tab 并关闭移动端菜单
const selectTab = (id) => {
    currentTab.value = id
    isMobileMenuOpen.value = false // 选中后自动收起菜单
}
</script>

<template>
  <div id="app-container">
    <!-- ========================================== -->
    <!-- 移动端顶部导航栏 (仅在小屏幕显示) -->
    <!-- ========================================== -->
    <div class="mobile-header">
        <div class="brand-mobile">
            <div class="brand-icon">S</div>
            上校工具箱
        </div>
        <button class="menu-toggle" @click="toggleMobileMenu">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="3" y1="12" x2="21" y2="12"></line><line x1="3" y1="6" x2="21" y2="6"></line><line x1="3" y1="18" x2="21" y2="18"></line></svg>
        </button>
    </div>

    <!-- 移动端菜单遮罩层 -->
    <div v-if="isMobileMenuOpen" class="mobile-overlay" @click="isMobileMenuOpen = false"></div>

    <!-- ========================================== -->
    <!-- 侧边栏导航 (Sidebar) -->
    <!-- ========================================== -->
    <nav class="sidebar" :class="{ 'mobile-open': isMobileMenuOpen }">
      <div class="brand" @click="selectTab('announcement')">
        <div class="brand-icon">S</div>
        上校工具箱
      </div>
      <ul class="nav-menu">
        <!-- 遍历生成菜单项 -->
        <li v-for="tab in tabs" :key="tab.id" class="nav-item" :class="{ active: currentTab === tab.id }" @click="selectTab(tab.id)">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" v-html="tab.icon"></svg> {{ tab.name }}
        </li>
      </ul>
      <div style="padding:16px; font-size:11px; color:#9ca3af; text-align:center; margin-top: auto;">v1.0.6</div>
    </nav>

    <!-- ========================================== -->
    <!-- 主内容区域 (Main Content) -->
    <!-- ========================================== -->
    <main class="main-content">
        <!-- 动态渲染当前选中的组件 -->
        <Announcement v-if="currentTab === 'announcement'" />
        <DiffTool v-if="currentTab === 'diff'" :report-event="reportEvent" />
        <TimestampTool v-if="currentTab === 'timestamp'" :report-event="reportEvent" />
        <JsonTool v-if="currentTab === 'json'" :report-event="reportEvent" />
        <CronTool v-if="currentTab === 'cron'" :report-event="reportEvent" />
        <UrlTool v-if="currentTab === 'url'" :report-event="reportEvent" />
        <CryptoTool v-if="currentTab === 'crypto'" :report-event="reportEvent" />

        <RegexTool v-if="currentTab === 'regex'" :report-event="reportEvent" />
        <UuidTool v-if="currentTab === 'uuid'" :report-event="reportEvent" />
        <HashTool v-if="currentTab === 'hash'" :report-event="reportEvent" />

        <Stats v-if="currentTab === 'stats'" />

        <SqlTool v-if="currentTab === 'sql'" :report-event="reportEvent" />
        <NumberTool v-if="currentTab === 'number'" :report-event="reportEvent" />
        <QrcodeTool v-if="currentTab === 'qrcode'" :report-event="reportEvent" />
        <ColorTool v-if="currentTab === 'color'" :report-event="reportEvent" />
        <TextTool v-if="currentTab === 'text'" :report-event="reportEvent" />
    </main>
  </div>
</template>

<style>
/* ========================================== */
/* 全局样式与变量定义 */
/* ========================================== */
:root {
  --primary: #007aff;       /* Apple Blue */
  --primary-hover: #0071e3;
  --bg-page: #f5f5f7;       /* Apple Light Gray */
  --bg-sidebar: #ffffff;
  --bg-card: #ffffff;
  --text-primary: #1d1d1f;  /* Apple Dark */
  --text-secondary: #86868b;
  --border-color: #d1d1d6;
  --accent: #f59e0b;
  --radius-sm: 8px;
  --radius-md: 12px;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Helvetica, Arial, sans-serif;
  background-color: var(--bg-page);
  color: var(--text-primary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale; /* Add this for better font rendering */
  text-rendering: optimizeLegibility; /* Improve readability */
}

#app-container {
  display: flex;
  min-height: 100vh;
}

/* ========================================== */
/* 侧边栏样式 */
/* ========================================== */
.sidebar {
  width: 260px;
  background: var(--bg-sidebar);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  z-index: 100;
  transition: transform 0.3s ease;
}

.brand {
  height: 64px;
  display: flex;
  align-items: center;
  padding: 0 24px;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--text-primary); /* Change from var(--primary) to var(--text-primary) which is black */
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
}

.brand-icon {
    width: 32px;
    height: 32px;
    background: #1d1d1f; /* Apple Black */
    color: white;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 12px;
    font-weight: 800;
    font-size: 1.2rem;
}

.nav-menu {
  list-style: none;
  padding: 16px 12px;
  margin: 0;
  flex: 1;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  margin-bottom: 4px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  color: var(--text-secondary);
  font-weight: 500;
  transition: all 0.2s;
  font-size: 0.95rem;
}

.nav-item svg {
  width: 20px;
  height: 20px;
  margin-right: 12px;
  stroke-width: 2;
}

.nav-item:hover {
  background-color: #f9fafb;
  color: var(--text-primary);
}

.nav-item.active {
  background-color: #eff6ff; /* 选中状态淡蓝背景 */
  color: var(--primary);
}

/* ========================================== */
/* 主内容区域样式 */
/* ========================================== */
.main-content {
  flex: 1;
  margin-left: 260px; /* 留出侧边栏宽度 */
  padding: 32px;
  max-width: 100%;
}

/* ========================================== */
/* 移动端适配样式 (Max-width: 768px) */
/* ========================================== */
.mobile-header {
    display: none;
    height: 60px;
    background: white;
    border-bottom: 1px solid var(--border-color);
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 40;
}

.brand-mobile {
    display: flex;
    align-items: center;
    font-weight: 700;
    color: var(--primary);
    font-size: 1.1rem;
}
.brand-mobile .brand-icon {
    width: 28px;
    height: 28px;
    font-size: 1rem;
    margin-right: 8px;
}

.menu-toggle {
    background: none;
    border: none;
    color: var(--text-primary);
    cursor: pointer;
    padding: 4px;
}

.mobile-overlay {
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0,0,0,0.5);
    z-index: 45;
}

@media (max-width: 768px) {
    .sidebar {
        transform: translateX(-100%); /* 默认隐藏侧边栏 */
        box-shadow: 2px 0 8px rgba(0,0,0,0.1);
    }
    .sidebar.mobile-open {
        transform: translateX(0); /* 打开时滑入 */
    }
    .main-content {
        margin-left: 0;
        margin-top: 60px; /* 留出顶部导航高度 */
        padding: 16px;
    }
    .mobile-header {
        display: flex;
    }
    .mobile-overlay {
        display: block;
    }
}
</style>
