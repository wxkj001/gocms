<template>
  <CommonPage>
    <template #action>
      <NButton v-permission="'AddUser'" type="primary" @click="handleAdd()">
        <i class="i-material-symbols:add mr-4 text-18" />
        创建新用户
      </NButton>
    </template>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1200"
      :columns="columns"
      :get-data="api.list"
    >
    <!--   <MeQueryItem label="用户名" :label-width="50">
        <n-input
          v-model:value="queryItems.username"
          type="text"
          placeholder="请输入用户名"
          clearable
        />
      </MeQueryItem>

      <MeQueryItem label="性别" :label-width="50">
        <n-select v-model:value="queryItems.gender" clearable :options="genders" />
      </MeQueryItem> -->

      <!--  <MeQueryItem label="状态" :label-width="50">
        <n-select
          v-model:value="queryItems.enable"
          clearable
          :options="[
            { label: '启用', value: 1 },
            { label: '停用', value: 0 },
          ]"
        />
      </MeQueryItem> -->
    </MeCrud>
  </CommonPage>
</template>

<script setup>
import { MeCrud, MeModal, MeQueryItem } from '@/components'
import { useCrud } from '@/composables'
import { withPermission } from '@/directives'
import { formatDateTime } from '@/utils'
import { NAvatar, NButton, NInput, NSwitch, NTag } from 'naive-ui'
import api from './api'

defineOptions({ name: 'configMgt' })

const $table = ref(null)
/** QueryBar筛选参数（可选） */
const queryItems = ref({})

onMounted(() => {
  $table.value?.handleSearch()
})
const columns = [
  { title: '设置名称', key: 'config_name', width: 50, ellipsis: { tooltip: true } },
  { title: '设置项', key: 'config_key', width: 50 },
  {
    title: '参数',
    key: 'config_value',
    width: 20,
    render: (row) => {
      switch (row.config_type) {
        case 'input':
          return h(
            NInput,
            {
              size: 'small',
              value: row.config_value,
              placeholder: '请输入参数',
              onUpdateValue: (value) => {
                row.config_value = value
              },
              onBlur: () => {
                handleSaveConfig(row)
              },
            },
          )
        case 'switch':
          return h(
            NSwitch,
            {
              size: 'small',
              rubberBand: false,
              value: row.config_value,
              loading: !!row.enableLoading,
              checkedValue: '1',
              uncheckedValue: '0',
              onUpdateValue: () => {},
            },
            {
              checked: () => '启用',
              unchecked: () => '停用',
            },
          )
        case 'select':
          return '下拉框'
      }
    },

  },
]
async function handleSaveConfig(row) {
  try {
    await api.update(row)
    $message.success('操作成功')
    $table.value?.handleSearch()
  }
  catch (error) {
    console.error(error)
  }
}
/*
const genders = [
  { label: '男', value: 1 },
  { label: '女', value: 2 },
]
const roles = ref([])
api.getAllRoles().then(({ data = [] }) => (roles.value = data))

const {
  modalRef,
  modalFormRef,
  modalForm,
  modalAction,
  handleAdd,
  handleDelete,
  handleOpen,
  handleSave,
} = useCrud({
  name: '用户',
  initForm: { enable: true },
  doCreate: api.create,
  doDelete: api.delete,
  doUpdate: api.update,
  refresh: () => $table.value?.handleSearch(),
})

async function handleEnable(row) {
  row.enableLoading = true
  try {
    await api.update({ id: row.id, enable: !row.enable })
    row.enableLoading = false
    $message.success('操作成功')
    $table.value?.handleSearch()
  }
  catch (error) {
    console.error(error)
    row.enableLoading = false
  }
}

function handleOpenRolesSet(row) {
  const roleIds = row.roles.map(item => item.id)
  handleOpen({
    action: 'setRole',
    title: '分配角色',
    row: { id: row.id, username: row.username, roleIds },
    onOk: onSave,
  })
}

function onSave() {
  if (modalAction.value === 'setRole') {
    return handleSave({
      api: () => api.update(modalForm.value),
      cb: () => $message.success('分配成功'),
    })
  }
  else if (modalAction.value === 'reset') {
    return handleSave({
      api: () => api.resetPwd(modalForm.value.id, modalForm.value),
      cb: () => $message.success('密码重置成功'),
    })
  }
  handleSave()
} */
</script>
