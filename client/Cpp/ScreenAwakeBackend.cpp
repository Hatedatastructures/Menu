#include "ScreenAwakeBackend.hpp"

#include <QtCore/qglobal.h>

#if defined(Q_OS_ANDROID)
#include <QJniObject>
#include <QtCore/qcoreapplication_platform.h>
#endif

namespace Menu::Client {
namespace {

#if defined(Q_OS_ANDROID)

class AndroidScreenAwakeBackend final : public ScreenAwakeBackend {
public:
    void Apply(bool Enabled) override {
        const auto Context = QNativeInterface::QAndroidApplication::context();
        if (Context.object() == nullptr) {
            return;
        }
        QJniObject::callStaticMethod<void>(
            "com/menu/cookflow/SystemUiBridge",
            "setKeepScreenOn",
            "(Landroid/content/Context;Z)V",
            Context.object<jobject>(), static_cast<jboolean>(Enabled));
    }
};

#else

class DesktopScreenAwakeBackend final : public ScreenAwakeBackend {
public:
    void Apply(bool) override {}
};

#endif

}  // namespace

std::unique_ptr<ScreenAwakeBackend> CreateScreenAwakeBackend() {
#if defined(Q_OS_ANDROID)
    return std::make_unique<AndroidScreenAwakeBackend>();
#else
    return std::make_unique<DesktopScreenAwakeBackend>();
#endif
}

}  // namespace Menu::Client
