<template>
  <div>
    <el-page-header content="管理操作日志"/>
    <el-form :inline="true" class="mb-4">
      <el-form-item label="管理员">
        <el-input v-model="adminId" placeholder="ID/昵称"></el-input>
      </el-form-item>
      <el-form-item label="类型">
        <el-select v-model="action" placeholder="全部">
          <el-option label="全部" value=""></el-option>
          <el-option label="通过" value="approve"></el-option>
          <el-option label="驳回" value="reject"></el-option>
          <el-option label="删除" value="delete"></el-option>
          <el-option label="登录" value="login"></el-option>
          <!-- 可扩展其它类型 -->
        </el-select>
      </el-form-item>
      <el-form-item label="对象">
        <el-select v-model="target" placeholder="全部">
          <el-option label="全部" value=""></el-option>
          <el-option label="评论" value="comment"></el-option>
          <el-option label="金句" value="quote"></el-option>
          <el-option label="用户" value="user"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="目标ID">
        <el-input v-model="targetId" placeholder="内容/用户ID"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="fetchList">搜索</el-button>
      </el-form-item>
    </el-form>

    <el-table
        :data="logList"
        border
        style="width:100%"
        v-loading="loading"
    >
      <el-table-column label="时间" width="170">
        <template #default="scope">
          {{ formatDate(scope.row.created) }}
        </template>
      </el-table-column>
      <el-table-column label="管理员" prop="adminId" width="120"/>
      <el-table-column label="操作" width="80">
        <template #default="scope">
          <span v-if="scope.row.action==='approve'">通过</span>
          <span v-else-if="scope.row.action==='reject'">驳回</span>
          <span v-else-if="scope.row.action==='delete'">删除</span>
          <span v-else-if="scope.row.action==='login'">登录</span>
          <span v-else>{{ scope.row.action }}</span>
        </template>
      </el-table-column>
      <el-table-column label="对象类型" width="80">
        <template #default="scope">
          <span v-if="scope.row.target==='comment'">评论</span>
          <span v-else-if="scope.row.target==='quote'">金句</span>
          <span v-else-if="scope.row.target==='user'">用户</span>
          <span v-else>{{ scope.row.target }}</span>
        </template>
      </el-table-column>
      <el-table-column label="目标ID" prop="targetId" width="100"/>
      <el-table-column label="说明" prop="detail" min-width="140" :show-overflow-tooltip="true"/>
    </el-table>
    <el-pagination
        class="mt-4"
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="onPageChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/api/request'

const logList = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 20

const adminId = ref('')
const action = ref('')
const target = ref('')
const targetId = ref('')

const loading = ref(false)

function formatDate(ts: number) {
  if (!ts) return '-'
  if (ts > 1e10) ts = Math.floor(ts/1000)
  return new Date(ts*1000).toLocaleString()
}

function fetchList() {
  loading.value = true
  request.get('/api/admin/logs', {
    params: {
      offset: (page.value-1)*pageSize,
      limit: pageSize,
      adminId: adminId.value,
      action: action.value,
      target: target.value,
      targetId: targetId.value
    }
  }).then(res => {
    logList.value = res.data.list
    total.value = res.data.total
  }).catch(() => {
    ElMessage.error('获取日志失败')
  }).finally(() => {
    loading.value = false
  })
}

function onPageChange(p:number) {
  page.value = p
  fetchList()
}

onMounted(fetchList)
</script>