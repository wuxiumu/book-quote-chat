<template>
  <div class="fixed inset-0 z-40 flex items-center justify-center bg-black/50" @click.self="$emit('close')">
    <div class="bg-white rounded-xl shadow-xl p-6 min-w-[260px] max-w-[90vw] relative">
      <button class="absolute right-3 top-3 text-xl opacity-50 hover:opacity-90" @click="$emit('close')">&times;</button>
      <div class="flex flex-col items-center gap-2">
        <img :src="user.avatar" class="w-16 h-16 rounded-full shadow mb-1" alt="用户头像" />
        <div class="font-bold text-lg">{{ user.user || user.name }}</div>
        <div class="text-sm text-gray-400 mb-2">ID: {{ user.userId || user.id }}</div>
        <div class="flex gap-2 mt-2">
          <button class="btn btn-outline btn-sm" @click="handleAddFriend(user)">加好友</button>
          <button class="btn btn-primary btn-sm" @click="handlePrivateChat(user)">私聊</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ user: any }>();
const emit = defineEmits(['close']);

function handleAddFriend(user) {
  // 这里调用后端接口即可，示例弹窗
  alert(`已向 ${user.user || user.name} 发起好友申请！`);
}
function handlePrivateChat(user) {
  // 这里路由跳转即可
  window.location.href = `/chat/private/${user.userId || user.id}`;
}
</script>