---
name: new-admin-page
description: Создать новую страницу админки в Nuxt 4. Триггер — когда пользователь просит добавить раздел в админку (например, "управление авторами", "страница настроек", "форма редактирования книги"). Создаёт страницу с auth-проверкой, admin layout, типизированным API-клиентом, формами через shadcn-vue.
---

# new-admin-page

Создаёт страницу в админ-панели Nuxt 4 (под `/admin/*`).

## Когда использовать

- "добавь страницу в админку для X"
- "нужна форма управления Y в админке"
- "создай раздел админки ..."

## Что создать

1. **Страница** в `/frontend/pages/admin/{section}/[action].vue` (или `index.vue` для списка):
   - `definePageMeta({ layout: 'admin', middleware: ['auth'] })`
   - Использует `<script setup lang="ts">`
   - Использует компоненты shadcn-vue (Button, Input, Dialog, DropdownMenu, DataTable и т.д.)
2. **Composable** в `/frontend/composables/use{Entity}.ts`:
   - `useFetch`/`$fetch` обёртка для API
   - Типы для request/response из `/shared/types/`
3. **Pinia store** (если данные нужны глобально) в `/frontend/stores/{entity}.ts`
4. **i18n ключи** во всех 3 файлах локали (`/frontend/locales/tm.json`, `ru.json`, `en.json`):
   - Префикс `admin.{section}.{key}`
   - Запустить `/i18n-add-key` для добавления

## Стандартные паттерны страницы списка

```vue
<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: ['auth'] })
const { data, pending, refresh } = await useFetch('/api/admin/...')
// фильтры, сортировка, пагинация
</script>

<template>
  <div>
    <PageHeader :title="$t('...')" />
    <Filters />
    <DataTable :columns :data />
    <Pagination />
  </div>
</template>
```

## Стандартные паттерны формы

- Использовать `vee-validate` + `zod` для валидации
- Использовать shadcn-vue Form компоненты
- Поля на 3 языках — tabs (TM/RU/EN) или 3 input в столбик
- Файлы — drag-and-drop зона с прогресс-баром (через `tus-js-client`)

## Проверка прав на фронте

- `auth` middleware проверяет JWT и наличие пользователя
- Для страниц только-суперадмина — добавить redirect, если `user.role !== 'superadmin'`
- Для страниц с книгами — отфильтровать разделы по `user.allowed_sections` (приходят с бэкенда при логине)

## Соглашения

- Имена файлов — kebab-case (`book-form.vue`, `section-tree.vue`)
- Компоненты — PascalCase в шаблонах
- API client использует автогенерированные типы из бэкенда (или вручную в `/shared/types/`)
