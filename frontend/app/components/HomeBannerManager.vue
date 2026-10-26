<script setup lang="ts">
import { Plus, Upload, Trash2, GripVertical } from 'lucide-vue-next'

// Универсальный менеджер баннеров главной — обслуживает и слайдер, и промо-блоки
// (одинаковая структура полей). endpoint: /admin/home/slides | /admin/home/promos.
const props = defineProps<{
  endpoint: string
  title: string
  description?: string
}>()

const { t } = useI18n()
const api = useApi()
const toast = useToast()
const { ask } = useConfirm()
const fileUrl = useFileUrl()

interface Banner {
  id: number
  image_path: string
  link_url: string
  sort_order: number
  is_active: boolean
  title_tm: string
  title_ru: string
  title_en: string
  subtitle_tm: string
  subtitle_ru: string
  subtitle_en: string
}

const { data, refresh, pending } = await useAsyncData(`banners-${props.endpoint}`, () =>
  api.get<Banner[]>(props.endpoint).catch(() => [])
)
const items = computed(() => data.value ?? [])

const busy = ref<number | 'new' | null>(null)

async function addItem() {
  busy.value = 'new'
  try {
    await api.post(props.endpoint, {
      is_active: true,
      sort_order: items.value.length,
      link_url: '',
      title_tm: '',
      title_ru: '',
      title_en: '',
      subtitle_tm: '',
      subtitle_ru: '',
      subtitle_en: ''
    })
    await refresh()
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error'))
  } finally {
    busy.value = null
  }
}

async function save(item: Banner) {
  busy.value = item.id
  try {
    await api.put(`${props.endpoint}/${item.id}`, {
      is_active: item.is_active,
      sort_order: item.sort_order,
      link_url: item.link_url,
      title_tm: item.title_tm,
      title_ru: item.title_ru,
      title_en: item.title_en,
      subtitle_tm: item.subtitle_tm,
      subtitle_ru: item.subtitle_ru,
      subtitle_en: item.subtitle_en
    })
    toast.success(t('admin.common.saved'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error'))
  } finally {
    busy.value = null
  }
}

async function remove(item: Banner) {
  const ok = await ask({
    title: t('admin.home.confirm_delete_block_title'),
    description: t('admin.home.confirm_delete_block_description'),
    confirmText: t('admin.common.delete'),
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`${props.endpoint}/${item.id}`)
    await refresh()
    toast.success(t('admin.home.block_deleted'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error'))
  }
}

async function uploadImage(item: Banner, e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  busy.value = item.id
  try {
    await api.uploadField(`${props.endpoint}/${item.id}/image`, 'image', file)
    await refresh()
    toast.success(t('admin.home.image_uploaded'))
  } catch (err: any) {
    toast.error(err.data?.error || t('admin.common.error_upload'))
  } finally {
    busy.value = null
    input.value = ''
  }
}
</script>

<template>
  <Card>
    <CardHeader class="flex-row items-center justify-between space-y-0">
      <div>
        <CardTitle class="text-lg">{{ title }}</CardTitle>
        <CardDescription v-if="description">{{ description }}</CardDescription>
      </div>
      <Button type="button" size="sm" :disabled="busy === 'new'" @click="addItem">
        <Plus class="w-4 h-4" /> {{ t('admin.common.add') }}
      </Button>
    </CardHeader>
    <CardContent class="space-y-4">
      <p v-if="pending" class="text-sm text-muted-foreground">{{ t('admin.common.loading') }}</p>
      <p v-else-if="!items.length" class="text-sm text-muted-foreground">{{ t('admin.common.no_data') }}</p>

      <div
        v-for="item in items"
        :key="item.id"
        class="rounded-lg border border-border p-4 space-y-4"
      >
        <div class="flex gap-4">
          <!-- Картинка -->
          <div class="w-40 shrink-0 space-y-2">
            <div class="aspect-video rounded-md border border-border bg-muted overflow-hidden flex items-center justify-center">
              <img
                v-if="item.image_path"
                :src="fileUrl(item.image_path)"
                alt=""
                class="w-full h-full object-cover"
              />
              <Upload v-else class="w-6 h-6 text-muted-foreground" />
            </div>
            <label class="btn-outline btn-sm w-full cursor-pointer justify-center">
              <Upload class="w-3.5 h-3.5" /> {{ t('admin.common.image') }}
              <input type="file" accept="image/*" class="hidden" @change="uploadImage(item, $event)" />
            </label>
          </div>

          <!-- Поля -->
          <div class="flex-1 min-w-0 space-y-3">
            <div class="grid grid-cols-3 gap-2">
              <div>
                <label class="text-xs text-muted-foreground">{{ t('admin.common.title') }} {{ t('admin.common.lang_tm') }}</label>
                <Input v-model="item.title_tm" class="mt-1" />
              </div>
              <div>
                <label class="text-xs text-muted-foreground">{{ t('admin.common.lang_ru') }}</label>
                <Input v-model="item.title_ru" class="mt-1" />
              </div>
              <div>
                <label class="text-xs text-muted-foreground">{{ t('admin.common.lang_en') }}</label>
                <Input v-model="item.title_en" class="mt-1" />
              </div>
            </div>
            <div class="grid grid-cols-3 gap-2">
              <div>
                <label class="text-xs text-muted-foreground">{{ t('admin.home.subtitle') }} {{ t('admin.common.lang_tm') }}</label>
                <Input v-model="item.subtitle_tm" class="mt-1" />
              </div>
              <div>
                <label class="text-xs text-muted-foreground">{{ t('admin.common.lang_ru') }}</label>
                <Input v-model="item.subtitle_ru" class="mt-1" />
              </div>
              <div>
                <label class="text-xs text-muted-foreground">{{ t('admin.common.lang_en') }}</label>
                <Input v-model="item.subtitle_en" class="mt-1" />
              </div>
            </div>
            <div class="grid grid-cols-[1fr_auto_auto] gap-3 items-end">
              <div>
                <label class="text-xs text-muted-foreground">{{ t('admin.home.link_optional') }}</label>
                <Input v-model="item.link_url" class="mt-1 font-mono text-xs" placeholder="/sections/economics" />
              </div>
              <div>
                <label class="text-xs text-muted-foreground">{{ t('admin.common.sort_order') }}</label>
                <Input v-model.number="item.sort_order" type="number" class="mt-1 w-20" />
              </div>
              <label class="flex items-center gap-2 h-10 cursor-pointer">
                <Checkbox v-model="item.is_active" />
                <span class="text-sm">{{ t('admin.common.active') }}</span>
              </label>
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-2">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            class="text-destructive hover:text-destructive"
            @click="remove(item)"
          >
            <Trash2 class="w-4 h-4" /> {{ t('admin.common.delete') }}
          </Button>
          <Button type="button" size="sm" :disabled="busy === item.id" @click="save(item)">
            {{ t('admin.common.save') }}
          </Button>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
