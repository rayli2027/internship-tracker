<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ArrowLeft, ArrowRight, Calendar, Clock, Link, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import AppHeader from '../components/AppHeader.vue'
import ApplicationDrawer from '../components/ApplicationDrawer.vue'
import { api } from '../api'
import type { Application, ApplicationEvent } from '../types'
import { useI18n } from '../i18n'
import { applicationToForm, toPayload } from '../utils'
import { localDateKey } from '../planner'

type CalendarItem = {
  id: string
  applicationId: string
  type: 'deadline' | 'interview' | 'workflow'
  eventType?: string
  at: string
  company: string
  position: string
  note: string
  sourceEvent?: ApplicationEvent
}

const applications = ref<Application[]>([])
const events = ref<ApplicationEvent[]>([])
const loading = ref(false)
const cursor = ref(new Date(new Date().getFullYear(), new Date().getMonth(), 1))
const selectedKey = ref('')
const drawerOpen = ref(false)
const current = ref<Application | null>(null)
const viewMode = ref<'month' | 'week' | 'agenda'>('month')
const filters = ref(['deadline', 'interview', 'workflow'])
const eventDialog = ref(false)
const eventSaving = ref(false)
const monthMenuOpen = ref(false)
const pickerYear = ref(cursor.value.getFullYear())
const eventForm = ref({ application_id: '', event_type: 'interview', title: '', content: '', event_at: '' })
const { language } = useI18n()
const en = computed(() => language.value === 'en')

const keyOf = (date: Date) => localDateKey(date)
const parseKey = (value: string) => {
  const d = new Date(value)
  return Number.isNaN(d.getTime()) ? '' : keyOf(d)
}
const isRealDate = (value: string) => {
  const d = new Date(value)
  return !Number.isNaN(d.getTime()) && d.getFullYear() > 1971
}
const todayKey = keyOf(new Date())

const allItems = computed<CalendarItem[]>(() =>
  [
    ...applications.value
      .filter((i) => i.deadline_at)
      .map((i) => ({
        id: `deadline-${i.id}`,
        applicationId: i.id,
        type: 'deadline' as const,
        at: i.deadline_at,
        company: i.company_name,
        position: i.position_name,
        note: i.notes,
      })),
    ...applications.value
      .filter(
        (i) =>
          isRealDate(i.applied_at) &&
          !events.value.some(
            (e) =>
              e.application_id === i.id &&
              (e.event_type === 'applied' || (e.event_type === 'status_change' && e.to_status === 'Applied')) &&
              parseKey(e.event_at) === parseKey(i.applied_at),
          ),
      )
      .map((i) => ({
        id: `applied-${i.id}`,
        applicationId: i.id,
        type: 'workflow' as const,
        eventType: 'applied',
        at: i.applied_at,
        company: i.company_name,
        position: i.position_name,
        note: i.notes,
      })),
    ...events.value
      .filter((i) => i.event_type !== 'created')
      .map((i) => ({
        id: `event-${i.id}`,
        applicationId: i.application_id,
        type: (i.event_type === 'interview' ? 'interview' : 'workflow') as 'interview' | 'workflow',
        eventType: i.event_type === 'status_change' ? i.to_status || i.event_type : i.event_type,
        at: i.event_at,
        company: i.company_name || '',
        position: i.position_name || '',
        note: [i.title, i.content, ...Object.values(i.details || {}).map(String)].filter(Boolean).join(' · '),
        sourceEvent: i,
      })),
  ].sort((a, b) => new Date(a.at).getTime() - new Date(b.at).getTime()),
)

const filteredItems = computed(() => allItems.value.filter((i) => filters.value.includes(i.type)))
const itemsByDay = computed(() => {
  const map = new Map<string, CalendarItem[]>()
  for (const item of filteredItems.value) {
    const key = parseKey(item.at)
    if (!map.has(key)) map.set(key, [])
    map.get(key)!.push(item)
  }
  return map
})
const days = computed(() => {
  const first = new Date(cursor.value.getFullYear(), cursor.value.getMonth(), 1)
  const start = new Date(first)
  start.setDate(1 - ((first.getDay() + 6) % 7))
  return Array.from({ length: 42 }, (_, index) => {
    const date = new Date(start)
    date.setDate(start.getDate() + index)
    const key = keyOf(date)
    return {
      key,
      date,
      current: date.getMonth() === cursor.value.getMonth(),
      today: key === todayKey,
      items: itemsByDay.value.get(key) || [],
    }
  })
})
const selectedItems = computed(() => itemsByDay.value.get(selectedKey.value) || [])
const shownDays = computed(() =>
  viewMode.value === 'month'
    ? days.value
    : days.value.filter((d) => {
        if (viewMode.value === 'agenda') return d.key === selectedKey.value
        const selected = new Date(`${selectedKey.value}T12:00:00`)
        const start = new Date(selected)
        start.setDate(selected.getDate() - ((selected.getDay() + 6) % 7))
        const end = new Date(start)
        end.setDate(start.getDate() + 6)
        return d.date >= start && d.date <= end
      }),
)
const locale = computed(() => (language.value === 'en' ? 'en-US' : 'zh-CN'))
const monthTitle = computed(() => new Intl.DateTimeFormat(locale.value, { month: 'long', year: 'numeric' }).format(cursor.value))
const monthNames = computed(() =>
  Array.from({ length: 12 }, (_, month) => new Intl.DateTimeFormat(locale.value, { month: 'short' }).format(new Date(2026, month, 1))),
)
const selectedTitle = computed(() =>
  selectedKey.value
    ? new Intl.DateTimeFormat(locale.value, { weekday: 'long', month: 'long', day: 'numeric' }).format(new Date(`${selectedKey.value}T12:00:00`))
    : '',
)
const weekDays = computed(() => (en.value ? ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] : ['一', '二', '三', '四', '五', '六', '日']))
const time = (value: string) => new Intl.DateTimeFormat(locale.value, { hour: '2-digit', minute: '2-digit' }).format(new Date(value))

const eventLabel = (item: CalendarItem) => {
  if (item.type === 'deadline') return en.value ? 'Deadline' : '截止'
  if (item.type === 'interview') return en.value ? 'Interview' : '面试'
  const labels: Record<string, [string, string]> = {
    oa: ['Assessment', '笔试'],
    applied: ['Apply', '投递'],
    offer: ['Offer', 'Offer'],
    memo: ['Reminder', '提醒'],
    rejected: ['Outcome', '结果'],
    withdrawn: ['Withdrawn', '已放弃'],
    Applied: ['Apply', '投递'],
    OA: ['Assessment', '笔试'],
    Interview: ['Interview', '面试'],
    Offer: ['Offer', 'Offer'],
    Rejected: ['Outcome', '结果'],
    Withdrawn: ['Withdrawn', '已放弃'],
  }
  const pair = labels[item.eventType || '']
  return pair ? pair[en.value ? 0 : 1] : en.value ? 'Workflow' : '流程'
}

const typeOptions = computed(() => [
  { value: 'interview', label: en.value ? 'Interview' : '面试' },
  { value: 'oa', label: en.value ? 'Assessment' : '笔试' },
  { value: 'memo', label: en.value ? 'Reminder' : '提醒' },
  { value: 'offer', label: 'Offer' },
])

async function load() {
  loading.value = true
  try {
    const [a, e] = await Promise.all([api.listApplications(), api.listTimelineEvents()])
    applications.value = a.items
    events.value = e.items
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : en.value ? 'Unable to load calendar' : '无法加载日历')
  } finally {
    loading.value = false
  }
}

function moveMonth(delta: number) {
  cursor.value = new Date(cursor.value.getFullYear(), cursor.value.getMonth() + delta, 1)
}
function openMonthMenu() {
  pickerYear.value = cursor.value.getFullYear()
}
function selectMonth(month: number) {
  cursor.value = new Date(pickerYear.value, month, 1)
  selectedKey.value = keyOf(new Date(pickerYear.value, month, 1))
  monthMenuOpen.value = false
}
function goToday() {
  const now = new Date()
  cursor.value = new Date(now.getFullYear(), now.getMonth(), 1)
  selectedKey.value = todayKey
}
function openApplication(id: string) {
  current.value = applications.value.find((item) => item.id === id) || null
  if (current.value) drawerOpen.value = true
}
function openCreate() {
  eventForm.value = {
    application_id: applications.value[0]?.id || '',
    event_type: 'interview',
    title: '',
    content: '',
    event_at: `${selectedKey.value}T09:00`,
  }
  eventDialog.value = true
}
async function createEvent() {
  if (!eventForm.value.application_id || !eventForm.value.title) return
  eventSaving.value = true
  try {
    await api.createApplicationEvent(eventForm.value.application_id, {
      event_type: eventForm.value.event_type as 'interview',
      title: eventForm.value.title,
      content: eventForm.value.content,
      event_at: eventForm.value.event_at,
      details: {},
    })
    eventDialog.value = false
    await load()
    ElMessage.success(en.value ? 'Schedule created' : '日程已创建')
  } finally {
    eventSaving.value = false
  }
}
function dragStart(item: CalendarItem, event: DragEvent) {
  event.stopPropagation()
  event.dataTransfer?.setData('text/calendar-item', item.id)
}
async function dropOnDay(key: string, event: DragEvent) {
  event.preventDefault()
  const id = event.dataTransfer?.getData('text/calendar-item')
  const item = allItems.value.find((i) => i.id === id)
  if (!item) return
  const old = new Date(item.at)
  const next = `${key}T${String(old.getHours()).padStart(2, '0')}:${String(old.getMinutes()).padStart(2, '0')}`
  if (item.sourceEvent) {
    const e = item.sourceEvent
    await api.updateApplicationEvent(item.applicationId, e.id, {
      event_type: (e.event_type === 'status_change' ? 'memo' : e.event_type) as 'memo',
      title: e.title || eventLabel(item),
      content: e.content,
      event_at: next,
      details: e.details,
    })
  } else {
    const app = applications.value.find((i) => i.id === item.applicationId)
    if (!app) return
    const form = applicationToForm(app)
    if (item.type === 'deadline') form.deadline_at = key
    else form.applied_at = next
    await api.updateApplication(app.id, toPayload(form))
  }
  await load()
  ElMessage.success(en.value ? 'Schedule moved' : '日程已移动')
}

onMounted(() => {
  selectedKey.value = todayKey
  load()
})
</script>

<template>
  <div class="app-shell">
    <AppHeader />
    <main class="workspace calendar-page" v-loading="loading">
      <section class="board-toolbar-strip calendar-toolbar-strip">
        <div class="view-switch" role="group">
          <button type="button" :class="{ active: viewMode === 'month' }" @click="viewMode = 'month'">{{ en ? 'Month' : '月' }}</button>
          <button type="button" :class="{ active: viewMode === 'week' }" @click="viewMode = 'week'">{{ en ? 'Week' : '周' }}</button>
          <button type="button" :class="{ active: viewMode === 'agenda' }" @click="viewMode = 'agenda'">{{ en ? 'Agenda' : '日' }}</button>
        </div>
        <el-checkbox-group v-model="filters" size="small" class="calendar-filters">
          <el-checkbox-button value="deadline">{{ en ? 'Deadlines' : '截止' }}</el-checkbox-button>
          <el-checkbox-button value="interview">{{ en ? 'Interviews' : '面试' }}</el-checkbox-button>
          <el-checkbox-button value="workflow">{{ en ? 'Workflow' : '流程' }}</el-checkbox-button>
        </el-checkbox-group>
        <div class="toolbar-spacer"></div>
        <button class="toolbar-add" type="button" @click="openCreate">
          <el-icon><Plus /></el-icon>
          <span>{{ en ? 'Add schedule' : '新增日程' }}</span>
        </button>
      </section>

      <section class="calendar-layout" :class="`view-${viewMode}`">
        <div class="page-panel calendar-card">
          <header class="calendar-toolbar">
            <span></span>
            <div class="month-navigation">
              <el-button circle :icon="ArrowLeft" :aria-label="en ? 'Previous month' : '上个月'" @click="moveMonth(-1)" />
              <el-popover v-model:visible="monthMenuOpen" placement="bottom" :width="286" trigger="click" @show="openMonthMenu">
                <template #reference>
                  <button class="month-title-button" type="button">{{ monthTitle }}</button>
                </template>
                <div class="month-jump">
                  <div class="month-jump-year">
                    <el-button text :icon="ArrowLeft" @click="pickerYear--" />
                    <strong>{{ pickerYear }}</strong>
                    <el-button text :icon="ArrowRight" @click="pickerYear++" />
                  </div>
                  <div class="month-jump-grid">
                    <button
                      v-for="(name, index) in monthNames"
                      :key="name"
                      type="button"
                      :class="{ active: pickerYear === cursor.getFullYear() && index === cursor.getMonth() }"
                      @click="selectMonth(index)"
                    >
                      {{ name }}
                    </button>
                  </div>
                </div>
              </el-popover>
              <el-button circle :icon="ArrowRight" :aria-label="en ? 'Next month' : '下个月'" @click="moveMonth(1)" />
            </div>
            <el-button @click="goToday">{{ en ? 'Today' : '今天' }}</el-button>
          </header>
          <div v-if="viewMode !== 'agenda'" class="calendar-week">
            <span v-for="day in weekDays" :key="day">{{ day }}</span>
          </div>
          <div class="calendar-grid">
            <button
              v-for="day in shownDays"
              :key="day.key"
              type="button"
              class="calendar-day"
              :class="{ muted: !day.current, today: day.today, selected: selectedKey === day.key }"
              @click="selectedKey = day.key"
              @dragover.prevent
              @drop="dropOnDay(day.key, $event)"
            >
              <span class="day-number">{{ day.date.getDate() }}</span>
              <span
                v-for="item in day.items.slice(0, 3)"
                :key="item.id"
                :class="`calendar-event ${item.type}`"
                draggable="true"
                @dragstart="dragStart(item, $event)"
              >
                <i></i>
                <b>{{ item.company }}</b>
                <small>{{ item.type === 'interview' ? time(item.at) : eventLabel(item) }}</small>
              </span>
              <span v-if="day.items.length > 3" class="more-events">+{{ day.items.length - 3 }} {{ en ? 'more' : '项' }}</span>
            </button>
          </div>
        </div>

        <aside class="page-panel day-agenda">
          <div class="agenda-heading">
            <span><el-icon><Calendar /></el-icon></span>
            <div>
              <p>{{ en ? 'Day plan' : '当日日程' }}</p>
              <h2>{{ selectedTitle }}</h2>
            </div>
          </div>
          <div v-if="!selectedItems.length" class="agenda-empty">
            <strong>{{ en ? 'No events' : '当天暂无安排' }}</strong>
            <p>{{ en ? 'Pick another date or add a schedule.' : '选择其他日期，或新增日程。' }}</p>
          </div>
          <button v-for="item in selectedItems" :key="item.id" class="agenda-item" type="button" @click="openApplication(item.applicationId)">
            <span :class="`agenda-type ${item.type}`">{{ eventLabel(item) }}</span>
            <h3>{{ item.company }}</h3>
            <p>{{ item.position }}</p>
            <span class="agenda-time"><el-icon><Clock /></el-icon>{{ time(item.at) }}</span>
            <small v-if="item.note">{{ item.note }}</small>
            <span class="agenda-open"><el-icon><Link /></el-icon>{{ en ? 'Open' : '查看岗位' }}</span>
          </button>
        </aside>
      </section>
    </main>

    <ApplicationDrawer v-model="drawerOpen" :item="current" @saved="load" @deleted="load" />
    <el-dialog v-model="eventDialog" :title="en ? 'Add schedule' : '新增日程'" width="min(480px, 92vw)">
      <el-form label-position="top">
        <el-form-item :label="en ? 'Application' : '岗位'">
          <el-select v-model="eventForm.application_id" style="width: 100%">
            <el-option v-for="item in applications" :key="item.id" :value="item.id" :label="`${item.company_name} · ${item.position_name}`" />
          </el-select>
        </el-form-item>
        <el-form-item :label="en ? 'Type' : '类型'">
          <el-select v-model="eventForm.event_type" style="width: 100%">
            <el-option v-for="opt in typeOptions" :key="opt.value" :value="opt.value" :label="opt.label" />
          </el-select>
        </el-form-item>
        <el-form-item :label="en ? 'Date and time' : '日期与时间'">
          <el-date-picker v-model="eventForm.event_at" type="datetime" value-format="YYYY-MM-DDTHH:mm" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="en ? 'Title' : '标题'"><el-input v-model="eventForm.title" /></el-form-item>
        <el-form-item :label="en ? 'Notes' : '备注'"><el-input v-model="eventForm.content" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="eventDialog = false">{{ en ? 'Cancel' : '取消' }}</el-button>
        <el-button type="primary" :loading="eventSaving" @click="createEvent">{{ en ? 'Save' : '保存' }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>
