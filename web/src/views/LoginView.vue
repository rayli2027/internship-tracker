<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { api } from '../api'
import { useAuth } from '../auth'
import { useI18n } from '../i18n'
import AuthSplit from '../components/AuthSplit.vue'
import { isValidEmail } from '../utils'

const REMEMBER_KEY = 'it_remember_account'
const SAVED_EMAIL_KEY = 'it_saved_email'

const router = useRouter()
const route = useRoute()
const { setSession } = useAuth()
const { t, language } = useI18n()
const loading = ref(false)
const rememberAccount = ref(false)
const form = reactive({ email: '', password: '', preferred_language: language.value })

onMounted(() => {
  rememberAccount.value = localStorage.getItem(REMEMBER_KEY) === '1'
  if (rememberAccount.value) {
    form.email = localStorage.getItem(SAVED_EMAIL_KEY) || ''
  }
})

function persistRememberedEmail() {
  if (rememberAccount.value) {
    localStorage.setItem(REMEMBER_KEY, '1')
    localStorage.setItem(SAVED_EMAIL_KEY, form.email)
  } else {
    localStorage.removeItem(REMEMBER_KEY)
    localStorage.removeItem(SAVED_EMAIL_KEY)
  }
}

async function submit() {
  if (!isValidEmail(form.email)) {
    ElMessage.error(language.value === 'en' ? 'Enter a valid email address.' : '请输入有效的邮箱地址。')
    return
  }
  if (!form.password) {
    ElMessage.error(t('emailPasswordRequired'))
    return
  }
  loading.value = true
  try {
    form.preferred_language = language.value
    const data = await api.login(form)
    persistRememberedEmail()
    setSession(data.token, data.user)
    ElMessage.success(t('loginSuccess'))
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    router.push(data.user.role === 'admin' && redirect === '/' ? '/admin/overview' : redirect)
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : t('loginFailed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthSplit>
    <div class="auth-card-head">
      <h1>{{ t('login') }}</h1>
      <p class="auth-subtitle">{{ language === 'en' ? 'Sign in to your workspace' : '登录您的账户以继续' }}</p>
    </div>
    <form class="auth-form" @submit.prevent="submit">
      <el-form label-position="top">
        <el-form-item
          :label="t('email')"
          :error="form.email && !isValidEmail(form.email) ? (language === 'en' ? 'Invalid email format' : '邮箱格式不正确') : ''"
        >
          <el-input v-model.trim="form.email" type="email" autocomplete="username" />
        </el-form-item>
        <el-form-item :label="t('password')">
          <el-input v-model="form.password" type="password" show-password autocomplete="current-password" />
        </el-form-item>
        <el-form-item class="auth-remember-item">
          <el-checkbox v-model="rememberAccount">{{ t('rememberAccount') }}</el-checkbox>
        </el-form-item>
      </el-form>
      <el-button class="auth-submit" type="primary" native-type="submit" :loading="loading">{{ t('login') }}</el-button>
    </form>
    <p class="auth-foot">{{ t('noAccount') }} <router-link to="/register">{{ t('register') }}</router-link></p>
  </AuthSplit>
</template>
