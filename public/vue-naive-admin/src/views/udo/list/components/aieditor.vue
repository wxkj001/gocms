<template>
  <div ref="divRef" />
</template>

<script setup>
import { AiEditor } from 'aieditor'
import 'aieditor/dist/style.css'

const props = defineProps({
})
const divRef = ref()
let aiEditor = null
onMounted(() => {
  aiEditor = new AiEditor({
    element: divRef.value,
    placeholder: '',
    content: '',
  })
})
defineExpose({
  aiEditor,
})
</script>
