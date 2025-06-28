<template>
  <div>
    <el-button type="primary" @click="showAddUser" class="mb-4">新建用户</el-button>
    <el-page-header content="用户管理"/>
    <el-form :inline="true" class="mb-4">
      <el-form-item label="用户ID">
        <el-input v-model="searchId" placeholder="用户ID"></el-input>
      </el-form-item>
      <el-form-item label="昵称">
        <el-input v-model="searchName" placeholder="昵称"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="fetchList">搜索</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="userList" border v-loading="loading" style="width:100%">
      <el-table-column prop="id" label="ID" width="110" />
      <el-table-column prop="name" label="昵称" width="130" />
      <el-table-column prop="email" label="邮箱" width="170" />
      <el-table-column prop="status" label="状态" width="80">
        <template #default="scope">
          <el-tag v-if="scope.row.status === 'banned'" type="danger">封禁</el-tag>
          <el-tag v-else>正常</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250">
        <template #default="scope">
          <el-button size="small" @click="showEdit(scope.row)">编辑</el-button>
          <el-button v-if="scope.row.status !== 'banned'" size="small" type="danger" @click="banUser(scope.row)">封禁</el-button>
          <el-button v-else size="small" type="success" @click="unbanUser(scope.row)">解封</el-button>
          <el-button size="small" @click="viewReports(scope.row)">举报</el-button>
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

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editVisible" title="编辑用户" width="400px">
      <el-form :model="editUser">
        <el-form-item label="昵称">
          <el-input v-model="editUser.name" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="editUser.email" />
        </el-form-item>
        <!-- 可扩展更多字段 -->
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="saveUser">保存</el-button>
      </template>
    </el-dialog>

    <!-- 举报弹窗 -->
    <el-dialog v-model="reportVisible" title="举报记录" width="500px">
      <el-table :data="reports" border>
        <el-table-column prop="id" label="举报ID" width="100"/>
        <el-table-column prop="reason" label="原因"/>
        <el-table-column prop="created" label="时间" width="160">
          <template #default="scope">
            {{ formatDate(scope.row.created) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80"/>
      </el-table>
      <template #footer>
        <el-button @click="reportVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
  <!-- 新建用户弹窗 -->
  <el-dialog v-model="addVisible" title="新建用户" width="400px">
    <el-form :model="addUser">
      <el-form-item label="昵称">
        <el-input v-model="addUser.name" />
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input v-model="addUser.email" />
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="addUser.password" type="password" />
      </el-form-item>
      <!-- 可扩展更多字段 -->
    </el-form>
    <template #footer>
      <el-button @click="addVisible = false">取消</el-button>
      <el-button type="primary" @click="saveAddUser">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/api/request'

const userList = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 20

const searchId = ref('')
const searchName = ref('')

const loading = ref(false)

const editVisible = ref(false)
const editUser = ref<any>({})

const reportVisible = ref(false)
const reports = ref<any[]>([])

function fetchList() {
  loading.value = true
  request.get('/api/admin/users', {
    params: {
      offset: (page.value-1)*pageSize,
      limit: pageSize,
      id: searchId.value,
      name: searchName.value
    }
  }).then(res => {
    userList.value = res.data.list
    total.value = res.data.total
  }).catch(() => {
    ElMessage.error('获取用户失败')
  }).finally(() => {
    loading.value = false
  })
}

function onPageChange(p:number) {
  page.value = p
  fetchList()
}

function showEdit(row: any) {
  editUser.value = { ...row }
  editVisible.value = true
}

function saveUser() {
  request.post('/api/admin/update_user', {
    id: editUser.value.id,
    name: editUser.value.name,
    email: editUser.value.email
  }).then(() => {
    ElMessage.success('修改成功')
    editVisible.value = false
    fetchList()
  }).catch(() => {
    ElMessage.error('修改失败')
  })
}

function banUser(row: any) {
  ElMessageBox.confirm('确定要封禁该用户？', '警告').then(() => {
    request.post('/api/admin/ban_user', { id: row.id }).then(() => {
      ElMessage.success('已封禁')
      fetchList()
    })
  })
}
function unbanUser(row: any) {
  ElMessageBox.confirm('确定要解封该用户？', '提示').then(() => {
    request.post('/api/admin/unban_user', { id: row.id }).then(() => {
      ElMessage.success('已解封')
      fetchList()
    })
  })
}
function viewReports(row: any) {
  request.get('/api/admin/user_reports', {
    params: { userId: row.id }
  }).then(res => {
    reports.value = res.data.list
    reportVisible.value = true
  }).catch(() => {
    ElMessage.error('获取举报失败')
  })
}

function formatDate(ts: number) {
  if (!ts) return '-'
  if (ts > 1e10) ts = Math.floor(ts/1000)
  return new Date(ts*1000).toLocaleString()
}

const addVisible = ref(false)
const addUser = ref<any>({
  name: '',
  email: '',
  password: ''
})

function showAddUser() {
  addUser.value = { name: '', email: '', password: '' }
  addVisible.value = true
}

function saveAddUser() {
  if (!addUser.value.name || !addUser.value.email || !addUser.value.password) {
    ElMessage.error('请填写完整信息')
    return
  }
  request.post('/api/admin/create_user', addUser.value).then(() => {
    ElMessage.success('添加成功')
    addVisible.value = false
    fetchList()
  }).catch(() => {
    ElMessage.error('添加失败')
  })
}
onMounted(fetchList)
</script>