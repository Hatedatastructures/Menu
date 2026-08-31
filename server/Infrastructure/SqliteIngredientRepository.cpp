#include "Infrastructure/SqliteIngredientRepository.hpp"

#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure {
namespace {

class StatementGuard final {
public:
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

private:
    void Reset() noexcept {
        if (Statement != nullptr) {
            sqlite3_finalize(Statement);
            Statement = nullptr;
        }
    }

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

Foundation::Result<void> BindText(sqlite3_stmt* Statement, int Index, std::string_view Value) {
    const std::string Text(Value);
    if (sqlite3_bind_text(Statement, Index, Text.c_str(), -1, SQLITE_TRANSIENT) != SQLITE_OK) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 绑定失败"));
    }
    return Foundation::Result<void>();
}

Foundation::Result<std::vector<std::string>> ReadAliases(
    SqliteDatabase& Database,
    std::string_view IngredientId) {
    auto StatementResult = Prepare(
        Database,
        "SELECT Alias FROM IngredientAliases WHERE IngredientId = ? ORDER BY Alias;");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<std::vector<std::string>>::FromError(
            StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const auto BindResult = BindText(Statement.Get(), 1, IngredientId);
    if (!BindResult.HasValue()) {
        return Foundation::Result<std::vector<std::string>>::FromError(BindResult.ErrorValue());
    }
    std::vector<std::string> Aliases;
    while (true) {
        const int StepResult = sqlite3_step(Statement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<std::vector<std::string>>::FromError(
                StorageError(Database.NativeHandle()));
        }
        Aliases.push_back(ColumnText(Statement.Get(), 0));
    }
    return Aliases;
}

}  // namespace

Foundation::Result<std::vector<Domain::Ingredient>>
SqliteIngredientRepository::ListAll() {
    auto StatementResult = Prepare(
        Database,
        "SELECT Id, Name, Category, DefaultUnit, IsPantryStaple, SubstituteGroup, "
        "StoreSkuMapping FROM Ingredients ORDER BY Name, Id;");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<std::vector<Domain::Ingredient>>::FromError(
            StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    std::vector<Domain::Ingredient> Ingredients;
    while (true) {
        const int StepResult = sqlite3_step(Statement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<std::vector<Domain::Ingredient>>::FromError(
                StorageError(Database.NativeHandle()));
        }
        Domain::Ingredient IngredientValue;
        IngredientValue.Id = ColumnText(Statement.Get(), 0);
        IngredientValue.Name = ColumnText(Statement.Get(), 1);
        IngredientValue.Category = ColumnText(Statement.Get(), 2);
        IngredientValue.DefaultUnit = ColumnText(Statement.Get(), 3);
        IngredientValue.IsPantryStaple = sqlite3_column_int(Statement.Get(), 4) != 0;
        IngredientValue.SubstituteGroup = ColumnText(Statement.Get(), 5);
        IngredientValue.StoreSkuMapping = ColumnText(Statement.Get(), 6);
        const auto Aliases = ReadAliases(Database, IngredientValue.Id);
        if (!Aliases.HasValue()) {
            return Foundation::Result<std::vector<Domain::Ingredient>>::FromError(
                Aliases.ErrorValue());
        }
        IngredientValue.Aliases = Aliases.Value();
        Ingredients.push_back(std::move(IngredientValue));
    }
    return Ingredients;
}

}  // namespace Menu::Infrastructure
