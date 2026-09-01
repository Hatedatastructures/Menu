#pragma once

#include <Domain/CookingSession.hpp>
#include <Foundation/Result.hpp>

#include <string_view>

namespace Menu::Application {

class CookingSessionRepository {
public:
    virtual ~CookingSessionRepository() = default;

    virtual Foundation::Result<Domain::CookingSession> Create(
        std::string_view UserId,
        std::string_view RecipeId) = 0;

    virtual Foundation::Result<Domain::CookingSession> Update(
        std::string_view UserId,
        std::string_view SessionId,
        int CurrentStepOrder,
        std::string_view State) = 0;
};

}  // namespace Menu::Application
