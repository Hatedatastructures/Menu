#include "Infrastructure/SeedCatalog.hpp"

#include <utility>

namespace Menu::Infrastructure {
namespace {

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
        {2, "烹饪", std::move(FirstInstruction), CookMinutes * 60, CookMinutes > 0},
        {3, "收尾", std::move(SecondInstruction), 60, false},
    };
    return RecipeValue;
}

}  // namespace

std::vector<Domain::Ingredient> SeedCatalog::Ingredients() {
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

std::vector<Domain::Recipe> SeedCatalog::Recipes() {
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

}  // namespace Menu::Infrastructure
