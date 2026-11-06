<script setup lang="ts">
import { Edit, Star, StarFilled } from '@element-plus/icons-vue'
import type { Application, ApplicationStatus } from '../types'
import { fmtDate, fmtSalary } from '../utils'
import { useI18n } from '../i18n'
import { STATUS_META } from '../constants'

const props = defineProps<{
  item: Application
  index: number
}>()

const { language, statusLabel } = useI18n()

type DisplayStage = Exclude<ApplicationStatus, 'Rejected' | 'Withdrawn'> | 'Final'

function phaseOf(item: Application): DisplayStage {
  return ['Rejected', 'Withdrawn'].includes(item.status) ? 'Final' : (item.status as DisplayStage)
}

function stageLabel(stage: DisplayStage) {
  if (stage === 'Final') return language.value === 'en' ? 'Closed' : '已结束'
  return language.value === 'en' ? statusLabel(stage) : STATUS_META[stage].label
}

const emit = defineEmits<{
  open: []
  edit: []
  toggleLike: []
}>()
</script>

<template>
  <li class="application-list-row" @click="emit('open')">
    <span class="list-cell list-cell-index">{{ index }}</span>
    <span class="list-cell list-cell-company" :title="item.company_name">{{ item.company_name }}</span>
    <span class="list-cell list-cell-position" :title="item.position_name">{{ item.position_name }}</span>
    <span class="list-cell list-cell-stage">
      <span :class="`table-phase phase-${phaseOf(item).toLowerCase()}`">{{ stageLabel(phaseOf(item)) }}</span>
    </span>
    <span class="list-cell list-cell-location" :title="item.location || ''">{{ item.location || '—' }}</span>
    <span class="list-cell list-cell-salary">{{ fmtSalary(item.salary_min, item.salary_max, item.salary_unit) || '—' }}</span>
    <span class="list-cell list-cell-date">{{ fmtDate(item.deadline_at) || '—' }}</span>
    <span class="list-cell list-cell-date">{{ fmtDate(item.created_at) || '—' }}</span>
    <div class="list-cell list-cell-actions no-drag" @click.stop>
      <button
        class="ticket-icon-btn"
        type="button"
        :title="language === 'en' ? 'Edit' : '修改'"
        :aria-label="language === 'en' ? 'Edit' : '修改'"
        @click="emit('edit')"
      >
        <el-icon><Edit /></el-icon>
      </button>
      <button
        class="ticket-icon-btn"
        type="button"
        :aria-pressed="item.is_liked"
        :title="language === 'en' ? 'Favorite' : '收藏'"
        :aria-label="language === 'en' ? 'Favorite' : '收藏'"
        @click="emit('toggleLike')"
      >
        <el-icon>
          <StarFilled v-if="item.is_liked" />
          <Star v-else />
        </el-icon>
      </button>
    </div>
  </li>
</template>
