
<template>
  <CommonPage>
    <div class="flex">
      <n-spin size="small" :show="treeLoading">
        <MenuTree
          v-model:current-udo="currentUdo"
          class="w-320 shrink-0"
          :tree-data="treeData"
          @refresh="initData"
        />
      </n-spin>

      <div class="ml-40 w-0 flex-1">
        <template v-if="currentUdo">
          <div class="flex justify-between">
            <h3 class="mb-12">
              {{ currentUdo.name }}
            </h3>
            <NButton size="small" type="primary" @click="handleEdit(currentUdo)">
              <i class="i-material-symbols:edit-outline mr-4 text-14" />
              编辑
            </NButton>
          </div>
          <n-descriptions label-placement="left" bordered :column="2">
            <n-descriptions-item label="编码">
              {{ currentUdo.code }}
            </n-descriptions-item>
            <n-descriptions-item label="名称">
              {{ currentUdo.name }}
            </n-descriptions-item>
            <n-descriptions-item label="路由地址">
              
            </n-descriptions-item>
          </n-descriptions>


          <div class="mt-32 flex justify-between">
            <h3 class="mb-12">
              字段列表
            </h3>
            <NButton size="small" type="primary" @click="handleAddBtn">
              <i class="i-fe:plus mr-4 text-14" />
              新增
            </NButton>
          </div>

          <MeCrud
            ref="$table"
            :columns="btnsColumns"
            :scroll-x="-1"
            :get-data="api.getButtons"
            :query-items="{ parentId: currentUdo.id }"
          />
        </template>
        <n-empty v-else class="h-450 f-c-c" size="large" description="请选择菜单查看详情" />
      </div>
    </div>
    <!-- <ResAddOrEdit ref="modalRef" :menus="treeData" @refresh="initData" /> -->
    <FieldAddOrEdit ref="modalRef" :menus="treeData" @refresh="initData" />
  </CommonPage>
</template>

<script setup>
import { MeCrud } from '@/components'
import { NButton, NSwitch } from 'naive-ui'
import api from './api'
import MenuTree from './components/MenuTree.vue'
import FieldAddOrEdit from './components/FieldAddOrEdit.vue'

const treeData = ref([])
const treeLoading = ref(false)
const $table = ref(null)
const $tableApi = ref(null)
const currentUdo = ref(null)
async function initData(data) {
  treeLoading.value = true
  const res = await api.getUdoObject()
  console.log("udo",data);
  
  treeData.value = res?.data || []
  treeLoading.value = false

  if (data)
  currentUdo.value = data
}
initData()

const modalRef = ref(null)

const btnsColumns = [
  { title: '名称', key: 'name' },
  { title: '编码', key: 'code' },
  {
    title: '状态',
    key: 'enable',
    render: row =>
      h(
        NSwitch,
        {
          size: 'small',
          rubberBand: false,
          value: row.enable,
          loading: !!row.enableLoading,
          onUpdateValue: () => handleEnable(row),
        },
        {
          checked: () => '启用',
          unchecked: () => '停用',
        },
      ),
  },
  {
    title: '操作',
    key: 'actions',
    width: 320,
    align: 'right',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 12px;',
            onClick: () => handleEditBtn(row),
          },
          {
            default: () => '编辑',
            icon: () => h('i', { class: 'i-material-symbols:edit-outline text-14' }),
          },
        ),

        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            style: 'margin-left: 12px;',
            onClick: () => handleDeleteBtn(row.id),
          },
          {
            default: () => '删除',
            icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
          },
        ),
      ]
    },
  },
]
const apisColumns = [
  { title: '名称', key: 'name' },
  { title: '编码', key: 'code' },
  { title: '路由', key: 'path' },
  {
    title: '状态',
    key: 'enable',
    render: row =>
      h(
        NSwitch,
        {
          size: 'small',
          rubberBand: false,
          value: row.enable,
          loading: !!row.enableLoading,
          onUpdateValue: () => handleApiEnable(row),
        },
        {
          checked: () => '启用',
          unchecked: () => '停用',
        },
      ),
  },
  {
    title: '操作',
    key: 'actions',
    width: 320,
    align: 'right',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 12px;',
            onClick: () => handleEditApi(row),
          },
          {
            default: () => '编辑',
            icon: () => h('i', { class: 'i-material-symbols:edit-outline text-14' }),
          },
        ),

        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            style: 'margin-left: 12px;',
            onClick: () => handleDeleteApi(row.id),
          },
          {
            default: () => '删除',
            icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
          },
        ),
      ]
    },
  },
]
watch(
  () => currentUdo.value,
  async (v) => {
    await nextTick()
    if (v) {
      $table.value.handleSearch()
      $tableApi.value.handleSearch()
    }
  },
)

function handleAddBtn() {
  modalRef.value?.handleOpen({
    action: 'add',
    title: '新增按钮',
    row: { type: 'BUTTON', parentId: currentUdo.value.id },
    okText: '保存',
  })
}
function handleAddApi() {
  modalRef.value?.handleOpen({
    action: 'add',
    title: '新增按钮',
    row: { type: 'API', parentId: currentUdo.value.id },
    okText: '保存',
  })
}
function handleEditBtn(row) {
  modalRef.value?.handleOpen({
    action: 'edit',
    title: `编辑按钮 - ${row.name}`,
    row,
    okText: '保存',
  })
}
function handleEditApi(row) {
  modalRef.value?.handleOpen({
    action: 'edit',
    title: `编辑按钮 - ${row.name}`,
    row,
    okText: '保存',
  })
}
function handleDeleteBtn(id) {
  const d = $dialog.warning({
    content: '确定删除？',
    title: '提示',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        d.loading = true
        await api.deletePermission(id)
        $message.success('删除成功')
        $table.value.handleSearch()
        d.loading = false
      }
      catch (error) {
        console.error(error)
        d.loading = false
      }
    },
  })
}
function handleDeleteApi(id) {
  const d = $dialog.warning({
    content: '确定删除？',
    title: '提示',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        d.loading = true
        await api.deletePermission(id)
        $message.success('删除成功')
        $tableApi.value.handleSearch()
        d.loading = false
      }
      catch (error) {
        console.error(error)
        d.loading = false
      }
    },
  })
}
async function handleEnable(item) {
  try {
    item.enableLoading = true
    await api.savePermission(item.id, {
      enable: !item.enable,
    })
    $message.success('操作成功')
    $table.value?.handleSearch()
    item.enableLoading = false
  }
  catch (error) {
    console.error(error)
    item.enableLoading = false
  }
}
async function handleApiEnable(item) {
  try {
    item.enableLoading = true
    await api.savePermission(item.id, {
      enable: !item.enable,
    })
    $message.success('操作成功')
    $tableApi.value?.handleSearch()
    item.enableLoading = false
  }
  catch (error) {
    console.error(error)
    item.enableLoading = false
  }
}
</script>
