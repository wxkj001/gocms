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
              /pm/udo/info/{{ currentUdo.code }}
            </n-descriptions-item>
          </n-descriptions>

          <div class="mt-32 flex justify-between">
            <h3 class="mb-12">
              字段列表
            </h3>
            <NButton size="small" type="primary" @click="handleAddField">
              <i class="i-fe:plus mr-4 text-14" />
              新增
            </NButton>
          </div>

          <MeCrud
            ref="$table"
            :columns="btnsColumns"
            :scroll-x="-1"
            :get-data="api.getUdoFieldList"
            :query-items="{ object_id: currentUdo.id }"
          />
        </template>
        <n-empty v-else class="h-450 f-c-c" size="large" description="请选择菜单查看详情" />
      </div>
    </div>
    <!-- <ResAddOrEdit ref="modalRef" :menus="treeData" @refresh="initData" /> -->
    <FieldAddOrEdit ref="modalRef" :menus="treeData" @refresh="initfieldData" />
  </CommonPage>
</template>

<script setup>
import { MeCrud } from '@/components'
import { NButton, NSwitch, NTag } from 'naive-ui'
import api from './api'
import FieldAddOrEdit from './components/FieldAddOrEdit.vue'
import MenuTree from './components/MenuTree.vue'

const treeData = ref([])
const treeLoading = ref(false)
const $table = ref(null)
const $tableApi = ref(null)
const currentUdo = ref(null)
async function initData(data) {
  treeLoading.value = true
  const res = await api.getUdoObject()
  treeData.value = res?.data || []
  treeLoading.value = false

  if (data)
    currentUdo.value = data
}
initData()
async function initfieldData() {
  $table.value?.handleSearch()
}
const modalRef = ref(null)

// const columns = [
//   {
//     title: '字段编码',
//     key: 'code'
//   },
//   {
//     title: '字段名称',
//     key: 'name'
//   },

//   {
//     title: '操作',
//     key: 'actions',
//     render(row) {
//       return h(NSpace, null, {
//         default: () => [
//           h(
//             NButton,
//             {
//               size: 'small',
//               onClick: () => handleEdit(row)
//             },
//             { default: () => '编辑' }
//           ),
//           h(
//             NButton,
//             {
//               size: 'small',
//               type: 'error',
//               onClick: () => handleDelete(row)
//             },
//             { default: () => '删除' }
//           )
//         ]
//       })
//     }
//   }
// ]
const btnsColumns = [
  { title: '名称', key: 'name' },
  { title: '编码', key: 'code' },
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
        image: '图片',
      }
      return h(NTag, { type: 'info' }, { default: () => typeMap[row.field_type] || row.field_type })
    },
  },
  {
    title: '必填',
    key: 'is_required',
    render(row) {
      return h(NSwitch, {
        value: row.is_required,
        disabled: true,
        size: 'small',
      })
    },
  },
  {
    title: '描述',
    key: 'description',
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

      /*   h(
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
        ), */
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

function handleAddField() {
  modalRef.value?.handleOpen({
    action: 'add',
    title: '新增按钮',
    row: { object_id: currentUdo.value.id },
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
/* function handleDeleteBtn(id) {
  const d = $dialog.warning({
    content: '确定删除？',
    title: '提示',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        d.loading = true
        await api.getUdoFieldDel(id)
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
} */
</script>
