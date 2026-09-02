#include "ClientSessionState.hpp"

#include <utility>

namespace Menu::Client {

bool ClientSessionState::IsAuthenticated() const noexcept {
    return !AccessTokenValue.isEmpty();
}

bool ClientSessionState::HasRefreshToken() const noexcept {
    return !RefreshTokenValue.isEmpty();
}

const QString& ClientSessionState::AccessToken() const noexcept {
    return AccessTokenValue;
}

const QString& ClientSessionState::RefreshToken() const noexcept {
    return RefreshTokenValue;
}

const QString& ClientSessionState::DisplayName() const noexcept {
    return DisplayNameValue;
}

const QString& ClientSessionState::UserId() const noexcept {
    return UserIdValue;
}

void ClientSessionState::Set(
    QString AccessTokenInput,
    QString RefreshTokenInput,
    QString DisplayNameInput,
    QString UserIdInput) {
    AccessTokenValue = std::move(AccessTokenInput);
    RefreshTokenValue = std::move(RefreshTokenInput);
    DisplayNameValue = std::move(DisplayNameInput);
    UserIdValue = std::move(UserIdInput);
}

void ClientSessionState::Clear() noexcept {
    AccessTokenValue.clear();
    RefreshTokenValue.clear();
    DisplayNameValue.clear();
    UserIdValue.clear();
}

}  // namespace Menu::Client
