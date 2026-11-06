<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import AppHeader from '../components/AppHeader.vue'
import ApplicationDrawer from '../components/ApplicationDrawer.vue'
import { api } from '../api'
import { useI18n } from '../i18n'
import type { Application, ApplicationEvent } from '../types'
import { fmtTime } from '../utils'

const { language, statusLabel } = useI18n()
const en = computed(() => language.value === 'en')
const events = ref<ApplicationEvent[]>([])
const loading = ref(false)
const keyword = ref('')
const eventType = ref('')
const drawerOpen = ref(false)
const current = ref<Application | null>(null)

const groups = computed(() => {
  const query = keyword.value.trim().toLowerCase()
  const map = new Map<string, { applicationId: string; company: string; position: string; events: ApplicationEvent[] }>()
  for (const event of events.value) {
    if (eventType.value && event.event_type !== eventType.value) continue
    const searchable = `${event.company_name} ${event.position_name} ${event.title} ${event.content}`.toLowerCase()
    if (query && !searchable.includes(query)) continue
    if (!map.has(event.application_id)) {
      map.set(event.application_id, {
        applicationId: event.application_id,
        company: event.company_name || (en.value ? 'Unknown company' : '未知公司'),
        position: event.position_name || (en.value ? 'Unknown position' : '未知岗位'),
        events: [],
      })
    }
    map.get(event.application_id)!.events.push(event)
  }
  return [...map.values()]
    .map((group) => ({
      ...group,
      events: group.events.sort((a, b) => new Date(a.event_at).getTime() - new Date(b.event_at).getTime()),
    }))
    .sort(
      (a, b) =>
        new Date(b.events[b.events.length - 1]?.event_at || 0).getTime() -
        new Date(a.events[a.events.length - 1]?.event_at || 0).getTime(),
    )
})

const filterOptions = computed(() => [
  { v: 'created', en: 'Created', zh: '创建' },
  { v: 'status_change', en: 'Status change', zh: '状态变化' },
  { v: 'applied', en: 'Applied', zh: '已投递' },
  { v: 'oa', en: 'Assessment', zh: '笔试' },
  { v: 'interview', en: 'Interview', zh: '面试' },
  { v: 'offer', en: 'Offer', zh: 'Offer' },
  { v: 'rejected', en: 'Rejected', zh: '被拒' },
  { v: 'withdrawn', en: 'Withdrawn', zh: '主动放弃' },
  { v: 'memo', en: 'Note', zh: '备忘' },
])

async function load() {
  loading.value = true
  try {
    events.value = (await api.listTimelineEvents()).items
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : en.value ? 'Unable to load lifecycles' : '无法加载生命周期')
  } finally {
    loading.value = false
  }
}

async function openApplication(id: string) {
  try {
    current.value = (await api.getApplication(id)).item
    drawerOpen.value = true
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : en.value ? 'Application no longer exists' : '岗位不存在')
  }
}

function label(event: ApplicationEvent) {
  if (event.event_type === 'status_change' && event.to_status) {
    try {
      return statusLabel(event.to_status as any)
    } catch {
      return event.to_status
    }
  }
  const labels = en.value
    ? { created: 'Created', applied: 'Applied', oa: 'Assessment', interview: 'Interview', offer: 'Offer', rejected: 'Rejected', withdrawn: 'Withdrawn', memo: 'Note', status_change: 'Status change' }
    : { created: '创建', applied: '已投递', oa: '笔试', interview: '面试', offer: 'Offer', rejected: '被拒', withdrawn: '主动放弃', memo: '备忘', status_change: '状态变化' }
  return (labels as Record<string, string>)[event.event_type] || event.event_type
}

function detailText(event: ApplicationEvent) {
  return Object.entries(event.details || {})
    .filter(([, value]) => value)
    .map(([, value]) => String(value))
    .join(' · ')
}

onMounted(load)
</script>

<template>
  <div class="app-shell">
    <AppHeader />
    <main class="workspace timeline-page" v-loading="loading">
      <section class="board-toolbar-strip">
        <el-input
          v-model="keyword"
          class="board-search"
          clearable
          :prefix-icon="Search"
          :placeholder="en ? 'Search company, position, or notes' : '搜索公司、岗位或记录'"
        />
        <el-select v-model="eventType" clearable :placeholder="en ? 'All stages' : '全部阶段'" style="width: 160px">
          <el-option v-for="item in filterOptions" :key="item.v" :label="en ? item.en : item.zh" :value="item.v" />
        </el-select>
      </section>

      <section v-if="groups.length" class="application-lifecycle-list">
        <article v-for="group in groups" :key="group.applicationId" class="page-panel application-lifecycle">
          <button class="lifecycle-heading" type="button" @click="openApplication(group.applicationId)">
            <span>
              <strong>{{ group.company }}</strong>
              <small>{{ group.position }}</small>
            </span>
            <span class="lifecycle-open">{{ en ? 'Open →' : '查看 →' }}</span>
          </button>
          <div class="process-scroll">
            <ol class="process-line">
              <li v-for="event in group.events" :key="event.id" :class="`stage-${event.event_type}`">
                <span class="process-dot"></span>
                <time>{{ fmtTime(event.event_at).slice(0, 10) }}</time>
                <strong>{{ label(event) }}</strong>
                <p v-if="event.content || detailText(event)">
                  {{ detailText(event) }}
                  <template v-if="event.content && detailText(event)"> · </template>
                  {{ event.content }}
                </p>
              </li>
            </ol>
          </div>
        </article>
      </section>
      <div v-else class="page-panel list-empty">
        {{ en ? 'No lifecycle events match these filters.' : '没有符合条件的生命周期记录。' }}
      </div>
    </main>
    <ApplicationDrawer v-model="drawerOpen" :item="current" @saved="load" @deleted="load" />
  </div>
</template>
