<template>
  <div class="flex flex-col min-h-[80vh]">
    <!-- 顶部欢迎 -->
    <div class="text-center mt-10 flex-1">
      <h1 class="text-4xl font-extrabold mb-4">
        <template v-if="user">
          欢迎回来，{{ user.name }}！✨
        </template>
        <template v-else>
          欢迎来到 QuoteChat ✨
        </template>
      </h1>
      <p class="text-lg opacity-80">分享书中金句，结识志同道合的朋友。</p>

      <!-- 随机金句 -->
      <transition name="fade-move" mode="out-in">
        <div v-if="randomQuote" key="random-quote" class="my-10 mx-auto max-w-xl flex flex-col items-center">
          <div class="text-xl italic font-semibold text-pink-700 animate-bounce-in">
            {{ randomQuote.text }}
          </div>
          <div class="mt-1 text-sm text-gray-500">—— {{ randomQuote.book || '佚名' }}</div>
        </div>
        <div v-else key="loading" class="my-10 text-gray-400 text-base">正在加载金句...</div>
      </transition>
      <button class="btn btn-sm btn-outline mb-2" @click="nextQuote">换一句</button>

      <!-- 掌声互动 -->
      <div class="my-4">
        <button class="btn btn-accent btn-lg flex items-center gap-2 relative" @click="clap">
          👏
          <span>{{ clapCount }}</span>
          <span v-if="showEffect" class="absolute -top-5 left-1/2 -translate-x-1/2 text-pink-500 text-xl animate-pop">+1</span>
        </button>
      </div>

      <div class="mt-8 join join-vertical lg:join-horizontal">
        <router-link to="/quotes" class="btn btn-primary join-item">开始浏览金句</router-link>
        <router-link to="/chat" class="btn btn-secondary join-item">加入群聊</router-link>
      </div>
    </div>

    <!-- 底部联系方式和友情链接 -->
    <footer class="mt-10 bg-base-100 text-center p-4 border-t">
      <div class="mb-2">© 2025 读书金句群聊 | 工程师：Aric | 邮箱：<a href="mailto:wuxiumu@163.com" class="link">wuxiumu@163.com</a></div>
      <div class="flex flex-wrap justify-center gap-2 items-center text-sm">
        <span>友情链接：</span>
        <template v-for="link in links" :key="link.url">
          <a :href="link.url" class="link link-info mx-1" target="_blank" rel="noopener">{{ link.name }}</a>
        </template>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import request from '@/api/request'
const user = ref<{name: string} | null>(null);
onMounted(() => {
  try {
    const saved = localStorage.getItem('user');
    if (saved) user.value = JSON.parse(saved);
  } catch {}
});

// 新版——请求后端接口获取10条金句

const randomQuotes = ref<{text:string,book:string}[]>([]);
const randomQuote = ref<{text:string,book:string}|null>(null);

async function fetchQuotes() {
  request.get('/api/quotes/random?limit=10').then(res => {
    randomQuotes.value =  res.data;
    if (randomQuotes.value.length) {
        randomQuote.value = randomQuotes.value[Math.floor(Math.random() * randomQuotes.value.length)];
    }
  })
}
function nextQuote() {
  if (!randomQuotes.value.length) return;
  let idx;
  do {
    idx = Math.floor(Math.random() * randomQuotes.value.length);
  } while (randomQuotes.value[idx] === randomQuote.value && randomQuotes.value.length > 1);
  randomQuote.value = randomQuotes.value[idx];
}



// 掌声互动
const clapCount = ref(0);
const showEffect = ref(false);
const clapped = ref(false);

async function fetchClapCount() {
  try {
    const res = await request.get('/api/clap/count');
    clapCount.value = res.data.count || 0;
  } catch {
    clapCount.value = 0;
  }
}
onMounted(() => {
  fetchQuotes();
  fetchLinks();
  fetchClapCount();
});

async function clap() {
  if (clapped.value) {
    clapCount.value += 1;
    return;
  }
  try {
    const res = await request.post('/api/clap');
    if (res.data && res.data.success) {
      clapCount.value = res.data.count;
      clapped.value = true;
      showEffect.value = true;
      setTimeout(() => (showEffect.value = false), 700);
    } else {
      clapCount.value += 1;
    }
  } catch (e) {
    clapCount.value += 1;
  }
}

// 友情链接
const links = ref<{name: string; url: string}[]>([]);
onMounted(() => {
  fetchQuotes();
  fetchLinks();
});

async function fetchLinks() {
  request.get('/api/friend-links').then(res => {
    links.value =  res.data;
  })
}
</script>

<style scoped>
.fade-move-enter-active, .fade-move-leave-active {
  transition: opacity 0.6s, transform 0.6s;
}
.fade-move-enter-from, .fade-move-leave-to {
  opacity: 0;
  transform: translateY(30px);
}
@keyframes pop {
  0% { opacity: 0; transform: scale(0.8);}
  50% { opacity: 1; transform: scale(1.2);}
  100% { opacity: 0; transform: scale(1);}
}
.animate-pop {
  animation: pop 0.7s;
}
@keyframes bounce-in {
  0% { opacity: 0; transform: scale(0.8);}
  80% { opacity: 1; transform: scale(1.05);}
  100% { opacity: 1; transform: scale(1);}
}
.animate-bounce-in {
  animation: bounce-in 0.8s;
}
</style>