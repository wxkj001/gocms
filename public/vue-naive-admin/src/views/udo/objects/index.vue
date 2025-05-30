<template>
  <n-card>
    <template #header>
      <div class="flex justify-between">
        <div class="text-16 font-bold">自定义对象管理</div>
        <n-button type="primary" size="large" @click="handleCreateObject" class="create-btn">
          <template #icon>
            <n-icon>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 5v14M5 12h14"></path>
              </svg>
            </n-icon>
          </template>
          新建对象
        </n-button>
      </div>
    </template>

    <!-- 使用描述文本提示用户操作流程 -->
    <div class="usage-tip">
      <n-alert type="info" closable>
        <template #icon>
          <n-icon>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <path d="M12 8v4M12 16h.01"></path>
            </svg>
          </n-icon>
        </template>
        <div>
          <strong>使用指南</strong>：
          <ol class="mt-2 ml-4">
            <li>首先点击"<strong>新建对象</strong>"创建一个对象</li>
            <li>创建对象后，使用"<strong>字段管理</strong>"按钮为对象添加字段</li>
            <li>最后，使用"<strong>数据管理</strong>"按钮录入具体数据</li>
          </ol>
        </div>
      </n-alert>
    </div>

    <n-data-table
      :columns="columns"
      :data="objectList"
      :loading="loading"
      :pagination="pagination"
      @update:page="handlePageChange"
    />

    <!-- 对象表单弹窗 -->
    <n-modal v-model:show="showObjectModal" :title="modalTitle">
      <n-card style="width: 600px">
        <n-form
          ref="formRef"
          :model="objectForm"
          :rules="objectRules"
          label-placement="left"
          label-width="80"
          require-mark-placement="right-hanging"
        >
          <n-form-item label="对象编码" path="code">
            <n-input v-model:value="objectForm.code" placeholder="请输入对象编码" :disabled="isEdit" />
            <template #help>
              <span class="text-xs text-gray-500">对象编码必须以字母开头，只能包含字母、数字和下划线</span>
            </template>
          </n-form-item>
          <n-form-item label="对象名称" path="name">
            <n-input v-model:value="objectForm.name" placeholder="请输入对象名称" />
          </n-form-item>
          <n-form-item label="描述" path="description">
            <n-input
              v-model:value="objectForm.description"
              type="textarea"
              placeholder="请输入描述"
            />
          </n-form-item>
        </n-form>
        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showObjectModal = false">取消</n-button>
            <n-button type="primary" @click="handleSaveObject">保存</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>

    <!-- 确认删除弹窗 -->
    <n-modal v-model:show="showDeleteModal" preset="dialog" title="确认删除">
      <template #content>
        确定要删除这个对象吗？此操作将删除对象的所有字段和数据记录，且无法恢复。
      </template>
      <template #action>
        <div class="flex justify-end gap-2">
          <n-button @click="showDeleteModal = false">取消</n-button>
          <n-button type="error" @click="confirmDelete">删除</n-button>
        </div>
      </template>
    </n-modal>
  </n-card>
</template>

<script setup>
import { h, ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NSpace, NTag } from 'naive-ui'
import * as UdoApi from '@/api/udo'

const router = useRouter()
const $message = window.$message

const loading = ref(false)
const objectList = ref([])
const showObjectModal = ref(false)
const showDeleteModal = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const modalTitle = ref('新建对象')

const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  onChange: (page) => {
    pagination.page = page
  },
  onUpdatePageSize: (pageSize) => {
    pagination.pageSize = pageSize
    pagination.page = 1
  }
})

const objectForm = reactive({
  code: '',
  name: '',
  description: ''
})

const objectRules = {
  code: [
    { required: true, message: '请输入对象编码', trigger: 'blur' },
    {
      pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/,
      message: '对象编码必须以字母开头，只能包含字母、数字和下划线',
      trigger: 'blur'
    }
  ],
  name: [{ required: true, message: '请输入对象名称', trigger: 'blur' }]
}

const formRef = ref(null)

const columns = [
  {
    title: '对象编码',
    key: 'code'
  },
  {
    title: '对象名称',
    key: 'name'
  },
  {
    title: '字段数',
    key: 'fieldsCount',
    render: (row) => {
      return row.fields?.length || 0
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
              onClick: () => handleEditFields(row)
            },
            { default: () => '字段管理' }
          ),
          h(
            NButton,
            {
              size: 'small',
              onClick: () => handleViewData(row)
            },
            { default: () => '数据管理' }
          ),
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
  loadObjects()
})

function loadObjects() {
  loading.value = true
  UdoApi.getObjects()
    .then((res) => {
      objectList.value = Array.isArray(res) ? res : []
    })
    .catch((err) => {
      $message.error(err.message || '获取对象列表失败')
    })
    .finally(() => {
      loading.value = false
    })
}

function handlePageChange(page) {
  pagination.page = page
}

function handleCreateObject() {
  isEdit.value = false
  modalTitle.value = '新建对象'
  objectForm.code = ''
  objectForm.name = ''
  objectForm.description = ''
  showObjectModal.value = true
}

function handleEdit(row) {
  isEdit.value = true
  modalTitle.value = '编辑对象'
  currentId.value = row.id
  objectForm.code = row.code
  objectForm.name = row.name
  objectForm.description = row.description
  showObjectModal.value = true
}

function handleSaveObject() {
  formRef.value?.validate((errors) => {
    if (errors) return

    if (isEdit.value) {
      // 更新对象
      UdoApi.updateObject(currentId.value, {
        name: objectForm.name,
        description: objectForm.description
      })
        .then(() => {
          $message.success('更新对象成功')
          showObjectModal.value = false
          loadObjects()
        })
        .catch((err) => {
          $message.error(err.message || '更新对象失败')
        })
    } else {
      // 创建对象
      UdoApi.createObject({
        code: objectForm.code,
        name: objectForm.name,
        description: objectForm.description
      })
        .then(() => {
          $message.success('创建对象成功')
          showObjectModal.value = false
          loadObjects()
        })
        .catch((err) => {
          $message.error(err.message || '创建对象失败')
        })
    }
  })
}

function handleEditFields(row) {
  router.push(`/udo/fields/${row.id}`)
}

function handleViewData(row) {
  router.push(`/udo/data/${row.id}`)
}

function handleDelete(row) {
  currentId.value = row.id
  showDeleteModal.value = true
}

function confirmDelete() {
  UdoApi.deleteObject(currentId.value)
    .then(() => {
      $message.success('删除对象成功')
      showDeleteModal.value = false
      loadObjects()
    })
    .catch((err) => {
      $message.error(err.message || '删除对象失败')
    })
}
</script>

<style scoped>
.create-btn {
  font-weight: bold;
  padding: 0 20px;
}

.usage-tip {
  margin-bottom: 20px;
}

.mt-2 {
  margin-top: 8px;
}

.ml-4 {
  margin-left: 16px;
}
</style> 