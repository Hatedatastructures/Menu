#pragma once

#include <string>

namespace Menu::Domain {

struct CookingSession {
    std::string Id;
    std::string UserId;
    std::string RecipeId;
    int CurrentStepOrder = 1;
    std::string State;
    std::string StartedAt;
    std::string UpdatedAt;
};

}  // namespace Menu::Domain
