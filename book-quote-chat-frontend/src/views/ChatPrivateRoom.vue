<template>
  <div class="chatroom" style="display:flex;flex-direction:column;height:100vh;">
    <div class="rtc-area" style="display: flex; gap: 16px; margin-bottom: 10px;">
      <video ref="localVideo" autoplay muted playsinline width="140" height="100" style="border-radius:8px; background:#333; cursor:pointer;" @click="showFullScreenVideo('local')" />
      <video ref="remoteVideo" autoplay playsinline width="140" height="100" style="border-radius:8px; background:#333; cursor:pointer;" @click="showFullScreenVideo('remote')" />
      <button
        :disabled="rtcState === 'calling' || rtcState === 'incall'"
        :class="['rtc-btn', rtcState]"
        @click="startCall"
        v-if="rtcState !== 'incall'"
      >
        <span v-if="rtcState === 'idle'">📞 发起通话</span>
        <span v-else-if="rtcState === 'calling'">
          <span class="loading"></span> 等待对方响应...
        </span>
      </button>
      <button
        class="rtc-btn hangup"
        v-if="rtcState === 'incall'"
        @click="hangupCall"
      >挂断</button>
    </div>
    <div v-if="rtcState === 'incall'" class="rtc-tip">通话中（可随时挂断）</div>
    <div v-if="rtcErr" class="rtc-err-tip">{{ rtcErr }}</div>
    <div class="header" style="padding:8px 16px;font-weight:bold;">
      超级简易私聊
    </div>
    <div class="flex-1 overflow-y-auto p-4 space-y-3 messages" ref="scrollArea">
      <transition-group name="fade-slide" tag="div">
        <ChatMessage
            v-for="msg in messages"
            :key="msg.id"
            :msg="msg"
            :self="msg.from === myId"
        />
      </transition-group>
    </div>
    <div class="footer flex gap-2 p-3 bg-base-100 border-t sticky bottom-0 z-10">
      <input v-model="input" @keyup.enter="sendText" placeholder="说点什么..." style="flex:1;padding:6px 10px;" maxlength="200" />
      <input type="file" accept="audio/*" @change="sendAudio" style="display:none;" ref="audioInput" />
      <button @click="audioInput.click()" title="发送语音">🎤</button>
      <input type="file" accept="video/*" @change="sendVideo" style="display:none;" ref="videoInput" />
      <button @click="videoInput.click()" title="发送视频">🎬</button>
      <button @click="sendText" style="padding:4px 14px;">发送</button>
    </div>
    <!-- 全屏视频遮罩 -->
    <div v-if="fullScreenVideo" class="fullscreen-video-mask" @click="closeFullScreenVideo">
      <video
        v-if="fullScreenVideo==='local'"
        ref="fullScreenLocalVideo"
        autoplay
        controls
        muted
        playsinline
        style="max-width:96vw;max-height:96vh;background:#111;border-radius:10px;"
        @click.stop
      ></video>
      <video
        v-if="fullScreenVideo==='remote'"
        ref="fullScreenRemoteVideo"
        autoplay
        controls
        playsinline
        style="max-width:96vw;max-height:96vh;background:#111;border-radius:10px;"
        @click.stop
      ></video>
      <button class="close-btn" @click.stop="closeFullScreenVideo">关闭</button>
    </div>
  </div>
</template>

<script setup>
import {ref, nextTick, onMounted, onBeforeUnmount, computed, watch} from 'vue'
import { useRoute } from 'vue-router'
import ChatMessage from '@/components/ChatMessage.vue'

const fullScreenVideo = ref('')
const fullScreenLocalVideo = ref(null)
const fullScreenRemoteVideo = ref(null)

function showFullScreenVideo(which) {
  fullScreenVideo.value = which
  setTimeout(() => {
    if (which === 'local' && fullScreenLocalVideo.value && localVideo.value) {
      fullScreenLocalVideo.value.srcObject = localVideo.value.srcObject
    }
    if (which === 'remote' && fullScreenRemoteVideo.value && remoteVideo.value) {
      fullScreenRemoteVideo.value.srcObject = remoteVideo.value.srcObject
    }
  }, 20)
}
function closeFullScreenVideo() {
  fullScreenVideo.value = ''
  if (fullScreenLocalVideo.value) fullScreenLocalVideo.value.srcObject = null
  if (fullScreenRemoteVideo.value) fullScreenRemoteVideo.value.srcObject = null
}

const scrollArea = ref<HTMLElement | null>(null);
const localVideo = ref(null)
const remoteVideo = ref(null)
let localStream = null
let remoteStream = null
let pc = null
let inCall = false

const rtcState = ref('idle') // idle, calling, incall
const rtcErr = ref('')
watch(rtcErr, val => { if (val) setTimeout(() => rtcErr.value = '', 3000) })

function startCall() {
  if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
    rtcErr.value = '当前浏览器不支持音视频通话'
    return
  }
  rtcState.value = 'calling'
  navigator.mediaDevices.getUserMedia({ video: true, audio: true })
    .then(stream => {
      rtcErr.value = ''
      localStream = stream
      localVideo.value.srcObject = stream
      pc = createPeerConnection()
      stream.getTracks().forEach(track => pc.addTrack(track, stream))
      pc.createOffer().then(offer => {
        pc.setLocalDescription(offer)
        ws.value?.send(JSON.stringify({ type: 'rtc-offer', offer, from: myId, to: otherId }))
      })
      setTimeout(() => {
        if (rtcState.value === 'calling') {
          hangupCall()
          rtcErr.value = '对方未响应，已自动挂断'
        }
      }, 15000)
    })
    .catch(err => {
      rtcErr.value = '获取摄像头/麦克风失败'
    })
}
function hangupCall() {
  console.log('[WebRTC] hangupCall 被调用')
  // 停止本地流
  if (localStream) {
    localStream.getTracks().forEach(track => track.stop())
    localStream = null
  }
  // 停止远端流
  if (remoteStream) {
    remoteStream.getTracks && remoteStream.getTracks().forEach(track => track.stop())
    remoteStream = null
  }
  // 关闭视频画面
  if (localVideo.value) localVideo.value.srcObject = null
  if (remoteVideo.value) remoteVideo.value.srcObject = null
  // 关闭 RTCPeerConnection
  if (pc) {
    pc.close()
    pc = null
  }
  inCall = false
  rtcState.value = 'idle'
  ws.value?.send(JSON.stringify({ type: 'rtc-hangup', from: myId, to: otherId }))
  console.log('[WebRTC] 发送 rtc-hangup', { from: myId, to: otherId })
}
function createPeerConnection() {
  console.log('[WebRTC] createPeerConnection()')
  const newPc = new RTCPeerConnection({
    iceServers: [{ urls: 'stun:stun.l.google.com:19302' }]
  })
  newPc.onicecandidate = (e) => {
    if (e.candidate) {
      ws.value?.send(JSON.stringify({ type: 'rtc-ice', candidate: e.candidate, from: myId, to: otherId }))
      console.log('[WebRTC] 发送 rtc-ice', { candidate: e.candidate, from: myId, to: otherId })
    }
  }
  newPc.oniceconnectionstatechange = () => {
    console.log('[WebRTC] ICE 连接状态', newPc.iceConnectionState)
  }
  newPc.ontrack = (e) => {
    remoteStream = e.streams[0]
    remoteVideo.value.srcObject = remoteStream
    console.log('[WebRTC] ontrack 收到远端流:', remoteStream)
  }
  return newPc
}

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';
const token = localStorage.getItem('token')
const wsUrl = `${WS_URL}?token=${encodeURIComponent(token||'')}`;

const input = ref('')
const messages = ref([])
const route = useRoute()
const id = route.params.id || ''
const [idA, idB] = id.split('_')   // 都是string
const sessionKey = genSessionKey(idA, idB)  // 传 string

const audioInput = ref()
const videoInput = ref()

// 取本地登录用户
const myUser = JSON.parse(localStorage.getItem('user') || '{}')
if (!myUser.id) {
  window.location.href = '/login'
}
const otherId = (idA === myUser.id) ? idB : idA
const myId = myUser.id
const ws = ref(null)

onMounted(() => {
  console.log('[WS] onMounted: 即将连接 WebSocket')
  ws.value = new window.WebSocket(
      `${wsUrl}&sessionKey=${sessionKey}`
  )
  ws.value.onopen = () => console.log('[WS] WebSocket已连接')
  ws.value.onclose = () => console.log('[WS] WebSocket已断开')
  ws.value.onerror = (e) => console.error('[WS] WebSocket出错', e)
  ws.value.onmessage = (e) => {
    const msg = JSON.parse(e.data)
    // 优先处理在线人数类型
    if (msg.type === 'online') {
      // console.log('[WS] 收到在线人数', msg)
      return
    }
    console.log('[WS] onmessage:', e.data)
    // 历史数组
    if (Array.isArray(msg)) {
      const map = new Map();
      msg.forEach(m => map.set(m.id, m));
      messages.value = Array.from(map.values())
          .sort((a, b) => new Date(a.created).getTime() - new Date(b.created).getTime());
      console.log('[WS] 收到历史消息', messages.value)
      return;
    }
    // 单条消息
    if (msg.type === 'chat') {
      console.log('[WS] 收到 chat 消息', msg)
      if (!msg.userid && msg.from) msg.userid = msg.from
      if (!messages.value.some(m => m.id === msg.id)) {
        messages.value.push(msg)
        messages.value.sort((a, b) => new Date(a.created).getTime() - new Date(b.created).getTime())
        console.log('[WS] 新增 chat 消息', msg)
      }
      nextTick(() => scrollArea.value?.scrollTo(0, scrollArea.value.scrollHeight));
      return;
    }
    if (msg.type === 'error') {
      console.error('[WS] 私聊错误', msg)
    }

    // RTC 信令
    console.log('[WebRTC] 收到信令', msg)
    if (msg.type === 'rtc-offer' && msg.to === myId) {
      console.log('[WebRTC] 收到 rtc-offer', msg)
      navigator.mediaDevices.getUserMedia({ video: true, audio: true })
          .then(stream => {
            localStream = stream
            localVideo.value.srcObject = stream
            pc = createPeerConnection()
            stream.getTracks().forEach(track => {
              pc.addTrack(track, stream)
              console.log('[WebRTC] 被叫 addTrack', track)
            })
            pc.setRemoteDescription(new RTCSessionDescription(msg.offer))
            pc.createAnswer().then(answer => {
              pc.setLocalDescription(answer)
              ws.value?.send(JSON.stringify({ type: 'rtc-answer', answer, from: myId, to: otherId }))
              console.log('[WebRTC] 发送 rtc-answer', { answer, from: myId, to: otherId })
            })
            inCall = true
          })
          .catch(err => {
            console.error('[WebRTC] 被叫 getUserMedia 失败', err)
          })
    }
    if (msg.type === 'rtc-answer' && msg.to === myId) {
      console.log('[WebRTC] 收到 rtc-answer', msg)
      pc?.setRemoteDescription(new RTCSessionDescription(msg.answer))
      rtcState.value = 'incall'
    }
    if (msg.type === 'rtc-ice' && msg.to === myId) {
      console.log('[WebRTC] 收到 rtc-ice', msg)
      pc?.addIceCandidate(new RTCIceCandidate(msg.candidate))
    }
    if (msg.type === 'rtc-hangup' && msg.to === myId) {
      console.log('[WebRTC] 收到 rtc-hangup', msg)
      hangupCall()
      rtcState.value = 'idle'
    }
  }
})
onBeforeUnmount(() => {
  ws.value?.close()
  console.log('[WS] 组件销毁，WebSocket 已关闭')
})

function sendText() {
  if (!input.value.trim()) return
  let msg = {
    sessionKey,
    id: (Date.now() + Math.random()).toString(),
    user: myUser.name || '游客',
    userid: myUser.id || '', // 关键修正！
    avatar: myUser.avatar || '/static/default-avatar.png',
    from: myUser.id,
    to: otherId,
    text: input.value,
    type: 'chat',
    subType: 'text',
    created: Date.now()
  }
  console.log('发送私聊消息', msg)
  ws.value?.send(JSON.stringify(msg))
  input.value = ''
}
function sendAudio(e) {
  const file = e.target.files[0]
  if (!file) return
  const formData = new FormData()
  formData.append('file', file)
  fetch('/api/uploads/', { method: 'POST', body: formData })
      .then(resp => resp.json())
      .then(data => {
        if (data.url) {
          ws.value?.send(JSON.stringify({
            sessionKey,
            from: myUser.id,
            to: otherId,
            text: data.url,
            type: 'chat',
            subType: 'audio',
            avatar: myUser.avatar || '/static/default-avatar.png',
            nickname: myUser.nickname || myUser.name || ''
          }))
        }
      })
}
function sendVideo(e) {
  const file = e.target.files[0]
  if (!file) return
  const formData = new FormData()
  formData.append('file', file)
  fetch('/api/uploads/', { method: 'POST', body: formData })
      .then(resp => resp.json())
      .then(data => {
        if (data.url) {
          ws.value?.send(JSON.stringify({
            sessionKey,
            from: myUser.id,
            to: otherId,
            text: data.url,
            type: 'chat',
            subType: 'video',
            avatar: myUser.avatar || '/static/default-avatar.png',
            nickname: myUser.nickname || myUser.name || ''
          }))
        }
      })
}
function addMsg({ id, text, type, from, avatar, nickname, created }) {
  messages.value.push({
    id: id ?? Date.now() + Math.random(),
    text: text ?? '',
    type: type || 'text',
    from: from ?? '',
    avatar: avatar || '',
    nickname: nickname || '',
    created: created || Date.now()
  })
  if (messages.value.length > 20) messages.value.splice(0, messages.value.length - 20)
  nextTick(() => {
    const container = document.querySelector('.flex-1')
    if (container) container.scrollTop = container.scrollHeight
  })
}
function genSessionKey(id1, id2) {
  // 保证 key 顺序唯一，小的在前
  return [id1, id2].sort().join('_')
}
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
.rtc-btn {
  font-size: 17px;
  border: none;
  padding: 9px 30px;
  border-radius: 22px;
  background: #53bb56;
  color: #fff;
  font-weight: bold;
  cursor: pointer;
  margin-right: 12px;
  transition: background 0.2s;
  box-shadow: 0 2px 8px #0002;
}
.rtc-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  opacity: .8;
}
.rtc-btn.hangup {
  background: #d44;
}
.rtc-tip {
  margin: -6px 0 8px 0;
  font-size: 15px;
  color: #2196f3;
}
.rtc-err-tip {
  margin-top: 6px;
  color: #fff;
  background: #e53935;
  padding: 4px 10px;
  border-radius: 8px;
  font-size: 14px;
  animation: shake 0.18s;
}
.loading {
  display:inline-block;
  width:18px;height:18px;
  border-radius:50%;
  border:3px solid #fff;
  border-top-color: #2196f3;
  margin-right:4px;
  vertical-align:-3px;
  animation:spin .8s linear infinite;
}
@keyframes spin { to { transform:rotate(360deg) } }
@keyframes shake { 10%,90% {transform: translateX(-1px);} 20%,80%{transform: translateX(2px);} 30%,50%,70%{transform: translateX(-4px);} 40%,60%{transform: translateX(4px);} }

.fullscreen-video-mask {
  position: fixed;
  left:0;top:0;right:0;bottom:0;
  background: rgba(0,0,0,0.88);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  z-index: 9999;
  animation: fadein 0.28s;
}
.close-btn {
  margin-top: 20px;
  padding: 8px 20px;
  font-size: 17px;
  background: #fff;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  opacity: 0.92;
}
@keyframes fadein { from { opacity:0; } to { opacity:1; } }
</style>