#include "Infrastructure/SqliteIngredientRepository.hpp"

#include "Infrastructure/SqliteIngredientRepositorySupport.hpp"
#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteIngredientRepositorySupport::ReadAliases;
using SqliteWorkflowSupport::ColumnText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

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
