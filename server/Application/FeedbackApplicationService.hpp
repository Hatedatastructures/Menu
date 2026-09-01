#pragma once

#include <Application/FeedbackRepository.hpp>

#include <memory>
#include <string_view>

namespace Menu::Application {

class FeedbackApplicationService final {
public:
    explicit FeedbackApplicationService(std::unique_ptr<FeedbackRepository> RepositoryValue);

    Foundation::Result<Domain::Feedback> Create(
        std::string_view UserId,
        std::string_view RecipeId,
        std::string_view Outcome,
        std::vector<std::string> Tags,
        std::string Comment);

private:
    std::unique_ptr<FeedbackRepository> Repository;
};

}  // namespace Menu::Application
