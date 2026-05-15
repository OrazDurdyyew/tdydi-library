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
  <div class="absolute inset-0 -z-10 overflow-hidden">
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
    <!-- Затемняющий слой для читаемости текста -->
    <div class="absolute inset-0 hero-overlay" />
  </div>
</template>

<style scoped>
.hero-overlay {
  background:
    linear-gradient(to right, rgba(15, 23, 42, 0.92) 0%, rgba(15, 23, 42, 0.7) 45%, rgba(15, 23, 42, 0.4) 100%);
}
</style>
