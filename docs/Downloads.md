# 下载、版本与缓存登记

更新时间：2026-09-01

## Phase 0 下载记录

Phase 0 没有下载新工具或库，没有使用代理，也没有修改全局代理配置。Phase 1 fresh configure 下载了下面两项，未使用代理；下载发生在当前 PowerShell 进程，未修改全局配置。下面的路径是后续唯一允许的项目缓存位置：

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
- Menu 状态：尚未下载到 Menu 缓存，Phase 1 configure 时重新校验。
