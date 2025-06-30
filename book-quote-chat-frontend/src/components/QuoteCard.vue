<template>
  <el-card class="quote-card flex flex-col gap-4 " style="height:100%;">
    <div class="card-main flex-1 w-full">
      <div class="flex gap-4 items-start w-full">
        <el-avatar
            :src="quote.avatar"
            :alt="quote.user"
            class="shadow"
            @error="onImgError"
            size="large"
        />
        <div class="flex-1">
          <div class="text-base font-semibold mb-2 cursor-pointer" @click="onCopy">{{ quote.text }}</div>
          <div class="text-xs text-gray-500 mb-1">—— {{ quote.book }}</div>
        </div>
      </div>
    </div>
    <div class="card-footer flex justify-between items-center w-full">
      <el-button
          class="btn-icon-text"
          size="large"
          :type="quote.liked ? 'success' : 'default'"
          :class="{ 'like-animate': likeAnimation[quote.id], 'liked-border': quote.liked }"
          @click.stop="handleLikeToggle"
      >👍 喜欢 {{ likeCountDisplay }}</el-button>
      <el-button
          type="primary"
          :disabled="!quote.commentable"
          class="btn-icon-text comment-btn"
          size="large"
          @click.stop="onComment"
      >💬 评论</el-button>
    </div>
    <el-tag v-if="copied" type="success" class="mt-2">已复制</el-tag>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
const props = defineProps<{ quote: any; defaultAvatar?: string }>()
const emit = defineEmits(['like', 'cancelLike', 'comment'])

const likeAnimation = ref<Record<string | number, boolean>>({})
const copied = ref(false)
const localLike = ref(props.quote.liked)
const localCount = ref(props.quote.likeCount || 0)

const likeCountDisplay = computed(() => localCount.value)

function onImgError(e: Event) {
  const img = e.target as HTMLImageElement
  img.src = props.defaultAvatar || '/static/default-avatar.png'
}

// 点赞/取消本地同步 + 动画
function handleLikeToggle() {
  const id = props.quote.id
  if (!localLike.value) {
    localLike.value = true
    localCount.value += 1
    likeAnimation.value[id] = true
    emit('like', props.quote)
    setTimeout(() => { likeAnimation.value[id] = false }, 400)
  } else {
    localLike.value = false
    localCount.value = Math.max(0, localCount.value - 1)
    emit('cancelLike', props.quote)
  }
}

function onComment() {
  emit('comment', props.quote)
}

async function onCopy() {
  try {
    await navigator.clipboard.writeText(props.quote.text)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 1500)
  } catch (e) {}
}
</script>

<style scoped>
@keyframes bounce {
  0% { transform: scale(1);}
  30% { transform: scale(1.25);}
  60% { transform: scale(0.95);}
  100% { transform: scale(1);}
}
.quote-card {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 180px; /* 确保内容再少也撑开卡片 */
}
.card-main {
  flex: 1 1 auto; /* 使内容区域弹性伸缩 */
  min-height: 70px;
  width: 100%;
}
.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  border-top: 1px solid #f0f0f0;
  padding-top: 8px;
  margin-top: auto; /* 保证按钮永远贴底部 */
  min-height: 54px;
}
/* 这样可彻底消除按钮下方多余空白 */

.like-animate {
  animation: bounce 0.4s;
}
.btn-icon-text {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  padding: 8px 16px;
  border-radius: 1px;
  cursor: pointer;
  user-select: none;
  white-space: nowrap;
  box-shadow: none;
  flex: 1 1 0%;
}


</style>