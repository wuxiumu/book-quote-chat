<template>
  <div>
    <h2>Agora 视频聊天 Demo</h2>
    <div>
      <button class="btn btn-sm btn-outline" @click="joinChannel" :disabled="joined">加入频道</button>
      <button class="btn btn-sm btn-outline" @click="leaveChannel" :disabled="!joined">离开频道</button>
    </div>
    <div style="display: flex;">
      <video ref="localVideo" autoplay muted playsinline style="width:320px; height:240px; border:1px solid #ccc; margin:8px;"></video>
      <video ref="remoteVideo" autoplay playsinline style="width:320px; height:240px; border:1px solid #ccc; margin:8px;"></video>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import AgoraRTC from 'agora-rtc-sdk-ng'

const APP_ID = import.meta.env.VITE_AGORA_APP_ID
const CHANNEL = import.meta.env.VITE_AGORA_CHANNEL
const TOKEN = import.meta.env.VITE_AGORA_TOKEN  // 本地可用 null，正式建议用 Token
const UID = Number(import.meta.env.VITE_AGORA_UID)

console.log('[Agora] APP_ID:', APP_ID, 'CHANNEL:', CHANNEL, 'TOKEN:', TOKEN, 'UID:', UID)

const joined = ref(false)
let client = null
let localTracks = []
const localVideo = ref(null)
const remoteVideo = ref(null)

async function joinChannel() {
  console.log('[Agora] 1. 创建 RTC client...')
  client = AgoraRTC.createClient({ mode: 'rtc', codec: 'vp8' })

  try {
    console.log('[Agora] 2. 加入频道:', { APP_ID, CHANNEL, TOKEN })
    const uid = await client.join(APP_ID, CHANNEL, TOKEN, null)
    console.log('[Agora] 3. 已加入频道, 当前 uid:', uid)

    localTracks = await AgoraRTC.createMicrophoneAndCameraTracks()
    console.log('[Agora] 4. 本地音视频轨道创建成功:', localTracks)

    await client.publish(localTracks)
    console.log('[Agora] 5. 已发布本地轨道:', localTracks)

    // 本地画面
    localTracks[1].play(localVideo.value)
    console.log('[Agora] 6. 播放本地视频')

    // 监听远端
    client.on('user-published', async (user, mediaType) => {
      console.log(`[Agora] 检测到远端用户 [uid=${user.uid}] 发布媒体类型:`, mediaType)
      await client.subscribe(user, mediaType)
      console.log(`[Agora] 已订阅远端用户 [uid=${user.uid}] 的媒体类型:`, mediaType)
      if (mediaType === 'video') {
        user.videoTrack.play(remoteVideo.value)
        console.log('[Agora] 播放远端视频流')
      }
      if (mediaType === 'audio') {
        user.audioTrack.play()
        console.log('[Agora] 播放远端音频流')
      }
    })

    joined.value = true
    console.log('[Agora] 7. 加入频道流程结束，已准备就绪！')
  } catch (error) {
    console.error('[Agora][错误] 加入频道流程异常:', error)
  }
}

async function leaveChannel() {
  if (!client) return
  try {
    console.log('[Agora] 离开频道：准备 unpublish 本地轨道')
    await client.unpublish(localTracks)
    console.log('[Agora] 已取消发布本地轨道')
    await client.leave()
    console.log('[Agora] 已离开频道')

    if (localTracks.length) {
      localTracks.forEach(track => track.close())
      localTracks = []
      console.log('[Agora] 已关闭本地音视频轨道')
    }
    if (localVideo.value) localVideo.value.srcObject = null
    if (remoteVideo.value) remoteVideo.value.srcObject = null
    console.log('[Agora] 已清理本地/远端视频画面')
    joined.value = false
  } catch (error) {
    console.error('[Agora][错误] 离开频道流程异常:', error)
  }
}
</script>