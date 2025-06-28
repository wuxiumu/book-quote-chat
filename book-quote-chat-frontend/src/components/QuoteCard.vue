<template>
  <div
      class="card bg-base-100 shadow-xl hover:shadow-2xl hover:ring-4 hover:ring-pink-200 hover:scale-105 hover:animate-wiggle transition duration-300"
      style="min-height:140px;"
  >
    <div class="relative px-4 pt-4">
      <!-- 头像浮动，文字环绕 -->
      <img
          :src="quote.avatar"
          class="float-left mr-3 mb-2 w-12 h-12 md:w-16 md:h-16 rounded-full border-2 border-pink-200 object-cover shadow-lg hover:shadow-2xl"
          :alt="quote.user"
          loading="lazy"
          style="shape-outside:circle(50%);"
      />
      <!-- 复制按钮浮动右上角 -->
      <button
          class="btn btn-xs btn-info absolute top-0 right-0 mt-4 mr-4 z-10"
          :class="{ 'animate-bounce': showCopyAnim }"
          @click="copy"
          title="复制金句"
      >
        <span v-if="!showCopyAnim">📋</span>
        <span v-else>✅</span>
      </button>
      <!-- 内容环绕头像 -->
      <div>
        <span class="text-pink-600 font-semibold text-base md:text-lg">{{ quote.user }}</span>
        <span class="ml-2 text-xs opacity-50">{{ new Date(quote.created).toLocaleDateString() }}</span>
        <p
            class="italic leading-relaxed mb-1 select-text break-words transition-all duration-200 text-base md:text-lg"
            :class="{ 'line-clamp-2': !expanded && isLong }"
            @click="toggleExpand"
            title="点击展开/收起"
        >“{{ quote.text }}”</p>
        <button
            v-if="isLong"
            class="btn btn-ghost btn-xs px-2 py-0 mb-1 text-xs text-blue-500 hover:underline"
            @click="toggleExpand"
        >
          {{ expanded ? '收起' : '展开全文' }}
        </button>
        <div class="flex justify-end mb-1">
          <span class="badge badge-outline text-xs md:text-sm">{{ quote.book }}</span>
        </div>
      </div>
      <div class="clearfix"></div>
    </div>
    <div class="card-actions justify-between mt-2 px-4 pb-3">
      <button
          class="btn btn-xs transition transform hover:scale-110 active:scale-95"
          :class="{ 'animate-bounce': bounceLike }"
          @click="like"
      >👍 喜欢 {{ likes }}</button>
      <button class="btn btn-xs" @click="showComment = !showComment">💬 评论</button>
    </div>
    <div v-if="showComment" class="px-4 pb-3">
      <input v-model="comment" class="input input-sm input-bordered w-full" placeholder="写下评论..." />
      <button class="btn btn-sm mt-1" @click="addComment">提交</button>
      <ul class="mt-2 space-y-1">
        <li v-for="(c, i) in comments" :key="i" class="text-sm opacity-70 transition-all duration-500 opacity-0 animate-fade-in">• {{ c }}</li>
      </ul>
    </div>
  </div>
</template>

<style>
@keyframes wiggle {
  0%, 100% { transform: rotate(-2deg); }
  50% { transform: rotate(2deg); }
}
.animate-wiggle {
  animation: wiggle 0.25s;
}
@keyframes fade-in {
  from { opacity: 0; transform: translateY(10px);}
  to { opacity: 1; transform: translateY(0);}
}
.animate-fade-in {
  animation: fade-in 0.5s forwards;
}
.line-clamp-2 {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  text-overflow: ellipsis;
}
/* 新增浮动兼容清除 */
.clearfix::after {
  content: "";
  display: table;
  clear: both;
}
</style>

<script setup lang="ts">
import { ref, computed } from 'vue';
const props = defineProps<{ quote: { id: string; text: string; book: string; user: string; avatar: string; created: string } }>();

const likes = ref(Math.floor(Math.random() * 20));
const bounceLike = ref(false);
const showComment = ref(false);
const comment = ref('');
const comments = ref<string[]>([]);

const showCopyAnim = ref(false);

const expanded = ref(false);
const isLong = computed(() => props.quote.text.length > 48);

function toggleExpand() {
  if (isLong.value) expanded.value = !expanded.value;
}

function like() {
  likes.value++;
  bounceLike.value = true;
  setTimeout(() => bounceLike.value = false, 400);
}
function addComment() {
  if (comment.value.trim()) {
    comments.value.push(comment.value.trim());
    comment.value = '';
  }
}
function copy() {
  navigator.clipboard.writeText(`“${props.quote.text}” ——${props.quote.user}《${props.quote.book}》`);
  showCopyAnim.value = true;
  setTimeout(() => (showCopyAnim.value = false), 900);
}
</script>

<style>
@keyframes wiggle {
  0%, 100% { transform: rotate(-2deg); }
  50% { transform: rotate(2deg); }
}
.animate-wiggle {
  animation: wiggle 0.25s;
}
@keyframes fade-in {
  from { opacity: 0; transform: translateY(10px);}
  to { opacity: 1; transform: translateY(0);}
}
.animate-fade-in {
  animation: fade-in 0.5s forwards;
}
.line-clamp-2 {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  text-overflow: ellipsis;
}
/* 新增浮动兼容清除 */
.clearfix::after {
  content: "";
  display: table;
  clear: both;
}
</style>