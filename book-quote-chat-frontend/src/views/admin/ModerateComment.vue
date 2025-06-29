<template>
  <div>

    <el-page-header
        title="后台"
        content="评论审核"
        @back="goAdmin"
    >
      <template #extra>
        <el-button @click="goToQuoteReview" type="primary" size="small">切换到金句审核</el-button>
      </template>
    </el-page-header>

    <el-form :inline="true" class="mb-2">
      <el-form-item>
        <el-input v-model="keyword" placeholder="搜索评论/用户/内容" @keyup.enter="fetchList" clearable/>
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
      </el-form-item>
      <el-form-item>
        <el-button type="success" :disabled="!multipleSelection.length" @click="handleBatchAudit('approved')">批量通过</el-button>
        <el-button type="danger" :disabled="!multipleSelection.length" @click="handleBatchAudit('rejected')">批量驳回</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border v-loading="loading" size="small" @selection-change="onSelectionChange">
      <el-table-column type="selection" width="45" />
      <el-table-column label="用户" width="120">
        <template #default="scope">
          <el-avatar :src="scope.row.avatar" size="small" class="mr-1"/>
          {{ scope.row.userName }}
        </template>
      </el-table-column>
      <el-table-column prop="content" label="内容" min-width="180">
        <template #default="scope">
          <span class="text-blue-500 cursor-pointer" @click="showDetail(scope.row)">
            {{ scope.row.content }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="created" label="时间" width="150">
        <template #default="scope">{{ formatDate(scope.row.created) }}</template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template #default="scope">
          <el-tag :type="scope.row.status==='approved'?'success':scope.row.status==='rejected'?'danger':'warning'">
            {{ statusLabel(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220">
        <template #default="scope">
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

    <!-- 评论详情弹窗 -->
    <el-dialog v-model="detailVisible" title="评论详情" width="400px">
      <div v-if="current">
        <p><b>用户：</b>{{ current.userName }}</p>
        <p><b>内容：</b>{{ current.content }}</p>
        <p><b>时间：</b>{{ formatDate(current.created) }}</p>
        <p><b>状态：</b>{{ statusLabel(current.status) }}</p>
        <p v-if="current.status==='rejected'"><b>驳回理由：</b>{{ current.rejectReason }}</p>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/api/request'

const router = useRouter()
const list = ref<any[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = 20
const keyword = ref('')
const status = ref('')

const detailVisible = ref(false)
const current = ref<any>(null)

const multipleSelection = ref<any[]>([])

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
  request.get('/api/admin/comments', {
    params: {
      offset: (page.value-1)*pageSize,
      limit: pageSize,
      keyword: keyword.value,
      status: status.value,
    }
  }).then(res => {
    list.value = res.data.list
    total.value = res.data.total
  }).finally(() => loading.value = false)
}

function handleAudit(row: any, auditStatus: 'approved' | 'rejected') {
  ElMessageBox.prompt(
      auditStatus === 'rejected' ? '请输入驳回理由' : '请确认通过',
      '审核',
      {
        inputType: auditStatus === 'rejected' ? 'textarea' : 'text',
        showCancelButton: true,
        inputPlaceholder: auditStatus === 'rejected' ? '请填写原因' : undefined
      }
  ).then(({ value }) => {
    request.post('/api/admin/comment/audit', {
      id: row.id,
      status: auditStatus,
      by: 'admin',
      reason: value || ''
    }).then(() => {
      ElMessage.success('操作成功')
      fetchList()
    })
  }).catch(() => {})
}
function handleDelete(row: any) {
  ElMessageBox.confirm('确定删除这条评论吗？', '提示', { type: 'warning' })
      .then(() => {
        request.post('/api/admin/comment/delete', { id: row.id })
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

function onSelectionChange(val: any[]) {
  multipleSelection.value = val
}

function handleBatchAudit(auditStatus: 'approved' | 'rejected') {
  if (multipleSelection.value.length === 0) return
  ElMessageBox.prompt(
    auditStatus === 'rejected' ? '请输入驳回理由' : `批量通过${multipleSelection.value.length}条`,
    '批量审核',
    {
      inputType: auditStatus === 'rejected' ? 'textarea' : 'text',
      showCancelButton: true,
      inputPlaceholder: auditStatus === 'rejected' ? '请填写原因' : undefined
    }
  ).then(({ value }) => {
    const ids = multipleSelection.value.map(item => item.id)
    request.post('/api/admin/comment/audit_batch', {
      ids,
      status: auditStatus,
      by: 'admin',
      reason: value || ''
    }).then(() => {
      ElMessage.success('批量操作成功')
      fetchList()
      multipleSelection.value = []
    })
  }).catch(() => {})
}

function goToQuoteReview() {
  router.push('/admin/moderate-quote')
}
function goAdmin() {
  router.push('/admin/')
}
fetchList()
</script>