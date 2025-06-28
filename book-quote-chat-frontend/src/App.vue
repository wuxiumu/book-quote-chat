<template>
  <div>
    <NavigationBar />
    <router-view class="p-4" />
  </div>
</template>

<script setup lang="ts">
import NavigationBar from '@/components/NavigationBar.vue';
</script>