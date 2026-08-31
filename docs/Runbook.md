# Menu 本地运行手册

状态：Phase 0 基线。Phase 1 完成 CMake 和服务端入口后，下面的构建命令才会产生可执行文件。

## 环境变量

在 PowerShell 中从仓库根目录执行：

```powershell
$env:MenuRoot = 'I:\code\Menu'
$env:MenuBuild = 'I:\code\Menu\build'
$env:FETCHCONTENT_BASE_DIR = 'I:\code\Menu\.cache\cmake-fetch'
$env:GRADLE_USER_HOME = 'I:\code\Menu\.gradle'
$env:npm_config_cache = 'I:\code\Menu\.npm-cache'
$env:ANDROID_HOME = 'I:\code\Menu\.tools\android-sdk'
$env:ANDROID_SDK_ROOT = $env:ANDROID_HOME
```

不要覆盖 `HOME` 或其他系统通用变量。当前没有 Java/ADB，因此不要设置不存在的 `JAVA_HOME`；Android 工具链安装后再写入实际路径。

## 服务端

Phase 1 目标命令：

```powershell
cmake --preset WindowsDebug
cmake --build --preset WindowsDebug --parallel 2
ctest --preset WindowsDebug --output-on-failure
.\build\WindowsDebug\MenuServer.exe --config .\config\Menu.json
```

启动后使用真实 HTTP 客户端验证：

```powershell
curl.exe --fail http://127.0.0.1:8080/healthz
curl.exe --fail http://127.0.0.1:8080/readyz
curl.exe --fail http://127.0.0.1:8080/api/v1/recipes
```

端口、数据库目录和静态资源目录由配置文件控制；生产默认不暴露管理 API 到公网。当前这些文件尚未创建，命令属于已记录的 Phase 1 验收入口。

## 管理面板

依赖和构建必须使用 I 盘 npm cache：

```powershell
npm --prefix .\admin\web ci --cache I:\code\Menu\.npm-cache
npm --prefix .\admin\web run build
npm --prefix .\admin\web run test
```

开发调试时先启动服务端，再启动 Vite dev server；面板必须通过真实 API 读取和写入，不以静态 JSON 作为最终验证。

## Qt Quick 桌面客户端

优先使用已有 Qt kit 的 CMake 工具：

```powershell
& 'C:\Qt\6.8.3\mingw_64\bin\qt-cmake.bat' -S . -B I:\code\Menu\build\QtDesktopDebug -DCMAKE_BUILD_TYPE=Debug
cmake --build I:\code\Menu\build\QtDesktopDebug --parallel 2
```

客户端运行时指向本机服务端。网络、磁盘和图片加载路径必须保持异步；UI 测试覆盖导航、勾选、退出恢复、计时器、离线状态和尺寸变化。

## Android

只有在以下工具实际安装并通过版本检查后执行：Qt Android kit、Android SDK platform/build-tools、NDK、JDK 和 Gradle。它们必须位于 I 盘；使用 Qt 生成的 deployment JSON，再调用 `androiddeployqt` 产出 debug APK。当前主机缺少这些前置条件，Android 构建是未验证状态。

## 验证与清理

每个实质改动后至少执行：

```powershell
git diff --check
cmake --build --preset WindowsDebug --parallel 2
ctest --preset WindowsDebug --output-on-failure
```

临时服务端进程用 PowerShell 记录 PID，验证结束后停止该 PID；不要用会误杀其他项目的宽泛进程名。构建和下载目录都在 I 盘，可按具体 preset 清理，不删除用户文件或其他项目目录。
