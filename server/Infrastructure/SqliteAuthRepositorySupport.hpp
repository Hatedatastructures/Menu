#pragma once

#include <Application/AuthRepository.hpp>
#include <Foundation/Result.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

#include <cstdint>
#include <optional>
#include <string_view>

namespace Menu::Infrastructure::SqliteAuthRepositorySupport {

Foundation::Result<std::optional<Application::AuthUser>> FindUser(
    SqliteDatabase& Database,
    std::string_view Sql,
    std::string_view TokenHash,
    std::int64_t Now,
    bool HasExpiry);

Foundation::Result<void> InsertToken(
    SqliteDatabase& Database,
    std::string_view Table,
    std::string_view UserId,
    std::string_view TokenHash,
    std::int64_t ExpiresAt);

Foundation::Result<void> RollbackWith(
    Foundation::Result<void> Result,
    SqliteDatabase& Database);

}  // namespace Menu::Infrastructure::SqliteAuthRepositorySupport
