---
name: new-public-page
description: Создать публичную страницу сайта в Nuxt 4 с SSR. Триггер — когда пользователь просит добавить страницу для посетителей (например, страница раздела, страница книги, страница новости, статичная страница). Создаёт страницу с SSR, i18n, SEO meta, breadcrumbs.
---

# new-public-page

Создаёт публичную страницу сайта в Nuxt 4 с SSR-рендерингом.

## Когда использовать

- "добавь публичную страницу X"
- "нужна страница раздела/книги/новости/..."
- "создай SEO-страницу для ..."

## Что создать

1. **Страница** в `/frontend/pages/...vue`:
   - `definePageMeta({ layout: 'default' })` (или без — default подразумевается)
   - НЕ ставить `ssr: false` — SEO важен
   - Использует `useAsyncData`/`useFetch` для SSR-fetch данных
   - Использует `useHead`/`useSeoMeta` для meta-тегов
2. **i18n** во всех 3 локалях
3. **SEO-meta**:
   ```ts
   useSeoMeta({
     title: () => data.value?.title,
     description: () => data.value?.annotation,
     ogTitle: () => data.value?.title,
     ogImage: () => data.value?.cover_url,
     ogType: 'article'
   })
   ```
4. **Breadcrumbs**: компонент `<Breadcrumbs />` с массивом сегментов
5. **JSON-LD** для seo (для книг — schema.org/Book, для новостей — Article)

## URL-структура

Все публичные URL начинаются с языкового префикса: `/tm/...`, `/ru/...`, `/en/...`.

Примеры:
- `/{lang}/sections/{slug}` — страница раздела (любого уровня)
- `/{lang}/books/{slug}` — страница книги
- `/{lang}/news/{slug}` — страница новости
- `/{lang}/periodicals/{slug}` — страница издания
- `/{lang}/periodicals/{slug}/issues/{issue-slug}` — страница выпуска

## Шаблонные элементы

- `<Container>` для контейнера с max-width
- `<Breadcrumbs>` хлебные крошки
- `<SectionTitle>` заголовок
- `<BookCard>`, `<SectionCard>`, `<NewsCard>` — карточки
- `<Pagination>` для списков

## Брендинг

CSS-переменные брендинга применяются автоматически через `app.vue` (загружаются из `/api/site-settings` при старте). Компоненты shadcn-vue используют эти переменные — отдельно ничего настраивать не нужно.

## Адаптивность

- Использовать Tailwind responsive prefixes (`md:`, `lg:`)
- Мобильный — 320px минимум
- Все картинки — через `<NuxtImg>` с `loading="lazy"` (кроме hero)
