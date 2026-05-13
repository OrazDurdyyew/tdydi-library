<script setup lang="ts">
import { Button } from '~/components/ui/button'
import { Input } from '~/components/ui/input'
import { Label } from '~/components/ui/label'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '~/components/ui/card'

definePageMeta({ layout: false })

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const { login } = useAuth()
const localePath = useLocalePath()

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await login(username.value, password.value)
    await navigateTo(localePath('/admin'))
  } catch (e: any) {
    error.value = e.data?.error || 'Ошибка входа'
  } finally {
    loading.value = false
  }
}

useHead({ title: 'Вход в админку' })
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background p-4">
    <Card class="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Добро пожаловать</CardTitle>
        <CardDescription>Войдите в админ-панель</CardDescription>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="submit" class="space-y-5">
          <div class="space-y-2">
            <Label for="username">Логин</Label>
            <Input
              id="username"
              v-model="username"
              type="text"
              required
              autocomplete="username"
              placeholder="admin"
            />
          </div>
          <div class="space-y-2">
            <Label for="password">Пароль</Label>
            <Input
              id="password"
              v-model="password"
              type="password"
              required
              autocomplete="current-password"
              placeholder="••••••••"
            />
          </div>
          <div
            v-if="error"
            class="text-sm text-destructive px-3 py-2 bg-destructive/10 rounded-lg"
          >
            {{ error }}
          </div>
          <Button type="submit" :disabled="loading" class="w-full">
            {{ loading ? 'Вход…' : 'Войти' }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
