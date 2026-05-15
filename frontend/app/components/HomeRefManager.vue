<script setup lang="ts">
import { Plus, Trash2 } from 'lucide-vue-next'

// Менеджер блоков-ссылок главной: «какие разделы / подборки показывать на главной».
const props = defineProps<{
  endpoint: string // /admin/home/sections | /admin/home/collections
  sourceEndpoint: string // /sections | /collections
  refField: string // section_id | collection_id
  nestedKey: string // section | collection
  title: string
  description?: string
}>()

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()
const loc = useLocalized()

const { data, refresh } = await useAsyncData(`ref-${props.endpoint}`, () =>
  api.get<any[]>(props.endpoint).catch(() => [])
)
const { data: srcData } = await useAsyncData(`refsrc-${props.endpoint}`, () =>
  api.get<{ data: any[] }>(props.sourceEndpoint).catch(() => ({ data: [] }))
)

const items = computed(() => data.value ?? [])
const sources = computed(() => srcData.value?.data ?? [])
const usedIds = computed(() => new Set(items.value.map((i) => i[props.refField])))
const available = computed(() => sources.value.filter((s) => !usedIds.value.has(s.id)))

const picked = ref<string>('')

async function add() {
  if (!picked.value) return
  try {
    await api.post(props.endpoint, {
      [props.refField]: Number(picked.value),
      sort_order: items.value.length,
      is_active: true
    })
    picked.value = ''
    await refresh()
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  }
}

async function toggle(item: any) {
  try {
    await api.put(`${props.endpoint}/${item.id}`, {
      sort_order: item.sort_order,
      is_active: !item.is_active
    })
    await refresh()
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  }
}

async function remove(item: any) {
  const ok = await ask({
    title: 'Убрать с главной?',
    description: 'Блок будет убран с главной страницы.',
    confirmText: 'Убрать',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`${props.endpoint}/${item.id}`)
    await refresh()
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка')
  }
}
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-lg">{{ title }}</CardTitle>
      <CardDescription v-if="description">{{ description }}</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <!-- Добавить -->
      <div class="flex gap-2">
        <select v-model="picked" class="input flex-1">
          <option value="">— выбрать —</option>
          <option v-for="s in available" :key="s.id" :value="String(s.id)">
            {{ loc(s, 'name') }}
          </option>
        </select>
        <Button type="button" size="sm" :disabled="!picked" @click="add">
          <Plus class="w-4 h-4" /> Добавить
        </Button>
      </div>

      <!-- Список -->
      <p v-if="!items.length" class="text-sm text-muted-foreground">Пока пусто.</p>
      <ul v-else class="divide-y divide-border">
        <li v-for="item in items" :key="item.id" class="flex items-center gap-3 py-2.5">
          <span class="flex-1 font-medium text-sm">
            {{ item[nestedKey] ? loc(item[nestedKey], 'name') : '—' }}
          </span>
          <label class="flex items-center gap-2 cursor-pointer">
            <Checkbox :model-value="item.is_active" @update:model-value="toggle(item)" />
            <span class="text-xs text-muted-foreground">{{ item.is_active ? 'Активен' : 'Скрыт' }}</span>
          </label>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            class="text-destructive hover:text-destructive"
            @click="remove(item)"
          >
            <Trash2 class="w-4 h-4" />
          </Button>
        </li>
      </ul>
    </CardContent>
  </Card>
</template>
