<template>
  <CommonPage>
    <template #action>
      <NButton type="primary" @click="handleAdds()">
        <i class="i-material-symbols:add mr-4 text-18" />
        新增数据
      </NButton>
    </template>

    <MeCrud ref="$table" v-model:query-items="queryItems" :scroll-x="1200" :columns="columns" :get-data="getList">
      <!-- <MeQueryItem label="角色名" :label-width="50">
        <n-input v-model:value="queryItems.name" type="text" placeholder="请输入角色名" clearable />
      </MeQueryItem> -->
    </MeCrud>
    <MeModal ref="modalRef" width="100%" :modal-style="{ height: '100vh', 'align-self': 'auto' }">

      <n-scrollbar style="max-height: 80vh;">
        <n-form ref="modalFormRef" label-placement="left" label-align="left" :label-width="80" :model="modalForm">
          <n-form-item v-for="(item, key) in formFields" :key="key" :label="item.name" :path="item.code"
            :rule="ruleFun(item)">
            <AiEditor v-model:model-value="modalForm[item.code]" v-if="item.field_type === 'richtext'"></AiEditor>
            <n-input v-if="item.field_type === 'string'" v-model:value="modalForm[item.code]" />
            <n-input-number v-if="item.field_type === 'number'" v-model:value="modalForm[item.code]" clearable />
            <NSwitch v-if="item.field_type === 'boolean'" v-model:value="modalForm[item.code]" />
            <n-select v-if="item.field_type === 'enum'" v-model:value="modalForm[item.code]"
              :options="item.enum_options" />
            <n-date-picker v-if="item.field_type === 'datetime'" v-model:formatted-value="modalForm[item.code]" value-format="yyyy-MM-dd HH:mm:ss"
             type="datetime" clearable />
            <n-date-picker v-if="item.field_type === 'date'" v-model:value="modalForm[item.code]" type="date"
              clearable />
            <n-upload v-if="item.field_type === 'file'" multiple directory-dnd
              action="https://www.mocky.io/v2/5e4bafc63100007100d8b70f" :max="5">
              <n-upload-dragger>
                <div style="margin-bottom: 12px">
                  <n-icon size="48" :depth="3">
                    <ArchiveIcon />
                  </n-icon>
                </div>
                <n-text style="font-size: 16px">
                  点击或者拖动文件到该区域来上传
                </n-text>
              </n-upload-dragger>
            </n-upload>
          </n-form-item>
        </n-form>
      </n-scrollbar>
    </MeModal>
  </CommonPage>
</template>

<script setup>
import { MeCrud, MeModal, MeQueryItem } from '@/components'
import { useCrud } from '@/composables'
import { NButton, NSwitch } from 'naive-ui'
import { useRoute } from 'vue-router'; // 导入 useRoute
import api from './api'
import AiEditor from './components/aieditor.vue'

const route = useRoute();
// 通过 route.path 获取最后一级
const lastSegmentFromPath = computed(() => {
  const path = route.path;
  // 移除末尾的斜杠（如果有）
  const cleanPath = path.endsWith('/') && path.length > 1 ? path.slice(0, -1) : path;
  const segments = cleanPath.split('/');
  return segments[segments.length - 1] || '/'; // 如果是根路径 "/"，segments会是["", ""], 取第二个会是空字符串，所以这里特判
});

defineOptions({ name: 'RoleMgt' })

const $table = ref(null)
const formFields = ref([])
/** QueryBar筛选参数（可选） */
const queryItems = ref({})

const { modalRef, modalFormRef, modalAction, modalForm, handleAdd, handleDelete, handleEdit }
  = useCrud({
    name: '数据',
    doCreate: (data)=>{api.create(lastSegmentFromPath.value,data)},
    doDelete: (data)=>{api.delete(lastSegmentFromPath.value,data.id)},
    doUpdate: (data)=>{api.update(lastSegmentFromPath.value,data.id,data)},
    initForm: {},
    refresh: (_, keepCurrentPage) => $table.value?.handleSearch(keepCurrentPage),
  })
onMounted(() => {
  $table.value?.handleSearch()
})
function handleAdds() {
  handleAdd()
}
const columns = ref([])
api.fields(lastSegmentFromPath.value, {}).then(({ data = [] }) => {
  formFields.value = data
  for (const item of data) {
    columns.value.push({
      title: item.name,
      key: item.code,
    })
  }
  columns.value.push({
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
            disabled: row.code === 'SUPER_ADMIN',
            onClick: () => handleEdit(row),
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
            disabled: row.code === 'SUPER_ADMIN',
            onClick: () => handleDelete(row.id),
          },
          {
            default: () => '删除',
            icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
          },
        ),
      ]
    },
  })
})

async function getList(data) {
  return await api.list(lastSegmentFromPath.value, data)
}
function ruleFun(item) {
  switch (item.field_type) {
    case 'string':
      return {
        required: item.is_required,
        message: `请输入${item.name}`,
        trigger: ['blur', 'input'],
      }
    /* case 'datetime':
      return {
        type: 'number',
        required: item.is_required,
        message: `请输入${item.name}`,
        trigger: ['blur', 'input'],
      } */
    case 'date':
      return {
        type: 'number',
        required: item.is_required,
        message: `请输入${item.name}`,
        trigger: ['blur', 'input'],
      }
    case 'number':
      return {
        type: 'number',
        required: item.is_required,
        message: `请输入${item.name}`,
        trigger: ['blur', 'change'],
      }
    case 'boolean':
      return {
        type: 'boolean',
        required: item.is_required,
        message: `请输入${item.name}`,
        trigger: ['blur', 'change'],
      }
    default:
      return {}
  }
}
</script>
