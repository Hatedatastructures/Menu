# 下载、版本与缓存登记

更新时间：2026-09-01

## Phase 0 下载记录

Phase 0 没有下载新工具或库，没有使用代理，也没有修改全局代理配置。Phase 1 下载了下面三项，未使用代理；下载发生在当前 PowerShell 进程，未修改全局配置。下面的路径是后续唯一允许的项目缓存位置：

| 用途 | 目标路径 |
| --- | --- |
| CMake FetchContent 源码/构建缓存 | `I:\code\Menu\.cache\cmake-fetch` |
| 项目构建 | `I:\code\Menu\build` |
| 第三方工具 | `I:\code\Menu\.tools` |
| Android SDK/NDK/JDK | `I:\code\Menu\.tools\android-sdk` 或 `I:\android` |
| Gradle 用户目录 | `I:\code\Menu\.gradle` |
| npm 缓存 | `I:\code\Menu\.npm-cache` |

## 已下载记录

| 工具/库 | 版本 | URL | 文件路径 | 大小 | SHA256 | 代理 |
| --- | --- | --- | --- | ---: | --- | --- |
| Boost | 1.89.0 | https://archives.boost.io/release/1.89.0/source/boost_1_89_0.tar.bz2 | `I:\code\Menu\.cache\cmake-fetch\boostsource-subbuild\boostsource-populate-prefix\src\boost_1_89_0.tar.bz2` | 154699732 bytes | `85A33FA22621B4F314F8E85E1A5E2A9363D22E4F4992925D4BB3BC631B5A0C7A` | 未使用 |
| GoogleTest | 1.16.0 | https://github.com/google/googletest/archive/refs/tags/v1.16.0.tar.gz | `I:\code\Menu\.cache\cmake-fetch\googletest-subbuild\googletest-populate-prefix\src\v1.16.0.tar.gz` | 876245 bytes | `78C676FC63881529BF97BF9D45948D905A66833FBFA5318EA2CD7478CB98F399` | 未使用 |
| SQLite amalgamation | 3.53.4 | https://www.sqlite.org/2026/sqlite-amalgamation-3530400.zip | `I:\code\Menu\.cache\sqlite\sqlite-amalgamation-3530400.zip` | 2946650 bytes | SHA256 `1E71DDF93849C6A6ECF58B827C0692073D2DD7EE40196158068F7B29F422E87D`; 官方 SHA3-256 `628a44cfe82c66aed1ccbbe85a562d2e33ebe64b3288981ed76285612227934e` | 未使用 |
| Qt MinGW toolchain | GCC/G++ 13.1.0, MinGW-w64 x86_64 POSIX SEH | https://download.qt.io/online/qtsdkrepository/windows_x86/desktop/tools_mingw1310/qt.tools.win64_mingw1310/13.1.0-202407240918mingw1310.7z | `I:\code\Menu\.tools\mingw1310\qt.tools.win64_mingw1310-13.1.0.7z`; 解压到 `I:\code\Menu\.tools\mingw1310\Tools\mingw1310_64` | 112089140 bytes | SHA256 `A7B502294A903B64FCCD4A41BA39F48D3F6C9A5ECE167F0CFAEAC920D15B344F`; 官方 SHA1 `23c4bdecc1a2d2588f1b32b3e6cbf64d2a7ee852` | 未使用 |

## Android 与视觉验证工具

| 工具/库 | 版本 | URL | 文件路径/安装路径 | 大小 | SHA256 | 代理 |
| --- | --- | --- | --- | ---: | --- | --- |
| Temurin JDK | 17.0.20.1+1 | https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.20.1%2B1/OpenJDK17U-jdk_x64_windows_hotspot_17.0.20.1_1.zip | `I:\code\Menu\.tools\downloads\OpenJDK17U-jdk_x64_windows_hotspot_17.0.20.1_1.zip`; 解压 `I:\code\Menu\.tools\jdk17\jdk-17.0.20.1+1` | 190817615 bytes | `E53A79C3C3D86865BD7E787903884331068E71321714FFD44F145785AFFC7CB0` | 未使用 |
| Android command-line tools | 15859902 / sdkmanager 22.0 | https://dl.google.com/android/repository/commandlinetools-win-15859902_latest.zip | `I:\code\Menu\.tools\downloads\commandlinetools-win-15859902_latest.zip`; SDK `I:\code\Menu\.tools\android-sdk` | 155655386 bytes | `90AE805D20434428BFFCB699C290860F19BB5F66A67E6B330067E3DE801FB04A` | 未使用 |
| Android SDK platform-tools | 37.0.1 | https://developer.android.com/tools/releases/platform-tools | `I:\code\Menu\.tools\android-sdk\platform-tools` | `adb.exe` SHA256 `B4A6B455702684652CCCF7B46258B29E653538904359A58FD4931CF3EF286B3F` | 已安装 | 未使用 |
| Android SDK platform/build tools | android-36 / build-tools 36.0.0 | https://developer.android.com/tools/releases/platforms | `I:\code\Menu\.tools\android-sdk\platforms\android-36`; `...\build-tools\36.0.0` | `android.jar` SHA256 `D9EB9DA824D9E247A352F570F01E1169E725B2954BCA9E283A71786C59B59F9A`; `aapt2.exe` SHA256 `BABF3122E515DDB954C5AC4669E085CE990536C035E3072DE30127BDDD6E3608` | 已安装 | 未使用 |
| Android NDK | 26.1.10909125 | https://developer.android.com/ndk/downloads | `I:\code\Menu\.tools\android-sdk\ndk\26.1.10909125` | source.properties SHA256 `96CDDD3DEA11A24DC4F563280E350FE566F1477D64CC358967838A90C66A23D1` | 已安装 | 未使用 |
| Gradle | 8.10 | https://services.gradle.org/distributions/gradle-8.10-bin.zip | `I:\code\Menu\.tools\downloads\gradle-8.10-bin.zip`; 解压 `I:\code\Menu\.tools\gradle\gradle-8.10`; 用户缓存 `I:\code\Menu\.gradle` | 136713202 bytes | `5B9C5EB3F9FC2C94ABAEA57D90BD78747CA117DDBBF96C859D3741181A12BF2A` | 未使用 |
| Ninja | 1.13.1 | https://github.com/ninja-build/ninja/releases/download/v1.13.1/ninja-win.zip | `I:\code\Menu\.tools\downloads\ninja-win-v1.13.1.zip`; 解压 `I:\code\Menu\.tools\ninja-1.13.1` | 289808 bytes | `26A40FA8595694DEC2FAD4911E62D29E10525D2133C9A4230B66397774AE25BF` | 未使用 |
| Qt source | 6.8.3 | https://download.qt.io/official_releases/qt/6.8/6.8.3/single/qt-everywhere-src-6.8.3.tar.xz | `I:\code\Menu\.tools\downloads\qt-everywhere-src-6.8.3.tar.xz`; source `I:\code\Menu\.tools\qt-source\qt-everywhere-src-6.8.3`; Android prefixes `I:\code\Menu\.tools\qt-android\6.8.3\android_x86_64` and `I:\code\Menu\.tools\qt-android\6.8.3\android_arm64_v8a` | 994812276 bytes | `CDD3A69967208276BB01AF7ACE7DBA0BA53E679F886A4CBE624225C60FB73F2C` | 未使用 |
| Qt online installer | 6.8.3 installer | https://download.qt.io/official_releases/online_installers/qt-online-installer-windows-x64-online.exe | `I:\code\Menu\.tools\downloads\qt-online-installer-windows-x64-online.exe` | 58051784 bytes | `AE919BC9B224B8CCDADA69EC787A9F69330001F227F3FCBFB4A11A4ADB3786F6` | 未使用；无 Qt Account/许可证，未使用安装器 kit |
| Playwright npm | 1.62.1 | https://registry.npmjs.org/playwright/-/playwright-1.62.1.tgz | `I:\code\Menu\admin\web\node_modules\playwright`; npm cache `I:\code\Menu\.npm-cache` | integrity `sha512-0M+L3LAD8/nm554LOla9Ayx0j0tmFZ0FBcoQ7F1VuVHpM/XpiC8RcDzBQB8W5+hA8L22THxELzeF+2WcUzvcLg==` | 包管理器校验 | 未使用 |
| Playwright Chromium | 151.0.7922.34 / revision 1234 | https://cdn.playwright.dev/builds/cft/151.0.7922.34/win64/chrome-win64.zip | `I:\code\Menu\.tools\playwright-browsers\chromium-1234\chrome-win64` | chrome.exe 4024832 bytes | `409805A16D6416087E6B2F778DF1CF8F7BBB267D6B99F6B5BB0A618EACE234F2` | 未使用 |

## 平台错误资产

以下 5 个包确实来自下载页，但文件名明确是 Linux host，Windows 主机不能使用。它们保留用于审计，未解压、未进入 Qt prefix、未用于 APK 证据：

| 文件 | URL 基址 | 大小 | SHA256 | 状态 |
| --- | --- | ---: | --- | --- |
| `qtbase-Linux-RHEL_8_10-Clang-Android-Android_ANY-X86_64.7z` | https://download.qt.io/online/qtsdkrepository/all_os/android/qt6_683/qt6_683_x86_64/qt.qt6.683.android_x86_64/ | 16022863 | `15DAF309DA20896F59F66DD8E3EFF6D94C678609081DE0AA3AC5342DFA94A64C` | 错误平台，未使用 |
| `qtdeclarative-Linux-RHEL_8_10-Clang-Android-Android_ANY-X86_64.7z` | 同上 | 25675671 | `8DAB231F15F4BF0223606CD116B8217E9075A2F664B3D1AC25CE5D9D094165AC` | 错误平台，未使用 |
| `qttools-Linux-RHEL_8_10-Clang-Android-Android_ANY-X86_64.7z` | 同上 | 3854402 | `005C50E75A33461A1774E3702D74500290DF62B3ABD0C3AEC867D2218B42B30E` | 错误平台，未使用 |
| `qtsvg-Linux-RHEL_8_10-Clang-Android-Android_ANY-X86_64.7z` | 同上 | 185820 | `AABCE26EA46FC39DA852AFA3155964983306E521088B8E303D694E4623C1CCC7` | 错误平台，未使用 |
| `qttranslations-Linux-RHEL_8_10-Clang-Android-Android_ANY-X86_64.7z` | 同上 | 1783178 | `11775FF08EF15527555FEC5877D0B260C0C3126C4A6037DE7EE5421394222368` | 错误平台，未使用 |

Qt Android prefix 由源码构建生成；构建日志保留在 `build\QtAndroidGcc13`（错误的 win32-g++ 尝试）、`build\QtAndroidClang13`（x86_64）和 `build\QtAndroidArm64Clang13`（arm64-v8a）中。Qt cross-build 的 prefix 没有 `androiddeployqt` 是预期的 host/cross 边界，APK 使用已有 Windows host Qt 的 `androiddeployqt.exe`。

本地 QML smoke 字体不是下载物：`C:\Windows\Fonts\segoeui.ttf` 复制到 `I:\code\Menu\.tools\qt-smoke-fonts`，959752 bytes，SHA256 `8134DBCD09E7B123C9A7F229D49CFFBCB01352CC72EA5E1076B65D0DCA9F73CD`。

## 当前可复用环境

这些不是本项目在 Phase 0 下载的内容，只是盘点到的现有安装/缓存，不能代替 Menu 自己的锁定记录。

| 工具或库 | 版本/状态 | 路径 | 是否已纳入 Menu |
| --- | --- | --- | --- |
| Git | 2.53.0.windows.1 | `I:\Git\bin\git.exe` | 是，主机工具 |
| CMake | 4.3.0-rc1 | `I:\cmake\bin\cmake.exe` | 是，需持续验证 rc 兼容性 |
| Clang/Clang++ | 22.1.4 | `C:\msys64\ucrt64\bin` | 可选 Windows 编译器 |
| GCC/G++ | 16.1.0 | `C:\msys64\ucrt64\bin` | 可选 Windows 编译器 |
| Node/npm | 24.14.1 / 11.11.0 | `C:\Program Files\nodejs` | 管理面板工具链 |
| Qt | 6.8.3 desktop | `C:\Qt\6.8.3\mingw_64`、`C:\Qt\6.8.3\llvm-mingw_64` | 复用，不新增 C 盘安装 |
| SQLite CLI | 3.51.1 | `C:\msys64\ucrt64\bin\sqlite3.exe` | 仅用于诊断；生产链接将锁定依赖 |
| Boost | 1.89.0 源码缓存 | `I:\code\Prism\build\_deps\boost_src-src` | 未复制，Menu 单独锁定 |
| glaze | 现有源码缓存 | `I:\code\Prism\build\_deps\glaze-src` | 未复制，Menu 单独锁定 |
| GoogleTest | 现有源码缓存 | `I:\code\Prism\build\_deps\googletest-src` | 未复制，Menu 单独锁定 |

Qt `6.8.3/mingw_64` 的客户端测试必须使用上表 GCC/G++ 13.1.0 工具链；主机现有 MSYS2 GCC 16.1.0 与该 Qt kit ABI 不一致，不作为 Qt 客户端测试证据。工具链通过当前 PowerShell 的临时 `PATH` 使用，未修改系统环境变量。

## 后续下载登记规则

每次下载必须在本文件追加一行，记录：工具/库名称、精确版本、官方 URL、目标路径、文件大小、SHA256、下载是否使用本进程代理和用途。下载失败时只对当前 PowerShell 进程设置：

```powershell
$env:HTTP_PROXY = 'http://127.0.0.1:7890'
$env:HTTPS_PROXY = 'http://127.0.0.1:7890'
```

下载完成后记录 hash，再继续 configure；不把代理写入系统设置、Git 全局配置或 npm 全局配置。

## 锁定候选

Prism 当前 CMake 中已核对的 Boost 参考包为：

- URL：`https://archives.boost.io/release/1.89.0/source/boost_1_89_0.tar.bz2`
- SHA256：`85a33fa22621b4f314f8e85e1a5e2a9363d22e4f4992925d4bb3bc631b5a0c7a`
- Menu 状态：已下载并由 `I:\code\Menu\.cache\cmake-fetch` 的 configure/build 重新校验；版本与哈希与上表一致。
