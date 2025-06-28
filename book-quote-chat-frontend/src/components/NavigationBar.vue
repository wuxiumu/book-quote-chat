<template>
  <div class="navbar bg-gradient-to-r from-purple-500 to-pink-500 text-white shadow-lg">
    <div class="flex-1">
      <router-link to="/" class="text-2xl font-bold ml-2">📖 QuoteChat</router-link>
    </div>
    <div class="flex-none">
      <ul class="menu menu-horizontal px-1">
        <li><router-link to="/quotes">金句</router-link></li>
        <li><router-link to="/chat">群聊</router-link></li>
        <li><router-link to="/friends">好友</router-link></li>
      </ul>
    </div>
    <div>
      <button class="btn btn-sm btn-outline" @click="showContact = true">联系我们</button>
    </div>
    <!-- 联系我们弹窗 -->
    <div
        v-if="showContact"
        class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50"
        @click="showContact = false"
    >
      <div
          class="bg-white rounded-xl shadow-lg p-6 w-full max-w-xs max-h-[90vh] overflow-y-auto"
          @click.stop
      >
        <h3 class="text-lg font-bold mb-4 text-center text-pink-500">👨‍💻 联系工程师</h3>
        <ul class="mb-4 space-y-2">
          <li v-for="info in aiInfos" :key="info.label" class="flex items-center">
            <span class="font-semibold min-w-[72px] text-gray-600">{{ info.label }}：</span>
            <span class="flex-1 break-all text-gray-800">{{ info.value }}</span>
            <button v-if="info.copy" class="btn btn-ghost btn-xs ml-2" @click="copy(info.value)">复制</button>
          </li>
        </ul>
        <div class="mb-2 flex justify-center gap-3">
          <img
              referrerpolicy="no-referrer"
              :src="engineer.payWx"
              class="w-20 h-20 rounded-lg border object-contain cursor-pointer transition-transform hover:scale-110"
              alt="微信收款码"
              @click="previewImg = engineer.payWx"
          />
          <img
              referrerpolicy="no-referrer"
              :src="engineer.payAli"
              class="w-20 h-20 rounded-lg border object-contain cursor-pointer transition-transform hover:scale-110"
              alt="支付宝收款码"
              @click="previewImg = engineer.payAli"
          />
        </div>
        <div class="text-center mb-2">
          <a class="btn btn-sm btn-info" :href="engineer.payUrl" target="_blank">赞赏支持</a>
        </div>
        <div class="text-center mb-4">
          <a class="btn btn-sm btn-primary" href="mailto:wuxiumu@163.com?subject=友情链接申请" target="_blank">友情链接申请</a>
        </div>
        <button class="btn btn-error w-full mt-2 text-white" @click="showContact = false">关闭</button>
      </div>
    </div>
  </div>
  <div v-if="previewImg" class="fixed inset-0 bg-black bg-opacity-80 flex items-center justify-center z-50" @click="previewImg = ''">
    <img :src="previewImg" class="max-h-[80vh] max-w-[90vw] rounded-lg shadow-lg border-4 border-white" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
const previewImg = ref('');
const showContact = ref(false);

const aiInfos = [
  { label: '工程师', value: 'Aric', copy: true },
  { label: '电话', value: '18903676153', copy: true },
  { label: '邮箱', value: 'wuxiumu@163.com', copy: true },
  { label: '微信', value: 'qingbao199101', copy: true },
  { label: '技能栈', value: 'Vue,Go,Java,PHP', copy: false },
  { label: '一句话介绍', value: '十年码农，AI应用实践派', copy: false },
  { label: '兴趣爱好', value: '骑行、摄影、看科技八卦', copy: false },
  { label: '个人宣言', value: '代码是桥梁，AI是未来', copy: false },
  { label: '星座', value: '摩羯座♑️', copy: false },
  { label: '座右铭', value: 'Stay hungry, stay foolish.', copy: false }
];

const engineer = {
  payWx: 'https://archive.biliimg.com/bfs/archive/6c98caa50579e4df2ce124c5c296fed953ef0a0a.jpg',
  payAli: 'https://archive.biliimg.com/bfs/archive/0978ae41412bf85354a38a67f4ae9423ff19b1e4.jpg',
  payUrl: 'https://pay.gua.hk/'
};

// 复制功能
function copy(text: string) {
  navigator.clipboard.writeText(text);
}
</script>