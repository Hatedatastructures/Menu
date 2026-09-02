# Menu 架构基线

状态：Phase 0-5 已实现；当前 Windows 服务端、Qt 桌面、Android x86_64 AVD、ARM64 APK 和浏览器管理台均已完成 fresh 构建验证，证据与已知限制已归档。客户端连接、通知、主题和服务端启动配置已进一步模块化。

日期：2026-09-02

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

服务端是 C++20。客户端的 UI 用 Qt Quick，桌面和 Android 使用同一个 `MenuClientApp` 应用 target 与 `MenuClientCore` 业务边界；`MenuClientPlatform` 承载屏幕刷新率、主题持久化和 Android 通知等平台能力。管理面板是 TypeScript + React/Vite 静态包，直接调用同一套 API，不在前端复制领域规则。

## C++ 模块边界

目标名和源文件名统一使用 PascalCase；协议字段、JSON key、SQL 列名和第三方宏遵循外部约定。
项目自有客户端目录也使用 PascalCase：`client/Cpp`、`client/Qml`、`client/Tests`、`client/QmlTests`、`client/Android`。Android Java 包路径、资源文件名、Qt 生成的 `android-build` 目录和 SQL 迁移编号属于平台/工具约定，保留其格式。

下面箭头表示“右侧 target 依赖左侧 target”，而不是继承关系：

```text
MenuDomain          -> MenuFoundation
MenuApplication     -> MenuDomain, MenuFoundation
MenuInfrastructure  -> MenuApplication, MenuDomain, MenuFoundation
MenuTransportCore   -> MenuFoundation
MenuTransport       -> MenuTransportCore, MenuFoundation
MenuApi             -> MenuApplication, MenuTransportCore, MenuFoundation
MenuRuntime         -> MenuTransportCore, MenuFoundation
MenuComposition     -> MenuApi, MenuTransport, MenuInfrastructure, MenuApplication, MenuRuntime
MenuServer          -> MenuComposition
```

- `MenuFoundation`：配置、Result/Error、时钟、输入尺寸限制、日志脱敏和公共 ID。
- `MenuDomain`：User、Preference、Ingredient、Recipe、RecipeIngredient、RecipeStep、MealPlan、MealPlanItem、Reminder、CookingSession、Feedback、MediaAsset，以及份量换算、别名归一化、替代规则和提醒调度等无 I/O 规则。
- `MenuApplication`：用例服务、事务边界、仓储/时钟/推荐端口和 `RecommendationProvider`。确定性规则提供默认实现，模型实现只能通过端口接入并经过 JSON/过敏原/数量/单位/步骤/时间校验。
- `MenuInfrastructure`：SQLite 连接、迁移、种子数据、参数化 SQL 仓储、WAL checkpoint、密码哈希、token 存储和本地媒体。
- `MenuTransportCore`：只定义中立的 `HttpRequest`、`HttpResponse`、`HttpHandler`/router 注入接口，不包含菜谱业务，不依赖具体 listener。
- `MenuTransport`：用 Boost.Asio/Beast 实现 HTTP listener、连接生命周期和可选 WebSocket；只依赖 `MenuTransportCore`，不链接 `MenuApi`。
- `MenuApi`：实现路由和用例适配，消费 `MenuTransportCore` 的中立 HTTP 类型，负责 DTO、Boost.JSON 序列化、认证授权、CORS、统一错误响应和输入校验；不创建或持有具体 listener。
- `MenuRuntime`：读取并校验服务端 JSON 配置，只负责启动所需的监听参数、路径和限制，不依赖业务服务或具体存储实现。
- `MenuComposition`：可复用的 composition root，负责目录、数据库迁移/种子、依赖组装，把 `MenuApi` 的 handler 注入 `MenuTransport` listener。
- `MenuServer`：只负责命令行参数、配置读取、启动/运行/停止 `MenuComposition`。

配置阶段递归检查 target link closure：`MenuTransportCore` 和 `MenuTransport` 的 closure 不得出现 `MenuApi`，`MenuApi` 的 closure 只能出现 `MenuTransportCore` 而不能出现具体 `MenuTransport`，`MenuRuntime` 不得依赖业务/API/存储 target，`MenuComposition` 才能连接 API、Transport、Infrastructure、Application 和 Runtime，`MenuServer` 只连接 `MenuComposition`。Runtime 源码的 include 方向也在配置期检查。该检查用于防止跨层循环依赖，而不是依赖人工阅读 CMake。

生产目标不依赖测试支持；测试目标按 Production、Preview/Client contract 和 API integration 分开，避免测试辅助库反向污染生产依赖。

## 并发与 I/O 语义

- 一个 `boost::asio::io_context` 由受控数量的工作线程驱动；连接状态和写队列绑定到 per-connection strand。
- HTTP handler 只做边界解析和调度。SQLite、文件媒体读取、密码哈希等可能阻塞的操作投递到明确的 storage/worker executor，完成后回到请求 executor。
- 当前所有仓储共享一个 SQLite 连接，因此 `MenuServer` 的 storage executor 固定为单 worker，确保跨多条语句的事务不会在线程之间交错；引入读连接池前不得提高该数量。所有写入使用参数化 SQL 和事务。
- UI 的网络、磁盘和图片解码不能同步阻塞 GUI 线程。Qt 模型批量更新，ListView delegate 尺寸稳定且不在 delegate 中保存业务状态。
- 客户端离线 Recipe 快照必须保存每个食材的 `IngredientId`、`IngredientName`、`IngredientCategory`、`IngredientDefaultUnit`、`IngredientIsPantryStaple`、数量、单位、处理方式和必需标记；离线渲染不能再依赖一次额外的食材请求。`GET /api/v1/ingredients` 使用同一套显示字段供食材页和缓存预热使用。
- WebSocket 只在计划/提醒同步有真实需求时启用；首版 HTTP API 和离线缓存优先，避免为展示而增加协议复杂度。

## 数据与安全基线

SQLite 文件位于运行目录下的 `server/Data`，默认启用 WAL、busy timeout 和受控 checkpoint。WAL 仅用于同一主机上的数据库进程；不把数据库文件放网络文件系统。

认证采用短期 access token 加可撤销 refresh token。密码采用 OpenSSL 提供的 scrypt 参数化派生（等价强密码哈希），不把 secret 或 API key 放客户端。管理 API 默认需要管理员身份，登录和 refresh 有限速，输入有大小/字段校验，日志不记录密码、token 和完整个人偏好。

首版 API 契约以 `shared/schema` 的 OpenAPI/JSON Schema 为准。每个 DTO 保留版本号或兼容演进策略，数据库用版本化迁移；基础设施可以替换为 PostgreSQL adapter 而不修改 Domain/Application 接口。

## 已验证的用户闭环

当前实现已经把以下路径接到真实服务端，而不是静态数组或内部函数调用：

1. `MenuServer` 启动迁移、WAL 和种子数据，健康检查、菜谱、食材和推荐端点可用。
2. Qt 客户端通过 `ClientApi` 请求今晚推荐、食材和计划；Recipe DTO 携带中文食材名称、类别、单位、常备状态和处理方式。
3. 管理面板使用真实登录和 Bearer 认证完成菜谱/食材 CRUD、草稿发布、预览、删除确认和错误态。
4. 客户端支持一周计划的添加/移除和合并清单、做饭模式步骤勾选/并行计时器/反馈、偏好设置和用户隔离缓存；登录后的 access token 可通过 `RefreshSession` 使用 refresh token 轮换。
5. 网络失败时使用带 schema/version 2 的本地缓存；401/登出会清除私有模型、活动菜谱、会话和对应缓存，并阻止旧异步写入复活。客户端可在本地设置服务端地址，切换地址会终止旧请求并切换独立的服务端缓存目录。

这些功能分别由单元、API 集成、Qt/QML、真实 HTTP 和 Playwright 测试覆盖；未测量的跨平台/高刷新率项目仍在文档末尾明确列出。

## 文件职责拆分

- `server/Api/ApiParsing*.cpp` 按 query、通用 JSON、认证、菜谱/食材、推荐和工作流 payload 拆分；`ApiRouter.cpp` 只负责组合、CORS 和请求分派。
- `client/Cpp/ClientApi*.cpp` 按状态、缓存、请求、认证和工作流拆分；模型只做稳定 role 映射和批量转换。`ConnectionEndpoint`、`ConnectionSettingsStore`、`PreferencesStore`、`DisplayMetrics`、`ThemeController`、`CookingTimerCoordinator`、`NotificationController`、`NotificationBackend` 和 `SystemUiBackend` 是独立的平台/配置边界；JNI 只存在于 Android 平台 backend。
- `client/Qml/Components` 按认证、账户、偏好、连接、设备通知和通用控件拆分；`ProfilePage` 通过账户/偏好/设备分区管理设置，`CookingPage` 只发计时语义命令，平台通知由 coordinator/backend 处理。Android Java 按 `MenuActivity`、`NotificationPermissions`、`NotificationChannels`、`NotificationScheduler`、`NotificationRenderer`、`AndroidActivityResolver` 和 `SystemUiBridge` 拆分，`NotificationBridge` 只保留 JNI 入口。
- `client/Cpp` 的缓存、偏好、连接、传输、认证和平台能力分别落在独立实现单元；`SystemUiBackend`、`NotificationBackend` 和 `DisplayMetrics` 先包含 Qt 平台宏再选择 Android 实现，桌面仍使用无副作用 backend。
- `server/Infrastructure` 的菜谱、食材、认证和计划仓储按 Read/Write/Support 实现拆分；种子导入按食材/菜谱行写入拆分，启动编排只负责事务和依赖顺序。
- `admin/web/src/Components`、`admin/web/src/Pages` 按登录入口、工作区、页面、编辑器、预览和基础控件拆分；`App.tsx` 不承载 CRUD 业务。

## 客户端交互基线

- 窄窗口使用底部导航，宽窗口使用导航栏；页面内容采用稳定的最大宽度、统一 8px 圆角和可扫描的分组间距。今晚页优先展示日期、筛选摘要和可执行方案，设置页把服务地址放在“设备”分区，避免连接配置混入烹饪主流程。
- Android target SDK 35 的 edge-to-edge 窗口由 `DisplayMetrics.TopInset` 和 `MenuActivity/SystemUiBridge` 共同处理：QML 头部绘制在状态栏背景下方，内容避开状态栏 inset，系统图标外观跟随浅色/深色主题。
- `DisplayMetrics` 读取 `QScreen::refreshRate()`；无效指标只显示 `120 Hz 预算`，有效设备显示实际 Hz 和帧预算。这里是渲染目标和可观测指标，不把 120Hz 当成未经测量的性能承诺。

拆分不是新的运行时边界：所有请求仍通过 `ClientApi`/`MenuApi`，便于在 fresh 构建中检查依赖和行为。

## 环境决策

- 本机复用已有 `C:\Qt\6.8.3` 桌面 kit，不在 C 盘安装第二套 Qt。
- C++ 依赖使用带 SHA256 的 CMake FetchContent 锁定；源码、构建和下载缓存通过 `I:\code\Menu\.cache`、`I:\code\Menu\build` 和 `I:\code\Menu\.tools` 管理。若后续 vcpkg 能在 I 盘可复现安装，再评估迁移，不改变 target 边界。
- npm cache 为 `I:\code\Menu\.npm-cache`，Gradle 为 `I:\code\Menu\.gradle`，Android SDK/NDK/JDK 只允许放 `I:\code\Menu\.tools` 或 `I:\android`。
- 代理只在下载当前进程中按需设置 `HTTP_PROXY`/`HTTPS_PROXY=http://127.0.0.1:7890`，不写入系统或用户全局配置。
- Qt Android 目标使用 `I:\code\Menu\.tools\qt-android\6.8.3\android_x86_64` 和 `I:\code\Menu\.tools\qt-android\6.8.3\android_arm64_v8a` 的源码构建 `android-clang` prefix；跨编译 prefix 按 Qt 设计不包含 host `androiddeployqt`，APK 由已验证的 `C:\Qt\6.8.3\mingw_64\bin\androiddeployqt.exe` 作为 `QT_HOST_PATH` 工具生成。
- 当前证据来自 Windows 桌面、Windows 服务端、Android x86_64 AVD 和 ARM64 APK 静态验收；ARM Android 真机实际安装、Linux 和 iOS 仍不能由这些结果推断通过。

## 当前风险

| 风险 | 当前状态 | 处理方式 |
| --- | --- | --- |
| Android Qt kit、JDK、SDK、NDK、Gradle 不在 PATH | 已满足（临时进程环境） | 所有新增工具在 I 盘；Qt target prefix 和 host deployment tool 分开记录，最终命令显式设置环境 |
| CMake 是 4.3.0-rc1 | 可用但非稳定版 | 优先验证；若出现工具链兼容问题，在 I 盘准备稳定版并记录替换原因 |
| Windows 没有 Ninja/MSVC | MinGW/LLVM-MinGW 可用 | 首版使用已有 MinGW 或 Qt kit 的编译器；CI 另行提供 Linux/Windows 证据 |
| 食谱图片授权与离线包大小 | 已建立 | 使用仓库内生成位图，记录尺寸/hash/来源；不依赖外链 |
| SQLite 单写边界与异步请求生命周期 | 已实现并测试 | 共享连接固定由单 worker 执行；异常和 completion dispatch failure 有回归测试；读连接池尚未引入 |
| 120Hz 证据 | 未满足 | AVD 报告 60Hz；保留首屏 gfxinfo，不能宣称真实设备 120fps |
| ARM Android 真机安装、Linux/iOS 构建 | 未验证 | ARM64 APK 已静态验收；实际设备和其他平台仍需要对应工具链与设备，不把 Windows/AVD 结果外推 |
