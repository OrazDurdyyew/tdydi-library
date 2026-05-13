<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { admin, fetchMe } = useAuth()
if (!admin.value) await fetchMe()

useHead({ title: 'Дашборд' })
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-2">Дашборд</h1>
    <p v-if="admin" class="text-muted-foreground mb-8">
      Добро пожаловать, <strong class="text-foreground">{{ admin.username }}</strong>
      <span class="text-sm">· {{ admin.role }}</span>
    </p>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-4xl">
      <NuxtLink to="/admin/sections" class="card card-hover group">
        <div class="text-3xl mb-3">📚</div>
        <div class="font-semibold text-lg group-hover:text-primary transition-colors">Разделы</div>
        <div class="text-sm text-muted-foreground mt-1">Управление меню сайта и кафедрами</div>
      </NuxtLink>

      <NuxtLink to="/admin/books" class="card card-hover group">
        <div class="text-3xl mb-3">📖</div>
        <div class="font-semibold text-lg group-hover:text-primary transition-colors">Книги</div>
        <div class="text-sm text-muted-foreground mt-1">Каталог книг и загрузка файлов</div>
      </NuxtLink>
    </div>
  </div>
</template>
