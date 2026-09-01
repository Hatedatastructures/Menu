#include "Application/FeedbackApplicationService.hpp"

#include <utility>

namespace Menu::Application {
namespace {

Foundation::Error InvalidFeedback(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

}  // namespace

FeedbackApplicationService::FeedbackApplicationService(
    std::unique_ptr<FeedbackRepository> RepositoryValue)
    : Repository(std::move(RepositoryValue)) {}

Foundation::Result<Domain::Feedback> FeedbackApplicationService::Create(
    std::string_view UserId,
    std::string_view RecipeId,
    std::string_view Outcome,
    std::vector<std::string> Tags,
    std::string Comment) {
    if (!Repository || UserId.empty() || RecipeId.empty() || RecipeId.size() > 128U ||
        (Outcome != "made" && Outcome != "skipped") || Tags.size() > 8U ||
        Comment.size() > 2000U) {
        return Foundation::Result<Domain::Feedback>::FromError(
            InvalidFeedback("反馈信息无效"));
    }
    for (const std::string& Tag : Tags) {
        if (Tag.empty() || Tag.size() > 64U) {
            return Foundation::Result<Domain::Feedback>::FromError(
                InvalidFeedback("反馈标签无效"));
        }
    }
    return Repository->Create(
        UserId, RecipeId, Outcome, std::move(Tags), std::move(Comment));
}

}  // namespace Menu::Application
