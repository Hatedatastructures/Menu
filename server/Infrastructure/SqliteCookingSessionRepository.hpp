#pragma once

#include <Application/CookingSessionRepository.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

namespace Menu::Infrastructure {

class SqliteCookingSessionRepository final : public Application::CookingSessionRepository {
public:
    explicit SqliteCookingSessionRepository(SqliteDatabase& DatabaseValue)
        : Database(DatabaseValue) {}

    Foundation::Result<Domain::CookingSession> Create(
        std::string_view UserId,
        std::string_view RecipeId) override;

    Foundation::Result<Domain::CookingSession> Update(
        std::string_view UserId,
        std::string_view SessionId,
        int CurrentStepOrder,
        std::string_view State) override;

private:
    Foundation::Result<Domain::CookingSession> Read(
        std::string_view UserId,
        std::string_view SessionId);

    SqliteDatabase& Database;
};

}  // namespace Menu::Infrastructure
