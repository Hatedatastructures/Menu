#include "Infrastructure/SqliteRecipeRepository.hpp"

#include <boost/json/array.hpp>
#include <boost/json/parse.hpp>
#include <boost/json/serialize.hpp>

#include <array>
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

Foundation::Result<void> BindInteger(sqlite3_stmt* Statement, int Index, int Value) {
    if (sqlite3_bind_int(Statement, Index, Value) != SQLITE_OK) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 绑定失败"));
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> BindDouble(sqlite3_stmt* Statement, int Index, double Value) {
    if (sqlite3_bind_double(Statement, Index, Value) != SQLITE_OK) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 绑定失败"));
    }
    return Foundation::Result<void>();
}

std::string SerializeStrings(const std::vector<std::string>& Values) {
    boost::json::array Array;
    for (const std::string& Value : Values) {
        Array.emplace_back(Value);
    }
    return boost::json::serialize(Array);
}

Foundation::Result<void> StepDone(SqliteDatabase& Database, sqlite3_stmt* Statement) {
    if (sqlite3_step(Statement) != SQLITE_DONE) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> InsertIngredientRows(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue) {
    for (const Domain::RecipeIngredient& IngredientValue : RecipeValue.Ingredients) {
        auto StatementResult = Prepare(
            Database,
            "INSERT INTO RecipeIngredients "
            "(RecipeId, IngredientId, Quantity, Unit, ServingFactor, Preparation, Required) "
            "VALUES (?, ?, ?, ?, ?, ?, ?);");
        if (!StatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
        }
        StatementGuard Statement = std::move(StatementResult).Value();
        const auto RecipeBind = BindText(Statement.Get(), 1, RecipeValue.Id);
        const auto IngredientBind = BindText(Statement.Get(), 2, IngredientValue.IngredientId);
        const auto QuantityBind = BindDouble(Statement.Get(), 3, IngredientValue.Quantity);
        const auto UnitBind = BindText(Statement.Get(), 4, IngredientValue.Unit);
        const auto FactorBind = BindDouble(Statement.Get(), 5, IngredientValue.ServingFactor);
        const auto PreparationBind = BindText(Statement.Get(), 6, IngredientValue.Preparation);
        const auto RequiredBind = BindInteger(
            Statement.Get(), 7, IngredientValue.Required ? 1 : 0);
        if (!RecipeBind.HasValue() || !IngredientBind.HasValue() || !QuantityBind.HasValue() ||
            !UnitBind.HasValue() || !FactorBind.HasValue() || !PreparationBind.HasValue() ||
            !RequiredBind.HasValue()) {
            const Foundation::Error ErrorValue = !RecipeBind.HasValue()
                                                     ? RecipeBind.ErrorValue()
                                                     : !IngredientBind.HasValue()
                                                           ? IngredientBind.ErrorValue()
                                                           : !QuantityBind.HasValue()
                                                                 ? QuantityBind.ErrorValue()
                                                                 : !UnitBind.HasValue()
                                                                       ? UnitBind.ErrorValue()
                                                                       : !FactorBind.HasValue()
                                                                             ? FactorBind.ErrorValue()
                                                                             : !PreparationBind.HasValue()
                                                                                   ? PreparationBind.ErrorValue()
                                                                                   : RequiredBind.ErrorValue();
            return Foundation::Result<void>::FromError(ErrorValue);
        }
        const auto Result = StepDone(Database, Statement.Get());
        if (!Result.HasValue()) {
            return Result;
        }
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> InsertStepRows(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue) {
    for (const Domain::RecipeStep& StepValue : RecipeValue.Steps) {
        auto StatementResult = Prepare(
            Database,
            "INSERT INTO RecipeSteps "
            "(RecipeId, StepOrder, Title, Instruction, DurationSeconds, HasTimer) "
            "VALUES (?, ?, ?, ?, ?, ?);");
        if (!StatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
        }
        StatementGuard Statement = std::move(StatementResult).Value();
        const auto RecipeBind = BindText(Statement.Get(), 1, RecipeValue.Id);
        const auto OrderBind = BindInteger(Statement.Get(), 2, StepValue.StepOrder);
        const auto TitleBind = BindText(Statement.Get(), 3, StepValue.Title);
        const auto InstructionBind = BindText(Statement.Get(), 4, StepValue.Instruction);
        const auto DurationBind = BindInteger(Statement.Get(), 5, StepValue.DurationSeconds);
        const auto TimerBind = BindInteger(Statement.Get(), 6, StepValue.HasTimer ? 1 : 0);
        if (!RecipeBind.HasValue() || !OrderBind.HasValue() || !TitleBind.HasValue() ||
            !InstructionBind.HasValue() || !DurationBind.HasValue() || !TimerBind.HasValue()) {
            const Foundation::Error ErrorValue = !RecipeBind.HasValue()
                                                     ? RecipeBind.ErrorValue()
                                                     : !OrderBind.HasValue()
                                                           ? OrderBind.ErrorValue()
                                                           : !TitleBind.HasValue()
                                                                 ? TitleBind.ErrorValue()
                                                                 : !InstructionBind.HasValue()
                                                                       ? InstructionBind.ErrorValue()
                                                                       : !DurationBind.HasValue()
                                                                             ? DurationBind.ErrorValue()
                                                                             : TimerBind.ErrorValue();
            return Foundation::Result<void>::FromError(ErrorValue);
        }
        const auto Result = StepDone(Database, Statement.Get());
        if (!Result.HasValue()) {
            return Result;
        }
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> DeleteRecipeChildren(
    SqliteDatabase& Database,
    std::string_view RecipeId) {
    for (const char* Sql : {
             "DELETE FROM RecipeIngredients WHERE RecipeId = ?;",
             "DELETE FROM RecipeSteps WHERE RecipeId = ?;"}) {
        auto StatementResult = Prepare(Database, Sql);
        if (!StatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
        }
        StatementGuard Statement = std::move(StatementResult).Value();
        const auto BindResult = BindText(Statement.Get(), 1, RecipeId);
        if (!BindResult.HasValue()) {
            return BindResult;
        }
        const auto Result = StepDone(Database, Statement.Get());
        if (!Result.HasValue()) {
            return Result;
        }
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> SaveRecipeBase(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue,
    bool Update) {
    const char* Sql = Update
                          ? "UPDATE Recipes SET Slug = ?, Name = ?, Cuisine = ?, Description = ?, "
                            "PrepMinutes = ?, CookMinutes = ?, Servings = ?, Difficulty = ?, "
                            "ImagePath = ?, Status = ?, AllergensJson = ?, CookwareJson = ?, "
                            "UpdatedAt = CURRENT_TIMESTAMP WHERE Id = ?;"
                          : "INSERT INTO Recipes "
                            "(Id, Slug, Name, Cuisine, Description, PrepMinutes, CookMinutes, "
                            "Servings, Difficulty, ImagePath, Status, AllergensJson, CookwareJson) "
                            "VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);";
    auto StatementResult = Prepare(Database, Sql);
    if (!StatementResult.HasValue()) {
        return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    int Index = 1;
    if (!Update) {
        const auto IdBind = BindText(Statement.Get(), Index++, RecipeValue.Id);
        if (!IdBind.HasValue()) {
            return IdBind;
        }
    }
    const std::array<std::string_view, 5> TextValues = {
        RecipeValue.Slug, RecipeValue.Name, RecipeValue.Cuisine,
        RecipeValue.Description, RecipeValue.ImagePath};
    const auto SlugBind = BindText(Statement.Get(), Index++, TextValues[0]);
    const auto NameBind = BindText(Statement.Get(), Index++, TextValues[1]);
    const auto CuisineBind = BindText(Statement.Get(), Index++, TextValues[2]);
    const auto DescriptionBind = BindText(Statement.Get(), Index++, TextValues[3]);
    if (!SlugBind.HasValue() || !NameBind.HasValue() || !CuisineBind.HasValue() ||
        !DescriptionBind.HasValue()) {
        const Foundation::Error ErrorValue = !SlugBind.HasValue()
                                                 ? SlugBind.ErrorValue()
                                                 : !NameBind.HasValue()
                                                       ? NameBind.ErrorValue()
                                                       : !CuisineBind.HasValue()
                                                             ? CuisineBind.ErrorValue()
                                                             : DescriptionBind.ErrorValue();
        return Foundation::Result<void>::FromError(ErrorValue);
    }
    const auto PrepBind = BindInteger(Statement.Get(), Index++, RecipeValue.PrepMinutes);
    const auto CookBind = BindInteger(Statement.Get(), Index++, RecipeValue.CookMinutes);
    const auto ServingBind = BindInteger(Statement.Get(), Index++, RecipeValue.Servings);
    const auto DifficultyBind = BindInteger(Statement.Get(), Index++, RecipeValue.Difficulty);
    if (!PrepBind.HasValue() || !CookBind.HasValue() || !ServingBind.HasValue() ||
        !DifficultyBind.HasValue()) {
        const Foundation::Error ErrorValue = !PrepBind.HasValue()
                                                 ? PrepBind.ErrorValue()
                                                 : !CookBind.HasValue()
                                                       ? CookBind.ErrorValue()
                                                       : !ServingBind.HasValue()
                                                             ? ServingBind.ErrorValue()
                                                             : DifficultyBind.ErrorValue();
        return Foundation::Result<void>::FromError(ErrorValue);
    }
    const auto ImageBind = BindText(Statement.Get(), Index++, RecipeValue.ImagePath);
    const auto StatusBind = BindText(Statement.Get(), Index++, RecipeValue.Status);
    const std::string AllergensJson = SerializeStrings(RecipeValue.Allergens);
    const std::string CookwareJson = SerializeStrings(RecipeValue.Cookware);
    const auto AllergensBind = BindText(Statement.Get(), Index++, AllergensJson);
    const auto CookwareBind = BindText(Statement.Get(), Index++, CookwareJson);
    if (!ImageBind.HasValue() || !StatusBind.HasValue() || !AllergensBind.HasValue() ||
        !CookwareBind.HasValue()) {
        const Foundation::Error ErrorValue = !ImageBind.HasValue()
                                                 ? ImageBind.ErrorValue()
                                                 : !StatusBind.HasValue()
                                                       ? StatusBind.ErrorValue()
                                                       : !AllergensBind.HasValue()
                                                             ? AllergensBind.ErrorValue()
                                                             : CookwareBind.ErrorValue();
        return Foundation::Result<void>::FromError(ErrorValue);
    }
    if (Update) {
        const auto IdBind = BindText(Statement.Get(), Index, RecipeValue.Id);
        if (!IdBind.HasValue()) {
            return IdBind;
        }
    }
    const auto Result = StepDone(Database, Statement.Get());
    if (!Result.HasValue()) {
        return Result;
    }
    if (Update && sqlite3_changes(Database.NativeHandle()) != 1) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "菜谱不存在"));
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
    if (StepResult == SQLITE_DONE) {
        return Foundation::Result<Domain::Recipe>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "菜谱不存在"));
    }
    if (StepResult != SQLITE_ROW) {
        return Foundation::Result<Domain::Recipe>::FromError(
            StorageError(Database.NativeHandle()));
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

Foundation::Result<std::vector<Domain::Recipe>> SqliteRecipeRepository::ListAll() {
    auto StatementResult = Prepare(
        Database,
        "SELECT Id, Slug, Name, Cuisine, Description, PrepMinutes, CookMinutes, Servings, "
        "Difficulty, ImagePath, Status, AllergensJson, CookwareJson FROM Recipes "
        "ORDER BY UpdatedAt DESC, Id;");
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

Foundation::Result<Domain::Recipe> SqliteRecipeRepository::Create(
    const Domain::Recipe& RecipeValue) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(BeginResult.ErrorValue());
    }
    const auto BaseResult = SaveRecipeBase(Database, RecipeValue, false);
    if (!BaseResult.HasValue()) {
        Database.Rollback();
        if (sqlite3_errcode(Database.NativeHandle()) == SQLITE_CONSTRAINT) {
            return Foundation::Result<Domain::Recipe>::FromError(
                Foundation::Error(Foundation::ErrorCode::Conflict, "菜谱 ID 或 slug 已存在"));
        }
        return Foundation::Result<Domain::Recipe>::FromError(BaseResult.ErrorValue());
    }
    const auto IngredientResult = InsertIngredientRows(Database, RecipeValue);
    if (!IngredientResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(IngredientResult.ErrorValue());
    }
    const auto StepResult = InsertStepRows(Database, RecipeValue);
    if (!StepResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(StepResult.ErrorValue());
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(CommitResult.ErrorValue());
    }
    return RecipeValue;
}

Foundation::Result<Domain::Recipe> SqliteRecipeRepository::Update(
    std::string_view Id,
    const Domain::Recipe& RecipeValue) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(BeginResult.ErrorValue());
    }
    const auto BaseResult = SaveRecipeBase(Database, RecipeValue, true);
    if (!BaseResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(BaseResult.ErrorValue());
    }
    if (RecipeValue.Id != Id) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(
            Foundation::Error(Foundation::ErrorCode::InvalidArgument, "菜谱 ID 不一致"));
    }
    const auto DeleteResult = DeleteRecipeChildren(Database, Id);
    if (!DeleteResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(DeleteResult.ErrorValue());
    }
    const auto IngredientResult = InsertIngredientRows(Database, RecipeValue);
    if (!IngredientResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(IngredientResult.ErrorValue());
    }
    const auto StepResult = InsertStepRows(Database, RecipeValue);
    if (!StepResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(StepResult.ErrorValue());
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(CommitResult.ErrorValue());
    }
    return RecipeValue;
}

Foundation::Result<void> SqliteRecipeRepository::Delete(std::string_view Id) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return BeginResult;
    }
    auto StatementResult = Prepare(Database, "DELETE FROM Recipes WHERE Id = ?;");
    if (!StatementResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const auto BindResult = BindText(Statement.Get(), 1, Id);
    if (!BindResult.HasValue()) {
        Database.Rollback();
        return BindResult;
    }
    const auto DeleteResult = StepDone(Database, Statement.Get());
    if (!DeleteResult.HasValue()) {
        Database.Rollback();
        return DeleteResult;
    }
    if (sqlite3_changes(Database.NativeHandle()) == 0) {
        Database.Rollback();
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "菜谱不存在"));
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return CommitResult;
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure
