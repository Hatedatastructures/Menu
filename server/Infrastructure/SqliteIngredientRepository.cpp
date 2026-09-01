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

Foundation::Result<void> BindInteger(sqlite3_stmt* Statement, int Index, int Value) {
    if (sqlite3_bind_int(Statement, Index, Value) != SQLITE_OK) {
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

Foundation::Result<Domain::Ingredient> ReadIngredient(
    SqliteDatabase& Database,
    std::string_view IngredientId) {
    auto StatementResult = Prepare(
        Database,
        "SELECT Id, Name, Category, DefaultUnit, IsPantryStaple, SubstituteGroup, "
        "StoreSkuMapping FROM Ingredients WHERE Id = ? LIMIT 1;");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const auto BindResult = BindText(Statement.Get(), 1, IngredientId);
    if (!BindResult.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(BindResult.ErrorValue());
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult == SQLITE_DONE) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "食材不存在"));
    }
    if (StepResult != SQLITE_ROW) {
        return Foundation::Result<Domain::Ingredient>::FromError(
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
        return Foundation::Result<Domain::Ingredient>::FromError(Aliases.ErrorValue());
    }
    IngredientValue.Aliases = Aliases.Value();
    return IngredientValue;
}

Foundation::Result<void> InsertAliases(
    SqliteDatabase& Database,
    const Domain::Ingredient& IngredientValue) {
    for (const std::string& Alias : IngredientValue.Aliases) {
        auto StatementResult = Prepare(
            Database,
            "INSERT INTO IngredientAliases (IngredientId, Alias) VALUES (?, ?);");
        if (!StatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
        }
        StatementGuard Statement = std::move(StatementResult).Value();
        const auto IngredientBind = BindText(Statement.Get(), 1, IngredientValue.Id);
        const auto AliasBind = BindText(Statement.Get(), 2, Alias);
        if (!IngredientBind.HasValue() || !AliasBind.HasValue()) {
            return Foundation::Result<void>::FromError(
                !IngredientBind.HasValue() ? IngredientBind.ErrorValue()
                                           : AliasBind.ErrorValue());
        }
        if (sqlite3_step(Statement.Get()) != SQLITE_DONE) {
            return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
        }
    }
    return Foundation::Result<void>();
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

Foundation::Result<Domain::Ingredient> SqliteIngredientRepository::Create(
    const Domain::Ingredient& IngredientValue) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(BeginResult.ErrorValue());
    }
    auto StatementResult = Prepare(
        Database,
        "INSERT INTO Ingredients "
        "(Id, Name, Category, DefaultUnit, IsPantryStaple, SubstituteGroup, StoreSkuMapping) "
        "VALUES (?, ?, ?, ?, ?, ?, ?);");
    if (!StatementResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const auto IdBind = BindText(Statement.Get(), 1, IngredientValue.Id);
    const auto NameBind = BindText(Statement.Get(), 2, IngredientValue.Name);
    const auto CategoryBind = BindText(Statement.Get(), 3, IngredientValue.Category);
    const auto UnitBind = BindText(Statement.Get(), 4, IngredientValue.DefaultUnit);
    const auto PantryBind = BindInteger(
        Statement.Get(), 5, IngredientValue.IsPantryStaple ? 1 : 0);
    const auto GroupBind = BindText(Statement.Get(), 6, IngredientValue.SubstituteGroup);
    const auto MappingBind = BindText(Statement.Get(), 7, IngredientValue.StoreSkuMapping);
    if (!IdBind.HasValue() || !NameBind.HasValue() || !CategoryBind.HasValue() ||
        !UnitBind.HasValue() || !PantryBind.HasValue() || !GroupBind.HasValue() ||
        !MappingBind.HasValue()) {
        Database.Rollback();
        const Foundation::Error ErrorValue = !IdBind.HasValue()
            ? IdBind.ErrorValue()
            : !NameBind.HasValue() ? NameBind.ErrorValue()
            : !CategoryBind.HasValue() ? CategoryBind.ErrorValue()
            : !UnitBind.HasValue() ? UnitBind.ErrorValue()
            : !PantryBind.HasValue() ? PantryBind.ErrorValue()
            : !GroupBind.HasValue() ? GroupBind.ErrorValue() : MappingBind.ErrorValue();
        return Foundation::Result<Domain::Ingredient>::FromError(ErrorValue);
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult != SQLITE_DONE) {
        Database.Rollback();
        if (StepResult == SQLITE_CONSTRAINT) {
            return Foundation::Result<Domain::Ingredient>::FromError(
                Foundation::Error(Foundation::ErrorCode::Conflict, "食材已存在"));
        }
        return Foundation::Result<Domain::Ingredient>::FromError(
            StorageError(Database.NativeHandle()));
    }
    const auto AliasResult = InsertAliases(Database, IngredientValue);
    if (!AliasResult.HasValue()) {
        Database.Rollback();
        if (sqlite3_errcode(Database.NativeHandle()) == SQLITE_CONSTRAINT) {
            return Foundation::Result<Domain::Ingredient>::FromError(
                Foundation::Error(Foundation::ErrorCode::Conflict, "食材别名已存在"));
        }
        return Foundation::Result<Domain::Ingredient>::FromError(AliasResult.ErrorValue());
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(CommitResult.ErrorValue());
    }
    return ReadIngredient(Database, IngredientValue.Id);
}

Foundation::Result<Domain::Ingredient> SqliteIngredientRepository::Update(
    std::string_view Id,
    const Domain::Ingredient& IngredientValue) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(BeginResult.ErrorValue());
    }
    auto StatementResult = Prepare(
        Database,
        "UPDATE Ingredients SET Name = ?, Category = ?, DefaultUnit = ?, "
        "IsPantryStaple = ?, SubstituteGroup = ?, StoreSkuMapping = ? WHERE Id = ?;");
    if (!StatementResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const auto NameBind = BindText(Statement.Get(), 1, IngredientValue.Name);
    const auto CategoryBind = BindText(Statement.Get(), 2, IngredientValue.Category);
    const auto UnitBind = BindText(Statement.Get(), 3, IngredientValue.DefaultUnit);
    const auto PantryBind = BindInteger(
        Statement.Get(), 4, IngredientValue.IsPantryStaple ? 1 : 0);
    const auto GroupBind = BindText(Statement.Get(), 5, IngredientValue.SubstituteGroup);
    const auto MappingBind = BindText(Statement.Get(), 6, IngredientValue.StoreSkuMapping);
    const auto IdBind = BindText(Statement.Get(), 7, Id);
    if (!NameBind.HasValue() || !CategoryBind.HasValue() || !UnitBind.HasValue() ||
        !PantryBind.HasValue() || !GroupBind.HasValue() || !MappingBind.HasValue() ||
        !IdBind.HasValue()) {
        Database.Rollback();
        const Foundation::Error ErrorValue = !NameBind.HasValue()
            ? NameBind.ErrorValue()
            : !CategoryBind.HasValue() ? CategoryBind.ErrorValue()
            : !UnitBind.HasValue() ? UnitBind.ErrorValue()
            : !PantryBind.HasValue() ? PantryBind.ErrorValue()
            : !GroupBind.HasValue() ? GroupBind.ErrorValue()
            : !MappingBind.HasValue() ? MappingBind.ErrorValue() : IdBind.ErrorValue();
        return Foundation::Result<Domain::Ingredient>::FromError(ErrorValue);
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult != SQLITE_DONE) {
        Database.Rollback();
        if (StepResult == SQLITE_CONSTRAINT) {
            return Foundation::Result<Domain::Ingredient>::FromError(
                Foundation::Error(Foundation::ErrorCode::Conflict, "食材名称已存在"));
        }
        return Foundation::Result<Domain::Ingredient>::FromError(
            StorageError(Database.NativeHandle()));
    }
    if (sqlite3_changes(Database.NativeHandle()) == 0) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "食材不存在"));
    }
    auto DeleteResult = Prepare(
        Database, "DELETE FROM IngredientAliases WHERE IngredientId = ?;");
    if (!DeleteResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(DeleteResult.ErrorValue());
    }
    StatementGuard DeleteStatement = std::move(DeleteResult).Value();
    const auto DeleteBind = BindText(DeleteStatement.Get(), 1, Id);
    if (!DeleteBind.HasValue() || sqlite3_step(DeleteStatement.Get()) != SQLITE_DONE) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(
            !DeleteBind.HasValue() ? DeleteBind.ErrorValue()
                                   : StorageError(Database.NativeHandle()));
    }
    const auto AliasResult = InsertAliases(Database, IngredientValue);
    if (!AliasResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(AliasResult.ErrorValue());
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(CommitResult.ErrorValue());
    }
    return ReadIngredient(Database, Id);
}

Foundation::Result<void> SqliteIngredientRepository::Delete(std::string_view Id) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return BeginResult;
    }
    auto StatementResult = Prepare(
        Database, "DELETE FROM Ingredients WHERE Id = ?;");
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
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult != SQLITE_DONE) {
        Database.Rollback();
        if (StepResult == SQLITE_CONSTRAINT) {
            return Foundation::Result<void>::FromError(
                Foundation::Error(Foundation::ErrorCode::Conflict, "食材已被菜谱使用"));
        }
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    if (sqlite3_changes(Database.NativeHandle()) == 0) {
        Database.Rollback();
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "食材不存在"));
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return CommitResult;
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure
