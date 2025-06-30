<template>
  <div :class="self ? 'chat chat-end' : 'chat chat-start'">
    <div class="chat-image avatar">
      <template v-if="msg.avatar">
        <img
            :src="msg.avatar"
            :title="`查看${msg.user}的资料`"
            @click="$emit('showProfile', { ...msg, self })"
            class="w-10 h-10 rounded-full shadow object-cover cursor-pointer hover:scale-110 hover:ring-2 hover:ring-blue-400 transition duration-150"
        />
      </template>
      <template v-else>
        <div class="w-10 h-10 rounded-full bg-neutral-focus text-neutral-content flex items-center justify-center shadow">
          {{ ( msg.user || '系统').charAt(0).toUpperCase() }}
        </div>
      </template>
    </div>
    <div class="chat-header flex items-center gap-2">
      <span class="font-semibold">{{ msg.user }}</span>
      <span class="text-xs opacity-50">{{ new Date(msg.created).toLocaleTimeString() }}</span>
    </div>
    <div class="chat-bubble animate-fade-in markdown-body" v-html="renderedText" @click="onBubbleClick"></div>
  </div>
  <!-- 全屏图片/视频预览蒙层 -->
  <div
      v-if="previewUrl"
      class="fixed inset-0 z-50 bg-black/80 flex items-center justify-center"
      @click="closePreview"
      style="cursor: zoom-out;"
  >
    <template v-if="previewType === 'image'">
      <img :src="previewUrl" class="max-w-full max-h-full rounded-lg shadow-2xl transition-all duration-200" />
    </template>
    <template v-else-if="previewType === 'video'">
      <video
          :src="previewUrl"
          class="max-w-full max-h-full rounded-lg shadow-2xl"
          controls
          autoplay
          style="background:#222;"
          @click.stop
      />
    </template>
  </div>
</template>
<script setup lang="ts">
import { defineProps, computed,onMounted, ref, nextTick } from 'vue';
import { marked } from 'marked';
const props = defineProps<{ msg: { user: string; avatar?: string; text: string; created: string; image?: string }, self?: boolean }>()
const renderedText = computed(() => marked.parse(props.msg.text || ''))

const previewUrl = ref('');
const previewType = ref(''); // 'image' 或 'video'

onMounted(() => {
  nextTick(() => {
    document.querySelectorAll('.markdown-body img').forEach(img => {
      const image = img as HTMLImageElement
      image.style.cursor = 'pointer';
      image.onclick = () => openPreview('image', image.src);
    });
    document.querySelectorAll('.markdown-body video').forEach(video => {
      const v = video as HTMLVideoElement
      v.style.cursor = 'pointer';
      v.onclick = () => openPreview('video', v.src || (v.querySelector('source')?.src ?? ''));
    });
  });
});

// 打开全屏预览
function openPreview(type: string, url: string) {
  previewType.value = type;
  previewUrl.value = url;
}
// 关闭
function closePreview() {
  previewUrl.value = '';
  previewType.value = '';
}

function onBubbleClick(e: MouseEvent) {
  // 只处理 video 标签点击
  if ((e.target as HTMLElement).tagName === 'VIDEO') {
    const video = e.target as HTMLVideoElement;
    // 切换播放/暂停
    if (video.paused) {
      video.play();
      video.muted = true;
    } else {
      video.muted = !video.muted;
      if (video.paused) video.play();
    }
    e.stopPropagation();
  }
}
</script>
<style>
@keyframes fade-in {
  from { opacity: 0; transform: translateY(10px);}
  to { opacity: 1; transform: translateY(0);}
}
.animate-fade-in {
  animation: fade-in 0.5s forwards;
}
.markdown-body video {
  max-width: 80vw;
  max-height: 80vh;
  width: auto;
  height: auto;
  display: block;
  border-radius: 10px;
  margin: 8px 0;
  background: #111;
}
.fixed[inset-0] img,
.fixed[inset-0] video {
  animation: fade-in 0.3s;
}
</style>