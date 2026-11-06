import { describe, expect, it } from 'vitest'
import { calendarKind, localDateKey, statusOptionLabel } from './planner'

describe('planner regression checks', () => {
  it('keeps the clicked local calendar date unchanged', () => {
    expect(localDateKey(new Date(2026,8,1,12))).toBe('2026-09-01')
    expect(localDateKey(new Date(2026,8,30,12))).toBe('2026-09-30')
  })
  it('synchronizes scheduled workflow events with calendar categories', () => {
    expect(calendarKind({event_type:'interview'})).toBe('interview')
    expect(calendarKind({event_type:'oa'})).toBe('workflow')
    expect(calendarKind({event_type:'offer'})).toBe('workflow')
    expect(calendarKind({event_type:'created'})).toBeNull()
  })
  it('does not duplicate bilingual labels in English', () => {
    expect(statusOptionLabel('Wishlist','en','想投')).toBe('Wishlist')
    expect(statusOptionLabel('OA','en','笔试')).toBe('Online assessment')
    expect(statusOptionLabel('Wishlist','zh','想投')).toBe('Wishlist · 想投')
  })
})
