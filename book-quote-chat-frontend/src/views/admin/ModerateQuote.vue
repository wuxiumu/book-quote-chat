<template>
  <div>
    <el-page-header content="金句审核" class="mb-4">
      <template #extra>
        <el-button @click="goToComment" type="primary" size="small">去评论审核</el-button>
      </template>
    </el-page-header>
    <el-form :inline="true" class="mb-2">
      <el-form-item>
        <el-input v-model="keyword" placeholder="搜索金句/用户/书名/标签" @keyup.enter="fetchList" clearable/>
      </el-form-item>
      <el-form-item>
        <el-select v-model="status" placeholder="筛选状态" clearable style="width:120px" @change="fetchList">
          <el-option label="全部" value=""></el-option>
          <el-option label="待审核" value="pending"></el-option>
          <el-option label="已通过" value="approved"></el-option>
          <el-option label="已驳回" value="rejected"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button @click="fetchList" type="primary">搜索</el-button>
        <el-button @click="showAdd">新增金句</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border v-loading="loading" size="small">
      <el-table-column prop="text" label="内容" min-width="180">
        <template #default="scope">
          <span class="text-blue-500 cursor-pointer" @click="showDetail(scope.row)">
            {{ scope.row.text }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="user" label="用户" width="90"/>
      <el-table-column prop="book" label="书名" width="120"/>
      <el-table-column prop="tags" label="标签" min-width="80">
        <template #default="scope">
          <el-tag v-for="tag in scope.row.tags" :key="tag" size="small">{{ tag }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template #default="scope">
          <el-tag :type="scope.row.status==='approved'?'success':scope.row.status==='rejected'?'danger':'warning'">
            {{ statusLabel(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="210">
        <template #default="scope">
          <el-button size="small" @click="showEdit(scope.row)">编辑</el-button>
          <el-button size="small" type="primary" v-if="scope.row.status==='pending'" @click="handleAudit(scope.row, 'approved')">通过</el-button>
          <el-button size="small" type="danger" v-if="scope.row.status==='pending'" @click="handleAudit(scope.row, 'rejected')">驳回</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        @current-change="fetchList"
        layout="total, prev, pager, next"
        class="mt-2"
    />

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="editVisible" :title="editForm.id ? '编辑金句' : '新增金句'" width="480px">
      <el-form :model="editForm" label-width="60px">
        <el-form-item label="内容">
          <el-input v-model="editForm.text" type="textarea" :rows="3"/>
        </el-form-item>
        <el-form-item label="书名">
          <el-input v-model="editForm.book"/>
        </el-form-item>
        <el-form-item label="标签">
          <el-select v-model="editForm.tags" multiple filterable allow-create>
            <el-option v-for="t in tagOptions" :key="t" :label="t" :value="t"/>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="saveQuote">保存</el-button>
      </template>
    </el-dialog>

    <!-- 金句详情弹窗 -->
    <el-dialog v-model="detailVisible" title="金句详情" width="400px">
      <div v-if="current">
        <p><b>内容：</b>{{ current.text }}</p>
        <p><b>用户：</b>{{ current.user }}</p>
        <p><b>书名：</b>{{ current.book }}</p>
        <p><b>标签：</b>
          <el-tag v-for="tag in current.tags" :key="tag" size="small">{{ tag }}</el-tag>
        </p>
        <p><b>状态：</b>{{ statusLabel(current.status) }}</p>
        <p><b>创建时间：</b>{{ formatDate(current.created) }}</p>
        <p v-if="current.status==='rejected'"><b>驳回理由：</b>{{ current.rejectReason }}</p>
        <p v-if="current.status==='approved'"><b>审核人：</b>{{ current.auditBy }}</p>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/api/request'

const status = ref('')
const list = ref<any[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = 20
const keyword = ref('')

const editVisible = ref(false)
const editForm = reactive<any>({})
const detailVisible = ref(false)
const current = ref<any>(null)

// 标签选项（可根据实际调整或从接口获取）
const tagOptions = ['情感', '哲理', '成长', '生活', '青春', '孤独']

function statusLabel(status: string) {
  if (status === 'approved') return '已通过'
  if (status === 'rejected') return '已驳回'
  return '待审核'
}

function formatDate(ts: number) {
  if (!ts) return ''
  return new Date(ts * 1000).toLocaleString()
}

function fetchList() {
  loading.value = true
  request.get('/api/admin/quotes', {
    params: {
      offset: (page.value - 1) * pageSize,
      limit: pageSize,
      keyword: keyword.value,
      status: status.value,
      sort: 'created'
    }
  }).then(res => {
    list.value = res.data.list
    total.value = res.data.total
  }).finally(() => loading.value = false)
}

function showAdd() {
  Object.assign(editForm, { id: '', text: '', book: '', tags: [] })
  editVisible.value = true
}
function showEdit(row: any) {
  Object.assign(editForm, { ...row })
  editVisible.value = true
}
function saveQuote() {
  const url = editForm.id ? '/api/admin/quote/edit' : '/api/admin/quote/add'
  request.post(url, editForm).then(() => {
    ElMessage.success('保存成功')
    editVisible.value = false
    fetchList()
  }).catch(() => ElMessage.error('保存失败'))
}
function handleAudit(row: any, status: 'approved' | 'rejected') {
  ElMessageBox.prompt(status === 'rejected' ? '请输入驳回理由' : '请确认通过', '审核', {
    inputType: status === 'rejected' ? 'textarea' : 'text',
    showCancelButton: true,
    inputPlaceholder: status === 'rejected' ? '请填写原因' : undefined
  }).then(({ value }) => {
    request.post('/api/admin/quote/audit', {
      id: row.id,
      status,
      by: 'admin', // 真实项目请替换为当前管理员
      reason: value || ''
    }).then(() => {
      ElMessage.success('操作成功')
      fetchList()
    })
  }).catch(() => {})
}
function handleDelete(row: any) {
  ElMessageBox.confirm('确定删除这条金句吗？', '提示', { type: 'warning' })
      .then(() => {
        request.post('/api/admin/quote/delete', { id: row.id })
            .then(() => {
              ElMessage.success('已删除')
              fetchList()
            })
      })
}
function showDetail(row: any) {
  current.value = row
  detailVisible.value = true
}
function goToComment() {
  // 跳转到评论审核页（根据你的路由配置调整）
  window.location.href = '/admin/moderate-comment'
}

fetchList()
</script>