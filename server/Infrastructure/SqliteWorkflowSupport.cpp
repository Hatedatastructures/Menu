#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <utility>

#include <chrono>

namespace Menu::Infrastructure::SqliteWorkflowSupport {

StatementGuard::StatementGuard(sqlite3_stmt* StatementValue)
    : Statement(StatementValue) {}

StatementGuard::StatementGuard(StatementGuard&& Other) noexcept
    : Statement(std::exchange(Other.Statement, nullptr)) {}

StatementGuard& StatementGuard::operator=(StatementGuard&& Other) noexcept {
    if (this != &Other) {
        Reset();
        Statement = std::exchange(Other.Statement, nullptr);
    }
    return *this;
}

StatementGuard::~StatementGuard() {
    Reset();
}

sqlite3_stmt* StatementGuard::Get() const noexcept {
    return Statement;
}

void StatementGuard::Reset() noexcept {
    if (Statement != nullptr) {
        sqlite3_finalize(Statement);
        Statement = nullptr;
    }
}

Foundation::Error StorageError(sqlite3* Handle) {
    const char* Message = Handle == nullptr ? nullptr : sqlite3_errmsg(Handle);
    return Foundation::Error(
        Foundation::ErrorCode::StorageUnavailable,
        Message == nullptr ? "SQLite 查询失败" : std::string(Message));
}

Foundation::Result<StatementGuard> Prepare(
    SqliteDatabase& Database,
    std::string_view Sql) {
    sqlite3_stmt* RawStatement = nullptr;
    const std::string SqlText(Sql);
    if (sqlite3_prepare_v2(
            Database.NativeHandle(), SqlText.c_str(), static_cast<int>(SqlText.size()),
            &RawStatement, nullptr) != SQLITE_OK) {
        return Foundation::Result<StatementGuard>::FromError(
            StorageError(Database.NativeHandle()));
    }
    return StatementGuard(RawStatement);
}

Foundation::Result<void> BindText(
    sqlite3_stmt* Statement,
    int Index,
    std::string_view Value) {
    const std::string Text(Value);
    if (sqlite3_bind_text(Statement, Index, Text.c_str(), -1, SQLITE_TRANSIENT) != SQLITE_OK) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 绑定失败"));
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> BindInteger(
    sqlite3_stmt* Statement,
    int Index,
    std::int64_t Value) {
    if (sqlite3_bind_int64(Statement, Index, Value) != SQLITE_OK) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 绑定失败"));
    }
    return Foundation::Result<void>();
}

std::string ColumnText(sqlite3_stmt* Statement, int Column) {
    const unsigned char* TextValue = sqlite3_column_text(Statement, Column);
    return TextValue == nullptr ? std::string() : reinterpret_cast<const char*>(TextValue);
}

Foundation::Result<std::vector<Domain::RecipeIngredient>> ReadRecipeIngredients(
    SqliteDatabase& Database,
    std::string_view RecipeId,
    int RecipeServings,
    int RequestedServings) {
    auto StatementResult = Prepare(
        Database,
        "SELECT ri.IngredientId, i.Name, i.Category, i.DefaultUnit, "
        "i.IsPantryStaple, ri.Quantity, ri.Unit, ri.ServingFactor, "
        "ri.Preparation, ri.Required "
        "FROM RecipeIngredients ri JOIN Ingredients i "
        "ON i.Id = ri.IngredientId WHERE ri.RecipeId = ? "
        "ORDER BY ri.rowid;");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<std::vector<Domain::RecipeIngredient>>::FromError(
            StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const auto BindResult = BindText(Statement.Get(), 1, RecipeId);
    if (!BindResult.HasValue()) {
        return Foundation::Result<std::vector<Domain::RecipeIngredient>>::FromError(
            BindResult.ErrorValue());
    }
    std::vector<Domain::RecipeIngredient> Ingredients;
    while (true) {
        const int StepResult = sqlite3_step(Statement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<std::vector<Domain::RecipeIngredient>>::FromError(
                StorageError(Database.NativeHandle()));
        }
        Domain::RecipeIngredient Ingredient;
        Ingredient.IngredientId = ColumnText(Statement.Get(), 0);
        Ingredient.IngredientName = ColumnText(Statement.Get(), 1);
        Ingredient.IngredientCategory = ColumnText(Statement.Get(), 2);
        Ingredient.IngredientDefaultUnit = ColumnText(Statement.Get(), 3);
        Ingredient.IngredientIsPantryStaple = sqlite3_column_int(Statement.Get(), 4) != 0;
        const double BaseQuantity = sqlite3_column_double(Statement.Get(), 5);
        Ingredient.Quantity = RecipeServings > 0
            ? BaseQuantity * static_cast<double>(RequestedServings) /
                  static_cast<double>(RecipeServings)
            : BaseQuantity;
        Ingredient.Unit = ColumnText(Statement.Get(), 6);
        Ingredient.ServingFactor = sqlite3_column_double(Statement.Get(), 7);
        Ingredient.Preparation = ColumnText(Statement.Get(), 8);
        Ingredient.Required = sqlite3_column_int(Statement.Get(), 9) != 0;
        Ingredients.push_back(std::move(Ingredient));
    }
    return Ingredients;
}

std::string CurrentTimestamp() {
    using namespace std::chrono;
    return std::to_string(duration_cast<seconds>(
                               system_clock::now().time_since_epoch())
                               .count());
}

}  // namespace Menu::Infrastructure::SqliteWorkflowSupport
