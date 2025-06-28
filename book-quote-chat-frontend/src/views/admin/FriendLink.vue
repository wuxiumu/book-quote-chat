<template>
  <div>
    <el-page-header content="友链 / 配置管理"/>
    <el-button type="primary" @click="showAdd" class="mb-4">新增</el-button>
    <el-upload
        class="ml-2"
        action=""
        :before-upload="beforeImport"
        :show-file-list="false"
        accept=".json"
    >
      <el-button>批量导入</el-button>
    </el-upload>
    <el-button @click="exportConfig" class="ml-2">导出全部</el-button>
    <el-table :data="list" border style="width: 100%" class="mt-4" v-loading="loading">
      <el-table-column prop="name" label="名称" width="120"/>
      <el-table-column prop="value" label="链接/内容" min-width="200"/>
      <el-table-column prop="desc" label="备注" min-width="120"/>
      <el-table-column label="操作" width="180">
        <template #default="scope">
          <el-button size="small" @click="showEdit(scope.row)">编辑</el-button>
          <el-button size="small" type="danger" @click="delItem(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="editVisible" :title="editForm.id ? '编辑' : '新增'" width="400px">
      <el-form :model="editForm" label-width="60px">
        <el-form-item label="名称">
          <el-input v-model="editForm.name"/>
        </el-form-item>
        <el-form-item label="链接">
          <el-input v-model="editForm.value"/>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="editForm.desc"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="saveItem">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/api/request'

const list = ref<any[]>([])
const loading = ref(false)
const editVisible = ref(false)
const editForm = ref<any>({ id: '', name: '', value: '', desc: '' })

function fetchList() {
  loading.value = true
  request.get('/api/admin/config')
      .then(res => list.value = res.data)
      .catch(() => ElMessage.error('获取失败'))
      .finally(() => loading.value = false)
}

function showAdd() {
  editForm.value = { id: '', name: '', value: '', desc: '' }
  editVisible.value = true
}

function showEdit(row: any) {
  editForm.value = { ...row }
  editVisible.value = true
}

function saveItem() {
  const data = { ...editForm.value }
  const url = data.id ? '/api/admin/config/update' : '/api/admin/config/add'
  request.post(url, data)
      .then(() => {
        ElMessage.success('保存成功')
        editVisible.value = false
        fetchList()
      })
      .catch(() => ElMessage.error('保存失败'))
}

function delItem(row: any) {
  ElMessageBox.confirm('确定要删除吗？', '警告').then(() => {
    request.get('/api/admin/config/delete', { params: { id: row.id } })
        .then(() => {
          ElMessage.success('已删除')
          fetchList()
        })
  })
}

function beforeImport(file: File) {
  const reader = new FileReader()
  reader.onload = function(e) {
    try {
      const arr = JSON.parse(e.target?.result as string)
      if (!Array.isArray(arr)) throw new Error()
      request.post('/api/admin/config/import', arr).then(() => {
        ElMessage.success('导入成功')
        fetchList()
      }).catch(() => ElMessage.error('导入失败'))
    } catch {
      ElMessage.error('文件格式错误')
    }
  }
  reader.readAsText(file)
  return false // 阻止默认上传
}

function exportConfig() {
  const dataStr = JSON.stringify(list.value, null, 2)
  const blob = new Blob([dataStr], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'config.json'
  a.click()
  URL.revokeObjectURL(url)
}

onMounted(fetchList)
</script>