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
const isDark = useDark()
const emit = defineEmits(['update:modelValue'])
const divRef = ref()
let aiEditor = null
onMounted(() => {
  console.log(isDark.value);
  aiEditor = new AiEditor({
    element: divRef.value,
    placeholder: '',
    content: props.modelValue,
    theme: isDark.value ? 'dark' : 'light',
    onChange:(aiEditor)=>{
        emit('update:modelValue', aiEditor.getHtml())
    }
  })
})
watch(() => props.modelValue, (newValue) => {
  if (aiEditor && aiEditor.getHtml() !== newValue) {
    aiEditor.setHtml(newValue)
  }
})
onUnmounted(() => {
  aiEditor && aiEditor.destroy();
})
defineExpose({
  aiEditor,
})
</script>
