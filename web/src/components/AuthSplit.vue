<script setup lang="ts">
import { computed } from 'vue'
import { ArrowDown, Briefcase, Moon, Sunny } from '@element-plus/icons-vue'
import { useI18n } from '../i18n'
import { useTheme } from '../theme'
import type { PreferredLanguage } from '../types'

const { language, setGuestLanguage } = useI18n()
const { isDark, toggleMode } = useTheme()

const languageLabel = computed(() =>
  language.value === 'en' ? 'Language - English' : '语言 - 中文',
)

const brandHeadline = computed(() =>
  language.value === 'en'
    ? 'Manage internship applications with clarity'
    : '系统化追踪实习投递全流程',
)

const brandCopy = computed(() =>
  language.value === 'en'
    ? 'A professional workspace for students and recruiters to track every application from wishlist to offer.'
    : '面向学生与招聘团队的专业工作区，从意向岗位到 Offer 全程可视、可追踪。',
)

const features = computed(() =>
  language.value === 'en'
    ? ['Board and list views for pipeline tracking', 'Deadline and lifecycle reminders', 'Resume workspace and data export']
    : ['看板与列表双视图，阶段流转清晰可查', '截止日期与生命周期事件提醒', '简历管理与数据导出'],
)

function onLanguageChange(value: string | number | object) {
  setGuestLanguage(value as PreferredLanguage)
}
</script>

<template>
  <main class="auth-split">
    <div class="auth-toolbar">
      <el-dropdown trigger="click" @command="onLanguageChange">
        <button class="auth-lang-trigger" type="button" :aria-label="language === 'en' ? 'Switch language' : '切换语言'">
          <span>{{ languageLabel }}</span>
          <el-icon><ArrowDown /></el-icon>
        </button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="zh">中文</el-dropdown-item>
            <el-dropdown-item command="en">English</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <button
        class="theme-toggle"
        type="button"
        :aria-label="isDark ? (language === 'en' ? 'Use light mode' : '切换浅色模式') : (language === 'en' ? 'Use dark mode' : '切换深色模式')"
        :title="isDark ? (language === 'en' ? 'Light mode' : '浅色模式') : (language === 'en' ? 'Dark mode' : '深色模式')"
        @click="toggleMode"
      >
        <el-icon><Sunny v-if="isDark" /><Moon v-else /></el-icon>
      </button>
    </div>

    <aside class="auth-brand-panel">
      <div class="auth-brand-panel-inner">
        <div class="auth-brand auth-brand--panel">
          <span class="brand-logo brand-logo--panel"><el-icon><Briefcase /></el-icon></span>
          <span class="brand-name">Internship Tracker</span>
        </div>
        <h2 class="auth-brand-headline">{{ brandHeadline }}</h2>
        <p class="auth-brand-copy">{{ brandCopy }}</p>
        <ul class="auth-feature-list">
          <li v-for="item in features" :key="item">{{ item }}</li>
        </ul>
      </div>
      <p class="auth-brand-foot">{{ language === 'en' ? 'Enterprise-grade application tracking' : '企业级投递管理与协作' }}</p>
    </aside>

    <section class="auth-form-area">
      <div class="auth-form-panel">
        <slot />
      </div>
    </section>
  </main>
</template>
