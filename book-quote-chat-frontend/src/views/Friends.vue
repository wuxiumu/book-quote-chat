<template>
  <div class="max-w-4xl mx-auto pt-8 px-2">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold">👥 好友列表</h2>
      <button class="btn btn-info btn-sm" @click="showAdd = true">➕ 添加好友</button>
    </div>
    <div v-if="friends.length === 0" class="text-center opacity-60 py-12">还没有好友，快去添加一个吧！</div>
    <div class="grid gap-4 sm:grid-cols-2 md:grid-cols-3">
      <FriendCard v-for="f in friends" :key="f.id" :friend="f" @chat="chatTo" />
    </div>

    <!-- 添加好友弹窗 -->
    <div v-if="showAdd" class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-30">
      <div class="bg-white rounded-xl shadow-lg p-6 w-full max-w-xs">
        <h3 class="text-lg font-bold mb-3">添加好友</h3>
        <input v-model="newFriend" class="input input-bordered w-full mb-2" placeholder="输入好友昵称…" @keyup.enter="addFriend" />
        <button class="btn btn-primary w-full" @click="addFriend" :disabled="!newFriend.trim()">确认添加</button>
        <button class="btn btn-ghost w-full mt-2" @click="showAdd = false">取消</button>
        <div v-if="error" class="text-xs text-red-500 mt-2">{{ error }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
const router = useRouter();
import FriendCard from '@/components/FriendCard.vue';


interface Friend { id: string; name: string; avatar: string }
const friends = ref<Friend[]>([
  { id: 'f1', name: '书友A', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=A' },
  { id: 'f2', name: '书友B', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=B' },
  { id: 'f3', name: '书友C', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=C' },
]);
const showAdd = ref(false);
const newFriend = ref('');
const error = ref('');

function addFriend() {
  const name = newFriend.value.trim();
  if (!name) return;
  if (friends.value.some(f => f.name === name)) {
    error.value = '该昵称已在好友列表！';
    return;
  }
  friends.value.push({
    id: Date.now() + '',
    name,
    avatar: `https://api.dicebear.com/7.x/miniavs/svg?seed=${encodeURIComponent(name)}`
  });
  newFriend.value = '';
  error.value = '';
  showAdd.value = false;
}

function chatTo(friend: Friend) {
  router.push({ name: 'ChatRoom', params: { id: friend.id } });
}
</script>