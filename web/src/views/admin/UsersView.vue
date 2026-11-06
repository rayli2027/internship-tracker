<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../../api'
import type { User } from '../../types'
import { fmtTime } from '../../utils'

const loading = ref(false)
const items = ref<User[]>([])
const total = ref(0)
const filters = reactive({
  keyword: '',
  page: 1,
  page_size: 20,
})

async function load() {
  loading.value = true
  try {
    const data = await api.adminUsers(filters)
    items.value = data.items
    total.value = data.total
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '加载失败')
  } finally {
    loading.value = false
  }
}

function search() {
  filters.page = 1
  load()
}

function reset() {
  filters.keyword = ''
  filters.page = 1
  filters.page_size = 20
  load()
}

async function toggle(row: User, enabled: boolean) {
  try {
    await api.adminSetUserEnabled(row.id, enabled)
    ElMessage.success(enabled ? '已启用' : '已禁用')
    await load()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '更新失败')
  }
}

onMounted(load)
</script>

<template>
  <div>
    <p class="eyebrow">Users</p>
    <h1 class="page-title">用户管理</h1>
    <p class="subtitle">搜索后点查询。禁用后该用户不能登录。</p>
    <div class="filter-bar">
      <el-input v-model="filters.keyword" placeholder="邮箱或用户名" clearable style="width: 240px" @keyup.enter="search" />
      <el-button type="primary" @click="search">查询</el-button>
      <el-button @click="reset">重置</el-button>
    </div>
    <el-table
      :data="items"
      :empty-text="loading ? ' ' : '暂无数据'"
      v-loading="loading"
      :cell-style="{ fontSize: '12px' }"
      :header-cell-style="{ fontSize: '12px', fontWeight: 600 }"
    >
      <el-table-column prop="id" label="用户 ID" width="180" align="center" show-overflow-tooltip />
      <el-table-column prop="email" label="邮箱" min-width="200" show-overflow-tooltip />
      <el-table-column prop="username" label="用户名" width="140" align="center" show-overflow-tooltip />
      <el-table-column prop="role" label="角色" width="90" align="center" show-overflow-tooltip />
      <el-table-column label="状态" width="90" align="center" show-overflow-tooltip>
        <template #default="{ row }">{{ row.is_enabled ? '启用' : '禁用' }}</template>
      </el-table-column>
      <el-table-column label="创建时间" width="168" align="center" show-overflow-tooltip>
        <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="120" align="center">
        <template #default="{ row }">
          <el-button v-if="row.is_enabled" size="small" @click="toggle(row, false)">禁用</el-button>
          <el-button v-else size="small" type="primary" @click="toggle(row, true)">启用</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="pagination-bar">
      <span class="pagination-total">总计 {{ total }} 条</span>
      <el-pagination
        v-model:current-page="filters.page"
        v-model:page-size="filters.page_size"
        layout="prev, pager, next, sizes"
        :page-sizes="[20, 50, 100]"
        :total="total"
        @current-change="load"
        @size-change="search"
      />
    </div>
  </div>
</template>
