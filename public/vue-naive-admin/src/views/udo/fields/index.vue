<template>
  <div>
    <n-card class="header-card">
      <div class="text-16 font-bold">字段管理 {{ objectInfo ? `- ${objectInfo.name}` : '' }}</div>
    </n-card>

    <div v-if="!objectInfo" class="mt-4">
      <n-alert type="warning" class="tips-alert">
        请先选择一个对象，然后才能管理其字段
        <template #action>
          <n-button text @click="showSearchForm = !showSearchForm">
            {{ showSearchForm ? '隐藏搜索' : '按编码搜索' }}
          </n-button>
        </template>
      </n-alert>

      <!-- 对象搜索卡片，默认隐藏 -->
      <n-card v-if="showSearchForm" class="search-card">
        <n-form inline>
          <n-form-item label="选择对象" class="mr-4">
            <n-select 
              v-model:value="selectedObjectId" 
              :options="objectOptions" 
              placeholder="请选择对象"
              :loading="loadingObjects"
              @update:value="handleObjectSelect"
              style="width: 300px"
            />
          </n-form-item>
        </n-form>
      </n-card>
    </div>

    <n-card v-if="objectInfo" class="mt-4">
      <div class="flex justify-between mb-4">
        <div>
          <p class="text-16 font-bold">字段列表</p>
          <p class="text-gray-400">对象编码: {{ objectInfo?.code }}</p>
        </div>
        <n-button type="primary" @click="handleCreateField">添加字段</n-button>
      </div>

      <n-data-table
        :columns="columns"
        :data="fieldList"
        :loading="loading"
      />

      <!-- 字段表单弹窗 -->
      <n-modal v-model:show="showFieldModal" :title="modalTitle" :style="{width: '800px'}">
        <n-card>
          <n-form
            ref="formRef"
            :model="fieldForm"
            :rules="fieldRules"
            label-placement="left"
            label-width="120"
            require-mark-placement="right-hanging"
          >
            <n-grid :cols="2" :x-gap="24">
              <n-grid-item>
                <n-form-item label="字段编码" path="code">
                  <n-input v-model:value="fieldForm.code" placeholder="请输入字段编码" :disabled="isEdit" />
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="字段名称" path="name">
                  <n-input v-model:value="fieldForm.name" placeholder="请输入字段名称" />
                </n-form-item>
              </n-grid-item>
            </n-grid>

            <n-grid :cols="2" :x-gap="24">
              <n-grid-item>
                <n-form-item label="字段类型" path="field_type">
                  <n-select
                    v-model:value="fieldForm.field_type"
                    :options="fieldTypeOptions"
                    placeholder="请选择字段类型"
                    :disabled="isEdit"
                  />
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="是否必填" path="is_required">
                  <n-switch v-model:value="fieldForm.is_required" />
                </n-form-item>
              </n-grid-item>
            </n-grid>

            <n-form-item label="默认值" path="default_value">
              <n-input v-model:value="fieldForm.default_value" placeholder="请输入默认值" />
            </n-form-item>

            <!-- 不同字段类型的特定属性 -->
            <template v-if="fieldForm.field_type === 'string' || fieldForm.field_type === 'text'">
              <n-grid :cols="2" :x-gap="24">
                <n-grid-item>
                  <n-form-item label="最小长度" path="min_length">
                    <n-input-number v-model:value="fieldForm.min_length" placeholder="请输入最小长度" clearable />
                  </n-form-item>
                </n-grid-item>
                <n-grid-item>
                  <n-form-item label="最大长度" path="max_length">
                    <n-input-number v-model:value="fieldForm.max_length" placeholder="请输入最大长度" clearable />
                  </n-form-item>
                </n-grid-item>
              </n-grid>

              <n-form-item label="正则表达式" path="regex_pattern">
                <n-input v-model:value="fieldForm.regex_pattern" placeholder="请输入正则表达式" />
              </n-form-item>

              <n-form-item label="正则错误提示" path="regex_message">
                <n-input v-model:value="fieldForm.regex_message" placeholder="请输入正则表达式验证失败时的提示信息" />
              </n-form-item>
            </template>

            <template v-if="fieldForm.field_type === 'number'">
              <n-grid :cols="2" :x-gap="24">
                <n-grid-item>
                  <n-form-item label="最小值" path="min_value">
                    <n-input-number v-model:value="fieldForm.min_value" placeholder="请输入最小值" clearable />
                  </n-form-item>
                </n-grid-item>
                <n-grid-item>
                  <n-form-item label="最大值" path="max_value">
                    <n-input-number v-model:value="fieldForm.max_value" placeholder="请输入最大值" clearable />
                  </n-form-item>
                </n-grid-item>
              </n-grid>
            </template>

            <template v-if="fieldForm.field_type === 'enum'">
              <n-form-item label="枚举选项" path="enum_options">
                <div class="enum-options-container">
                  <div v-for="(option, index) in fieldForm.enum_options" :key="index" class="enum-option-item">
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
                  <n-switch v-model:value="fieldForm.is_unique" />
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="是否可搜索" path="is_searchable">
                  <n-switch v-model:value="fieldForm.is_searchable" />
                </n-form-item>
              </n-grid-item>
            </n-grid>

            <n-form-item label="占位符文本" path="placeholder">
              <n-input v-model:value="fieldForm.placeholder" placeholder="请输入占位符文本" />
            </n-form-item>

            <n-form-item label="帮助文本" path="help_text">
              <n-input v-model:value="fieldForm.help_text" placeholder="请输入帮助文本" />
            </n-form-item>

            <n-form-item label="字段描述" path="description">
              <n-input
                v-model:value="fieldForm.description"
                type="textarea"
                placeholder="请输入字段描述"
              />
            </n-form-item>
          </n-form>
          <template #footer>
            <div class="flex justify-end gap-2">
              <n-button @click="showFieldModal = false">取消</n-button>
              <n-button type="primary" @click="handleSaveField">保存</n-button>
            </div>
          </template>
        </n-card>
      </n-modal>

      <!-- 确认删除弹窗 -->
      <n-modal v-model:show="showDeleteModal" preset="dialog" title="确认删除">
        <template #content>
          确定要删除这个字段吗？如果该字段已经被使用，对应的数据将永久丢失。
        </template>
        <template #action>
          <div class="flex justify-end gap-2">
            <n-button @click="showDeleteModal = false">取消</n-button>
            <n-button type="error" @click="confirmDelete">删除</n-button>
          </div>
        </template>
      </n-modal>
    </n-card>
  </div>
</template>

<script setup>
import { h, ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NSpace, NSwitch, NTag, NIcon, NSelect, NDivider } from 'naive-ui'
import * as UdoApi from '@/api/udo'

const route = useRoute()
const router = useRouter()
const $message = window.$message

const loading = ref(false)
const fieldList = ref([])
const objectInfo = ref(null)
const showFieldModal = ref(false)
const showDeleteModal = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const modalTitle = ref('添加字段')
const objectCode = ref('')
const objectId = ref(null)
const showSearchForm = ref(true)
const selectedObjectId = ref(null)
const loadingObjects = ref(false)
const objectsList = ref([])

// 对象选项
const objectOptions = computed(() => {
  if (!objectsList.value || !Array.isArray(objectsList.value)) {
    return [];
  }
  return objectsList.value.map(obj => ({
    label: `${obj.name} (${obj.code})`,
    value: obj.id
  }))
})

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
  { label: '图片', value: 'image' }
]

const fieldForm = reactive({
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
  sort_order: 0
})

const fieldRules = {
  code: [
    { required: true, message: '请输入字段编码', trigger: 'blur' },
    {
      pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/,
      message: '字段编码必须以字母开头，只能包含字母、数字和下划线',
      trigger: 'blur'
    }
  ],
  name: [{ required: true, message: '请输入字段名称', trigger: 'blur' }],
  field_type: [{ required: true, message: '请选择字段类型', trigger: 'change' }]
}

const formRef = ref(null)

const columns = [
  {
    title: '字段编码',
    key: 'code'
  },
  {
    title: '字段名称',
    key: 'name'
  },
  {
    title: '字段类型',
    key: 'field_type',
    render(row) {
      const typeMap = {
        string: '文本',
        text: '长文本',
        richtext: '富文本',
        number: '数字',
        boolean: '布尔',
        date: '日期',
        datetime: '日期时间',
        enum: '枚举',
        file: '文件',
        image: '图片'
      }
      return h(NTag, { type: 'info' }, { default: () => typeMap[row.field_type] || row.field_type })
    }
  },
  {
    title: '必填',
    key: 'is_required',
    render(row) {
      return h(NSwitch, {
        value: row.is_required,
        disabled: true,
        size: 'small'
      })
    }
  },
  {
    title: '描述',
    key: 'description'
  },
  {
    title: '操作',
    key: 'actions',
    render(row) {
      return h(NSpace, null, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              onClick: () => handleEdit(row)
            },
            { default: () => '编辑' }
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              onClick: () => handleDelete(row)
            },
            { default: () => '删除' }
          )
        ]
      })
    }
  }
]

onMounted(() => {
  // 检查URL中是否有ID参数
  const routeId = route.params.id
  if (routeId && routeId !== 'undefined') {
    objectId.value = routeId
    loadObjectById(routeId)
  }
  
  // 加载所有对象列表
  loadAllObjects()
})

// 通过ID加载对象
function loadObjectById(id) {
  loading.value = true
  
  // 获取对象基本信息
  UdoApi.getObjectById(id)
    .then((res) => {
      objectInfo.value = res
      objectCode.value = res.code
    })
    .catch((err) => {
      $message.error(err.message || '获取对象信息失败')
    })
    .finally(() => {
      loading.value = false
    })
  
  // 直接获取字段列表
  UdoApi.getObjectFields(id)
    .then((res) => {
      fieldList.value = Array.isArray(res) ? res : []
    })
    .catch((err) => {
      $message.error(err.message || '获取字段列表失败')
    })
}

// 通过编码搜索对象
function searchObject() {
  if (!objectCode.value) {
    $message.warning('请输入对象编码')
    return
  }
  
  loading.value = true
  UdoApi.getObjectByCode(objectCode.value)
    .then((res) => {
      objectInfo.value = res
      objectId.value = res.id
      
      // 直接获取字段列表
      return UdoApi.getObjectFields(res.id)
    })
    .then((fields) => {
      if (Array.isArray(fields)) {
        fieldList.value = fields
      }
    })
    .catch((err) => {
      $message.error(err.message || '获取对象信息失败')
    })
    .finally(() => {
      loading.value = false
    })
}

function resetFieldForm() {
  fieldForm.code = ''
  fieldForm.name = ''
  fieldForm.description = ''
  fieldForm.field_type = 'string'
  fieldForm.is_required = false
  fieldForm.min_length = null
  fieldForm.max_length = null
  fieldForm.regex_pattern = ''
  fieldForm.regex_message = ''
  fieldForm.min_value = null
  fieldForm.max_value = null
  fieldForm.enum_options = [{ value: '', label: '' }]
  fieldForm.default_value = ''
  fieldForm.placeholder = ''
  fieldForm.help_text = ''
  fieldForm.is_unique = false
  fieldForm.is_searchable = false
  fieldForm.sort_order = fieldList.value.length * 10
}

function handleCreateField() {
  if (!objectId.value) {
    $message.warning('请先选择一个对象')
    return
  }
  
  isEdit.value = false
  modalTitle.value = '添加字段'
  resetFieldForm()
  showFieldModal.value = true
}

function handleEdit(row) {
  isEdit.value = true
  modalTitle.value = '编辑字段'
  currentId.value = row.id
  
  // 填充表单
  fieldForm.code = row.code
  fieldForm.name = row.name
  fieldForm.description = row.description
  fieldForm.field_type = row.field_type
  fieldForm.is_required = row.is_required
  fieldForm.min_length = row.min_length
  fieldForm.max_length = row.max_length
  fieldForm.regex_pattern = row.regex_pattern
  fieldForm.regex_message = row.regex_message
  fieldForm.min_value = row.min_value
  fieldForm.max_value = row.max_value
  fieldForm.default_value = row.default_value
  fieldForm.placeholder = row.placeholder
  fieldForm.help_text = row.help_text
  fieldForm.is_unique = row.is_unique
  fieldForm.is_searchable = row.is_searchable
  fieldForm.sort_order = row.sort_order

  // 处理枚举选项
  if (row.enum_options && row.enum_options.length > 0) {
    fieldForm.enum_options = [...row.enum_options]
  } else {
    fieldForm.enum_options = [{ value: '', label: '' }]
  }

  showFieldModal.value = true
}

function addOption() {
  fieldForm.enum_options.push({ value: '', label: '' })
}

function removeOption(index) {
  if (fieldForm.enum_options.length > 1) {
    fieldForm.enum_options.splice(index, 1)
  }
}

function handleSaveField() {
  if (!objectId.value) {
    $message.warning('请先选择一个对象')
    return
  }
  
  formRef.value?.validate((errors) => {
    if (errors) return

    // 过滤空的enum选项
    if (fieldForm.field_type === 'enum') {
      fieldForm.enum_options = fieldForm.enum_options.filter(
        option => option.value.trim() !== '' && option.label.trim() !== ''
      )
      if (fieldForm.enum_options.length === 0) {
        $message.error('枚举选项不能为空')
        return
      }
    }

    const formData = {
      code: fieldForm.code,
      name: fieldForm.name,
      description: fieldForm.description,
      field_type: fieldForm.field_type,
      is_required: fieldForm.is_required,
      default_value: fieldForm.default_value,
      placeholder: fieldForm.placeholder,
      help_text: fieldForm.help_text,
      is_unique: fieldForm.is_unique,
      is_searchable: fieldForm.is_searchable,
      sort_order: fieldForm.sort_order
    }

    // 根据字段类型添加特定属性
    if (fieldForm.field_type === 'string' || fieldForm.field_type === 'text') {
      formData.min_length = fieldForm.min_length
      formData.max_length = fieldForm.max_length
      formData.regex_pattern = fieldForm.regex_pattern
      formData.regex_message = fieldForm.regex_message
    } else if (fieldForm.field_type === 'number') {
      formData.min_value = fieldForm.min_value
      formData.max_value = fieldForm.max_value
    } else if (fieldForm.field_type === 'enum') {
      formData.enum_options = fieldForm.enum_options
    }

    if (isEdit.value) {
      // 更新字段
      UdoApi.updateField(currentId.value, formData)
        .then(() => {
          $message.success('更新字段成功')
          showFieldModal.value = false
          // 重新加载当前对象的字段
          loadObjectById(objectId.value)
        })
        .catch((err) => {
          $message.error(err.message || '更新字段失败')
        })
    } else {
      // 创建字段
      UdoApi.createField(objectId.value, formData)
        .then(() => {
          $message.success('创建字段成功')
          showFieldModal.value = false
          // 重新加载当前对象的字段
          loadObjectById(objectId.value)
        })
        .catch((err) => {
          $message.error(err.message || '创建字段失败')
        })
    }
  })
}

function handleDelete(row) {
  currentId.value = row.id
  showDeleteModal.value = true
}

function confirmDelete() {
  UdoApi.deleteField(currentId.value)
    .then(() => {
      $message.success('删除字段成功')
      showDeleteModal.value = false
      // 重新加载当前对象的字段
      loadObjectById(objectId.value)
    })
    .catch((err) => {
      $message.error(err.message || '删除字段失败')
    })
}

function handleObjectSelect(id) {
  if (!id) {
    return;
  }
  
  objectId.value = id
  loadingObjects.value = true
  
  // 获取对象基本信息
  UdoApi.getObjectById(id)
    .then((res) => {
      objectInfo.value = res
      objectCode.value = res.code
    })
    .catch((err) => {
      $message.error(err.message || '获取对象信息失败')
    })
  
  // 直接获取字段列表
  UdoApi.getObjectFields(id)
    .then((res) => {
      fieldList.value = Array.isArray(res) ? res : []
    })
    .catch((err) => {
      $message.error(err.message || '获取字段列表失败')
    })
    .finally(() => {
      loadingObjects.value = false
    })
}

// 加载所有对象
function loadAllObjects() {
  loadingObjects.value = true
  UdoApi.getObjects()
    .then((res) => {
      objectsList.value = Array.isArray(res) ? res : []
    })
    .catch((err) => {
      $message.error(err.message || '获取对象列表失败')
    })
    .finally(() => {
      loadingObjects.value = false
    })
}
</script>

<style scoped>
.header-card {
  margin-bottom: 16px;
}

.mt-4 {
  margin-top: 16px;
}

.tips-alert {
  margin-bottom: 16px;
}

.search-card {
  margin-bottom: 16px;
  background-color: #fafafa;
}

.enum-options-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.enum-option-item {
  border-bottom: 1px dashed #eee;
  padding-bottom: 8px;
}
</style> 