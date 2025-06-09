<template>
  <MeModal ref="modalRef">
    <n-form
      ref="modalFormRef"
      label-placement="left"
      require-mark-placement="left"
      :rules="fieldRules"
      :label-width="100"
      :model="modalForm">
      <n-grid :cols="2" :x-gap="24">
        <n-grid-item>
          <n-form-item label="字段编码" path="code" >
            <n-input v-model:value="modalForm.code" placeholder="请输入字段编码" :disabled="modalAction === 'edit'" />
          </n-form-item>
        </n-grid-item>
        <n-grid-item>
          <n-form-item label="字段名称" path="name" >
            <n-input v-model:value="modalForm.name" placeholder="请输入字段名称" />
          </n-form-item>
        </n-grid-item>
      </n-grid>

      <n-grid :cols="2" :x-gap="24">
        <n-grid-item>
          <n-form-item label="字段类型" path="field_type">
            <n-select v-model:value="modalForm.field_type" :options="fieldTypeOptions" placeholder="请选择字段类型"
                      :disabled="modalAction === 'edit'" />
          </n-form-item>
        </n-grid-item>
        <n-grid-item>
          <n-form-item label="是否必填" path="is_required">
            <n-switch v-model:value="modalForm.is_required" :disabled="modalAction === 'edit'" />
          </n-form-item>
        </n-grid-item>
      </n-grid>

      <n-form-item label="默认值" path="default_value">
        <n-input v-model:value="modalForm.default_value" placeholder="请输入默认值" />
      </n-form-item>

      <!-- 不同字段类型的特定属性 -->
      <template v-if="modalForm.field_type === 'string' || modalForm.field_type === 'text'">
        <n-grid :cols="2" :x-gap="24">
          <n-grid-item>
            <n-form-item label="最小长度" path="min_length">
              <n-input-number v-model:value="modalForm.min_length" placeholder="请输入最小长度" clearable />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="最大长度" path="max_length">
              <n-input-number v-model:value="modalForm.max_length" placeholder="请输入最大长度" clearable />
            </n-form-item>
          </n-grid-item>
        </n-grid>

        <n-form-item label="正则表达式" path="regex_pattern">
          <n-input v-model:value="modalForm.regex_pattern" placeholder="请输入正则表达式" />
        </n-form-item>

        <n-form-item label="正则错误提示" path="regex_message">
          <n-input v-model:value="modalForm.regex_message" placeholder="请输入正则表达式验证失败时的提示信息" />
        </n-form-item>
      </template>

      <template v-if="modalForm.field_type === 'number'">
        <n-grid :cols="2" :x-gap="24">
          <n-grid-item>
            <n-form-item label="最小值" path="min_value">
              <n-input-number v-model:value="modalForm.min_value" placeholder="请输入最小值" clearable />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="最大值" path="max_value">
              <n-input-number v-model:value="modalForm.max_value" placeholder="请输入最大值" clearable />
            </n-form-item>
          </n-grid-item>
        </n-grid>
      </template>

      <template v-if="modalForm.field_type === 'enum'">
        <n-form-item label="枚举选项" path="enum_options">
          <div class="enum-options-container">
            <div v-for="(option, index) in modalForm.enum_options" :key="index" class="enum-option-item">
              <div class="flex items-center mb-2">
                <n-input v-model:value="option.value" placeholder="选项值" class="mr-2" />
                <n-input v-model:value="option.label" placeholder="选项标签" class="mr-2" />
                <n-button type="error" @click="removeOption(index)">删除</n-button>
              </div>
            </div>
            <n-button type="info" @click="addOption">添加选项</n-button>
          </div>
        </n-form-item>
      </template>

      <n-grid :cols="2" :x-gap="24">
        <n-grid-item>
          <n-form-item label="是否唯一" path="is_unique">
            <n-switch v-model:value="modalForm.is_unique" :disabled="modalAction === 'edit'" />
          </n-form-item>
        </n-grid-item>
        <n-grid-item>
          <n-form-item label="是否可搜索" path="is_searchable">
            <n-switch v-model:value="modalForm.is_searchable" />
          </n-form-item>
        </n-grid-item>
      </n-grid>

      <n-form-item label="占位符文本" path="placeholder">
        <n-input v-model:value="modalForm.placeholder" placeholder="请输入占位符文本" />
      </n-form-item>

      <n-form-item label="帮助文本" path="help_text">
        <n-input v-model:value="modalForm.help_text" placeholder="请输入帮助文本" />
      </n-form-item>

      <n-form-item label="字段描述" path="description">
        <n-input v-model:value="modalForm.description" type="textarea" placeholder="请输入字段描述" />
      </n-form-item>
    </n-form>
  </MeModal>
</template>

<script setup>
import { MeModal } from '@/components'
import { useForm, useModal } from '@/composables'

import api from '../api'

/* const props = defineProps({
  menus: {
    type: Array,
    required: true,
  },
}) */
const emit = defineEmits(['refresh'])
// 字段类型选项
const fieldTypeOptions = [
  { label: '文本', value: 'string' },
  { label: '长文本', value: 'text' },
  { label: '富文本', value: 'richtext' },
  { label: '数字', value: 'number' },
  { label: '布尔', value: 'boolean' },
  { label: '日期', value: 'date' },
  { label: '日期时间', value: 'datetime' },
  { label: '枚举', value: 'enum' },
  { label: '文件', value: 'file' },
  { label: '图片', value: 'image' },
]
const defaultForm = {
  code: '',
  name: '',
  description: '',
  field_type: 'string',
  is_required: false,
  min_length: null,
  max_length: null,
  regex_pattern: '',
  regex_message: '',
  min_value: null,
  max_value: null,
  enum_options: [{ value: '', label: '' }],
  default_value: '',
  placeholder: '',
  help_text: '',
  is_unique: false,
  is_searchable: false,
  sort_order: 0,
  status: 1,
}

const [modalFormRef, modalForm, validation] = useForm()
const [modalRef, okLoading] = useModal()
const fieldRules = {
  code: [
    { required: true, message: '请输入字段编码', trigger: 'blur' },
    {
      pattern: /^[a-z]\w*$/i,
      message: '字段编码必须以字母开头，只能包含字母、数字和下划线',
      trigger: 'blur',
    },
  ],
  name: [{ required: true, message: '请输入字段名称', trigger: 'blur' }],
  field_type: [{ required: true, message: '请选择字段类型', trigger: 'change' }],
}
const modalAction = ref('')
function handleOpen(options = {}) {
  const { action, row = {}, ...rest } = options
  modalAction.value = action
  modalForm.value = { ...defaultForm, ...row }
  // parentIdDisabled.value = !!row.parentId && (row.type === 'BUTTON' || row.type === 'API')
  modalRef.value.open({ ...rest, onOk: onSave })
}
function addOption() {
  modalForm.value.enum_options.push({ value: '', label: '' })
}
function removeOption(index) {
  if (modalForm.value.enum_options.length > 1) {
    modalForm.value.enum_options.splice(index, 1)
  }
}
async function onSave() {
  await validation()
  okLoading.value = true
  try {
    let newFormData
    if (modalAction.value === 'add') {
      const res = await api.addUdoField(modalForm.value)
      newFormData = res.data
    }
    else if (modalAction.value === 'edit') {
      // await api.savePermission(modalForm.value.id, modalForm.value)
      await api.getUdoFieldUp(modalForm.value)
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
