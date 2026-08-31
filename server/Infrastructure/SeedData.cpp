#include "Infrastructure/SeedData.hpp"

#include <Domain/Ingredient.hpp>
#include <Domain/Recipe.hpp>

#include <boost/json/array.hpp>
#include <boost/json/serialize.hpp>

#include <array>
#include <string>
#include <string_view>
#include <utility>
#include <vector>

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
            if (Statement != nullptr) {
                sqlite3_finalize(Statement);
            }
            Statement = std::exchange(Other.Statement, nullptr);
        }
        return *this;
    }

    ~StatementGuard() {
        if (Statement != nullptr) {
            sqlite3_finalize(Statement);
        }
    }

    [[nodiscard]] sqlite3_stmt* Get() const noexcept {
        return Statement;
    }

private:
    sqlite3_stmt* Statement = nullptr;
};

Foundation::Error StorageError(sqlite3* Handle) {
    const char* Message = Handle == nullptr ? nullptr : sqlite3_errmsg(Handle);
    return Foundation::Error(
        Foundation::ErrorCode::StorageUnavailable,
        Message == nullptr ? "SQLite 写入失败" : std::string(Message));
}

Foundation::Result<StatementGuard> Prepare(SqliteDatabase& Database, const char* Sql) {
    sqlite3_stmt* Statement = nullptr;
    if (sqlite3_prepare_v2(Database.NativeHandle(), Sql, -1, &Statement, nullptr) != SQLITE_OK) {
        return Foundation::Result<StatementGuard>::FromError(
            StorageError(Database.NativeHandle()));
    }
    return StatementGuard(Statement);
}

Foundation::Result<void> BindText(sqlite3_stmt* Statement, int Index, std::string_view Value) {
    const std::string Text(Value);
    if (sqlite3_bind_text(Statement, Index, Text.c_str(), -1, SQLITE_TRANSIENT) != SQLITE_OK) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 绑定失败"));
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> BindInt(sqlite3_stmt* Statement, int Index, int Value) {
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

Foundation::Result<void> StepDone(
    SqliteDatabase& Database,
    sqlite3_stmt* Statement) {
    if (sqlite3_step(Statement) != SQLITE_DONE) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
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

std::vector<Domain::Ingredient> BuildIngredients() {
    return {
        {"ingredient.rice", "大米", {"米饭"}, "主食", "g", true, "", ""},
        {"ingredient.chicken", "鸡腿肉", {"鸡肉"}, "肉禽", "g", false, "", ""},
        {"ingredient.beef", "牛肉", {"牛里脊"}, "肉禽", "g", false, "", ""},
        {"ingredient.egg", "鸡蛋", {"蛋"}, "蛋奶", "个", true, "", ""},
        {"ingredient.tomato", "番茄", {"西红柿", "蕃茄"}, "蔬菜", "g", false, "", ""},
        {"ingredient.onion", "洋葱", {"葱头"}, "蔬菜", "g", true, "", ""},
        {"ingredient.potato", "土豆", {"马铃薯"}, "蔬菜", "g", false, "", ""},
        {"ingredient.carrot", "胡萝卜", {"红萝卜"}, "蔬菜", "g", false, "", ""},
        {"ingredient.spinach", "菠菜", {"菠薐菜"}, "蔬菜", "g", false, "", ""},
        {"ingredient.tofu", "北豆腐", {"豆腐"}, "豆制品", "g", false, "", ""},
        {"ingredient.salmon", "三文鱼", {"鲑鱼"}, "水产", "g", false, "", ""},
        {"ingredient.pasta", "意大利面", {"通心粉"}, "主食", "g", true, "", ""},
        {"ingredient.garlic", "大蒜", {"蒜"}, "调味", "g", true, "", ""},
        {"ingredient.ginger", "生姜", {"姜"}, "调味", "g", true, "", ""},
        {"ingredient.soy_sauce", "生抽", {"酱油"}, "调味", "ml", true, "", ""},
        {"ingredient.salt", "食盐", {"盐"}, "调味", "g", true, "", ""},
        {"ingredient.oil", "食用油", {"油"}, "调味", "ml", true, "", ""},
        {"ingredient.sesame", "白芝麻", {"芝麻"}, "调味", "g", true, "", ""},
        {"ingredient.lemon", "柠檬", {}, "水果", "个", false, "", ""},
        {"ingredient.cream", "淡奶油", {"奶油"}, "蛋奶", "ml", false, "", ""},
    };
}

Domain::Recipe BuildRecipe(
    std::string Id,
    std::string Slug,
    std::string Name,
    std::string Cuisine,
    std::string Description,
    int PrepMinutes,
    int CookMinutes,
    int Servings,
    int Difficulty,
    std::string ImagePath,
    std::vector<std::string> Allergens,
    std::vector<std::string> Cookware,
    std::vector<Domain::RecipeIngredient> Ingredients,
    std::string FirstInstruction,
    std::string SecondInstruction) {
    Domain::Recipe RecipeValue;
    RecipeValue.Id = std::move(Id);
    RecipeValue.Slug = std::move(Slug);
    RecipeValue.Name = std::move(Name);
    RecipeValue.Cuisine = std::move(Cuisine);
    RecipeValue.Description = std::move(Description);
    RecipeValue.PrepMinutes = PrepMinutes;
    RecipeValue.CookMinutes = CookMinutes;
    RecipeValue.Servings = Servings;
    RecipeValue.Difficulty = Difficulty;
    RecipeValue.ImagePath = std::move(ImagePath);
    RecipeValue.Status = "published";
    RecipeValue.Allergens = std::move(Allergens);
    RecipeValue.Cookware = std::move(Cookware);
    RecipeValue.Ingredients = std::move(Ingredients);
    RecipeValue.Steps = {
        {1, "准备", "清洗、切配并称量食材。", PrepMinutes * 60, false},
        {2, "烹饪", std::move(FirstInstruction), CookMinutes * 60, true},
        {3, "收尾", std::move(SecondInstruction), 60, false},
    };
    return RecipeValue;
}

std::vector<Domain::Recipe> BuildRecipes() {
    return {
        BuildRecipe("recipe.tomato_egg", "tomato-egg", "番茄炒蛋", "中餐",
                     "酸甜下饭，适合工作日晚餐。", 10, 12, 2, 1,
                     "assets/media/tomato-egg.png", {}, {"炒锅"},
                     {{"ingredient.tomato", 300.0, "g", true},
                      {"ingredient.egg", 3.0, "个", true},
                      {"ingredient.oil", 20.0, "ml", true},
                      {"ingredient.salt", 2.0, "g", true}},
                     "先炒鸡蛋至凝固后盛出，再把番茄炒出汁。",
                     "倒回鸡蛋快速翻匀，调味后关火。"),
        BuildRecipe("recipe.soy_chicken", "soy-chicken", "葱油鸡腿", "中餐",
                     "一锅完成的嫩鸡腿，配米饭即可。", 12, 25, 2, 2,
                     "assets/media/soy-chicken.png", {"大豆"}, {"炒锅"},
                     {{"ingredient.chicken", 400.0, "g", true},
                      {"ingredient.onion", 100.0, "g", true},
                      {"ingredient.ginger", 10.0, "g", true},
                      {"ingredient.soy_sauce", 25.0, "ml", true}},
                     "鸡腿肉擦干，皮面煎至金黄。",
                     "加入葱姜和调味汁，加盖焖熟后切片。"),
        BuildRecipe("recipe.beef_onion", "beef-onion", "洋葱黑椒牛肉", "中餐",
                     "快手大火炒制，牛肉软嫩不拖沓。", 15, 10, 2, 2,
                     "assets/media/beef-onion.png", {"大豆"}, {"炒锅"},
                     {{"ingredient.beef", 300.0, "g", true},
                      {"ingredient.onion", 180.0, "g", true},
                      {"ingredient.garlic", 10.0, "g", true},
                      {"ingredient.soy_sauce", 20.0, "ml", true}},
                     "牛肉切薄片，洋葱切块并分开放置。",
                     "炒香牛肉后加入洋葱和酱汁，大火翻匀。"),
        BuildRecipe("recipe.tofu_spinach", "tofu-spinach", "菠菜豆腐汤", "中餐",
                     "清淡暖胃，十五分钟内完成。", 8, 12, 2, 1,
                     "assets/media/tofu-spinach.png", {}, {"汤锅"},
                     {{"ingredient.tofu", 250.0, "g", true},
                      {"ingredient.spinach", 200.0, "g", true},
                      {"ingredient.egg", 1.0, "个", false},
                      {"ingredient.salt", 2.0, "g", true}},
                     "汤锅加水煮开，放入豆腐煮透。",
                     "加入菠菜和蛋液，调味后立即关火。"),
        BuildRecipe("recipe.pasta_tomato", "pasta-tomato", "番茄奶油意面", "西餐",
                     "酸香顺滑的一锅意面，适合快速晚餐。", 10, 20, 2, 2,
                     "assets/media/pasta-tomato.png", {"乳制品", "小麦"}, {"汤锅", "炒锅"},
                     {{"ingredient.pasta", 220.0, "g", true},
                      {"ingredient.tomato", 250.0, "g", true},
                      {"ingredient.cream", 100.0, "ml", true},
                      {"ingredient.garlic", 10.0, "g", true}},
                     "煮意面至略有硬芯，同时炒香蒜末和番茄。",
                     "加入淡奶油和意面水，拌入意面收汁。"),
        BuildRecipe("recipe.lemon_salmon", "lemon-salmon", "柠檬烤三文鱼", "西餐",
                     "烤箱少油料理，准备好即可等待出炉。", 8, 22, 2, 1,
                     "assets/media/lemon-salmon.png", {"鱼类"}, {"烤箱"},
                     {{"ingredient.salmon", 360.0, "g", true},
                      {"ingredient.lemon", 1.0, "个", true},
                      {"ingredient.salt", 2.0, "g", true},
                      {"ingredient.oil", 10.0, "ml", true}},
                     "三文鱼擦干，放在烤盘上并挤入柠檬汁。",
                     "烤箱烤至中心刚熟，出炉静置后分切。"),
        BuildRecipe("recipe.cream_spinach", "cream-spinach", "奶油菠菜", "西餐",
                     "适合作为肉类主菜的配菜。", 5, 10, 2, 1,
                     "assets/media/cream-spinach.png", {"乳制品"}, {"炒锅"},
                     {{"ingredient.spinach", 300.0, "g", true},
                      {"ingredient.cream", 120.0, "ml", true},
                      {"ingredient.garlic", 8.0, "g", true},
                      {"ingredient.salt", 2.0, "g", true}},
                     "菠菜焯水挤干，炒香蒜末。",
                     "加入菠菜和淡奶油煮至浓稠，调味即可。"),
        BuildRecipe("recipe.oyakodon", "oyakodon", "亲子丼", "日系",
                     "鸡肉与鸡蛋盖饭，温和饱腹。", 12, 18, 2, 2,
                     "assets/media/oyakodon.png", {"大豆", "蛋类"}, {"炒锅"},
                     {{"ingredient.chicken", 280.0, "g", true},
                      {"ingredient.egg", 3.0, "个", true},
                      {"ingredient.onion", 150.0, "g", true},
                      {"ingredient.rice", 360.0, "g", true},
                      {"ingredient.soy_sauce", 25.0, "ml", true}},
                     "洋葱和鸡肉放入调味汁中煮至熟透。",
                     "淋入蛋液，半凝固时盖在热米饭上。"),
        BuildRecipe("recipe.tofu_teriyaki", "tofu-teriyaki", "照烧豆腐", "日系",
                     "外脆里嫩，植物性蛋白的简单做法。", 10, 15, 2, 1,
                     "assets/media/tofu-teriyaki.png", {"大豆"}, {"炒锅"},
                     {{"ingredient.tofu", 350.0, "g", true},
                      {"ingredient.soy_sauce", 30.0, "ml", true},
                      {"ingredient.garlic", 8.0, "g", false},
                      {"ingredient.sesame", 5.0, "g", false}},
                     "豆腐切块并擦干，平底锅煎至四面金黄。",
                     "倒入照烧汁收浓，撒芝麻后装盘。"),
        BuildRecipe("recipe.ginger_beef", "ginger-beef", "姜汁牛肉盖饭", "日系",
                     "姜香明显的快手牛肉盖饭。", 10, 16, 2, 2,
                     "assets/media/ginger-beef.png", {"大豆"}, {"炒锅"},
                     {{"ingredient.beef", 300.0, "g", true},
                      {"ingredient.ginger", 20.0, "g", true},
                      {"ingredient.onion", 100.0, "g", true},
                      {"ingredient.rice", 360.0, "g", true},
                      {"ingredient.soy_sauce", 25.0, "ml", true}},
                     "牛肉切片，洋葱切丝并调好姜汁酱汁。",
                     "洋葱炒软后加入牛肉和酱汁，煮至刚熟盖饭。"),
        BuildRecipe("recipe.tomato_rice", "tomato-rice", "番茄焖饭", "中餐",
                     "电饭锅友好的一锅主食，准备后即可等待。", 10, 35, 2, 1,
                     "assets/media/tomato-rice.png", {}, {"电饭锅"},
                     {{"ingredient.rice", 300.0, "g", true},
                      {"ingredient.tomato", 250.0, "g", true},
                      {"ingredient.carrot", 80.0, "g", true},
                      {"ingredient.egg", 2.0, "个", false}},
                     "大米淘洗后放入电饭锅，铺上番茄和切丁蔬菜。",
                     "按米饭程序煮熟，拌匀后按需加入溏心蛋。"),
        BuildRecipe("recipe.tofu_salad", "tofu-salad", "芝麻豆腐沙拉", "西餐",
                     "无需开火的清爽晚餐配菜。", 12, 0, 2, 1,
                     "assets/media/tofu-salad.png", {"大豆"}, {},
                     {{"ingredient.tofu", 250.0, "g", true},
                      {"ingredient.spinach", 120.0, "g", true},
                      {"ingredient.sesame", 10.0, "g", false},
                      {"ingredient.lemon", 0.5, "个", false}},
                     "豆腐沥水切块，菠菜洗净并充分沥干。",
                     "混合柠檬汁和芝麻，淋在食材上即可。"),
    };
}

Foundation::Result<void> InsertIngredient(
    SqliteDatabase& Database,
    const Domain::Ingredient& IngredientValue) {
    auto StatementResult = Prepare(
        Database,
        "INSERT OR IGNORE INTO Ingredients "
        "(Id, Name, Category, DefaultUnit, IsPantryStaple, SubstituteGroup, StoreSkuMapping) "
        "VALUES (?, ?, ?, ?, ?, ?, ?);");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const std::array<std::string_view, 4> TextValues = {
        IngredientValue.Id, IngredientValue.Name, IngredientValue.Category,
        IngredientValue.DefaultUnit};
    for (std::size_t Index = 0; Index < TextValues.size(); ++Index) {
        const auto BindResult = BindText(Statement.Get(), static_cast<int>(Index + 1), TextValues[Index]);
        if (!BindResult.HasValue()) {
            return BindResult;
        }
    }
    const auto PantryResult = BindInt(Statement.Get(), 5, IngredientValue.IsPantryStaple ? 1 : 0);
    const auto GroupResult = BindText(Statement.Get(), 6, IngredientValue.SubstituteGroup);
    const auto SkuResult = BindText(Statement.Get(), 7, IngredientValue.StoreSkuMapping);
    if (!PantryResult.HasValue()) {
        return PantryResult;
    }
    if (!GroupResult.HasValue()) {
        return GroupResult;
    }
    if (!SkuResult.HasValue()) {
        return SkuResult;
    }
    const auto InsertResult = StepDone(Database, Statement.Get());
    if (!InsertResult.HasValue()) {
        return InsertResult;
    }

    for (const std::string& Alias : IngredientValue.Aliases) {
        auto AliasStatementResult = Prepare(
            Database,
            "INSERT OR IGNORE INTO IngredientAliases (IngredientId, Alias) VALUES (?, ?);");
        if (!AliasStatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(AliasStatementResult.ErrorValue());
        }
        StatementGuard AliasStatement = std::move(AliasStatementResult).Value();
        const auto IngredientBind = BindText(AliasStatement.Get(), 1, IngredientValue.Id);
        const auto AliasBind = BindText(AliasStatement.Get(), 2, Alias);
        if (!IngredientBind.HasValue()) {
            return IngredientBind;
        }
        if (!AliasBind.HasValue()) {
            return AliasBind;
        }
        const auto AliasResult = StepDone(Database, AliasStatement.Get());
        if (!AliasResult.HasValue()) {
            return AliasResult;
        }
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> InsertRecipe(SqliteDatabase& Database, const Domain::Recipe& RecipeValue) {
    auto StatementResult = Prepare(
        Database,
        "INSERT OR IGNORE INTO Recipes "
        "(Id, Slug, Name, Cuisine, Description, PrepMinutes, CookMinutes, Servings, Difficulty, "
        "ImagePath, Status, AllergensJson, CookwareJson) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const std::array<std::string_view, 5> TextValues = {
        RecipeValue.Id, RecipeValue.Slug, RecipeValue.Name, RecipeValue.Cuisine,
        RecipeValue.Description};
    for (std::size_t Index = 0; Index < TextValues.size(); ++Index) {
        const auto BindResult = BindText(Statement.Get(), static_cast<int>(Index + 1), TextValues[Index]);
        if (!BindResult.HasValue()) {
            return BindResult;
        }
    }
    const std::array<int, 4> IntValues = {
        RecipeValue.PrepMinutes, RecipeValue.CookMinutes, RecipeValue.Servings,
        RecipeValue.Difficulty};
    for (std::size_t Index = 0; Index < IntValues.size(); ++Index) {
        const auto BindResult = BindInt(Statement.Get(), static_cast<int>(Index + 6), IntValues[Index]);
        if (!BindResult.HasValue()) {
            return BindResult;
        }
    }
    const auto ImageResult = BindText(Statement.Get(), 10, RecipeValue.ImagePath);
    const auto StatusResult = BindText(Statement.Get(), 11, RecipeValue.Status);
    const std::string AllergensJson = SerializeStrings(RecipeValue.Allergens);
    const std::string CookwareJson = SerializeStrings(RecipeValue.Cookware);
    const auto AllergensResult = BindText(Statement.Get(), 12, AllergensJson);
    const auto CookwareResult = BindText(Statement.Get(), 13, CookwareJson);
    if (!ImageResult.HasValue()) {
        return ImageResult;
    }
    if (!StatusResult.HasValue()) {
        return StatusResult;
    }
    if (!AllergensResult.HasValue()) {
        return AllergensResult;
    }
    if (!CookwareResult.HasValue()) {
        return CookwareResult;
    }
    const auto RecipeResult = StepDone(Database, Statement.Get());
    if (!RecipeResult.HasValue()) {
        return RecipeResult;
    }

    for (const Domain::RecipeIngredient& IngredientValue : RecipeValue.Ingredients) {
        auto IngredientStatementResult = Prepare(
            Database,
            "INSERT OR IGNORE INTO RecipeIngredients "
            "(RecipeId, IngredientId, Quantity, Unit, ServingFactor, Preparation, Required) "
            "VALUES (?, ?, ?, ?, ?, ?, ?);");
        if (!IngredientStatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(IngredientStatementResult.ErrorValue());
        }
        StatementGuard IngredientStatement = std::move(IngredientStatementResult).Value();
        const auto RecipeBind = BindText(IngredientStatement.Get(), 1, RecipeValue.Id);
        const auto IngredientBind = BindText(IngredientStatement.Get(), 2, IngredientValue.IngredientId);
        const auto QuantityBind = BindDouble(IngredientStatement.Get(), 3, IngredientValue.Quantity);
        const auto UnitBind = BindText(IngredientStatement.Get(), 4, IngredientValue.Unit);
        const auto FactorBind = BindDouble(IngredientStatement.Get(), 5, IngredientValue.ServingFactor);
        const auto PreparationBind = BindText(IngredientStatement.Get(), 6, IngredientValue.Preparation);
        const auto RequiredBind = BindInt(IngredientStatement.Get(), 7, IngredientValue.Required ? 1 : 0);
        if (!RecipeBind.HasValue()) {
            return RecipeBind;
        }
        if (!IngredientBind.HasValue()) {
            return IngredientBind;
        }
        if (!QuantityBind.HasValue()) {
            return QuantityBind;
        }
        if (!UnitBind.HasValue()) {
            return UnitBind;
        }
        if (!FactorBind.HasValue()) {
            return FactorBind;
        }
        if (!PreparationBind.HasValue()) {
            return PreparationBind;
        }
        if (!RequiredBind.HasValue()) {
            return RequiredBind;
        }
        const auto IngredientResult = StepDone(Database, IngredientStatement.Get());
        if (!IngredientResult.HasValue()) {
            return IngredientResult;
        }
    }

    for (const Domain::RecipeStep& StepValue : RecipeValue.Steps) {
        auto StepStatementResult = Prepare(
            Database,
            "INSERT OR IGNORE INTO RecipeSteps "
            "(RecipeId, StepOrder, Title, Instruction, DurationSeconds, HasTimer) "
            "VALUES (?, ?, ?, ?, ?, ?);");
        if (!StepStatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(StepStatementResult.ErrorValue());
        }
        StatementGuard StepStatement = std::move(StepStatementResult).Value();
        const auto RecipeBind = BindText(StepStatement.Get(), 1, RecipeValue.Id);
        const auto OrderBind = BindInt(StepStatement.Get(), 2, StepValue.StepOrder);
        const auto TitleBind = BindText(StepStatement.Get(), 3, StepValue.Title);
        const auto InstructionBind = BindText(StepStatement.Get(), 4, StepValue.Instruction);
        const auto DurationBind = BindInt(StepStatement.Get(), 5, StepValue.DurationSeconds);
        const auto TimerBind = BindInt(StepStatement.Get(), 6, StepValue.HasTimer ? 1 : 0);
        if (!RecipeBind.HasValue()) {
            return RecipeBind;
        }
        if (!OrderBind.HasValue()) {
            return OrderBind;
        }
        if (!TitleBind.HasValue()) {
            return TitleBind;
        }
        if (!InstructionBind.HasValue()) {
            return InstructionBind;
        }
        if (!DurationBind.HasValue()) {
            return DurationBind;
        }
        if (!TimerBind.HasValue()) {
            return TimerBind;
        }
        const auto StepResult = StepDone(Database, StepStatement.Get());
        if (!StepResult.HasValue()) {
            return StepResult;
        }
    }
    return Foundation::Result<void>();
}

}  // namespace

Foundation::Result<void> SeedData::InsertIfEmpty(SqliteDatabase& Database) {
    const auto IngredientCount = Database.ScalarInt("SELECT COUNT(*) FROM Ingredients;");
    const auto RecipeCount = Database.ScalarInt("SELECT COUNT(*) FROM Recipes;");
    if (!IngredientCount.HasValue()) {
        return Foundation::Result<void>::FromError(IngredientCount.ErrorValue());
    }
    if (!RecipeCount.HasValue()) {
        return Foundation::Result<void>::FromError(RecipeCount.ErrorValue());
    }
    if (IngredientCount.Value() > 0 && RecipeCount.Value() > 0) {
        return Foundation::Result<void>();
    }

    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return BeginResult;
    }
    for (const Domain::Ingredient& IngredientValue : BuildIngredients()) {
        const auto Result = InsertIngredient(Database, IngredientValue);
        if (!Result.HasValue()) {
            Database.Rollback();
            return Result;
        }
    }
    for (const Domain::Recipe& RecipeValue : BuildRecipes()) {
        const auto Result = InsertRecipe(Database, RecipeValue);
        if (!Result.HasValue()) {
            Database.Rollback();
            return Result;
        }
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return CommitResult;
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure
