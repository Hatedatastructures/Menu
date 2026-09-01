#pragma once

#include <Domain/Feedback.hpp>
#include <Foundation/Result.hpp>

#include <string_view>

namespace Menu::Application {

class FeedbackRepository {
public:
    virtual ~FeedbackRepository() = default;

    virtual Foundation::Result<Domain::Feedback> Create(
        std::string_view UserId,
        std::string_view RecipeId,
        std::string_view Outcome,
        std::vector<std::string> Tags,
        std::string Comment) = 0;
};

}  // namespace Menu::Application
