<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()

useHead({ title: () => t('admin.home.title') })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-2">{{ t('admin.home.title') }}</h1>
    <p class="text-sm text-muted-foreground mb-8">
      {{ t('admin.home.description') }}
    </p>

    <div class="space-y-6 max-w-4xl">
      <HomeBannerManager
        endpoint="/admin/home/slides"
        :title="t('admin.home.slides')"
        :description="t('admin.home.slides_description')"
      />

      <HomeBannerManager
        endpoint="/admin/home/promos"
        :title="t('admin.home.promos')"
        :description="t('admin.home.promos_description')"
      />

      <HomeRefManager
        endpoint="/admin/home/sections"
        source-endpoint="/sections"
        ref-field="section_id"
        nested-key="section"
        :title="t('admin.home.section_blocks')"
        :description="t('admin.home.section_blocks_description')"
      />

      <HomeRefManager
        endpoint="/admin/home/collections"
        source-endpoint="/collections"
        ref-field="collection_id"
        nested-key="collection"
        :title="t('admin.home.home_collections')"
        :description="t('admin.home.home_collections_description')"
      />
    </div>
  </div>
</template>
