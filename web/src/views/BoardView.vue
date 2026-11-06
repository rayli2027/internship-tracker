<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { Grid, List, Plus, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import AppHeader from '../components/AppHeader.vue'
import ApplicationCard from '../components/ApplicationCard.vue'
import ApplicationListRow from '../components/ApplicationListRow.vue'
import ApplicationDrawer from '../components/ApplicationDrawer.vue'
import { api } from '../api'
import { STATUS_META } from '../constants'
import type { Application, ApplicationStatus } from '../types'
import { useI18n } from '../i18n'

type DisplayStage = Exclude<ApplicationStatus, 'Rejected' | 'Withdrawn'> | 'Final'
const BOARD_STAGES: DisplayStage[] = ['Wishlist', 'Applied', 'OA', 'Interview', 'Offer', 'Final']

const items = ref<Application[]>([])
const keyword = ref('')
const displayMode = ref<'board' | 'list'>(localStorage.getItem('it_display_mode') === 'list' ? 'list' : 'board')
const phaseFilter = ref<DisplayStage | 'All'>('All')
const listPage = ref(1)
const listPageSize = ref(Number(localStorage.getItem('it_list_page_size') || 20))
const loading = ref(false)
const drawerOpen = ref(false)
const drawerMode = ref<'create' | 'edit' | 'detail'>('create')
const current = ref<Application | null>(null)
const draggingId = ref<string | null>(null)
const dragOverStatus = ref<DisplayStage | null>(null)
const statusDialogOpen = ref(false)
const pendingStatus = ref<ApplicationStatus>('Wishlist')
const pendingItem = ref<Application | null>(null)
const statusMemo = ref('')
const statusWhen = ref('')
const { language, statusLabel, t } = useI18n()

const stageLabel = (stage: DisplayStage) =>
  stage === 'Final'
    ? language.value === 'en'
      ? 'Closed Without Offer'
      : '未获 Offer 结束'
    : language.value === 'en'
      ? statusLabel(stage)
      : STATUS_META[stage].label

const stageHint = (stage: DisplayStage) =>
  stage === 'Final'
    ? language.value === 'en'
      ? 'Rejected or withdrawn'
      : '被拒或主动放弃'
    : STATUS_META[stage].hint

const phaseOf = (item: Application): DisplayStage =>
  ['Rejected', 'Withdrawn'].includes(item.status) ? 'Final' : (item.status as DisplayStage)

const sortedItems = computed(() =>
  [...items.value].sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()),
)

const phaseCounts = computed(
  () =>
    Object.fromEntries(
      BOARD_STAGES.map((stage) => [stage, sortedItems.value.filter((item) => phaseOf(item) === stage).length]),
    ) as Record<DisplayStage, number>,
)

const listItems = computed(() =>
  phaseFilter.value === 'All'
    ? sortedItems.value
    : sortedItems.value.filter((item) => phaseOf(item) === phaseFilter.value),
)

const pagedListItems = computed(() => {
  const start = (listPage.value - 1) * listPageSize.value
  return listItems.value.slice(start, start + listPageSize.value)
})

let timer = 0
let suppressClick = false

const grouped = computed(() => {
  const map = Object.fromEntries(BOARD_STAGES.map((status) => [status, [] as Application[]])) as Record<DisplayStage, Application[]>
  for (const item of sortedItems.value) {
    const stage: DisplayStage = ['Rejected', 'Withdrawn'].includes(item.status) ? 'Final' : (item.status as DisplayStage)
    map[stage].push(item)
  }
  return map
})

async function load() {
  loading.value = true
  try {
    const data = await api.listApplications(keyword.value.trim(), false)
    items.value = data.items
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '加载失败')
  } finally {
    loading.value = false
  }
}

function scheduleLoad() {
  window.clearTimeout(timer)
  timer = window.setTimeout(load, 280)
}

function openCreate() {
  current.value = null
  drawerMode.value = 'create'
  drawerOpen.value = true
}

function openDetail(item: Application) {
  if (suppressClick) return
  current.value = item
  drawerMode.value = 'detail'
  drawerOpen.value = true
}

function openEdit(item: Application) {
  current.value = item
  drawerMode.value = 'edit'
  drawerOpen.value = true
}

function onDragStart(item: Application, event: DragEvent) {
  draggingId.value = item.id
  event.dataTransfer?.setData('text/plain', item.id)
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
  }
}

function onDragEnd() {
  if (draggingId.value) {
    suppressClick = true
    window.setTimeout(() => {
      suppressClick = false
    }, 80)
  }
  draggingId.value = null
  dragOverStatus.value = null
}

function onDragOver(status: DisplayStage, event: DragEvent) {
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'move'
  }
  dragOverStatus.value = status
}

async function onDrop(status: DisplayStage, event: DragEvent) {
  event.preventDefault()
  const id = event.dataTransfer?.getData('text/plain') || draggingId.value
  const item = items.value.find((row) => row.id === id)
  dragOverStatus.value = null
  draggingId.value = null
  if (!item || item.status === status || (status === 'Final' && ['Rejected', 'Withdrawn'].includes(item.status))) {
    return
  }
  await changeStatus(item, status === 'Final' ? 'Rejected' : status)
}

async function changeStatus(item: Application, status: ApplicationStatus) {
  if (item.status === status) return
  pendingItem.value = item
  pendingStatus.value = status
  statusMemo.value = ''
  statusWhen.value = ''
  statusDialogOpen.value = true
}

async function applyStatusChange() {
  const item = pendingItem.value
  if (!item) return
  const status = pendingStatus.value
  try {
    await api.changeApplicationStatus(item.id, {
      status,
      content: statusMemo.value,
      event_at: statusWhen.value,
    })
    ElMessage.success(language.value === 'en' ? `Moved to ${statusLabel(status)}` : `已移到 ${statusLabel(status)}`)
    statusDialogOpen.value = false
    await load()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '改状态失败')
  }
}

async function toggleLike(item: Application) {
  try {
    await api.setApplicationLiked(item.id, !item.is_liked)
    await load()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '更新喜爱失败')
  }
}

function onKeydown(event: KeyboardEvent) {
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'k') {
    event.preventDefault()
    document.querySelector<HTMLInputElement>('.board-search input')?.focus()
  }
}

onMounted(() => {
  load()
  window.addEventListener('keydown', onKeydown)
})

watch(displayMode, (value) => localStorage.setItem('it_display_mode', value))
watch(listPageSize, (value) => localStorage.setItem('it_list_page_size', String(value)))
watch([phaseFilter, listPageSize], () => {
  listPage.value = 1
})
watch(listItems, (rows) => {
  const maxPage = Math.max(1, Math.ceil(rows.length / listPageSize.value))
  if (listPage.value > maxPage) listPage.value = maxPage
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
  window.clearTimeout(timer)
})
</script>

<template>
  <div class="app-shell">
    <AppHeader />
    <main class="workspace">
      <section class="board-toolbar-strip">
        <div class="view-switch" role="group" :aria-label="language === 'en' ? 'View mode' : '视图模式'">
          <button type="button" :class="{ active: displayMode === 'board' }" :title="language === 'en' ? 'Board' : '看板'" @click="displayMode = 'board'">
            <el-icon><Grid /></el-icon>
          </button>
          <button type="button" :class="{ active: displayMode === 'list' }" :title="language === 'en' ? 'List' : '列表'" @click="displayMode = 'list'">
            <el-icon><List /></el-icon>
          </button>
        </div>
        <el-input
          v-model="keyword"
          class="board-search"
          clearable
          :placeholder="t('search')"
          :prefix-icon="Search"
          @input="scheduleLoad"
          @clear="load"
        />
        <button class="toolbar-add" type="button" @click="openCreate">
          <el-icon><Plus /></el-icon>
          <span>{{ t('newApplication') }}</span>
        </button>
      </section>

      <section v-if="displayMode === 'board'" class="board" v-loading="loading">
        <article
          v-for="status in BOARD_STAGES"
          :key="status"
          class="lane"
          :class="{ 'is-drop-target': dragOverStatus === status }"
          @dragover="onDragOver(status, $event)"
          @drop="onDrop(status, $event)"
        >
          <header>
            <div>
              <div class="lane-title">{{ stageLabel(status) }}</div>
              <div v-if="language !== 'en'" class="lane-hint">{{ status }}</div>
            </div>
            <el-badge :value="grouped[status].length" type="info" />
          </header>
          <ul class="lane-list">
            <li v-for="item in grouped[status]" :key="item.id">
              <ApplicationCard
                :item="item"
                :dragging="draggingId === item.id"
                @open="openDetail(item)"
                @edit="openEdit(item)"
                @toggle-like="toggleLike(item)"
                @drag-start="(event) => onDragStart(item, event)"
                @drag-end="onDragEnd"
              />
            </li>
            <li v-if="grouped[status].length === 0" class="empty-card">
              <strong>{{ dragOverStatus === status ? '松开放到这一列' : keyword ? '没有匹配的卡片' : '这一列还是空的' }}</strong>
              <span>{{ dragOverStatus === status ? stageHint(status) : keyword ? '换个公司名试试' : '拖卡片进来，或点右上角新建' }}</span>
            </li>
          </ul>
        </article>
      </section>

      <section v-else class="list-view" v-loading="loading">
        <nav class="phase-tabs">
          <button type="button" :class="{ active: phaseFilter === 'All' }" @click="phaseFilter = 'All'">
            {{ language === 'en' ? 'All' : '全部' }} <span>{{ items.length }}</span>
          </button>
          <button
            v-for="stage in BOARD_STAGES"
            :key="stage"
            type="button"
            :class="{ active: phaseFilter === stage }"
            @click="phaseFilter = stage"
          >
            {{ stageLabel(stage) }} <span>{{ phaseCounts[stage] }}</span>
          </button>
        </nav>

        <article class="list-card">
          <div class="application-list-head">
            <span>{{ language === 'en' ? '#' : '序号' }}</span>
            <span>{{ language === 'en' ? 'Company' : '公司' }}</span>
            <span>{{ language === 'en' ? 'Position' : '职位' }}</span>
            <span>{{ language === 'en' ? 'Stage' : '阶段' }}</span>
            <span>{{ language === 'en' ? 'Location' : '地点' }}</span>
            <span>{{ language === 'en' ? 'Pay' : '薪酬' }}</span>
            <span>{{ language === 'en' ? 'Deadline' : '截止日期' }}</span>
            <span>{{ language === 'en' ? 'Created' : '新增时间' }}</span>
            <span>{{ language === 'en' ? 'Actions' : '操作' }}</span>
          </div>
          <ul v-if="pagedListItems.length" class="application-list">
            <ApplicationListRow
              v-for="(item, idx) in pagedListItems"
              :key="item.id"
              :item="item"
              :index="(listPage - 1) * listPageSize + idx + 1"
              @open="openDetail(item)"
              @edit="openEdit(item)"
              @toggle-like="toggleLike(item)"
            />
          </ul>
          <div v-else class="list-empty">
            {{ language === 'en' ? 'No applications match your filters.' : '没有符合条件的岗位。' }}
          </div>

          <div class="list-pagination">
            <el-pagination
              v-model:current-page="listPage"
              v-model:page-size="listPageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="listItems.length"
              background
              layout="total, sizes, prev, pager, next"
            />
          </div>
        </article>
      </section>
    </main>

    <ApplicationDrawer v-model="drawerOpen" :item="current" :mode="drawerMode" @saved="load" @deleted="load" />
    <el-dialog
      v-model="statusDialogOpen"
      :title="language === 'en' ? 'Move stage' : '阶段流转'"
      width="min(420px, 92vw)"
    >
      <p class="subtitle">
        {{ pendingItem?.company_name }} / {{ pendingItem?.position_name }}：
        {{ pendingItem ? statusLabel(pendingItem.status) : '' }} → {{ statusLabel(pendingStatus) }}
      </p>
      <el-form label-position="top">
        <el-form-item v-if="pendingStatus === 'Rejected' || pendingStatus === 'Withdrawn'" :label="language === 'en' ? 'Outcome' : '结果'">
          <el-radio-group v-model="pendingStatus">
            <el-radio-button value="Rejected">{{ language === 'en' ? 'Rejected' : '被拒' }}</el-radio-button>
            <el-radio-button value="Withdrawn">{{ language === 'en' ? 'Withdrawn' : '主动放弃' }}</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="language === 'en' ? 'When' : '时间'">
          <el-date-picker
            v-model="statusWhen"
            type="datetime"
            value-format="YYYY-MM-DDTHH:mm"
            :placeholder="language === 'en' ? 'Defaults to now' : '默认当前时间'"
            style="width:100%"
          />
        </el-form-item>
        <el-form-item :label="language === 'en' ? 'Notes' : '备注'">
          <el-input
            v-model="statusMemo"
            type="textarea"
            :rows="3"
            :placeholder="language === 'en' ? 'Optional' : '可选'"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="statusDialogOpen = false">{{ language === 'en' ? 'Cancel' : '取消' }}</el-button>
        <el-button type="primary" @click="applyStatusChange">{{ language === 'en' ? 'Confirm' : '确认' }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>
