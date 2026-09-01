# Android 验证证据

更新时间：2026-09-01

## 工具链

- Qt 6.8.3 target prefix：`I:\code\Menu\.tools\qt-android\6.8.3\android_x86_64`。
- target 配置：`android-clang`、x86_64、API 28、Clang 17.0.2；Qt source build 和 `cmake --install` 均为退出码 0，已检查 `bin\qmake.bat`、`bin\qt-cmake.bat` 和 `lib\cmake\Qt6\Qt6Config.cmake`。
- Qt cross-build 按设计不生成 host `androiddeployqt`；APK 使用已安装的 Windows host kit `C:\Qt\6.8.3\mingw_64\bin\androiddeployqt.exe`。
- JDK 17.0.20.1、Android SDK platform-tools 37.0.1、platform android-36、build-tools 36.0.0、NDK 26.1.10909125、Gradle 8.10 和 Ninja 1.13.1 均在 I 盘。

## APK

```text
build/QtAndroidFinal/client/android-build/build/outputs/apk/debug/android-build-debug.apk
size: 27,214,651 bytes
sha256: FFBE9972FE920957E637385B4A29BB9B3E830FA3B375E9F621B092E716F2FFAE
package: com.menu.cookflow
abi: x86_64
```

`QtAndroidFinal` 的 Ninja 构建（56/56）和 Gradle `assembleDebug` 均退出码 0；旧 debug 包签名冲突时 `adb install -r` 的退出码为 1，随后精确卸载并重新安装当前包，`adb uninstall` 和 `adb install` 均退出码 0，详见 `docs/evidence-android-install-final.txt`。

## 窗口化 AVD

- AVD：`MenuApi36`，Android 36 Google APIs x86_64，窗口化参数 `-no-audio -no-boot-anim -gpu host`。
- 本轮窗口化 emulator PID：`3164`；应用进程 PID：`3711`（仅对应本轮最终 APK 采集）。
- `androidboot.qemu.vsync=60`，因此只能验证 60Hz AVD 行为，不能宣称 120fps。gfxinfo 是短时启动样本，不是完整滚动性能剖析。
- 早期 `-no-window` 黑帧只作为失败背景保留；当前窗口化截图有真实 Qt 内容，排除了将黑屏直接归因于 screencap 的判断。

## 当前 APK 截图

- 首屏在线推荐：`docs/evidence-android-final-clean-first-screen.png`，SHA256 `1BD781B8B4D5EEC1FEFFA5E18241A210D8DC5D7D240E705877002D92B8E1343A`。
- 登录后设置：`docs/evidence-android-final-profile-authenticated.png`，SHA256 `717AB3B3476D83D5DA129A92F45DB0CEFA344696639EF58E5939A4635174E580`。
- 登录后今晚：`docs/evidence-android-final-authenticated-tonight.png`，SHA256 `1F28B573D911D8FEE9D1141CABFBB82B0C4C90E71267AFDDC1EF4CAACE4A2A3D`。
- 本周方案：`docs/evidence-android-final-week-options.png`，SHA256 `6D57B1BE6EEE6BB2BD5BFD57ABC9FC3DF47A187D4A092C80E5895EF7E721321A`。
- 两道菜和合并清单：`docs/evidence-android-final-week-merged-two-recipes.png`，SHA256 `6933E267213FAC5CF8F60C3E777EC01CB04DF1EE3D6CBA01EF047A6A219EC314`。
- 做饭页：`docs/evidence-android-final-cooking.png`，SHA256 `D5121D312DDA0491C2713EAD4047A935015051C21FDD8AF82B2E6BFC3F231A3E`。
- 计时运行中：`docs/evidence-android-final-cooking-timer.png`，SHA256 `EC930674AD096B5657326D2D947E7A6BB68295D8446CDC5B6D35A6CD064347CF`。
- 第一步完成：`docs/evidence-android-final-cooking-complete.png`，SHA256 `7FC856E7334A26C2FC26DF09EB97CC738CC4DA257A942B15AE8320B5BEDF2B27`。

## 日志

- 首屏 clean logcat：`docs/evidence-android-final-clean-logcat.txt`，SHA256 `22261512CEC7A48E452FC1DE4BBADF64D4B5C49AE364224AABD3C38D9BC181E4`。
- 首屏 gfxinfo：`docs/evidence-android-final-clean-gfxinfo.txt`，SHA256 `E02925BCD6F887DC32255CD5059464803B2EE39CD40AB96D8D8668F7F21FEEE5`。
- 认证日志：`docs/evidence-android-final-auth-logcat.txt`，SHA256 `E9E417594D0F09C0EB3C92A61B30FFDD28B770ADD3E6D1F4CCB0CB66FA415C8F`。
- 做饭交互日志：`docs/evidence-android-final-cooking-logcat.txt`，SHA256 `51AE5E5B21FBCC78ECD2406AB065AA979CF9131A5C0E1A44DF862547A312F113`。
- 应用日志中未发现 `QmlWarning`、QML binding/Connections warning、QObject 跨线程操作、崩溃或 fatal exception；原始 PID 日志保留了 Android/Qt 的 `Qt A11Y`、HWUI 和冷启动帧提示，这些是平台运行时提示，不是 QML 绑定错误。完整 warning 取舍以此说明为准。

安装和打包的退出码记录在 `docs/evidence-android-install-final.txt` 与 `docs/evidence-android-build-closeout.txt`：精确卸载和安装均为 0，CMake `apk` 目标为 0，Gradle `assembleDebug` 为 0。

## 未验证

- ARM Android 真机、真实 120Hz 设备、Android release 签名包和 iOS 没有当前主机证据。
- QML Profiler/Perfetto 的完整首屏与滚动帧时间尚未建立；当前只有桌面 smoke 和 Android gfxinfo 启动样本。
