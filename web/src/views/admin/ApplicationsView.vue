<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../../api'
import type { Application, User } from '../../types'
import { fmtSalary, fmtTime } from '../../utils'

const loading = ref(false)
const items = ref<Application[]>([])
const users = ref<User[]>([])
const total = ref(0)
const filters = reactive({
  keyword: '',
  user_id: '',
  page: 1,
  page_size: 20,
})

async function loadUsers() {
  const data = await api.adminUsers({ page: 1, page_size: 100 })
  users.value = data.items
}

async function load() {
  loading.value = true
  try {
    const data = await api.adminApplications(filters)
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
  filters.user_id = ''
  filters.page = 1
  filters.page_size = 20
  load()
}

onMounted(async () => {
  await loadUsers()
  await load()
})
</script>

<template>
  <div>
    <p class="eyebrow">Inspection</p>
    <h1 class="page-title">投递巡查</h1>
    <p class="subtitle">只读查看全站投递，不在后台改状态。</p>
    <div class="filter-bar">
      <el-select v-model="filters.user_id" clearable placeholder="按用户筛选" style="width: 220px">
        <el-option v-for="user in users" :key="user.id" :label="`${user.username} / ${user.email}`" :value="String(user.id)" />
      </el-select>
      <el-input v-model="filters.keyword" placeholder="公司、职位或用户名" clearable style="width: 240px" @keyup.enter="search" />
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
      <el-table-column prop="company_name" label="公司" min-width="160" show-overflow-tooltip />
      <el-table-column prop="position_name" label="职位" min-width="180" show-overflow-tooltip />
      <el-table-column prop="status" label="状态" width="110" align="center" show-overflow-tooltip />
      <el-table-column prop="user_name" label="用户" width="120" align="center" show-overflow-tooltip />
      <el-table-column label="薪酬" width="140" align="right" show-overflow-tooltip>
        <template #default="{ row }">
          <span class="tabular-nums">{{ fmtSalary(row.salary_min, row.salary_max, row.salary_unit) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="年限" width="78" align="center" show-overflow-tooltip>
        <template #default="{ row }">{{ row.experience_years || '不限' }}</template>
      </el-table-column>
      <el-table-column prop="education" label="学历" width="78" align="center" show-overflow-tooltip />
      <el-table-column prop="work_mode" label="办公" width="78" align="center" show-overflow-tooltip />
      <el-table-column label="喜爱" width="78" align="center" show-overflow-tooltip>
        <template #default="{ row }">{{ row.is_liked ? '是' : '否' }}</template>
      </el-table-column>
      <el-table-column prop="location" label="地点" min-width="140" show-overflow-tooltip />
      <el-table-column label="截止日期" width="168" align="center" show-overflow-tooltip>
        <template #default="{ row }">{{ fmtTime(row.deadline_at) }}</template>
      </el-table-column>
      <el-table-column prop="resume_version" label="简历版本" min-width="160" show-overflow-tooltip />
      <el-table-column label="更新时间" width="168" align="center" show-overflow-tooltip>
        <template #default="{ row }">{{ fmtTime(row.updated_at) }}</template>
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
