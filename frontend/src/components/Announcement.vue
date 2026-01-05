<script setup>
</script>

<template>
  <div class="announcement-container">
    <h1>公告</h1>
    
    <div class="update-card">
      <div class="version-tag">v1.0.6</div>
      <div class="date">2026-01-04</div>
      <ul class="change-list">
        <li>✨ <strong>全新日期选择器</strong>: 引入 Flatpickr，交互更流畅。</li>
        <li>📈 <strong>图表升级</strong>: 统计页趋势趋势优化，视觉更统一。</li>
      </ul>
    </div>

    <div class="update-card">
      <div class="version-tag">v1.0.5</div>
      <div class="date">2026-01-03</div>
      <ul class="change-list">
        <li>✨ <strong>UI 重塑</strong>: 回归中文，采用 Apple 风格 (San Francisco 字体 / 圆角 / 阴影)。</li>
        <li>🕒 <strong>时间戳增强</strong>: 增加日期在此转换功能 (Date -> Timestamp)。</li>
        <li>📊 <strong>图表优化</strong>: 升级为平滑曲线面积图 (Smooth Area Chart)。</li>
      </ul>
    </div>

    <div class="update-card">
      <div class="version-tag">v1.0.4</div>
      <div class="date">2025-12-29</div>
      <ul class="change-list">
        <li>✨ <strong>新工具</strong>: 上线 SQL 格式化、文本对比、二维码、颜色转换、进制转换、文本工具。</li>
        <li>🛠 <strong>优化</strong>: 导航栏重构，支持移动端适配。</li>
      </ul>
    </div>

  </div>
</template>

<style scoped>
.announcement-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

h1 {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 24px;
  color: #111827;
}

.update-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 20px;
  position: relative;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.version-tag {
  display: inline-block;
  background: #eff6ff;
  color: #2563eb;
  padding: 4px 12px;
  border-radius: 99px;
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 12px;
}

.date {
  position: absolute;
  top: 28px;
  right: 24px;
  color: #6b7280;
  font-size: 13px;
}

.change-list {
  margin: 0;
  padding-left: 20px;
  color: #374151;
}

.change-list li {
  margin-bottom: 8px;
  line-height: 1.6;
}

strong {
  color: #111827;
  font-weight: 600;
}
</style>