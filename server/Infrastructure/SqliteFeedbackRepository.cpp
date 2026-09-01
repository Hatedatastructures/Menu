#include "Infrastructure/SqliteFeedbackRepository.hpp"

#include "Infrastructure/SqliteWorkflowSupport.hpp"
#include "Infrastructure/TokenService.hpp"

#include <boost/json/array.hpp>
#include <boost/json/serialize.hpp>

#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

Foundation::Error NotFound(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::NotFound, std::move(Message));
}

Foundation::Result<void> StepDone(SqliteDatabase& Database, sqlite3_stmt* Statement) {
    if (sqlite3_step(Statement) != SQLITE_DONE) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    return Foundation::Result<void>();
}

std::string SerializeTags(const std::vector<std::string>& Tags) {
    boost::json::array Values;
    for (const std::string& Tag : Tags) {
        Values.emplace_back(Tag);
    }
    return boost::json::serialize(Values);
}

}  // namespace

Foundation::Result<Domain::Feedback> SqliteFeedbackRepository::Create(
    std::string_view UserId,
    std::string_view RecipeId,
    std::string_view Outcome,
    std::vector<std::string> Tags,
    std::string Comment) {
    auto RecipeResult = Prepare(
        Database, "SELECT 1 FROM Recipes WHERE Id = ? LIMIT 1;");
    if (!RecipeResult.HasValue()) {
        return Foundation::Result<Domain::Feedback>::FromError(RecipeResult.ErrorValue());
    }
    StatementGuard RecipeStatement = std::move(RecipeResult).Value();
    const auto RecipeBind = BindText(RecipeStatement.Get(), 1, RecipeId);
    if (!RecipeBind.HasValue()) {
        return Foundation::Result<Domain::Feedback>::FromError(RecipeBind.ErrorValue());
    }
    const int RecipeStep = sqlite3_step(RecipeStatement.Get());
    if (RecipeStep == SQLITE_DONE) {
        return Foundation::Result<Domain::Feedback>::FromError(
            NotFound("菜谱不存在"));
    }
    if (RecipeStep != SQLITE_ROW) {
        return Foundation::Result<Domain::Feedback>::FromError(
            StorageError(Database.NativeHandle()));
    }

    const auto Token = TokenService::CreateToken(16U);
    if (!Token.HasValue()) {
        return Foundation::Result<Domain::Feedback>::FromError(Token.ErrorValue());
    }
    const std::string FeedbackId = "feedback." + Token.Value();
    const std::string TagsJson = SerializeTags(Tags);
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::Feedback>::FromError(BeginResult.ErrorValue());
    }
    auto InsertResult = Prepare(
        Database,
        "INSERT INTO Feedback (Id, UserId, RecipeId, Outcome, TagsJson, Comment) "
        "VALUES (?, ?, ?, ?, ?, ?);");
    if (!InsertResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Feedback>::FromError(InsertResult.ErrorValue());
    }
    StatementGuard InsertStatement = std::move(InsertResult).Value();
    const auto IdBind = BindText(InsertStatement.Get(), 1, FeedbackId);
    const auto UserBind = BindText(InsertStatement.Get(), 2, UserId);
    const auto RecipeIdBind = BindText(InsertStatement.Get(), 3, RecipeId);
    const auto OutcomeBind = BindText(InsertStatement.Get(), 4, Outcome);
    const auto TagsBind = BindText(InsertStatement.Get(), 5, TagsJson);
    const auto CommentBind = BindText(InsertStatement.Get(), 6, Comment);
    if (!IdBind.HasValue() || !UserBind.HasValue() || !RecipeIdBind.HasValue() ||
        !OutcomeBind.HasValue() || !TagsBind.HasValue() || !CommentBind.HasValue()) {
        Database.Rollback();
        const Foundation::Error ErrorValue = !IdBind.HasValue()
            ? IdBind.ErrorValue()
            : !UserBind.HasValue() ? UserBind.ErrorValue()
            : !RecipeIdBind.HasValue() ? RecipeIdBind.ErrorValue()
            : !OutcomeBind.HasValue() ? OutcomeBind.ErrorValue()
            : !TagsBind.HasValue() ? TagsBind.ErrorValue() : CommentBind.ErrorValue();
        return Foundation::Result<Domain::Feedback>::FromError(ErrorValue);
    }
    const auto InsertStep = StepDone(Database, InsertStatement.Get());
    if (!InsertStep.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Feedback>::FromError(InsertStep.ErrorValue());
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Feedback>::FromError(CommitResult.ErrorValue());
    }

    Domain::Feedback FeedbackValue;
    FeedbackValue.Id = FeedbackId;
    FeedbackValue.UserId = std::string(UserId);
    FeedbackValue.RecipeId = std::string(RecipeId);
    FeedbackValue.Outcome = std::string(Outcome);
    FeedbackValue.Tags = std::move(Tags);
    FeedbackValue.Comment = std::move(Comment);
    FeedbackValue.CreatedAt = SqliteWorkflowSupport::CurrentTimestamp();
    return FeedbackValue;
}

}  // namespace Menu::Infrastructure
