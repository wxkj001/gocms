<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2024/04/01 15:52:31
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <MeModal ref="modalRef">
    <n-form
      ref="modalFormRef"
      label-placement="left"
      require-mark-placement="left"
      :label-width="100"
      :model="modalForm"
    >
      <n-grid :cols="24" :x-gap="24">
       
        <n-form-item-gi :span="12" path="name" :rule="required">
          <template #label>
            <QuestionLabel label="名称" content="标题" />
          </template>
          <n-input v-model:value="modalForm.name" />
        </n-form-item-gi>
        <n-form-item-gi :span="12" path="code" :rule="required">
          <template #label>
            <QuestionLabel label="编码" content="如果是菜单则对应前端路由的name，使用大驼峰" />
          </template>
          <n-input v-model:value="modalForm.code" />
        </n-form-item-gi>
      </n-grid>
    </n-form>
  </MeModal>
</template>

<script setup>
import { MeModal } from '@/components'
import { useForm, useModal } from '@/composables'
import icons from 'isme:icons'
import pagePathes from 'isme:page-pathes'
import api from '../api'
import QuestionLabel from './QuestionLabel.vue'

const props = defineProps({
  menus: {
    type: Array,
    required: true,
  },
})
const emit = defineEmits(['refresh'])

const menuOptions = computed(() => {
  return [{ name: '根菜单', id: '', children: props.menus || [] }]
})
const componentOptions = pagePathes.map(path => ({ label: path, value: path }))
const iconOptions = icons.map(item => ({
  label: () =>
    h('span', { class: 'flex items-center' }, [h('i', { class: `${item} text-18 mr-8` }), item]),
  value: item,
}))

const required = {
  required: true,
  message: '此为必填项',
  trigger: ['blur', 'change'],
}

const defaultForm = { enable: true, show: true, layout: '' }
const [modalFormRef, modalForm, validation] = useForm()
const [modalRef, okLoading] = useModal()

const modalAction = ref('')
const parentIdDisabled = ref(false)
function handleOpen(options = {}) {
  const { action, row = {}, ...rest } = options
  modalAction.value = action
  if (row.method !== '') {
    row.method = row.method?.split('|') || ''
  }

  modalForm.value = { ...defaultForm, ...row }
  parentIdDisabled.value = !!row.parentId && (row.type === 'BUTTON' || row.type === 'API')
  modalRef.value.open({ ...rest, onOk: onSave })
}

async function onSave() {
  await validation()
  okLoading.value = true
  try {
    let newFormData
    if (modalAction.value === 'add') {
     
      const res = await api.addUdoObject(modalForm.value)
      newFormData = res.data
    }
    else if (modalAction.value === 'edit') {
      
      
      await api.savePermission(modalForm.value.id, modalForm.value)
    }
    okLoading.value = false
    $message.success('保存成功')
    emit('refresh', modalAction.value === 'add' ? newFormData : modalForm.value)
  }
  catch (error) {
    console.error(error)
    okLoading.value = false
    return false
  }
}

defineExpose({
  handleOpen,
})
</script>
