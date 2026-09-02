#pragma once

#include <QObject>
#include <QPointer>

class QScreen;
class QWindow;

namespace Menu::Client {

constexpr qreal FallbackRefreshRateHz = 120.0;

[[nodiscard]] qreal EffectiveRefreshRate(qreal ReportedRefreshRate) noexcept;
[[nodiscard]] qreal FrameBudgetMilliseconds(qreal RefreshRate) noexcept;

class DisplayMetrics final : public QObject {
    Q_OBJECT
    Q_PROPERTY(qreal RefreshRateHz READ RefreshRateHz NOTIFY Changed)
    Q_PROPERTY(qreal ReportedRefreshRateHz READ ReportedRefreshRateHz NOTIFY Changed)
    Q_PROPERTY(qreal FrameBudgetMilliseconds READ FrameBudgetMilliseconds NOTIFY Changed)
    Q_PROPERTY(qreal TopInset READ TopInset NOTIFY Changed)
    Q_PROPERTY(bool UsingFallback READ UsingFallback NOTIFY Changed)

public:
    explicit DisplayMetrics(QObject* Parent = nullptr);

    [[nodiscard]] qreal RefreshRateHz() const noexcept;
    [[nodiscard]] qreal ReportedRefreshRateHz() const noexcept;
    [[nodiscard]] qreal FrameBudgetMilliseconds() const noexcept;
    [[nodiscard]] qreal TopInset() const noexcept;
    [[nodiscard]] bool UsingFallback() const noexcept;

    void AttachWindow(QWindow* Window);

public slots:
    void Refresh();

signals:
    void Changed();

private:
    QPointer<QWindow> WindowValue;
    QPointer<QScreen> ScreenValue;
    qreal RefreshRateValue = FallbackRefreshRateHz;
    qreal ReportedRefreshRateValue = 0.0;
    qreal TopInsetValue = 0.0;
    bool UsingFallbackValue = true;
};

}  // namespace Menu::Client
