# Menu 本地运行手册

更新时间：2026-09-02

所有构建、下载、Gradle 和 npm 缓存都放在 `I:\code\Menu`。PowerShell 变量使用项目专用名称，不修改全局代理或系统环境变量。

## 服务端

从仓库根目录执行 fresh configure、构建和测试：

```powershell
cmake --fresh --preset WindowsDebug
cmake --build --preset WindowsDebug --parallel 2
ctest --preset WindowsDebug --output-on-failure
```

启动服务端：

```powershell
& 'I:\code\Menu\build\WindowsDebug\bin\MenuServer.exe' --config 'I:\code\Menu\config\Menu.example.json'
```

启动后使用真实 HTTP 客户端：

```powershell
curl.exe --fail http://127.0.0.1:8080/healthz
curl.exe --fail http://127.0.0.1:8080/readyz
curl.exe --fail http://127.0.0.1:8080/api/v1/recipes
curl.exe --fail http://127.0.0.1:8080/api/v1/ingredients
```

服务端默认使用本地 SQLite、WAL 和种子媒体；配置文件控制监听端口、数据库和媒体目录。默认数据库目录为 PascalCase 的 `server/Data`。管理 API 需要管理员 Bearer token。

## Qt 桌面客户端

`C:\Qt\6.8.3\mingw_64` 是 GCC 13.1 构建的 Qt kit，使用 `I:\code\Menu\.tools\mingw1310` 中的匹配编译器。不要把 MSYS2 GCC 16.1 与该 Qt kit 混用。

推荐使用 preset：

```powershell
cmake --preset QtDesktopGcc13Debug
cmake --build --preset QtDesktopGcc13Debug --target MenuClientApp MenuClientApiTests MenuClientModelTests MenuClientPlanModelTests MenuClientQmlTests --parallel 2
ctest --test-dir I:\code\Menu\build\QtDesktopGcc13Debug --output-on-failure
```

运行桌面客户端：

```powershell
$env:Path = 'I:\code\Menu\.tools\mingw1310\Tools\mingw1310_64\bin;C:\Qt\6.8.3\mingw_64\bin;' + $env:Path
$env:MENU_API_BASE_URL = 'http://127.0.0.1:8080'
& 'I:\code\Menu\build\QtDesktopGcc13Debug\client\MenuClient.exe'
```

运行无窗口 QML smoke 时，必须提供本地字体并用外部超时检查退出码：

```powershell
$env:QT_QPA_PLATFORM = 'offscreen'
$env:QT_QUICK_BACKEND = 'software'
$env:QT_QUICK_CONTROLS_STYLE = 'Basic'
$env:QT_QPA_FONTDIR = 'I:\code\Menu\.tools\qt-smoke-fonts'
$env:MENU_QML_SMOKE = '1'
$SmokeProcess = Start-Process -FilePath 'I:\code\Menu\build\QtDesktopGcc13Debug\client\MenuClient.exe' -PassThru
if (-not $SmokeProcess.WaitForExit(10000)) { Stop-Process -Id $SmokeProcess.Id -Force; throw 'QML smoke timeout' }
if ($SmokeProcess.ExitCode -ne 0) { throw "QML smoke exit code $($SmokeProcess.ExitCode)" }
```

## Android

Android 工具的临时环境如下；不要将这些变量写入系统环境：

```powershell
$env:JAVA_HOME = 'I:\code\Menu\.tools\jdk17\jdk-17.0.20.1+1'
$env:ANDROID_HOME = 'I:\code\Menu\.tools\android-sdk'
$env:ANDROID_SDK_ROOT = $env:ANDROID_HOME
$env:ANDROID_USER_HOME = 'I:\code\Menu\.tools\android-user-home'
$env:ANDROID_AVD_HOME = 'I:\code\Menu\.tools\android-user-home\avd'
$env:GRADLE_USER_HOME = 'I:\code\Menu\.gradle'
$env:Path = "$env:JAVA_HOME\bin;I:\code\Menu\.tools\gradle\gradle-8.10\bin;I:\code\Menu\.tools\ninja-1.13.1;$env:Path"
```

Qt Android target prefix：`I:\code\Menu\.tools\qt-android\6.8.3\android_x86_64`（x86_64 AVD）和 `I:\code\Menu\.tools\qt-android\6.8.3\android_arm64_v8a`（ARM64 真机）。两个 prefix 都由 Qt 6.8.3 源码以 `android-clang`、API 28 构建并安装，包含 `qmake.bat`、`qt-cmake.bat` 和 `lib\cmake\Qt6\Qt6Config.cmake`。Qt 跨编译 prefix 按设计不生成 host `androiddeployqt`；CMake 使用 `QT_HOST_PATH=C:\Qt\6.8.3\mingw_64` 的已验证 `androiddeployqt.exe` 打包。

配置、编译、打包和安装：

```powershell
$env:JAVA_HOME = 'I:\code\Menu\.tools\jdk17\jdk-17.0.20.1+1'
$env:GRADLE_USER_HOME = 'I:\code\Menu\.gradle'
cmake --fresh --preset QtAndroidDebug
cmake --build --preset QtAndroidDebug --target apk --parallel 2
$AndroidApk = 'I:\code\Menu\build\QtAndroidDebug\client\android-build\build\outputs\apk\debug\android-build-debug.apk'
& 'I:\code\Menu\.tools\android-sdk\platform-tools\adb.exe' install -r $AndroidApk
& 'I:\code\Menu\.tools\android-sdk\platform-tools\adb.exe' shell monkey -p com.menu.cookflow 1
```

ARM64 Android 16 真机使用独立 preset 和 APK，不能拿上面的 x86_64 AVD 包安装：

```powershell
$env:JAVA_HOME = 'I:\code\Menu\.tools\jdk17\jdk-17.0.20.1+1'
$env:GRADLE_USER_HOME = 'I:\code\Menu\.gradle'
cmake --fresh --preset QtAndroidArm64Debug
cmake --build --preset QtAndroidArm64Debug --target apk --parallel 2
$AndroidApk = 'I:\code\Menu\build\QtAndroidArm64Debug\client\android-build\build\outputs\apk\debug\android-build-debug.apk'
$Adb = 'I:\code\Menu\.tools\android-sdk\platform-tools\adb.exe'
& $Adb shell getprop ro.product.cpu.abilist
& $Adb install -r $AndroidApk
& $Adb shell monkey -p com.menu.cookflow 1
```

若设备已安装旧 debug 包且提示 `INSTALL_FAILED_UPDATE_INCOMPATIBLE`，先确认无需保留旧数据，再执行 `& $Adb uninstall com.menu.cookflow` 后重新安装。当前 ARM64 包的静态 ABI 验收记录在 `docs/AndroidEvidence.md`。

首次在 Android 13 及以上设备使用提醒时，在“我的 → 设备与通知”点“允许通知”；“通知设置”会打开系统渠道页面，做饭计时在后台通过 Android AlarmManager 触发。当前 target SDK 为 35，Android 15/16 的 edge-to-edge 窗口已在 Android 36 AVD 验证。系统对声音、震动和深色/浅色呈现拥有最终控制权。

历史窗口化 AVD 验收包来自 fresh `QtAndroidFinal` 构建；ARM64 手机验收包来自 `QtAndroidArm64Debug` 构建。AVD 访问宿主服务时，默认 `Menu.example.json` 的 `127.0.0.1` 只适合宿主机；仅在本地窗口化 AVD 验证期间使用临时副本将监听地址改为 `0.0.0.0`，验证后删除副本，不把该绑定用于部署。

当前可复现实例是 AVD `MenuApi36`、Android 36 Google APIs x86_64。窗口化实例适合截图；完成后使用精确设备命令停止：

```powershell
& 'I:\code\Menu\.tools\android-sdk\platform-tools\adb.exe' -s emulator-5554 emu kill
```

## 管理面板与 E2E

```powershell
$env:npm_config_cache = 'I:\code\Menu\.npm-cache'
npm --prefix .\admin\web ci --cache I:\code\Menu\.npm-cache
npm --prefix .\admin\web run build
npm --prefix .\admin\web run test
$env:PLAYWRIGHT_BROWSERS_PATH = 'I:\code\Menu\.tools\playwright-browsers'
npm --prefix .\admin\web run test:e2e
```

E2E 只使用本机服务和脚本中的本地测试账号；不会访问外部账户、支付或真实数据。截图保存在 `docs\evidence-admin-*.png`。

## 压测与证据

服务端运行后：

```powershell
node .\scripts\RunApiBenchmark.mjs
```

该脚本固定请求数量和并发度，输出错误率、吞吐、p50/p95/p99。Android/桌面截图、gfxinfo、logcat、Qt smoke 日志和下载哈希分别记录在 `docs` 与 `docs\Downloads.md`。

## 清理

验证结束后只停止本轮记录的服务端、Vite 和 AVD PID/会话；不使用宽泛的进程名、不删除工作树、不清理其他项目缓存。`build`、`.tools`、`.gradle` 和 `.npm-cache` 是可复用的 I 盘目录。
