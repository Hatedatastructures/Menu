#include "ThemeController.hpp"

#include <QGuiApplication>
#include <QSettings>
#include <QStyleHints>

namespace Menu::Client {

ThemeController::ThemeController(QObject* Parent, BackendFactory Factory)
    : QObject(Parent), Backend(Factory ? Factory() : CreateSystemUiBackend()) {
    const QSettings Settings;
    const QString StoredMode = Settings.value(
        QStringLiteral("theme/mode"), ModeValue).toString();
    if (StoredMode == QStringLiteral("light") || StoredMode == QStringLiteral("dark")) {
        ModeValue = StoredMode;
    }
    if (auto* Hints = QGuiApplication::styleHints()) {
        connect(Hints, &QStyleHints::colorSchemeChanged, this,
            [this](Qt::ColorScheme) {
                if (ModeValue == QStringLiteral("system")) {
                    ApplySystemBars();
                    emit DarkChanged();
                }
            });
    }
    ApplySystemBars();
}

ThemeController::~ThemeController() = default;

QString ThemeController::Mode() const {
    return ModeValue;
}

bool ThemeController::Dark() const noexcept {
    return ModeValue == QStringLiteral("dark") ||
        (ModeValue == QStringLiteral("system") && SystemDark());
}

void ThemeController::SetMode(const QString& ModeValueInput) {
    const QString NormalizedMode =
        ModeValueInput == QStringLiteral("light") || ModeValueInput == QStringLiteral("dark")
        ? ModeValueInput : QStringLiteral("system");
    if (ModeValue == NormalizedMode) {
        return;
    }
    const bool PreviousDark = Dark();
    ModeValue = NormalizedMode;
    Persist();
    emit ModeChanged();
    if (PreviousDark != Dark()) {
        ApplySystemBars();
        emit DarkChanged();
    }
}

void ThemeController::RefreshSystemBars() {
    ApplySystemBars();
}

bool ThemeController::SystemDark() const noexcept {
    const auto* Hints = QGuiApplication::styleHints();
    return Hints != nullptr && Hints->colorScheme() == Qt::ColorScheme::Dark;
}

void ThemeController::Persist() const {
    QSettings Settings;
    Settings.setValue(QStringLiteral("theme/mode"), ModeValue);
}

void ThemeController::ApplySystemBars() const {
    if (Backend) {
        Backend->ApplyTheme(Dark());
    }
}

}  // namespace Menu::Client
