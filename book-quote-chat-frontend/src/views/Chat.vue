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
      <transition-group name="fade-slide" tag="div">
      <ChatMessage
          v-for="m in messages"
          :key="m.id"
          :msg="m"
          :self="m.userid === myId"
          @showProfile="openProfile"
      />
      </transition-group>
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
      <button type="button" class="btn btn-square btn-ghost" @click="showEmoji = !showEmoji" title="表情">
        😊
      </button>
      <input type="file" class="hidden" id="fileInput" :accept="acceptStr" @change="handleFileChange" />
      <label for="fileInput" class="btn btn-square btn-ghost" title="上传资源">📷</label>
      <button class="btn btn-primary" type="submit" :disabled="!input.trim()">发送</button>
    </form>

  </div>

  <!-- 名片弹窗 -->
  <ProfileDialog
    v-if="showProfile"
    :user="profileUser"
    @close="showProfile = false"
  />
  <div
      v-show="showEmoji"
      ref="emojiMenuRef"
      class="fixed left-1/2 bottom-28 z-50 -translate-x-1/2 rounded-2xl bg-white/95 border shadow-lg flex flex-col"
      style="width: 100%; height: 410px;"
  >
    <EmojiPicker
        @select="addEmoji"
        :disable-group="false"
        :disable-skin="false"
        style="width:100%;height:370px;overflow-y:auto;"
    />
    <button class="w-full btn btn-outline rounded-b-2xl" @click="showEmoji=false" style="height:40px;">
      关闭表情菜单
    </button>
  </div>
  <div v-if="loginDialog" class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center">
    <div class="bg-white p-6 rounded-xl shadow-xl flex flex-col items-center">
      <span class="mb-4 text-pink-600 text-lg font-bold">请先登录</span>
      <button class="btn btn-primary w-full" @click="loginDialog=false;$router.push('/user/login')">去登录</button>
      <button class="mt-2 btn btn-ghost" @click="loginDialog=false">取消</button>
    </div>
  </div>
</template>

<script setup lang="ts">
// 读取上传类型配置
const uploadAccept = import.meta.env.VITE_UPLOAD_ACCEPT || 'jpg,jpeg,png,gif,webp,bmp,mp4,webm,ogg';
const acceptStr = uploadAccept
  .split(',')
  .map(ext => '.' + ext.trim())
  .join(',');
// 获取图片完整 URL
function getImageUrl(url: string) {
  if (!url) return '';
  if (/^https?:\/\//.test(url)) return url;
  const base = import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, '') || '';
  return base + (url.startsWith('/') ? url : '/' + url);
}
import { ref, onMounted, nextTick, computed,watch } from 'vue';
import request from '@/api/request'
import ChatMessage from '@/components/ChatMessage.vue';
import ProfileDialog from '@/components/ProfileDialog.vue';
import EmojiPicker from 'vue3-emoji-picker'
const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';
const token = localStorage.getItem('token')
const wsUrl = `${WS_URL}?token=${encodeURIComponent(token||'')}`;
const loginDialog = ref(false)

interface Msg {
  id: string;
  user: string;
  userid: string;
  type?: string;
  avatar: string;
  text: string;
  created: string;
  image?: string;
}

const user = ref<{ name: string; avatar: string } | null>(null)
const myName = computed(() => user.value?.name || '游客');
const myId = computed(() => user.value?.userid || '');
function loadUser() {
  try {
    const u = JSON.parse(localStorage.getItem('user') || 'null')
    if (u && u.name && u.avatar) {
      user.value = { name: u.name, avatar: u.avatar, userid: u.id }
    } else {
      user.value = null
    }
  } catch {
    user.value = null
  }
}
let ws: WebSocket | null = null;
const messages = ref<Msg[]>([]);
const connected = ref(false);
const loading = ref(false);
const input = ref('');
const scrollArea = ref<HTMLElement | null>(null);


const onlineCount = ref(0);

const showProfile = ref(false)
const profileUser = ref(null)
function openProfile(user) {
  profileUser.value = user
  showProfile.value = true
}

const showEmoji = ref(false)
function addEmoji(emoji) {
  input.value += emoji.i
  showEmoji.value = false
}
function handleFileChange(e) {
  const file = e.target.files[0]
  if (file) sendImage(file)
}

async function sendImage(file) {
  const url = await uploadFile(file);
  if (!url) return;

  // 判断文件类型，生成 markdown
  const ext = file.name.split('.').pop()?.toLowerCase() || '';
  let markdown = '';
  if (['jpg','jpeg','png','gif','webp','bmp'].includes(ext)) {
    markdown = `![图片](${url})`;
  } else if (['mp4','webm','ogg'].includes(ext)) {
    // 自定义视频标签（大部分 markdown 解析器不支持视频，可加自定义 html 标签）
    markdown = `<video controls style="max-width:80vw;max-height:80vh;background:#111;border-radius:10px;" src="${url}" />`;
  } else {
    // 其它文件类型用链接
    markdown = `[文件](${url})`;
  }
  // 自动补到输入框
  input.value += (input.value ? '\n' : '') + markdown + '\n';
}

function connectWS() {
  loading.value = true;
  ws = new WebSocket(wsUrl);
  ws.onopen = () => {
    connected.value = true;
    loading.value = false;
  };
  ws.onclose = () => {
    connected.value = false;
    loading.value = false;
    setTimeout(() => {
      connectWS();
    }, 3000);
  };
  ws.onerror = (err) => {
    console.error('WebSocket error:', err);
  };
  ws.onmessage = async (event) => {
    try {
      const msg = JSON.parse(event.data);
      // 新增：处理在线人数消息
      if (msg.type === 'online') {
        onlineCount.value = msg.count || 0;
        return; // 避免继续走后面的消息逻辑
      }
      if (Array.isArray(msg)) {
        // 彻底去重并按时间排序
        const map = new Map();
        msg.forEach(m => map.set(m.id, m));
        messages.value = Array.from(map.values())
            .sort((a, b) => new Date(a.created).getTime() - new Date(b.created).getTime());
      } else if (msg && msg.id) {
        // 新消息：如果已存在则跳过，否则 push
        if (!messages.value.some(m => m.id === msg.id)) {
          messages.value.push(msg);
          // 可选：保证数组顺序
          messages.value.sort((a, b) => new Date(a.created).getTime() - new Date(b.created).getTime());
        }
      }
      await nextTick();
      scrollArea.value?.scrollTo(0, scrollArea.value.scrollHeight);
    } catch (e) {
      console.error('Failed to parse message:', e);
    }
  };
}

function send() {
  if (!user.value) {
    connected.value = false;
    return; // 未登录不连 ws
  }
  if (!input.value.trim()) return;
  if (ws && ws.readyState === WebSocket.OPEN) { // 0 CONNECTING 1 OPEN  2 CLOSING 3 CLOSED
    let msg;
    msg = {
      id: (Date.now() + Math.random()).toString(),
      type: 'chat',
      user: user.value?.name || '游客',
      userid: user.value?.userid || '', // 关键修正！
      avatar: user.value?.avatar || '/static/default-avatar.png',
      text: input.value,
      created: new Date().toISOString()
    };
    ws.send(JSON.stringify(msg)); // 发送消息
    input.value = '';
  }
}

onMounted(() => {
  loadUser();
  if (user.value) {
    connectWS();
  } else {
    loginDialog.value = true; // <-- 新增这一行，未登录自动弹窗
  }
});
watch(
    () => messages.value.length,
    async () => {
      // 新消息后自动滚动到底
      await nextTick();
      if (scrollArea.value) {
        scrollArea.value.scrollTo({
          top: scrollArea.value.scrollHeight,
          behavior: 'smooth',
        });
      }
    }
);

// 新增文件上传方法
// 新增文件上传方法（用 request 适配你的 API 工具）
async function uploadFile(file) {
  const form = new FormData();
  form.append('file', file);
  try {
    const data = await request.post('/api/upload/', form, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': 'Bearer ' + (localStorage.getItem('token') || '')
      }
    });
    let url = data.data.url || '';
    if (url && !/^https?:\/\//.test(url)) {
      // 自动拼接环境变量前缀，去掉多余 /
      const base = import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, '') || '';
      url = base + (url.startsWith('/') ? url : '/' + url);
    }
    return url;
  } catch (e) {
    console.error('图片上传失败', e);
    return '';
  }
}

</script>
<style>
.fade-slide-enter-active {
  transition: all 0.4s cubic-bezier(.61,-0.04,.36,1.07);
}
.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(30px) scale(0.95);
}
.fade-slide-enter-to {
  opacity: 1;
  transform: translateY(0) scale(1);
}
</style>