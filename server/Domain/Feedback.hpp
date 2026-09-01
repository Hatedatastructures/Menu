#pragma once

#include <string>
#include <vector>

namespace Menu::Domain {

struct Feedback {
    std::string Id;
    std::string UserId;
    std::string RecipeId;
    std::string Outcome;
    std::vector<std::string> Tags;
    std::string Comment;
    std::string CreatedAt;
};

}  // namespace Menu::Domain
