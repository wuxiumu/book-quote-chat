<template>
  <div class="flex flex-col h-[100dvh] bg-base-200">
    <!-- 群聊顶部 -->
    <div class="sticky top-0 z-10 bg-base-100/80 backdrop-blur border-b px-4 py-2 flex items-center gap-4">
      <span class="font-bold text-pink-600 text-lg">📚 书友群聊</span>
      <span class="text-xs text-gray-400">共 {{ onlineCount }} 人在线</span>
      <span class="ml-auto text-xs text-pink-400 truncate">公告：文明发言，推荐你最近喜欢的书！</span>
    </div>

    <!-- 聊天消息区 -->
    <div class="flex-1 overflow-y-auto p-4 space-y-3" ref="scrollArea">
      <ChatMessage
          v-for="m in messages"
          :key="m.id"
          :msg="m"
          :self="m.user === myName"
      />
    </div>

    <!-- 输入区 -->
    <form
        class="flex gap-2 p-3 bg-base-100 border-t sticky bottom-0 z-10"
        @submit.prevent="send"
        @keydown.enter.exact.prevent="send"
    >
      <input
          v-model="input"
          class="input input-bordered flex-1"
          :placeholder="myName ? '说点什么...' : '请输入昵称后发言'"
          maxlength="100"
          autocomplete="off"
          @keydown.ctrl.enter.stop
      />
      <button class="btn btn-primary" type="submit" :disabled="!input.trim()">发送</button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, computed } from 'vue';
import ChatMessage from '@/components/ChatMessage.vue';

interface Msg {
  id: string;
  user: string;
  avatar: string;
  text: string;
  created: string;
}
const chatTexts = [
  '早上好，大家读什么书？', '今天阳光不错，适合看书。',
  '有人看过《活着》吗？', '推荐一本最近喜欢的书吧。',
  '这句金句很棒：“人生没有白走的路，每一步都算数。”',
  '谁最近在追《遥远的救世主》？',
  '有没有喜欢科幻的朋友？',
  '你们都用什么读书笔记工具？',
  '现在有什么书友打卡群吗？',
  '明天你们打算看什么书？',
  '听说豆瓣新热评很精彩！',
  '你们觉得读电子书还是纸质书更好？',
  '有没有喜欢诗词的？',
  '今天完成了10页，成就感爆棚！',
  '新朋友快来介绍下自己吧~',
  '我在复读《小王子》，还是感动。',
  '有书友喜欢写作吗？',
  '今晚有群聊活动吗？',
  '一起加油，早睡早起看书！',
  '请大家文明发言，群聊更和谐。',
];

const chatUsers = [
  { name: '小张', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=zhang' },
  { name: '小王', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=wang' },
  { name: '阿狸', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=ali' },
  { name: '小红', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=xiaohong' },
  { name: '小刚', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=gang' },
  { name: '小李', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=li' },
  { name: '小南', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=nan' },
  { name: '书友A', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=a' },
  { name: '书友B', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=b' },
];

const myName = 'Aric';
const myAvatar = 'https://api.dicebear.com/7.x/miniavs/svg?seed=amu';
function randomMsg(): Msg {
  const user = chatUsers[Math.floor(Math.random() * chatUsers.length)];
  const text = chatTexts[Math.floor(Math.random() * chatTexts.length)];
  const time = new Date(Date.now() - Math.random() * 3 * 86400000).toISOString();
  return {
    id: (Date.now() + Math.random()).toString(),
    user: user.name,
    avatar: user.avatar,
    text,
    created: time,
  };
}
const messages = ref<Msg[]>(Array.from({ length: 50 }, () => randomMsg()));

const input = ref('');
const scrollArea = ref<HTMLElement | null>(null);

// 假设有 23 人在线，后期可以从后端获取
const onlineCount = computed(() => {
  const set = new Set(messages.value.map(m => m.user));
  return set.size + 7; // 假定还没说话的也有
});

function send() {
  if (!input.value.trim()) return;
  messages.value.push({
    id: (Date.now() + Math.random()).toString(),
    user: myName,
    avatar: myAvatar,
    text: input.value,
    created: new Date().toISOString()
  });
  input.value = '';
  nextTick(() => scrollArea.value?.scrollTo(0, scrollArea.value.scrollHeight));
}

// 自动滚动到底
onMounted(() => {
  setInterval(() => {
    const user = chatUsers.filter(u => ['小李', '小红', '小刚'].includes(u.name))[Math.floor(Math.random() * 3)];
    const text = ['早', '分享个金句', '今天状态真好'][Math.floor(Math.random() * 3)];
    messages.value.push({
      id: (Date.now() + Math.random()).toString(),
      user: user.name,
      avatar: user.avatar,
      text,
      created: new Date().toISOString()
    });
    nextTick(() => scrollArea.value?.scrollTo(0, scrollArea.value.scrollHeight));
  }, 5000);
});
</script>