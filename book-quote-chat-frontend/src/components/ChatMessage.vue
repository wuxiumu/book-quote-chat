<template>
  <div :class="self ? 'chat chat-end' : 'chat chat-start'">
    <div class="chat-image avatar">
      <div class="w-10 rounded-full bg-neutral-focus text-neutral-content flex items-center justify-center shadow">
        {{ msg.user.charAt(0).toUpperCase() }}
      </div>
    </div>
    <div class="chat-header flex items-center gap-2">
      <span class="font-semibold">{{ msg.user }}</span>
      <span class="text-xs opacity-50">{{ new Date(msg.created).toLocaleTimeString() }}</span>
    </div>
    <div class="chat-bubble animate-fade-in">{{ msg.text }}</div>
  </div>
</template>
<script setup lang="ts">
import { defineProps } from 'vue';
defineProps<{ msg: { user: string; text: string; created: string }, self?: boolean }>()
</script>
<style>
@keyframes fade-in {
  from { opacity: 0; transform: translateY(10px);}
  to { opacity: 1; transform: translateY(0);}
}
.animate-fade-in {
  animation: fade-in 0.5s forwards;
}
</style>