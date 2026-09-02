#include "SystemUiBackend.hpp"

#include <QtCore/qglobal.h>

#if defined(Q_OS_ANDROID)
#include <QJniObject>
#include <QtCore/qcoreapplication_platform.h>
#endif

namespace Menu::Client {
namespace {

#if defined(Q_OS_ANDROID)

class AndroidSystemUiBackend final : public SystemUiBackend {
public:
    void ApplyTheme(bool Dark) override {
        const auto Context = QNativeInterface::QAndroidApplication::context();
        if (Context.object() == nullptr) {
            return;
        }
        QJniObject::callStaticMethod<void>(
            "com/menu/cookflow/SystemUiBridge",
            "applySystemBars",
            "(Landroid/content/Context;Z)V",
            Context.object<jobject>(), static_cast<jboolean>(Dark));
    }
};

#else

class DesktopSystemUiBackend final : public SystemUiBackend {
public:
    void ApplyTheme(bool) override {}
};

#endif

}  // namespace

std::unique_ptr<SystemUiBackend> CreateSystemUiBackend() {
#if defined(Q_OS_ANDROID)
    return std::make_unique<AndroidSystemUiBackend>();
#else
    return std::make_unique<DesktopSystemUiBackend>();
#endif
}

}  // namespace Menu::Client
