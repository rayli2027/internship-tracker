<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Calendar, Clock } from '@element-plus/icons-vue'
import { api } from '../api'
import type { Application, ApplicationEvent } from '../types'
import { useI18n } from '../i18n'

const applications = ref<Application[]>([])
const events = ref<ApplicationEvent[]>([])
const { language } = useI18n()
const startOfToday = () => { const date = new Date(); date.setHours(0, 0, 0, 0); return date }
const dateKey = (value: string) => { const date = new Date(value); return Number.isNaN(date.getTime()) ? '' : `${date.getFullYear()}-${String(date.getMonth()+1).padStart(2,'0')}-${String(date.getDate()).padStart(2,'0')}` }
const todayKey = computed(() => dateKey(startOfToday().toISOString()))
const isRealDate = (value:string) => { const date=new Date(value); return !Number.isNaN(date.getTime()) && date.getFullYear()>1971 }
const reminders = computed(() => {
  const today = startOfToday().getTime(), limit = today + 7 * 86400000
  return [
    ...applications.value.filter(item => item.deadline_at).map(item => ({ id:`d-${item.id}`, kind:'deadline', at:item.deadline_at, company:item.company_name })),
    ...applications.value.filter(item => isRealDate(item.applied_at) && !events.value.some(event => event.application_id===item.id && (event.event_type==='applied'||(event.event_type==='status_change'&&event.to_status==='Applied')) && dateKey(event.event_at)===dateKey(item.applied_at))).map(item => ({ id:`a-${item.id}`, kind:'workflow', at:item.applied_at, company:item.company_name })),
    ...events.value.filter(item => item.event_type !== 'created').map(item => ({ id:`e-${item.id}`, kind:item.event_type === 'interview' ? 'interview' : 'workflow', at:item.event_at, company:item.company_name || '' })),
  ].filter(item => { const value = new Date(item.at).getTime(); return value >= today && value < limit }).sort((a,b) => new Date(a.at).getTime()-new Date(b.at).getTime())
})
const today = computed(() => reminders.value.filter(item => dateKey(item.at) === todayKey.value))
const upcoming = computed(() => reminders.value.filter(item => dateKey(item.at) !== todayKey.value).slice(0, 4))
const format = (value: string) => new Intl.DateTimeFormat(language.value === 'en' ? 'en-US' : 'zh-CN', { month:'short', day:'numeric', hour:'2-digit', minute:'2-digit' }).format(new Date(value))
onMounted(async () => { try { const [a,e] = await Promise.all([api.listApplications(), api.listTimelineEvents()]); applications.value=a.items; events.value=e.items } catch { /* Board has its own error state. */ } })
</script>

<template>
  <section class="reminder-strip">
    <div class="reminder-intro"><span class="reminder-icon"><el-icon><Calendar /></el-icon></span><div><p class="eyebrow">{{ language === 'en' ? 'YOUR DAY' : '今日安排' }}</p><h2>{{ language === 'en' ? 'Today & coming up' : '今日待办与临期提醒' }}</h2></div></div>
    <div class="reminder-group"><strong>{{ language === 'en' ? 'Today' : '今天' }}</strong><span v-if="!today.length" class="reminder-empty">{{ language === 'en' ? 'Nothing scheduled' : '暂无安排' }}</span><router-link v-for="item in today" :key="item.id" to="/calendar" class="reminder-pill"><span :class="`reminder-dot ${item.kind}`"></span>{{ item.company }} · {{ item.kind === 'deadline' ? (language === 'en' ? 'Deadline' : '截止') : item.kind === 'interview' ? (language === 'en' ? 'Interview' : '面试') : (language === 'en' ? 'Workflow' : '流程') }}</router-link></div>
    <div class="reminder-group"><strong>{{ language === 'en' ? 'Next 7 days' : '未来 7 天' }}</strong><span v-if="!upcoming.length" class="reminder-empty">{{ language === 'en' ? 'You are all clear' : '暂无临期事项' }}</span><router-link v-for="item in upcoming" :key="item.id" to="/calendar" class="reminder-pill"><el-icon><Clock /></el-icon><span>{{ format(item.at) }}</span><b>{{ item.company }}</b></router-link></div>
    <router-link class="reminder-open" to="/calendar">{{ language === 'en' ? 'Open calendar →' : '打开日历 →' }}</router-link>
  </section>
</template>
