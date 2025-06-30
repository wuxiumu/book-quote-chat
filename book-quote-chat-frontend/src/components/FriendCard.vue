<template>
  <div class="card bg-base-100 shadow-md hover:shadow-xl hover:ring-2 hover:ring-pink-200 transition duration-200 flex flex-col items-center py-4 px-3 group relative">
    <div class="relative">
      <img
          :src="friend.avatar"
          class="w-16 h-16 rounded-full border-2 border-pink-200 shadow mb-2 object-cover transition-transform duration-500 group-hover:scale-110 group-hover:rotate-3"
          :alt="friend.name"
      />
      <span
        class="absolute bottom-0 right-0 w-4 h-4 rounded-full border-2 border-white"
        :class="friend.online ? 'bg-green-500' : 'bg-gray-400'"
        title="Online status"
      ></span>
    </div>
    <div class="font-semibold text-pink-600 text-base truncate">{{ friend.name }}</div>
    <button class="btn btn-xs mt-2 btn-outline btn-primary" @click="$emit('chat', friend)">
      私聊
    </button>
    <div class="mt-3 w-full px-2">
      <label class="block text-xs text-gray-500 mb-1" for="remark-input">备注</label>
      <input
        id="remark-input"
        v-model="editableRemark"
        @keydown.enter.prevent="saveRemark"
        @blur="saveRemark"
        type="text"
        placeholder="点击编辑备注"
        class="input input-sm input-bordered w-full"
      />
    </div>
    <div class="mt-2 w-full px-2">
      <label class="block text-xs text-gray-500 mb-1">分组</label>
      <div class="text-sm text-gray-700 truncate" :title="friend.group">{{ friend.group || '未分组' }}</div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, watch } from 'vue';

interface Friend {
  id: string;
  name: string;
  avatar: string;
  group: string;
  remark?: string; // 可选
  online: boolean;
}
const props = defineProps<{
  friend: Friend
  onChat?: (friend: Friend) => any
  "onUpdate-remark"?: (payload: { id: string; remark: string }) => any
}>()

const emit = defineEmits<{
  (e: 'chat', friend: Friend): void;
  (e: 'update-remark', payload: { id: string; remark: string }): void;
}>();

const editableRemark = ref(props.friend.remark ?? '');


watch(() => props.friend.remark, (newRemark) => {
  editableRemark.value = newRemark ?? '';
});

function saveRemark() {
  const trimmed = editableRemark.value.trim();
  if (trimmed !== (props.friend.remark ?? '')) {
    emit('update-remark', { id: props.friend.id, remark: trimmed });
  }
  editableRemark.value = trimmed;
}
</script>