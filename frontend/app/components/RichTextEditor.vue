<script setup lang="ts">
import { Extension } from '@tiptap/core'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Image from '@tiptap/extension-image'
import Underline from '@tiptap/extension-underline'
import { TextStyle } from '@tiptap/extension-text-style'
import Color from '@tiptap/extension-color'
import Highlight from '@tiptap/extension-highlight'
import FontFamily from '@tiptap/extension-font-family'
import TextAlign from '@tiptap/extension-text-align'
import Placeholder from '@tiptap/extension-placeholder'
import {
  Bold, Italic, Underline as UnderlineIcon, Strikethrough,
  Heading1, Heading2, Heading3, Pilcrow,
  List, ListOrdered, Quote, Code, Code2, Minus,
  AlignLeft, AlignCenter, AlignRight, AlignJustify,
  Link as LinkIcon, Image as ImageIcon,
  Undo2, Redo2, Eraser, FileCode2, Palette, Type
} from 'lucide-vue-next'

interface Props {
  modelValue: string
  placeholder?: string
  minHeight?: string
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '',
  minHeight: '320px'
})
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const api = useApi()
const toast = useToast()
const config = useRuntimeConfig()

const showSource = ref(false)
const sourceHtml = ref('')
const fileInput = ref<HTMLInputElement | null>(null)
const showColors = ref(false)
const showFontSize = ref(false)
const showFontFamily = ref(false)

const TEXT_COLORS = [
  '#000000', '#374151', '#6b7280', '#9ca3af', '#e5e7eb', '#ffffff',
  '#dc2626', '#ea580c', '#d97706', '#ca8a04', '#16a34a', '#059669',
  '#0891b2', '#2563eb', '#4f46e5', '#7c3aed', '#c026d3', '#db2777'
]

const FONT_SIZES = ['12px', '14px', '16px', '18px', '20px', '24px', '28px', '32px', '40px', '48px']

const FONT_FAMILIES = [
  { label: 'По умолчанию', value: '' },
  { label: 'Sans-serif', value: 'Inter, system-ui, sans-serif' },
  { label: 'Serif', value: '"Source Serif 4", Georgia, serif' },
  { label: 'Monospace', value: 'ui-monospace, "SFMono-Regular", Menlo, monospace' }
]

const FontSize = Extension.create({
  name: 'fontSize',
  addOptions() {
    return { types: ['textStyle'] as string[] }
  },
  addGlobalAttributes() {
    return [
      {
        types: this.options.types,
        attributes: {
          fontSize: {
            default: null,
            parseHTML: (element: HTMLElement) =>
              element.style.fontSize?.replace(/['"]+/g, '') || null,
            renderHTML: (attributes: Record<string, any>) => {
              if (!attributes.fontSize) return {}
              return { style: `font-size: ${attributes.fontSize}` }
            }
          }
        }
      }
    ]
  },
  addCommands() {
    return {
      setFontSize:
        (size: string) =>
        ({ chain }: any) =>
          chain().setMark('textStyle', { fontSize: size }).run(),
      unsetFontSize:
        () =>
        ({ chain }: any) =>
          chain().setMark('textStyle', { fontSize: null }).removeEmptyTextStyle().run()
    }
  }
}) as any

const editor = useEditor({
  content: props.modelValue || '',
  extensions: [
    StarterKit.configure({
      heading: { levels: [1, 2, 3] },
      codeBlock: { HTMLAttributes: { class: 'rte-codeblock' } }
    }),
    Underline,
    Link.configure({
      openOnClick: false,
      autolink: true,
      HTMLAttributes: { rel: 'noopener noreferrer', target: '_blank' }
    }),
    Image.configure({
      inline: false,
      allowBase64: false,
      HTMLAttributes: { class: 'rte-image' }
    }),
    TextStyle,
    Color,
    Highlight.configure({ multicolor: true }),
    FontFamily,
    FontSize,
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Placeholder.configure({ placeholder: props.placeholder })
  ],
  editorProps: {
    attributes: {
      class: 'rte-content prose-content focus:outline-none',
      spellcheck: 'false'
    }
  },
  onUpdate: ({ editor }) => {
    const html = editor.getHTML()
    emit('update:modelValue', html === '<p></p>' ? '' : html)
  }
})

watch(
  () => props.modelValue,
  (val) => {
    const e = editor.value
    if (!e) return
    if (val !== e.getHTML()) {
      e.commands.setContent(val || '', { emitUpdate: false })
    }
  }
)

onBeforeUnmount(() => {
  editor.value?.destroy()
})

function setLink() {
  const e = editor.value
  if (!e) return
  const previous = e.getAttributes('link').href as string | undefined
  const url = window.prompt('Адрес ссылки:', previous || 'https://')
  if (url === null) return
  if (url === '') {
    e.chain().focus().extendMarkRange('link').unsetLink().run()
    return
  }
  e.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
}

function pickImage() {
  fileInput.value?.click()
}

async function onImageFile(ev: Event) {
  const target = ev.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  try {
    const res = await api.uploadField<{ path: string; url: string }>('/admin/content/image', 'file', file)
    const fullUrl = res.url.startsWith('http')
      ? res.url
      : config.public.apiBase.replace(/\/api\/?$/, '') + res.url
    editor.value?.chain().focus().setImage({ src: fullUrl, alt: file.name }).run()
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка загрузки изображения')
  } finally {
    if (fileInput.value) fileInput.value.value = ''
  }
}

function setColor(c: string) {
  editor.value?.chain().focus().setColor(c).run()
  showColors.value = false
}
function clearColor() {
  editor.value?.chain().focus().unsetColor().run()
  showColors.value = false
}
function setHighlight(c: string) {
  editor.value?.chain().focus().toggleHighlight({ color: c }).run()
}
function setFontSize(size: string) {
  editor.value?.chain().focus().setFontSize(size).run()
  showFontSize.value = false
}
function clearFontSize() {
  editor.value?.chain().focus().unsetFontSize().run()
  showFontSize.value = false
}
function setFontFamily(family: string) {
  if (!family) editor.value?.chain().focus().unsetFontFamily().run()
  else editor.value?.chain().focus().setFontFamily(family).run()
  showFontFamily.value = false
}

function toggleSource() {
  const e = editor.value
  if (!e) return
  if (showSource.value) {
    e.commands.setContent(sourceHtml.value || '', { emitUpdate: true })
    showSource.value = false
  } else {
    sourceHtml.value = e.getHTML()
    showSource.value = true
  }
}

function onSourceInput() {
  emit('update:modelValue', sourceHtml.value)
}

function clearFormatting() {
  editor.value?.chain().focus().clearNodes().unsetAllMarks().run()
}

const isActive = (name: string, attrs?: Record<string, any>) =>
  !!editor.value?.isActive(name, attrs)
const isAlign = (align: 'left' | 'center' | 'right' | 'justify') =>
  !!(editor.value as any)?.isActive({ textAlign: align })
</script>

<template>
  <div class="rte" :class="{ 'rte--source': showSource }">
    <div class="rte-toolbar" v-if="editor">
      <div class="rte-group">
        <button type="button" class="rte-btn" :class="{ active: isActive('bold') }" @click="editor.chain().focus().toggleBold().run()" title="Жирный (Ctrl+B)">
          <Bold class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isActive('italic') }" @click="editor.chain().focus().toggleItalic().run()" title="Курсив (Ctrl+I)">
          <Italic class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isActive('underline') }" @click="editor.chain().focus().toggleUnderline().run()" title="Подчёркнутый (Ctrl+U)">
          <UnderlineIcon class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isActive('strike') }" @click="editor.chain().focus().toggleStrike().run()" title="Зачёркнутый">
          <Strikethrough class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isActive('code') }" @click="editor.chain().focus().toggleCode().run()" title="Моноширинный">
          <Code class="w-4 h-4" />
        </button>
      </div>

      <div class="rte-sep" />

      <div class="rte-group">
        <button type="button" class="rte-btn" :class="{ active: isActive('paragraph') && !isActive('heading') }" @click="editor.chain().focus().setParagraph().run()" title="Абзац">
          <Pilcrow class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isActive('heading', { level: 1 }) }" @click="editor.chain().focus().toggleHeading({ level: 1 }).run()" title="Заголовок 1">
          <Heading1 class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isActive('heading', { level: 2 }) }" @click="editor.chain().focus().toggleHeading({ level: 2 }).run()" title="Заголовок 2">
          <Heading2 class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isActive('heading', { level: 3 }) }" @click="editor.chain().focus().toggleHeading({ level: 3 }).run()" title="Заголовок 3">
          <Heading3 class="w-4 h-4" />
        </button>
      </div>

      <div class="rte-sep" />

      <div class="rte-group">
        <div class="rte-dropdown">
          <button type="button" class="rte-btn" @click="showFontSize = !showFontSize; showColors = false; showFontFamily = false" title="Размер шрифта">
            <Type class="w-4 h-4" /><span class="rte-caret">▾</span>
          </button>
          <div v-if="showFontSize" class="rte-menu">
            <button type="button" class="rte-menu-item" @click="clearFontSize">По умолчанию</button>
            <button
              v-for="size in FONT_SIZES"
              :key="size"
              type="button"
              class="rte-menu-item"
              :style="{ fontSize: size }"
              @click="setFontSize(size)"
            >{{ size }}</button>
          </div>
        </div>

        <div class="rte-dropdown">
          <button type="button" class="rte-btn" @click="showFontFamily = !showFontFamily; showColors = false; showFontSize = false" title="Шрифт">
            <FileCode2 class="w-4 h-4" /><span class="rte-caret">▾</span>
          </button>
          <div v-if="showFontFamily" class="rte-menu rte-menu--wide">
            <button
              v-for="f in FONT_FAMILIES"
              :key="f.label"
              type="button"
              class="rte-menu-item"
              :style="{ fontFamily: f.value || 'inherit' }"
              @click="setFontFamily(f.value)"
            >{{ f.label }}</button>
          </div>
        </div>

        <div class="rte-dropdown">
          <button type="button" class="rte-btn" @click="showColors = !showColors; showFontSize = false; showFontFamily = false" title="Цвет текста">
            <Palette class="w-4 h-4" /><span class="rte-caret">▾</span>
          </button>
          <div v-if="showColors" class="rte-menu rte-menu--palette">
            <div class="rte-menu-label">Цвет текста</div>
            <div class="rte-swatches">
              <button
                v-for="c in TEXT_COLORS"
                :key="'t' + c"
                type="button"
                class="rte-swatch"
                :style="{ background: c }"
                @click="setColor(c)"
                :title="c"
              />
            </div>
            <button type="button" class="rte-menu-item rte-menu-item--small" @click="clearColor">Сбросить цвет</button>
            <div class="rte-menu-label">Подсветка</div>
            <div class="rte-swatches">
              <button
                type="button"
                class="rte-swatch rte-swatch--off"
                @click="editor?.chain().focus().unsetHighlight().run()"
                title="Без подсветки"
              >×</button>
              <button
                v-for="c in ['#fef08a','#fde68a','#fecaca','#fbcfe8','#ddd6fe','#bae6fd','#bbf7d0']"
                :key="'h' + c"
                type="button"
                class="rte-swatch"
                :style="{ background: c }"
                @click="setHighlight(c)"
              />
            </div>
          </div>
        </div>
      </div>

      <div class="rte-sep" />

      <div class="rte-group">
        <button type="button" class="rte-btn" :class="{ active: isAlign('left') }" @click="editor.chain().focus().setTextAlign('left').run()" title="По левому краю">
          <AlignLeft class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isAlign('center') }" @click="editor.chain().focus().setTextAlign('center').run()" title="По центру">
          <AlignCenter class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isAlign('right') }" @click="editor.chain().focus().setTextAlign('right').run()" title="По правому краю">
          <AlignRight class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isAlign('justify') }" @click="editor.chain().focus().setTextAlign('justify').run()" title="По ширине">
          <AlignJustify class="w-4 h-4" />
        </button>
      </div>

      <div class="rte-sep" />

      <div class="rte-group">
        <button type="button" class="rte-btn" :class="{ active: isActive('bulletList') }" @click="editor.chain().focus().toggleBulletList().run()" title="Маркированный список">
          <List class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isActive('orderedList') }" @click="editor.chain().focus().toggleOrderedList().run()" title="Нумерованный список">
          <ListOrdered class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isActive('blockquote') }" @click="editor.chain().focus().toggleBlockquote().run()" title="Цитата">
          <Quote class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" :class="{ active: isActive('codeBlock') }" @click="editor.chain().focus().toggleCodeBlock().run()" title="Блок кода">
          <Code2 class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" @click="editor.chain().focus().setHorizontalRule().run()" title="Разделитель">
          <Minus class="w-4 h-4" />
        </button>
      </div>

      <div class="rte-sep" />

      <div class="rte-group">
        <button type="button" class="rte-btn" :class="{ active: isActive('link') }" @click="setLink" title="Ссылка">
          <LinkIcon class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" @click="pickImage" title="Изображение">
          <ImageIcon class="w-4 h-4" />
        </button>
        <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="onImageFile" />
      </div>

      <div class="rte-sep" />

      <div class="rte-group">
        <button type="button" class="rte-btn" @click="clearFormatting" title="Очистить форматирование">
          <Eraser class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" @click="editor.chain().focus().undo().run()" :disabled="!editor.can().undo()" title="Отменить (Ctrl+Z)">
          <Undo2 class="w-4 h-4" />
        </button>
        <button type="button" class="rte-btn" @click="editor.chain().focus().redo().run()" :disabled="!editor.can().redo()" title="Повторить (Ctrl+Y)">
          <Redo2 class="w-4 h-4" />
        </button>
      </div>

      <div class="rte-sep" />

      <div class="rte-group">
        <button type="button" class="rte-btn" :class="{ active: showSource }" @click="toggleSource" title="HTML-исходник">
          &lt;/&gt;
        </button>
      </div>
    </div>

    <div v-if="!showSource" class="rte-body" :style="{ minHeight }">
      <EditorContent :editor="editor" />
    </div>
    <textarea
      v-else
      v-model="sourceHtml"
      class="rte-source"
      :style="{ minHeight }"
      spellcheck="false"
      @input="onSourceInput"
    />
  </div>
</template>

<style scoped>
.rte {
  border: 1px solid rgb(var(--border));
  border-radius: 0.5rem;
  background: rgb(var(--background));
  overflow: hidden;
}

.rte-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.25rem;
  padding: 0.5rem;
  border-bottom: 1px solid rgb(var(--border));
  background: rgb(var(--muted) / 0.4);
}

.rte-group {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

.rte-sep {
  width: 1px;
  align-self: stretch;
  background: rgb(var(--border));
  margin: 0.125rem 0.25rem;
}

.rte-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 2rem;
  height: 2rem;
  padding: 0 0.4rem;
  border-radius: 0.375rem;
  font-size: 0.85rem;
  color: rgb(var(--foreground));
  background: transparent;
  border: 1px solid transparent;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s, color 0.15s;
}
.rte-btn:hover:not(:disabled) {
  background: rgb(var(--muted));
}
.rte-btn.active {
  background: rgb(var(--primary) / 0.12);
  color: rgb(var(--primary));
  border-color: rgb(var(--primary) / 0.4);
}
.rte-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.rte-caret {
  margin-left: 2px;
  font-size: 0.7rem;
  opacity: 0.7;
}

.rte-dropdown {
  position: relative;
}
.rte-menu {
  position: absolute;
  z-index: 30;
  top: calc(100% + 4px);
  left: 0;
  min-width: 10rem;
  max-height: 280px;
  overflow-y: auto;
  padding: 0.25rem;
  background: rgb(var(--popover, var(--background)));
  border: 1px solid rgb(var(--border));
  border-radius: 0.5rem;
  box-shadow: 0 4px 16px rgb(0 0 0 / 0.08);
}
.rte-menu--wide { min-width: 14rem; }
.rte-menu--palette { min-width: 16rem; padding: 0.5rem; }
.rte-menu-item {
  display: block;
  width: 100%;
  text-align: left;
  padding: 0.4rem 0.6rem;
  border-radius: 0.25rem;
  font-size: 0.85rem;
  background: transparent;
  border: 0;
  cursor: pointer;
}
.rte-menu-item:hover { background: rgb(var(--muted)); }
.rte-menu-item--small { font-size: 0.75rem; opacity: 0.8; }
.rte-menu-label {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: rgb(var(--muted-foreground));
  padding: 0.25rem 0.25rem;
  margin-top: 0.25rem;
}
.rte-swatches {
  display: grid;
  grid-template-columns: repeat(9, 1fr);
  gap: 4px;
  padding: 0.25rem;
}
.rte-swatch {
  width: 1.4rem;
  height: 1.4rem;
  border: 1px solid rgb(var(--border));
  border-radius: 4px;
  cursor: pointer;
  transition: transform 0.1s;
}
.rte-swatch:hover { transform: scale(1.1); }
.rte-swatch--off {
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  color: rgb(var(--muted-foreground));
}

.rte-body {
  padding: 1rem 1.25rem;
}
.rte-body :deep(.rte-content) {
  min-height: inherit;
  outline: none;
  line-height: 1.65;
  font-size: 0.95rem;
}
.rte-body :deep(.rte-content > * + *) { margin-top: 0.75em; }
.rte-body :deep(h1) { font-size: 1.75rem; font-weight: 700; }
.rte-body :deep(h2) { font-size: 1.4rem; font-weight: 700; }
.rte-body :deep(h3) { font-size: 1.15rem; font-weight: 600; }
.rte-body :deep(p.is-editor-empty:first-child::before) {
  content: attr(data-placeholder);
  color: rgb(var(--muted-foreground));
  float: left;
  height: 0;
  pointer-events: none;
}
.rte-body :deep(ul) { list-style: disc; padding-left: 1.5rem; }
.rte-body :deep(ol) { list-style: decimal; padding-left: 1.5rem; }
.rte-body :deep(blockquote) {
  border-left: 3px solid rgb(var(--primary));
  padding-left: 1rem;
  color: rgb(var(--muted-foreground));
  font-style: italic;
}
.rte-body :deep(code) {
  background: rgb(var(--muted));
  padding: 0.1rem 0.35rem;
  border-radius: 4px;
  font-size: 0.85em;
}
.rte-body :deep(.rte-codeblock),
.rte-body :deep(pre) {
  background: rgb(var(--muted));
  padding: 0.75rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.85em;
  overflow-x: auto;
}
.rte-body :deep(a) {
  color: rgb(var(--primary));
  text-decoration: underline;
}
.rte-body :deep(.rte-image),
.rte-body :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.5rem;
  display: block;
  margin: 0.75rem 0;
}
.rte-body :deep(mark) {
  padding: 0 0.15rem;
  border-radius: 2px;
}
.rte-body :deep(hr) {
  border: 0;
  border-top: 1px solid rgb(var(--border));
  margin: 1rem 0;
}

.rte-source {
  width: 100%;
  padding: 1rem 1.25rem;
  font-family: ui-monospace, "SFMono-Regular", Menlo, monospace;
  font-size: 0.8rem;
  line-height: 1.5;
  border: 0;
  outline: 0;
  resize: vertical;
  background: rgb(var(--background));
  color: rgb(var(--foreground));
}

.hidden { display: none; }
</style>
