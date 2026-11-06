<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import AppHeader from '../components/AppHeader.vue'
import { api } from '../api'
import { useAuth } from '../auth'
import { useI18n } from '../i18n'
import { useTheme } from '../theme'

const { user, setSession } = useAuth()
const { language } = useI18n()
const saving = ref(false)
const { mode, setMode } = useTheme()
const form = reactive({
  username: user.value?.username || '',
  password: '',
  preferred_language: user.value?.preferred_language || language.value,
})
const en = () => language.value === 'en'

async function save() {
  saving.value = true
  try {
    const data = await api.updateProfile(form)
    setSession(localStorage.getItem('it_token') || '', data.user)
    form.password = ''
    ElMessage.success(en() ? 'Settings saved.' : '设置已保存。')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : (en() ? 'Unable to save settings.' : '保存失败。'))
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="app-shell">
    <AppHeader />
    <main class="workspace settings-page">
      <header class="settings-hero">
        <h1>{{ en() ? 'Profile settings' : '个人设置' }}</h1>
        <p class="subtitle">{{ en() ? 'Manage your profile and workspace preferences.' : '管理个人资料与工作区偏好。' }}</p>
      </header>

      <section class="settings-section">
        <div class="settings-section-head">
          <h2>{{ en() ? 'Appearance' : '外观' }}</h2>
          <p>{{ en() ? 'Choose the color scheme that is easiest on your eyes.' : '选择更适合你的界面配色。' }}</p>
        </div>
        <div class="theme-options">
          <button type="button" :class="{ active: mode === 'light' }" @click="setMode('light')">
            <span class="theme-preview light"><i></i><i></i><i></i></span>
            {{ en() ? 'Light' : '浅色' }}
          </button>
          <button type="button" :class="{ active: mode === 'dark' }" @click="setMode('dark')">
            <span class="theme-preview dark"><i></i><i></i><i></i></span>
            {{ en() ? 'Dark' : '深色' }}
          </button>
        </div>
      </section>

      <section class="settings-section">
        <div class="settings-section-head">
          <h2>{{ en() ? 'Account' : '账户信息' }}</h2>
          <p>{{ en() ? 'Update how you appear in the workspace and your sign-in details.' : '更新工作区中的显示名称与登录信息。' }}</p>
        </div>
        <el-form label-position="top" class="settings-form" @submit.prevent="save">
          <div class="settings-form-grid">
            <el-form-item :label="en() ? 'Display name' : '用户名'">
              <el-input v-model="form.username" />
            </el-form-item>
            <el-form-item :label="en() ? 'New password' : '新密码'">
              <el-input
                v-model="form.password"
                type="password"
                show-password
                :placeholder="en() ? 'Leave empty to keep your password' : '留空则不修改密码'"
              />
            </el-form-item>
            <el-form-item :label="en() ? 'Interface language' : '界面语言'">
              <el-radio-group v-model="form.preferred_language">
                <el-radio value="en">English</el-radio>
                <el-radio value="zh">中文</el-radio>
              </el-radio-group>
            </el-form-item>
          </div>
          <div class="settings-actions">
            <el-button type="primary" native-type="submit" :loading="saving">{{ en() ? 'Save settings' : '保存设置' }}</el-button>
            <el-button @click="$router.push('/')">{{ en() ? 'Back to board' : '返回看板' }}</el-button>
          </div>
        </el-form>
      </section>
    </main>
  </div>
</template>
