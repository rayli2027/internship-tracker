import type { ApplicationStatus } from './types'

const EMPTY_DEADLINE_PREFIX = '1970-01-01'

export function fmtTime(value: string | null | undefined) {
  if (!value || value.startsWith(EMPTY_DEADLINE_PREFIX)) {
    return ''
  }
  return value.replace('T', ' ').slice(0, 19)
}

export function fmtDate(value: string | null | undefined) {
  const full = fmtTime(value)
  return full ? full.slice(0, 10) : ''
}

export function isEmptyDeadline(value: string | null | undefined) {
  return !value || value.startsWith(EMPTY_DEADLINE_PREFIX)
}

export function daysUntil(value: string | null | undefined) {
  if (isEmptyDeadline(value)) {
    return null
  }
  const deadline = new Date(value as string)
  if (Number.isNaN(deadline.getTime())) {
    return null
  }
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  deadline.setHours(0, 0, 0, 0)
  return Math.round((deadline.getTime() - today.getTime()) / 86400000)
}

export function emptyApplicationForm() {
  return {
    company_name: '',
    position_name: '',
    status: 'Wishlist' as ApplicationStatus,
    apply_link: '',
    location: '',
    deadline_at: '',
    jd_text: '',
    notes: '',
    resume_version: '',
    priority: '',
    salary_min: 0,
    salary_max: 0,
    salary_unit: 'USD/hr',
    experience_years: 0,
    education: '',
    work_mode: '',
    days_per_week: 0,
    intern_months: 0,
    can_convert: false,
    tech_stack: '',
    is_liked: false,
    source_type: '',
    source_name: '',
    source_url: '',
    employment_type: '实习',
    department: '',
    requirements: '',
    responsibilities: '',
    contact_name: '',
    contact_email: '',
    applied_at: '',
  }
}

export function applicationToForm(item: {
  company_name: string
  position_name: string
  status: ApplicationStatus
  apply_link: string
  location: string
  deadline_at: string
  jd_text: string
  notes: string
  resume_version: string
  priority: string
  salary_min?: number
  salary_max?: number
  salary_unit?: string
  experience_years?: number
  education?: string
  work_mode?: string
  days_per_week?: number
  intern_months?: number
  can_convert?: boolean
  tech_stack?: string
  is_liked?: boolean
  source_type?: string
  source_name?: string
  source_url?: string
  employment_type?: string
  department?: string
  requirements?: string
  responsibilities?: string
  contact_name?: string
  contact_email?: string
  applied_at?: string
}) {
  const form = emptyApplicationForm()
  Object.assign(form, {
    company_name: item.company_name,
    position_name: item.position_name,
    status: item.status,
    apply_link: item.apply_link,
    location: item.location,
    deadline_at: isEmptyDeadline(item.deadline_at) ? '' : item.deadline_at.slice(0, 10),
    jd_text: item.jd_text,
    notes: item.notes,
    resume_version: item.resume_version,
    priority: item.priority,
    salary_min: item.salary_min ?? 0,
    salary_max: item.salary_max ?? 0,
    salary_unit: item.salary_unit || 'USD/hr',
    experience_years: item.experience_years ?? 0,
    education: item.education ?? '',
    work_mode: item.work_mode ?? '',
    days_per_week: item.days_per_week ?? 0,
    intern_months: item.intern_months ?? 0,
    can_convert: Boolean(item.can_convert),
    tech_stack: item.tech_stack ?? '',
    is_liked: Boolean(item.is_liked),
    source_type: item.source_type ?? '',
    source_name: item.source_name ?? '',
    source_url: item.source_url ?? '',
    employment_type: item.employment_type || '实习',
    department: item.department ?? '',
    requirements: item.requirements ?? '',
    responsibilities: item.responsibilities ?? '',
    contact_name: item.contact_name ?? '',
    contact_email: item.contact_email ?? '',
    applied_at: isEmptyDeadline(item.applied_at) ? '' : (item.applied_at?.slice(0, 16) ?? ''),
  })
  return form
}

export function toPayload(form: ReturnType<typeof emptyApplicationForm>) {
  return {
    ...form,
    deadline_at: form.deadline_at ? `${form.deadline_at} 00:00:00` : '',
    applied_at: form.applied_at ? `${form.applied_at.replace('T', ' ')}:00` : '',
  }
}

export function isValidEmail(value: string) {
  const email = value.trim()
  return email.length <= 300 && /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/.test(email)
}

export function fmtSalary(min: number, max: number, unit: string) {
  if (!min && !max) {
    return ''
  }
  const suffix = unit || 'USD/hr'
  if (min && max && min !== max) {
    return `${trimNum(min)}–${trimNum(max)} ${suffix}`
  }
  return `${trimNum(min || max)} ${suffix}`
}

export function fmtExperience(years: number) {
  if (!years) {
    return '年限不限'
  }
  return `要求 ${years} 年`
}

function trimNum(value: number) {
  return Number.isInteger(value) ? String(value) : value.toFixed(1)
}
