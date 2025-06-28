<template>
  <div class="max-w-4xl mx-auto pt-8 px-2">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold">👥 好友列表</h2>
      <button class="btn btn-info btn-sm" @click="showAdd = true">➕ 添加好友</button>
    </div>
    <div class="mb-4">
      <label for="groupFilter" class="mr-2 font-medium">分组筛选：</label>
      <select id="groupFilter" v-model="selectedGroup" class="select select-bordered select-sm w-40">
        <option value="">全部分组</option>
        <option v-for="group in groups" :key="group" :value="group">{{ group }}</option>
      </select>
    </div>
    <div v-if="filteredFriends.length === 0" class="text-center opacity-60 py-12">还没有好友，快去添加一个吧！</div>
    <div class="grid gap-6 sm:grid-cols-2 md:grid-cols-3">
      <FriendCard
          v-for="f in filteredFriends"
          :key="f.id"
          :friend="f"
          @chat="chatTo"
          @update-remark="updateRemark"
      />
    </div>

    <!-- 添加好友弹窗 -->
    <div v-if="showAdd" class="fixed inset-0 bg-black bg-opacity-20 flex items-center justify-center z-30">
      <div class="bg-white rounded-xl shadow-lg p-6 w-full max-w-xs">
        <h3 class="text-lg font-bold mb-3">添加好友</h3>
        <input ref="addInput" v-model="newName" class="input input-bordered w-full mb-2" placeholder="输入好友昵称…" @keyup.enter="addFriend" @keyup.esc="showAdd = false" />
        <button class="btn btn-primary w-full" @click="addFriend" :disabled="!newName.trim()">确认添加</button>
        <button class="btn btn-ghost w-full mt-2" @click="showAdd = false">取消</button>
        <div v-if="error" class="text-xs text-red-500 mt-2">{{ error }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch, computed } from 'vue';
import FriendCard from '@/components/FriendCard.vue';

interface Friend { id: string; name: string; avatar: string; group: string; remark?: string; online: boolean }
const friends = ref<Friend[]>([]);

const showAdd = ref(false);
const newName = ref('');
const error = ref('');
const addInput = ref<HTMLInputElement | null>(null);

const selectedGroup = ref('');

onMounted(async () => {
  // 本地 mock 数据，可替换成 fetch
  friends.value = [
    { id: '1', name: '阿狸', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=ali', group: '同事', remark: '前端达人', online: true },
    { id: '2', name: '小王', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=wang', group: '大学同学', remark: '一块读过《活着》', online: false },
    { id: '3', name: 'Lina', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=lina', group: '群友', remark: '诗词爱好者', online: true },
  ];
});

const groups = computed(() => {
  const set = new Set<string>();
  friends.value.forEach(f => {
    if (f.group) set.add(f.group);
  });
  return Array.from(set);
});

const filteredFriends = computed(() => {
  if (!selectedGroup.value) return friends.value;
  return friends.value.filter(f => f.group === selectedGroup.value);
});

function addFriend() {
  const nameTrimmed = newName.value.trim();
  if (!nameTrimmed) return;
  if (friends.value.some(f => f.name === nameTrimmed)) {
    error.value = '该昵称已在好友列表！';
    return;
  }
  friends.value.push({
    id: Date.now() + '',
    name: nameTrimmed,
    avatar: `https://api.dicebear.com/7.x/miniavs/svg?seed=${encodeURIComponent(nameTrimmed)}`,
    group: '群友',
    remark: '',
    online: false,
  });
  newName.value = '';
  showAdd.value = false;
  error.value = '';
}

function chatTo(friend: Friend) {
  alert(`进入和【${friend.name}】的私聊页面（你可以跳转到私聊路由）`);
}

function updateRemark(payload: { id: string; remark: string }) {
  const friend = friends.value.find(f => f.id === payload.id);
  if (friend) {
    friend.remark = payload.remark;
  }
}

watch(showAdd, (val) => {
  if (val) {
    nextTick(() => {
      addInput.value?.focus();
      error.value = '';
    });
  }
});
</script>