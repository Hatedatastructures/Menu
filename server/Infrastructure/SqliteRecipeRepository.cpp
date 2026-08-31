#include "Infrastructure/SqliteRecipeRepository.hpp"

#include <boost/json/parse.hpp>
#include <boost/json/serialize.hpp>

#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure {
namespace {

class StatementGuard final {
public:
    StatementGuard() = default;

    explicit StatementGuard(sqlite3_stmt* StatementValue)
        : Statement(StatementValue) {}

    StatementGuard(const StatementGuard&) = delete;
    StatementGuard& operator=(const StatementGuard&) = delete;

    StatementGuard(StatementGuard&& Other) noexcept
        : Statement(std::exchange(Other.Statement, nullptr)) {}

    StatementGuard& operator=(StatementGuard&& Other) noexcept {
        if (this != &Other) {
            Reset();
            Statement = std::exchange(Other.Statement, nullptr);
        }
        return *this;
    }

    ~StatementGuard() {
        Reset();
    }

    [[nodiscard]] sqlite3_stmt* Get() const noexcept {
        return Statement;
    }

    void Reset() noexcept {
        if (Statement != nullptr) {
            sqlite3_finalize(Statement);
            Statement = nullptr;
        }
    }

private:
    sqlite3_stmt* Statement = nullptr;
};

Foundation::Error StorageError(sqlite3* Handle) {
    const char* Message = Handle == nullptr ? nullptr : sqlite3_errmsg(Handle);
    return Foundation::Error(
        Foundation::ErrorCode::StorageUnavailable,
        Message == nullptr ? "SQLite 查询失败" : std::string(Message));
}

Foundation::Result<StatementGuard> Prepare(SqliteDatabase& Database, std::string_view Sql) {
    sqlite3_stmt* RawStatement = nullptr;
    const std::string SqlText(Sql);
    const int ResultCode = sqlite3_prepare_v2(
        Database.NativeHandle(), SqlText.c_str(), static_cast<int>(SqlText.size()),
        &RawStatement, nullptr);
    if (ResultCode != SQLITE_OK) {
        return Foundation::Result<StatementGuard>::FromError(
            StorageError(Database.NativeHandle()));
    }
    return StatementGuard(RawStatement);
}

std::string ColumnText(sqlite3_stmt* Statement, int Column) {
    const unsigned char* TextValue = sqlite3_column_text(Statement, Column);
    return TextValue == nullptr ? std::string() : reinterpret_cast<const char*>(TextValue);
}

Foundation::Result<std::vector<std::string>> ParseStringArray(
    sqlite3* Handle,
    std::string_view JsonText) {
    boost::system::error_code Error;
    const boost::json::value Parsed = boost::json::parse(JsonText, Error);
    if (Error || !Parsed.is_array()) {
        return Foundation::Result<std::vector<std::string>>::FromError(StorageError(Handle));
    }

    std::vector<std::string> Values;
    for (const boost::json::value& Value : Parsed.as_array()) {
        if (!Value.is_string()) {
            return Foundation::Result<std::vector<std::string>>::FromError(StorageError(Handle));
        }
        Values.emplace_back(Value.as_string().c_str());
    }
    return Values;
}

Foundation::Result<void> BindText(sqlite3_stmt* Statement, int Index, std::string_view Value) {
    const std::string Text(Value);
    if (sqlite3_bind_text(Statement, Index, Text.c_str(), -1, SQLITE_TRANSIENT) != SQLITE_OK) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 绑定失败"));
    }
    return Foundation::Result<void>();
}

Foundation::Result<Domain::Recipe> ReadRecipeRow(
    SqliteDatabase& Database,
    sqlite3_stmt* Statement) {
    Domain::Recipe RecipeValue;
    RecipeValue.Id = ColumnText(Statement, 0);
    RecipeValue.Slug = ColumnText(Statement, 1);
    RecipeValue.Name = ColumnText(Statement, 2);
    RecipeValue.Cuisine = ColumnText(Statement, 3);
    RecipeValue.Description = ColumnText(Statement, 4);
    RecipeValue.PrepMinutes = sqlite3_column_int(Statement, 5);
    RecipeValue.CookMinutes = sqlite3_column_int(Statement, 6);
    RecipeValue.Servings = sqlite3_column_int(Statement, 7);
    RecipeValue.Difficulty = sqlite3_column_int(Statement, 8);
    RecipeValue.ImagePath = ColumnText(Statement, 9);
    RecipeValue.Status = ColumnText(Statement, 10);

    const auto Allergens = ParseStringArray(Database.NativeHandle(), ColumnText(Statement, 11));
    if (!Allergens.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(Allergens.ErrorValue());
    }
    RecipeValue.Allergens = Allergens.Value();

    const auto Cookware = ParseStringArray(Database.NativeHandle(), ColumnText(Statement, 12));
    if (!Cookware.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(Cookware.ErrorValue());
    }
    RecipeValue.Cookware = Cookware.Value();

    auto IngredientStatementResult = Prepare(
        Database,
        "SELECT RecipeIngredients.IngredientId, RecipeIngredients.Quantity, "
        "RecipeIngredients.Unit, RecipeIngredients.Required, RecipeIngredients.ServingFactor, "
        "RecipeIngredients.Preparation, Ingredients.Name, Ingredients.Category, "
        "Ingredients.DefaultUnit, Ingredients.IsPantryStaple "
        "FROM RecipeIngredients JOIN Ingredients ON Ingredients.Id = RecipeIngredients.IngredientId "
        "WHERE RecipeIngredients.RecipeId = ? ORDER BY RecipeIngredients.IngredientId;");
    if (!IngredientStatementResult.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(
            IngredientStatementResult.ErrorValue());
    }
    StatementGuard IngredientStatement = std::move(IngredientStatementResult).Value();
    const auto IngredientBind = BindText(IngredientStatement.Get(), 1, RecipeValue.Id);
    if (!IngredientBind.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(IngredientBind.ErrorValue());
    }
    while (true) {
        const int StepResult = sqlite3_step(IngredientStatement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<Domain::Recipe>::FromError(
                StorageError(Database.NativeHandle()));
        }
        RecipeValue.Ingredients.emplace_back(
            ColumnText(IngredientStatement.Get(), 0),
            sqlite3_column_double(IngredientStatement.Get(), 1),
            ColumnText(IngredientStatement.Get(), 2),
            sqlite3_column_int(IngredientStatement.Get(), 3) != 0,
            sqlite3_column_double(IngredientStatement.Get(), 4),
            ColumnText(IngredientStatement.Get(), 5));
        Domain::RecipeIngredient& IngredientValue = RecipeValue.Ingredients.back();
        IngredientValue.IngredientName = ColumnText(IngredientStatement.Get(), 6);
        IngredientValue.IngredientCategory = ColumnText(IngredientStatement.Get(), 7);
        IngredientValue.IngredientDefaultUnit = ColumnText(IngredientStatement.Get(), 8);
        IngredientValue.IngredientIsPantryStaple =
            sqlite3_column_int(IngredientStatement.Get(), 9) != 0;
    }

    auto StepStatementResult = Prepare(
        Database,
        "SELECT StepOrder, Title, Instruction, DurationSeconds, HasTimer "
        "FROM RecipeSteps WHERE RecipeId = ? ORDER BY StepOrder;");
    if (!StepStatementResult.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(StepStatementResult.ErrorValue());
    }
    StatementGuard StepStatement = std::move(StepStatementResult).Value();
    const auto StepBind = BindText(StepStatement.Get(), 1, RecipeValue.Id);
    if (!StepBind.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(StepBind.ErrorValue());
    }
    while (true) {
        const int StepResult = sqlite3_step(StepStatement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<Domain::Recipe>::FromError(
                StorageError(Database.NativeHandle()));
        }
        RecipeValue.Steps.push_back(Domain::RecipeStep{
            sqlite3_column_int(StepStatement.Get(), 0),
            ColumnText(StepStatement.Get(), 1),
            ColumnText(StepStatement.Get(), 2),
            sqlite3_column_int(StepStatement.Get(), 3),
            sqlite3_column_int(StepStatement.Get(), 4) != 0});
    }

    return RecipeValue;
}

Foundation::Result<Domain::Recipe> ReadPublishedRecipe(
    SqliteDatabase& Database,
    std::string_view Id) {
    auto StatementResult = Prepare(
        Database,
        "SELECT Id, Slug, Name, Cuisine, Description, PrepMinutes, CookMinutes, Servings, "
        "Difficulty, ImagePath, Status, AllergensJson, CookwareJson FROM Recipes "
        "WHERE Id = ? AND Status = 'published';");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const auto BindResult = BindText(Statement.Get(), 1, Id);
    if (!BindResult.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(BindResult.ErrorValue());
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult != SQLITE_ROW) {
        return Foundation::Result<Domain::Recipe>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "菜谱不存在"));
    }
    return ReadRecipeRow(Database, Statement.Get());
}

}  // namespace

Foundation::Result<std::vector<Domain::Recipe>> SqliteRecipeRepository::ListPublished() {
    auto StatementResult = Prepare(
        Database,
        "SELECT Id, Slug, Name, Cuisine, Description, PrepMinutes, CookMinutes, Servings, "
        "Difficulty, ImagePath, Status, AllergensJson, CookwareJson FROM Recipes "
        "WHERE Status = 'published' ORDER BY Id;");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<std::vector<Domain::Recipe>>::FromError(
            StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    std::vector<Domain::Recipe> Recipes;
    while (true) {
        const int StepResult = sqlite3_step(Statement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<std::vector<Domain::Recipe>>::FromError(
                StorageError(Database.NativeHandle()));
        }
        auto RecipeResult = ReadRecipeRow(Database, Statement.Get());
        if (!RecipeResult.HasValue()) {
            return Foundation::Result<std::vector<Domain::Recipe>>::FromError(
                RecipeResult.ErrorValue());
        }
        Recipes.push_back(std::move(RecipeResult).Value());
    }
    return Recipes;
}

Foundation::Result<std::optional<Domain::Recipe>>
SqliteRecipeRepository::FindPublishedById(std::string_view Id) {
    auto RecipeResult = ReadPublishedRecipe(Database, Id);
    if (!RecipeResult.HasValue()) {
        if (RecipeResult.ErrorValue().CodeValue() == Foundation::ErrorCode::NotFound) {
            return std::optional<Domain::Recipe>();
        }
        return Foundation::Result<std::optional<Domain::Recipe>>::FromError(
            RecipeResult.ErrorValue());
    }
    return std::optional<Domain::Recipe>(std::move(RecipeResult).Value());
}

}  // namespace Menu::Infrastructure
