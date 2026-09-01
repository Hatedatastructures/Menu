# 视觉与交互验收

更新时间：2026-09-01

## 管理面板

Playwright 使用项目本地 Chromium（`PLAYWRIGHT_BROWSERS_PATH=I:\code\Menu\.tools\playwright-browsers`）和本地服务端测试账号，覆盖登录、菜谱 CRUD、发布、预览、应用内删除确认的取消/确认路径、错误状态和食材 CRUD。截图：

- `docs/evidence-admin-1440x900.png`
- `docs/evidence-admin-1280x800.png`
- `docs/evidence-admin-390x844.png`
- `docs/evidence-admin-playwright-preview.png`

预览中的食材已由 Ingredient catalog enrich 为“番茄”等中文显示字段，测试断言预览不包含 `ingredient.` 原始 ID。移动布局已收紧为单一清晰头部，检查了横向溢出、裁切、低对比度和按钮文字溢出。

## Qt/Android

Qt 桌面 smoke 使用软件渲染、Basic controls、本地字体和 15 秒外部超时；本轮 `docs/evidence-qml-smoke-final.txt` 记录退出码 0，并完整出现 `timer-fired`。QML 测试覆盖导航状态和本地日期边界，QML 警告连接器把运行时 warning 输出为失败证据。客户端缓存 envelope 使用 schema version 2，旧版本会被拒绝。

Android 窗口化 AVD 截图：`docs/evidence-android-final-clean-first-screen.png`；认证、今晚推荐、计划合并清单、做饭页、计时和步骤完成态分别记录在 `docs/evidence-android-final-profile-authenticated.png`、`docs/evidence-android-final-authenticated-tonight.png`、`docs/evidence-android-final-week-merged-two-recipes.png`、`docs/evidence-android-final-cooking.png`、`docs/evidence-android-final-cooking-timer.png` 和 `docs/evidence-android-final-cooking-complete.png`。截图中可见系统状态栏、Menu 页面、真实推荐菜名、本地位图、中文食材、逐项清单和底部导航，没有空白应用内容、重叠或外链图片。设备刷新率为 60Hz，未宣称 120fps。

## 未验证

- Playwright 没有替代 Android 真机触控测试；当前 Android 证据是 x86_64 AVD。
- QML Profiler/Perfetto 的完整首屏与滚动帧时间尚未建立；gfxinfo 仅作为启动样本记录。
- Linux、ARM Android、iOS 和真实 120Hz 设备没有当前主机证据。
