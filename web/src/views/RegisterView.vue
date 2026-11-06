<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { api } from '../api'
import { useAuth } from '../auth'
import { useI18n } from '../i18n'
import AuthSplit from '../components/AuthSplit.vue'
import type { PreferredLanguage } from '../types'
import { isValidEmail } from '../utils'

const router = useRouter()
const { setSession } = useAuth()
const { t, language } = useI18n()
const loading = ref(false)
const form = reactive<{ email: string; username: string; password: string; preferred_language: PreferredLanguage }>({ email: '', username: '', password: '', preferred_language: language.value })

async function submit() {
  if (!isValidEmail(form.email)) { ElMessage.error(language.value === 'en' ? 'Enter a valid email address.' : '请输入有效的邮箱地址。'); return }
  if (form.password.length < 6) { ElMessage.error(t('registrationRequired')); return }
  form.preferred_language = language.value
  loading.value = true
  try {
    const data = await api.register(form)
    setSession(data.token, data.user)
    ElMessage.success(t('registrationSuccess'))
    router.push('/')
  } catch (error) { ElMessage.error(error instanceof Error ? error.message : t('registrationFailed')) }
  finally { loading.value = false }
}
</script>

<template>
  <AuthSplit>
    <div class="auth-card-head">
      <h1>{{ t('register') }}</h1>
      <p class="auth-subtitle">{{ language === 'en' ? 'Create an account to start tracking' : '创建账户，开始管理投递' }}</p>
    </div>
    <form class="auth-form" @submit.prevent="submit">
      <el-form label-position="top">
        <el-form-item :label="t('email')" :error="form.email && !isValidEmail(form.email) ? (language === 'en' ? 'Invalid email format' : '邮箱格式不正确') : ''"><el-input v-model.trim="form.email" type="email" autocomplete="username" /></el-form-item>
        <el-form-item :label="t('username')"><el-input v-model="form.username" :placeholder="t('usernameHint')" /></el-form-item>
        <el-form-item :label="t('password')"><el-input v-model="form.password" type="password" show-password autocomplete="new-password" /></el-form-item>
      </el-form>
      <el-button class="auth-submit" type="primary" native-type="submit" :loading="loading">{{ t('createAccount') }}</el-button>
    </form>
    <p class="auth-foot">{{ t('haveAccount') }} <router-link to="/login">{{ t('login') }}</router-link></p>
  </AuthSplit>
</template>
