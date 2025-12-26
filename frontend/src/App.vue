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
</script>

<template>
  <div id="app-container">
    <nav class="sidebar">
      <div class="brand" @click="currentTab = 'announcement'">
        <div class="brand-icon">S</div>
        上校工具箱
      </div>
      <ul class="nav-menu">
        <li v-for="tab in tabs" :key="tab.id" class="nav-item" :class="{ active: currentTab === tab.id }" @click="currentTab = tab.id">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" v-html="tab.icon"></svg> {{ tab.name }}
        </li>
      </ul>
      <div style="padding:16px; font-size:11px; color:#9ca3af; text-align:center; margin-top: auto;">v1.0.2 (Modular)</div>
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