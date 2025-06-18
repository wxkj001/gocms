<template>
  <div ref="divRef" style="height: 550px;" />
</template>

<script setup>
import { useDark } from '@vueuse/core'
import { AiEditor } from 'aieditor'
import 'aieditor/dist/style.css'

const props = defineProps({
  modelValue: {
    type: String, // Or whatever type your AiEditor content is (e.g., Object)
    default: '',
  },
})
const emit = defineEmits(['update:modelValue'])
const isDark = useDark()
const divRef = ref()
let aiEditor = null
onMounted(() => {
  aiEditor = new AiEditor({
    element: divRef.value,
    placeholder: '',
    content: props.modelValue,
    theme: isDark.value ? 'dark' : 'light',
    onChange: (aiEditor) => {
      emit('update:modelValue', aiEditor.getHtml())
    },
    ai: {
      models: {
        openai: {
          endpoint: 'https://api.moonshot.cn',
          model: 'moonshot-v1-8k',
          apiKey: 'sk-alQ96zb******',
        },
      },
    },
  })
})
watch(() => props.modelValue, (newValue) => {
  if (aiEditor && aiEditor.getHtml() !== newValue) {
    aiEditor.setHtml(newValue)
  }
})
onUnmounted(() => {
  aiEditor && aiEditor.destroy()
})
defineExpose({
  aiEditor,
})
</script>
