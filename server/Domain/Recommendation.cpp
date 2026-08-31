#include "Domain/Recommendation.hpp"

#include <algorithm>
#include <ranges>
#include <tuple>
#include <unordered_set>
#include <utility>

namespace Menu::Domain {
namespace {

bool ContainsValue(const std::vector<std::string>& Values, std::string_view Value) {
    return std::ranges::any_of(Values, [Value](const std::string& Candidate) {
        return Candidate == Value;
    });
}

bool Intersects(const std::vector<std::string>& Left, const std::vector<std::string>& Right) {
    return std::ranges::any_of(Left, [&Right](const std::string& Value) {
        return ContainsValue(Right, Value);
    });
}

bool CookwareAvailable(const Recipe& RecipeValue, const RecommendationRequest& Request) {
    return std::ranges::all_of(RecipeValue.Cookware,
                               [&Request](const std::string& Cookware) {
                                   return ContainsValue(Request.Cookware, Cookware);
                               });
}

bool IsRecent(const RecommendationRequest& Request, std::string_view RecipeId) {
    return ContainsValue(Request.RecentRecipeIds, RecipeId);
}

Foundation::Error InvalidRequest(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

}  // namespace

Foundation::Result<std::vector<Recommendation>> RuleBasedRecommendation::Rank(
    const RecommendationRequest& Request,
    const std::vector<Recipe>& Recipes) {
    if (Request.Servings <= 0 || Request.AvailableMinutes < 0) {
        return Foundation::Result<std::vector<Recommendation>>::FromError(
            InvalidRequest("推荐请求的份量或可用时间无效"));
    }

    std::vector<Recommendation> Recommendations;
    for (const Recipe& RecipeValue : Recipes) {
        if (RecipeValue.Status != "published" ||
            RecipeValue.PrepMinutes + RecipeValue.CookMinutes > Request.AvailableMinutes ||
            Intersects(RecipeValue.Allergens, Request.Allergies) ||
            !CookwareAvailable(RecipeValue, Request)) {
            continue;
        }
        if (!Request.Cuisines.empty() && !ContainsValue(Request.Cuisines, RecipeValue.Cuisine)) {
            continue;
        }

        const auto ValidationResult = RecipeRules::ValidateRecipe(RecipeValue);
        if (!ValidationResult.HasValue()) {
            continue;
        }

        const IngredientAvailability Availability = RecipeRules::BuildIngredientAvailability(
            RecipeValue, Request.PantryIngredientIds);
        int Score = static_cast<int>(Availability.AvailableCount * 10U);
        if (ContainsValue(Request.Cuisines, RecipeValue.Cuisine)) {
            Score += 100;
        }
        if (IsRecent(Request, RecipeValue.Id)) {
            Score -= 25;
        }

        Recommendations.push_back(
            Recommendation{RecipeValue, Availability, Score});
    }

    std::ranges::sort(Recommendations, [](const Recommendation& Left,
                                          const Recommendation& Right) {
        if (Left.Score != Right.Score) {
            return Left.Score > Right.Score;
        }
        if (Left.Availability.AvailableCount != Right.Availability.AvailableCount) {
            return Left.Availability.AvailableCount > Right.Availability.AvailableCount;
        }
        if (Left.Availability.MissingCount != Right.Availability.MissingCount) {
            return Left.Availability.MissingCount < Right.Availability.MissingCount;
        }
        if (Left.RecipeValue.Difficulty != Right.RecipeValue.Difficulty) {
            return Left.RecipeValue.Difficulty < Right.RecipeValue.Difficulty;
        }
        return Left.RecipeValue.Id < Right.RecipeValue.Id;
    });

    if (Recommendations.size() > 3U) {
        Recommendations.resize(3U);
    }
    return Recommendations;
}

}  // namespace Menu::Domain
