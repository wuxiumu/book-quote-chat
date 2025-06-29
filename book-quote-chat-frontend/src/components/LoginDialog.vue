<template>
  <el-dialog v-model="visible" title="登录" width="300px" :close-on-click-modal="false">
    <el-form :model="form">
      <el-form-item>
        <el-input v-model="form.username" placeholder="用户名"/>
      </el-form-item>
      <el-form-item>
        <el-input v-model="form.password" placeholder="密码" type="password"/>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button type="primary" @click="onLogin">登录</el-button>
    </template>
  </el-dialog>
</template>
<script setup lang="ts">
import { ref, watch } from 'vue'
import request from '@/api/request'
const props = defineProps(['modelValue'])
const emit = defineEmits(['update:modelValue', 'success'])
const visible = ref(false)
const form = ref({ username: '', password: '' })
watch(() => props.modelValue, v => visible.value = v)
watch(visible, v => emit('update:modelValue', v))
async function onLogin() {
  try {
    const res = await request.post('/api/login', {
      mode: 'password',
      username: form.value.username,
      password: form.value.password
    })
    if (res.data && res.data.token) {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      emit('success', res.data.user)
      visible.value = false
    } else {
      alert('登录失败')
    }
  } catch {
    alert('登录失败')
  }
}
</script>