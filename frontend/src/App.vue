<script setup>
import { ref } from 'vue'
import Announcement from './components/Announcement.vue'
import TimestampTool from './components/TimestampTool.vue'
import JsonTool from './components/JsonTool.vue'
import CronTool from './components/CronTool.vue'
import UrlTool from './components/UrlTool.vue'
import CryptoTool from './components/CryptoTool.vue'
import Stats from './components/Stats.vue'
import Base64Tool from './components/Base64Tool.vue'
import RegexTool from './components/RegexTool.vue'

const currentTab = ref('announcement')
const isMobileMenuOpen = ref(false)

const tabs = [
  { id: 'announcement', name: '公告页', icon: '<path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>' },
  { id: 'timestamp', name: '时间戳工具', icon: '<circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>' },
  { id: 'url', name: 'URL 编解码', icon: '<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>' },
  { id: 'json', name: 'JSON 工具', icon: '<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>' },
  { id: 'cron', name: 'Crond 生成器', icon: '<path d="M5 22h14"/><path d="M5 2h14"/><path d="M17 22v-4a2 2 0 0 0-2-2H9a2 2 0 0 0-2 2v4"/>' },
  { id: 'base64', name: 'Base64 转换', icon: '<path d="M4 7V4h16v3M9 20h6M12 4v16"/>' },
  { id: 'regex', name: '正则测试', icon: '<path d="M21 12H3M7 8l-4 4 4 4M17 8l4 4-4 4"/>' },
  { id: 'crypto', name: '加解密', icon: '<rect x="3" y="11" width="18" height="11" rx="2"></rect>' },
  { id: 'stats', name: '工具调用统计', icon: '<circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>' }
]

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

const toggleMobileMenu = () => {
    isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const selectTab = (id) => {
    currentTab.value = id
    isMobileMenuOpen.value = false // Close menu on selection
}
</script>

<template>
  <div id="app-container">
    <!-- Mobile Header -->
    <div class="mobile-header">
        <div class="brand-mobile">
            <div class="brand-icon">S</div>
            上校工具箱
        </div>
        <button class="menu-toggle" @click="toggleMobileMenu">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="3" y1="12" x2="21" y2="12"></line><line x1="3" y1="6" x2="21" y2="6"></line><line x1="3" y1="18" x2="21" y2="18"></line></svg>
        </button>
    </div>

    <!-- Mobile Overlay -->
    <div v-if="isMobileMenuOpen" class="mobile-overlay" @click="isMobileMenuOpen = false"></div>

    <nav class="sidebar" :class="{ 'mobile-open': isMobileMenuOpen }">
      <div class="brand" @click="selectTab('announcement')">
        <div class="brand-icon">S</div>
        上校工具箱
      </div>
      <ul class="nav-menu">
        <li v-for="tab in tabs" :key="tab.id" class="nav-item" :class="{ active: currentTab === tab.id }" @click="selectTab(tab.id)">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" v-html="tab.icon"></svg> {{ tab.name }}
        </li>
      </ul>
      <div style="padding:16px; font-size:11px; color:#9ca3af; text-align:center; margin-top: auto;">v1.0.3 (Modular)</div>
    </nav>

    <main class="main-content">
        <Announcement v-if="currentTab === 'announcement'" />
        <TimestampTool v-if="currentTab === 'timestamp'" :report-event="reportEvent" />
        <JsonTool v-if="currentTab === 'json'" :report-event="reportEvent" />
        <CronTool v-if="currentTab === 'cron'" :report-event="reportEvent" />
        <UrlTool v-if="currentTab === 'url'" :report-event="reportEvent" />
        <CryptoTool v-if="currentTab === 'crypto'" :report-event="reportEvent" />
        <Base64Tool v-if="currentTab === 'base64'" :report-event="reportEvent" />
        <RegexTool v-if="currentTab === 'regex'" :report-event="reportEvent" />
        <Stats v-if="currentTab === 'stats'" />
    </main>
  </div>
</template>

<style>
/* Global Resets & Variables */
:root {
  --primary: #2563eb;
  --primary-hover: #1d4ed8;
  --bg-page: #f3f4f6;
  --bg-sidebar: #ffffff;
  --bg-card: #ffffff;
  --text-primary: #111827;
  --text-secondary: #6b7280;
  --border-color: #e5e7eb;
  --accent: #f59e0b; /* Amber for highlights */
  --radius-sm: 6px;
  --radius-md: 8px;
}

body {
  margin: 0;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  background-color: var(--bg-page);
  color: var(--text-primary);
  -webkit-font-smoothing: antialiased;
}

#app-container {
  display: flex;
  min-height: 100vh;
}

/* Sidebar */
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
  z-index: 50;
  transition: transform 0.3s ease;
}

.brand {
  height: 64px;
  display: flex;
  align-items: center;
  padding: 0 24px;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--primary);
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
}

.brand-icon {
    width: 32px;
    height: 32px;
    background: var(--primary);
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
  background-color: #eff6ff; /* Light blue */
  color: var(--primary);
}

/* Main Content */
.main-content {
  flex: 1;
  margin-left: 260px; /* Width of sidebar */
  padding: 32px;
  max-width: 100%;
}

/* Mobile Specifics */
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

/* Responsive Media Queries */
@media (max-width: 768px) {
    .sidebar {
        transform: translateX(-100%); /* Hide by default */
        box-shadow: 2px 0 8px rgba(0,0,0,0.1);
    }
    .sidebar.mobile-open {
        transform: translateX(0);
    }
    .main-content {
        margin-left: 0;
        margin-top: 60px; /* Header height */
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