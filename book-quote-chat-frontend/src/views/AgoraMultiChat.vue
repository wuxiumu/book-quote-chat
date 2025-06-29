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
      <div class="video-box">
        <div class="username">我 (本地)</div>
        <video ref="localVideo" autoplay muted playsinline></video>
      </div>
      <div v-for="user in remoteUsers" :key="user.uid" class="video-box">
        <div class="username">用户: {{ user.uid }}</div>
        <video :ref="setRemoteVideoRef(user.uid)" autoplay playsinline></video>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, nextTick } from 'vue'
import AgoraRTC from 'agora-rtc-sdk-ng'
import { ElMessage } from 'element-plus'

const APP_ID = import.meta.env.VITE_AGORA_APP_ID
const CHANNEL = import.meta.env.VITE_AGORA_CHANNEL
const TOKEN = import.meta.env.VITE_AGORA_TOKEN
// 随机生成一个 UID
const UID =  Math.floor(Math.random() * 1000000)

const cameras = ref([])
const microphones = ref([])
const selectedCameraId = ref('')
const selectedMicId = ref('')
const selectedResolution = ref('480p')

const joined = ref(false)
let client = null
let localTracks = []
const localVideo = ref(null)

// 用数组管理远端用户
const remoteUsers = reactive([])

// 用Map存每个远端用户的 video dom ref
const remoteVideoRefs = reactive(new Map())
const setRemoteVideoRef = uid => el => {
  if (el) remoteVideoRefs.set(uid, el)
}

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
  localTracks[1].play(localVideo.value)
  console.log('[Agora] 本地轨道已切换', config)
}

async function joinChannel() {
  client = AgoraRTC.createClient({ mode: 'rtc', codec: 'vp8' })

  // 订阅远端
  client.on('user-published', async (user, mediaType) => {
    await client.subscribe(user, mediaType)
    let exists = remoteUsers.find(u => u.uid === user.uid)
    if (!exists) {
      remoteUsers.push(user)
      ElMessage.success(`用户 ${user.uid} 加入频道`)
    }

    // 等dom生成
    await nextTick()
    if (mediaType === 'video' && user.videoTrack && remoteVideoRefs.has(user.uid)) {
      user.videoTrack.play(remoteVideoRefs.get(user.uid))
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
      remoteVideoRefs.delete(user.uid) // 清理 ref
      ElMessage.info(`用户 ${user.uid} 离开频道`)
    }
  })

  await client.join(APP_ID, CHANNEL, TOKEN, UID)
  joined.value = true
  await recreateLocalTracks()
  console.log('[Agora] 已加入频道:', CHANNEL, '我的UID:', UID)
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
  // 彻底清理所有远端窗口
  remoteUsers.splice(0)
  remoteVideoRefs.clear()
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
}
.video-box {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.video-box video {
  width: 280px;
  height: 210px;
  border-radius: 10px;
  background: #111;
  border: 1.5px solid #aaa;
  margin-bottom: 6px;
}
.username {
  color: #22c55e;
  font-size: 14px;
  text-align: center;
  margin-bottom: 2px;
}
@media (max-width: 600px) {
  .video-grid {
    flex-direction: column;
    gap: 8px;
  }
  .video-box video {
    width: 98vw;
    height: 180px;
  }
}
</style>