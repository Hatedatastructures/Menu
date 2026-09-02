#pragma once

#include "SystemUiBackend.hpp"

#include <QObject>
#include <QString>

#include <functional>

namespace Menu::Client {

class ThemeController final : public QObject {
    Q_OBJECT
    Q_PROPERTY(QString Mode READ Mode WRITE SetMode NOTIFY ModeChanged)
    Q_PROPERTY(bool Dark READ Dark NOTIFY DarkChanged)

public:
    using BackendFactory = std::function<std::unique_ptr<SystemUiBackend>()>;

    explicit ThemeController(
        QObject* Parent = nullptr,
        BackendFactory Factory = {});
    ~ThemeController() override;

    [[nodiscard]] QString Mode() const;
    [[nodiscard]] bool Dark() const noexcept;

    void SetMode(const QString& ModeValue);
    Q_INVOKABLE void RefreshSystemBars();

signals:
    void ModeChanged();
    void DarkChanged();

private:
    [[nodiscard]] bool SystemDark() const noexcept;
    void Persist() const;
    void ApplySystemBars() const;

    QString ModeValue = QStringLiteral("system");
    std::unique_ptr<SystemUiBackend> Backend;
};

}  // namespace Menu::Client
