<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../../api'
import type { Overview } from '../../types'

const data = ref<Overview | null>(null)
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    data.value = await api.adminOverview()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div v-loading="loading">
    <p class="eyebrow">概览</p>
    <h1 class="page-title">运营概览</h1>
    <p class="subtitle">看用户和投递规模，不做审批。</p>
    <div class="kpi-grid" style="margin-top: 20px">
      <article class="kpi">
        <span>用户数</span>
        <strong>{{ data?.user_count ?? 0 }}</strong>
      </article>
      <article class="kpi">
        <span>投递总数</span>
        <strong>{{ data?.application_count ?? 0 }}</strong>
      </article>
      <article class="kpi">
        <span>今日新增用户</span>
        <strong>{{ data?.today_user_count ?? 0 }}</strong>
      </article>
      <article class="kpi">
        <span>今日新增投递</span>
        <strong>{{ data?.today_application_count ?? 0 }}</strong>
      </article>
    </div>
  </div>
</template>
