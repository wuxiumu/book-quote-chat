// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router';

// 路由配置
const routes = [
    { path: '/', name: 'Home', component: () => import('@/views/Home.vue') },
    {
        path: '/user/login',
        component: () => import('@/views/user/LoginRegister.vue'),
        meta: { title: '用户登录' }
    },
    { path: '/quotes', name: 'Quotes', component: () => import('@/views/Quotes.vue') },
    { path: '/chat', name: 'Chat', component: () => import('@/views/Chat.vue') },
    {
        path: '/chat/private/:id',
        name: 'PrivateChatRoom',
        component: () => import('@/views/ChatPrivateRoom.vue'),
    },
    {
        path: '/agora-chat',
        name: 'AgoraChat',
        component: () => import('@/views/AgoraChat.vue')
    },
    {
        path: '/agora-chats',
        name: 'AgoraMultiChat',
        component: () => import('@/views/AgoraMultiChat.vue')
    },
    { path: '/friends', name: 'Friends', component: () => import('@/views/Friends.vue') },
    {
        path: '/chat/:id',
        name: 'ChatRoom',
        component: () => import('@/views/ChatRoom.vue'),
    },
    {
        path: '/admin',
        component: () => import('@/views/admin/AdminLayout.vue'),
        children: [
            { path: '', redirect: '/admin/moderate' },
            { path: 'moderate', component: () => import('@/views/admin/Moderate.vue'), name: 'Moderate' },
            { path: 'user', component: () => import('@/views/admin/UserManage.vue'), name: 'UserManage' },
            { path: 'log', component: () => import('@/views/admin/LogList.vue'), name: 'LogList' },
            { path: 'friendlink', component: () => import('@/views/admin/FriendLink.vue'), name: 'FriendLink' },
            { path: 'stat', component: () => import('@/views/admin/Stat.vue'),meta: { title: '数据统计' }, name: 'Stat' }
        ]
    },
    {
        path: '/admin/moderate-comment',
        component: () => import('@/views/admin/ModerateComment.vue'),
        name: 'ModerateComment',
        meta: { title: '评论审核' }
    },
    {
        path: '/admin/moderate-quote',
        component: () => import('@/views/admin/ModerateQuote.vue'),
        name: 'ModerateQuote',
        meta: { title: '金句审核' }
    }

];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;