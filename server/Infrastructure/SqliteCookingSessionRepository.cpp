#include "Infrastructure/SqliteCookingSessionRepository.hpp"

#include "Infrastructure/SqliteWorkflowSupport.hpp"
#include "Infrastructure/TokenService.hpp"

#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::ColumnText;
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

}  // namespace

Foundation::Result<Domain::CookingSession> SqliteCookingSessionRepository::Create(
    std::string_view UserId,
    std::string_view RecipeId) {
    auto RecipeResult = Prepare(
        Database, "SELECT 1 FROM Recipes WHERE Id = ? LIMIT 1;");
    if (!RecipeResult.HasValue()) {
        return Foundation::Result<Domain::CookingSession>::FromError(
            RecipeResult.ErrorValue());
    }
    StatementGuard RecipeStatement = std::move(RecipeResult).Value();
    const auto RecipeBind = BindText(RecipeStatement.Get(), 1, RecipeId);
    if (!RecipeBind.HasValue()) {
        return Foundation::Result<Domain::CookingSession>::FromError(RecipeBind.ErrorValue());
    }
    const int RecipeStep = sqlite3_step(RecipeStatement.Get());
    if (RecipeStep == SQLITE_DONE) {
        return Foundation::Result<Domain::CookingSession>::FromError(
            NotFound("菜谱不存在"));
    }
    if (RecipeStep != SQLITE_ROW) {
        return Foundation::Result<Domain::CookingSession>::FromError(
            StorageError(Database.NativeHandle()));
    }

    const auto Token = TokenService::CreateToken(16U);
    if (!Token.HasValue()) {
        return Foundation::Result<Domain::CookingSession>::FromError(Token.ErrorValue());
    }
    const std::string SessionId = "session." + Token.Value();
    const std::string Timestamp = SqliteWorkflowSupport::CurrentTimestamp();
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::CookingSession>::FromError(BeginResult.ErrorValue());
    }
    auto InsertResult = Prepare(
        Database,
        "INSERT INTO CookingSessions "
        "(Id, UserId, RecipeId, CurrentStepOrder, State, StartedAt, UpdatedAt) "
        "VALUES (?, ?, ?, 1, 'active', ?, ?);");
    if (!InsertResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::CookingSession>::FromError(
            InsertResult.ErrorValue());
    }
    StatementGuard InsertStatement = std::move(InsertResult).Value();
    const auto IdBind = BindText(InsertStatement.Get(), 1, SessionId);
    const auto UserBind = BindText(InsertStatement.Get(), 2, UserId);
    const auto RecipeIdBind = BindText(InsertStatement.Get(), 3, RecipeId);
    const auto StartedBind = BindText(InsertStatement.Get(), 4, Timestamp);
    const auto UpdatedBind = BindText(InsertStatement.Get(), 5, Timestamp);
    if (!IdBind.HasValue() || !UserBind.HasValue() || !RecipeIdBind.HasValue() ||
        !StartedBind.HasValue() || !UpdatedBind.HasValue()) {
        Database.Rollback();
        const Foundation::Error ErrorValue = !IdBind.HasValue()
            ? IdBind.ErrorValue()
            : !UserBind.HasValue() ? UserBind.ErrorValue()
            : !RecipeIdBind.HasValue() ? RecipeIdBind.ErrorValue()
            : !StartedBind.HasValue() ? StartedBind.ErrorValue()
                                      : UpdatedBind.ErrorValue();
        return Foundation::Result<Domain::CookingSession>::FromError(ErrorValue);
    }
    const auto InsertStep = StepDone(Database, InsertStatement.Get());
    if (!InsertStep.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::CookingSession>::FromError(InsertStep.ErrorValue());
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::CookingSession>::FromError(CommitResult.ErrorValue());
    }
    return Read(UserId, SessionId);
}

Foundation::Result<Domain::CookingSession> SqliteCookingSessionRepository::Update(
    std::string_view UserId,
    std::string_view SessionId,
    int CurrentStepOrder,
    std::string_view State) {
    auto UpdateResult = Prepare(
        Database,
        "UPDATE CookingSessions SET CurrentStepOrder = ?, State = ?, UpdatedAt = ? "
        "WHERE Id = ? AND UserId = ?;");
    if (!UpdateResult.HasValue()) {
        return Foundation::Result<Domain::CookingSession>::FromError(
            UpdateResult.ErrorValue());
    }
    StatementGuard Statement = std::move(UpdateResult).Value();
    const auto StepBind = SqliteWorkflowSupport::BindInteger(
        Statement.Get(), 1, CurrentStepOrder);
    const auto StateBind = BindText(Statement.Get(), 2, State);
    const std::string Timestamp = SqliteWorkflowSupport::CurrentTimestamp();
    const auto TimeBind = BindText(Statement.Get(), 3, Timestamp);
    const auto IdBind = BindText(Statement.Get(), 4, SessionId);
    const auto UserBind = BindText(Statement.Get(), 5, UserId);
    if (!StepBind.HasValue() || !StateBind.HasValue() || !TimeBind.HasValue() ||
        !IdBind.HasValue() || !UserBind.HasValue()) {
        const Foundation::Error ErrorValue = !StepBind.HasValue()
            ? StepBind.ErrorValue()
            : !StateBind.HasValue() ? StateBind.ErrorValue()
            : !TimeBind.HasValue() ? TimeBind.ErrorValue()
            : !IdBind.HasValue() ? IdBind.ErrorValue() : UserBind.ErrorValue();
        return Foundation::Result<Domain::CookingSession>::FromError(ErrorValue);
    }
    const auto StepResult = StepDone(Database, Statement.Get());
    if (!StepResult.HasValue()) {
        return Foundation::Result<Domain::CookingSession>::FromError(StepResult.ErrorValue());
    }
    if (sqlite3_changes(Database.NativeHandle()) == 0) {
        return Foundation::Result<Domain::CookingSession>::FromError(
            NotFound("做饭会话不存在"));
    }
    return Read(UserId, SessionId);
}

Foundation::Result<Domain::CookingSession> SqliteCookingSessionRepository::Read(
    std::string_view UserId,
    std::string_view SessionId) {
    auto Result = Prepare(
        Database,
        "SELECT Id, RecipeId, CurrentStepOrder, State, StartedAt, UpdatedAt "
        "FROM CookingSessions WHERE Id = ? AND UserId = ? LIMIT 1;");
    if (!Result.HasValue()) {
        return Foundation::Result<Domain::CookingSession>::FromError(Result.ErrorValue());
    }
    StatementGuard Statement = std::move(Result).Value();
    const auto IdBind = BindText(Statement.Get(), 1, SessionId);
    const auto UserBind = BindText(Statement.Get(), 2, UserId);
    if (!IdBind.HasValue() || !UserBind.HasValue()) {
        return Foundation::Result<Domain::CookingSession>::FromError(
            !IdBind.HasValue() ? IdBind.ErrorValue() : UserBind.ErrorValue());
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult == SQLITE_DONE) {
        return Foundation::Result<Domain::CookingSession>::FromError(
            NotFound("做饭会话不存在"));
    }
    if (StepResult != SQLITE_ROW) {
        return Foundation::Result<Domain::CookingSession>::FromError(
            StorageError(Database.NativeHandle()));
    }
    Domain::CookingSession Session;
    Session.Id = ColumnText(Statement.Get(), 0);
    Session.UserId = std::string(UserId);
    Session.RecipeId = ColumnText(Statement.Get(), 1);
    Session.CurrentStepOrder = sqlite3_column_int(Statement.Get(), 2);
    Session.State = ColumnText(Statement.Get(), 3);
    Session.StartedAt = ColumnText(Statement.Get(), 4);
    Session.UpdatedAt = ColumnText(Statement.Get(), 5);
    return Session;
}

}  // namespace Menu::Infrastructure
