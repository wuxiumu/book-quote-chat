# book-quote-chat · 书摘金句聊天室

> 围绕书摘/金句的 **Vue 聊天与收藏**小产品，适合读书社群。

**定制：** 微信 **qingbao199101**（备注：书摘聊天）｜ [wuxiumu@163.com](mailto:wuxiumu@163.com)

> 让灵魂碰撞，让金句流传。这里，每一条神评，都是你我的共鸣！



------



## **项目简介**



**QuoteChat** 是一个极简但酷炫的**读书金句分享+实时群聊交友**平台。

无论你是孤独的读书人，还是八卦的社交动物，这里都有你的“一席之地”。



**特色：**



- ✨ 金句/神评一键点赞、评论、收藏、转发
- 🧑‍🤝‍🧑 好友系统，群聊与私聊，边交友边吐槽
- 🚀 实时 WebSocket 群聊 + 私聊，体验如微信一样丝滑
- 🖼️ 图文/语音/视频消息、表情包，畅聊不卡壳
- 💬 管理后台，金句/评论/友链/用户/数据统计，一站掌控
- 🔒 JWT 登录、权限验证，支持多种登录模式
- 🔥 支持 WebRTC 音视频通话（可集成云厂商接口）



------



## **快速上手**



### **1. 环境准备**



- Node.js >= 18
- Go >= 1.18
- npm >= 10.8
- 推荐操作系统：macOS/Windows/Linux 皆可



### **2. 前端运行**



```
cd book-quote-chat-frontend
npm install
npm run dev         # 本地开发
npm run build       # 生产打包
```



### **3. 后端运行**



```
cd book-quote-chat-backend
go run main.go      # 本地后端服务
```



### **4. .env 配置参考**



```
VITE_API_BASE_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws
ALLOWED_ORIGINS=http://localhost:5173,http://yourdomain.com
UPLOAD_EXTS=jpg,jpeg,png,gif,mp4,mp3
UPLOAD_MAX_SIZE_MB=20
```



------



## **主要功能一览**



- **读书金句/神评**：浏览、点赞、评论、复制，AI 神评论
- **群聊/私聊**：实时 WebSocket，消息动画，表情包、图片、语音/视频、文件
- **好友系统**：加好友、分组、备注、在线状态
- **WebRTC 音视频通话**：1对1通话，支持全屏、挂断、状态显示
- **管理后台**：内容审核、用户管理、日志、数据统计、友链配置
- **权限系统**：注册/登录/JWT，接口鉴权
- **限流防刷**：点赞、评论、举报等敏感操作限流
- **极简美观 UI**：自适应移动端，酷炫动画、支持主题切换
- **多语言/多端兼容**：PC/Mac/手机均流畅



------



## **技术栈**



- 前端：Vue3 + Vite + TailwindCSS + Pinia + WebSocket
- 后端：Go + 原生 HTTP + JWT + 文件存储
- 特色：无需数据库，所有核心数据 JSON 文件化，支持热加载



------



## **项目亮点**



- 前后端完全分离，接口文档完善，MVP 极易二次开发
- 所有关键接口都有自动化/集成测试，稳定可扩展
- UI/UX 追求极致，交互动画媲美大厂社交产品
- 支持扩展外部 AI/云音视频 SDK，兼容未来商业化



------



## **未来计划**



- 移动端 App（Flutter/H5混合）
- AI 智能推荐金句/好友
- 付费音视频分钟充值/会员体系
- 更多主题皮肤&自定义表情包
- 支持多语言国际化
- 支持插件市场，二次开发更简单



------



## **致谢&联系方式**





> 有问题请提 issue，欢迎 PR！

>

> 工程师：**Aric**

> 邮箱： wuxiumu@163.com

> 微信：qingbao199101

>

> 一句话宣言：“代码是桥梁，AI是未来。”



------



**点亮 Star，世界因你不同！🌟**

## SEO 字段预留（产品页 HTML 统一生成时用）

| 字段 | 用途 |
|---|---|
| `seo_title` | `<title>` / OG title |
| `seo_description` | meta description |
| `canonical_url` | 规范 URL |
| `og_image` | 分享图 |
| `robots` | 如 `index,follow` |
| `json_ld_type` | SoftwareApplication / WebPage / FAQPage |
| `hreflang` | zh-CN / en |
| `site_name` | 白标品牌名 |

> **GEO ≠ 地理定位。** GEO = Generative Engine Optimization（生成式引擎优化），见下一节。

## GEO（生成式引擎优化）

面向 ChatGPT / Perplexity / Gemini 等生成式答案引擎的可引用结构：中文为主、双语关键词、事实摘要前置、FAQ 收尾。

### 关键词（中英）

- 中文：书摘, 金句聊天室, 读书社群, Vue聊天
- English: book quotes, quote chat room, reading community

### 事实摘要

- 产品：书摘金句向的聊天/互动应用。
- 适合：读书会、知识社群、内容号副产品。
- 交付：皮肤、房间玩法、部署。
- 不做：不是完整社交 App 替代。

### FAQ

**Q：这是什么？**  
A：见事实摘要「产品」。

**Q：如何定制？**  
A：微信 `qingbao199101`，备注项目名。

**Q：是否含密钥？**  
A：默认不含，自行配置。

## 定制合作

房间玩法、品牌皮肤、私有化。加微信备注「书摘聊天」。

---

## 联系方式

定制、演示或合作：

- 邮箱：[wuxiumu@163.com](mailto:wuxiumu@163.com)
- 微信：**qingbao199101**（添加请备注来意）

### 赞助支持

如果内容对你有帮助，欢迎自愿赞助。

| 微信赞赏码 | 支付宝收款码 |
|---|---|
| ![微信赞赏码](https://chiguashentan-test.oss-cn-beijing.aliyuncs.com/wxprogram/zan/20261002185236_31_227.png) | ![支付宝收款码](https://chiguashentan-test.oss-cn-beijing.aliyuncs.com/wxprogram/zan/20261002185309_32_227.jpg) |

