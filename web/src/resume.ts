export interface ResumeFields {
  full_name: string
  email: string
  phone: string
  location: string
  links: string
  summary: string
  education: string
  experience: string
  projects: string
  skills: string
  extras: string
}

export type ResumeDocument = { v: 2 } & ResumeFields

export function emptyResumeFields(): ResumeFields {
  return {
    full_name: '',
    email: '',
    phone: '',
    location: '',
    links: '',
    summary: '',
    education: '',
    experience: '',
    projects: '',
    skills: '',
    extras: '',
  }
}

export function parseResumeContent(raw: string): ResumeFields {
  const trimmed = raw.trim()
  if (!trimmed) return emptyResumeFields()
  try {
    const parsed = JSON.parse(trimmed) as Partial<ResumeDocument>
    if (parsed && parsed.v === 2) {
      const base = emptyResumeFields()
      for (const key of Object.keys(base) as (keyof ResumeFields)[]) {
        base[key] = typeof parsed[key] === 'string' ? parsed[key]! : ''
      }
      return base
    }
  } catch {
    // legacy plain-text resume
  }
  return { ...emptyResumeFields(), extras: raw }
}

export function serializeResumeContent(fields: ResumeFields): string {
  const doc: ResumeDocument = { v: 2, ...fields }
  return JSON.stringify(doc)
}

export function resumeHasContent(fields: ResumeFields): boolean {
  return Object.values(fields).some((value) => value.trim().length > 0)
}

export function composeResumePlainText(fields: ResumeFields, labels?: Partial<Record<keyof ResumeFields, string>>): string {
  const L = {
    full_name: labels?.full_name || 'Name',
    email: labels?.email || 'Email',
    phone: labels?.phone || 'Phone',
    location: labels?.location || 'Location',
    links: labels?.links || 'Links',
    summary: labels?.summary || 'Summary',
    education: labels?.education || 'Education',
    experience: labels?.experience || 'Experience',
    projects: labels?.projects || 'Projects',
    skills: labels?.skills || 'Skills',
    extras: labels?.extras || 'Additional',
  }
  const lines: string[] = []
  if (fields.full_name.trim()) lines.push(fields.full_name.trim())
  const contact = [fields.email, fields.phone, fields.location, fields.links].map((v) => v.trim()).filter(Boolean)
  if (contact.length) lines.push(contact.join(' · '))
  const sections: (keyof ResumeFields)[] = ['summary', 'education', 'experience', 'projects', 'skills', 'extras']
  for (const key of sections) {
    const value = fields[key].trim()
    if (!value) continue
    lines.push('')
    lines.push(`${L[key]}`)
    lines.push(value)
  }
  return lines.join('\n').trim()
}
