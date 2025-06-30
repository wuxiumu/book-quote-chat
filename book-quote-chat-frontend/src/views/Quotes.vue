<template>
  <div>
    <h2 class="text-2xl font-bold mb-4">📚 书摘金句</h2>
    <div
        class="grid gap-4 md:grid-cols-2 lg:grid-cols-3"
        ref="quoteGrid"
        @scroll.passive="onScroll"
        style="max-height: 80vh; overflow-y: auto;"
    >
      <QuoteCard
          v-for="q in quotes"
          :key="q.id"
          :quote="q"
          :user="user"
          :default-avatar="defaultAvatar"
          @like="handleLike"
          @cancelLike="handleCancelLike"
          @comment="handleComment"
          :like-animating="likeAnimation[q.id]"
      />
    </div>
    <div v-if="loading" class="my-6 text-center text-gray-400">加载中...</div>
    <div v-if="noMore" class="my-6 text-center text-gray-400">没有更多啦~</div>
    <LoginDialog v-model="showLogin" @success="onLoginSuccess"/>
    <CommentDialog v-model="showComment" :quote="commentTarget" @submit="onCommentSubmit" :visible.sync="showComment"/>
  </div>
</template>

<script setup lang="ts">
import request from '@/api/request'
import { ref, onMounted, nextTick } from 'vue';
import QuoteCard from '@/components/QuoteCard.vue';
import LoginDialog from '@/components/LoginDialog.vue';
import CommentDialog from '@/components/CommentDialog.vue';

interface Quote {
  id: string;
  text: string;
  book: string;
  user: string;
  avatar: string;
  created: string;
  liked?: boolean;
  commentable?: boolean;
  likeAnimating?: boolean;
  likeCount?: number;    // ← 补充这一行！
}

const defaultAvatar = '/static/default-avatar.png'; // 你项目的默认头像路径
const quotes = ref<Quote[]>([]);
const offset = ref(0);
const limit = 30;
const loading = ref(false);
const noMore = ref(false);
const user = ref<{ id: string, name: string } | null>(null);

const showLogin = ref(false)
const showComment = ref(false)
const commentTarget = ref<Quote | null>(null)

const likedIds = ref<Set<string>>(new Set()) // 这样无论分页加载多少条，用户已点过的金句都能准确标记“已点赞”。

// 点赞动画状态，点赞后按钮带动画
const likeAnimation = ref<Record<string, boolean>>({})

const quoteGrid = ref<HTMLElement | null>(null)

async function fetchLikedMap() {
  try {
    const res = await request.get('/api/likes', { params: { targetType: 'quote' } });
    if (res.status === 200 && res.data && typeof res.data === 'object') {
      likedIds.value = new Set(Object.keys(res.data));
    } else {
      likedIds.value = new Set();
    }
  } catch {
    likedIds.value = new Set();
  }
}

// 这样每行的卡片高度始终一致，按钮对齐底部。
function quoteCardMainHeightSync() {
  if (!quoteGrid.value) return;
  // 根据响应式设计，lg:grid-cols-3 最大列数为3
  // md:grid-cols-2 中等屏幕为2列，小屏幕为1列
  // 这里取最大列数3
  const perRow = 3;
  const cardMains = quoteGrid.value.querySelectorAll<HTMLElement>('.card-main');
  for (let i = 0; i < cardMains.length; i += perRow) {
    let maxHeight = 0;
    const group = Array.from(cardMains).slice(i, i + perRow);
    group.forEach(el => {
      el.style.height = 'auto'; // 先重置高度
      if (el.offsetHeight > maxHeight) maxHeight = el.offsetHeight;
    });
    group.forEach(el => {
      el.style.height = maxHeight + 'px';
    });
  }
}

onMounted(async () => {
  try {
    const u = localStorage.getItem('user')
    if (u) user.value = JSON.parse(u)
  } catch {}
  if (user.value) {
    await fetchLikedMap();
  }
  await loadMore();
  await nextTick();
  quoteCardMainHeightSync();
  window.addEventListener('resize', quoteCardMainHeightSync);
});

async function loadMore() {
  if (loading.value || noMore.value) return;
  loading.value = true;
  try {
    const res = await request.get('/api/quotes', { params: { limit, offset: offset.value } })
    let list = res.data.list || res.data || []
    list = list.map((q: any) => ({
      ...q,
      avatar: q.avatar || defaultAvatar,
      liked: likedIds.value.has(q.id),
      commentable: q.commentable !== false,
      likeAnimating: likeAnimation.value[q.id] // 赋值点赞动画状态
    }))
    if (list.length < limit) noMore.value = true
    quotes.value.push(...list)
    offset.value += list.length
  } catch (e) {
    noMore.value = true
  }
  loading.value = false
  await nextTick();
  quoteCardMainHeightSync();
}

function onScroll(e: Event) {
  const el = e.target as HTMLElement
  if (el.scrollTop + el.clientHeight >= el.scrollHeight - 100) {
    loadMore();
  }
}

// 点赞
async function handleLike(quote: Quote) {
  if (!user.value) {
    showLogin.value = true
    return
  }
  if (quote.liked) return // 已点赞不重复发请求
  // 本地乐观 +1
  quote.liked = true
  likedIds.value.add(quote.id)
  quote.likeCount = (quote.likeCount || 0) + 1
  likeAnimation.value[quote.id] = true
  setTimeout(() => {
    likeAnimation.value[quote.id] = false
  }, 500)
  try {
    const res = await request.post('/api/likes', { targetType: 'quote', targetId: quote.id })
    if (res.status !== 200) {
      // 后端失败则回滚
      quote.liked = false
      likedIds.value.delete(quote.id)
      quote.likeCount = Math.max(0, (quote.likeCount || 1) - 1)
      console.error(res.data.message || '点赞失败')
    }
  } catch {
    quote.liked = false
    likedIds.value.delete(quote.id)
    quote.likeCount = Math.max(0, (quote.likeCount || 1) - 1)
    console.error('点赞失败')
  }
}

// 取消点赞
async function handleCancelLike(quote: Quote) {
  if (!user.value) {
    showLogin.value = true
    return
  }
  if (!quote.liked) return // 未点赞时不发请求
  // 本地乐观 -1
  quote.liked = false
  likedIds.value.delete(quote.id)
  quote.likeCount = Math.max(0, (quote.likeCount || 1) - 1)
  try {
    const res = await request.delete('/api/likes', { data: { targetType: 'quote', targetId: quote.id } })
    if (res.status !== 204) {
      // 后端失败则回滚
      quote.liked = true
      likedIds.value.add(quote.id)
      quote.likeCount = (quote.likeCount || 0) + 1
      console.error(res.data.message || '取消点赞失败')
    }
  } catch {
    quote.liked = true
    likedIds.value.add(quote.id)
    quote.likeCount = (quote.likeCount || 0) + 1
    console.error('取消点赞失败')
  }
}

// 评论
function handleComment(quote: Quote) {
  if (!user.value) {
    showLogin.value = true;
    return;
  }
  if (!quote.commentable) {
    console.error('没有评论权限');
    return;
  }
  commentTarget.value = quote
  showComment.value = true
}
function onLoginSuccess(userInfo: any) {
  user.value = userInfo
  showLogin.value = false
}
function onCommentSubmit(_: any) {
  // 可以添加评论成功后的刷新逻辑
  showComment.value = false
}
</script>