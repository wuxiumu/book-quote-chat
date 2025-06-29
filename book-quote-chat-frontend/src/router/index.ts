// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router';
import Home from '@/views/Home.vue';
import Quotes from '@/views/Quotes.vue';
import Chat from '@/views/Chat.vue';
import Friends from '@/views/Friends.vue';
import ChatRoom from '@/views/ChatRoom.vue';
import ChatPrivateRoom from '@/views/ChatPrivateRoom.vue';

const routes = [
    { path: '/', name: 'Home', component: Home },
    {
        path: '/user/login',
        component: () => import('@/views/user/LoginRegister.vue'),
        meta: { title: '用户登录' }
    },
    { path: '/quotes', name: 'Quotes', component: Quotes },
    { path: '/chat', name: 'Chat', component: Chat },
    {
        path: '/chat/private/:id',
        name: 'PrivateChatRoom',
        component: ChatPrivateRoom,
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
    { path: '/friends', name: 'Friends', component: Friends },
    {
        path: '/chat/:id',
        name: 'ChatRoom',
        component: ChatRoom,
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