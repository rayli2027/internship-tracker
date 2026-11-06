<script setup lang="ts">
import AppHeader from '../../components/AppHeader.vue'
</script>

<template>
  <div class="app-shell">
    <AppHeader />
    <div class="admin-layout">
      <aside class="admin-nav">
        <p class="eyebrow">管理</p>
        <router-link to="/admin/overview">概览</router-link>
        <router-link to="/admin/users">用户管理</router-link>
        <router-link to="/admin/applications">投递巡查</router-link>
        <router-link to="/">返回看板</router-link>
      </aside>
      <section class="admin-main">
        <router-view />
      </section>
    </div>
  </div>
</template>
