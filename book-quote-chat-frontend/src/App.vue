<template>
  <div>
    <NavigationBar />
    <div class="p-4">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import NavigationBar from '@/components/NavigationBar.vue';
</script>