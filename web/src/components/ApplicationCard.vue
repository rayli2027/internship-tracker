<script setup lang="ts">
import { Edit, Star, StarFilled } from '@element-plus/icons-vue'
import type { Application } from '../types'
import { fmtDate } from '../utils'
import { useI18n } from '../i18n'
import { STATUS_META } from '../constants'

const props = defineProps<{
  item: Application
  dragging?: boolean
}>()

const { language, statusLabel } = useI18n()
const stageText = () => {
  if (['Rejected', 'Withdrawn'].includes(props.item.status)) {
    return language.value === 'en' ? 'Closed' : '已结束'
  }
  return language.value === 'en'
    ? statusLabel(props.item.status)
    : STATUS_META[props.item.status].label
}

const emit = defineEmits<{
  open: []
  edit: []
  toggleLike: []
  dragStart: [event: DragEvent]
  dragEnd: []
}>()

function onDragStart(event: DragEvent) {
  const target = event.target as HTMLElement
  if (target.closest('.no-drag')) {
    event.preventDefault()
    return
  }
  emit('dragStart', event)
}

function onCardClick(event: MouseEvent) {
  if ((event.target as HTMLElement).closest('.no-drag')) return
  emit('open')
}
</script>

<template>
  <article
    class="ticket ticket-compact"
    :class="{ 'is-dragging': dragging, 'is-liked': item.is_liked }"
    :data-status="item.status"
    draggable="true"
    tabindex="0"
    role="button"
    @click="onCardClick"
    @dragstart="onDragStart"
    @dragend="emit('dragEnd')"
  >
    <div class="ticket-row-main">
      <div class="ticket-title-block">
        <span class="ticket-company">{{ item.company_name }}</span>
        <span class="ticket-role">{{ item.position_name }}</span>
      </div>
      <div class="ticket-icon-actions no-drag" @mousedown.stop>
        <button
          class="ticket-icon-btn"
          type="button"
          :title="language === 'en' ? 'Edit' : '修改'"
          :aria-label="language === 'en' ? 'Edit' : '修改'"
          @click.stop="emit('edit')"
        >
          <el-icon><Edit /></el-icon>
        </button>
        <button
          class="ticket-icon-btn"
          type="button"
          :aria-pressed="item.is_liked"
          :title="language === 'en' ? 'Favorite' : '收藏'"
          :aria-label="language === 'en' ? 'Favorite' : '收藏'"
          @click.stop="emit('toggleLike')"
        >
          <el-icon>
            <StarFilled v-if="item.is_liked" />
            <Star v-else />
          </el-icon>
        </button>
      </div>
    </div>
    <p class="ticket-meta-line">
      <span class="ticket-stage">{{ stageText() }}</span>
      <span v-if="item.location">{{ item.location }}</span>
      <span v-if="fmtDate(item.deadline_at)">{{ language === 'en' ? 'Due' : '截止' }} {{ fmtDate(item.deadline_at) }}</span>
      <span v-else-if="fmtDate(item.applied_at)">{{ language === 'en' ? 'Applied' : '投递' }} {{ fmtDate(item.applied_at) }}</span>
    </p>
  </article>
</template>
