#include "Application/CookingSessionApplicationService.hpp"

#include <string>
#include <utility>

namespace Menu::Application {
namespace {

Foundation::Error InvalidSession(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

bool IsValidState(std::string_view State) {
    return State == "active" || State == "paused" || State == "completed" ||
        State == "abandoned";
}

}  // namespace

CookingSessionApplicationService::CookingSessionApplicationService(
    std::unique_ptr<CookingSessionRepository> RepositoryValue)
    : Repository(std::move(RepositoryValue)) {}

Foundation::Result<Domain::CookingSession> CookingSessionApplicationService::Create(
    std::string_view UserId,
    std::string_view RecipeId) {
    if (!Repository || UserId.empty() || RecipeId.empty() || RecipeId.size() > 128U) {
        return Foundation::Result<Domain::CookingSession>::FromError(
            InvalidSession("做饭会话信息无效"));
    }
    return Repository->Create(UserId, RecipeId);
}

Foundation::Result<Domain::CookingSession> CookingSessionApplicationService::Update(
    std::string_view UserId,
    std::string_view SessionId,
    int CurrentStepOrder,
    std::string_view State) {
    if (!Repository || UserId.empty() || SessionId.empty() || SessionId.size() > 128U ||
        CurrentStepOrder < 1 || CurrentStepOrder > 128 || !IsValidState(State)) {
        return Foundation::Result<Domain::CookingSession>::FromError(
            InvalidSession("做饭会话更新无效"));
    }
    return Repository->Update(UserId, SessionId, CurrentStepOrder, State);
}

}  // namespace Menu::Application
