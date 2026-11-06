<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Moon, Setting, Sunny, SwitchButton, User } from '@element-plus/icons-vue'
import { useAuth } from '../auth'
import { useI18n } from '../i18n'
import { useTheme } from '../theme'

const router = useRouter()
const { user, isAdmin, clearSession } = useAuth()
const { t, language } = useI18n()
const { isDark, toggleMode } = useTheme()
const initial = computed(() => (user.value?.username || user.value?.email || 'U').slice(0, 1).toUpperCase())

function command(value: string) {
  if (value === 'logout') {
    clearSession()
    router.push('/login')
    return
  }
  router.push(value)
}
</script>

<template>
  <header class="topbar">
    <div class="brand-cluster">
      <router-link class="brand" to="/">
        <span class="brand-name">Internship Tracker</span>
      </router-link>
      <nav class="primary-nav" :aria-label="language === 'en' ? 'Main features' : '主要功能'">
        <router-link to="/">{{ language === 'en' ? 'Board' : '看板' }}</router-link>
        <router-link to="/resume">{{ language === 'en' ? 'Resume' : '简历' }}</router-link>
        <router-link to="/calendar">{{ language === 'en' ? 'Calendar' : '日历' }}</router-link>
        <router-link to="/timeline">{{ language === 'en' ? 'Lifecycle' : '生命周期' }}</router-link>
      </nav>
    </div>
    <div class="topbar-actions">
      <button
        class="theme-toggle"
        type="button"
        :aria-label="isDark ? (language === 'en' ? 'Use light mode' : '切换浅色模式') : (language === 'en' ? 'Use dark mode' : '切换深色模式')"
        @click="toggleMode"
      >
        <el-icon><Sunny v-if="isDark" /><Moon v-else /></el-icon>
      </button>
      <el-button v-if="!isAdmin && $route.path.startsWith('/admin')" text @click="router.push('/')">{{ t('board') }}</el-button>
      <el-dropdown trigger="click" @command="command">
        <button class="user-menu-trigger" type="button" :aria-label="language === 'en' ? 'Open account menu' : '打开账户菜单'">
          <span class="user-avatar">{{ initial }}</span>
          <span class="user-summary"><strong>{{ user?.username || user?.email }}</strong><small>{{ user?.email }}</small></span>
        </button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="/settings" :icon="Setting">{{ language === 'en' ? 'Settings' : '设置' }}</el-dropdown-item>
            <el-dropdown-item v-if="isAdmin" command="/admin/overview" :icon="User">{{ t('admin') }}</el-dropdown-item>
            <el-dropdown-item divided command="logout" :icon="SwitchButton">{{ t('logout') }}</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>
