# Android 验证证据

更新时间：2026-09-02

## 工具链

- Qt 6.8.3 target prefix：`I:\code\Menu\.tools\qt-android\6.8.3\android_x86_64`（x86_64 AVD）和 `I:\code\Menu\.tools\qt-android\6.8.3\android_arm64_v8a`（ARM64 真机）。
- target 配置：`android-clang`、API 28、Clang 17.0.2；x86_64 与 `arm64-v8a` 的 Qt source build 和 `cmake --install` 均为退出码 0，两个 prefix 均已检查 `bin\qmake.bat`、`bin\qt-cmake.bat` 和 `lib\cmake\Qt6\Qt6Config.cmake`。
- Qt cross-build 按设计不生成 host `androiddeployqt`；APK 使用已安装的 Windows host kit `C:\Qt\6.8.3\mingw_64\bin\androiddeployqt.exe`。
- JDK 17.0.20.1、Android SDK platform-tools 37.0.1、platform android-36、build-tools 36.0.0、NDK 26.1.10909125、Gradle 8.10 和 Ninja 1.13.1 均在 I 盘。

## APK

```text
build/QtAndroidDebug/client/android-build/build/outputs/apk/debug/android-build-debug.apk
size: 28,758,117 bytes
sha256: A5A0A2C8FCA5534D733B27719733180FC12332B92DE1C6B1674CB7B2AD7E1CD2
package: com.menu.cookflow
abi: x86_64
```

`QtAndroidDebug` 的 `MenuClientApp` Ninja 构建和 Gradle `assembleDebug` 均退出码 0；当前 debug 包使用 v2 签名，x86_64 AVD 安装和启动均通过。

## ARM64 真机包

```text
build/QtAndroidArm64Debug/client/android-build/build/outputs/apk/debug/android-build-debug.apk
size: 28,122,475 bytes
sha256: ACD2CC349D0FB9F2BD9DF70838E7FA261725ABD9C8A10B3E680241AB05CB0B57
package: com.menu.cookflow
abi: arm64-v8a
minSdk: 28
targetSdk: 35
```

`QtAndroidArm64Debug` 的 CMake `apk` 目标和 Gradle `assembleDebug` 均退出码 0。`aapt2 dump badging` 报告 `native-code: 'arm64-v8a'`；APK 内 82 个 native 条目全部位于 `lib/arm64-v8a`，未发现 x86/x86_64 条目；`POST_NOTIFICATIONS`、`RECEIVE_BOOT_COMPLETED`、`SCHEDULE_EXACT_ALARM`、通知 receiver、白色小图标、主题/系统栏和分模块 QML 均已合入。debug manifest 允许局域网 HTTP，release manifest 已验证 `usesCleartextTraffic=false`。当前没有连接实体手机，因此尚未在用户的 Android 16 设备上执行 adb 安装。

## 本轮 UI 与通知回归

- x86_64 AVD fresh APK：`build/QtAndroidDebug/client/android-build/build/outputs/apk/debug/android-build-debug.apk`，size `28,758,117` bytes，sha256 `A5A0A2C8FCA5534D733B27719733180FC12332B92DE1C6B1674CB7B2AD7E1CD2`；本轮 `MenuClientApp` target 重命名、`MenuActivity`/系统栏适配和 UI 调整后卸载、全新安装、启动和在线首屏回归通过。常亮控制已接入客户端，AVD 窗口 flag 和桌面单测均已验证。
- 当前 UI 截图：`build/QtAndroidDebug/ui-final-current.png`、`build/QtAndroidDebug/ui-final-profile-device.png`、`build/QtAndroidDebug/ui-final-profile-device-dark.png`、`build/QtAndroidDebug/ui-final-cooking.png`、`build/QtAndroidDebug/ui-live-awake-detail-new.png`；浅色主题使用深色状态栏图标、深色主题使用白色状态栏图标，服务器地址设置、刷新率、通知权限、通知设置和做饭页面均已实测。最新做饭页截图确认无时长步骤不显示 0 秒计时。
- 最新 AVD 常亮取证：进入做饭页后 `dumpsys window windows` 的应用窗口包含 `fl=KEEP_SCREEN_ON`，返回后该 flag 消失；退出做饭页的清理逻辑由 `ScreenAwakeController` 单测覆盖。截图 `ui-live-awake-detail-new.png`，SHA256 `F11F793159A6B2A3BB56AFABA06BBF852430F4D20B6E62566D25FD02A5F90494`。
- `dumpsys notification` 看到 `cooking_timers`（高优先级、双震动）与 `meal_reminders`（默认优先级）两个渠道，并实际发布 `Menu 通知测试`；通知应用总开关/渠道关闭时 `PermissionGranted` 会变为 false。本轮重新打包后，x86_64 AVD 卸载/安装/启动与测试通知回归均通过，日志未出现 `AndroidRuntime`、`FATAL EXCEPTION`、`UnsatisfiedLinkError` 或 `QmlWarning`。

## 窗口化 AVD

- AVD：`MenuApi36`，Android 36 Google APIs x86_64，窗口化参数 `-no-audio -no-boot-anim -gpu host`；target SDK 35 的 edge-to-edge 窗口已实测。
- 本轮使用窗口化 emulator 采集截图；进程 PID 属于临时运行态，不作为可复用证据。
- `androidboot.qemu.vsync=60`，因此只能验证 60Hz AVD 行为，不能宣称 120fps。gfxinfo 是短时启动样本，不是完整滚动性能剖析。
- 早期 `-no-window` 黑帧只作为失败背景保留；当前窗口化截图有真实 Qt 内容，排除了将黑屏直接归因于 screencap 的判断。

## 当前 APK 截图

- 在线首屏推荐：`build/QtAndroidDebug/ui-final-current.png`，SHA256 `542674ED4D382BF1B9593773F15542462B6810AC47424D69C647126BC5BA88EE`。
- 我的页与服务端设置：`build/QtAndroidDebug/ui-final-current-profile.png`，SHA256 `FAF36D07671ADDE73D6EA276A8D15564F1281BF1B44277A4264C145A20FC03D4`。
- 深色主题与刷新率/通知设置：`build/QtAndroidDebug/ui-final-current-profile-dark-notifications.png`，SHA256 `0D1D28BF9629BF4D9A8B779A3495A600B88868D653BE27F716178138F16D1A74`。
- 深色主题设置中间状态：`build/QtAndroidDebug/ui-final-current-profile-dark.png`，SHA256 `C20DAE38B09C9CE5EBD663C92546DF47070E850FDB7C4BF2E7C93F9857BEEEA4`。

## 日志

- 首屏 clean logcat：`docs/evidence-android-final-clean-logcat.txt`，SHA256 `22261512CEC7A48E452FC1DE4BBADF64D4B5C49AE364224AABD3C38D9BC181E4`。
- 首屏 gfxinfo：`docs/evidence-android-final-clean-gfxinfo.txt`，SHA256 `E02925BCD6F887DC32255CD5059464803B2EE39CD40AB96D8D8668F7F21FEEE5`。
- 认证日志：`docs/evidence-android-final-auth-logcat.txt`，SHA256 `E9E417594D0F09C0EB3C92A61B30FFDD28B770ADD3E6D1F4CCB0CB66FA415C8F`。
- 做饭交互日志：`docs/evidence-android-final-cooking-logcat.txt`，SHA256 `51AE5E5B21FBCC78ECD2406AB065AA979CF9131A5C0E1A44DF862547A312F113`。
- 应用日志中未发现 `QmlWarning`、QML binding/Connections warning、QObject 跨线程操作、崩溃或 fatal exception；原始 PID 日志保留了 Android/Qt 的 `Qt A11Y`、HWUI 和冷启动帧提示，这些是平台运行时提示，不是 QML 绑定错误。完整 warning 取舍以此说明为准。

安装和打包的退出码记录在 `docs/evidence-android-install-final.txt` 与 `docs/evidence-android-build-closeout.txt`：精确卸载和安装均为 0，CMake `apk` 目标为 0，Gradle `assembleDebug` 为 0。

## 未验证

- ARM Android 16 真机的实际安装和运行、真实 120Hz 设备、Android release 签名包和 iOS 没有当前主机证据；ARM64 APK 本身已完成静态验收。设备页会将实测屏幕刷新率与缺失指标时的 `120 Hz 预算` 分开显示，不把预算冒充硬件刷新率。
- QML Profiler/Perfetto 的完整首屏与滚动帧时间尚未建立；当前只有桌面 smoke 和 Android gfxinfo 启动样本。
