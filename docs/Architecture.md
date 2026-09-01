# Menu 架构基线

状态：Phase 0-5 已实现；当前 Windows 服务端、Qt 桌面、Android x86_64 AVD 和浏览器管理台均已完成 fresh 验证，证据与已知限制已归档。

日期：2026-09-01

## 产品边界

Menu 是一个家庭做饭执行器。核心闭环是：用户在 30 秒内从今晚的可行方案中做决定，看到准确的食材和缺口，在合适的时间完成备菜，然后按步骤把一顿饭做完。

首版包含：

- 首次设置：人数、菜系偏好、忌口/过敏、可用时间、厨具和常备食材。
- “今晚”：三个确定性推荐方案，展示总时长、份量、难度、食材数量和已有/缺少。
- 菜谱详情：本地成品图、份量换算、规范化食材、处理方式、替代食材和分步耗时。
- 计划：今天和一周计划、重复食材合并、采购/备菜清单。
- 做饭模式：分步完成、并行计时器、解冻/腌制提醒、退出恢复和屏幕常亮策略。
- 反馈、可调提醒、离线读取已打开菜谱/今日计划/做饭步骤。
- 桌面管理面板：菜谱、步骤、食材、分类、发布状态和基础统计。

首版明确不做：朴朴或其他商超购物车、社交社区、聊天、短视频流、版权内容抓取、医疗营养诊断和无约束 AI 菜谱生成。Ingredient 保留可为空的 StoreSkuMapping 字段，为后续商超适配预留稳定边界。

## 目标运行形态

```text
Android Qt Quick / Desktop Qt Quick
              |
          HTTPS/JSON
              |
       MenuApi + MenuTransport
              |
       MenuApplication (use cases)
          /             \
 MenuDomain       MenuInfrastructure
                         |
                  SQLite + WAL + Media

Admin React/Vite ------ HTTPS/JSON ------ MenuApi
```

服务端是 C++20。客户端的 UI 用 Qt Quick，桌面版先作为可运行验证端，Android 版使用相同的 `MenuClientCore` 业务边界。管理面板是 TypeScript + React/Vite 静态包，直接调用同一套 API，不在前端复制领域规则。

## C++ 模块边界

目标名和源文件名统一使用 PascalCase；协议字段、JSON key、SQL 列名和第三方宏遵循外部约定。

下面箭头表示“右侧 target 依赖左侧 target”，而不是继承关系：

```text
MenuDomain          -> MenuFoundation
MenuApplication     -> MenuDomain, MenuFoundation
MenuInfrastructure  -> MenuApplication, MenuDomain, MenuFoundation
MenuTransportCore   -> MenuFoundation
MenuTransport       -> MenuTransportCore, MenuFoundation
MenuApi             -> MenuApplication, MenuTransportCore, MenuFoundation
MenuServer          -> MenuApi, MenuTransport, MenuInfrastructure, MenuApplication
```

- `MenuFoundation`：配置、Result/Error、时钟、输入尺寸限制、日志脱敏和公共 ID。
- `MenuDomain`：User、Preference、Ingredient、Recipe、RecipeIngredient、RecipeStep、MealPlan、MealPlanItem、Reminder、CookingSession、Feedback、MediaAsset，以及份量换算、别名归一化、替代规则和提醒调度等无 I/O 规则。
- `MenuApplication`：用例服务、事务边界、仓储/时钟/推荐端口和 `RecommendationProvider`。确定性规则提供默认实现，模型实现只能通过端口接入并经过 JSON/过敏原/数量/单位/步骤/时间校验。
- `MenuInfrastructure`：SQLite 连接、迁移、种子数据、参数化 SQL 仓储、WAL checkpoint、密码哈希、token 存储和本地媒体。
- `MenuTransportCore`：只定义中立的 `HttpRequest`、`HttpResponse`、`HttpHandler`/router 注入接口，不包含菜谱业务，不依赖具体 listener。
- `MenuTransport`：用 Boost.Asio/Beast 实现 HTTP listener、连接生命周期和可选 WebSocket；只依赖 `MenuTransportCore`，不链接 `MenuApi`。
- `MenuApi`：实现路由和用例适配，消费 `MenuTransportCore` 的中立 HTTP 类型，负责 DTO、Boost.JSON 序列化、认证授权、CORS、统一错误响应和输入校验；不创建或持有具体 listener。
- `MenuServer`：唯一 composition root，负责配置读取、依赖组装，把 `MenuApi` 的 handler 注入 `MenuTransport` listener，然后监听和优雅退出。

配置阶段递归检查 target link closure：`MenuTransportCore` 和 `MenuTransport` 的 closure 不得出现 `MenuApi`，`MenuApi` 的 closure 只能出现 `MenuTransportCore` 而不能出现具体 `MenuTransport`，`MenuServer` 必须同时连接两者。该检查用于防止 `MenuApi <-> MenuTransport` 循环依赖，而不是依赖人工阅读 CMake。

生产目标不依赖测试支持；测试目标按 Production、Preview/Client contract 和 API integration 分开，避免测试辅助库反向污染生产依赖。

## 并发与 I/O 语义

- 一个 `boost::asio::io_context` 由受控数量的工作线程驱动；连接状态和写队列绑定到 per-connection strand。
- HTTP handler 只做边界解析和调度。SQLite、文件媒体读取、密码哈希等可能阻塞的操作投递到明确的 storage/worker executor，完成后回到请求 executor。
- SQLite 写入由单一写 actor/strand 串行化；读操作可使用受控连接池。所有写入使用参数化 SQL 和事务。
- UI 的网络、磁盘和图片解码不能同步阻塞 GUI 线程。Qt 模型批量更新，ListView delegate 尺寸稳定且不在 delegate 中保存业务状态。
- 客户端离线 Recipe 快照必须保存每个食材的 `IngredientId`、`IngredientName`、`IngredientCategory`、`IngredientDefaultUnit`、`IngredientIsPantryStaple`、数量、单位、处理方式和必需标记；离线渲染不能再依赖一次额外的食材请求。`GET /api/v1/ingredients` 使用同一套显示字段供食材页和缓存预热使用。
- WebSocket 只在计划/提醒同步有真实需求时启用；首版 HTTP API 和离线缓存优先，避免为展示而增加协议复杂度。

## 数据与安全基线

SQLite 文件位于运行目录下的 `server/data`，默认启用 WAL、busy timeout 和受控 checkpoint。WAL 仅用于同一主机上的数据库进程；不把数据库文件放网络文件系统。

认证采用短期 access token 加可撤销 refresh token。密码采用 OpenSSL 提供的 scrypt 参数化派生（等价强密码哈希），不把 secret 或 API key 放客户端。管理 API 默认需要管理员身份，登录和 refresh 有限速，输入有大小/字段校验，日志不记录密码、token 和完整个人偏好。

首版 API 契约以 `shared/schema` 的 OpenAPI/JSON Schema 为准。每个 DTO 保留版本号或兼容演进策略，数据库用版本化迁移；基础设施可以替换为 PostgreSQL adapter 而不修改 Domain/Application 接口。

## 已验证的用户闭环

当前实现已经把以下路径接到真实服务端，而不是静态数组或内部函数调用：

1. `MenuServer` 启动迁移、WAL 和种子数据，健康检查、菜谱、食材和推荐端点可用。
2. Qt 客户端通过 `ClientApi` 请求今晚推荐、食材和计划；Recipe DTO 携带中文食材名称、类别、单位、常备状态和处理方式。
3. 管理面板使用真实登录和 Bearer 认证完成菜谱/食材 CRUD、草稿发布、预览、删除确认和错误态。
4. 客户端支持一周计划的添加/移除和合并清单、做饭模式步骤勾选/并行计时器/反馈、偏好设置和用户隔离缓存；登录后的 access token 可通过 `RefreshSession` 使用 refresh token 轮换。
5. 网络失败时使用带 schema/version 2 的本地缓存；401/登出会清除私有模型、活动菜谱、会话和对应缓存，并阻止旧异步写入复活。

这些功能分别由单元、API 集成、Qt/QML、真实 HTTP 和 Playwright 测试覆盖；未测量的跨平台/高刷新率项目仍在文档末尾明确列出。

## 文件职责拆分

- `server/Api/ApiParsing*.cpp` 按 query、通用 JSON、认证、菜谱/食材、推荐和工作流 payload 拆分；`ApiRouter.cpp` 只负责组合、CORS 和请求分派。
- `client/cpp/ClientApi*.cpp` 按状态、缓存、请求、认证和工作流拆分；模型只做稳定 role 映射和批量转换。
- `admin/web/src` 按登录入口、工作区、页面、编辑器、预览和基础控件拆分；`App.tsx` 不承载 CRUD 业务。

拆分不是新的运行时边界：所有请求仍通过 `ClientApi`/`MenuApi`，便于在 fresh 构建中检查依赖和行为。

## 环境决策

- 本机复用已有 `C:\Qt\6.8.3` 桌面 kit，不在 C 盘安装第二套 Qt。
- C++ 依赖使用带 SHA256 的 CMake FetchContent 锁定；源码、构建和下载缓存通过 `I:\code\Menu\.cache`、`I:\code\Menu\build` 和 `I:\code\Menu\.tools` 管理。若后续 vcpkg 能在 I 盘可复现安装，再评估迁移，不改变 target 边界。
- npm cache 为 `I:\code\Menu\.npm-cache`，Gradle 为 `I:\code\Menu\.gradle`，Android SDK/NDK/JDK 只允许放 `I:\code\Menu\.tools` 或 `I:\android`。
- 代理只在下载当前进程中按需设置 `HTTP_PROXY`/`HTTPS_PROXY=http://127.0.0.1:7890`，不写入系统或用户全局配置。
- Qt Android 目标使用 `I:\code\Menu\.tools\qt-android\6.8.3\android_x86_64` 的源码构建 `android-clang` prefix；跨编译 prefix 按 Qt 设计不包含 host `androiddeployqt`，APK 由已验证的 `C:\Qt\6.8.3\mingw_64\bin\androiddeployqt.exe` 作为 `QT_HOST_PATH` 工具生成。
- 当前证据来自 Windows 桌面、Windows 服务端和 Android x86_64 AVD；Linux、ARM Android 真机和 iOS 不能由这些结果推断通过。

## 当前风险

| 风险 | 当前状态 | 处理方式 |
| --- | --- | --- |
| Android Qt kit、JDK、SDK、NDK、Gradle 不在 PATH | 已满足（临时进程环境） | 所有新增工具在 I 盘；Qt target prefix 和 host deployment tool 分开记录，最终命令显式设置环境 |
| CMake 是 4.3.0-rc1 | 可用但非稳定版 | 优先验证；若出现工具链兼容问题，在 I 盘准备稳定版并记录替换原因 |
| Windows 没有 Ninja/MSVC | MinGW/LLVM-MinGW 可用 | 首版使用已有 MinGW 或 Qt kit 的编译器；CI 另行提供 Linux/Windows 证据 |
| 食谱图片授权与离线包大小 | 已建立 | 使用仓库内生成位图，记录尺寸/hash/来源；不依赖外链 |
| SQLite 写 actor 与异步请求生命周期 | 已实现并测试 | SQLite 工作投递到 storage executor，写事务串行化；异常和 completion dispatch failure 有回归测试 |
| 120Hz 证据 | 未满足 | AVD 报告 60Hz；保留首屏 gfxinfo，不能宣称真实设备 120fps |
| Linux/ARM/iOS 构建 | 未验证 | 需要对应工具链和设备；当前不把 Windows/Android x86_64 结果外推 |
