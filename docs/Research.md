# Menu 调研与技术依据

更新时间：2026-09-02

本文件记录已核对的官方资料、对实现的约束和仍需实测的事项。网页内容会更新，构建时以本机 Qt 6.8.3、实际 Boost 版本和锁定下载哈希为准。

当前服务端嵌入 SQLite 3.53.4 amalgamation；它来自 SQLite 官方 2026 下载页，构建使用本地 SHA256 锁定值。系统 PATH 中的 SQLite 3.51.1 仅用于环境诊断，不作为 Menu 的生产链接库。

## 官方资料

| 主题 | 来源 | 对本项目的结论 |
| --- | --- | --- |
| Qt Quick 性能 | https://doc.qt.io/qt-6/qtquick-performance.html | 网络/重计算使用异步事件驱动和 worker；避免在一帧中执行阻塞操作；用 QML Profiler 找真正的瓶颈；ListView delegate 保持简单。 |
| Qt Quick Image | https://doc.qt.io/qt-6/qml-qtquick-image.html | 大图使用 `asynchronous: true`，设置 `sourceSize` 限制内存；必要时使用加载中保留旧图，避免厨房页面闪烁。 |
| Qt Quick ListView | https://doc.qt.io/qt-6/qml-qtquick-listview.html | delegate 会按需创建和销毁，业务状态放模型/会话而不是 delegate；使用稳定尺寸和适量缓存，避免动态 role。 |
| Boost.Asio | https://www.boost.org/doc/libs/latest/doc/html/boost_asio/overview.html | 用 io_context、executor、C++20 coroutine 组织跨平台异步 I/O；线程数量受控。 |
| Boost.Asio strand | https://www.boost.org/doc/libs/latest/doc/html/boost_asio/reference/strand.html | strand 提供序列化 handler 调用，作为连接写队列和 SQLite writer actor 的并发语义基础。 |
| Boost.Beast WebSocket | https://www.boost.org/doc/libs/latest/libs/beast/doc/html/beast/using_websocket.html | Beast WebSocket 建立在 Asio 的异步模型上；stream 非线程安全，必须绑定 strand 或单一 executor。首版只在同步需求成立时启用。 |
| SQLite WAL | https://www.sqlite.org/wal.html | WAL 提供读写并发，但所有进程必须在同一主机；要处理 checkpoint、busy 和 WAL 文件大小，不能放网络文件系统。 |
| Qt Android 部署 | https://doc.qt.io/qt-6/android-deploy-qt-tool.html | CMake/qmake 先生成 deployment JSON，再由 `androiddeployqt` 产出 APK；Android kit、SDK、NDK、JDK 仍需独立安装与验证。 |

## 当前环境实测

- Qt 6.8.3 Android target 使用官方源码 `qt-everywhere-src-6.8.3`，`android-clang`、API 28、Clang 17.0.2，x86_64 与 `arm64-v8a` 两套构建和 `cmake --install` 的退出码均为 0。安装 prefix 分别是 `I:\code\Menu\.tools\qt-android\6.8.3\android_x86_64` 和 `I:\code\Menu\.tools\qt-android\6.8.3\android_arm64_v8a`，已检查 `qmake.bat`、`qt-cmake.bat` 和 `Qt6Config.cmake`。
- Qt 的 `androiddeployqt` 特性要求非 cross compile；因此 Android target prefix 不会包含该 host 工具。当前使用 Qt 6.8.3 host kit 的 `C:\Qt\6.8.3\mingw_64\bin\androiddeployqt.exe`，这也是 Qt Android CMake macros 的 `QT_HOST_PATH` 选择，并已成功生成 APK。
- `MenuApi36` AVD 的 `androidboot.qemu.vsync=60`，不是 120Hz 设备。窗口化对照可以通过 `adb screencap` 看到完整 Qt 页面；`-no-window` headless 实例只适合作为非视觉启动测试，不能把黑帧当作产品截图。
- Android 首屏第一次出现 `loadFromModule` 不返回、窗口化截图全黑；将 `Main.qml` 的五页 eager instantiation 改为当前页首次加载并保留实例的 `Loader` 后，Android 日志出现 `qml-loaded 1`，窗口化截图显示真实推荐、中文文案、本地位图和底部导航。
- 客户端连接设置使用严格的 `http/https + host + port` 解析，地址写入本地 `connection.json`；切换地址会清理内存会话并切换独立的服务端缓存目录，健康检查请求 `/healthz`。Android manifest 允许显式配置的局域网 HTTP，同时 UI 对明文连接给出警告。
- Android 通知由 `MenuClientPlatform` 的 `NotificationController` 通过独立 `NotificationBackend` 调用 Java bridge，`CookingTimerCoordinator` 负责稳定的菜谱/步骤 ID；系统栏由独立 `SystemUiBackend/SystemUiBridge` 管理，`MenuActivity` 在 resume/focus 生命周期重新应用系统栏外观。Android 8+ 创建“做饭计时/备菜提醒”渠道，Android 13+ 请求 `POST_NOTIFICATIONS`，计时使用 `AlarmManager` 的 elapsed realtime，并持久化 wall-clock deadline 以便设备重启恢复；系统负责深浅色、声音和用户对渠道震动设置的最终控制。
- Android 官方自适应布局建议按窗口尺寸切换底部导航、导航栏和列表/详情结构；本客户端以 720/980 logical px 断点切换导航和菜谱网格，并把状态保存在 C++ 模型而不是 delegate。edge-to-edge 场景通过平台状态栏高度 inset 避免内容被系统栏覆盖。

## 关键设计推导

### UI 帧预算

120Hz 的目标预算是 8.33 ms/frame，但这是测量目标而不是产品承诺。客户端读取 `QScreen::refreshRate()`，无有效值时使用 120Hz fallback；视觉帧由 Qt 的 vsync/`QWindow::requestUpdate()` 驱动，不用定时器模拟刷新率。首屏和滚动路径上禁止同步网络、SQLite、磁盘和大图解码；图片使用本地资源和明确的 sourceSize；列表只渲染需要的 delegate。设备若只支持 60Hz 或 90Hz，就记录实际刷新率和帧时间，不宣称 120fps。

### SQLite 与异步服务端

SQLite 的事务和 WAL 能满足本地首版，但 SQLite API 本身是同步调用。因此 HTTP event loop 只负责接收请求、校验和投递，真正的数据库调用运行在 storage executor。写入通过一个串行 actor 维护事务顺序，读连接数量受控；响应回到原请求 executor 前检查请求是否已经取消。

### 确定性推荐

推荐先用可测试规则：过滤过敏原/忌口，按可用时间与厨具过滤，再按已有食材覆盖率、偏好菜系、难度和近期重复度排序，返回最多三个方案。`RecommendationProvider` 是 Application 端口，未来受控模型只能返回 schema 化 JSON，服务端要校验数量、单位、过敏原、步骤先后和时间可行性。

### 图片与版权

首版图片放在仓库 `assets/media`，不依赖外链；每个 `MediaAsset` 保存来源、授权说明、宽高和 hash。没有可用照片时使用本地生成的位图占位图，并在 `docs/Downloads.md` 和资源清单说明来源，客户端离线仍能打开。

## 仍未完成的验证

- Linux 编译器、CTest 和 Linux 运行时验证。
- QML Profiler/Perfetto 的完整首屏与滚动帧时间；当前只保留桌面 smoke 和 Android gfxinfo 启动样本。
- ARM Android 真机、120Hz 设备和 Android release 签名包；当前 debug 包已使用 target SDK 35 并在 Android 36 AVD 验证 edge-to-edge。
- 生产 TLS 证书/反向代理部署；当前本地 HTTP 验证不等于公网部署验证。
