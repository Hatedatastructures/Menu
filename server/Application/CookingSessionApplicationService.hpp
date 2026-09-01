#pragma once

#include <Application/CookingSessionRepository.hpp>

#include <memory>
#include <string_view>

namespace Menu::Application {

class CookingSessionApplicationService final {
public:
    explicit CookingSessionApplicationService(
        std::unique_ptr<CookingSessionRepository> RepositoryValue);

    Foundation::Result<Domain::CookingSession> Create(
        std::string_view UserId,
        std::string_view RecipeId);

    Foundation::Result<Domain::CookingSession> Update(
        std::string_view UserId,
        std::string_view SessionId,
        int CurrentStepOrder,
        std::string_view State);

private:
    std::unique_ptr<CookingSessionRepository> Repository;
};

}  // namespace Menu::Application
