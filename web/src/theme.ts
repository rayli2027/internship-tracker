import { computed, ref } from 'vue'

export type ColorMode = 'light' | 'dark'
const STORAGE_KEY = 'it_color_mode'
const mode = ref<ColorMode>(localStorage.getItem(STORAGE_KEY) === 'dark' ? 'dark' : 'light')

function applyMode(value: ColorMode) {
  document.documentElement.dataset.theme = value
  document.documentElement.style.colorScheme = value
}

applyMode(mode.value)

export function useTheme() {
  const isDark = computed(() => mode.value === 'dark')
  function setMode(value: ColorMode) {
    mode.value = value
    localStorage.setItem(STORAGE_KEY, value)
    applyMode(value)
  }
  function toggleMode() { setMode(isDark.value ? 'light' : 'dark') }
  return { mode, isDark, setMode, toggleMode }
}
