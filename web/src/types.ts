export type UserRole = 'user' | 'admin'
export type PreferredLanguage = 'zh' | 'en'

export type ApplicationStatus = 'Wishlist' | 'Applied' | 'OA' | 'Interview' | 'Offer' | 'Rejected' | 'Withdrawn'

export type EventType = 'created' | 'applied' | 'oa' | 'interview' | 'offer' | 'rejected' | 'withdrawn' | 'memo' | 'status_change'

export interface ApplicationEvent {
  id: string
  application_id: string
  event_type: EventType
  from_status: string
  to_status: string
  title: string
  content: string
  details: Record<string, unknown>
  event_at: string
  created_at: string
  updated_at: string
  company_name?: string
  position_name?: string
}

export interface ApplicationEventPayload {
  event_type: Exclude<EventType, 'status_change' | 'created'>
  title: string
  content: string
  event_at: string
  details: Record<string, unknown>
}

export type Priority = '' | 'High' | 'Medium' | 'Low'

export interface User {
  id: string
  email: string
  username: string
  preferred_language: PreferredLanguage
  role: UserRole
  is_enabled: boolean
  created_at: string
  updated_at: string
}

export interface Application {
  id: string
  company_name: string
  position_name: string
  user_id: string
  user_name: string
  status: ApplicationStatus
  apply_link: string
  location: string
  deadline_at: string
  jd_text: string
  notes: string
  resume_version: string
  priority: Priority | string
  salary_min: number
  salary_max: number
  salary_unit: string
  experience_years: number
  education: string
  work_mode: string
  days_per_week: number
  intern_months: number
  can_convert: boolean
  tech_stack: string
  is_liked: boolean
  source_type: string
  source_name: string
  source_url: string
  employment_type: string
  department: string
  requirements: string
  responsibilities: string
  contact_name: string
  contact_email: string
  applied_at: string
  created_at: string
  updated_at: string
}

export interface ApplicationPayload {
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
  salary_min: number
  salary_max: number
  salary_unit: string
  experience_years: number
  education: string
  work_mode: string
  days_per_week: number
  intern_months: number
  can_convert: boolean
  tech_stack: string
  is_liked: boolean
  source_type: string
  source_name: string
  source_url: string
  employment_type: string
  department: string
  requirements: string
  responsibilities: string
  contact_name: string
  contact_email: string
  applied_at: string
}

export interface Overview {
  user_count: number
  application_count: number
  today_user_count: number
  today_application_count: number
}

export interface Paged<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export interface ResumeProfile {
  user_id: string
  title: string
  source_name: string
  content: string
  created_at: string
  updated_at: string
}

export interface ResumeAnalysisResult {
  score: number
  summary: string
  strengths: string[]
  gaps: string[]
  missing_keywords: string[]
  improvements: string[]
}

export interface ResumeAnalysis {
  id: string
  application_id: string
  company_name: string
  position_name: string
  model: string
  result: ResumeAnalysisResult
  created_at: string
}
