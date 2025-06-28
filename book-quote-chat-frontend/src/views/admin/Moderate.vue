<template>
  <div>
    <el-page-header content="评论/金句审核"/>
    <el-form :inline="true" class="mb-4">
      <el-form-item label="类型">
        <el-radio-group v-model="targetType" size="small" @change="onTypeChange">
          <el-radio-button label="quote">金句</el-radio-button>
          <el-radio-button label="comment">评论</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="关键词">
        <el-input v-model="keyword" placeholder="内容/用户ID"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="fetchList">搜索</el-button>
      </el-form-item>
      <el-form-item>
        <el-button @click="batchApprove" :disabled="!multipleSelection.length">批量通过</el-button>
        <el-button type="danger" @click="batchDelete" :disabled="!multipleSelection.length">批量删除</el-button>
      </el-form-item>
    </el-form>

    <el-table
        :data="dataList"
        border
        style="width:100%"
        @selection-change="multipleSelection = $event"
        highlight-current-row
        v-loading="loading"
    >
      <el-table-column type="selection" width="40"/>
      <el-table-column label="类型" width="80">
        <template #default="scope">
          <span v-if="scope.row.targetType === 'quote'">金句</span>
          <span v-else-if="scope.row.targetType === 'comment'">评论</span>
          <span v-else>未知</span>
        </template>
      </el-table-column>
      <el-table-column label="内容" prop="content" min-width="180" :show-overflow-tooltip="true">
        <template #default="scope">
        <span
            class="table-content-link"
            style="cursor:pointer;color:#409eff;text-decoration:underline;"
            @click="showDetail(scope.row)"
        >
          {{ scope.row.content }}
        </span>
        </template>
      </el-table-column>
      <el-table-column label="用户" prop="userName" width="100"/>
      <el-table-column label="时间" prop="created" width="140">
        <template #default="scope">
          {{ formatDate(scope.row.created) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180">
        <template #default="scope">
          <el-button size="small" type="success" @click="approve(scope.row)">通过</el-button>
          <el-button size="small" type="danger" @click="deleteItem(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
        class="mt-4"
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="onPageChange"
    />

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="内容详情" width="400px">
      <div v-if="current">
        <p><b>类型：</b>
          <span v-if="current.targetType === 'quote'">金句</span>
          <span v-else-if="current.targetType === 'comment'">评论</span>
        </p>
        <p><b>内容：</b>{{ current.content }}</p>
        <p><b>用户：</b>{{ current.userName }} ({{ current.userId }})</p>
        <p v-if="current.book"><b>金句来源：</b>{{ current.book }}</p>
        <p><b>时间：</b>{{ formatDate(current.created) }}</p>

        <el-divider>举报历史</el-divider>
        <div v-if="reportLoading">
          <el-skeleton :rows="2" animated/>
        </div>
        <div v-else-if="reports.length === 0" class="text-gray-400 text-sm">无举报记录</div>
        <el-table v-else :data="reports" size="small" border style="width:100%">
          <el-table-column prop="reporterName" label="举报人" width="90"/>
          <el-table-column prop="reason" label="原因" min-width="80"/>
          <el-table-column prop="created" label="时间" width="150">
            <template #default="scope">{{ formatDate(scope.row.created) }}</template>
          </el-table-column>
          <el-table-column prop="status" label="处理状态" width="90">
            <template #default="scope">
              <el-tag v-if="scope.row.status==='pending'">待处理</el-tag>
              <el-tag v-else-if="scope.row.status==='processed'" type="success">已处理</el-tag>
              <span v-else>{{ scope.row.status }}</span>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/api/request'


const dataList = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 20

const targetType = ref('quote') // 默认金句
function onTypeChange() {
  page.value = 1
  fetchList()
}
const keyword = ref('')
const multipleSelection = ref<any[]>([])

const detailVisible = ref(false)
const current = ref<any>(null)
const loading = ref(false)

const reports = ref<any[]>([])
const reportLoading = ref(false)

function formatDate(ts: number) {
  if (!ts) return '-'
  if (ts > 1e10) ts = Math.floor(ts/1000)
  return new Date(ts*1000).toLocaleString()
}

function fetchList() {
  loading.value = true
  const url = targetType.value === 'comment'
    ? '/api/admin/comments'
    : '/api/admin/quotes'
  request.get(url, {
    params: {
      offset: (page.value-1)*pageSize,
      limit: pageSize,
      keyword: keyword.value,
    }
  }).then(res => {
    dataList.value = res.data.list
    total.value = res.data.total
  }).catch(() => {
    ElMessage.error('获取列表失败')
  }).finally(() => {
    loading.value = false
  })
}

function onPageChange(p:number) {
  page.value = p
  fetchList()
}

function showDetail(row:any) {
  current.value = row
  detailVisible.value = true
  // 自动查举报记录
  reports.value = []
  reportLoading.value = true
  request.get('/api/admin/reports', { params: { targetId: row.id } })
    .then(res => {
      reports.value = res.data.list || []
    })
    .catch(() => {
      reports.value = []
    })
    .finally(() => {
      reportLoading.value = false
    })
}

function approve(row:any) {
  ElMessageBox.confirm('确定通过该条内容审核？', '提示').then(() => {
    request.post('/api/admin/audit_comment', {
      adminId: 'admin-demo',
      commentId: row.id,
      action: 'approve',
      detail: ''
    }).then(() => {
      ElMessage.success('审核通过！')
      fetchList()
    })
  })
}
function deleteItem(row:any) {
  ElMessageBox.confirm('确定要删除吗？', '警告').then(() => {
    request.post('/api/admin/audit_comment', {
      adminId: 'admin-demo',
      commentId: row.id,
      action: 'delete',
      detail: ''
    }).then(() => {
      ElMessage.success('已删除')
      fetchList()
    })
  })
}
function batchApprove() {
  if (!multipleSelection.value.length) return
  ElMessageBox.confirm('确定批量通过所选内容？', '提示').then(() => {
    Promise.all(multipleSelection.value.map(row =>
        request.post('/api/admin/audit_comment', {
          adminId: 'admin-demo',
          commentId: row.id,
          action: 'approve',
          detail: ''
        })
    )).then(() => {
      ElMessage.success('操作完成')
      fetchList()
    })
  })
}
function batchDelete() {
  if (!multipleSelection.value.length) return
  ElMessageBox.confirm('确定批量删除？', '警告').then(() => {
    Promise.all(multipleSelection.value.map(row =>
        request.post('/api/admin/audit_comment', {
          adminId: 'admin-demo',
          commentId: row.id,
          action: 'delete',
          detail: ''
        })
    )).then(() => {
      ElMessage.success('已删除')
      fetchList()
    })
  })
}

onMounted(fetchList)
</script>