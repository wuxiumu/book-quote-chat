<template>
  <div class="p-6">
    <!-- 统计区 -->
    <el-row :gutter="24" class="mb-6">
      <el-col :span="6">
        <el-card shadow="hover" class="!bg-gradient-to-r !from-blue-50 !to-blue-100">
          <div class="text-base text-blue-800">今日新增金句</div>
          <div class="font-extrabold text-3xl my-1 text-blue-600">{{ stats.quotesToday }}</div>
          <div class="text-xs text-gray-400">累计：{{ stats.quotesTotal }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="!bg-gradient-to-r !from-orange-50 !to-orange-100">
          <div class="text-base text-orange-800">待审核金句</div>
          <div class="font-extrabold text-3xl my-1 text-orange-600">{{ stats.quotesPending }}</div>
          <div class="text-xs text-gray-400">昨日通过：{{ stats.quotesApprovedYesterday }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="!bg-gradient-to-r !from-green-50 !to-green-100">
          <div class="text-base text-green-800">今日新增评论</div>
          <div class="font-extrabold text-3xl my-1 text-green-600">{{ stats.commentsToday }}</div>
          <div class="text-xs text-gray-400">累计：{{ stats.commentsTotal }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="!bg-gradient-to-r !from-pink-50 !to-pink-100">
          <div class="text-base text-pink-800">待审核评论</div>
          <div class="font-extrabold text-3xl my-1 text-pink-600">{{ stats.commentsPending }}</div>
          <div class="text-xs text-gray-400">昨日通过：{{ stats.commentsApprovedYesterday }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 趋势图表区 -->
    <el-card shadow="never" class="mb-8">
      <template #header>近七日数据趋势</template>
      <div style="height:320px">
        <v-chart :option="chartOption" autoresize />
      </div>
    </el-card>

    <!-- 跳转按钮区 -->
    <div class="flex gap-8 justify-center mt-8">
      <el-button type="primary" size="large" class="!h-16 !text-xl !px-12 font-bold shadow-xl"
                 @click="goToQuote">
        <i class="el-icon-document-checked mr-2"></i>金句审核
      </el-button>
      <el-button type="success" size="large" class="!h-16 !text-xl !px-12 font-bold shadow-xl"
                 @click="goToComment">
        <i class="el-icon-message mr-2"></i>评论审核
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import request from '@/api/request'
// ECharts 支持（需提前安装 echarts + v-charts 或 vue-echarts）
import VChart from 'vue-echarts'

const router = useRouter()
const stats = ref({
  quotesToday: 0,
  quotesPending: 0,
  quotesTotal: 0,
  quotesApprovedYesterday: 0,
  commentsToday: 0,
  commentsPending: 0,
  commentsTotal: 0,
  commentsApprovedYesterday: 0,
  chartDates: [],
  chartQuotes: [],
  chartComments: []
})

function goToQuote() {
  router.push('/admin/moderate-quote')
}
function goToComment() {
  router.push('/admin/moderate-comment')
}

function fetchStats() {
  request.get('/api/admin/stat_audit_overview')
      .then(res => {
        stats.value = res.data
        chartOption.value = {
          tooltip: { trigger: 'axis' },
          legend: { data: ['金句', '评论'] },
          xAxis: { type: 'category', data: stats.value.chartDates },
          yAxis: { type: 'value' },
          series: [
            {
              name: '金句',
              type: 'line',
              smooth: true,
              data: stats.value.chartQuotes
            },
            {
              name: '评论',
              type: 'line',
              smooth: true,
              data: stats.value.chartComments
            }
          ]
        }
      })
      .catch(err => {
        // 避免异常导致整个 admin 区域挂掉
        console.error('统计接口出错', err)
        // 可选：弹窗提示
        // ElMessage.error('无法加载统计数据')
      })
}

const chartOption = ref({})
onMounted(fetchStats)
</script>