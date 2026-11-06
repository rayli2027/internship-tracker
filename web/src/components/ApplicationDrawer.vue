<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Application, ApplicationEvent, ApplicationEventPayload, ApplicationStatus, EventType } from '../types'
import { EDUCATIONS, EMPLOYMENT_TYPES, PRIORITIES, SALARY_UNITS, SOURCE_TYPES, STATUSES, STATUS_META, WORK_MODES } from '../constants'
import { applicationToForm, emptyApplicationForm, fmtTime, isValidEmail, toPayload } from '../utils'
import { api } from '../api'
import { useI18n } from '../i18n'

const props = withDefaults(
  defineProps<{
    modelValue: boolean
    item: Application | null
    mode?: 'create' | 'edit' | 'detail'
  }>(),
  { mode: 'detail' },
)
const { language, statusLabel } = useI18n()
const emit = defineEmits<{ 'update:modelValue': [value: boolean]; saved: []; deleted: [] }>()
const form = reactive(emptyApplicationForm())
const saving = ref(false)
const deleting = ref(false)
const events = ref<ApplicationEvent[]>([])
const eventsLoading = ref(false)
const detailTab = ref<'overview' | 'lifecycle'>('overview')

const title = computed(() => {
  if (props.mode === 'create') return language.value === 'en' ? 'New application' : '新建岗位'
  if (props.mode === 'edit') return language.value === 'en' ? 'Edit application' : '修改岗位'
  if (props.item) return `${props.item.company_name} · ${props.item.position_name}`
  return language.value === 'en' ? 'Application' : '岗位详情'
})
const drawerSize = computed(() =>
  props.mode === 'detail'
    ? (language.value === 'en' ? 'min(720px, 96vw)' : 'min(640px, 94vw)')
    : (language.value === 'en' ? 'min(560px, 94vw)' : 'min(520px, 92vw)'),
)

watch(
  () => [props.modelValue, props.item, props.mode] as const,
  ([open, item, mode]) => {
    if (!open) return
    if (mode === 'create') {
      Object.assign(form, emptyApplicationForm())
      events.value = []
      return
    }
    if (item && mode === 'edit') {
      Object.assign(form, applicationToForm(item))
      return
    }
    if (item && mode === 'detail') {
      detailTab.value = 'overview'
      loadEvents()
      return
    }
  },
  { immediate: true },
)

async function loadEvents() {
  if (!props.item) return
  eventsLoading.value = true
  try {
    events.value = (await api.listApplicationEvents(props.item.id)).items
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '加载生命周期失败')
  } finally {
    eventsLoading.value = false
  }
}

function close() {
  emit('update:modelValue', false)
}

function validURL(value: string) {
  if (!value.trim()) return true
  try {
    return ['http:', 'https:'].includes(new URL(value).protocol)
  } catch {
    return false
  }
}

async function save() {
  if (!form.company_name.trim() || !form.position_name.trim()) return ElMessage.error('公司名称和岗位名称为必填项')
  if (!validURL(form.apply_link) || !validURL(form.source_url)) return ElMessage.error('链接需以 http:// 或 https:// 开头')
  if (form.contact_email && !isValidEmail(form.contact_email)) return ElMessage.error('请输入有效的联系人邮箱')
  if (saving.value) return
  saving.value = true
  try {
    if (props.item) {
      const nextStatus = form.status
      const payload = toPayload(form)
      payload.status = props.item.status
      await api.updateApplication(props.item.id, payload)
      if (nextStatus !== props.item.status) await api.changeApplicationStatus(props.item.id, { status: nextStatus })
    } else {
      await api.createApplication(toPayload(form))
    }
    ElMessage.success(props.item ? '已保存更改' : '已新建岗位')
    emit('saved')
    close()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function remove() {
  if (!props.item || deleting.value) return
  try {
    await ElMessageBox.confirm(
      `确认删除“${props.item.company_name} / ${props.item.position_name}”？数据会安全保留，但不再显示在看板中。`,
      '删除岗位',
      { confirmButtonText: '确认删除', cancelButtonText: '取消', type: 'warning' },
    )
    deleting.value = true
    await api.deleteApplication(props.item.id)
    ElMessage.success('岗位已从看板移除')
    emit('deleted')
    close()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') ElMessage.error(error instanceof Error ? error.message : '删除失败')
  } finally {
    deleting.value = false
  }
}
</script>

<template>
  <el-drawer :model-value="modelValue" :title="title" :size="drawerSize" destroy-on-close @close="close">
    <section v-if="mode === 'detail' && item" class="detail-drawer">
      <el-tabs v-model="detailTab" class="detail-tabs">
        <el-tab-pane :label="language === 'en' ? 'Overview' : '基本信息'" name="overview">
          <section class="application-overview" v-loading="eventsLoading">
            <div class="overview-status">
              <span>{{ language === 'en' ? 'Current stage' : '当前阶段' }}</span>
              <strong>{{ statusLabel(item.status) }}</strong>
            </div>
            <div class="overview-grid overview-grid--stack">
          <div v-if="item.location">
            <span>{{ language === 'en' ? 'Location' : '地点' }}</span>
            <strong>{{ item.location }}</strong>
          </div>
          <div v-if="item.work_mode">
            <span>{{ language === 'en' ? 'Work mode' : '办公方式' }}</span>
            <strong>{{ item.work_mode }}</strong>
          </div>
          <div v-if="item.department">
            <span>{{ language === 'en' ? 'Department / team' : '部门 / 团队' }}</span>
            <strong>{{ item.department }}</strong>
          </div>
          <div v-if="item.employment_type">
            <span>{{ language === 'en' ? 'Employment type' : '用工类型' }}</span>
            <strong>{{ item.employment_type }}</strong>
          </div>
          <div v-if="item.priority">
            <span>{{ language === 'en' ? 'Priority' : '优先级' }}</span>
            <strong>{{ item.priority }}</strong>
          </div>
          <div v-if="item.education">
            <span>{{ language === 'en' ? 'Education' : '学历要求' }}</span>
            <strong>{{ item.education }}</strong>
          </div>
          <div v-if="item.experience_years">
            <span>{{ language === 'en' ? 'Experience' : '经验要求' }}</span>
            <strong>{{ item.experience_years }} {{ language === 'en' ? 'years' : '年' }}</strong>
          </div>
          <div v-if="item.days_per_week">
            <span>{{ language === 'en' ? 'Schedule' : '到岗安排' }}</span>
            <strong>{{ item.days_per_week }} {{ language === 'en' ? 'days/week' : '天/周' }} · {{ item.intern_months }} {{ language === 'en' ? 'months' : '个月' }}</strong>
          </div>
          <div v-if="fmtTime(item.deadline_at)">
            <span>{{ language === 'en' ? 'Application deadline' : '截止日期' }}</span>
            <strong>{{ fmtTime(item.deadline_at) }}</strong>
          </div>
          <div v-if="fmtTime(item.applied_at)">
            <span>{{ language === 'en' ? 'Applied on' : '投递时间' }}</span>
            <strong>{{ fmtTime(item.applied_at) }}</strong>
          </div>
          <div v-if="item.salary_min || item.salary_max">
            <span>{{ language === 'en' ? 'Compensation' : '薪酬' }}</span>
            <strong>{{ item.salary_min }}–{{ item.salary_max }} {{ item.salary_unit }}</strong>
          </div>
          <div v-if="item.resume_version">
            <span>{{ language === 'en' ? 'Resume used' : '使用简历' }}</span>
            <strong>{{ item.resume_version }}</strong>
          </div>
          <div v-if="item.source_name || item.source_type">
            <span>{{ language === 'en' ? 'Source' : '来源' }}</span>
            <strong>{{ item.source_name || item.source_type }}</strong>
          </div>
          <div v-if="item.contact_name || item.contact_email">
            <span>{{ language === 'en' ? 'Contact' : '联系人' }}</span>
            <strong>
              {{ item.contact_name }}
              <template v-if="item.contact_name && item.contact_email"> · </template>
              {{ item.contact_email }}
            </strong>
          </div>
        </div>
        <div v-if="item.tech_stack" class="overview-notes">
          <span>{{ language === 'en' ? 'Skills / tech stack' : '技能 / 技术栈' }}</span>
          <p>{{ item.tech_stack }}</p>
        </div>
        <div v-if="item.requirements" class="overview-notes">
          <span>{{ language === 'en' ? 'Requirements' : '岗位要求' }}</span>
          <p>{{ item.requirements }}</p>
        </div>
        <div v-if="item.responsibilities" class="overview-notes">
          <span>{{ language === 'en' ? 'Responsibilities' : '岗位职责' }}</span>
          <p>{{ item.responsibilities }}</p>
        </div>
        <div v-if="item.notes" class="overview-notes">
          <span>{{ language === 'en' ? 'Personal notes' : '个人笔记' }}</span>
          <p>{{ item.notes }}</p>
        </div>
            <div v-if="item.apply_link || item.source_url" class="overview-actions">
              <el-button tag="a" :href="item.apply_link || item.source_url" target="_blank">
                {{ language === 'en' ? 'Open posting' : '打开岗位链接' }}
              </el-button>
            </div>
          </section>
        </el-tab-pane>
        <el-tab-pane :label="language === 'en' ? 'Lifecycle' : '生命周期'" name="lifecycle">
          <section class="detail-lifecycle-panel" v-loading="eventsLoading">
            <el-timeline v-if="events.length" class="event-timeline">
              <el-timeline-item v-for="event in events" :key="event.id" :timestamp="fmtTime(event.event_at)" placement="top">
                <article class="event-card event-card--readonly">
                  <div>
                    <strong>{{ event.title }}</strong>
                    <span v-if="event.from_status || event.to_status" class="event-status">{{ event.from_status }} → {{ event.to_status }}</span>
                  </div>
                  <p v-if="event.content">{{ event.content }}</p>
                  <div class="event-details">
                    <span v-for="(value, key) in event.details" v-show="value" :key="key">{{ key }}: {{ value }}</span>
                  </div>
                </article>
              </el-timeline-item>
            </el-timeline>
            <p v-else class="detail-lifecycle-empty">{{ language === 'en' ? 'No lifecycle events yet.' : '暂无生命周期记录。' }}</p>
          </section>
        </el-tab-pane>
      </el-tabs>
    </section>

    <form v-else class="drawer-form drawer-form-stack" @submit.prevent="save">
      <section class="form-section">
        <h3>{{ language === 'en' ? 'Job information' : '岗位信息' }}</h3>
        <div class="form-stack">
          <el-form-item label="公司名称" required><el-input v-model="form.company_name" maxlength="120" placeholder="例如 Stripe" /></el-form-item>
          <el-form-item label="岗位名称" required><el-input v-model="form.position_name" maxlength="120" placeholder="例如 Software Engineer Intern" /></el-form-item>
          <el-form-item label="部门 / 团队"><el-input v-model="form.department" placeholder="例如 Infrastructure" /></el-form-item>
          <el-form-item label="用工类型"><el-select v-model="form.employment_type"><el-option v-for="v in EMPLOYMENT_TYPES" :key="v" :label="v" :value="v" /></el-select></el-form-item>
          <el-form-item label="状态"><el-select v-model="form.status"><el-option v-for="s in STATUSES" :key="s" :label="language === 'en' ? statusLabel(s as ApplicationStatus) : `${s} · ${STATUS_META[s as ApplicationStatus].label}`" :value="s" /></el-select></el-form-item>
          <el-form-item label="地点"><el-input v-model="form.location" placeholder="例如 上海 / Remote" /></el-form-item>
          <el-form-item label="办公方式"><el-select v-model="form.work_mode"><el-option v-for="v in WORK_MODES" :key="v.value" :label="v.label" :value="v.value" /></el-select></el-form-item>
          <el-form-item label="优先级"><el-select v-model="form.priority"><el-option v-for="v in PRIORITIES" :key="v.value" :label="v.label" :value="v.value" /></el-select></el-form-item>
          <el-form-item label="申请链接"><el-input v-model="form.apply_link" placeholder="https://" /></el-form-item>
        </div>
      </section>

      <section class="form-section">
        <h3>{{ language === 'en' ? 'Requirements' : '岗位要求' }}</h3>
        <div class="form-stack">
          <el-form-item label="学历要求"><el-select v-model="form.education"><el-option v-for="v in EDUCATIONS" :key="v.value" :label="v.label" :value="v.value" /></el-select></el-form-item>
          <el-form-item label="经验年限"><el-input-number v-model="form.experience_years" :min="0" :max="30" /></el-form-item>
          <el-form-item label="每周到岗天数"><el-input-number v-model="form.days_per_week" :min="0" :max="7" /></el-form-item>
          <el-form-item label="最短实习月数"><el-input-number v-model="form.intern_months" :min="0" :max="36" /></el-form-item>
          <el-form-item label="薪酬">
            <div class="salary-row">
              <el-input-number v-model="form.salary_min" :min="0" :controls="false" placeholder="最低" />
              <span>至</span>
              <el-input-number v-model="form.salary_max" :min="0" :controls="false" placeholder="最高" />
              <el-select v-model="form.salary_unit" style="width:120px"><el-option v-for="v in SALARY_UNITS" :key="v.value" :label="v.label" :value="v.value" /></el-select>
            </div>
          </el-form-item>
          <el-form-item label="硬性要求 / 技术栈"><el-input v-model="form.tech_stack" placeholder="例如 Python, PyTorch" /></el-form-item>
          <el-form-item label="岗位职责"><el-input v-model="form.responsibilities" type="textarea" :rows="4" /></el-form-item>
          <el-form-item label="岗位要求"><el-input v-model="form.requirements" type="textarea" :rows="4" /></el-form-item>
          <el-form-item label="原始 JD"><el-input v-model="form.jd_text" type="textarea" :rows="5" /></el-form-item>
          <el-form-item><el-switch v-model="form.can_convert" active-text="提供转正机会" /></el-form-item>
        </div>
      </section>

      <section class="form-section">
        <h3>{{ language === 'en' ? 'Application & contact' : '投递与联系人' }}</h3>
        <div class="form-stack">
          <el-form-item label="来源类型"><el-select v-model="form.source_type" clearable allow-create filterable placeholder="请选择或输入"><el-option v-for="v in SOURCE_TYPES" :key="v" :label="v" :value="v" /></el-select></el-form-item>
          <el-form-item label="来源名称"><el-input v-model="form.source_name" placeholder="例如校招公众号 / 推荐人" /></el-form-item>
          <el-form-item label="截止日期"><el-date-picker v-model="form.deadline_at" type="date" value-format="YYYY-MM-DD" style="width:100%" /></el-form-item>
          <el-form-item label="实际投递时间"><el-date-picker v-model="form.applied_at" type="datetime" value-format="YYYY-MM-DDTHH:mm" style="width:100%" /></el-form-item>
          <el-form-item label="联系人"><el-input v-model="form.contact_name" /></el-form-item>
          <el-form-item label="联系人邮箱" :error="form.contact_email && !isValidEmail(form.contact_email) ? '邮箱格式不正确' : ''"><el-input v-model="form.contact_email" type="email" /></el-form-item>
          <el-form-item label="来源链接"><el-input v-model="form.source_url" placeholder="https://" /></el-form-item>
        </div>
      </section>

      <section class="form-section">
        <h3>{{ language === 'en' ? 'Personal records' : '个人记录' }}</h3>
        <div class="form-stack">
          <el-form-item label="使用的简历版本"><el-input v-model="form.resume_version" placeholder="例如 resume-v3-systems.pdf" /></el-form-item>
          <el-form-item label="笔记"><el-input v-model="form.notes" type="textarea" :rows="4" /></el-form-item>
          <el-form-item><el-switch v-model="form.is_liked" active-text="标记为喜欢" /></el-form-item>
        </div>
      </section>

      <p v-if="item" class="lane-hint">上次更新 {{ fmtTime(item.updated_at) }}</p>
      <div class="drawer-actions">
        <el-button type="primary" native-type="submit" :loading="saving">{{ item ? '保存更改' : '新建岗位' }}</el-button>
        <el-button v-if="item" type="danger" plain :loading="deleting" @click="remove">删除岗位</el-button>
      </div>
    </form>
  </el-drawer>
</template>
