<template>
  <div>
    <el-page-header content="数据统计大盘"/>
    <el-row :gutter="20" class="mb-4">
      <el-col :span="6"><el-card>用户总数<br><b style="font-size:2em">{{ stat.userCount }}</b></el-card></el-col>
      <el-col :span="6"><el-card>金句总数<br><b style="font-size:2em">{{ stat.quoteCount }}</b></el-card></el-col>
      <el-col :span="6"><el-card>评论总数<br><b style="font-size:2em">{{ stat.commentCount }}</b></el-card></el-col>
      <el-col :span="6" class="flex items-center justify-end">
        <el-button @click="exportData" type="primary">导出当前数据</el-button>
      </el-col>
    </el-row>
    <el-row :gutter="16" class="mb-2">
      <el-col :span="8">
        <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            size="small"
            @change="fetchStat"
        />
      </el-col>
      <el-col :span="16">
        <el-radio-group v-model="chartType" size="small" @change="drawAllCharts">
          <el-radio-button label="line">折线</el-radio-button>
          <el-radio-button label="bar">柱状</el-radio-button>
          <el-radio-button label="area">面积</el-radio-button>
          <el-radio-button label="pie">饼图</el-radio-button>
        </el-radio-group>
      </el-col>
    </el-row>
    <el-row :gutter="24">
      <el-col :span="12"><el-card><div ref="chartUser" style="height:300px"/></el-card></el-col>
      <el-col :span="12"><el-card><div ref="chartActive" style="height:300px"/></el-card></el-col>
    </el-row>
    <el-row :gutter="24" class="mt-4">
      <el-col :span="12"><el-card><div ref="chartQuote" style="height:300px"/></el-card></el-col>
      <el-col :span="12"><el-card><div ref="chartComment" style="height:300px"/></el-card></el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import request from '@/api/request'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'

const stat = ref({
  userCount: 0,
  quoteCount: 0,
  commentCount: 0,
  userTrend: [],
  commentTrend: [],
  quoteTrend: [],
  activeUserTrend: []
})

const chartUser = ref<HTMLDivElement|null>(null)
const chartActive = ref<HTMLDivElement|null>(null)
const chartQuote = ref<HTMLDivElement|null>(null)
const chartComment = ref<HTMLDivElement|null>(null)

const chartType = ref('line')
const dateRange = ref<[Date, Date] | null>(null)

// 默认近7天
onMounted(() => {
  const end = dayjs().endOf('day').toDate()
  const start = dayjs().subtract(6, 'day').startOf('day').toDate()
  dateRange.value = [start, end]
  fetchStat()
})

async function fetchStat() {
  try {
    const [start, end] = dateRange.value || []
    const params: any = {}
    if (start && end) {
      params.start = dayjs(start).format('YYYY-MM-DD')
      params.end = dayjs(end).format('YYYY-MM-DD')
    }
    const res = await request.get('/api/admin/stat_overview', { params })
    stat.value = res.data
    nextTick(drawAllCharts)
  } catch {
    ElMessage.error('数据获取失败')
  }
}

function drawAllCharts() {
  drawChart(chartUser.value, stat.value.userTrend, '用户注册数', '用户注册趋势')
  drawChart(chartActive.value, stat.value.activeUserTrend, '日活跃用户', '活跃用户趋势')
  drawChart(chartQuote.value, stat.value.quoteTrend, '新增金句', '金句新增趋势')
  drawChart(chartComment.value, stat.value.commentTrend, '评论数', '评论趋势')
}

function drawChart(dom: HTMLDivElement|null, trend: any[], seriesName: string, title: string) {
  if (!dom) return
  const dates = (trend || []).map((p:any) => p.date)
  const data = (trend || []).map((p:any) => p.count)
  const chart = echarts.init(dom)
  if (chartType.value === 'pie') {
    chart.setOption({
      title: { text: title + '（饼图）', left: 'center', top: 8 },
      tooltip: { trigger: 'item' },
      series: [{
        name: seriesName,
        type: 'pie',
        radius: '65%',
        data: dates.map((d, i) => ({ name: d, value: data[i] })),
        label: { show: true, formatter: '{b}: {c}' }
      }]
    })
    return
  }
  chart.setOption({
    title: { text: title, left: 'center', top: 8 },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: dates },
    yAxis: { type: 'value' },
    series: [{
      name: seriesName,
      data,
      type: chartType.value === 'area' ? 'line' : chartType.value,
      smooth: true,
      areaStyle: chartType.value === 'area' ? {} : undefined,
    }]
  })
}

function exportData() {
  const data = {
    ...stat.value,
    dateRange: dateRange.value
        ? [
          dayjs(dateRange.value[0]).format('YYYY-MM-DD'),
          dayjs(dateRange.value[1]).format('YYYY-MM-DD')
        ] : []
  }
  const str = JSON.stringify(data, null, 2)
  const blob = new Blob([str], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'stat.json'
  a.click()
  URL.revokeObjectURL(url)
}
</script>