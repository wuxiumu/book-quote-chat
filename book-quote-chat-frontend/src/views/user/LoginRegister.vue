<template>
  <div class="flex min-h-screen items-center justify-center bg-gradient-to-br from-blue-100 to-purple-100">
    <el-card class="w-96 shadow-xl">
      <div class="text-center text-2xl font-bold mb-6">
        {{ tab==='login' ? '用户登录' : '注册新账号' }}
      </div>
      <el-tabs v-model="tab" stretch>
        <el-tab-pane label="登录" name="login">
          <el-form :model="loginForm" @submit.native.prevent="onLogin">
            <el-form-item prop="username">
              <el-input v-model="loginForm.username" placeholder="用户名/手机号" autofocus/>
            </el-form-item>
            <el-form-item prop="password">
              <el-input v-model="loginForm.password" placeholder="密码" type="password"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" class="w-full" @click="onLogin" :loading="loading">登录</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="注册" name="register">
          <el-form :model="registerForm" @submit.native.prevent="onRegister">
            <el-form-item prop="name">
              <el-input v-model="registerForm.name" placeholder="用户名"/>
            </el-form-item>
            <el-form-item prop="password">
              <el-input v-model="registerForm.password" placeholder="密码" type="password"/>
            </el-form-item>
            <el-form-item prop="email">
              <el-input v-model="registerForm.email" placeholder="邮箱"/>
            </el-form-item>
            <el-form-item label="头像">
              <div class="flex gap-2">
                <el-avatar
                  v-for="img in avatarList"
                  :key="img"
                  :src="img"
                  :class="registerForm.avatar === img ? 'ring-2 ring-blue-500' : 'opacity-70 hover:opacity-100 cursor-pointer'"
                  size="large"
                  @click="registerForm.avatar = img"
                />
              </div>
            </el-form-item>
            <el-form-item>
              <el-button type="success" class="w-full" @click="onRegister" :loading="loading">注册</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/api/request'
import { useRouter } from 'vue-router'

const router = useRouter()
const tab = ref('login')
const loading = ref(false)

const loginForm = ref({ username: '', password: '' })
const registerForm = ref({ name: '', password: '', email: '', avatar: '', group: 'user' })

const avatarList = [
  'https://randomuser.me/api/portraits/men/32.jpg',
  'https://randomuser.me/api/portraits/women/44.jpg',
  'https://randomuser.me/api/portraits/men/11.jpg',
  'https://randomuser.me/api/portraits/women/68.jpg',
  'https://randomuser.me/api/portraits/men/18.jpg'
]

async function onLogin() {
  if (!loginForm.value.username || !loginForm.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const res = await request.post('/api/login', {
      name: loginForm.value.username,
      password: loginForm.value.password
    })
    // 存储 token/user 信息，跳转
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('user', JSON.stringify(res.data.user))
    ElMessage.success('登录成功')
    window.location.href = '/'
  } catch (e: any) {
    ElMessage.error(e?.response?.data || '登录失败')
  } finally {
    loading.value = false
  }
}

async function onRegister() {
  loading.value = true
  if (!registerForm.value.name || registerForm.value.name.length < 2) {
    ElMessage.warning('用户名不少于2位')
    loading.value = false
    return
  }
  if (!registerForm.value.password || registerForm.value.password.length < 6) {
    ElMessage.warning('密码不少于6位')
    loading.value = false
    return
  }
  if (
      !registerForm.value.email ||
      !/^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$/.test(registerForm.value.email)
  ) {
    ElMessage.warning('请输入有效邮箱')
    loading.value = false
    return
  }
  if (!registerForm.value.avatar) {
    ElMessage.warning('请选择头像')
    loading.value = false
    return
  }
  try {
    const res = await request.post('/api/register', {
      name: registerForm.value.name,
      password: registerForm.value.password,
      email: registerForm.value.email,
      avatar: registerForm.value.avatar,
      group: registerForm.value.group
    })
    ElMessage.success('注册成功，请登录')
    tab.value = 'login'
    loginForm.value.username = registerForm.value.name
    loginForm.value.password = ''
  } catch (e: any) {
    ElMessage.error(e?.response?.data || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>