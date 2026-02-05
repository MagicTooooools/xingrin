# LunaFox Nudge Proxy

一个轻量级的 API 中转服务，用于为 LunaFox 用户生成 AI 驱动的关怀弹窗（Nudges），同时避免在开源前端代码中暴露 API Key。

## 🚀 部署到 Zeabur

1. **推送到 GitHub**: 确保 `tools/nudge-proxy` 目录已包含在你的 GitHub 仓库中。
2. **登录 Zeabur**: 访问 [Zeabur 控制台](https://zeabur.com)。
3. **创建项目**: 创建一个新项目（例如 `lunafox-services`）。
4. **部署服务**:
   - 点击 **新建服务** -> **Git**。
   - 选择你的仓库。
   - **根目录 (Root Directory)**: 设置为 `tools/nudge-proxy`。
   - Zeabur 会自动识别为 Node.js 项目并部署。
5. **配置环境变量**:
   - 进入服务的 **设置** -> **变量**。
   - 添加 `AI_API_KEY`: 你的 OpenAI 或 DeepSeek API Key。
   - (可选) `AI_BASE_URL`: 例如 `https://api.deepseek.com/v1`。
6. **启用域名**:
   - 进入 **网络 (Networking)**。
   - 启用 "公网域名" (你会获得一个类似 `xxx.zeabur.app` 的地址)。

## 🔗 前端集成

在 LunaFox 前端 (`frontend`) 的环境变量中设置：
```env
NEXT_PUBLIC_NUDGE_API_URL=https://your-service.zeabur.app/generate
```

`useAiNudge` Hook 会检测此变量，如果存在则自动调用该接口，否则降级使用本地静态文案。
