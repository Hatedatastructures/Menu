#pragma once

#include <QString>

namespace Menu::Client {

class ClientSessionState final {
public:
    [[nodiscard]] bool IsAuthenticated() const noexcept;
    [[nodiscard]] bool HasRefreshToken() const noexcept;
    [[nodiscard]] const QString& AccessToken() const noexcept;
    [[nodiscard]] const QString& RefreshToken() const noexcept;
    [[nodiscard]] const QString& DisplayName() const noexcept;
    [[nodiscard]] const QString& UserId() const noexcept;

    void Set(
        QString AccessTokenValue,
        QString RefreshTokenValue,
        QString DisplayNameValue,
        QString UserIdValue);
    void Clear() noexcept;

private:
    QString AccessTokenValue;
    QString RefreshTokenValue;
    QString DisplayNameValue;
    QString UserIdValue;
};

}  // namespace Menu::Client
