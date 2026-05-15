<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

interface Admin {
  id: number
  username: string
  email: string
  role: string
  all_sections: boolean
  created_at: string
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()
const { admin: me } = useAuth()

const { data, refresh, pending } = await useAsyncData('admin-admins', () =>
  api.get<Admin[]>('/admin/admins')
)

const admins = computed(() => data.value ?? [])

async function remove(id: number, username: string) {
  const ok = await ask({
    title: 'Удалить администратора?',
    description: `Аккаунт «${username}» будет удалён.`,
    confirmText: 'Удалить',
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/admins/${id}`)
    await refresh()
    toast.success('Администратор удалён')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка удаления')
  }
}

const roleLabel = (r: string) => (r === 'superadmin' ? 'Суперадмин' : 'Контент-менеджер')

useHead({ title: 'Администраторы' })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">Администраторы</h1>
      <NuxtLink to="/admin/admins/new">
        <Button><Plus class="w-4 h-4" /> Добавить</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">Загрузка…</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Логин</TableHead>
            <TableHead>Email</TableHead>
            <TableHead>Роль</TableHead>
            <TableHead>Разделы</TableHead>
            <TableHead class="text-right">Действия</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="a in admins" :key="a.id">
            <TableCell class="font-mono">
              <div class="flex items-center gap-2">
                <span class="font-medium">{{ a.username }}</span>
                <Badge v-if="me?.id === a.id" variant="success">вы</Badge>
              </div>
            </TableCell>
            <TableCell class="text-muted-foreground">{{ a.email || '—' }}</TableCell>
            <TableCell>
              <Badge :variant="a.role === 'superadmin' ? 'warning' : 'secondary'">
                {{ roleLabel(a.role) }}
              </Badge>
            </TableCell>
            <TableCell class="text-muted-foreground">
              {{ a.all_sections ? 'Все разделы' : 'Выборочные' }}
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/admins/${a.id}`">
                  <Button variant="ghost" size="sm">Изменить</Button>
                </NuxtLink>
                <Button
                  v-if="me?.id !== a.id"
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(a.id, a.username)"
                >
                  Удалить
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="admins.length === 0" :colspan="5">Нет данных</TableEmpty>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>
