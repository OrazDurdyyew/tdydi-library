// Транслитерация совпадает с api/internal/handlers/utils.go slugify():
// после смены логики на фронте — синхронизировать бэк, иначе бэк нормализует slug
// и фронтовая подсказка разойдётся с тем, что реально сохранится.
const TRANSLIT: Record<string, string> = {
  а: 'a', б: 'b', в: 'v', г: 'g', д: 'd', е: 'e', ё: 'yo',
  ж: 'zh', з: 'z', и: 'i', й: 'j', к: 'k', л: 'l', м: 'm',
  н: 'n', о: 'o', п: 'p', р: 'r', с: 's', т: 't', у: 'u',
  ф: 'f', х: 'kh', ц: 'ts', ч: 'ch', ш: 'sh', щ: 'shch',
  ъ: '', ы: 'y', ь: '', э: 'e', ю: 'yu', я: 'ya',
  ä: 'ae', ň: 'ng', ö: 'o', ü: 'u', ý: 'y', ş: 'sh', ž: 'zh',
  á: 'a', à: 'a', â: 'a', ã: 'a', å: 'a',
  é: 'e', è: 'e', ê: 'e', ë: 'e',
  í: 'i', ì: 'i', î: 'i', ï: 'i',
  ó: 'o', ò: 'o', ô: 'o', õ: 'o', ø: 'o',
  ú: 'u', ù: 'u', û: 'u',
  ñ: 'n', ç: 'c', ß: 'ss'
}

export function slugify(input: string): string {
  if (!input) return ''
  let out = ''
  for (const ch of input.toLowerCase()) {
    out += ch in TRANSLIT ? TRANSLIT[ch] : ch
  }
  return out.replace(/[^a-z0-9]+/g, '-').replace(/-{2,}/g, '-').replace(/^-|-$/g, '')
}

// useAutoSlug — связывает поле-источник (название) и поле slug.
// Авто-режим работает пока target пустой ИЛИ совпадает со slug(прошлого source).
// Как только пользователь правит slug руками — автогенерация отключается до тех пор,
// пока slug снова не станет пустым или не совпадёт с автогенерацией.
export function useAutoSlug(source: Ref<string>, target: Ref<string>) {
  watch(source, (val, oldVal) => {
    if (target.value === '' || target.value === slugify(oldVal || '')) {
      target.value = slugify(val)
    }
  })
}
