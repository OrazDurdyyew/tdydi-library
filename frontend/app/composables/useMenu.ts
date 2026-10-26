export interface MenuSection {
  id: number
  slug: string
  name_tm: string
  name_ru: string
  name_en: string
  children?: MenuSection[]
}

export interface MenuPage {
  slug: string
  title_tm: string
  title_ru: string
  title_en: string
}

// Общий запрос меню сайта. Ключ useAsyncData один для всех потребителей
// (Header, SectionsMenu, Footer) — данные тянутся один раз за рендер.
export const useMenu = () => {
  const api = useApi()
  return useAsyncData('public-menu', () =>
    api
      .get<{ sections: MenuSection[]; pages: MenuPage[] }>('/menu')
      .catch(() => ({ sections: [], pages: [] }))
  )
}
