<template>
  <div class="agora-demo">
    <h2>Agora 多人视频通话 Demo</h2>
    <div class="toolbar">
      <select v-model="selectedCameraId">
        <option v-for="cam in cameras" :key="cam.deviceId" :value="cam.deviceId">{{ cam.label || '摄像头' }}</option>
      </select>
      <select v-model="selectedMicId">
        <option v-for="mic in microphones" :key="mic.deviceId" :value="mic.deviceId">{{ mic.label || '麦克风' }}</option>
      </select>
      <select v-model="selectedResolution">
        <option value="360p">360p</option>
        <option value="480p">480p</option>
        <option value="720p">720p</option>
        <option value="1080p">1080p</option>
      </select>
      <button class="btn btn-sm btn-outline" @click="joinChannel" :disabled="joined">加入频道</button>
      <button class="btn btn-sm btn-outline" @click="leaveChannel" :disabled="!joined">离开频道</button>
    </div>

    <div class="video-grid">
      <!-- 本地视频窗口 -->
      <div
          class="video-box"
          :class="{ fullscreen: fullscreenId === 'local' }"
          @click="toggleFullscreen('local')"
      >
        <div class="username">我 (本地)</div>
        <video
            ref="localVideo"
            autoplay
            playsinline
            :muted="isMuted['local']"
            :style="{ transform: `rotate(${rotateMap['local'] || 0}deg)` }"
        ></video>
        <button
            class="mute-btn"
            :class="{ on: !isMuted['local'] }"
            @click.stop="toggleMute('local')"
        >
          <span v-if="isMuted['local']">🔇</span>
          <span v-else>🔊</span>
          {{ isMuted['local'] ? '静音' : '已开声' }}
        </button>
        <button class="rotate-btn" @click.stop="rotateVideo('local')">⟳</button>
        <button class="fullscreen-btn" @click.stop="enterRealFullscreen('local')">⛶</button>
      </div>
      <!-- 远端用户窗口 -->
      <div
          v-for="user in remoteUsers"
          :key="user.uid"
          class="video-box"
          :class="{ fullscreen: fullscreenId === user.uid }"
          @click="toggleFullscreen(user.uid)"
      >
        <div class="username">用户: {{ user.uid }}</div>
        <video
            :ref="setRemoteVideoRef(user.uid)"
            autoplay
            playsinline
            :muted="isMuted[user.uid] ?? true"
            :style="{ transform: `rotate(${rotateMap[user.uid] || 0}deg)` }"
        ></video>
        <button
            class="mute-btn"
            :class="{ on: !isMuted[user.uid] }"
            @click.stop="toggleMute(user.uid)"
        >
          <span v-if="isMuted[user.uid] ?? true">🔇</span>
          <span v-else>🔊</span>
          {{ isMuted[user.uid] ?? true ? '静音' : '已开声' }}
        </button>
        <button class="rotate-btn" @click.stop="rotateVideo(user.uid)">⟳</button>
        <button class="fullscreen-btn" @click.stop="enterRealFullscreen(user.uid)">⛶</button>
      </div>
    </div>
    <!-- 遮罩层（全屏时显示） -->
    <div v-if="fullscreenId" class="fullscreen-mask" @click="toggleFullscreen()"></div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, nextTick } from 'vue'
import AgoraRTC from 'agora-rtc-sdk-ng'
import { ElMessage } from 'element-plus'

const APP_ID = import.meta.env.VITE_AGORA_APP_ID
const CHANNEL = import.meta.env.VITE_AGORA_CHANNEL
const TOKEN = import.meta.env.VITE_AGORA_TOKEN
const UID = Math.floor(Math.random() * 1000000)

const cameras = ref([])
const microphones = ref([])
const selectedCameraId = ref('')
const selectedMicId = ref('')
const selectedResolution = ref('480p')

const joined = ref(false)
let client = null
let localTracks = []
const localVideo = ref(null)

const remoteUsers = reactive([])
const remoteVideoRefs = reactive(new Map())
const setRemoteVideoRef = uid => el => { if (el) remoteVideoRefs.set(uid, el) }

// 静音状态，key = uid 或 'local'
const isMuted = reactive({ local: true })
// 全屏窗口 id（'local' 或远端 uid），null 为无
const fullscreenId = ref(null)

// 旋转角度映射，key = uid 或 'local'
const rotateMap = reactive({})

const RESOLUTION_CONFIG = {
  '360p': { width: 640, height: 360, frameRate: 15 },
  '480p': { width: 854, height: 480, frameRate: 15 },
  '720p': { width: 1280, height: 720, frameRate: 15 },
  '1080p': { width: 1920, height: 1080, frameRate: 15 },
}

onMounted(async () => {
  cameras.value = await AgoraRTC.getCameras()
  microphones.value = await AgoraRTC.getMicrophones()
  if (cameras.value.length) selectedCameraId.value = cameras.value[0].deviceId
  if (microphones.value.length) selectedMicId.value = microphones.value[0].deviceId
})

watch([selectedCameraId, selectedMicId, selectedResolution], async () => {
  if (!joined.value) return
  await recreateLocalTracks()
})

async function recreateLocalTracks() {
  if (localTracks.length) {
    await client.unpublish(localTracks)
    localTracks.forEach(track => track.close())
    localTracks = []
  }
  const config = RESOLUTION_CONFIG[selectedResolution.value]
  localTracks = await AgoraRTC.createMicrophoneAndCameraTracks(
      { microphoneId: selectedMicId.value },
      { cameraId: selectedCameraId.value, encoderConfig: config }
  )
  await client.publish(localTracks)
  // 播放本地视频
  localTracks[1].play(localVideo.value)
  // 设置本地静音状态
  isMuted['local'] = true
  localVideo.value.muted = isMuted['local']
  // 初始化旋转角度为0
  rotateMap['local'] = 0
  console.log('[Agora] 本地轨道已切换', config)
}

async function joinChannel() {
  client = AgoraRTC.createClient({ mode: 'rtc', codec: 'vp8' })

  client.on('user-published', async (user, mediaType) => {
    await client.subscribe(user, mediaType)
    let exists = remoteUsers.find(u => u.uid === user.uid)
    if (!exists) {
      remoteUsers.push(user)
      ElMessage.success(`用户 ${user.uid} 加入频道`)
      // 新用户默认静音
      isMuted[user.uid] = true
      // 初始化旋转角度为0
      rotateMap[user.uid] = 0
    }
    await nextTick()
    if (mediaType === 'video' && user.videoTrack && remoteVideoRefs.has(user.uid)) {
      user.videoTrack.play(remoteVideoRefs.get(user.uid))
      // 设置静音状态
      remoteVideoRefs.get(user.uid).muted = isMuted[user.uid]
    }
    if (mediaType === 'audio' && user.audioTrack) {
      user.audioTrack.play()
    }
    console.log(`[Agora] 远端用户加入: ${user.uid}, type: ${mediaType}`)
  })

  client.on('user-unpublished', user => {
    let idx = remoteUsers.findIndex(u => u.uid === user.uid)
    if (idx !== -1) {
      remoteUsers.splice(idx, 1)
      remoteVideoRefs.delete(user.uid)
      delete isMuted[user.uid]
      delete rotateMap[user.uid]
      ElMessage.info(`用户 ${user.uid} 离开频道`)
    }
  })

  await client.join(APP_ID, CHANNEL, TOKEN, UID)
  joined.value = true
  await recreateLocalTracks()
  console.log('[Agora] 已加入频道:', CHANNEL, '我的UID:', UID)
}

// 切换某窗口静音
function toggleMute(uid) {
  isMuted[uid] = !isMuted[uid]
  // 实时切换 video dom muted 属性
  if (uid === 'local' && localVideo.value) {
    localVideo.value.muted = isMuted[uid]
  }
  if (uid !== 'local' && remoteVideoRefs.has(uid)) {
    const video = remoteVideoRefs.get(uid)
    if (video) video.muted = isMuted[uid]
  }
}

// 切换全屏窗口
function toggleFullscreen(uid) {
  if (fullscreenId.value === uid) {
    fullscreenId.value = null
  } else {
    fullscreenId.value = uid
  }
}

// 旋转视频90度
function rotateVideo(uid) {
  if (!(uid in rotateMap)) {
    rotateMap[uid] = 0
  }
  rotateMap[uid] = (rotateMap[uid] + 90) % 360
}

// 进入真全屏
function enterRealFullscreen(uid) {
  let videoEl = null
  if (uid === 'local') {
    videoEl = localVideo.value
  } else if (remoteVideoRefs.has(uid)) {
    videoEl = remoteVideoRefs.get(uid)
  }
  if (!videoEl) return
  if (videoEl.requestFullscreen) {
    videoEl.requestFullscreen()
  } else if (videoEl.webkitRequestFullscreen) { /* Safari */
    videoEl.webkitRequestFullscreen()
  } else if (videoEl.msRequestFullscreen) { /* IE11 */
    videoEl.msRequestFullscreen()
  }
}

async function leaveChannel() {
  if (!client) return
  await client.unpublish(localTracks)
  await client.leave()
  if (localTracks.length) {
    localTracks.forEach(track => track.close())
    localTracks = []
  }
  if (localVideo.value) localVideo.value.srcObject = null
  remoteUsers.splice(0)
  remoteVideoRefs.clear()
  // 静音/全屏/旋转状态清空
  Object.keys(isMuted).forEach(k => delete isMuted[k])
  Object.keys(rotateMap).forEach(k => delete rotateMap[k])
  fullscreenId.value = null
  joined.value = false
  ElMessage.info('[Agora] 已离开频道')
}
</script>

<style scoped>
.agora-demo {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}
.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  margin-bottom: 16px;
}
.video-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  justify-content: center;
  position: relative;
  z-index: 1;
}
.video-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  transition: all 0.2s;
  cursor: pointer;
  background: #f5f7fa;
  border-radius: 10px;
  border: 2px solid #aaa;
  box-shadow: 0 2px 10px 0 rgba(80,80,80,0.10);
}
.video-box.fullscreen {
  position: fixed;
  top: 50%;
  left: 50%;
  z-index: 99;
  width: 70vw;
  height: 75vh;
  transform: translate(-50%, -50%) scale(1.15);
  background: #222;
  box-shadow: 0 8px 30px 0 rgba(60,60,60,0.35);
  border: 3px solid #22c55e;
}
.video-box video {
  width: 280px;
  height: 210px;
  border-radius: 10px;
  background: #111;
  border: 1.5px solid #aaa;
  margin-bottom: 6px;
  transition: box-shadow .15s, transform 0.3s;
}
.video-box.fullscreen video {
  width: 100%;
  height: 100%;
  border: 0;
  margin-bottom: 0;
}
.username {
  color: #22c55e;
  font-size: 14px;
  text-align: center;
  margin-bottom: 2px;
  user-select: none;
}
.mute-btn {
  position: absolute;
  right: 10px;
  top: 14px;
  z-index: 2;
  font-size: 1.25rem;
  background: #fff;
  color: #333;
  border: 2px solid #22c55e;
  border-radius: 32px;
  padding: 6px 14px;
  cursor: pointer;
  box-shadow: 0 0 10px rgba(30,150,100,0.08);
  font-weight: bold;
  outline: none;
  transition: all .18s;
}
.mute-btn.on {
  color: #22c55e;
  border-color: #22c55e;
  background: #f0fff3;
}
.mute-btn:hover {
  filter: brightness(1.15);
}
.video-box.fullscreen .mute-btn {
  font-size: 2rem;
  right: 25px;
  top: 32px;
  padding: 10px 26px;
}

/* 新增旋转按钮和真全屏按钮样式 */
.rotate-btn, .fullscreen-btn {
  position: absolute;
  bottom: 12px;
  z-index: 3;
  font-size: 1.5rem;
  background: #fff;
  color: #333;
  border: 2px solid #22c55e;
  border-radius: 32px;
  padding: 8px 14px;
  cursor: pointer;
  box-shadow: 0 0 10px rgba(30,150,100,0.08);
  font-weight: bold;
  outline: none;
  transition: all .18s;
  user-select: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  line-height: 1;
}
.rotate-btn:hover, .fullscreen-btn:hover {
  filter: brightness(1.15);
}
.rotate-btn {
  right: 60px;
}
.fullscreen-btn {
  right: 10px;
}
.video-box.fullscreen .rotate-btn,
.video-box.fullscreen .fullscreen-btn {
  font-size: 2.5rem;
  width: 56px;
  height: 56px;
  bottom: 28px;
}
.video-box.fullscreen .rotate-btn {
  right: 90px;
}
.video-box.fullscreen .fullscreen-btn {
  right: 30px;
}

/* 移动端适配 - 按钮更大，位置稍调 */
@media (max-width: 600px) {
  .video-grid {
    flex-direction: column;
    gap: 8px;
  }
  .video-box video {
    width: 98vw;
    height: 180px;
  }
  .video-box.fullscreen {
    width: 97vw;
    height: 62vw;
    min-height: 40vw;
  }
  .rotate-btn, .fullscreen-btn {
    width: 56px;
    height: 56px;
    font-size: 2rem;
    padding: 0;
    bottom: 14px;
  }
  .rotate-btn {
    right: 80px;
  }
  .fullscreen-btn {
    right: 20px;
  }
  .video-box.fullscreen .rotate-btn,
  .video-box.fullscreen .fullscreen-btn {
    font-size: 3rem;
    width: 68px;
    height: 68px;
    bottom: 36px;
  }
  .video-box.fullscreen .rotate-btn {
    right: 110px;
  }
  .video-box.fullscreen .fullscreen-btn {
    right: 40px;
  }
}
.fullscreen-mask {
  position: fixed;
  z-index: 98;
  left: 0; top: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.25);
}
</style>