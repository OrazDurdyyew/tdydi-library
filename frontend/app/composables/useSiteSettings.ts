export interface SiteSettings {
  id: number
  name_tm: string
  name_ru: string
  name_en: string
  logo_light_path?: string
  logo_dark_path?: string
  favicon_path?: string
  primary_color: string
  secondary_color: string
  accent_color: string
  background_color: string
  foreground_color: string
  muted_color: string
  border_color: string
  ring_color: string
  font_family: string
  contact_address?: string
  contact_phone?: string
  contact_email?: string
  social_links?: Record<string, string>
  footer_tm: string
  footer_ru: string
  footer_en: string
}

export const useSiteSettings = () => {
  return useState<SiteSettings | null>('site_settings', () => null)
}
