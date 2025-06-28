<template>
  <div class="flex flex-col h-[100dvh] bg-base-200">
    <!-- 顶部 -->
    <div class="sticky top-0 z-10 bg-base-100/90 backdrop-blur border-b px-4 py-2 flex items-center gap-3">
      <img :src="friend?.avatar" class="w-9 h-9 rounded-full border object-cover" />
      <span class="font-bold text-lg">{{ friend?.name || '好友' }}</span>
      <button class="btn btn-xs btn-ghost ml-auto" @click="goBack">返回</button>
      <button class="btn btn-xs btn-ghost" @click="loadHistory">加载历史</button>
    </div>
    <!-- 消息区 -->
    <div ref="msgArea" class="flex-1 overflow-y-auto px-3 py-4 space-y-3">
      <template v-for="msg in messages" :key="msg.id">
        <div
            :class="msg.self ? 'flex justify-end' : 'flex justify-start'"
            class="items-end group relative"
            @contextmenu.prevent="showRevokeMenu($event, msg)"
            @touchstart="handleTouchStart($event, msg)"
            @touchend="handleTouchEnd"
        >
          <img v-if="!msg.self" :src="friend?.avatar" class="w-7 h-7 rounded-full mr-2" />
          <div
            :class="[
              msg.self ? 'bg-pink-100 text-right' : 'bg-white text-left',
              'rounded-xl px-4 py-2 shadow max-w-[70%] break-words relative',
              msg.revoke ? 'italic text-gray-500' : ''
            ]"
            style="word-break:break-word;"
          >
            <template v-if="msg.revoke">
              消息已撤回
            </template>
            <template v-else>
              <template v-if="msg.type === 'text'">
                {{ msg.text }}
              </template>
              <template v-else-if="msg.type === 'emoji'">
                <span class="text-3xl leading-none select-none">{{ msg.text }}</span>
              </template>
              <template v-else-if="msg.type === 'image'">
                <img :src="msg.text" alt="图片" class="max-w-full rounded-md cursor-pointer" @click="previewImage(msg.text)" />
              </template>
            </template>
            <div v-if="msg.self" class="text-xs text-gray-400 mt-1 select-none">
              {{ msg.read ? '已读' : '未读' }}
            </div>
            <div class="text-xs text-gray-400 mt-1 select-none">
              {{ formatTime(msg.time) }}
            </div>
          </div>
          <img v-if="msg.self" :src="myAvatar" class="w-7 h-7 rounded-full ml-2" />
        </div>
      </template>
    </div>
    <!-- 输入区 -->
    <form
        class="flex gap-2 p-3 bg-base-100 border-t sticky bottom-0 z-10 flex-wrap items-center"
        @submit.prevent="send"
    >
      <button type="button" class="btn btn-ghost btn-sm" @click="toggleEmojiPicker" title="表情">
        😊
      </button>
      <button type="button" class="btn btn-ghost btn-sm" @click="openImageDialog" title="图片">
        🖼️
      </button>
      <input
          v-model="input"
          class="input input-bordered flex-1 min-w-[120px]"
          placeholder="说点什么..."
          maxlength="200"
          autocomplete="off"
          ref="inputRef"
        />
      <button class="btn btn-primary" type="submit" :disabled="!input.trim()">发送</button>
    </form>
    <!-- emoji选择弹窗 -->
    <div v-if="showEmoji" class="absolute bottom-[56px] left-4 bg-white border rounded-md shadow-lg p-2 max-w-xs max-h-40 overflow-y-auto z-20 flex flex-wrap gap-1">
      <button
        v-for="emoji in emojis"
        :key="emoji"
        type="button"
        class="text-2xl p-1 hover:bg-gray-100 rounded"
        @click="selectEmoji(emoji)"
      >{{ emoji }}</button>
    </div>
    <!-- 图片上传弹窗 -->
    <dialog ref="imageDialog" class="modal">
      <form method="dialog" class="modal-box p-4">
        <h3 class="font-bold text-lg mb-2">发送图片</h3>
        <input type="file" accept="image/*" @change="handleImageUpload" class="mb-2" />
        <input
          v-model="imageUrl"
          type="text"
          placeholder="或输入图片链接"
          class="input input-bordered w-full mb-4"
        />
        <div class="modal-action flex justify-end gap-2">
          <button type="button" class="btn btn-sm btn-ghost" @click="closeImageDialog">取消</button>
          <button type="button" class="btn btn-sm btn-primary" @click="sendImage" :disabled="!imageUrl.trim()">发送</button>
        </div>
      </form>
    </dialog>
    <!-- 撤回菜单 -->
    <div
      v-if="revokeMenu.visible"
      :style="{ top: revokeMenu.y + 'px', left: revokeMenu.x + 'px' }"
      class="absolute bg-white border rounded shadow p-1 z-30"
      @click="revokeMessage(revokeMenu.msg)"
      @mousedown.prevent
    >
      <button class="btn btn-xs btn-ghost">撤回</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
const myAvatar = 'https://api.dicebear.com/7.x/miniavs/svg?seed=Aric';

// 假好友列表
const allFriends = [
  { id: 'f1', name: '书友A', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=A' },
  { id: 'f2', name: '书友B', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=B' },
  { id: 'f3', name: '书友C', avatar: 'https://api.dicebear.com/7.x/miniavs/svg?seed=C' },
];

const friend = computed(() =>
    allFriends.find(f => f.id === route.params.id)
);

interface Msg {
  id: string;
  text: string;
  self: boolean;
  type: 'text' | 'image' | 'emoji';
  read: boolean;
  revoke: boolean;
  time: number;
}

const messages = ref<Msg[]>([
  { id: 'm1', text: '你好呀，这里是私聊！', self: false, type: 'text', read: true, revoke: false, time: Date.now() - 600000 },
  { id: 'm2', text: 'Hi！很高兴认识你！', self: true, type: 'text', read: true, revoke: false, time: Date.now() - 590000 },
  { id: 'm3', text: '最近在读什么书？', self: false, type: 'text', read: false, revoke: false, time: Date.now() - 580000 },
  { id: 'm4', text: '刚翻完《活着》，超级感慨~', self: true, type: 'text', read: false, revoke: false, time: Date.now() - 570000 },
]);

const input = ref('');
const msgArea = ref<HTMLElement | null>(null);
const inputRef = ref<HTMLInputElement | null>(null);

function send() {
  if (!input.value.trim()) return;
  messages.value.push({
    id: 'm' + Date.now(),
    text: input.value,
    self: true,
    type: 'text',
    read: false,
    revoke: false,
    time: Date.now(),
  });
  input.value = '';
  nextTick(() => {
    scrollToBottom();
  });
}

function scrollToBottom() {
  if (msgArea.value) {
    msgArea.value.scrollTo({ top: msgArea.value.scrollHeight, behavior: 'smooth' });
  }
}

function goBack() {
  router.back();
}

onMounted(() => {
  nextTick(() => {
    if (msgArea.value) {
      msgArea.value.scrollTo(0, msgArea.value.scrollHeight);
    }
  });
});

// 时间格式化
function formatTime(timestamp: number) {
  const date = new Date(timestamp);
  const now = new Date();
  const diff = now.getTime() - timestamp;
  if (diff < 60000) return '刚刚';
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前';
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

// 撤回消息菜单及逻辑
const revokeMenu = ref<{ visible: boolean; x: number; y: number; msg: Msg | null }>({ visible: false, x: 0, y: 0, msg: null });
let touchTimer: number | null = null;

function showRevokeMenu(e: MouseEvent, msg: Msg) {
  if (!msg.self || msg.revoke) return;
  revokeMenu.value = { visible: true, x: e.clientX, y: e.clientY, msg };
  window.addEventListener('click', hideRevokeMenu);
}

function hideRevokeMenu() {
  revokeMenu.value.visible = false;
  revokeMenu.value.msg = null;
  window.removeEventListener('click', hideRevokeMenu);
}

function revokeMessage(msg: Msg | null) {
  if (!msg) return;
  const target = messages.value.find(m => m.id === msg.id);
  if (target) {
    target.revoke = true;
    target.text = '';
  }
  hideRevokeMenu();
}

// 触摸长按处理
function handleTouchStart(e: TouchEvent, msg: Msg) {
  if (!msg.self || msg.revoke) return;
  touchTimer = window.setTimeout(() => {
    const touch = e.touches[0];
    revokeMenu.value = { visible: true, x: touch.clientX, y: touch.clientY, msg };
  }, 600);
}
function handleTouchEnd() {
  if (touchTimer) {
    clearTimeout(touchTimer);
    touchTimer = null;
  }
}

// 加载历史消息（模拟）
function loadHistory() {
  const now = Date.now();
  const historyMsgs: Msg[] = [
    { id: 'h' + (now - 300000), text: '这是历史消息1', self: false, type: 'text', read: true, revoke: false, time: now - 300000 },
    { id: 'h' + (now - 290000), text: '😀', self: true, type: 'emoji', read: true, revoke: false, time: now - 290000 },
    { id: 'h' + (now - 280000), text: 'https://picsum.photos/200/150?random=1', self: false, type: 'image', read: true, revoke: false, time: now - 280000 },
  ];
  messages.value = [...historyMsgs, ...messages.value];
  nextTick(() => {
    if (msgArea.value) {
      msgArea.value.scrollTop = 0;
    }
  });
}

// 预览图片（简单打开新窗口）
function previewImage(src: string) {
  window.open(src, '_blank');
}

// 图片发送相关
const imageDialog = ref<HTMLDialogElement | null>(null);
const imageUrl = ref('');

function openImageDialog() {
  imageUrl.value = '';
  imageDialog.value?.showModal();
}
function closeImageDialog() {
  imageDialog.value?.close();
}
function sendImage() {
  if (!imageUrl.value.trim()) return;
  messages.value.push({
    id: 'm' + Date.now(),
    text: imageUrl.value.trim(),
    self: true,
    type: 'image',
    read: false,
    revoke: false,
    time: Date.now(),
  });
  imageUrl.value = '';
  closeImageDialog();
  nextTick(() => {
    scrollToBottom();
  });
}
function handleImageUpload(e: Event) {
  const target = e.target as HTMLInputElement;
  if (target.files && target.files.length > 0) {
    const file = target.files[0];
    const reader = new FileReader();
    reader.onload = (ev) => {
      if (typeof ev.target?.result === 'string') {
        imageUrl.value = ev.target.result;
      }
    };
    reader.readAsDataURL(file);
  }
}

// emoji发送相关
const showEmoji = ref(false);
const emojis = ['😀','😂','😍','😎','😭','👍','🙏','🎉','🔥','💯','🥳','🤔','😴','😡','🤩'];

function toggleEmojiPicker() {
  showEmoji.value = !showEmoji.value;
}

function selectEmoji(emoji: string) {
  // 如果输入框有焦点，则插入光标位置，否则直接发送
  if (inputRef.value && document.activeElement === inputRef.value) {
    const el = inputRef.value;
    const start = el.selectionStart || 0;
    const end = el.selectionEnd || 0;
    input.value = input.value.slice(0, start) + emoji + input.value.slice(end);
    nextTick(() => {
      el.selectionStart = el.selectionEnd = start + emoji.length;
      el.focus();
    });
  } else {
    messages.value.push({
      id: 'm' + Date.now(),
      text: emoji,
      self: true,
      type: 'emoji',
      read: false,
      revoke: false,
      time: Date.now(),
    });
    nextTick(() => {
      scrollToBottom();
    });
  }
  showEmoji.value = false;
}

</script>