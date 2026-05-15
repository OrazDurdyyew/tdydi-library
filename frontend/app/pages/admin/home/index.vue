<script setup lang="ts">
definePageMeta({ layout: 'admin' })

useHead({ title: 'Главная страница' })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-2">Главная страница</h1>
    <p class="text-sm text-muted-foreground mb-8">
      Управление баннерами, промо-блоками и блоками разделов/подборок на главной.
    </p>

    <div class="space-y-6 max-w-4xl">
      <HomeBannerManager
        endpoint="/admin/home/slides"
        title="Баннер-карусель"
        description="Слайды вверху главной страницы. Несколько слайдов листаются автоматически."
      />

      <HomeBannerManager
        endpoint="/admin/home/promos"
        title="Промо-блоки"
        description="Блоки «картинка + заголовок + ссылка» под баннером."
      />

      <HomeRefManager
        endpoint="/admin/home/sections"
        source-endpoint="/sections"
        ref-field="section_id"
        nested-key="section"
        title="Блоки разделов на главной"
        description="Разделы, чьи книги показываются на главной отдельным блоком (с баннером раздела, если у него есть картинка)."
      />

      <HomeRefManager
        endpoint="/admin/home/collections"
        source-endpoint="/collections"
        ref-field="collection_id"
        nested-key="collection"
        title="Блоки подборок на главной"
        description="Виртуальные подборки, отображаемые на главной странице."
      />
    </div>
  </div>
</template>
