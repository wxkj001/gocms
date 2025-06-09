<template>
  <div>
    <n-space vertical :size="12">
      <h3>数据对象</h3>
      <div class="flex">
        <n-input v-model:value="pattern" placeholder="搜索" clearable />
        <NButton class="ml-12" type="primary" @click="handleAdd()">
          <i class="i-material-symbols:add mr-4 text-14" />
          新增
        </NButton>
      </div>

      <n-tree
        :show-irrelevant-nodes="false"
        :pattern="pattern"
        :data="treeData"
        :selected-keys="[currentUdo?.code]"
        :render-prefix="renderPrefix"
        :render-suffix="renderSuffix"
        :on-update:selected-keys="onSelect"
        key-field="code"
        label-field="name"

        block-line default-expand-all
      />
    </n-space>

    <ResAddOrEdit ref="modalRef" :menus="treeData" @refresh="(data) => emit('refresh', data)" />
  </div>
</template>

<script setup>
import { NButton } from 'naive-ui'
import { withModifiers } from 'vue'
import api from '../api'
import ResAddOrEdit from './ResAddOrEdit.vue'

defineProps({
  treeData: {
    type: Array,
    default: () => [],
  },
  currentUdo: {
    type: Object,
    default: () => null,
  },
})
const emit = defineEmits(['refresh', 'update:currentUdo'])

const pattern = ref('')

const modalRef = ref(null)
async function handleAdd(data = {}) {
  modalRef.value?.handleOpen({
    action: 'add',
    title: '新增对象',
    row: data,
    okText: '保存',
  })
}

function onSelect(keys, option, { node }) {
  // console.log('currentUdo', keys, option, node)

  emit('update:currentUdo', node)
}

function renderPrefix({ option }) {
  return h('i', { class: `${option.icon}?mask text-16` })
}

function renderSuffix({ option }) {
  return [
    h(
      NButton,
      {
        text: true,
        type: 'error',
        size: 'tiny',
        style: 'margin-left: 12px;',
        onClick: withModifiers(() => handleDelete(option), ['stop']),
      },
      { default: () => '删除' },
    ),
  ]
}

function handleDelete(item) {
  $dialog.confirm({
    content: `确认删除【${item.name}】？数据无价请确认备份后再删除！`,
    async confirm() {
      try {
        $message.loading('正在删除', { key: 'deleteMenu' })
        await api.getUdoObjectDel(item.id)
        $message.success('删除成功', { key: 'deleteMenu' })
        emit('refresh')
        emit('update:currentUdo', null)
      }
      catch (error) {
        console.error(error)
        $message.destroy('deleteMenu')
      }
    },
  })
}
</script>
