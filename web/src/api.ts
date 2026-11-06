import type { Application, ApplicationEvent, ApplicationEventPayload, ApplicationPayload, ApplicationStatus, Overview, Paged, PreferredLanguage, ResumeAnalysis, ResumeProfile, User } from './types'

function getToken() {
  return localStorage.getItem('it_token') || ''
}

export async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers)
  if (!headers.has('Content-Type') && options.body && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json')
  }
  const token = getToken()
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(path, { ...options, headers })
  const data = await response.json().catch(() => ({}))
  if (!response.ok) {
    throw new Error(data.error || `请求失败（${response.status}）`)
  }
  return data as T
}

export const api = {
  register(payload: { email: string; username: string; password: string; preferred_language: PreferredLanguage }) {
    return request<{ token: string; user: User }>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },
  login(payload: { email: string; password: string; preferred_language?: PreferredLanguage }) {
    return request<{ token: string; user: User }>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },
  me() {
    return request<{ user: User }>('/api/auth/me')
  },
  updateProfile(payload: { username: string; password: string; preferred_language: PreferredLanguage }) {
    return request<{ user: User }>('/api/auth/profile', { method: 'PATCH', body: JSON.stringify(payload) })
  },
  listApplications(keyword = '', likedOnly = false) {
    const search = new URLSearchParams()
    if (keyword) search.set('keyword', keyword)
    if (likedOnly) search.set('liked', '1')
    const query = search.toString()
    return request<{ items: Application[] }>(`/api/applications${query ? `?${query}` : ''}`)
  },
  setApplicationLiked(id: string, isLiked: boolean) {
    return request<{ item: Application }>(`/api/applications/${id}/liked`, {
      method: 'PATCH',
      body: JSON.stringify({ is_liked: isLiked }),
    })
  },
  createApplication(payload: ApplicationPayload) {
    return request<{ item: Application }>('/api/applications', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },
  updateApplication(id: string, payload: ApplicationPayload) {
    return request<{ item: Application }>(`/api/applications/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },
  getApplication(id: string) {
    return request<{ item: Application }>(`/api/applications/${id}`)
  },
  changeApplicationStatus(id: string, payload: { status: ApplicationStatus; title?: string; content?: string; event_at?: string; details?: Record<string, unknown> }) {
    return request<{ item: Application }>(`/api/applications/${id}/status`, { method: 'PATCH', body: JSON.stringify(payload) })
  },
  listApplicationEvents(id: string) {
    return request<{ items: ApplicationEvent[] }>(`/api/applications/${id}/events`)
  },
  listTimelineEvents() {
    return request<{ items: ApplicationEvent[] }>('/api/applications/timeline')
  },
  createApplicationEvent(id: string, payload: ApplicationEventPayload) {
    return request<{ item: ApplicationEvent }>(`/api/applications/${id}/events`, { method: 'POST', body: JSON.stringify(payload) })
  },
  updateApplicationEvent(applicationId: string, eventId: string, payload: ApplicationEventPayload) {
    return request<{ item: ApplicationEvent }>(`/api/applications/${applicationId}/events/${eventId}`, { method: 'PUT', body: JSON.stringify(payload) })
  },
  deleteApplicationEvent(applicationId: string, eventId: string) {
    return request<{ deleted: boolean }>(`/api/applications/${applicationId}/events/${eventId}`, { method: 'DELETE' })
  },
  deleteApplication(id: string) {
    return request<{ deleted: boolean }>(`/api/applications/${id}`, { method: 'DELETE' })
  },
  getResume() {
    return request<{ item: ResumeProfile | null }>('/api/resume')
  },
  saveResume(payload: { title: string; source_name: string; content: string }) {
    return request<{ item: ResumeProfile }>('/api/resume', { method: 'PUT', body: JSON.stringify(payload) })
  },
  extractResume(file: File) {
    const body = new FormData()
    body.append('file', file)
    return request<{ content: string; source_name: string }>('/api/resume/extract', { method: 'POST', body })
  },
  listResumeAnalyses() {
    return request<{ items: ResumeAnalysis[] }>('/api/resume/analyses')
  },
  analyzeResume(applicationId: string) {
    return request<{ item: ResumeAnalysis }>('/api/resume/analyze', { method: 'POST', body: JSON.stringify({ application_id: applicationId }) })
  },
  adminOverview() {
    return request<Overview>('/api/admin/overview')
  },
  adminUsers(params: { keyword?: string; page?: number; page_size?: number }) {
    const search = new URLSearchParams()
    if (params.keyword) search.set('keyword', params.keyword)
    search.set('page', String(params.page ?? 1))
    search.set('page_size', String(params.page_size ?? 20))
    return request<Paged<User>>(`/api/admin/users?${search.toString()}`)
  },
  adminSetUserEnabled(id: string, isEnabled: boolean) {
    return request<{ user: User }>(`/api/admin/users/${id}/enabled`, {
      method: 'PATCH',
      body: JSON.stringify({ is_enabled: isEnabled }),
    })
  },
  adminApplications(params: { keyword?: string; user_id?: string; page?: number; page_size?: number }) {
    const search = new URLSearchParams()
    if (params.keyword) search.set('keyword', params.keyword)
    if (params.user_id) search.set('user_id', params.user_id)
    search.set('page', String(params.page ?? 1))
    search.set('page_size', String(params.page_size ?? 20))
    return request<Paged<Application>>(`/api/admin/applications?${search.toString()}`)
  },
}
