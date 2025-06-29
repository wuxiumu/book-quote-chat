<template>
  <div class="fixed inset-0 z-40 flex items-center justify-center bg-black/50" @click.self="$emit('close')">
    <div class="bg-white rounded-xl shadow-xl p-6 min-w-[260px] max-w-[90vw] relative">
      <button class="absolute right-3 top-3 text-xl opacity-50 hover:opacity-90" @click="$emit('close')">&times;</button>
      <div class="flex flex-col items-center gap-2">
        <img :src="user.avatar" class="w-16 h-16 rounded-full shadow mb-1" alt="用户头像" />
        <div class="font-bold text-lg">{{ user.user || user.name }}</div>
        <div class="text-sm text-gray-400 mb-2">
          <template v-if="myUser && user.userid === myUser.id">
            你的ID: {{ user.userid || user.id }}
            <button @click="copyId(user.userid || user.id)" class="ml-2 text-xs btn btn-link">复制</button>
          </template>
          <template v-else>
            Ta的ID: {{ user.userid || user.id }}
            <button @click="copyId(user.userid || user.id)" class="ml-2 text-xs btn btn-link">复制</button>
          </template>
        </div>
        <div class="flex gap-2 mt-2">
          <button class="btn btn-outline btn-sm" @click="handleAddFriend(user)" :disabled="myUser && user.userid === myUser.id">加好友</button>
          <button class="btn btn-primary btn-sm" @click="handlePrivateChat(user)">私聊</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { message } from '@/utils/message'

const props = defineProps<{ user: any; myUser?: any }>();
const emit = defineEmits(['close']);
const myUser = JSON.parse(localStorage.getItem('user') || '{}')

function copyId(id) {
  navigator.clipboard.writeText(id)
    .then(() => {
      message.success('已复制ID')
    });
}
function handleAddFriend(user) {
  console.log("myUser",myUser)
  console.log("user",user)
  if (myUser && user.userid === myUser.id) {
    message.info('不能加自己为好友')
    return;
  }
  message.success(`已向 ${user.user || user.name} 发起好友申请！`)
}
function handlePrivateChat(user) {
  if (myUser && user.userid === myUser.id) {
    message.info('不能和自己私聊')
    return;
  }
  window.location.href = `/chat/private/`+genSessionKey(user.userid, myUser.id);
}
function genSessionKey(id1, id2) {
  // 保证 key 顺序唯一，小的在前
  return [id1, id2].sort().join('_')
}
</script>