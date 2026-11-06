<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { DocumentAdd, MagicStick, UploadFilled, View } from '@element-plus/icons-vue'
import AppHeader from '../components/AppHeader.vue'
import { api } from '../api'
import { useI18n } from '../i18n'
import type { Application, ResumeAnalysis } from '../types'
import {
  composeResumePlainText,
  emptyResumeFields,
  parseResumeContent,
  resumeHasContent,
  serializeResumeContent,
  type ResumeFields,
} from '../resume'

const { language } = useI18n()
const en = computed(() => language.value === 'en')
const title = ref(en.value ? 'My Resume' : '我的简历')
const sourceName = ref('')
const fields = ref<ResumeFields>(emptyResumeFields())
const panel = ref<'edit' | 'preview'>('edit')
const applications = ref<Application[]>([])
const analyses = ref<ResumeAnalysis[]>([])
const selectedApplication = ref('')
const activeAnalysis = ref<ResumeAnalysis | null>(null)
const loading = ref(true)
const saving = ref(false)
const analyzing = ref(false)
const extracting = ref(false)
const saved = ref(false)

const fieldLabels = computed(() =>
  en.value
    ? {
        full_name: 'Full name',
        email: 'Email',
        phone: 'Phone',
        location: 'City / location',
        links: 'Links',
        summary: 'Summary',
        education: 'Education',
        experience: 'Experience',
        projects: 'Projects',
        skills: 'Skills',
        extras: 'Certificates / other',
      }
    : {
        full_name: '姓名',
        email: '邮箱',
        phone: '手机',
        location: '城市 / 地点',
        links: '链接',
        summary: '个人简介',
        education: '教育经历',
        experience: '实习 / 工作经历',
        projects: '项目经历',
        skills: '技能',
        extras: '证书 / 其他',
      },
)

const previewText = computed(() => composeResumePlainText(fields.value, fieldLabels.value))
const previewSections = computed(() => {
  const L = fieldLabels.value
  const f = fields.value
  return [
    { key: 'summary', title: L.summary, body: f.summary },
    { key: 'education', title: L.education, body: f.education },
    { key: 'experience', title: L.experience, body: f.experience },
    { key: 'projects', title: L.projects, body: f.projects },
    { key: 'skills', title: L.skills, body: f.skills },
    { key: 'extras', title: L.extras, body: f.extras },
  ].filter((section) => section.body.trim())
})

const eligibleApplications = computed(() =>
  applications.value.filter((item) =>
    [item.jd_text, item.requirements, item.responsibilities, item.tech_stack].some((value) => value?.trim()),
  ),
)

function markChanged() {
  saved.value = false
}

async function load() {
  loading.value = true
  try {
    const [resumeResponse, applicationResponse, analysisResponse] = await Promise.all([
      api.getResume(),
      api.listApplications(),
      api.listResumeAnalyses(),
    ])
    if (resumeResponse.item) {
      title.value = resumeResponse.item.title
      sourceName.value = resumeResponse.item.source_name
      fields.value = parseResumeContent(resumeResponse.item.content)
      saved.value = true
    } else {
      title.value = en.value ? 'My Resume' : '我的简历'
    }
    applications.value = applicationResponse.items
    analyses.value = analysisResponse.items
    activeAnalysis.value = analyses.value[0] || null
    selectedApplication.value = activeAnalysis.value?.application_id || eligibleApplications.value[0]?.id || ''
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : en.value ? 'Unable to load resume.' : '无法加载简历。')
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!resumeHasContent(fields.value)) {
    ElMessage.warning(en.value ? 'Fill in at least one resume field.' : '请至少填写一个简历字段。')
    return false
  }
  saving.value = true
  try {
    await api.saveResume({
      title: title.value,
      source_name: sourceName.value,
      content: serializeResumeContent(fields.value),
    })
    saved.value = true
    ElMessage.success(en.value ? 'Resume saved.' : '简历已保存。')
    return true
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : en.value ? 'Unable to save resume.' : '保存失败。')
    return false
  } finally {
    saving.value = false
  }
}

async function analyze() {
  if (!selectedApplication.value) {
    ElMessage.warning(en.value ? 'Select an application with job information.' : '请选择包含岗位信息的申请。')
    return
  }
  if (!saved.value && !(await save())) return
  analyzing.value = true
  try {
    const response = await api.analyzeResume(selectedApplication.value)
    activeAnalysis.value = response.item
    analyses.value = [response.item, ...analyses.value]
    ElMessage.success(en.value ? 'Compatibility analysis completed.' : '匹配度分析已完成。')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : en.value ? 'Unable to run analysis.' : '分析失败。')
  } finally {
    analyzing.value = false
  }
}

async function importResume(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  extracting.value = true
  try {
    const result = await api.extractResume(file)
    fields.value = { ...emptyResumeFields(), extras: result.content }
    sourceName.value = result.source_name
    if (title.value === 'My Resume' || title.value === '我的简历') {
      title.value = file.name.replace(/\.[^.]+$/, '')
    }
    saved.value = false
    panel.value = 'edit'
    ElMessage.success(
      en.value
        ? 'Text extracted into Additional. Move it into the right fields, then save.'
        : '文字已提取到「证书 / 其他」，请整理到对应字段后再保存。',
    )
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : en.value ? 'Unable to read this file.' : '无法读取该文件。')
  } finally {
    extracting.value = false
    input.value = ''
  }
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat(en.value ? 'en-US' : 'zh-CN', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
}

onMounted(load)
</script>

<template>
  <div class="app-shell">
    <AppHeader />
    <main class="workspace resume-page" v-loading="loading">
      <section class="board-toolbar-strip">
        <div class="view-switch" role="group" :aria-label="en ? 'Resume panel' : '简历面板'">
          <button type="button" :class="{ active: panel === 'edit' }" @click="panel = 'edit'">
            {{ en ? 'Edit' : '编辑' }}
          </button>
          <button type="button" :class="{ active: panel === 'preview' }" @click="panel = 'preview'">
            <el-icon><View /></el-icon>
            {{ en ? 'Preview' : '预览' }}
          </button>
        </div>
        <el-input
          v-model="title"
          class="board-search"
          :placeholder="en ? 'Resume name' : '简历名称'"
          @input="markChanged"
        />
        <span class="save-state" :class="{ saved }">{{ saved ? (en ? 'Saved' : '已保存') : (en ? 'Unsaved' : '未保存') }}</span>
        <button class="toolbar-add" type="button" :disabled="saving" @click="save">
          <el-icon><DocumentAdd /></el-icon>
          <span>{{ saving ? (en ? 'Saving…' : '保存中…') : (en ? 'Save' : '保存') }}</span>
        </button>
      </section>

      <section class="resume-layout">
        <article class="page-panel resume-editor-panel">
          <template v-if="panel === 'edit'">
            <div class="resume-import-row">
              <label class="text-import-button" :class="{ disabled: extracting }">
                <el-icon class="upload-icon"><UploadFilled /></el-icon>
                {{ extracting ? (en ? 'Reading…' : '读取中…') : (en ? 'Upload file' : '上传文件') }}
                <input
                  type="file"
                  :disabled="extracting"
                  accept=".pdf,.docx,.txt,.md,application/pdf,application/vnd.openxmlformats-officedocument.wordprocessingml.document,text/plain,text/markdown"
                  @change="importResume"
                />
              </label>
              <span>{{ sourceName || (en ? 'PDF / DOCX / TXT / MD · max 10 MB' : '支持 PDF、DOCX、TXT、MD · 最大 10 MB') }}</span>
            </div>

            <div class="resume-form-grid">
              <label class="resume-field">
                <span>{{ fieldLabels.full_name }}</span>
                <el-input v-model="fields.full_name" @input="markChanged" />
              </label>
              <label class="resume-field">
                <span>{{ fieldLabels.email }}</span>
                <el-input v-model="fields.email" @input="markChanged" />
              </label>
              <label class="resume-field">
                <span>{{ fieldLabels.phone }}</span>
                <el-input v-model="fields.phone" @input="markChanged" />
              </label>
              <label class="resume-field">
                <span>{{ fieldLabels.location }}</span>
                <el-input v-model="fields.location" @input="markChanged" />
              </label>
              <label class="resume-field resume-field--full">
                <span>{{ fieldLabels.links }}</span>
                <el-input
                  v-model="fields.links"
                  :placeholder="en ? 'Portfolio / GitHub / LinkedIn' : '作品集 / GitHub / LinkedIn'"
                  @input="markChanged"
                />
              </label>
              <label class="resume-field resume-field--full">
                <span>{{ fieldLabels.summary }}</span>
                <el-input v-model="fields.summary" type="textarea" :rows="3" @input="markChanged" />
              </label>
              <label class="resume-field resume-field--full">
                <span>{{ fieldLabels.education }}</span>
                <el-input
                  v-model="fields.education"
                  type="textarea"
                  :rows="4"
                  :placeholder="en ? 'School · Major · Dates · Highlights' : '学校 · 专业 · 时间 · 要点'"
                  @input="markChanged"
                />
              </label>
              <label class="resume-field resume-field--full">
                <span>{{ fieldLabels.experience }}</span>
                <el-input
                  v-model="fields.experience"
                  type="textarea"
                  :rows="5"
                  :placeholder="en ? 'Company · Role · Dates · Impact' : '公司 · 职位 · 时间 · 成果'"
                  @input="markChanged"
                />
              </label>
              <label class="resume-field resume-field--full">
                <span>{{ fieldLabels.projects }}</span>
                <el-input v-model="fields.projects" type="textarea" :rows="5" @input="markChanged" />
              </label>
              <label class="resume-field resume-field--full">
                <span>{{ fieldLabels.skills }}</span>
                <el-input
                  v-model="fields.skills"
                  type="textarea"
                  :rows="3"
                  :placeholder="en ? 'Languages, frameworks, tools' : '语言、框架、工具'"
                  @input="markChanged"
                />
              </label>
              <label class="resume-field resume-field--full">
                <span>{{ fieldLabels.extras }}</span>
                <el-input v-model="fields.extras" type="textarea" :rows="3" @input="markChanged" />
              </label>
            </div>
          </template>

          <article v-else class="resume-preview">
            <header class="resume-preview-head">
              <h2>{{ fields.full_name || title || (en ? 'Untitled resume' : '未命名简历') }}</h2>
              <p v-if="[fields.email, fields.phone, fields.location, fields.links].some((v) => v.trim())">
                {{ [fields.email, fields.phone, fields.location, fields.links].filter((v) => v.trim()).join(' · ') }}
              </p>
            </header>
            <section v-for="section in previewSections" :key="section.key" class="resume-preview-section">
              <h3>{{ section.title }}</h3>
              <pre>{{ section.body }}</pre>
            </section>
            <p v-if="!previewText" class="list-empty">{{ en ? 'Nothing to preview yet.' : '还没有可预览的内容。' }}</p>
          </article>
        </article>

        <aside class="resume-side">
          <article class="page-panel">
            <h2 class="panel-title">{{ en ? 'AI match' : 'AI 匹配' }}</h2>
            <p class="panel-desc">{{ en ? 'Compare this resume with a tracked role.' : '将当前简历与已跟踪岗位进行匹配。' }}</p>
            <el-select
              v-model="selectedApplication"
              filterable
              :placeholder="en ? 'Select an application' : '选择岗位'"
              style="width: 100%"
            >
              <el-option
                v-for="item in eligibleApplications"
                :key="item.id"
                :value="item.id"
                :label="`${item.company_name} — ${item.position_name}`"
              />
            </el-select>
            <p v-if="!eligibleApplications.length" class="analysis-warning">
              {{ en ? 'No roles have job details yet. Add a JD on the Board first.' : '暂无含岗位信息的申请，请先在看板补充 JD。' }}
            </p>
            <el-button
              class="analyze-button"
              type="primary"
              :icon="MagicStick"
              :loading="analyzing"
              :disabled="!eligibleApplications.length"
              @click="analyze"
            >
              {{ en ? 'Analyze' : '分析匹配度' }}
            </el-button>
            <small class="panel-hint">{{ en ? 'AI advice can be wrong — review before editing.' : 'AI 建议可能有误，修改前请自行核对。' }}</small>
          </article>

          <article v-if="activeAnalysis" class="page-panel analysis-result">
            <div class="score-row">
              <div class="score-ring" :style="{ '--score': activeAnalysis.result.score }">
                <strong>{{ activeAnalysis.result.score }}</strong>
                <span>/100</span>
              </div>
              <div>
                <p class="panel-kicker">{{ activeAnalysis.company_name }}</p>
                <h2 class="panel-title">{{ activeAnalysis.position_name }}</h2>
                <small>{{ formatDate(activeAnalysis.created_at) }}</small>
              </div>
            </div>
            <p class="analysis-summary">{{ activeAnalysis.result.summary }}</p>
            <div class="analysis-section strengths">
              <h3>{{ en ? 'Strengths' : '优势' }}</h3>
              <ul><li v-for="item in activeAnalysis.result.strengths" :key="item">{{ item }}</li></ul>
            </div>
            <div class="analysis-section gaps">
              <h3>{{ en ? 'Gaps' : '差距' }}</h3>
              <ul><li v-for="item in activeAnalysis.result.gaps" :key="item">{{ item }}</li></ul>
            </div>
            <div class="analysis-section">
              <h3>{{ en ? 'Missing keywords' : '缺少关键词' }}</h3>
              <div class="keyword-list"><span v-for="item in activeAnalysis.result.missing_keywords" :key="item">{{ item }}</span></div>
            </div>
            <div class="analysis-section improvements">
              <h3>{{ en ? 'Improvements' : '改进建议' }}</h3>
              <ol><li v-for="item in activeAnalysis.result.improvements" :key="item">{{ item }}</li></ol>
            </div>
          </article>

          <article v-if="analyses.length" class="page-panel">
            <h2 class="panel-title">{{ en ? 'History' : '历史分析' }}</h2>
            <div class="history-list compact">
              <button
                v-for="item in analyses"
                :key="item.id"
                type="button"
                :class="{ active: activeAnalysis?.id === item.id }"
                @click="activeAnalysis = item; selectedApplication = item.application_id"
              >
                <strong>{{ item.result.score }}</strong>
                <span>{{ item.company_name }}<small>{{ item.position_name }} · {{ formatDate(item.created_at) }}</small></span>
              </button>
            </div>
          </article>
        </aside>
      </section>
    </main>
  </div>
</template>
