import { computed, nextTick, ref, watch } from 'vue'
import { useAuth } from './auth'
import type { ApplicationStatus, PreferredLanguage } from './types'

const messages = {
  en: {
    login: 'Sign in', register: 'Create account', email: 'Email', username: 'Username', password: 'Password', preferredLanguage: 'Preferred language', chinese: 'Chinese', english: 'English', createAccount: 'Create account', haveAccount: 'Already have an account?', noAccount: 'New here?', usernameHint: 'Optional — your email will be used by default.', registrationRequired: 'Enter your email and a password of at least 6 characters.', registrationSuccess: 'Account created.', registrationFailed: 'Registration failed', loginSuccess: 'Signed in.', loginFailed: 'Sign-in failed', emailPasswordRequired: 'Enter your email and password.', rememberAccount: 'Remember account', brand: 'Internship Tracker', admin: 'Admin', board: 'Board', logout: 'Sign out', search: 'Search by company or role', favorites: 'Favorites only', allApplications: 'All applications', newApplication: 'New application', loadFailed: 'Unable to load data.',
  },
  zh: {
    login: '登录', register: '注册', email: '邮箱', username: '用户名', password: '密码', preferredLanguage: '首选语言', chinese: '中文', english: 'English', createAccount: '创建账号', haveAccount: '已有账号？', noAccount: '还没有账号？', usernameHint: '可留空，默认使用邮箱。', registrationRequired: '请填写邮箱，密码至少 6 位。', registrationSuccess: '注册成功。', registrationFailed: '注册失败', loginSuccess: '已登录。', loginFailed: '登录失败', emailPasswordRequired: '请填写邮箱和密码。', rememberAccount: '记住账号', brand: '实习投递看板', admin: '管理端', board: '返回看板', logout: '退出', search: '按公司名或职位搜索', favorites: '只看喜欢', allApplications: '全部岗位', newApplication: '新建投递', loadFailed: '加载失败。',
  },
} as const

type MessageKey = keyof typeof messages.en
const guestLanguage = ref<PreferredLanguage>((localStorage.getItem('it_guest_language') === 'en' || navigator.language.toLowerCase().startsWith('en')) ? 'en' : 'zh')

export function useI18n() {
  const { user } = useAuth()
  const language = computed<PreferredLanguage>(() => user.value?.preferred_language ?? guestLanguage.value)
  const setGuestLanguage = (value: PreferredLanguage) => { guestLanguage.value = value; localStorage.setItem('it_guest_language', value) }
  const t = (key: MessageKey) => messages[language.value][key]
  const statusLabel = (status: ApplicationStatus) => ({
    en: { Wishlist: 'Wishlist', Applied: 'Applied', OA: 'Online assessment', Interview: 'Interview', Offer: 'Offer', Rejected: 'Rejected', Withdrawn: 'Withdrawn' },
    zh: { Wishlist: '想投', Applied: '已投', OA: '笔试', Interview: '面试', Offer: 'Offer', Rejected: '被拒', Withdrawn: '已放弃' },
  }[language.value][status])
  return { language, setGuestLanguage, t, statusLabel }
}

// Existing screens use static Chinese copy. Translate their rendered labels for English users
// while new UI uses the t() helper above.
export function installPageLocalization() {
  const { language } = useI18n()
  const english: Record<string, string> = {
    '管理端': 'Admin', '返回看板': 'Board', '退出': 'Sign out', '新建投递': 'New application', '全部岗位': 'All applications', '只看喜欢': 'Favorites only',
    '概览': 'Overview', '用户管理': 'User management', '投递巡查': 'Application inspection', '运营概览': 'Operations overview', '用户数': 'Users', '投递总数': 'Applications',
    '今日新增用户': 'New users today', '今日新增投递': 'New applications today', '查询': 'Search', '重置': 'Reset', '暂无数据': 'No data', '启用': 'Enabled', '禁用': 'Disabled',
    '公司名': 'Company', '职位名': 'Position', '状态': 'Status', '地点': 'Location', '薪酬': 'Salary', '优先级': 'Priority', '截止日期': 'Deadline', '笔记': 'Notes',
    '保存更改': 'Save changes', '删除投递': 'Delete application', '喜欢': 'Favorite', '可转正': 'Conversion opportunity', '登录': 'Sign in', '注册': 'Create account',
    '基本信息': 'Basic information', 'AI 匹配': 'AI match', '岗位信息': 'Job information', '公司名称': 'Company name', '岗位名称': 'Position title', '部门 / 团队': 'Department / Team', '用工类型': 'Employment type',
    '申请链接': 'Application link', '办公方式': 'Work arrangement', '岗位要求': 'Job requirements', '学历要求': 'Education requirement', '经验年限': 'Experience required(years)', '每周到岗天数': 'Days per week', '最短实习月数': 'Minimum internship length',
    '最低': 'Minimum', '最高': 'Maximum', '至': 'to', '硬性要求 / 技术栈': 'Required skills / Tech stack', '岗位职责': 'Responsibilities', '原始 JD': 'Original job description', '提供转正机会': 'Conversion opportunity',
    '投递与联系人': 'Application and contact', '来源类型': 'Source type', '来源名称': 'Source name', '实际投递时间': 'Application time', '联系人': 'Contact person', '联系人邮箱': 'Contact email', '来源链接': 'Source link',
    '个人记录': 'Personal notes', '使用的简历版本': 'Resume version used', '标记为喜欢': 'Mark as favorite', '上次更新': 'Last updated', '新建岗位': 'Create application', '删除岗位': 'Delete application',
    '实习': 'Internship', '全职': 'Full-time', '兼职': 'Part-time', '合同制': 'Contract', '公司官网': 'Company website', '校园招聘': 'Campus recruiting', '内推': 'Referral', '其他': 'Other', 'Boss直聘': 'Boss Zhipin',
    '不限': 'No requirement', '大专': 'Associate degree', '本科': "Bachelor's", '硕士': "Master's", '博士': 'Doctorate', '未设': 'Not set', '高': 'High', '中': 'Medium', '低': 'Low', '线下': 'On-site', '混合': 'Hybrid', '远程': 'Remote',
    '想投': 'Wishlist', '已投': 'Applied', '笔试': 'Online assessment', '面试': 'Interview', '被拒': 'Rejected', '已放弃': 'Withdrawn', '结束': 'Rejected', '主动结束流程': 'Application withdrawn', '还没提交': 'Not submitted', '材料已送出': 'Application submitted', '已拿到': 'Received',
    '投递生命周期': 'Application lifecycle', '沿着一条时间线回顾每个岗位从发现、投递到最终结果的完整过程。': 'Review every application from discovery through the final outcome on one timeline.', '全部事件': 'All events', '新增岗位': 'Application added', '状态变化': 'Status change', '投递': 'Applied', '主动放弃': 'Withdrawn', '备忘录': 'Memo',
    '最新优先': 'Newest first', '最早优先': 'Oldest first', '查看岗位详情 →': 'View application →', '还没有符合条件的生命周期记录': 'No lifecycle events match these filters',
    '记录状态变化': 'Record status change', '发生时间（可选）': 'Time (optional)', '薪资和回复期限': 'Compensation and response deadline', '面试轮次 / 方式': 'Interview round / format', '笔试形式 / 重点': 'Assessment format / focus', '原因': 'Reason', '相关信息': 'Details', '备忘录 / 复盘（可选）': 'Memo / reflection (optional)', '取消': 'Cancel', '跳过备忘并更新': 'Skip memo and update', '保存并更新': 'Save and update',
    '只看喜爱': 'Favorites only', '没有匹配的卡片': 'No matching applications', '这一列还是空的': 'This column is empty', '松开放到这一列': 'Drop into this column', '换个公司名试试': 'Try another company name', '拖卡片进来，或点右上角新建': 'Drag a card here or create a new application',
    '年限不限': 'No experience requirement', '天/周': 'days/week', '已过期': 'Overdue by', '剩余': 'Remaining', '天': 'days', '查看岗位': 'View posting', '设置': 'Settings', '生命周期': 'Lifecycle', '看板': 'Board',
    '搜索公司、岗位或记录': 'Search company, position, or notes', '请选择或输入': 'Select or enter', '邮箱格式不正确': 'Invalid email format', '例如 上海 / Remote': 'e.g. Shanghai / Remote', '例如校招公众号 / 推荐人': 'e.g. campus channel / referrer',
    '添加了': 'Added', '未知公司': 'Unknown company', '未知岗位': 'Unknown position', '从': 'from', '进入': 'moved to', '线上作业 / OA': 'Online assignment / OA', '电话 / 视频 / onsite': 'Phone / video / on-site',
    '例如 Stripe': 'e.g. Stripe', '例如 Infrastructure': 'e.g. Infrastructure', '例如 Python, PyTorch': 'e.g. Python, PyTorch', '例如 resume-v3-systems.pdf': 'e.g. resume-v3-systems.pdf',
    '按公司名或职位搜索': 'Search by company or position', '例如 Software Engineer Intern': 'e.g. Software Engineer Intern',
    '公司名称和岗位名称为必填项': 'Company name and position title are required', '链接需以 http:// 或 https:// 开头': 'Links must begin with http:// or https://', '请输入有效的联系人邮箱': 'Enter a valid contact email address',
    '普通备忘': 'General note', '已投递': 'Applied', '编辑事件': 'Edit event', '添加生命周期事件': 'Add lifecycle event', '事件类型': 'Event type', '发生时间': 'Date and time', '标题': 'Title', '轮次': 'Round',
    '例如：技术一面': 'e.g. First technical interview', '例如：第二轮': 'e.g. Second round', '保存事件': 'Save event', '添加事件': 'Add event', '取消编辑': 'Cancel editing', '编辑': 'Edit', '删除': 'Delete',
    'AI 匹配分析': 'AI match analysis', '将在个人简历模块完成后接入。届时可选择简历版本，与当前岗位 JD 生成匹配报告。': 'Available after the resume module is complete. You will be able to compare a resume version with the current job description.',
    '按用户筛选': 'Filter by user', '公司、职位或用户名': 'Company, position, or username', '公司': 'Company', '职位': 'Position', '更新时间': 'Updated at', '创建时间': 'Created at',
    '搜索后点查询。禁用后该用户不能登录。': 'Enter a search and select Search. Disabled users cannot sign in.',
    '请填写事件标题': 'Enter an event title', '事件已更新': 'Event updated', '事件已添加': 'Event added', '保存事件失败': 'Unable to save event', '删除事件': 'Delete event', '删除事件失败': 'Unable to delete event',
    '已保存更改': 'Changes saved', '已新建岗位': 'Application created', '保存失败': 'Unable to save', '确认删除': 'Confirm deletion', '确认删除岗位': 'Delete application', '数据会安全保留，但不再显示在看板中。': 'The data will be retained safely but will no longer appear on the board.', '岗位已从看板移除': 'Application removed from the board', '删除失败': 'Unable to delete',
    '形式 / 地点': 'Format / Location', '线上 / 现场 / HackerRank': 'Online / On-site / HackerRank', '结果': 'Result', '薪资 / 回复期限': 'Compensation / Response deadline', '记录 / 复盘': 'Notes / Reflection', '还没有生命周期记录': 'No lifecycle events yet',
    '邮箱或用户名': 'Email or username', '用户 ID': 'User ID', '角色': 'Role', '操作': 'Actions', '总计': 'Total', '条': 'records', '已启用': 'User enabled', '已禁用': 'User disabled', '更新失败': 'Unable to update',
    '只读查看全站投递，不在后台改状态。': 'Review all applications in read-only mode. Statuses cannot be changed here.', '用户': 'User', '年限': 'Experience', '学历': 'Education', '办公': 'Work mode', '喜爱': 'Favorite', '是': 'Yes', '否': 'No', '简历版本': 'Resume version',
  }
  const originalText = new WeakMap<Text, string>()
  const originalAttributes = new WeakMap<HTMLElement, Record<string, string>>()
  const translate = (source: string) => {
    const trimmed = source.trim()
    if (english[trimmed]) return source.replace(trimmed, english[trimmed])
    let result = source
    Object.entries(english).filter(([zh]) => zh.length >= 2).sort((a, b) => b[0].length - a[0].length).forEach(([zh, en]) => { result = result.split(zh).join(en) })
    return result
  }
  const apply = () => nextTick(() => {
    const walker = document.createTreeWalker(document.body, NodeFilter.SHOW_TEXT)
    const nodes: Text[] = []; let node: Node | null
    while ((node = walker.nextNode())) nodes.push(node as Text)
    nodes.forEach((text) => {
      if (!originalText.has(text)) originalText.set(text, text.nodeValue || '')
      let original = originalText.get(text) || ''
      const current = text.nodeValue || ''
      // Vue and Element Plus reuse text nodes (notably calendar day cells). If the
      // component changed the source text, refresh our cache instead of restoring
      // a stale value from the previous render/month.
      if (current !== original && current !== translate(original)) {
        original = current
        originalText.set(text, current)
      }
      const value = language.value === 'en' ? translate(original) : original
      if (value !== text.nodeValue) text.nodeValue = value
    })
    document.querySelectorAll<HTMLElement>('[placeholder],[title],[aria-label]').forEach((element) => {
      if (!originalAttributes.has(element)) originalAttributes.set(element, {})
      const originals = originalAttributes.get(element)!
      for (const attribute of ['placeholder', 'title', 'aria-label']) {
        const current = element.getAttribute(attribute); if (!current) continue
        if (!(attribute in originals)) originals[attribute] = current
        else if (current !== originals[attribute] && current !== translate(originals[attribute])) originals[attribute] = current
        const value = language.value === 'en' ? translate(originals[attribute]) : originals[attribute]
        if (value !== current) element.setAttribute(attribute, value)
      }
    })
  })
  watch(language, apply, { immediate: true })
  new MutationObserver(apply).observe(document.body, { childList: true, subtree: true })
}
