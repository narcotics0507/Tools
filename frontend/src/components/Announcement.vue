<script setup>
const logs = [
  { version: '1.0.4', date: '2025-12-27', changes: ['修复移动端显示异常问题','新增UUID生成器'] },
  { version: '1.0.3', date: '2025-12-26', changes: ['新增“工具调用统计”功能', 'JSON 工具改为左右分栏布局', '新增 Base64、正则测试工具'] },
  { version: '1.0.2', date: '2025-12-25', changes: ['Crond支持反解析功能'] },
  { version: '1.0.1', date: '2025-12-22', changes: ['Crond 每月执行日期回归为 31 天网格点击选择模式', '集成：URL、JSON、Crypto 模块完整逻辑接入'] },
  { version: '1.0.0', date: '2025-12-02', changes: ['上校工具箱 v1.0.0 正式上线'] }
]
</script>

<template>
  <div class="card">
    <div class="section-header">系统公告</div>
    <div v-for="log in logs" :key="log.version" style="margin-bottom: 20px; border-left: 2px solid var(--accent); padding-left: 15px;">
      <div style="font-weight:700;">Version {{ log.version }}</div>
      <div style="font-size:11px; color:#999; margin-bottom:8px;">{{ log.date }}</div>
      <ul style="margin:0; padding-left:15px; color:#555;">
        <li v-for="item in log.changes">{{ item }}</li>
      </ul>
    </div>
  </div>
</template>