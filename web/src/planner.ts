import type { ApplicationEvent, ApplicationStatus, PreferredLanguage } from './types'

export function localDateKey(value: string | Date) {
  const date = value instanceof Date ? value : new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return `${date.getFullYear()}-${String(date.getMonth()+1).padStart(2,'0')}-${String(date.getDate()).padStart(2,'0')}`
}

export function calendarKind(event: Pick<ApplicationEvent,'event_type'>) {
  if (event.event_type === 'created') return null
  return event.event_type === 'interview' ? 'interview' : 'workflow'
}

export function statusOptionLabel(status:ApplicationStatus, language:PreferredLanguage, chineseLabel:string) {
  if (language === 'en') return status === 'OA' ? 'Online assessment' : status
  return `${status} · ${chineseLabel}`
}
