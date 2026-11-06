<script setup lang="ts">
import { computed } from 'vue'
import en from 'element-plus/es/locale/lang/en'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { installPageLocalization } from './i18n'
import { useI18n } from './i18n'
installPageLocalization()
const { language } = useI18n()
const elementLocale = computed(() => language.value === 'en' ? en : zhCn)
</script>

<template>
  <el-config-provider :locale="elementLocale"><router-view /></el-config-provider>
</template>
