import type { ApplicationStatus } from './types'

export const STATUSES: ApplicationStatus[] = [
  'Wishlist',
  'Applied',
  'OA',
  'Interview',
  'Offer',
  'Rejected',
  'Withdrawn',
]

export const STATUS_META: Record<ApplicationStatus, { label: string; hint: string }> = {
  Wishlist: { label: '想投', hint: '还没提交' },
  Applied: { label: '已投', hint: '材料已送出' },
  OA: { label: '笔试', hint: '线上作业 / OA' },
  Interview: { label: '面试', hint: '电话 / 视频 / onsite' },
  Offer: { label: 'Offer', hint: '已拿到' },
  Rejected: { label: '结束', hint: '被拒或主动放弃' },
  Withdrawn: { label: '已放弃', hint: '主动结束流程' },
}

export const PRIORITIES = [
  { value: '', label: '未设' },
  { value: 'High', label: '高' },
  { value: 'Medium', label: '中' },
  { value: 'Low', label: '低' },
]

export const SALARY_UNITS = [
  { value: '元/天', label: '元/天' },
  { value: '元/月', label: '元/月' },
  { value: 'k/月', label: 'k/月' },
  { value: 'USD/hr', label: 'USD/hr' },
  { value: 'USD/month', label: 'USD/month' },
  { value: 'USD/year', label: 'USD/year' },
]

export const EDUCATIONS = [
  { value: '', label: '不限' },
  { value: '大专', label: '大专' },
  { value: '本科', label: '本科' },
  { value: '硕士', label: '硕士' },
  { value: '博士', label: '博士' },
]

export const WORK_MODES = [
  { value: '', label: '未设' },
  { value: '线下', label: '线下' },
  { value: '混合', label: '混合' },
  { value: '远程', label: '远程' },
]

export const SOURCE_TYPES = ['公司官网', 'LinkedIn', 'Indeed', 'Boss直聘', '校园招聘', '内推', '其他']
export const EMPLOYMENT_TYPES = ['实习', '全职', '兼职', '合同制']
