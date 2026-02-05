# LunaFox AI Nudge Proxy (Deno Version)

一个极其轻量的 API 代理服务，运行在 Deno Deploy 上，用于为 LunaFox 客户端生成 AI 关怀文案，同时隐藏你的 API Key。

## 为什么选择 Deno Deploy?
- **完全免费**: Free Tier 每月 10万次请求。
- **无需绑卡**: 只要有 GitHub 账号就能用。
- **全球加速**: 基于边缘计算，速度快。
- **单文件**: 只需要一个 `main.ts`。

## 🚀 部署指南 (3 分钟)

1. **推送代码**: 确保你的 `tools/nudge-proxy/main.ts` 已提交到 GitHub 仓库。
2. **登录 Deno Deploy**: 访问 [dash.deno.com](https://dash.deno.com)。
3. **新建项目**:
   - 点击 **"New Project"**。
   - 选择 **"Select a repository"** -> 找到你的 `lunafox` 仓库。
   - **Branch**: 选择你的分支 (如 `002-server-orchestration` 或 `main`)。
   - **Entrypoint**: 选择 `tools/nudge-proxy/main.ts`。
   - 点击 **"Link"**。
4. **配置环境变量 (Environment Variables)**:
   - 项目创建后，进入 **Settings** -> **Environment Variables**。
   - 添加 `AI_API_KEY`: 填入你的 DeepSeek / OpenAI API Key。
   - (可选) 添加 `AI_BASE_URL`: 如果用 DeepSeek，填 `https://api.deepseek.com/v1`。
5. **获取域名**:
   - Deno 会自动给你分配一个 `xxxx.deno.dev` 的域名。
   - 复制这个域名。

## 🔗 前端集成

在 LunaFox 前端 (`frontend`) 的环境变量 (`.env`) 中设置：

```env
NEXT_PUBLIC_NUDGE_API_URL=https://your-project-name.deno.dev
```

注意：Deno Deploy 默认不需要 `/generate` 路径后缀（取决于 `main.ts` 怎么写，目前的 `main.ts` 会直接处理根路径 POST 请求，或者你可以加上 `/generate` 但代码里要适配路由）。

**注意**: 当前 `main.ts` 简单地处理所有 POST 请求。你可以直接填 `https://your-project-name.deno.dev` 作为 API 地址。
