<template>
  <el-dialog v-model="visible" title="发表评论" width="400px" :close-on-click-modal="true">
    <div v-if="quote" class="mb-3">
      <div class="text-xs text-gray-400 mb-2">金句：</div>
      <div class="text-sm text-gray-600 italic">“{{ quote.text }}”</div>
    </div>
    <div v-if="loading && page === 1" class="mb-3 text-center text-gray-500">评论加载中...</div>
    <div v-else-if="commentList.length" ref="commentBox" @scroll.passive="handleScroll" class="mb-3 max-h-40 overflow-auto border border-gray-200 rounded p-2">
      <div v-for="comment in commentList" :key="comment.id" class="flex items-start mb-2">
        <img :src="comment.user?.avatar || '/static/default-avatar.png'" alt="头像" class="w-6 h-6 rounded-full mr-2" />
        <div class="flex-1">
          <div class="text-xs font-semibold text-gray-700">{{ comment.user?.nickname || '匿名' }}</div>
          <div class="text-sm text-gray-800">{{ comment.content }}</div>
          <div class="text-xs text-gray-400">{{ new Date(comment.created * 1000).toLocaleString() }}</div>
        </div>
      </div>
      <div v-if="loading && page > 1" class="text-center text-gray-500">加载中...</div>
      <div v-else-if="!hasMore" class="text-center text-gray-400">没有更多了</div>
    </div>
    <div v-else class="mb-3 text-center text-gray-400">暂无评论</div>
    <el-input
        v-model="content"
        type="textarea"
        rows="4"
        placeholder="请输入评论内容"
        maxlength="200"
        show-word-limit
    />
    <template #footer>
      <el-button type="primary" @click="submit" :disabled="!content.trim()">发表评论</el-button>
    </template>
  </el-dialog>
</template>
<script setup lang="ts">
import { ref, watch } from 'vue'
import request from '@/api/request'
const props = defineProps(['modelValue', 'quote'])
const emit = defineEmits(['update:modelValue', 'submit'])
const visible = ref(false)
const content = ref('')
const loading = ref(false)
const page = ref(1)
const hasMore = ref(true)
const PAGE_SIZE = 10
const commentBox = ref()

interface Comment {
  id: number;
  user?: {
    avatar?: string;
    nickname?: string;
  };
  content: string;
  created: number;
}

const commentList = ref<Comment[]>([])

async function loadComments(id: number, append = false) {
  if (!hasMore.value && append) return
  loading.value = true
  try {
    const res = await request.get('/api/comments', {
      params: {
        targetType: 'quote',
        targetId: id,
        offset: (page.value - 1) * PAGE_SIZE,
        limit: PAGE_SIZE
      }
    })
    const data = res.data.list || []
    if (append) {
      commentList.value = commentList.value.concat(data)
    } else {
      commentList.value = data
    }
    hasMore.value = data.length === PAGE_SIZE
  } catch {
    if (!append) commentList.value = []
    hasMore.value = false
  } finally {
    loading.value = false
  }
}

function handleScroll() {
  const el = commentBox.value
  if (!el || loading.value || !hasMore.value) return
  if (el.scrollTop + el.clientHeight >= el.scrollHeight - 10) {
    page.value++
    if (props.quote && props.quote.id) {
      loadComments(props.quote.id, true)
    }
  }
}

watch(() => props.modelValue, v => visible.value = v)
watch(visible, v => emit('update:modelValue', v))
watch(() => props.quote, (newQuote) => {
  content.value = ''
  page.value = 1
  hasMore.value = true
  if (newQuote && newQuote.id) {
    loadComments(newQuote.id)
  } else {
    commentList.value = []
  }
})
async function submit() {
  if (!props.quote || !content.value.trim()) return
  try {
    await request.post('/api/comments', {
      targetType: 'quote',
      targetId: props.quote.id,
      content: content.value.trim()
    })
    emit('submit', { quoteId: props.quote.id, content: content.value.trim() })
    visible.value = false
    page.value = 1
    hasMore.value = true
    await loadComments(props.quote.id)
  } catch (e) {
    console.error('评论失败', e)
  }
}
</script>