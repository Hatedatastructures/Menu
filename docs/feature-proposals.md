# CS Board 白板声画工坊 · 功能规划文档

> 基于 codegraph 全量代码阅读整理。**后端主链：Go（`server/internal/service/workshop/`，端口 8888）**，Python FastAPI 版仅作行为参考。目标：先讲清项目现状，再给出以 **Remotion 渲染** 为核心的功能扩展方案与实施路线。

---

## 一、项目现状速览

### 1.1 整体架构

```text
┌─────────────────────────── 浏览器前端 ───────────────────────────┐
│ web/ (React 19 + vinext/Vite, 端口 13000)                        │
│   app/page.tsx        —— 单页工作台（文案/风格/人物/设置/历史）      │
│   src/api/{job,config,health,preference}/  —— API 封装            │
│   src/utils/request.ts —— fetch 封装（{code,data} 信封解析）        │
└──────────────┬────────────────────────────────────────────────────┘
               │ /api/* （vite proxy → 127.0.0.1:8888 = Go 后端）
┌──────────────▼────────────────────────────────────────────────────┐
│ Go 后端（主链）                                                     │
│   internal/service/job/job.go     —— 服务层（@path 注解路由约定）     │
│   internal/service/workshop/                                        │
│     create.go   创建任务 / CreateRerender（队列上限 20）              │
│     voice.go    edge-tts 配音阶段                                    │
│     plan.go / prompts.go / text.go   分镜规划                        │
│     images.go / providers.go    图片模型调用（3 次自动重试）           │
│     jobs.go     内存任务表 + 三级队列 EnsureWorkers(:605)             │
│     render.go   手绘渲染 / 合成（RenderGeneratedJob :236 等）         │
└──────────────┬────────────────────────────────────────────────────┘
               │ runCmd 子进程调用
┌──────────────▼────────────────────────────────────────────────────┐
│ 渲染脚本层 scripts/（被 Go 通过 .venv Python 调用）                   │
│   stream_render.py           手绘模拟渲染核心（骨架追踪笔迹+笔尖贴图） │
│   render_stream_whiteboard.py 入口封装                             │
│   merge_scenes.py            多段 MP4 拼接（ffmpeg concat 优先）     │
│   add_key_text.py            分镜图顶部中文重点词叠加                │
└────────────────────────────────────────────────────────────────────┘
```

> 注：`start-webapp.ps1` 目前仍启动 Python uvicorn(:18765)，而前端代理指向 ：8888（Go）。以 Go 为准后该脚本后续需改为启动 Go 服务（不在本档范围内展开）。

### 1.2 一条任务的完整流水线（Go 实现）

| 阶段 | Go 函数 | 进度 | 产物 |
| --- | --- | --- | --- |
| 配音 | `VoiceStage`（voice.go） | 0→14% | `voice.mp3` |
| 分镜规划 | `ModelStage` 内调 `make_plan` 复刻（plan.go） | 14→22% | `plan.json`（title/key_text/concept/elements/text/duration_ms） |
| 图片生成 | `ModelStage`（images.go） | 22→78% | `board-NN.png` + `.source.png` + `boards.json` |
| 手绘渲染 | `renderBoardVideos`（[render.go:130](../server/internal/service/workshop/render.go#L130)） | 78→90% | `board-NN.mp4` + `.annotation.json` |
| 合成 | `composeFinalVideo`（[render.go:177](../server/internal/service/workshop/render.go#L177)） | 90→100% | `silent.mp4` → `final.mp4`（+字幕烧录） |

关键机制：
- **断点恢复**：每阶段产物落盘后写 `checkpoint`；重启时按 checkpoint 续跑（`RetryJob`）。
- **重渲染不花钱**：`RerenderJob`（[render.go:263](../server/internal/service/workshop/render.go#L263)）只复用 `voice.mp3 + plan.json + board-*.png` 重做本地渲染。
- **任务目录**：所有中间产物都是 `.webapp/jobs/<job_id>/` 下的文件——这是接入任何新渲染引擎的理想数据契约。
- **路由约定**：服务层方法用 `@path` 注解声明路由（如 [job.go:89](../server/internal/service/job/job.go#L89) 的 `/api/jobs/:id/rerender`），经 req/res 模型 + `snapTo` 转换为 `{code,data}` 信封响应。

### 1.3 关键文件索引

| 文件 | 职责 |
| --- | --- |
| [server/internal/service/workshop/render.go](../server/internal/service/workshop/render.go) | 手绘渲染 + 合成（Remotion 主要挂载点） |
| [server/internal/service/workshop/jobs.go](../server/internal/service/workshop/jobs.go) | 任务表、三级队列、EnsureWorkers |
| [server/internal/service/workshop/create.go](../server/internal/service/workshop/create.go) | 创建任务 / 重渲染入口（参数归一化在此） |
| [server/internal/service/job/job.go](../server/internal/service/job/job.go) | HTTP 服务层（新增端点照此模式写） |
| [scripts/stream_render.py](../scripts/stream_render.py) | 手绘动画核心算法（保留不动） |
| [web/app/page.tsx](../web/app/page.tsx) | 前端整页（状态、轮询、表单） |
| [web/src/api/job/index.ts](../web/src/api/job/index.ts) | 前端 Job API 封装 |

---

## 二、核心提案：Remotion 渲染引擎

### 2.1 定位判断（重要）

现有手绘引擎的价值在「真实笔迹模拟」（骨架追踪、笔尖贴图、停顿呼吸），**Remotion 不应替换它，而是补位它做不到的三件事**：

1. **浏览器内实时预览** —— `@remotion/player` 直接在网页里播放成片效果，不用等几分钟渲染完才知道结果；
2. **包装合成层** —— 片头卡、片尾订阅卡、转场、进度条、BGM 波形，这些「动效包装」用手绘引擎做很别扭，用 Remotion 是声明式几行组件的事；
3. **零图片模型的动态图文风格** —— 新增一种「动态信息图」风格：直接读 `plan.json` 用程序化动效（文字弹入、SVG 描边、图标入场）出片，**不调用图片模型，几十秒出一支视频**，作为低成本快产线。

### 2.2 三种接入形态

| 形态 | 用什么 | 输入 | 输出 | 成本 | 建议 |
| --- | --- | --- | --- | --- | --- |
| A. 预览播放器 | `@remotion/player`（前端依赖） | `plan.json` + board 图（经新 API 暴露） | 浏览器内即时播放 | 无渲染成本 | ✅ 第一优先，UX 提升最大 |
| B. 包装合成器 | `remotion render`（服务端 Node） | 各段 `board-NN.mp4` + 配置 | 带片头片尾转场的完整片 | 一次额外渲染 | ✅ 第二步 |
| C. 动态图文风格 | `remotion render` | 仅 `plan.json`（无图片模型） | 全程序化动效视频 | 极低 | ✅ 第三步，新风格选项 |

### 2.3 数据契约（现成的，不用改流水线）

Remotion 的 props 直接由任务目录产物组装：

```jsonc
{
  "scenes": [ /* plan.json 原样：title, key_text, concept, elements, text, duration_ms */ ],
  "boards": [
    { "index": 1, "image": "board-01.png", "video": "board-01.mp4",
      "durationMs": 12400, "sceneNumbers": [1, 2] }
  ],
  "audio": { "src": "voice.mp3", "durationMs": 98000 },
  "subtitles": [ /* WriteSubtitles 已能产出 cue 列表 */ ],
  "brand": { "penText": "账号名", "style": "极简粗线简笔白板风" },
  "layout": { "aspect": "16:9", "width": 1920, "height": 1080 }
}
```

### 2.4 目录结构建议

```text
cs-board/
├── remotion/                    # 新增：独立 npm 包
│   ├── package.json             # remotion, @remotion/player, @remotion/cli
│   ├── remotion.config.ts
│   ├── src/
│   │   ├── Root.tsx             # <Composition id="WhiteboardPack" .../>
│   │   ├── schema.ts            # zod schema = 上面的 props 契约
│   │   ├── compositions/
│   │   │   ├── PackagedVideo.tsx      # 形态B：串 board 视频 + 片头尾 + 转场
│   │   │   ├── MotionInfographic.tsx  # 形态C：纯程序化动效
│   │   │   └── pieces/          # IntroCard / EndCard / ProgressBar / KeyTextPop ...
│   │   └── styles.ts            # 与 11 种画面风格联动的配色 token
│   └── public/                  # 渲染时拷贝/软链 job 产物的临时目录
├── web/
│   ├── app/page.tsx             # 加 <Player> 预览面板 + 渲染器选择器
│   └── src/api/job/index.ts     # 加 previewData(id)、assetUrl(id, name)
└── server/internal/service/workshop/
    └── remotion.go              # 新增：组 props JSON → 调 node remotion CLI
```

### 2.5 Go 服务端集成点（主链）

1. **renderer 参数贯通**：仿照 `include_subtitles` 的既有链路——
   `CreateJob` 表单解析（create.go 归一化处）→ 存入 job map → `task` 结构体加字段（jobs.go）→ `ModelStage` 尾部按 renderer 分流：
   - `"handdrawn"`（默认）：现有路径不变；
   - `"handdrawn-pack"`：`renderBoardVideos` 完成各段 mp4 后追加 `remotionPack()`（片头尾+转场包装）；
   - `"motion"`：跳过手绘与图片模型，规划完成后直接走 `remotionMotion()`。
   
   同步修改：`CreateRerender`（create.go:376）接受并透传 renderer；断点恢复入队处按 renderer 分流。
2. **新增 [remotion.go](../server/internal/service/workshop/)（照 render.go 风格写）**：
   ```go
   // buildPreviewProps(jobID) 组装 2.3 契约 JSON（读 plan.json / boards.json / probeDuration）
   // runRemotionRender(jobID, composition, outPartial string) error：
   //   cmd := []string{"node", remotionCLI, "render", composition, out,
   //                   "--props", propsPath, "--concurrency", ...}
   //   产物命名 final.partial.mp4 → os.Rename → final.mp4，复刻 composeFinalVideo 的 partial 模式
   ```
   并发控制：用带缓冲 channel 信号量限制同时 1–2 个 remotion 进程（对齐现有 RENDER 并发策略）。
3. **预览静态产物端点**（形态A必需，照 [job.go](../server/internal/service/job/job.go) 现有模式加注解方法）：
   - `GET /api/jobs/:id/preview-data` → 返回组装好的 props JSON（data 字段直出）
   - `GET /api/jobs/:id/assets/:name` → 白名单校验文件名（仅 `board-*.png|mp4`、`voice.mp3`）防目录穿越后 `c.File(...)`，参考现有 `DownloadJob` 的文件流返回。
4. **检查点兼容**：remotion 产物走 `final.partial.mp4 → final.mp4` + `validMediaFile` 校验，断点恢复逻辑自动生效，无需新 checkpoint 类型（可另存 `checkpoint: "remotion"` 便于续跑提示）。

### 2.6 前端集成点

- [page.tsx](../web/app/page.tsx) `controlPanel` 区加「渲染引擎」选择（标准手绘 / 带包装 / 动态图文），随 `create` FormData 提交 `renderer` 字段，并存入本地偏好 `LOCAL_PREFERENCES_KEY`。
- 任务完成卡片加「预览」按钮 → 动态 `import('@remotion/player')` 懒加载（打包体积约几百 KB），`inputProps` 来自新的 preview-data 端点，素材 URL 走 assets 端点。
- [web/src/api/job/index.ts](../web/src/api/job/index.ts) 增加 `previewData(id)` 与 `assetUrl(id, name)` 工具函数。

### 2.7 环境注意

- Remotion 源码可得但非纯开源：个人与 ≤3 人公司免费；当前本地自用无碍。
- `remotion render` 需要无头 Chrome（首次自动下载）；本机 Node ≥22.13 满足要求；Windows 支持良好。
- Go 调用时建议锁定 `node` 绝对路径或依赖 PATH，并在 health 端点暴露 `remotion: ok/missing` 便于诊断。

---

## 三、其他功能建议（按优先级）

### P0 — 低成本高回报，建议随 Remotion 一起做

| 功能 | 说明 | 主要改动点 |
| --- | --- | --- |
| 背景音乐混音 | 内置免版权 BGM 库，音量可调；合成时 ffmpeg `amix` 配音+BGM（配音 ducking 可选） | `composeFinalVideo`（render.go）+ 前端选择器 + config 增加 bgm_volume |
| 竖屏/横屏输出 | 9:16 抖音版 / 16:9 B站版一键切换（图片模型按比例出图 + remotion layout 参数） | create 表单字段 → 图片 prompt 附比例 → 渲染布局 |
| 单分镜重新生成 | 历史任务里点某张 `board-NN` 只重跑该图模型调用，其余复用 | 新端点 `POST /api/jobs/:id/regenerate-board`，从 `ModelStage` 循环体抽函数 |
| NVENC 硬件编码 | 有 NVIDIA 卡时 `-c:v h264_nvenc` 替代 libx264，合成速度提升数倍 | `merge_scenes.py` / `composeFinalVideo` 启动时探测一次 `ffmpeg -encoders` 缓存结果 |

### P1 — 体验深化

| 功能 | 说明 | 主要改动点 |
| --- | --- | --- |
| 分镜编辑器 | 出图前可改 `key_text`/文案分组/删减场景，避免整单重来 | 前端新面板 + `PATCH /api/jobs/:id/plan`（写回 plan.json 后从 model 阶段续跑） |
| 多 TTS 供应商 | Azure Speech、GPT-SoVITS（本地）等做成可选节点，health 显示多节点 | voice.go 抽 provider 列表（providers.go 已有雏形） |
| 成片库页面 | 独立 `/library` 路由：网格浏览所有 done 任务、收藏、对比两版渲染 | 前端新路由 + listJobs 扩展过滤参数 |
| 字幕样式可选 | 字号/描边/位置预设 3 套（现在 force_style 写死在 `ffmpegSubtitleFilter`） | 抽配置函数 + config 字段 |
| 批量生产 | 粘贴多条文案（或上传 txt/csv）建一批任务进队列 | 前端批量入口 + 循环调 createJob（队列上限 20 已有保护） |

### P2 — 远期方向

- **数字人口播合成**：手绘片段 + 真人半身像开场（第三方推理服务）。
- **双语字幕/翻译发布**：plan.json 文本送文本模型译英，生成双 SRT。
- **开放 Webhook/API Token**：让外部自动化（剪映草稿、公众号定时发布）拉取成片。

---

## 四、建议实施路线图

```text
阶段一（约 2–4 天）Remotion 地基 + 快赢项
├── 搭 remotion/ 包 + zod schema 契约
├── 形态A：<Player> 预览（Go 新增 preview-data / assets 端点 + 前端预览面板）
├── P0：BGM 混音、NVENC 探测
└── 验收：任意历史任务可在浏览器秒开预览

阶段二（约 1 周）Remotion 包装层上线
├── 形态B：IntroCard / EndCard / 转场 / 进度条组件
├── renderer 参数全链贯通（create → task → ModelStage 分流 → rerender/断点恢复）
├── P0：竖屏比例、单分镜重生
└── 验收：创建任务时可选「带片头片尾包装」，断点恢复正常

阶段三（约 1–2 周）动态图文风格 + 生产化
├── 形态C：MotionInfographic 全程序化风格（11 种风格配色映射）
├── P1：分镜编辑器、成片库
└── 验收：一条纯文案 → 60 秒内出片的快产线可用
```
