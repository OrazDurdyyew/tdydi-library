<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const api = useApi()
const toast = useToast()
const localePath = useLocalePath()

const { data: sectionsResp } = await useAsyncData('sections-for-admin', () =>
  api.get<{ data: any[] }>('/sections')
)
const sections = computed(() => sectionsResp.value?.data ?? [])

const form = reactive({
  username: '',
  password: '',
  email: '',
  role: 'content_manager',
  all_sections: true,
  section_ids: [] as number[],
})

const saving = ref(false)

async function submit() {
  if (!form.username || !form.password) {
    toast.error('Логин и пароль обязательны')
    return
  }
  saving.value = true
  try {
    await api.post('/admin/admins', form)
    toast.success('Администратор создан')
    await navigateTo(localePath('/admin/admins'))
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка сохранения')
  } finally {
    saving.value = false
  }
}

useHead({ title: 'Добавить администратора' })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/admins" class="text-muted-foreground hover:text-foreground text-sm">← Администраторы</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif">Добавить администратора</h1>
    </div>

    <form @submit.prevent="submit" class="space-y-6 max-w-xl">
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Данные аккаунта</CardTitle>
        </CardHeader>
        <CardContent class="space-y-5">
          <div class="space-y-2">
            <Label for="username">Логин <span class="text-destructive">*</span></Label>
            <Input id="username" v-model="form.username" required class="font-mono" autocomplete="off" />
          </div>
          <div class="space-y-2">
            <Label for="password">Пароль <span class="text-destructive">*</span></Label>
            <Input id="password" v-model="form.password" required type="password" autocomplete="new-password" />
          </div>
          <div class="space-y-2">
            <Label for="email">Email</Label>
            <Input id="email" v-model="form.email" type="email" />
          </div>
          <div class="space-y-2">
            <Label for="role">Роль</Label>
            <Select id="role" v-model="form.role">
              <option value="content_manager">Контент-менеджер</option>
              <option value="superadmin">Суперадмин</option>
            </Select>
          </div>
        </CardContent>
      </Card>

      <Card v-if="form.role === 'content_manager'">
        <CardHeader>
          <CardTitle class="text-lg">Доступ к разделам</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <label class="flex items-center gap-3 cursor-pointer">
            <Checkbox v-model="form.all_sections" />
            <span class="text-sm">Все разделы</span>
          </label>
          <div v-if="!form.all_sections" class="border border-border rounded-lg p-3 max-h-56 overflow-y-auto space-y-1.5">
            <label
              v-for="s in sections"
              :key="s.id"
              class="flex items-center gap-3 cursor-pointer hover:bg-muted px-2 py-1 rounded-md"
            >
              <Checkbox :value="s.id" v-model="form.section_ids" />
              <span :class="s.parent_id === null ? 'font-medium text-sm' : 'text-sm pl-4 text-muted-foreground'">
                {{ s.name_ru || s.name_tm }}
              </span>
            </label>
          </div>
        </CardContent>
      </Card>

      <div class="flex gap-3">
        <Button type="submit" :disabled="saving">
          {{ saving ? 'Сохранение…' : 'Создать' }}
        </Button>
        <NuxtLink to="/admin/admins">
          <Button type="button" variant="outline">Отмена</Button>
        </NuxtLink>
      </div>
    </form>
  </div>
</template>
