<template>
  <div>
    <n-card class="header-card">
      <div class="text-16 font-bold">数据管理 {{ objectInfo ? `- ${objectInfo.name}` : '' }}</div>
    </n-card>

    <div v-if="!objectInfo" class="mt-4">
      <n-alert type="warning" class="tips-alert">
        请先选择一个对象，然后才能管理其数据
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
          <p class="text-16 font-bold">数据列表</p>
          <p class="text-gray-400">对象编码: {{ objectInfo?.code }}</p>
        </div>
        <n-button type="primary" @click="handleCreateData">添加数据</n-button>
      </div>

      <!-- 搜索区域 -->
      <n-card class="mb-4" size="small">
        <div class="mb-2">
          <n-space align="center">
            <span>过滤条件:</span>
            <n-button size="small" @click="addFilter" type="info">添加过滤条件</n-button>
            <n-button size="small" @click="clearFilters" v-if="filters.length > 0">清空</n-button>
          </n-space>
        </div>

        <n-space v-for="(filter, index) in filters" :key="index" class="mb-2" vertical>
          <div class="flex items-center gap-2">
            <n-select
              v-model:value="filter.field"
              :options="fieldOptions"
              placeholder="选择字段"
              style="width: 160px"
              @update:value="fieldChanged(index)"
            />

            <n-select
              v-model:value="filter.operator"
              :options="getOperatorOptions(filter.field)"
              placeholder="操作符"
              style="width: 120px"
            />

            <!-- 根据字段类型显示不同的输入控件 -->
            <template v-if="getFieldType(filter.field) === 'string' || getFieldType(filter.field) === 'text' || getFieldType(filter.field) === 'richtext'">
              <n-input v-model:value="filter.value" placeholder="输入值" style="width: 200px" />
            </template>

            <template v-else-if="getFieldType(filter.field) === 'number'">
              <n-input-number v-model:value="filter.value" placeholder="输入数值" style="width: 200px" />
            </template>

            <template v-else-if="getFieldType(filter.field) === 'boolean'">
              <n-select
                v-model:value="filter.value"
                :options="[
                  { label: '是', value: true },
                  { label: '否', value: false }
                ]"
                placeholder="选择"
                style="width: 200px"
              />
            </template>

            <template v-else-if="getFieldType(filter.field) === 'date'">
              <n-date-picker v-model:value="filter.value" type="date" style="width: 200px" />
            </template>

            <template v-else-if="getFieldType(filter.field) === 'datetime'">
              <n-date-picker v-model:value="filter.value" type="datetime" style="width: 200px" />
            </template>

            <template v-else-if="getFieldType(filter.field) === 'enum'">
              <n-select
                v-model:value="filter.value"
                :options="getEnumOptions(filter.field)"
                placeholder="选择"
                style="width: 200px"
              />
            </template>

            <n-button @click="removeFilter(index)" type="error" circle>
              <template #icon>
                <n-icon>
                  <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 24 24"><path d="M19 6.41L17.59 5L12 10.59L6.41 5L5 6.41L10.59 12L5 17.59L6.41 19L12 13.41L17.59 19L19 17.59L13.41 12L19 6.41z" fill="currentColor"></path></svg>
                </n-icon>
              </template>
            </n-button>
          </div>
        </n-space>

        <div class="flex justify-end mt-2">
          <n-button type="primary" @click="search" :disabled="filters.length === 0">搜索</n-button>
        </div>
      </n-card>

      <!-- 数据表格 -->
      <n-data-table
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :pagination="pagination"
        @update:page="handlePageChange"
      />

      <!-- 数据表单弹窗 -->
      <n-modal v-model:show="showDataModal" :title="modalTitle" style="width: 800px">
        <n-card>
          <n-form
            ref="formRef"
            :model="dataForm"
            :rules="formRules"
            label-placement="left"
            label-width="100"
            require-mark-placement="right-hanging"
          >
            <template v-for="field in objectInfo?.fields" :key="field.id">
              <!-- 文本输入框 -->
              <n-form-item v-if="field.field_type === 'string'" :label="field.name" :path="field.code">
                <n-input
                  v-model:value="dataForm[field.code]"
                  :placeholder="field.placeholder || '请输入' + field.name"
                />
                <template #help>
                  {{ field.help_text }}
                </template>
              </n-form-item>

              <!-- 长文本输入框 -->
              <n-form-item v-else-if="field.field_type === 'text'" :label="field.name" :path="field.code">
                <n-input
                  v-model:value="dataForm[field.code]"
                  type="textarea"
                  :placeholder="field.placeholder || '请输入' + field.name"
                />
                <template #help>
                  {{ field.help_text }}
                </template>
              </n-form-item>

              <!-- 富文本编辑器 -->
              <n-form-item v-else-if="field.field_type === 'richtext'" :label="field.name" :path="field.code">
                <n-input
                  v-model:value="dataForm[field.code]"
                  type="textarea"
                  :placeholder="field.placeholder || '请输入' + field.name"
                />
                <template #help>
                  {{ field.help_text }}
                </template>
              </n-form-item>

              <!-- 数字输入框 -->
              <n-form-item v-else-if="field.field_type === 'number'" :label="field.name" :path="field.code">
                <n-input-number
                  v-model:value="dataForm[field.code]"
                  :placeholder="field.placeholder || '请输入' + field.name"
                  :min="field.min_value"
                  :max="field.max_value"
                />
                <template #help>
                  {{ field.help_text }}
                </template>
              </n-form-item>

              <!-- 布尔开关 -->
              <n-form-item v-else-if="field.field_type === 'boolean'" :label="field.name" :path="field.code">
                <n-switch v-model:value="dataForm[field.code]" />
                <template #help>
                  {{ field.help_text }}
                </template>
              </n-form-item>

              <!-- 日期选择器 -->
              <n-form-item v-else-if="field.field_type === 'date'" :label="field.name" :path="field.code">
                <n-date-picker
                  v-model:value="dataForm[field.code]"
                  type="date"
                  :placeholder="field.placeholder || '请选择' + field.name"
                />
                <template #help>
                  {{ field.help_text }}
                </template>
              </n-form-item>

              <!-- 日期时间选择器 -->
              <n-form-item v-else-if="field.field_type === 'datetime'" :label="field.name" :path="field.code">
                <n-date-picker
                  v-model:value="dataForm[field.code]"
                  type="datetime"
                  :placeholder="field.placeholder || '请选择' + field.name"
                />
                <template #help>
                  {{ field.help_text }}
                </template>
              </n-form-item>

              <!-- 枚举选择器 -->
              <n-form-item v-else-if="field.field_type === 'enum'" :label="field.name" :path="field.code">
                <n-select
                  v-model:value="dataForm[field.code]"
                  :options="field.enum_options.map(opt => ({ label: opt.label, value: opt.value }))"
                  :placeholder="field.placeholder || '请选择' + field.name"
                />
                <template #help>
                  {{ field.help_text }}
                </template>
              </n-form-item>

              <!-- 文件上传 -->
              <n-form-item v-else-if="field.field_type === 'file'" :label="field.name" :path="field.code">
                <n-input
                  v-model:value="dataForm[field.code]"
                  :placeholder="field.placeholder || '请输入' + field.name"
                />
                <template #help>
                  {{ field.help_text || '请输入文件URL' }}
                </template>
              </n-form-item>

              <!-- 图片上传 -->
              <n-form-item v-else-if="field.field_type === 'image'" :label="field.name" :path="field.code">
                <n-input
                  v-model:value="dataForm[field.code]"
                  :placeholder="field.placeholder || '请输入' + field.name"
                />
                <template #help>
                  {{ field.help_text || '请输入图片URL' }}
                </template>
              </n-form-item>
            </template>
          </n-form>
          <template #footer>
            <div class="flex justify-end gap-2">
              <n-button @click="showDataModal = false">取消</n-button>
              <n-button type="primary" @click="handleSaveData">保存</n-button>
            </div>
          </template>
        </n-card>
      </n-modal>

      <!-- 确认删除弹窗 -->
      <n-modal v-model:show="showDeleteModal" preset="dialog" title="确认删除">
        <template #content>
          确定要删除这条数据吗？此操作无法恢复。
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
import { h, ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NSpace, NTag, NIcon, NDivider } from 'naive-ui'
import * as UdoApi from '@/api/udo'

const route = useRoute()
const router = useRouter()
const $message = window.$message
const formRef = ref(null)

const loading = ref(false)
const dataList = ref([])
const objectInfo = ref(null)
const showDataModal = ref(false)
const showDeleteModal = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const modalTitle = ref('添加数据')
const objectCode = ref('')
const objectId = ref(null)
const showSearchForm = ref(true)
const selectedObjectId = ref(null)
const loadingObjects = ref(false)
const objectsList = ref([])

// 查询过滤条件
const filters = ref([])

// 分页配置
const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  itemCount: 0,
  onChange: (page) => {
    pagination.page = page
    search()
  },
  onUpdatePageSize: (pageSize) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    search()
  }
})

const dataForm = reactive({})
const formRules = ref({})

// 可选择的字段
const fieldOptions = computed(() => {
  if (!objectInfo.value?.fields || !Array.isArray(objectInfo.value.fields)) {
    return []
  }
  // 返回所有字段，不再过滤 is_searchable
  return objectInfo.value.fields.map(field => ({
    label: field.name,
    value: field.code
  }))
})

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

// 动态表格列
const columns = computed(() => {
  if (!objectInfo.value || !objectInfo.value.fields || !Array.isArray(objectInfo.value.fields)) {
    return []
  }
  
  const fieldColumns = objectInfo.value.fields.map(field => ({
    title: field.name,
    key: field.code,
    render: (row) => formatFieldValue(row.data[field.code], field)
  }))

  return [
    ...fieldColumns.slice(0, 5), // 只显示前5个字段作为列
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
})

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
  UdoApi.getObjectById(id)
    .then((res) => {
      objectInfo.value = res
      objectCode.value = res.code // 更新输入框
      
      // 生成表单验证规则
      generateFormRules(res.fields)
      
      // 加载数据
      search()
    })
    .catch((err) => {
      $message.error(err.message || '获取对象信息失败')
    })
    .finally(() => {
      loading.value = false
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
      
      // 生成表单验证规则
      generateFormRules(res.fields)
    })
    .catch((err) => {
      $message.error(err.message || '获取对象信息失败')
    })
    .finally(() => {
      loading.value = false
    })
}

// 生成表单验证规则
function generateFormRules(fields) {
  if (!fields || !Array.isArray(fields)) {
    formRules.value = {};
    return;
  }
  
  const rules = {}
  fields.forEach(field => {
    const rule = []
    
    // 必填规则
    if (field.is_required) {
      rule.push({ required: true, message: `${field.name}为必填项`, trigger: 'blur' })
    }
    
    // 字符串长度规则
    if ((field.field_type === 'string' || field.field_type === 'text') && 
        (field.min_length !== null || field.max_length !== null)) {
      const lenRule = { trigger: 'blur' }
      
      if (field.min_length !== null) {
        lenRule.min = field.min_length
        lenRule.message = `${field.name}长度不能小于${field.min_length}`
      }
      
      if (field.max_length !== null) {
        lenRule.max = field.max_length
        lenRule.message = field.min_length !== null
          ? `${field.name}长度必须在${field.min_length}到${field.max_length}之间`
          : `${field.name}长度不能超过${field.max_length}`
      }
      
      rule.push(lenRule)
    }
    
    // 正则表达式规则
    if (field.regex_pattern) {
      rule.push({
        pattern: new RegExp(field.regex_pattern),
        message: field.regex_message || `${field.name}格式不正确`,
        trigger: 'blur'
      })
    }
    
    // 数值范围规则
    if (field.field_type === 'number' && 
        (field.min_value !== null || field.max_value !== null)) {
      const numRule = { trigger: 'change', type: 'number' }
      
      if (field.min_value !== null) {
        numRule.min = field.min_value
        numRule.message = `${field.name}不能小于${field.min_value}`
      }
      
      if (field.max_value !== null) {
        numRule.max = field.max_value
        numRule.message = field.min_value !== null
          ? `${field.name}必须在${field.min_value}到${field.max_value}之间`
          : `${field.name}不能超过${field.max_value}`
      }
      
      rule.push(numRule)
    }
    
    if (rule.length > 0) {
      rules[field.code] = rule
    }
  })
  
  formRules.value = rules
}

function search() {
  if (!objectId.value) {
    $message.warning('请先选择一个对象')
    return
  }
  
  loading.value = true
  const query = {
    page: pagination.page,
    page_size: pagination.pageSize,
    filters: filters.value
      .filter(f => f.field && f.operator && f.value !== undefined && f.value !== '')
      .map(f => ({
        field: f.field,
        operator: f.operator,
        value: f.value
      }))
  }
  
  UdoApi.queryData(objectId.value, query)
    .then((res) => {
      dataList.value = res.items || []
      pagination.itemCount = res.total
    })
    .catch((err) => {
      $message.error(err.message || '查询数据失败')
    })
    .finally(() => {
      loading.value = false
    })
}

function handlePageChange(page) {
  pagination.page = page
  search()
}

function addFilter() {
  filters.value.push({
    field: '',
    operator: 'eq',
    value: ''
  })
}

function removeFilter(index) {
  filters.value.splice(index, 1)
}

function clearFilters() {
  filters.value = []
}

function fieldChanged(index) {
  // 切换字段时重置操作符和值
  filters.value[index].operator = 'eq'
  filters.value[index].value = ''
}

function getFieldType(fieldCode) {
  if (!objectInfo.value?.fields) return 'string'
  const field = objectInfo.value.fields.find(f => f.code === fieldCode)
  return field ? field.field_type : 'string'
}

function getOperatorOptions(fieldCode) {
  const fieldType = getFieldType(fieldCode)
  
  // 通用操作符
  const commonOps = [
    { label: '等于', value: 'eq' },
    { label: '不等于', value: 'neq' }
  ]
  
  // 字符串特有操作符
  const stringOps = [
    { label: '包含', value: 'contains' },
    { label: '开头是', value: 'starts_with' },
    { label: '结尾是', value: 'ends_with' }
  ]
  
  // 数值特有操作符
  const numberOps = [
    { label: '大于', value: 'gt' },
    { label: '大于等于', value: 'gte' },
    { label: '小于', value: 'lt' },
    { label: '小于等于', value: 'lte' }
  ]
  
  if (['string', 'text', 'richtext'].includes(fieldType)) {
    return [...commonOps, ...stringOps]
  } else if (fieldType === 'number') {
    return [...commonOps, ...numberOps]
  } else {
    return commonOps
  }
}

function getEnumOptions(fieldCode) {
  if (!objectInfo.value?.fields || !Array.isArray(objectInfo.value.fields)) {
    return []
  }
  const field = objectInfo.value.fields.find(f => f.code === fieldCode)
  if (!field || !field.enum_options || !Array.isArray(field.enum_options)) {
    return []
  }
  return field.enum_options.map(opt => ({ label: opt.label, value: opt.value }))
}

function formatFieldValue(value, field) {
  if (value === undefined || value === null) {
    return '-'
  }
  
  switch (field.field_type) {
    case 'boolean':
      return h(NTag, { type: value ? 'success' : 'error' }, { default: () => value ? '是' : '否' })
    case 'enum':
      const option = field.enum_options?.find(opt => opt.value === value)
      return option ? option.label : value
    case 'date':
      return value ? new Date(value).toLocaleDateString() : '-'
    case 'datetime':
      return value ? new Date(value).toLocaleString() : '-'
    case 'image':
      return h('img', { src: value, alt: field.name, style: 'max-width: 50px; max-height: 50px;' })
    default:
      // 截取长文本
      if (typeof value === 'string' && value.length > 50) {
        return value.substring(0, 50) + '...'
      }
      return value
  }
}

function resetDataForm() {
  for (const key in dataForm) {
    delete dataForm[key]
  }
  
  // 设置默认值
  if (objectInfo.value && objectInfo.value.fields) {
    objectInfo.value.fields.forEach(field => {
      if (field.default_value) {
        dataForm[field.code] = field.default_value
      } else {
        // 根据字段类型设置默认空值
        switch (field.field_type) {
          case 'boolean':
            dataForm[field.code] = false
            break
          case 'number':
            dataForm[field.code] = null
            break
          default:
            dataForm[field.code] = ''
        }
      }
    })
  }
}

function handleCreateData() {
  if (!objectId.value) {
    $message.warning('请先选择一个对象')
    return
  }
  
  isEdit.value = false
  modalTitle.value = '添加数据'
  resetDataForm()
  showDataModal.value = true
}

function handleEdit(row) {
  isEdit.value = true
  modalTitle.value = '编辑数据'
  currentId.value = row.id
  
  resetDataForm()
  // 填充表单数据
  for (const key in row.data) {
    dataForm[key] = row.data[key]
  }
  
  showDataModal.value = true
}

function handleSaveData() {
  if (!objectId.value) {
    $message.warning('请先选择一个对象')
    return
  }
  
  formRef.value?.validate((errors) => {
    if (errors) return

    if (isEdit.value) {
      // 更新数据
      UdoApi.updateData(currentId.value, { data: { ...dataForm } })
        .then(() => {
          $message.success('更新数据成功')
          showDataModal.value = false
          search()
        })
        .catch((err) => {
          $message.error(err.message || '更新数据失败')
        })
    } else {
      // 创建数据
      UdoApi.createData(objectId.value, { data: { ...dataForm } })
        .then(() => {
          $message.success('创建数据成功')
          showDataModal.value = false
          search()
        })
        .catch((err) => {
          $message.error(err.message || '创建数据失败')
        })
    }
  })
}

function handleDelete(row) {
  currentId.value = row.id
  showDeleteModal.value = true
}

function confirmDelete() {
  UdoApi.deleteData(currentId.value)
    .then(() => {
      $message.success('删除数据成功')
      showDeleteModal.value = false
      search()
    })
    .catch((err) => {
      $message.error(err.message || '删除数据失败')
    })
}

function handleObjectSelect(id) {
  if (!id) {
    return;
  }
  
  objectId.value = id
  loadingObjects.value = true
  UdoApi.getObjectById(id)
    .then((res) => {
      objectInfo.value = res
      objectCode.value = res.code // 更新输入框
      
      // 生成表单验证规则
      generateFormRules(res.fields)
      
      // 重置过滤条件和分页
      filters.value = []
      pagination.page = 1
      
      // 加载数据
      search()
    })
    .catch((err) => {
      $message.error(err.message || '获取对象信息失败')
    })
    .finally(() => {
      loadingObjects.value = false
    })
}

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
</style> 