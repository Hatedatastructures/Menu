#include "DisplayMetrics.hpp"

#include <QGuiApplication>
#include <QScreen>
#include <QWindow>

#include <cmath>

#if defined(Q_OS_ANDROID)
#include <QJniObject>
#include <QtCore/qcoreapplication_platform.h>
#endif

namespace Menu::Client {

qreal EffectiveRefreshRate(qreal ReportedRefreshRate) noexcept {
    return std::isfinite(ReportedRefreshRate) && ReportedRefreshRate > 1.0
        ? ReportedRefreshRate
        : FallbackRefreshRateHz;
}

qreal FrameBudgetMilliseconds(qreal RefreshRate) noexcept {
    return 1000.0 / EffectiveRefreshRate(RefreshRate);
}

DisplayMetrics::DisplayMetrics(QObject* Parent)
    : QObject(Parent) {
    Refresh();
}

qreal DisplayMetrics::RefreshRateHz() const noexcept {
    return RefreshRateValue;
}

qreal DisplayMetrics::ReportedRefreshRateHz() const noexcept {
    return ReportedRefreshRateValue;
}

qreal DisplayMetrics::FrameBudgetMilliseconds() const noexcept {
    return Menu::Client::FrameBudgetMilliseconds(RefreshRateValue);
}

qreal DisplayMetrics::TopInset() const noexcept {
    return TopInsetValue;
}

bool DisplayMetrics::UsingFallback() const noexcept {
    return UsingFallbackValue;
}

void DisplayMetrics::AttachWindow(QWindow* Window) {
    if (WindowValue == Window) {
        Refresh();
        return;
    }
    if (WindowValue != nullptr) {
        disconnect(WindowValue, &QWindow::screenChanged,
                   this, &DisplayMetrics::Refresh);
    }
    WindowValue = Window;
    if (WindowValue != nullptr) {
        connect(WindowValue, &QWindow::screenChanged,
                this, &DisplayMetrics::Refresh);
    }
    Refresh();
}

void DisplayMetrics::Refresh() {
    QScreen* NextScreen = WindowValue != nullptr && WindowValue->screen() != nullptr
        ? WindowValue->screen()
        : QGuiApplication::primaryScreen();
    if (ScreenValue != NextScreen) {
        if (ScreenValue != nullptr) {
            disconnect(ScreenValue, &QScreen::refreshRateChanged,
                       this, &DisplayMetrics::Refresh);
        }
        ScreenValue = NextScreen;
        if (ScreenValue != nullptr) {
            connect(ScreenValue, &QScreen::refreshRateChanged,
                    this, &DisplayMetrics::Refresh);
        }
    }

    const qreal ReportedRefreshRate = ScreenValue == nullptr
        ? 0.0
        : ScreenValue->refreshRate();
    const qreal NextRefreshRate = EffectiveRefreshRate(ReportedRefreshRate);
    const bool NextUsingFallback = NextRefreshRate == FallbackRefreshRateHz &&
        (!std::isfinite(ReportedRefreshRate) || ReportedRefreshRate <= 1.0);
#if defined(Q_OS_ANDROID)
    qreal NextTopInset = 0.0;
    const auto Context = QNativeInterface::QAndroidApplication::context();
    if (Context.object() != nullptr) {
        const jint Pixels = QJniObject::callStaticMethod<jint>(
            "com/menu/cookflow/SystemUiBridge",
            "statusBarHeight",
            "(Landroid/content/Context;)I",
            Context.object<jobject>());
        const qreal DevicePixelRatio = WindowValue != nullptr
            ? WindowValue->devicePixelRatio()
            : (ScreenValue != nullptr ? ScreenValue->devicePixelRatio() : 1.0);
        if (Pixels > 0 && DevicePixelRatio > 0.0) {
            NextTopInset = static_cast<qreal>(Pixels) / DevicePixelRatio;
        }
    }
#else
    const qreal NextTopInset = 0.0;
#endif
    if (RefreshRateValue == NextRefreshRate &&
        ReportedRefreshRateValue == ReportedRefreshRate &&
        UsingFallbackValue == NextUsingFallback &&
        qFuzzyCompare(TopInsetValue + 1.0, NextTopInset + 1.0)) {
        return;
    }
    RefreshRateValue = NextRefreshRate;
    ReportedRefreshRateValue = ReportedRefreshRate;
    UsingFallbackValue = NextUsingFallback;
    TopInsetValue = NextTopInset;
    emit Changed();
}

}  // namespace Menu::Client
