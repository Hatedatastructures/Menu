#include "Domain/RecipeRules.hpp"

#include <algorithm>
#include <array>
#include <cctype>
#include <ranges>
#include <string>
#include <unordered_set>

namespace Menu::Domain {
namespace {

std::string Trim(std::string_view Value) {
    std::size_t First = 0;
    while (First < Value.size() &&
           std::isspace(static_cast<unsigned char>(Value[First])) != 0) {
        ++First;
    }

    std::size_t Last = Value.size();
    while (Last > First &&
           std::isspace(static_cast<unsigned char>(Value[Last - 1])) != 0) {
        --Last;
    }

    return std::string(Value.substr(First, Last - First));
}

std::string LowerAscii(std::string Value) {
    for (char& Character : Value) {
        Character = static_cast<char>(
            std::tolower(static_cast<unsigned char>(Character)));
    }
    return Value;
}

bool MatchesNormalized(std::string_view Left, std::string_view Right) {
    return LowerAscii(Trim(Left)) == LowerAscii(Trim(Right));
}

bool IsSupportedUnit(std::string_view Unit) {
    static constexpr std::array<std::string_view, 8> Units = {
        "g", "kg", "ml", "l", "个", "份", "勺", "茶匙"};
    return std::ranges::find(Units, Unit) != Units.end();
}

Foundation::Error InvalidRecipe(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

}  // namespace

Foundation::Result<void> RecipeRules::ValidateRecipe(const Recipe& RecipeValue) {
    if (RecipeValue.Id.empty() || RecipeValue.Name.empty() || RecipeValue.Cuisine.empty()) {
        return Foundation::Result<void>::FromError(
            InvalidRecipe("菜谱必须包含 ID、名称和菜系"));
    }
    if (RecipeValue.Servings <= 0) {
        return Foundation::Result<void>::FromError(
            InvalidRecipe("菜谱基准份量必须大于零"));
    }
    if (RecipeValue.PrepMinutes < 0 || RecipeValue.CookMinutes < 0) {
        return Foundation::Result<void>::FromError(
            InvalidRecipe("菜谱时长不能为负数"));
    }
    if (RecipeValue.Difficulty < 1 || RecipeValue.Difficulty > 5) {
        return Foundation::Result<void>::FromError(
            InvalidRecipe("菜谱难度必须在一到五之间"));
    }
    if (RecipeValue.Ingredients.empty()) {
        return Foundation::Result<void>::FromError(
            InvalidRecipe("菜谱至少需要一种食材"));
    }

    for (const RecipeIngredient& IngredientValue : RecipeValue.Ingredients) {
        if (IngredientValue.IngredientId.empty() || IngredientValue.Quantity <= 0.0 ||
            IngredientValue.ServingFactor <= 0.0 ||
            !IsSupportedUnit(IngredientValue.Unit)) {
            return Foundation::Result<void>::FromError(
                InvalidRecipe("食材必须包含正数量和支持的单位"));
        }
    }

    for (std::size_t Index = 0; Index < RecipeValue.Steps.size(); ++Index) {
        const RecipeStep& Step = RecipeValue.Steps[Index];
        if (Step.StepOrder != static_cast<int>(Index + 1) || Step.Title.empty() ||
            Step.Instruction.empty() || Step.DurationSeconds < 0) {
            return Foundation::Result<void>::FromError(
                InvalidRecipe("步骤顺序、说明和耗时无效"));
        }
    }

    return Foundation::Result<void>();
}

Foundation::Result<std::vector<RecipeIngredient>> RecipeRules::ScaleRecipeIngredients(
    const Recipe& RecipeValue,
    int TargetServings) {
    if (TargetServings <= 0 || RecipeValue.Servings <= 0) {
        return Foundation::Result<std::vector<RecipeIngredient>>::FromError(
            InvalidRecipe("目标份量必须大于零"));
    }

    const double ServingRatio =
        static_cast<double>(TargetServings) / static_cast<double>(RecipeValue.Servings);
    std::vector<RecipeIngredient> ScaledIngredients = RecipeValue.Ingredients;
    for (RecipeIngredient& IngredientValue : ScaledIngredients) {
        IngredientValue.Quantity *= ServingRatio * IngredientValue.ServingFactor;
    }
    return ScaledIngredients;
}

Foundation::Result<std::string> RecipeRules::NormalizeIngredientAlias(
    std::string_view Alias,
    const std::vector<Ingredient>& Ingredients) {
    const std::string NormalizedAlias = LowerAscii(Trim(Alias));
    if (NormalizedAlias.empty()) {
        return Foundation::Result<std::string>::FromError(
            InvalidRecipe("食材别名不能为空"));
    }

    for (const Ingredient& IngredientValue : Ingredients) {
        if (MatchesNormalized(NormalizedAlias, IngredientValue.Name) ||
            std::ranges::any_of(IngredientValue.Aliases,
                                [&NormalizedAlias](const std::string& Candidate) {
                                    return MatchesNormalized(NormalizedAlias, Candidate);
                                })) {
            return IngredientValue.Id;
        }
    }

    return Foundation::Result<std::string>::FromError(
        InvalidRecipe("未找到匹配的食材别名"));
}

IngredientAvailability RecipeRules::BuildIngredientAvailability(
    const Recipe& RecipeValue,
    const std::vector<std::string>& PantryIngredientIds) {
    const std::unordered_set<std::string> PantryIds(
        PantryIngredientIds.begin(), PantryIngredientIds.end());
    IngredientAvailability Availability;

    for (const RecipeIngredient& IngredientValue : RecipeValue.Ingredients) {
        if (PantryIds.contains(IngredientValue.IngredientId)) {
            ++Availability.AvailableCount;
            Availability.AvailableIngredientIds.push_back(IngredientValue.IngredientId);
        } else if (IngredientValue.Required) {
            ++Availability.MissingCount;
            Availability.MissingIngredientIds.push_back(IngredientValue.IngredientId);
        }
    }
    return Availability;
}

}  // namespace Menu::Domain
