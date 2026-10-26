<script setup lang="ts">
import { Swiper, SwiperSlide } from 'swiper/vue'
import { Autoplay, EffectFade } from 'swiper/modules'
import 'swiper/css'
import 'swiper/css/effect-fade'

// Фоновая карусель hero-блока. Картинки берутся из слайдов главной (home_slides).
const props = defineProps<{ slides: any[] }>()

const fileUrl = useFileUrl()
const images = computed(() =>
  (props.slides ?? []).map((s) => fileUrl(s.image_path)).filter(Boolean)
)
</script>

<template>
  <!--
    Mobile: обычный блок-карточка сверху с аспектом 16:9, без затемнения — длинный
    landscape-баннер не растягивается в высоту и не перекрывается тёмным фильтром.
    md+: абсолютный фон-слой с градиентным overlay для контраста с текстом.
  -->
  <div
    class="hero-carousel relative w-full aspect-[16/9] overflow-hidden border-b border-border
           md:absolute md:inset-0 md:-z-10 md:aspect-auto md:border-0"
  >
    <ClientOnly>
      <Swiper
        v-if="images.length"
        :modules="[Autoplay, EffectFade]"
        effect="fade"
        :fade-effect="{ crossFade: true }"
        :autoplay="{ delay: 6000, disableOnInteraction: false }"
        :loop="images.length > 1"
        :allow-touch-move="false"
        class="h-full w-full"
      >
        <SwiperSlide v-for="(img, i) in images" :key="i">
          <img :src="img" alt="" class="w-full h-full object-cover" />
        </SwiperSlide>
      </Swiper>
      <template #fallback>
        <img v-if="images[0]" :src="images[0]" alt="" class="w-full h-full object-cover" />
      </template>
    </ClientOnly>
    <!-- Затемняющий слой только на md+, где карусель — фон под текстом. -->
    <div class="absolute inset-0 hero-overlay hidden md:block" />
  </div>
</template>

<style scoped>
.hero-overlay {
  background:
    linear-gradient(to right, rgba(15, 23, 42, 0.92) 0%, rgba(15, 23, 42, 0.7) 45%, rgba(15, 23, 42, 0.4) 100%);
}
</style>
