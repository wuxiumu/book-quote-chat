<template>
  <el-container>
    <el-aside width="180px">
      <el-menu router>
        <el-menu-item index="/admin/moderate">审核</el-menu-item>
        <el-menu-item index="/admin/user">用户管理</el-menu-item>
        <el-menu-item index="/admin/log">操作日志</el-menu-item>
        <el-menu-item index="/admin/friendlink">友链管理</el-menu-item>
        <el-menu-item index="/admin/stat">多指标数据统计</el-menu-item>
      </el-menu>
    </el-aside>
    <el-main>
      <router-view/>
    </el-main>
  </el-container>
</template>