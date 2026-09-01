#pragma once

#include <Application/FeedbackRepository.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

namespace Menu::Infrastructure {

class SqliteFeedbackRepository final : public Application::FeedbackRepository {
public:
    explicit SqliteFeedbackRepository(SqliteDatabase& DatabaseValue)
        : Database(DatabaseValue) {}

    Foundation::Result<Domain::Feedback> Create(
        std::string_view UserId,
        std::string_view RecipeId,
        std::string_view Outcome,
        std::vector<std::string> Tags,
        std::string Comment) override;

private:
    SqliteDatabase& Database;
};

}  // namespace Menu::Infrastructure
